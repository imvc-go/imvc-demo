package calculator

import (
	"google.golang.org/grpc"
	"imvc-test/service/controllers/calculator"
)

func Register(server *grpc.Server) {
	calcSvc := calculator.NewCalculatorController()
	RegisterCalculatorServer(server, calcSvc)
}
