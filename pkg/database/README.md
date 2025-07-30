# Database Package

基于 GORM 的高性能数据库封装包，提供读写分离、链路追踪、慢查询监控等企业级功能。

## 功能特性

### 🚀 核心功能
- **读写分离**: 使用 GORM 官方 dbresolver 插件实现主从库分离
- **链路追踪**: 集成 OpenTelemetry，全链路 SQL 执行追踪
- **慢查询监控**: 可配置慢查询阈值，自动记录慢查询日志
- **多级别日志**: 支持 Silent、Error、Warn、Info 四个级别
- **连接池管理**: 灵活的连接池配置，支持连接生命周期管理
- **健康检查**: 主从库连接状态监控
- **事务支持**: 完整的事务管理功能

### 📊 监控与观测
- SQL 执行时间统计
- 连接池状态监控
- 慢查询日志记录
- OpenTelemetry 链路追踪
- 数据库连接统计信息

## 安装依赖

```bash
go get gorm.io/gorm
go get gorm.io/driver/mysql
go get gorm.io/plugin/dbresolver
go get go.opentelemetry.io/otel
go get go.opentelemetry.io/otel/trace
```

## 快速开始

### 基本使用

```go
package main

import (
    "context"
    "log"
    "time"
    
    "zero/pkg/database"
    "gorm.io/gorm/logger"
)

func main() {
    // 创建配置
    config := &database.Config{
        Master: database.MasterConfig{
            DSN: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
        },
        Slaves: []database.SlaveConfig{
            {DSN: "user:password@tcp(slave1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"},
        },
        Pool: database.PoolConfig{
            MaxIdleConns:    10,
            MaxOpenConns:    100,
            ConnMaxLifetime: time.Hour,
            ConnMaxIdleTime: time.Minute * 30,
        },
        Log: database.LogConfig{
            Level:    logger.Info,
            Colorful: true,
        },
        SlowQuery: database.SlowQueryConfig{
            Enabled:   true,
            Threshold: time.Millisecond * 200,
        },
        Tracing: database.TracingConfig{
            Enabled:     true,
            ServiceName: "my-service",
            RecordSQL:   true,
        },
    }
    
    // 创建客户端
    client, err := database.NewClient(config)
    if err != nil {
        log.Fatal(err)
    }
    defer client.Close()
    
    // 测试连接
    ctx := context.Background()
    if err := client.Ping(ctx); err != nil {
        log.Fatal(err)
    }
    
    // 使用数据库
    db := client.DB()
    // ... 你的业务逻辑
}
```

### 使用默认配置

```go
// 使用默认配置
config := database.DefaultConfig()
config.Master = "your-dsn-here"

client, err := database.NewClient(config)
if err != nil {
    log.Fatal(err)
}
defer client.Close()
```

## 详细功能说明

### 1. 读写分离

```go
// 自动读写分离（写操作自动使用主库，读操作使用从库）
db := client.DB()

// 创建操作（自动使用主库）
db.Create(&user)

// 查询操作（自动使用从库）
db.Find(&users)

// 强制使用主库
client.Master().Find(&users)

// 强制使用从库
client.Slave().Find(&users)
```

### 2. 事务管理

```go
err := client.Transaction(ctx, func(tx *gorm.DB) error {
    // 事务中的操作
    if err := tx.Create(&user1).Error; err != nil {
        return err
    }
    if err := tx.Create(&user2).Error; err != nil {
        return err
    }
    return nil
})
```

### 3. 健康检查

```go
// 检查主从库连接状态
if err := client.HealthCheck(ctx); err != nil {
    log.Printf("Database health check failed: %v", err)
}

// 获取连接池统计信息
stats, err := client.Stats()
if err == nil {
    fmt.Printf("Connection stats: %+v\n", stats)
}
```

### 4. 慢查询监控

```go
// 获取慢查询日志
slowQueries, err := client.GetSlowQueries(ctx, 10)
if err == nil {
    for _, query := range slowQueries {
        fmt.Printf("Slow query: %+v\n", query)
    }
}
```

### 5. 动态日志级别

```go
// 动态调整日志级别
client.SetLogLevel("warn")  // silent, error, warn, info
```

## 配置说明

### 主库配置 (MasterConfig)

```go
type MasterConfig struct {
    DSN string `json:"dsn" yaml:"dsn"` // 主库连接字符串
}
```

### 从库配置 (SlaveConfig)

```go
type SlaveConfig struct {
    DSN string `json:"dsn" yaml:"dsn"` // 从库连接字符串
}
```

### 连接池配置 (PoolConfig)

```go
type PoolConfig struct {
    MaxIdleConns    int           // 最大空闲连接数 (默认: 10)
    MaxOpenConns    int           // 最大连接数 (默认: 100)
    ConnMaxLifetime time.Duration // 连接最大生命周期 (默认: 1小时)
    ConnMaxIdleTime time.Duration // 连接最大空闲时间 (默认: 30分钟)
}
```

### 日志配置 (LogConfig)

```go
type LogConfig struct {
    Level                     logger.LogLevel // 日志级别
    Colorful                  bool           // 彩色输出
    IgnoreRecordNotFoundError bool           // 忽略记录未找到错误
    ParameterizedQueries      bool           // 参数化查询
}
```

**日志级别说明:**
- `logger.Silent`: 静默模式，不输出任何日志
- `logger.Error`: 只输出错误日志
- `logger.Warn`: 输出警告和错误日志
- `logger.Info`: 输出所有日志（包括 SQL 语句）

### 慢查询配置 (SlowQueryConfig)

```go
type SlowQueryConfig struct {
    Enabled   bool          // 是否启用慢查询监控
    Threshold time.Duration // 慢查询阈值 (默认: 200ms)
}
```

### OpenTelemetry 配置 (TracingConfig)

```go
type TracingConfig struct {
    Enabled      bool   // 是否启用链路追踪
    ServiceName  string // 服务名称
    RecordSQL    bool   // 是否记录 SQL 语句
    RecordParams bool   // 是否记录查询参数
}
```

## 最佳实践

### 1. 生产环境配置建议

```go
config := &database.Config{
    Master: database.MasterConfig{
        DSN: os.Getenv("DB_MASTER_DSN"),
    },
    Slaves: []database.SlaveConfig{
        {DSN: os.Getenv("DB_SLAVE1_DSN")},
        {DSN: os.Getenv("DB_SLAVE2_DSN")},
    },
    Pool: database.PoolConfig{
        MaxIdleConns:    20,
        MaxOpenConns:    200,
        ConnMaxLifetime: time.Hour * 2,
        ConnMaxIdleTime: time.Hour,
    },
    Log: database.LogConfig{
        Level:                     logger.Warn, // 生产环境建议使用 Warn 级别
        Colorful:                  false,       // 生产环境关闭颜色
        IgnoreRecordNotFoundError: true,
        ParameterizedQueries:      true,
    },
    SlowQuery: database.SlowQueryConfig{
        Enabled:   true,
        Threshold: time.Millisecond * 100, // 根据业务需求调整
    },
    Tracing: database.TracingConfig{
        Enabled:      true,
        ServiceName:  "production-service",
        RecordSQL:    false, // 生产环境可能不记录 SQL
        RecordParams: false,
    },
}
```

### 2. 开发环境配置建议

```go
config := &database.Config{
    Master: database.MasterConfig{
        DSN: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
    },
    Log: database.LogConfig{
        Level:    logger.Info, // 开发环境显示所有日志
        Colorful: true,        // 开启彩色输出
    },
    SlowQuery: database.SlowQueryConfig{
        Enabled:   true,
        Threshold: time.Millisecond * 50, // 开发环境更严格的慢查询阈值
    },
    Tracing: database.TracingConfig{
        Enabled:      true,
        ServiceName:  "dev-service",
        RecordSQL:    true, // 开发环境记录 SQL
        RecordParams: true, // 开发环境记录参数
    },
}
```

### 3. 错误处理

```go
// 统一错误处理
func handleDBError(err error) {
    if err != nil {
        if errors.Is(err, gorm.ErrRecordNotFound) {
            // 处理记录未找到
            log.Println("Record not found")
        } else {
            // 处理其他数据库错误
            log.Printf("Database error: %v", err)
        }
    }
}
```

### 4. 上下文使用

```go
// 始终使用上下文
ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
defer cancel()

// 在所有数据库操作中传递上下文
result := client.DB().WithContext(ctx).Find(&users)
```

## 注意事项

1. **DSN 格式**: 确保 DSN 包含 `parseTime=True` 参数以正确处理时间类型
2. **连接池配置**: 根据应用负载合理配置连接池参数
3. **慢查询阈值**: 根据业务需求设置合适的慢查询阈值
4. **日志级别**: 生产环境建议使用 Warn 或 Error 级别
5. **链路追踪**: 生产环境可能需要关闭 SQL 和参数记录以保护敏感信息
6. **健康检查**: 定期执行健康检查以监控数据库连接状态
7. **事务管理**: 合理使用事务，避免长时间持有连接

## 监控指标

包提供以下监控指标：

- 连接池状态（活跃连接数、空闲连接数等）
- SQL 执行时间
- 慢查询统计
- 错误率统计
- OpenTelemetry 链路追踪数据

这些指标可以集成到 Prometheus、Grafana 等监控系统中。