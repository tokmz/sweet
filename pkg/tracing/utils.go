package tracing

import (
	"context"
	"fmt"
	"net"
	"os"
	"path/filepath"
	"regexp"
	"runtime"
	"time"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

// IsInitialized 检查tracing是否已初始化
func IsInitialized() bool {
	return initialized && tracer != nil
}

// ShouldRecordHTTPHeaders 检查是否应该记录HTTP头
func ShouldRecordHTTPHeaders() bool {
	return initialized && globalConfig.RecordHTTPHeaders
}

// ShouldRecordHTTPBody 检查是否应该记录HTTP请求体
func ShouldRecordHTTPBody() bool {
	return initialized && globalConfig.RecordHTTPBody
}

// ShouldRecordSQL 检查是否应该记录SQL语句
func ShouldRecordSQL() bool {
	return initialized && globalConfig.RecordSQL
}

// IsTLSEnabled 检查是否启用了TLS
func IsTLSEnabled() bool {
	return initialized && globalConfig.Security.EnableTLS
}

// GetSamplingType 获取当前采样类型
func GetSamplingType() string {
	if !initialized {
		return ""
	}
	return globalConfig.Sampling.Type
}

// IsDynamicSamplingEnabled 检查是否启用了动态采样
func IsDynamicSamplingEnabled() bool {
	return initialized && globalConfig.Sampling.Type == "dynamic"
}

// GetConfig 获取当前配置（返回副本以避免外部修改）
func GetConfig() Config {
	if !initialized {
		return DefaultConfig()
	}
	return globalConfig
}

// ContextWithSpan 创建带有span的新上下文
func ContextWithSpan(ctx context.Context, span trace.Span) context.Context {
	return trace.ContextWithSpan(ctx, span)
}

// TraceFunction 追踪函数执行，自动记录执行时间和错误
func TraceFunction(ctx context.Context, name string) (context.Context, func(err *error)) {
	if !initialized || tracer == nil {
		return ctx, func(err *error) {}
	}

	// 获取调用者信息
	pc, file, line, ok := runtime.Caller(1)
	funcName := "unknown"
	if ok {
		funcInfo := runtime.FuncForPC(pc)
		if funcInfo != nil {
			funcName = funcInfo.Name()
		}
	}

	// 创建函数调用的span
	startTime := time.Now()
	ctx, span := tracer.Start(ctx, name,
		trace.WithAttributes(
			attribute.String("function.name", funcName),
			attribute.String("function.file", file),
			attribute.Int("function.line", line),
		),
	)

	// 返回函数结束时的回调，用于记录执行时间和错误
	return ctx, func(errPtr *error) {
		duration := time.Since(startTime)

		// 记录执行时间
		span.SetAttributes(attribute.Int64("function.duration_ms", duration.Milliseconds()))

		// 处理错误
		if errPtr != nil && *errPtr != nil {
			err := *errPtr
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		span.End()
	}
}

// WrapError 将错误包装并添加当前函数信息，用于增强错误上下文
func WrapError(err error, message string) error {
	if err == nil {
		return nil
	}

	_, file, line, ok := runtime.Caller(1)
	if !ok {
		file = "unknown"
		line = 0
	}

	// 获取短文件名
	shortFile := filepath.Base(file)

	return fmt.Errorf("%s at %s:%d: %w", message, shortFile, line, err)
}

// TraceMethodCall 追踪方法调用，常用于Repository、Service等层的方法追踪
func TraceMethodCall(ctx context.Context, component, method string, args ...interface{}) (context.Context, func(result interface{}, err *error)) {
	if !initialized || tracer == nil {
		return ctx, func(result interface{}, err *error) {}
	}

	// 构建参数属性
	attrs := make([]attribute.KeyValue, 0, len(args)+2)
	attrs = append(attrs,
		attribute.String("component", component),
		attribute.String("method", method),
	)

	// 添加参数信息
	for i, arg := range args {
		// 限制参数值长度，避免过大
		argStr := fmt.Sprintf("%+v", arg)
		argStr = LimitAttributeValue(argStr, globalConfig.MaxAttributeValueLength)
		attrs = append(attrs, attribute.String(fmt.Sprintf("arg.%d", i), argStr))
	}

	// 创建span
	startTime := time.Now()
	ctx, span := tracer.Start(ctx, fmt.Sprintf("%s.%s", component, method),
		trace.WithAttributes(attrs...),
	)

	// 返回方法结束时的回调
	return ctx, func(result interface{}, errPtr *error) {
		duration := time.Since(startTime)

		// 记录执行时间
		span.SetAttributes(attribute.Int64("method.duration_ms", duration.Milliseconds()))

		// 记录结果
		if result != nil {
			resultStr := fmt.Sprintf("%+v", result)
			resultStr = LimitAttributeValue(resultStr, globalConfig.MaxAttributeValueLength)
			span.SetAttributes(attribute.String("method.result", resultStr))
		}

		// 处理错误
		if errPtr != nil && *errPtr != nil {
			err := *errPtr
			span.RecordError(err)
			span.SetStatus(codes.Error, err.Error())
		} else {
			span.SetStatus(codes.Ok, "")
		}

		span.End()
	}
}

// WithTraceID 生成带有链路追踪信息的错误消息
func WithTraceID(ctx context.Context, message string) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().HasTraceID() {
		return message
	}

	traceID := span.SpanContext().TraceID().String()
	return fmt.Sprintf("%s [trace_id=%s]", message, traceID)
}

// GetTraceID 从上下文中获取链路追踪ID
func GetTraceID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().HasTraceID() {
		return ""
	}

	return span.SpanContext().TraceID().String()
}

// GetSpanID 从上下文中获取Span ID
func GetSpanID(ctx context.Context) string {
	span := trace.SpanFromContext(ctx)
	if !span.SpanContext().HasSpanID() {
		return ""
	}

	return span.SpanContext().SpanID().String()
}

// GetLocalIP 获取本机IP地址
func GetLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return ""
	}

	for _, addr := range addrs {
		if ipnet, ok := addr.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}

	return ""
}

// GetHostInfo 获取主机信息
func GetHostInfo() map[string]string {
	info := make(map[string]string)

	// 获取主机名
	hostname, err := os.Hostname()
	if err == nil {
		info["hostname"] = hostname
	}

	// 获取IP地址
	ip := GetLocalIP()
	if ip != "" {
		info["ip"] = ip
	}

	// 获取操作系统和架构
	info["os"] = runtime.GOOS
	info["arch"] = runtime.GOARCH

	return info
}

// IsSensitivePattern 使用正则表达式检查是否是敏感信息模式
func IsSensitivePattern(value string) bool {
	// 信用卡号模式 (简化版)
	creditCardPattern := regexp.MustCompile(`\b(?:\d{4}[- ]?){3}\d{4}\b`)
	if creditCardPattern.MatchString(value) {
		return true
	}

	// 中国身份证号模式 (18位)
	idCardPattern := regexp.MustCompile(`\b\d{17}[\dXx]\b`)
	if idCardPattern.MatchString(value) {
		return true
	}

	// 手机号模式
	phonePattern := regexp.MustCompile(`\b1[3-9]\d{9}\b`)
	if phonePattern.MatchString(value) {
		return true
	}

	// 邮箱地址模式
	emailPattern := regexp.MustCompile(`\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\.[A-Za-z]{2,}\b`)
	if emailPattern.MatchString(value) {
		return true
	}

	// JWT Token模式 (简化版)
	jwtPattern := regexp.MustCompile(`\bey[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\.[A-Za-z0-9_-]+\b`)
	if jwtPattern.MatchString(value) {
		return true
	}

	return false
}

// MaskSensitiveValue 对敏感值进行掩码处理
func MaskSensitiveValue(value string, maskChar string) string {
	if maskChar == "" {
		maskChar = "*"
	}

	// 信用卡号掩码 (保留前4位和后4位)
	creditCardPattern := regexp.MustCompile(`\b(\d{4})[- ]?\d{4}[- ]?\d{4}[- ]?(\d{4})\b`)
	value = creditCardPattern.ReplaceAllString(value, "$1-XXXX-XXXX-$2")

	// 手机号掩码 (保留前3位和后4位)
	phonePattern := regexp.MustCompile(`\b(1[3-9]\d{2})\d{4}(\d{4})\b`)
	value = phonePattern.ReplaceAllString(value, "$1****$2")

	// 邮箱地址掩码 (用户名部分保留首尾字符)
	emailPattern := regexp.MustCompile(`\b([A-Za-z0-9])[A-Za-z0-9._%+-]+([A-Za-z0-9])@([A-Za-z0-9.-]+\.[A-Za-z]{2,})\b`)
	value = emailPattern.ReplaceAllString(value, "$1***$2@$3")

	// 身份证号掩码 (保留前3位和后4位)
	idCardPattern := regexp.MustCompile(`\b(\d{3})\d{11}(\d{4}|\d{3}[Xx])\b`)
	value = idCardPattern.ReplaceAllString(value, "$1***********$2")

	return value
}
