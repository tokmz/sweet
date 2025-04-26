# Sweet Logger

基于zap和lumberjack的高性能日志包，支持同时输出到控制台（文本格式）和文件（JSON格式），并提供日志分割功能。支持结构化日志和Sugared风格两种API。

## 功能特性

- 基于uber-go/zap的高性能日志记录
- 支持控制台输出（文本格式）和文件输出（JSON格式）
- 使用lumberjack进行日志文件分割和轮转
- 双重API风格：结构化日志（类型安全）和Sugared风格（简洁灵活）
- 结构化日志记录，支持各种类型的字段
- 支持不同日志级别（Debug, Info, Warn, Error, Fatal）
- 支持命名日志记录器和预设字段
- 支持开发模式，显示调用者信息和堆栈跟踪
- 线程安全的实现，支持并发环境下的安全使用
- 自动创建日志目录，简化配置过程

## 依赖

- go.uber.org/zap
- gopkg.in/natefinch/lumberjack.v2

## 安装

```bash
go get -u go.uber.org/zap
go get -u gopkg.in/natefinch/lumberjack.v2
```

## 线程安全

Sweet Logger 实现了完整的线程安全机制，确保在并发环境下安全使用：

- 使用互斥锁保护全局默认日志记录器的访问
- 支持多次初始化，后续初始化会安全地替换之前的日志记录器
- 所有API调用都是线程安全的，可以在多个goroutine中并发使用
- 自动创建日志目录，无需手动创建

## 最佳实践

### 初始化

- 在应用启动时尽早初始化日志系统
- 使用`defer logger.Sync()`确保程序退出前刷新缓冲区
- 在高并发场景下，避免频繁调用`Init()`方法

### 性能考虑

- 对于高频日志，考虑使用`With()`创建预设字段的日志记录器
- 在性能敏感的代码路径上，可以先检查日志级别再构造复杂的日志消息
- 使用结构化日志而非字符串拼接，便于后期分析

### 日志级别使用建议

- `Debug`: 详细的开发调试信息，生产环境通常禁用
- `Info`: 常规操作信息，表示服务正常运行状态
- `Warn`: 潜在问题警告，不影响主要功能但需关注
- `Error`: 错误信息，表示功能受损但服务仍在运行
- `Fatal`: 严重错误，会导致应用终止运行

### 并发使用示例

```go
func worker(id int) {
    // 在多个goroutine中安全使用logger
    logger.Info("工作协程开始", logger.Int("worker_id", id))
    // 执行任务...
    logger.Info("工作协程结束", logger.Int("worker_id", id))
}

func main() {
    // 初始化日志
    config := logger.Config{
        Level:         logger.InfoLevel,
        EnableConsole: true,
        EnableFile:    true,
        Filename:      "./logs/app.log",
    }
    err := logger.Init(config)
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
    
    // 启动多个goroutine
    var wg sync.WaitGroup
    for i := 0; i < 10; i++ {
        wg.Add(1)
        go func(id int) {
            defer wg.Done()
            worker(id)
        }(i)
    }
    wg.Wait()
}
```

## 使用示例

### 初始化

```go
package main

import (
    "github.com/your-org/sweet/pkg/logger"
)

func main() {
    // 使用默认配置初始化
    config := logger.Config{
        Level:         logger.InfoLevel,
        EnableConsole: true,
        EnableFile:    true,
        Filename:      "./logs/app.log",
        MaxSize:       100, // MB
        MaxBackups:    10,
        MaxAge:        30,  // 天
        Compress:      true,
        Development:   false,
    }
    
    // 初始化日志系统，会自动创建日志目录
    err := logger.Init(config)
    if err != nil {
        panic(err)
    }
    defer logger.Sync() // 确保缓冲区数据写入
    
    // 开始使用
    logger.Info("日志系统已初始化")
}
```

### 使用结构化日志

```go
// 基本使用
logger.Debug("调试信息")
logger.Info("处理请求", logger.String("path", "/api/users"), logger.Int("status", 200))
logger.Warn("警告", logger.Err(err))
logger.Error("请求处理失败", logger.Err(err), logger.Int("status", 500))

// 创建子日志记录器
userLogger := logger.Named("user")
userLogger.Info("用户登录成功", logger.String("username", "zhang"))

// 创建带有预设字段的日志记录器
requestLogger := logger.With(
    logger.String("request_id", "req-123"),
    logger.String("user_agent", "Mozilla/5.0..."),
)
requestLogger.Info("收到请求") // 会自动包含预设字段
```

### 使用Sugared风格API

```go
// 获取Sugar日志记录器
sugar := logger.Sugar()

// 使用printf风格API
sugar.Debugf("调试信息: %s", message)
sugar.Infof("用户 %s 登录 IP: %s", username, ip)
sugar.Warnf("警告: %v", err)
sugar.Errorf("错误: %v", err)

// 使用键值对风格API
sugar.Infow("请求完成",
    "method", "GET",
    "path", "/api/users",
    "status", 200,
    "latency", 30.2,
)

// 全局Sugar风格API
logger.Infof("处理请求: %s", path)
logger.Warnw("请求超时",
    "path", path,
    "latency", latency,
)
```

### 使用常用字段辅助函数

```go
logger.Info("处理HTTP请求",
    logger.RequestID("req-123"),
    logger.Method("POST"),
    logger.URL("/api/users"),
    logger.StatusCode(201),
    logger.Latency(time.Millisecond * 35),
    logger.UserID(10086),
)
```

### 高级配置

```go
config := logger.Config{
    Level:         logger.InfoLevel,
    EnableConsole: true,
    EnableFile:    true,
    Filename:      "./logs/app.log", // 日志目录会自动创建
    MaxSize:       100,  // 单个日志文件最大尺寸，单位MB
    MaxBackups:    10,   // 保留的旧日志文件最大数量
    MaxAge:        30,   // 保留旧日志文件的最大天数
    Compress:      true, // 是否压缩旧日志文件
    
    // 日志格式控制
    DisableTimestamp: false, // 是否禁用时间戳
    DisableCaller:    false, // 是否禁用调用者信息
    DisableTrace:     false, // 是否禁用堆栈跟踪
    
    // 开发模式(会启用彩色日志和更详细的错误)
    Development: true,
    
    // 高级功能
    BufferSize:    256,               // 日志缓冲区大小
    FlushInterval: time.Second * 3,   // 日志刷新间隔
    AsyncWrite:    true,              // 是否启用异步写入
    
    // 采样配置(每秒只记录前100条，之后每10条记录1条)
    SamplingConfig: &logger.SamplingConfig{
        Initial:    100,
        Thereafter: 10,
        Interval:   time.Second,
    },
    
    // 过滤配置
    FilterConfig: &logger.FilterConfig{
        MinLevel: logger.WarnLevel,
        Keywords: []string{"敏感", "密码"},
    },
    
    // 安全配置
    SecurityConfig: &logger.SecurityConfig{
        EnableMasking: true,
        MaskingFields: []string{"password", "credit_card"},
    },
}

// 初始化日志系统，会自动创建日志目录并进行线程安全的初始化
err := logger.Init(config)
if err != nil {
    panic(err)
}
```

## 最佳实践

1. 在应用启动时初始化日志，在退出时调用`logger.Sync()`确保所有日志都被写入
2. 使用结构化日志代替简单的文本消息，便于后期分析
3. 使用命名日志记录器隔离不同模块的日志
4. 始终记录错误详情和上下文信息
5. 在生产环境禁用调试日志
6. 合理配置日志滚动策略，避免磁盘空间耗尽
7. 考虑对敏感信息进行脱敏处理