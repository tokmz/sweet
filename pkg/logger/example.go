package logger

import (
	"context"
	"errors"
	"time"

	"go.uber.org/zap"
)

// 定义自定义的context key类型，避免与其他包的key冲突
type contextKey string

const (
	requestIDKey contextKey = "request_id"
	userIDKey    contextKey = "user_id"
)

// User 示例用户结构
type User struct {
	ID   int64  `json:"id"`
	Name string `json:"name"`
	Age  int    `json:"age"`
}

// ExampleBasicUsage 基本使用示例
func ExampleBasicUsage() {
	// 使用默认配置初始化全局logger
	config := DefaultConfig()
	err := InitGlobalLogger(config)
	if err != nil {
		panic(err)
	}

	// 基本日志记录
	Info("应用启动", zap.String("version", "1.0.0"))
	Debug("调试信息", zap.Int("port", 8080))
	Warn("警告信息", zap.String("reason", "配置文件未找到，使用默认配置"))
	Error("错误信息", zap.Error(errors.New("数据库连接失败")))

	// 格式化日志
	Infof("用户 %s 登录成功", "张三")
	Debugf("处理请求耗时: %dms", 150)

	// 键值对日志
	Infow("用户操作",
		"action", "login",
		"user_id", 12345,
		"ip", "192.168.1.100",
		"timestamp", time.Now(),
	)

	// 同步日志
	_ = Sync()
}

// ExampleDevelopmentConfig 开发环境配置示例
func ExampleDevelopmentConfig() {
	// 开发环境配置
	config := DevelopmentConfig()
	logger, err := NewLogger(config)
	if err != nil {
		panic(err)
	}

	// 设置为全局logger
	SetGlobalLogger(logger)

	// 开发环境下会显示彩色输出和堆栈信息
	Debug("开发环境调试信息")
	Info("应用启动", zap.String("env", "development"))
	Warn("这是一个警告")

	// 错误会显示堆栈信息
	Error("这是一个错误", zap.Error(errors.New("示例错误")))
}

// ExampleProductionConfig 生产环境配置示例
func ExampleProductionConfig() {
	// 生产环境配置
	config := ProductionConfig()
	config.Filename = "logs/app.log"

	logger, err := NewLogger(config)
	if err != nil {
		panic(err)
	}

	SetGlobalLogger(logger)

	// 生产环境下输出JSON格式到文件
	Info("服务启动", zap.Int("port", 8080))

	Infow("数据库连接成功",
		"driver", "mysql",
		"host", "localhost",
		"database", "userdb",
	)
}

// ExampleCustomConfig 自定义配置示例
func ExampleCustomConfig() {
	// 自定义配置
	config := &Config{
		Level:            DebugLevel,
		Format:           JSONFormat,
		OutputMode:       BothMode, // 同时输出到控制台和文件
		Filename:         "logs/custom.log",
		EnableCaller:     true,
		EnableStacktrace: true,
		Development:      false,
		Async:            true,
		BufferSize:       512 * 1024, // 512KB缓冲区
		FlushInterval:    2 * time.Second,
		Rotate: RotateConfig{
			MaxSize:    50, // 50MB
			MaxAge:     7,  // 7天
			MaxBackups: 5,  // 保留5个文件
			Compress:   true,
			LocalTime:  true,
			RotateTime: 12 * time.Hour, // 每12小时分割
		},
	}

	logger, err := NewLogger(config)
	if err != nil {
		panic(err)
	}

	SetGlobalLogger(logger)

	Info("自定义配置日志记录器启动")
}

// ExampleWithFields 字段使用示例
func ExampleWithFields() {
	// 初始化logger
	config := DefaultConfig()
	_ = InitGlobalLogger(config)

	// 创建带有固定字段的logger
	userLogger := WithFields(
		zap.String("module", "user"),
		zap.String("component", "service"),
	)

	// 使用带字段的logger
	userLogger.Info("用户服务启动")
	userLogger.Debug("初始化用户缓存")

	// 添加更多字段
	requestLogger := userLogger.WithRequestID("req-123456")
	requestLogger.WithUserID("user-789").Info("处理用户请求")

	// 链式调用
	GetGlobalLogger().WithRequestID("req-789012").
		WithUserID("user-456").
		WithFields(zap.String("action", "update_profile")).
		Info("更新用户资料")
}

// ExampleWithContext 上下文使用示例
func ExampleWithContext() {
	// 初始化logger
	config := DefaultConfig()
	_ = InitGlobalLogger(config)

	// 创建带有请求ID的上下文
	ctx := context.WithValue(context.Background(), requestIDKey, "req-123456")
	ctx = context.WithValue(ctx, userIDKey, "user-789")

	// 使用上下文记录日志
	InfoCtx(ctx, "开始处理请求")
	DebugCtx(ctx, "验证用户权限")
	WarnCtx(ctx, "用户权限不足，使用默认权限")

	// 格式化日志
	InfofCtx(ctx, "处理用户 %s 的请求", "张三")

	// 键值对日志
	InfowCtx(ctx, "请求处理完成",
		"duration", "150ms",
		"status", "success",
		"response_size", 1024,
	)
}

// ExampleErrorHandling 错误处理示例
func ExampleErrorHandling() {
	// 初始化logger
	config := DefaultConfig()
	_ = InitGlobalLogger(config)

	// 创建上下文
	ctx := context.Background()

	// 模拟错误
	err := errors.New("数据库连接失败")

	// 记录错误
	ErrorCtx(ctx, "业务逻辑错误",
		zap.Error(err),
		zap.String("module", "user_service"),
		zap.Int64("user_id", 12345),
	)
}

// ExampleStructuredLogging 结构化日志示例
func ExampleStructuredLogging() {
	// 初始化logger
	config := DefaultConfig()
	config.Format = JSONFormat
	_ = InitGlobalLogger(config)

	// 记录用户对象
	user := User{
		ID:   12345,
		Name: "张三",
		Age:  30,
	}

	// 结构化日志记录
	Info("用户信息",
		zap.Int64("user_id", user.ID),
		zap.String("user_name", user.Name),
		zap.Int("user_age", user.Age),
		zap.Time("created_at", time.Now()),
		zap.Bool("is_active", true),
		zap.Float64("score", 95.5),
		zap.Strings("tags", []string{"vip", "premium"}),
	)

	// 嵌套对象
	Infow("订单信息",
		"order_id", "order-123456",
		"user", map[string]interface{}{
			"id":   user.ID,
			"name": user.Name,
		},
		"items", []map[string]interface{}{
			{"id": 1, "name": "商品A", "price": 99.99},
			{"id": 2, "name": "商品B", "price": 199.99},
		},
		"total", 299.98,
		"status", "completed",
	)
}

// ExamplePerformanceLogging 性能日志示例
func ExamplePerformanceLogging() {
	// 初始化异步logger
	config := DefaultConfig()
	config.Async = true
	config.BufferSize = 1024 * 1024 // 1MB缓冲区
	config.FlushInterval = time.Second
	_ = InitGlobalLogger(config)

	// 高频日志记录（模拟）
	start := time.Now()
	for i := 0; i < 10000; i++ {
		Debug("高频调试日志",
			zap.Int("iteration", i),
			zap.Time("timestamp", time.Now()),
			zap.String("operation", "batch_process"),
		)
	}
	duration := time.Since(start)

	Info("批量日志记录完成",
		zap.Int("count", 10000),
		zap.Duration("duration", duration),
		zap.Float64("ops_per_second", float64(10000)/duration.Seconds()),
	)

	// 确保所有日志都被写入
	_ = Sync()
}

// ExampleCleanup 清理示例
func ExampleCleanup() {
	// 初始化logger
	config := DefaultConfig()
	logger, err := NewLogger(config)
	if err != nil {
		panic(err)
	}

	// 使用logger
	logger.Info("应用启动")
	logger.Debug("初始化完成")

	// 应用关闭时的清理
	logger.Info("应用正在关闭")

	// 同步所有日志
	if err := logger.Sync(); err != nil {
		logger.Error("同步日志失败", zap.Error(err))
	}

	// 关闭logger
	if err := logger.Close(); err != nil {
		logger.Error("关闭日志记录器失败", zap.Error(err))
	}
}
