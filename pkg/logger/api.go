package logger

// 全局日志函数

// Debug 使用默认日志记录器输出调试日志
func Debug(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Debug(msg, fields...)
	}
}

// Info 使用默认日志记录器输出信息日志
func Info(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Info(msg, fields...)
	}
}

// Warn 使用默认日志记录器输出警告日志
func Warn(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Warn(msg, fields...)
	}
}

// Error 使用默认日志记录器输出错误日志
func Error(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Error(msg, fields...)
	}
}

// Fatal 使用默认日志记录器输出致命错误日志
func Fatal(msg string, fields ...Field) {
	if defaultLogger != nil {
		defaultLogger.Fatal(msg, fields...)
	}
}

// With 使用默认日志记录器创建带有预设字段的Logger
func With(fields ...Field) Logger {
	if defaultLogger != nil {
		return defaultLogger.With(fields...)
	}
	return nil
}

// Named 使用默认日志记录器创建具有指定名称的Logger
func Named(name string) Logger {
	if defaultLogger != nil {
		return defaultLogger.Named(name)
	}
	return nil
}

// Sugar 获取默认日志记录器的SugaredLogger
func Sugar() SugaredLogger {
	if defaultLogger != nil {
		return defaultLogger.Sugar()
	}
	return nil
}

// Sync 同步默认日志记录器的缓冲区
func Sync() error {
	if defaultLogger != nil {
		return defaultLogger.Sync()
	}
	return nil
}

// SetLevel 设置默认日志记录器的日志级别
func SetLevel(level Level) {
	if defaultLogger != nil {
		defaultLogger.SetLevel(level)
	}
}

// GetDefaultLogger 获取默认日志记录器
func GetDefaultLogger() Logger {
	return defaultLogger
}

// SugaredAPI 提供便捷的Sugar风格API

// Debugf 使用默认Sugar日志记录器输出格式化调试日志
func Debugf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Debugf(template, args...)
	}
}

// Infof 使用默认Sugar日志记录器输出格式化信息日志
func Infof(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Infof(template, args...)
	}
}

// Warnf 使用默认Sugar日志记录器输出格式化警告日志
func Warnf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Warnf(template, args...)
	}
}

// Errorf 使用默认Sugar日志记录器输出格式化错误日志
func Errorf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Errorf(template, args...)
	}
}

// Fatalf 使用默认Sugar日志记录器输出格式化致命错误日志
func Fatalf(template string, args ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Fatalf(template, args...)
	}
}

// Debugw 使用默认Sugar日志记录器输出带有键值对的调试日志
func Debugw(msg string, keysAndValues ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Debugw(msg, keysAndValues...)
	}
}

// Infow 使用默认Sugar日志记录器输出带有键值对的信息日志
func Infow(msg string, keysAndValues ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Infow(msg, keysAndValues...)
	}
}

// Warnw 使用默认Sugar日志记录器输出带有键值对的警告日志
func Warnw(msg string, keysAndValues ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Warnw(msg, keysAndValues...)
	}
}

// Errorw 使用默认Sugar日志记录器输出带有键值对的错误日志
func Errorw(msg string, keysAndValues ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Errorw(msg, keysAndValues...)
	}
}

// Fatalw 使用默认Sugar日志记录器输出带有键值对的致命错误日志
func Fatalw(msg string, keysAndValues ...any) {
	if defaultLogger != nil {
		defaultLogger.Sugar().Fatalw(msg, keysAndValues...)
	}
}
