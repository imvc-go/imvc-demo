package calculator

type Controller interface {
	Add(req *ReqAdd) *Response
	Plus(req *ReqPlus) *Response
	Multi(req *ReqMulti) *Response
	Divide(req *ReqDivide) *Response
}

type controller struct{}

func NewCalculatorController() Controller {
	return &controller{}
}

func (c *controller) Add(req *ReqAdd) *Response {
	return nil
}

func (c *controller) Plus(req *ReqPlus) *Response {
	return nil
}

func (c *controller) Multi(req *ReqMulti) *Response {
	return nil
}

func (c *controller) Divide(req *ReqDivide) *Response {
	return nil
}
