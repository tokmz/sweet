package logger

import (
	"time"

	"go.uber.org/zap/zapcore"
)

// Level 定义日志级别
type Level int8

const (
	// DebugLevel 调试级别
	DebugLevel Level = iota - 1
	// InfoLevel 信息级别
	InfoLevel
	// WarnLevel 警告级别
	WarnLevel
	// ErrorLevel 错误级别
	ErrorLevel
	// FatalLevel 致命错误级别
	FatalLevel
)

// Config 日志配置
type Config struct {
	// Level 日志级别
	Level Level

	// 输出目标控制
	EnableConsole bool
	EnableFile    bool

	// 日志文件配置
	Filename   string
	MaxSize    int  // 单个日志文件最大尺寸，单位MB
	MaxBackups int  // 保留的旧日志文件最大数量
	MaxAge     int  // 保留旧日志文件的最大天数
	Compress   bool // 是否压缩旧日志文件

	// 日志格式控制
	DisableTimestamp bool
	DisableCaller    bool
	DisableTrace     bool

	// 开发模式
	Development bool

	// 高级功能配置
	BufferSize     int             // 日志缓冲区大小
	FlushInterval  time.Duration   // 日志刷新间隔
	AsyncWrite     bool            // 是否启用异步写入
	SamplingConfig *SamplingConfig // 采样配置
	FilterConfig   *FilterConfig   // 过滤配置
	SecurityConfig *SecurityConfig // 安全配置
}

// SamplingConfig 采样配置
type SamplingConfig struct {
	Initial    int           // 初始采样数
	Thereafter int           // 后续采样数
	Interval   time.Duration // 采样间隔
}

// FilterConfig 过滤配置
type FilterConfig struct {
	EnableRegex    bool     // 是否启用正则过滤
	IncludePattern string   // 包含模式
	ExcludePattern string   // 排除模式
	MinLevel       Level    // 最小日志级别
	MaxLevel       Level    // 最大日志级别
	Keywords       []string // 关键词过滤
}

// SecurityConfig 安全配置
type SecurityConfig struct {
	EnableEncryption bool     // 是否启用加密
	EncryptionKey    string   // 加密密钥
	EnableMasking    bool     // 是否启用脱敏
	MaskingFields    []string // 需要脱敏的字段
}

// LogRotateConfig 日志轮转配置
type LogRotateConfig struct {
	Filename   string
	MaxSize    int
	MaxBackups int
	MaxAge     int
	Compress   bool
}

// Field 定义日志字段
type Field struct {
	Key   string
	Value any
}

// LevelToZapLevel 转换内部日志级别到zap日志级别
func (l Level) ToZapLevel() zapcore.Level {
	switch l {
	case DebugLevel:
		return zapcore.DebugLevel
	case InfoLevel:
		return zapcore.InfoLevel
	case WarnLevel:
		return zapcore.WarnLevel
	case ErrorLevel:
		return zapcore.ErrorLevel
	case FatalLevel:
		return zapcore.FatalLevel
	default:
		return zapcore.InfoLevel
	}
}

// Logger 定义Logger接口
type Logger interface {
	Debug(msg string, fields ...Field)
	Info(msg string, fields ...Field)
	Warn(msg string, fields ...Field)
	Error(msg string, fields ...Field)
	Fatal(msg string, fields ...Field)

	With(fields ...Field) Logger
	Named(name string) Logger

	// Sugar 返回一个SugaredLogger，提供更简洁的API
	Sugar() SugaredLogger

	// Sync 同步日志缓冲区
	Sync() error

	// SetLevel 动态设置日志级别
	SetLevel(level Level)

	// Export 导出日志
	Export(start, end time.Time, format string) (string, error)
}

// SugaredLogger 提供printf风格的API
type SugaredLogger interface {
	Debugf(template string, args ...any)
	Infof(template string, args ...any)
	Warnf(template string, args ...any)
	Errorf(template string, args ...any)
	Fatalf(template string, args ...any)

	Debugw(msg string, keysAndValues ...any)
	Infow(msg string, keysAndValues ...any)
	Warnw(msg string, keysAndValues ...any)
	Errorw(msg string, keysAndValues ...any)
	Fatalw(msg string, keysAndValues ...any)

	With(args ...any) SugaredLogger
	Named(name string) SugaredLogger
}

// TimeLayout 日志时间格式
const TimeLayout = time.RFC3339
