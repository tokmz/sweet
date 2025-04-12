package common

type Response struct {
	Code    int    `json:"code"`
	Msg     string `json:"msg"`
	Data    any    `json:"data"`
	TraceID string `json:"trace_id"`
}

func New(code int, msg string, data any) *Response {
	return &Response{Code: code, Msg: msg, Data: data}
}

type Page struct {
	List  any   `json:"list"`
	Total int64 `json:"total"`
}

func NewPage(list any, total int64) *Page {
	return &Page{List: list, Total: total}
}
