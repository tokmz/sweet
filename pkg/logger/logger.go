package logger

import (
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"sync"
	"time"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

var (
	// 默认日志对象
	defaultLogger *zapLogger
	// 保护defaultLogger的互斥锁
	loggerMu sync.RWMutex
	// 是否已初始化标志
	initialized bool
)

// zapLogger 是Logger接口的主要实现
type zapLogger struct {
	zap    *zap.Logger
	level  zap.AtomicLevel
	sugar  *zap.SugaredLogger
	config Config
}

// zapSugar 是SugaredLogger接口的实现
type zapSugar struct {
	sugar *zap.SugaredLogger
	base  *zapLogger
}

// Init 初始化默认日志
func Init(config Config) error {
	loggerMu.Lock()
	defer loggerMu.Unlock()

	// 如果已经初始化，先关闭之前的logger
	if initialized && defaultLogger != nil {
		_ = defaultLogger.Sync()
	}

	// 验证配置
	if err := validateConfig(&config); err != nil {
		return err
	}

	logger, err := NewLogger(config)
	if err != nil {
		return err
	}

	defaultLogger = logger.(*zapLogger)
	initialized = true
	return nil
}

// validateConfig 验证并设置配置的默认值
func validateConfig(config *Config) error {
	// 设置默认值
	if config.MaxSize == 0 {
		config.MaxSize = 100
	}
	if config.MaxBackups == 0 {
		config.MaxBackups = 10
	}
	if config.MaxAge == 0 {
		config.MaxAge = 30
	}

	// 验证日志文件配置
	if config.EnableFile && config.Filename != "" {
		// 确保日志目录存在
		logDir := filepath.Dir(config.Filename)
		if logDir != "." {
			if err := os.MkdirAll(logDir, 0755); err != nil {
				return fmt.Errorf("创建日志目录失败: %w", err)
			}
		}
	}

	return nil
}

// NewLogger 创建一个新的日志记录器
func NewLogger(config Config) (Logger, error) {
	// 配置已在validateConfig中验证和设置默认值

	// 创建原子级别
	level := zap.NewAtomicLevelAt(config.Level.ToZapLevel())

	// 创建编码器配置
	encoderConfig := zapcore.EncoderConfig{
		TimeKey:        "time",
		LevelKey:       "level",
		NameKey:        "logger",
		CallerKey:      "caller",
		FunctionKey:    zapcore.OmitKey,
		MessageKey:     "msg",
		StacktraceKey:  "stacktrace",
		LineEnding:     zapcore.DefaultLineEnding,
		EncodeLevel:    zapcore.CapitalLevelEncoder,
		EncodeTime:     zapcore.ISO8601TimeEncoder,
		EncodeDuration: zapcore.StringDurationEncoder,
		EncodeCaller:   zapcore.ShortCallerEncoder,
	}

	// 开发模式下使用彩色日志
	if config.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	}

	// 创建核心写入器
	var cores []zapcore.Core

	// 控制台输出
	if config.EnableConsole {
		consoleEncoder := zapcore.NewConsoleEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(
			consoleEncoder,
			zapcore.AddSync(os.Stdout),
			level,
		))
	}

	// 文件输出
	if config.EnableFile && config.Filename != "" {
		// 使用lumberjack进行日志轮转
		lumberJackLogger := &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.MaxSize,
			MaxBackups: config.MaxBackups,
			MaxAge:     config.MaxAge,
			Compress:   config.Compress,
		}

		fileEncoder := zapcore.NewJSONEncoder(encoderConfig)
		cores = append(cores, zapcore.NewCore(
			fileEncoder,
			zapcore.AddSync(lumberJackLogger),
			level,
		))
	}

	// 创建zap日志记录器
	if len(cores) == 0 {
		return nil, errors.New("没有可用的日志输出目标")
	}

	core := zapcore.NewTee(cores...)

	// 创建日志选项
	var opts []zap.Option
	if !config.DisableCaller {
		opts = append(opts, zap.AddCaller())
	}
	if !config.DisableTrace && (config.Development || level.Level() == zapcore.DebugLevel) {
		opts = append(opts, zap.AddStacktrace(zapcore.ErrorLevel))
	}
	if config.Development {
		opts = append(opts, zap.Development())
	}

	// 创建zap记录器
	logger := zap.New(core, opts...)

	// 创建包装器
	return &zapLogger{
		zap:    logger,
		level:  level,
		sugar:  logger.Sugar(),
		config: config,
	}, nil
}

// 将Field转换为zap.Field
func toZapFields(fields []Field) []zap.Field {
	zapFields := make([]zap.Field, len(fields))
	for i, field := range fields {
		zapFields[i] = zap.Any(field.Key, field.Value)
	}
	return zapFields
}

// Debug 输出调试日志
func (l *zapLogger) Debug(msg string, fields ...Field) {
	l.zap.Debug(msg, toZapFields(fields)...)
}

// Info 输出信息日志
func (l *zapLogger) Info(msg string, fields ...Field) {
	l.zap.Info(msg, toZapFields(fields)...)
}

// Warn 输出警告日志
func (l *zapLogger) Warn(msg string, fields ...Field) {
	l.zap.Warn(msg, toZapFields(fields)...)
}

// Error 输出错误日志
func (l *zapLogger) Error(msg string, fields ...Field) {
	l.zap.Error(msg, toZapFields(fields)...)
}

// Fatal 输出致命错误日志并退出程序
func (l *zapLogger) Fatal(msg string, fields ...Field) {
	l.zap.Fatal(msg, toZapFields(fields)...)
}

// With 创建带有预设字段的Logger
func (l *zapLogger) With(fields ...Field) Logger {
	return &zapLogger{
		zap:    l.zap.With(toZapFields(fields)...),
		level:  l.level,
		sugar:  l.zap.With(toZapFields(fields)...).Sugar(),
		config: l.config,
	}
}

// Named 创建具有指定名称的Logger
func (l *zapLogger) Named(name string) Logger {
	return &zapLogger{
		zap:    l.zap.Named(name),
		level:  l.level,
		sugar:  l.zap.Named(name).Sugar(),
		config: l.config,
	}
}

// Sugar 返回SugaredLogger接口
func (l *zapLogger) Sugar() SugaredLogger {
	return &zapSugar{
		sugar: l.sugar,
		base:  l,
	}
}

// Sync 同步日志缓冲区
func (l *zapLogger) Sync() error {
	return l.zap.Sync()
}

// SetLevel 动态设置日志级别
func (l *zapLogger) SetLevel(level Level) {
	l.level.SetLevel(level.ToZapLevel())
}

// Export 导出日志 - 简单实现
func (l *zapLogger) Export(start, end time.Time, format string) (string, error) {
	// 这里只是一个示例实现，实际上需要根据日志存储方式进行具体实现
	return fmt.Sprintf("日志导出功能未完全实现: %v 到 %v，格式：%s", start, end, format), nil
}

// Debugf 输出格式化调试日志
func (s *zapSugar) Debugf(template string, args ...any) {
	s.sugar.Debugf(template, args...)
}

// Infof 输出格式化信息日志
func (s *zapSugar) Infof(template string, args ...any) {
	s.sugar.Infof(template, args...)
}

// Warnf 输出格式化警告日志
func (s *zapSugar) Warnf(template string, args ...any) {
	s.sugar.Warnf(template, args...)
}

// Errorf 输出格式化错误日志
func (s *zapSugar) Errorf(template string, args ...any) {
	s.sugar.Errorf(template, args...)
}

// Fatalf 输出格式化致命错误日志并退出程序
func (s *zapSugar) Fatalf(template string, args ...any) {
	s.sugar.Fatalf(template, args...)
}

// Debugw 输出带有键值对的调试日志
func (s *zapSugar) Debugw(msg string, keysAndValues ...any) {
	s.sugar.Debugw(msg, keysAndValues...)
}

// Infow 输出带有键值对的信息日志
func (s *zapSugar) Infow(msg string, keysAndValues ...any) {
	s.sugar.Infow(msg, keysAndValues...)
}

// Warnw 输出带有键值对的警告日志
func (s *zapSugar) Warnw(msg string, keysAndValues ...any) {
	s.sugar.Warnw(msg, keysAndValues...)
}

// Errorw 输出带有键值对的错误日志
func (s *zapSugar) Errorw(msg string, keysAndValues ...any) {
	s.sugar.Errorw(msg, keysAndValues...)
}

// Fatalw 输出带有键值对的致命错误日志并退出程序
func (s *zapSugar) Fatalw(msg string, keysAndValues ...any) {
	s.sugar.Fatalw(msg, keysAndValues...)
}

// With 创建带有预设键值对的SugaredLogger
func (s *zapSugar) With(args ...any) SugaredLogger {
	return &zapSugar{
		sugar: s.sugar.With(args...),
		base:  s.base,
	}
}

// Named 创建具有指定名称的SugaredLogger
func (s *zapSugar) Named(name string) SugaredLogger {
	return &zapSugar{
		sugar: s.sugar.Named(name),
		base:  s.base,
	}
}
