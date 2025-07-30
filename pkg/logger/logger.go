package logger

import (
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"sync"

	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
	"gopkg.in/natefinch/lumberjack.v2"
)

// Logger 日志记录器接口
type Logger interface {
	Debug(msg string, fields ...zap.Field)
	Info(msg string, fields ...zap.Field)
	Warn(msg string, fields ...zap.Field)
	Error(msg string, fields ...zap.Field)
	Panic(msg string, fields ...zap.Field)
	Fatal(msg string, fields ...zap.Field)

	Debugf(template string, args ...interface{})
	Infof(template string, args ...interface{})
	Warnf(template string, args ...interface{})
	Errorf(template string, args ...interface{})
	Panicf(template string, args ...interface{})
	Fatalf(template string, args ...interface{})

	Debugw(msg string, keysAndValues ...interface{})
	Infow(msg string, keysAndValues ...interface{})
	Warnw(msg string, keysAndValues ...interface{})
	Errorw(msg string, keysAndValues ...interface{})
	Panicw(msg string, keysAndValues ...interface{})
	Fatalw(msg string, keysAndValues ...interface{})

	WithContext(ctx context.Context) Logger
	WithFields(fields ...zap.Field) Logger
	WithRequestID(requestID string) Logger
	WithUserID(userID string) Logger


	Sync() error
	Close() error
}

// ZapLogger zap日志记录器实现
type ZapLogger struct {
	logger *zap.Logger
	sugar  *zap.SugaredLogger
	config *Config
	mu     sync.RWMutex
	closed bool
}

var (
	// globalLogger 全局日志记录器
	globalLogger Logger
	// once 确保只初始化一次
	once sync.Once
)

// NewLogger 创建新的日志记录器
func NewLogger(config *Config) (Logger, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 创建编码器配置
	encoderConfig := getEncoderConfig(config)

	// 创建核心
	core, err := createCore(config, encoderConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to create logger core: %w", err)
	}

	// 创建选项
	options := getLoggerOptions(config)

	// 创建zap logger
	zapLogger := zap.New(core, options...)

	// 创建sugar logger
	sugar := zapLogger.Sugar()

	return &ZapLogger{
		logger: zapLogger,
		sugar:  sugar,
		config: config,
	}, nil
}

// getEncoderConfig 获取编码器配置
func getEncoderConfig(config *Config) zapcore.EncoderConfig {
	encoderConfig := zap.NewProductionEncoderConfig()

	if config.Development {
		encoderConfig = zap.NewDevelopmentEncoderConfig()
	}

	// 时间格式
	encoderConfig.TimeKey = "timestamp"
	encoderConfig.EncodeTime = zapcore.TimeEncoderOfLayout("2006-01-02 15:04:05.000")

	// 级别格式
	encoderConfig.LevelKey = "level"
	if config.Development {
		encoderConfig.EncodeLevel = zapcore.CapitalColorLevelEncoder
	} else {
		encoderConfig.EncodeLevel = zapcore.LowercaseLevelEncoder
	}

	// 调用者格式
	if config.EnableCaller {
		encoderConfig.CallerKey = "caller"
		encoderConfig.EncodeCaller = zapcore.ShortCallerEncoder
	}

	// 消息格式
	encoderConfig.MessageKey = "message"

	// 堆栈格式
	if config.EnableStacktrace {
		encoderConfig.StacktraceKey = "stacktrace"
	}

	return encoderConfig
}

// createCore 创建日志核心
func createCore(config *Config, encoderConfig zapcore.EncoderConfig) (zapcore.Core, error) {
	// 创建编码器
	var encoder zapcore.Encoder
	switch config.Format {
	case JSONFormat:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	case TextFormat:
		encoder = zapcore.NewConsoleEncoder(encoderConfig)
	default:
		encoder = zapcore.NewJSONEncoder(encoderConfig)
	}

	// 创建写入器
	writeSyncer, err := createWriteSyncer(config)
	if err != nil {
		return nil, err
	}

	// 创建级别
	level := getZapLevel(config.Level)

	// 创建核心
	core := zapcore.NewCore(encoder, writeSyncer, level)

	// 如果启用异步，包装为异步核心
	if config.Async {
		core = newAsyncCore(core, config.BufferSize, config.FlushInterval)
	}

	return core, nil
}

// createWriteSyncer 创建写入器
func createWriteSyncer(config *Config) (zapcore.WriteSyncer, error) {
	var writers []io.Writer

	// 控制台输出
	if config.OutputMode == ConsoleMode || config.OutputMode == BothMode {
		writers = append(writers, os.Stdout)
	}

	// 文件输出
	if config.OutputMode == FileMode || config.OutputMode == BothMode {
		// 确保目录存在
		dir := filepath.Dir(config.Filename)
		if err := os.MkdirAll(dir, 0755); err != nil {
			return nil, fmt.Errorf("failed to create log directory: %w", err)
		}

		// 创建lumberjack writer
		lumberjackLogger := &lumberjack.Logger{
			Filename:   config.Filename,
			MaxSize:    config.Rotate.MaxSize,
			MaxAge:     config.Rotate.MaxAge,
			MaxBackups: config.Rotate.MaxBackups,
			Compress:   config.Rotate.Compress,
			LocalTime:  config.Rotate.LocalTime,
		}

		writers = append(writers, lumberjackLogger)
	}

	if len(writers) == 0 {
		return nil, fmt.Errorf("no output writers configured")
	}

	// 如果只有一个writer，直接返回
	if len(writers) == 1 {
		return zapcore.AddSync(writers[0]), nil
	}

	// 多个writer，使用MultiWriteSyncer
	return zapcore.NewMultiWriteSyncer(
		zapcore.AddSync(io.MultiWriter(writers...)),
	), nil
}

// getLoggerOptions 获取logger选项
func getLoggerOptions(config *Config) []zap.Option {
	var options []zap.Option

	// 调用者信息
	if config.EnableCaller {
		options = append(options, zap.AddCaller())
		// 跳过2层调用栈：global.go函数 -> ZapLogger方法
		options = append(options, zap.AddCallerSkip(2))
	}

	// 堆栈信息
	if config.EnableStacktrace {
		options = append(options, zap.AddStacktrace(zapcore.ErrorLevel))
	}

	// 开发模式
	if config.Development {
		options = append(options, zap.Development())
	}

	// 添加服务信息字段
	if config.ServiceName != "" {
		options = append(options, zap.Fields(
			zap.String("service", config.ServiceName),
			zap.String("version", config.ServiceVersion),
		))
	}

	return options
}

// getZapLevel 转换日志级别
func getZapLevel(level LogLevel) zapcore.Level {
	switch level {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case PanicLevel:
		return zapcore.PanicLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// InitGlobalLogger 初始化全局日志记录器
func InitGlobalLogger(config *Config) error {
	var err error
	once.Do(func() {
		globalLogger, err = NewLogger(config)
	})
	return err
}

// GetGlobalLogger 获取全局日志记录器
func GetGlobalLogger() Logger {
	if globalLogger == nil {
		// 如果全局logger未初始化，使用默认配置初始化
		_ = InitGlobalLogger(DefaultConfig())
	}
	return globalLogger
}

// SetGlobalLogger 设置全局日志记录器
func SetGlobalLogger(logger Logger) {
	globalLogger = logger
}

// 实现Logger接口的方法

// Debug 记录调试级别日志
func (l *ZapLogger) Debug(msg string, fields ...zap.Field) {
	l.logger.Debug(msg, fields...)
}

// Info 记录信息级别日志
func (l *ZapLogger) Info(msg string, fields ...zap.Field) {
	l.logger.Info(msg, fields...)
}

// Warn 记录警告级别日志
func (l *ZapLogger) Warn(msg string, fields ...zap.Field) {
	l.logger.Warn(msg, fields...)
}

// Error 记录错误级别日志
func (l *ZapLogger) Error(msg string, fields ...zap.Field) {
	l.logger.Error(msg, fields...)
}

// Panic 记录恐慌级别日志
func (l *ZapLogger) Panic(msg string, fields ...zap.Field) {
	l.logger.Panic(msg, fields...)
}

// Fatal 记录致命级别日志
func (l *ZapLogger) Fatal(msg string, fields ...zap.Field) {
	l.logger.Fatal(msg, fields...)
}

// Debugf 格式化记录调试级别日志
func (l *ZapLogger) Debugf(template string, args ...interface{}) {
	l.sugar.Debugf(template, args...)
}

// Infof 格式化记录信息级别日志
func (l *ZapLogger) Infof(template string, args ...interface{}) {
	l.sugar.Infof(template, args...)
}

// Warnf 格式化记录警告级别日志
func (l *ZapLogger) Warnf(template string, args ...interface{}) {
	l.sugar.Warnf(template, args...)
}

// Errorf 格式化记录错误级别日志
func (l *ZapLogger) Errorf(template string, args ...interface{}) {
	l.sugar.Errorf(template, args...)
}

// Panicf 格式化记录恐慌级别日志
func (l *ZapLogger) Panicf(template string, args ...interface{}) {
	l.sugar.Panicf(template, args...)
}

// Fatalf 格式化记录致命级别日志
func (l *ZapLogger) Fatalf(template string, args ...interface{}) {
	l.sugar.Fatalf(template, args...)
}

// Debugw 键值对记录调试级别日志
func (l *ZapLogger) Debugw(msg string, keysAndValues ...interface{}) {
	l.sugar.Debugw(msg, keysAndValues...)
}

// Infow 键值对记录信息级别日志
func (l *ZapLogger) Infow(msg string, keysAndValues ...interface{}) {
	l.sugar.Infow(msg, keysAndValues...)
}

// Warnw 键值对记录警告级别日志
func (l *ZapLogger) Warnw(msg string, keysAndValues ...interface{}) {
	l.sugar.Warnw(msg, keysAndValues...)
}

// Errorw 键值对记录错误级别日志
func (l *ZapLogger) Errorw(msg string, keysAndValues ...interface{}) {
	l.sugar.Errorw(msg, keysAndValues...)
}

// Panicw 键值对记录恐慌级别日志
func (l *ZapLogger) Panicw(msg string, keysAndValues ...interface{}) {
	l.sugar.Panicw(msg, keysAndValues...)
}

// Fatalw 键值对记录致命级别日志
func (l *ZapLogger) Fatalw(msg string, keysAndValues ...interface{}) {
	l.sugar.Fatalw(msg, keysAndValues...)
}

// WithContext 添加上下文信息
func (l *ZapLogger) WithContext(ctx context.Context) Logger {
	fields := extractFieldsFromContext(ctx)
	return &ZapLogger{
		logger: l.logger.With(fields...),
		sugar:  l.logger.With(fields...).Sugar(),
		config: l.config,
	}
}

// WithFields 添加字段
func (l *ZapLogger) WithFields(fields ...zap.Field) Logger {
	return &ZapLogger{
		logger: l.logger.With(fields...),
		sugar:  l.logger.With(fields...).Sugar(),
		config: l.config,
	}
}

// WithRequestID 添加请求ID
func (l *ZapLogger) WithRequestID(requestID string) Logger {
	return l.WithFields(zap.String("request_id", requestID))
}

// WithUserID 添加用户ID
func (l *ZapLogger) WithUserID(userID string) Logger {
	return l.WithFields(zap.String("user_id", userID))
}



// Sync 同步日志
func (l *ZapLogger) Sync() error {
	l.mu.RLock()
	defer l.mu.RUnlock()

	if l.closed {
		return nil
	}

	return l.logger.Sync()
}

// Close 关闭日志记录器
func (l *ZapLogger) Close() error {
	l.mu.Lock()
	defer l.mu.Unlock()

	if l.closed {
		return nil
	}

	l.closed = true
	return l.logger.Sync()
}

// extractFieldsFromContext 从上下文中提取字段
func extractFieldsFromContext(ctx context.Context) []zap.Field {
	var fields []zap.Field

	// 提取请求ID（如果存在）
	if requestID := ctx.Value("request_id"); requestID != nil {
		if id, ok := requestID.(string); ok {
			fields = append(fields, zap.String("request_id", id))
		}
	}

	// 提取用户ID（如果存在）
	if userID := ctx.Value("user_id"); userID != nil {
		if id, ok := userID.(string); ok {
			fields = append(fields, zap.String("user_id", id))
		}
	}

	return fields
}