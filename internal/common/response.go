package common

type Response struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id,omitempty"`
}

func NewResponse(code int, msg string, data any, traceID string) *Response {
	return &Response{
		Code:    code,
		Msg:     msg,
		Data:    data,
		TraceID: traceID,
	}
}

type ListResp struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

func NewListResp(list any, total int64) *ListResp {
	return &ListResp{
		List:  list,
		Total: total,
	}
}
