package logger

import (
	"time"
)

// LogLevel 日志级别
type LogLevel string

const (
	DebugLevel LogLevel = "debug"
	InfoLevel  LogLevel = "info"
	WarnLevel  LogLevel = "warn"
	ErrorLevel LogLevel = "error"
	PanicLevel LogLevel = "panic"
	FatalLevel LogLevel = "fatal"
)

// OutputMode 输出模式
type OutputMode string

const (
	ConsoleMode OutputMode = "console"   // 控制台输出
	FileMode    OutputMode = "file"      // 文件输出
	BothMode    OutputMode = "both"      // 同时输出到控制台和文件
)

// Format 日志格式
type Format string

const (
	JSONFormat Format = "json" // JSON格式
	TextFormat Format = "text" // 文本格式
)

// RotateConfig 日志分割配置
type RotateConfig struct {
	// MaxSize 单个日志文件最大大小（MB）
	MaxSize int `json:"max_size" yaml:"max_size"`
	// MaxAge 日志文件保留天数
	MaxAge int `json:"max_age" yaml:"max_age"`
	// MaxBackups 保留的日志文件数量
	MaxBackups int `json:"max_backups" yaml:"max_backups"`
	// Compress 是否压缩历史日志文件
	Compress bool `json:"compress" yaml:"compress"`
	// LocalTime 是否使用本地时间
	LocalTime bool `json:"local_time" yaml:"local_time"`
	// RotateTime 按时间分割的间隔
	RotateTime time.Duration `json:"rotate_time" yaml:"rotate_time"`
}

// Config 日志配置
type Config struct {
	// Level 日志级别
	Level LogLevel `json:"level" yaml:"level"`
	// Format 日志格式
	Format Format `json:"format" yaml:"format"`
	// OutputMode 输出模式
	OutputMode OutputMode `json:"output_mode" yaml:"output_mode"`
	// Filename 日志文件路径
	Filename string `json:"filename" yaml:"filename"`
	// EnableCaller 是否显示调用者信息
	EnableCaller bool `json:"enable_caller" yaml:"enable_caller"`
	// EnableStacktrace 是否显示堆栈信息
	EnableStacktrace bool `json:"enable_stacktrace" yaml:"enable_stacktrace"`
	// Development 是否为开发模式
	Development bool `json:"development" yaml:"development"`
	// Async 是否启用异步写入
	Async bool `json:"async" yaml:"async"`
	// BufferSize 异步缓冲区大小
	BufferSize int `json:"buffer_size" yaml:"buffer_size"`
	// FlushInterval 刷新间隔
	FlushInterval time.Duration `json:"flush_interval" yaml:"flush_interval"`
	// Rotate 日志分割配置
	Rotate RotateConfig `json:"rotate" yaml:"rotate"`
	// ServiceName 服务名称
	ServiceName string `json:"service_name" yaml:"service_name"`
	// ServiceVersion 服务版本
	ServiceVersion string `json:"service_version" yaml:"service_version"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		Level:            InfoLevel,
		Format:           JSONFormat,
		OutputMode:       ConsoleMode,
		Filename:         "logs/app.log",
		EnableCaller:     true,
		EnableStacktrace: false,
		Development:      false,
		Async:            true,
		BufferSize:       256 * 1024, // 256KB
		FlushInterval:    time.Second,
		Rotate: RotateConfig{
			MaxSize:    100, // 100MB
			MaxAge:     30,  // 30天
			MaxBackups: 10,  // 保留10个文件
			Compress:   true,
			LocalTime:  true,
			RotateTime: 24 * time.Hour, // 每天分割
		},

	}
}

// DevelopmentConfig 返回开发环境配置
func DevelopmentConfig() *Config {
	config := DefaultConfig()
	config.Level = DebugLevel
	config.Format = TextFormat
	config.OutputMode = ConsoleMode
	config.Development = true
	config.EnableStacktrace = true
	config.Async = false
	return config
}

// ProductionConfig 返回生产环境配置
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.Level = InfoLevel
	config.Format = JSONFormat
	config.OutputMode = FileMode
	config.Development = false
	config.EnableStacktrace = false
	config.Async = true

	return config
}