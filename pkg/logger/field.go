package logger

import (
	"time"
)

// 通用字段创建函数

// String 创建字符串字段
func String(key string, val string) Field {
	return Field{Key: key, Value: val}
}

// Int 创建整数字段
func Int(key string, val int) Field {
	return Field{Key: key, Value: val}
}

// Int64 创建64位整数字段
func Int64(key string, val int64) Field {
	return Field{Key: key, Value: val}
}

// Float64 创建64位浮点数字段
func Float64(key string, val float64) Field {
	return Field{Key: key, Value: val}
}

// Bool 创建布尔字段
func Bool(key string, val bool) Field {
	return Field{Key: key, Value: val}
}

// Time 创建时间字段
func Time(key string, val time.Time) Field {
	return Field{Key: key, Value: val}
}

// Duration 创建时间间隔字段
func Duration(key string, val time.Duration) Field {
	return Field{Key: key, Value: val}
}

// Any 创建任意类型字段
func Any(key string, val interface{}) Field {
	return Field{Key: key, Value: val}
}

// Err 创建错误字段
func Err(err error) Field {
	return Field{Key: "error", Value: err}
}

// Stack 创建堆栈跟踪字段
func Stack() Field {
	return Field{Key: "stack", Value: "stack trace"}
}

// Caller 创建调用者信息字段
func Caller() Field {
	return Field{Key: "caller", Value: "caller info"}
}

// Namespace 创建命名空间字段
func Namespace(name string) Field {
	return Field{Key: "namespace", Value: name}
}

// 业务常用字段

// TraceID 创建跟踪ID字段
func TraceID(id string) Field {
	return Field{Key: "trace_id", Value: id}
}

// SpanID 创建跨度ID字段
func SpanID(id string) Field {
	return Field{Key: "span_id", Value: id}
}

// UserID 创建用户ID字段
func UserID(id interface{}) Field {
	return Field{Key: "user_id", Value: id}
}

// RequestID 创建请求ID字段
func RequestID(id string) Field {
	return Field{Key: "request_id", Value: id}
}

// IP 创建IP地址字段
func IP(ip string) Field {
	return Field{Key: "ip", Value: ip}
}

// Method 创建HTTP方法字段
func Method(method string) Field {
	return Field{Key: "method", Value: method}
}

// URL 创建URL字段
func URL(url string) Field {
	return Field{Key: "url", Value: url}
}

// StatusCode 创建HTTP状态码字段
func StatusCode(code int) Field {
	return Field{Key: "status_code", Value: code}
}

// Latency 创建延迟时间字段
func Latency(duration time.Duration) Field {
	return Field{Key: "latency", Value: duration}
}

// Module 创建模块名称字段
func Module(name string) Field {
	return Field{Key: "module", Value: name}
}

// Function 创建函数名称字段
func Function(name string) Field {
	return Field{Key: "function", Value: name}
}
