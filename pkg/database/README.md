# Database包

## 简介

Database包是一个基于GORM实现的MySQL数据库访问层，提供了以下主要功能：

- 支持主从读写分离配置
- 支持连接池管理
- 集成OpenTelemetry链路追踪
- 提供事务管理和自动重试机制
- 支持慢查询日志
- 支持上下文传递

## 配置项

数据库配置通过`Config`结构体定义：

```go
type Config struct {
    // 主库连接DSN
    Master string `json:"master"` // dsn格式：user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
    // 从库连接DSN数组
    Slaves []string `json:"slaves"`
    // 连接池配置
    Pool ConnPool `json:"pool"`
    // 是否开启日志
    EnableLog bool `json:"enable_log"`
    // 是否开启链路追踪
    EnableTrace bool `json:"enable_trace"`
    // 慢查询阈值（毫秒）
    SlowThreshold int `json:"slow_threshold"`
    // 数据库驱动类型
    DriverType string `json:"driver_type"`
    // 查询超时设置（毫秒），0表示不设置超时
    QueryTimeout int `json:"query_timeout"`
    // 最大重试次数
    MaxRetries int `json:"max_retries"`
    // 重试间隔（毫秒）
    RetryDelay int `json:"retry_delay"`
}

type ConnPool struct {
    MaxIdleConns    int `json:"max_idle_conns"`
    MaxOpenConns    int `json:"max_open_conns"`
    ConnMaxLifetime int `json:"conn_max_lifetime"`
}
```

## 初始化数据库

使用`Init`函数初始化数据库连接，它将返回配置好的`*gorm.DB`实例：

```go
import "your-project/pkg/database"

func main() {
    config := database.Config{
        Master:         "user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
        Slaves:         []string{"user:pass@tcp(127.0.0.1:3307)/dbname?charset=utf8mb4&parseTime=True&loc=Local"},
        EnableLog:      true,
        EnableTrace:    true,
        SlowThreshold:  500, // 慢查询阈值500毫秒
        DriverType:     "mysql",
        QueryTimeout:   3000, // 3秒超时
        MaxRetries:     3,
        RetryDelay:     100, // 100毫秒
        Pool: database.ConnPool{
            MaxIdleConns:    10,
            MaxOpenConns:    100,
            ConnMaxLifetime: 3600, // 1小时
        },
    }
    
    db, err := database.Init(config)
    if err != nil {
        panic(err)
    }
    
    // 使用返回的db实例进行数据库操作
}
```

## 基本使用

### 基本数据库操作

初始化后，可以使用返回的`*gorm.DB`实例进行数据库操作：

```go
import (
    "context"
    "gorm.io/gorm"
    "your-project/pkg/database"
)

func GetUser(ctx context.Context, db *gorm.DB, id uint) (*User, error) {
    var user User
    if err := db.WithContext(ctx).First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 上下文传递

推荐使用带上下文的方法：

```go
func GetUser(ctx context.Context, db *gorm.DB, id uint) (*User, error) {
    var user User
    if err := database.WithContext(db, ctx).First(&user, id).Error; err != nil {
        return nil, err
    }
    return &user, nil
}
```

### 主从库控制

强制使用主库或从库：

```go
// 强制使用主库
dbMaster := database.MasterDB(db, ctx)

// 强制使用从库
dbSlave := database.SlaveDB(db, ctx)
```

### 事务管理

使用事务并支持自动重试：

```go
err := database.Transaction(db, ctx, func(tx *gorm.DB) error {
    // 在事务中执行数据库操作
    if err := tx.Create(&user).Error; err != nil {
        return err
    }
    
    if err := tx.Create(&userProfile).Error; err != nil {
        return err
    }
    
    return nil
})
```

自定义重试次数：

```go
err := database.TransactionWithRetry(db, ctx, func(tx *gorm.DB) error {
    // 在事务中执行数据库操作
    return nil
}, 5) // 最多重试5次
```

## 链路追踪

当`EnableTrace`设置为`true`时，数据库操作将被自动记录到OpenTelemetry跟踪系统中，包括：

- 数据库操作类型（增删改查）
- SQL语句（当recordSQL=true）
- 影响的行数
- 错误信息
- 重试情况

## 最佳实践

1. 总是使用上下文传递，确保链路追踪和超时控制生效
2. 对于写操作，使用事务包装以确保数据一致性
3. 对于只读操作，可考虑使用`SlaveDB`减轻主库压力
4. 配置合理的连接池大小，避免连接资源不足或过多占用
5. 设置合适的慢查询阈值，以便及时发现性能问题
6. 在应用程序中妥善管理数据库实例，使用依赖注入等方式传递数据库实例 