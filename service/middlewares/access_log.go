package middlewares

import (
	"fmt"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"github.com/kataras/iris/v12/middleware/logger"
	"strconv"
	"time"
)

type AccessLogConfig struct {
	loggerCfg logger.Config
}

type accessLogMiddleware struct {
	config AccessLogConfig
}

// New creates and returns a new request logger middleware.
// Do not confuse it with the framework's Logger.
// This is for the http requests.
//
// Receives an optional configuation.
func NewAccessLogHandler(cfg ...AccessLogConfig) context.Handler {
	c := DefaultAccessLogConfig()
	if len(cfg) > 0 {
		c = cfg[0]
	}
	l := &accessLogMiddleware{config: c}

	return l.ServeHTTP
}

// DefaultConfig returns a default config
// that have all boolean fields to true except `Columns`,
// all strings are empty,
// LogFunc and Skippers to nil as well.
func DefaultAccessLogConfig() AccessLogConfig {
	return AccessLogConfig{
		loggerCfg: logger.Config{
			Status:             true,
			IP:                 true,
			Method:             true,
			Path:               true,
			Query:              false,
			LogFunc:            nil,
			LogFuncCtx:         nil,
			Skippers:           nil,
			MessageContextKeys: []string{"logger_message"},
			//如果不为空然后它的内容来自`ctx.GetHeader（“User-Agent”）
			MessageHeaderKeys: []string{"User-Agent"},
		},
	}
}

// Serve serves the middleware
func (mid *accessLogMiddleware) ServeHTTP(ctx iris.Context) {
	// all except latency to string
	var status, ip, method, path string
	var latency time.Duration
	var startTime, endTime time.Time
	startTime = time.Now()

	ctx.Next()

	// no time.Since in order to format it well after
	endTime = time.Now()
	latency = endTime.Sub(startTime)

	if mid.config.loggerCfg.Status {
		status = strconv.Itoa(ctx.GetStatusCode())
	}

	if mid.config.loggerCfg.IP {
		ip = ctx.RemoteAddr()
	}

	if mid.config.loggerCfg.Method {
		method = ctx.Method()
	}

	if mid.config.loggerCfg.Path {
		if mid.config.loggerCfg.Query {
			path = ctx.Request().URL.RequestURI()
		} else {
			path = ctx.Path()
		}
	}

	var message interface{}
	if ctxKeys := mid.config.loggerCfg.MessageContextKeys; len(ctxKeys) > 0 {
		for _, key := range ctxKeys {
			msg := ctx.Values().Get(key)
			if message == nil {
				message = msg
			} else {
				message = fmt.Sprintf(" %v %v", message, msg)
			}
		}
	}
	var headerMessage interface{}
	if headerKeys := mid.config.loggerCfg.MessageHeaderKeys; len(headerKeys) > 0 {
		for _, key := range headerKeys {
			msg := ctx.GetHeader(key)
			if headerMessage == nil {
				headerMessage = msg
			} else {
				headerMessage = fmt.Sprintf(" %v %v", headerMessage, msg)
			}
		}
	}

	// print the logs
	if logFunc := mid.config.loggerCfg.LogFunc; logFunc != nil {
		logFunc(endTime, latency, status, ip, method, path, message, headerMessage)
		return
	} else if logFuncCtx := mid.config.loggerCfg.LogFuncCtx; logFuncCtx != nil {
		logFuncCtx(ctx, latency)
		return
	}

	// no new line, the framework's logger is responsible how to render each log.
	line := fmt.Sprintf("%v %4v %s %s %s", status, latency, ip, method, path)
	if message != nil {
		line += fmt.Sprintf(" %v", message)
	}

	if headerMessage != nil {
		line += fmt.Sprintf(" %v", headerMessage)
	}
	ctx.Application().Logger().Info(line)
}
