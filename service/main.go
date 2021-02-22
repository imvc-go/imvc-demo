package main

import (
	"fmt"
	"github.com/kataras/golog"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/mvc"
	"imvc-test/service/libs"
	"imvc-test/service/middlewares"
	"imvc-test/service/utils"
	"imvc-test/service/views/http"
)

const (
	configKeyServerHost = "app.server.host"
	configKeyServerPort = "app.server.port"
)

func main() {
	app := iris.New()

	//utils init function
	utils.Init()

	//libs init function
	if err := libs.Init(); err != nil {
		libs.Close()
		golog.Fatalf("Init libs occured an error: %v", err)
	}
	iris.RegisterOnInterrupt(func() {
		defer libs.Close()
	})

	//tick service initial for prometheus pull
	app.Get("/api/v1/metrics", middlewares.TickMetricHandler())

	//business middleware configure
	mvc.Configure(app, middlewareConfigure)

	//here I add a calculator.proto service
	mvc.Configure(app.Party("/api/v1/calculator"), http.RegisterCalculatorService)

	mvc.New(app.Party("/api/v1/calculator", mvc.GRPC{
		Server:      app,
		ServiceName: configKeyServerHost,
		Strict:      false,
	}))

	_ = app.Run(iris.Addr(getServerAddr()), iris.WithRemoteAddrHeader("X-Forwarded-For"))
}

// 中间件的配置
// 注：先后顺序不得随意调整
func middlewareConfigure(app *mvc.Application) {
	// Tick
	app.Router.Use(middlewares.TickProm(middlewares.TickPromConfig{
		ServiceName: "service_name",
	}))

	// 把每个请求以info的优先级打印出来
	app.Router.Use(middlewares.NewAccessLogHandler())
}

func getServerAddr() string {
	return fmt.Sprintf(
		"%s:%d",
		utils.GetConfig().GetString(configKeyServerHost), utils.GetConfig().GetInt(configKeyServerPort))
}
