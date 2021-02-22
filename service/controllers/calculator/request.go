package calculator

type CommonReq struct {
	A int
	B int
}

type ReqAdd struct {
	CommonReq
}

type ReqPlus struct {
	CommonReq
}

type ReqMulti struct {
	CommonReq
}

type ReqDivide struct {
	CommonReq
	IgnoreError bool
}

type Response struct {
	Result int
}
