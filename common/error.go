package common

type Error struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (e *Error) Error() string {
	return e.Msg
}

// NewError 创建一个新的错误
// code 错误码，msg 错误消息
func NewError(code int, msg string) *Error {
	return &Error{
		Code: code,
		Msg:  msg,
	}
}

var (
	// ErrServer 服务器错误
	ErrServer = NewError(1001, "服务器错误")
	// ErrUnknown 未知错误
	ErrUnknown = NewError(1002, "未知错误")
	// ErrInvalidParam 无效参数错误
	ErrInvalidParam = NewError(1003, "无效参数")
	// ErrNotFound 资源未找到错误
	ErrNotFound = NewError(1004, "资源未找到")
	// ErrUnauthorized 未授权错误
	ErrUnauthorized = NewError(1005, "未授权")
	// ErrForbidden 禁止访问错误
	ErrForbidden = NewError(1006, "禁止访问")
)
