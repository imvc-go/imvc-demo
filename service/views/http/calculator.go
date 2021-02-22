package http

import (
	"github.com/kataras/iris/v12/mvc"
	"imvc-test/service/controllers/calculator"
)

type CalculatorView struct {
	Controller calculator.Controller
}

func RegisterCalculatorService(app *mvc.Application) {
	controller := calculator.NewCalculatorController()
	app.Register(controller)
	app.Handle(new(CalculatorView))
}
