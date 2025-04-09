package errs

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Error) Error() string {
	return e.Msg
}

func New(code int, msg string) *Error {
	return &Error{Code: code, Msg: msg}
}

var (
	ErrServer        = New(1000, "服务器错误")
	ErrAuthorization = New(1001, "授权失败")
	ErrPermission    = New(1002, "权限不足")
	ErrNotFound      = New(1003, "资源不存在")
	ErrBadRequest    = New(1004, "请求错误")
	ErrConflict      = New(1005, "资源冲突")
	ErrTimeout       = New(1007, "请求超时")
	ErrTooMany       = New(1008, "请求过多")
	ErrUnavailable   = New(1009, "服务不可用")
	ErrUnknown       = New(1010, "未知错误")
	ErrInvalid       = New(1013, "无效参数")
)
