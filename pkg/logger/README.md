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

## 依赖

- go.uber.org/zap
- gopkg.in/natefinch/lumberjack.v2

## 安装

```bash
go get -u go.uber.org/zap
go get -u gopkg.in/natefinch/lumberjack.v2
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
    
    err := logger.Init(config)
    if err != nil {
        panic(err)
    }
    defer logger.Sync()
    
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
    Filename:      "./logs/app.log",
    MaxSize:       100,
    MaxBackups:    10,
    MaxAge:        30,
    Compress:      true,
    
    // 日志格式控制
    DisableTimestamp: false,
    DisableCaller:    false,
    DisableTrace:     false,
    
    // 开发模式(会启用彩色日志和更详细的错误)
    Development: true,
    
    // 高级功能
    BufferSize:    256,
    FlushInterval: time.Second * 3,
    AsyncWrite:    true,
    
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
```

## 最佳实践

1. 在应用启动时初始化日志，在退出时调用`logger.Sync()`确保所有日志都被写入
2. 使用结构化日志代替简单的文本消息，便于后期分析
3. 使用命名日志记录器隔离不同模块的日志
4. 始终记录错误详情和上下文信息
5. 在生产环境禁用调试日志
6. 合理配置日志滚动策略，避免磁盘空间耗尽
7. 考虑对敏感信息进行脱敏处理 