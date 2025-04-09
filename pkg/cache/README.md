# Sweet Redis缓存包

## 概述
Sweet Redis缓存包是一个功能完备的Redis客户端封装，提供统一的接口支持单机模式、集群模式和哨兵模式，同时集成了OpenTelemetry链路追踪功能，帮助开发者快速构建高性能、可观测的分布式缓存系统。

## 特性
- **多模式支持**：支持Redis单机模式、集群模式和哨兵模式
- **统一接口**：不同模式下使用相同的接口，方便切换
- **链路追踪**：集成OpenTelemetry，支持分布式追踪
- **超时控制**：所有操作支持超时设置
- **连接池管理**：支持连接池配置，优化性能
- **错误处理**：统一的错误处理机制，增强可维护性
- **类型安全**：提供类型安全的API，减少运行时错误

## 安装
```bash
go get github.com/your-organization/sweet/pkg/cache
```

## 配置
创建Redis客户端需要提供配置，支持以下参数：

```go
// 创建Redis配置
config := cache.Config{
    // 基础配置
    Mode:         cache.ModeSingle,  // 支持 ModeSingle, ModeCluster, ModeSentinel
    Username:     "default",         // Redis用户名（可选）
    Password:     "password",        // Redis密码（可选）
    DB:           0,                 // 数据库编号
    EnableTrace:  true,              // 是否启用链路追踪
    ExecTimeout:  100,               // 执行超时时间（毫秒）
    
    // 连接配置
    ConnTimeout:    500,            // 连接超时（毫秒）
    ReadTimeout:    500,            // 读取超时（毫秒）
    WriteTimeout:   500,            // 写入超时（毫秒）
    PoolSize:       10,             // 连接池大小
    MinIdleConns:   2,              // 最小空闲连接数
    IdleTimeout:    60,             // 空闲超时（秒）
    
    // 重试配置
    MaxRetries:      3,             // 最大重试次数
    RetryDelay:      100,           // 重试延迟（毫秒）
    MinRetryBackoff: 8,             // 最小重试退避时间（毫秒）
    MaxRetryBackoff: 512,           // 最大重试退避时间（毫秒）
    
    // 单机模式配置
    Single: cache.SingleConfig{
        Addr: "localhost:6379",     // 单机地址
    },
    
    // 集群模式配置
    Cluster: cache.ClusterConfig{
        Addrs: []string{            // 集群节点地址
            "127.0.0.1:7001",
            "127.0.0.1:7002",
            "127.0.0.1:7003",
        },
    },
    
    // 哨兵模式配置
    Sentinel: cache.SentinelConfig{
        MasterName: "mymaster",     // 主节点名称
        Addrs: []string{            // 哨兵节点地址
            "127.0.0.1:26379",
            "127.0.0.1:26380",
            "127.0.0.1:26381",
        },
    },
}
```

## 使用示例

### 客户端创建
```go
import (
    "context"
    "fmt"
    "time"
    
    "github.com/your-organization/sweet/pkg/cache"
)

func main() {
    // 创建配置
    config := cache.Config{
        Mode: cache.ModeSingle,
        Single: cache.SingleConfig{
            Addr: "localhost:6379",
        },
    }
    
    // 创建客户端
    client, err := cache.NewClient(config)
    if err != nil {
        panic(err)
    }
    defer client.Close()
    
    // 使用客户端...
}
```

### 字符串操作
```go
// 设置值
err := client.Set(ctx, "key", "value", 10*time.Minute)
if err != nil {
    // 处理错误
}

// 获取值
val, err := client.Get(ctx, "key")
if err != nil {
    if errors.Is(err, cache.ErrKeyNotExists) {
        // 键不存在处理
    } else {
        // 其他错误处理
    }
}

// 检查并获取值
val, exists, err := client.GetWithExists(ctx, "key")
if err != nil {
    // 处理错误
} else if !exists {
    // 键不存在处理
} else {
    // 使用值
}
```

### 哈希表操作
```go
// 设置哈希表字段
n, err := client.HSet(ctx, "hash", "field1", "value1", "field2", "value2")
if err != nil {
    // 处理错误
}

// 获取哈希表字段
val, err := client.HGet(ctx, "hash", "field1")
if err != nil {
    // 处理错误
}

// 获取所有字段和值
fields, err := client.HGetAll(ctx, "hash")
if err != nil {
    // 处理错误
}
```

### 列表操作
```go
// 添加元素到列表头部
n, err := client.LPush(ctx, "list", "value1", "value2")
if err != nil {
    // 处理错误
}

// 从列表尾部弹出元素
val, err := client.RPop(ctx, "list")
if err != nil {
    // 处理错误
}

// 获取列表范围
vals, err := client.LRange(ctx, "list", 0, -1)
if err != nil {
    // 处理错误
}
```

### 集合操作
```go
// 添加元素到集合
n, err := client.SAdd(ctx, "set", "member1", "member2")
if err != nil {
    // 处理错误
}

// 获取集合成员
members, err := client.SMembers(ctx, "set")
if err != nil {
    // 处理错误
}

// 检查成员是否存在
exists, err := client.SIsMember(ctx, "set", "member1")
if err != nil {
    // 处理错误
}
```

### 有序集合操作
```go
// 添加元素到有序集合
members := []*cache.Z{
    {Score: 1.0, Member: "member1"},
    {Score: 2.0, Member: "member2"},
}
n, err := client.ZAdd(ctx, "zset", members...)
if err != nil {
    // 处理错误
}

// 获取有序集合范围（带分数）
membersWithScores, err := client.ZRangeWithScores(ctx, "zset", 0, -1)
if err != nil {
    // 处理错误
}
```

### 执行自定义命令
```go
// 执行PING命令
result, err := client.Do(ctx, "PING")
if err != nil {
    // 处理错误
}

// 执行INFO命令
info, err := client.Do(ctx, "INFO", "server")
if err != nil {
    // 处理错误
}
```

## 错误处理

包提供了以下错误类型：

- `ErrKeyNotExists`: 键不存在
- `ErrEmptyAddrs`: 未提供服务器地址
- `ErrEmptyMasterSet`: 哨兵模式未提供主节点名称
- `ErrInvalidMode`: 无效的Redis模式
- `ErrCommandFailed`: 命令执行失败
- `ConnectionError`: 连接错误，包含地址和底层错误
- `CommandError`: 命令错误，包含命令名称和底层错误

示例：
```go
val, err := client.Get(ctx, "non-existing-key")
if err != nil {
    if errors.Is(err, cache.ErrKeyNotExists) {
        // 键不存在处理
    } else if cmdErr, ok := err.(*cache.CommandError); ok {
        fmt.Printf("命令 %s 执行失败: %v\n", cmdErr.Command, cmdErr.Err)
    } else if connErr, ok := err.(*cache.ConnectionError); ok {
        fmt.Printf("连接 %s 失败: %v\n", connErr.Addr, connErr.Err)
    } else {
        // 其他错误处理
    }
}
```

## 链路追踪

启用链路追踪需设置`EnableTrace: true`，并确保应用初始化了OpenTelemetry：

```go
import (
    "go.opentelemetry.io/otel"
    "go.opentelemetry.io/otel/exporters/jaeger"
    "go.opentelemetry.io/otel/sdk/resource"
    sdktrace "go.opentelemetry.io/otel/sdk/trace"
    semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

func initTracer() {
    // 初始化Jaeger导出器
    exporter, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint("http://localhost:14268/api/traces")))
    if err != nil {
        panic(err)
    }
    
    // 创建资源
    res, err := resource.New(context.Background(),
        resource.WithAttributes(
            semconv.ServiceNameKey.String("service-name"),
        ),
    )
    if err != nil {
        panic(err)
    }
    
    // 注册提供者
    tp := sdktrace.NewTracerProvider(
        sdktrace.WithBatcher(exporter),
        sdktrace.WithResource(res),
    )
    otel.SetTracerProvider(tp)
}
```

## 性能优化

建议：
1. 根据实际需求调整连接池大小和空闲连接数
2. 设置合理的超时时间避免长时间阻塞
3. 大批量操作考虑使用Pipeline
4. 避免在热路径上创建和销毁客户端

## 许可证
MIT License 