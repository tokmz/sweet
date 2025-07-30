package logger

import (
	"context"

	"go.uber.org/zap"
)

// 全局便捷函数，使用全局logger实例

// Debug 记录调试级别日志
func Debug(msg string, fields ...zap.Field) {
	GetGlobalLogger().Debug(msg, fields...)
}

// Info 记录信息级别日志
func Info(msg string, fields ...zap.Field) {
	GetGlobalLogger().Info(msg, fields...)
}

// Warn 记录警告级别日志
func Warn(msg string, fields ...zap.Field) {
	GetGlobalLogger().Warn(msg, fields...)
}

// Error 记录错误级别日志
func Error(msg string, fields ...zap.Field) {
	GetGlobalLogger().Error(msg, fields...)
}

// Panic 记录恐慌级别日志
func Panic(msg string, fields ...zap.Field) {
	GetGlobalLogger().Panic(msg, fields...)
}

// Fatal 记录致命级别日志
func Fatal(msg string, fields ...zap.Field) {
	GetGlobalLogger().Fatal(msg, fields...)
}

// Debugf 格式化记录调试级别日志
func Debugf(template string, args ...interface{}) {
	GetGlobalLogger().Debugf(template, args...)
}

// Infof 格式化记录信息级别日志
func Infof(template string, args ...interface{}) {
	GetGlobalLogger().Infof(template, args...)
}

// Warnf 格式化记录警告级别日志
func Warnf(template string, args ...interface{}) {
	GetGlobalLogger().Warnf(template, args...)
}

// Errorf 格式化记录错误级别日志
func Errorf(template string, args ...interface{}) {
	GetGlobalLogger().Errorf(template, args...)
}

// Panicf 格式化记录恐慌级别日志
func Panicf(template string, args ...interface{}) {
	GetGlobalLogger().Panicf(template, args...)
}

// Fatalf 格式化记录致命级别日志
func Fatalf(template string, args ...interface{}) {
	GetGlobalLogger().Fatalf(template, args...)
}

// Debugw 键值对记录调试级别日志
func Debugw(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Debugw(msg, keysAndValues...)
}

// Infow 键值对记录信息级别日志
func Infow(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Infow(msg, keysAndValues...)
}

// Warnw 键值对记录警告级别日志
func Warnw(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Warnw(msg, keysAndValues...)
}

// Errorw 键值对记录错误级别日志
func Errorw(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Errorw(msg, keysAndValues...)
}

// Panicw 键值对记录恐慌级别日志
func Panicw(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Panicw(msg, keysAndValues...)
}

// Fatalw 键值对记录致命级别日志
func Fatalw(msg string, keysAndValues ...interface{}) {
	GetGlobalLogger().Fatalw(msg, keysAndValues...)
}

// WithContext 添加上下文信息
func WithContext(ctx context.Context) Logger {
	return GetGlobalLogger().WithContext(ctx)
}

// WithFields 添加字段
func WithFields(fields ...zap.Field) Logger {
	return GetGlobalLogger().WithFields(fields...)
}

// WithRequestID 添加请求ID
func WithRequestID(requestID string) Logger {
	return GetGlobalLogger().WithRequestID(requestID)
}

// WithUserID 添加用户ID
func WithUserID(userID string) Logger {
	return GetGlobalLogger().WithUserID(userID)
}



// Sync 同步日志
func Sync() error {
	return GetGlobalLogger().Sync()
}

// 上下文相关的便捷函数

// DebugCtx 使用上下文记录调试日志
func DebugCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Debug(msg, fields...)
}

// InfoCtx 使用上下文记录信息日志
func InfoCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Info(msg, fields...)
}

// WarnCtx 使用上下文记录警告日志
func WarnCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Warn(msg, fields...)
}

// ErrorCtx 使用上下文记录错误日志
func ErrorCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Error(msg, fields...)
}

// PanicCtx 使用上下文记录恐慌日志
func PanicCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Panic(msg, fields...)
}

// FatalCtx 使用上下文记录致命日志
func FatalCtx(ctx context.Context, msg string, fields ...zap.Field) {
	WithContext(ctx).Fatal(msg, fields...)
}

// DebugfCtx 使用上下文格式化记录调试日志
func DebugfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Debugf(template, args...)
}

// InfofCtx 使用上下文格式化记录信息日志
func InfofCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Infof(template, args...)
}

// WarnfCtx 使用上下文格式化记录警告日志
func WarnfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Warnf(template, args...)
}

// ErrorfCtx 使用上下文格式化记录错误日志
func ErrorfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Errorf(template, args...)
}

// PanicfCtx 使用上下文格式化记录恐慌日志
func PanicfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Panicf(template, args...)
}

// FatalfCtx 使用上下文格式化记录致命日志
func FatalfCtx(ctx context.Context, template string, args ...interface{}) {
	WithContext(ctx).Fatalf(template, args...)
}

// DebugwCtx 使用上下文键值对记录调试日志
func DebugwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	WithContext(ctx).Debugw(msg, keysAndValues...)
}

// InfowCtx 使用上下文键值对记录信息日志
func InfowCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	WithContext(ctx).Infow(msg, keysAndValues...)
}

// WarnwCtx 使用上下文键值对记录警告日志
func WarnwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	WithContext(ctx).Warnw(msg, keysAndValues...)
}

// ErrorwCtx 使用上下文键值对记录错误日志
func ErrorwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	WithContext(ctx).Errorw(msg, keysAndValues...)
}

// PanicwCtx 使用上下文键值对记录恐慌日志
func PanicwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	WithContext(ctx).Panicw(msg, keysAndValues...)
}

// FatalwCtx 使用上下文键值对记录致命日志
func FatalwCtx(ctx context.Context, msg string, keysAndValues ...interface{}) {
	WithContext(ctx).Fatalw(msg, keysAndValues...)
}