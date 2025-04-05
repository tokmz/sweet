# Redis缓存包

Redis缓存包是一个对Redis客户端的封装，支持单机模式、集群模式和哨兵模式，同时集成了链路追踪功能。

## 特性

- 支持三种Redis运行模式：单机、集群、哨兵
- 支持所有常用的Redis数据类型操作：String、Hash、List、Set、ZSet
- 集成OpenTelemetry链路追踪
- 统一的错误处理和空值处理
- 可配置的连接池管理
- 支持自定义命令执行

## 安装依赖

```bash
go get github.com/redis/go-redis/v9
go get go.opentelemetry.io/otel
```

## 快速开始

### 创建Redis客户端

```go
import (
    "context"
    "log"
    "time"
    
    "github.com/your-org/sweet/pkg/cache"
)

func main() {
    // 单机模式配置
    config := cache.Config{
        Mode: cache.ModeSingle,
        Single: struct {
            Addr string `json:"addr"`
        }{
            Addr: "localhost:6379",
        },
        Password:    "",
        DB:          0,
        PoolSize:    10,
        EnableTrace: true,
    }
    
    // 创建客户端
    client, err := cache.NewClient(config)
    if err != nil {
        log.Fatalf("创建Redis客户端失败: %v", err)
    }
    defer client.Close()
    
    // 使用客户端
    ctx := context.Background()
    
    // 设置值
    err = client.Set(ctx, "key", "value", time.Hour)
    if err != nil {
        log.Printf("设置值失败: %v", err)
    }
    
    // 获取值
    val, err := client.Get(ctx, "key")
    if err != nil {
        log.Printf("获取值失败: %v", err)
    } else {
        log.Printf("值: %s", val)
    }
}
```

### 集群模式

```go
config := cache.Config{
    Mode: cache.ModeCluster,
    Cluster: struct {
        Addrs []string `json:"addrs"`
    }{
        Addrs: []string{
            "127.0.0.1:7000",
            "127.0.0.1:7001",
            "127.0.0.1:7002",
        },
    },
    Password:    "",
    PoolSize:    10,
    EnableTrace: true,
}
```

### 哨兵模式

```go
config := cache.Config{
    Mode: cache.ModeSentinel,
    Sentinel: struct {
        MasterName string   `json:"master_name"`
        Addrs      []string `json:"addrs"`
    }{
        MasterName: "mymaster",
        Addrs: []string{
            "127.0.0.1:26379",
            "127.0.0.1:26380",
            "127.0.0.1:26381",
        },
    },
    Password:    "",
    DB:          0,
    PoolSize:    10,
    EnableTrace: true,
}
```

## 配置项

| 配置项 | 说明 | 默认值 |
|-------|------|--------|
| Mode | Redis模式：单机(single)、集群(cluster)、哨兵(sentinel) | 必填 |
| Single.Addr | 单机模式Redis地址 | 必填(单机模式) |
| Cluster.Addrs | 集群模式Redis节点地址列表 | 必填(集群模式) |
| Sentinel.MasterName | 哨兵模式主节点名称 | 必填(哨兵模式) |
| Sentinel.Addrs | 哨兵节点地址列表 | 必填(哨兵模式) |
| Username | 用户名 | "" |
| Password | 密码 | "" |
| DB | 数据库编号(仅单机和哨兵模式) | 0 |
| PoolSize | 连接池大小 | 10 |
| MinIdleConns | 最小空闲连接数 | 5 |
| IdleTimeout | 连接最大空闲时间(秒) | 300 |
| ConnTimeout | 连接超时时间(毫秒) | 5000 |
| ReadTimeout | 读取超时时间(毫秒) | 3000 |
| WriteTimeout | 写入超时时间(毫秒) | 3000 |
| MaxRetries | 最大重试次数 | 3 |
| RetryDelay | 重试间隔(毫秒) | 100 |
| MinRetryBackoff | 最小重试间隔(毫秒) | 8 |
| MaxRetryBackoff | 最大重试间隔(毫秒) | 512 |
| EnableTrace | 是否开启链路追踪 | false |

## 支持的Redis操作

### 字符串操作

- `Get`: 获取键值
- `Set`: 设置键值
- `SetNX`: 不存在时设置键值
- `Del`: 删除键
- `Exists`: 检查键是否存在
- `Expire`: 设置过期时间
- `TTL`: 获取剩余生存时间

### 哈希表操作

- `HGet`: 获取哈希表字段值
- `HSet`: 设置哈希表字段值
- `HGetAll`: 获取哈希表所有字段和值
- `HDel`: 删除哈希表字段
- `HExists`: 检查哈希表字段是否存在

### 列表操作

- `LPush`: 将值插入到列表头部
- `RPush`: 将值插入到列表尾部
- `LPop`: 移除并返回列表第一个元素
- `RPop`: 移除并返回列表最后一个元素
- `LRange`: 获取列表指定范围内的元素
- `LLen`: 获取列表长度

### 集合操作

- `SAdd`: 向集合添加一个或多个成员
- `SMembers`: 返回集合中的所有成员
- `SRem`: 移除集合中一个或多个成员
- `SIsMember`: 判断成员是否是集合的成员
- `SCard`: 获取集合的成员数

### 有序集合操作

- `ZAdd`: 向有序集合添加一个或多个成员
- `ZRange`: 返回有序集合中指定范围的成员
- `ZRangeWithScores`: 返回有序集合中指定范围的成员和分数
- `ZRem`: 移除有序集合中的一个或多个成员
- `ZCard`: 获取有序集合的成员数

### 自定义命令

- `Do`: 执行自定义命令 