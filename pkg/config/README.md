# Config - Viper配置管理包

这是一个基于 [Viper](https://github.com/spf13/viper) 的配置管理包，提供了简单易用的配置加载、管理和监听功能。

## 特性

- 🚀 **多格式支持**: 支持 YAML、JSON、TOML、Properties、INI 等多种配置文件格式
- 🌍 **环境变量**: 自动读取和绑定环境变量
- 📁 **多路径搜索**: 支持在多个路径中搜索配置文件
- 🔄 **配置监听**: 支持配置文件变化监听和热重载
- ☁️ **远程配置**: 支持从 etcd、Consul、Firestore 等远程配置中心加载配置
- 🏭 **工厂模式**: 提供工厂模式创建不同环境的配置管理器
- 🔧 **结构体绑定**: 支持将配置直接解析到 Go 结构体
- 🐛 **调试模式**: 内置调试功能，方便排查配置问题

## 安装

```bash
go get github.com/spf13/viper
go get github.com/fsnotify/fsnotify
go get github.com/spf13/pflag
```

## 快速开始

### 基本使用

```go
package main

import (
    "log"
    "your-project/pkg/config"
)

func main() {
    // 创建默认配置管理器
    manager, err := config.Default()
    if err != nil {
        log.Fatal(err)
    }

    // 加载配置
    if err := manager.Load(); err != nil {
        log.Fatal(err)
    }

    // 获取配置值
    host := manager.GetString("server.host")
    port := manager.GetInt("server.port")
    debug := manager.GetBool("debug")

    log.Printf("Server: %s:%d, Debug: %t", host, port, debug)
}
```

### 结构体绑定

```go
type AppConfig struct {
    Server struct {
        Host string `mapstructure:"host"`
        Port int    `mapstructure:"port"`
    } `mapstructure:"server"`
    Database struct {
        DSN string `mapstructure:"dsn"`
    } `mapstructure:"database"`
}

func main() {
    manager, _ := config.Development()
    manager.Load()

    var appConfig AppConfig
    if err := manager.Unmarshal(&appConfig); err != nil {
        log.Fatal(err)
    }

    log.Printf("Config: %+v", appConfig)
}
```

## 配置文件示例

### YAML 格式 (config.yaml)

```yaml
server:
  host: localhost
  port: 8080
  mode: debug
  read_timeout: 30s
  write_timeout: 30s

database:
  driver: mysql
  dsn: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
  max_idle_conns: 10
  max_open_conns: 100
  conn_max_lifetime: 1h

redis:
  addr: localhost:6379
  password: ""
  db: 0
  pool_size: 10

logger:
  level: info
  format: json
  output: both
  file_path: ./logs/app.log

auth:
  jwt_secret: your-secret-key
  access_token_expire: 24h
  refresh_token_expire: 168h
  token_issuer: your-app
```

### JSON 格式 (config.json)

```json
{
  "server": {
    "host": "localhost",
    "port": 8080,
    "mode": "debug"
  },
  "database": {
    "driver": "mysql",
    "dsn": "user:password@tcp(localhost:3306)/dbname"
  }
}
```

## 环境变量

配置管理器支持自动读取环境变量，环境变量的命名规则为：`{前缀}_{配置键}`

例如，配置键 `server.host` 对应的环境变量为 `APP_SERVER_HOST`（假设前缀为 `APP`）。

```bash
# 设置环境变量
export APP_SERVER_HOST=0.0.0.0
export APP_SERVER_PORT=9090
export APP_DATABASE_DSN="user:pass@tcp(db:3306)/mydb"
```

## 不同环境配置

### 开发环境

```go
manager, err := config.Development()
```

- 配置文件搜索路径：`[".", "./configs", "./config"]`
- 启用配置文件监听
- 启用调试模式

### 生产环境

```go
manager, err := config.Production()
```

- 配置文件搜索路径：`["/etc/app", "./configs"]`
- 禁用配置文件监听
- 禁用调试模式

### 测试环境

```go
manager, err := config.Testing()
```

- 配置文件搜索路径：`[".", "./testdata", "./configs"]`
- 禁用配置文件监听
- 启用调试模式

## 配置文件监听

```go
config := config.DevelopmentConfig()
config.WatchConfig = true
config.OnConfigChange = func() {
    log.Println("配置文件已更改，重新加载配置...")
    // 在这里处理配置变更逻辑
}

manager, err := config.NewManager(config)
if err != nil {
    log.Fatal(err)
}

manager.Load()
```

## 远程配置

### etcd

```go
manager, err := config.WithRemote(
    "etcd",
    "http://127.0.0.1:2379",
    "/config/myapp",
    config.YAMLFormat,
)
```

### Consul

```go
manager, err := config.WithRemote(
    "consul",
    "127.0.0.1:8500",
    "config/myapp",
    config.JSONFormat,
)
```

## 自定义配置

```go
customConfig := &config.Config{
    ConfigName: "myapp",
    ConfigType: config.YAMLFormat,
    ConfigPaths: []string{
        "./configs",
        "/etc/myapp",
        "$HOME/.myapp",
    },
    EnvPrefix:      "MYAPP",
    EnvKeyReplacer: "_",
    AutoEnv:        true,
    WatchConfig:    true,
    Debug:          true,
}

manager, err := config.NewManager(customConfig)
```

## 工厂模式

```go
factory := config.NewFactory()

// 根据环境变量创建
manager, err := factory.CreateWithEnv("APP_ENV")

// 使用指定配置文件创建
manager, err := factory.CreateWithFile("./configs/app.yaml")

// 创建自定义配置
manager, err := factory.CreateCustom(customConfig)
```

## 便捷函数

```go
// 创建并加载，失败时 panic
manager := config.LoadOrPanic(func() (*config.Manager, error) {
    return config.WithFile("./configs/app.yaml")
})

// 必须加载成功
manager := config.MustLoad(manager)
```

## API 参考

### 配置获取方法

```go
// 基本类型
manager.Get(key string) interface{}
manager.GetString(key string) string
manager.GetBool(key string) bool
manager.GetInt(key string) int
manager.GetInt32(key string) int32
manager.GetInt64(key string) int64
manager.GetUint(key string) uint
manager.GetUint32(key string) uint32
manager.GetUint64(key string) uint64
manager.GetFloat64(key string) float64

// 时间类型
manager.GetTime(key string) time.Time
manager.GetDuration(key string) time.Duration

// 切片类型
manager.GetIntSlice(key string) []int
manager.GetStringSlice(key string) []string

// 映射类型
manager.GetStringMap(key string) map[string]interface{}
manager.GetStringMapString(key string) map[string]string
manager.GetStringMapStringSlice(key string) map[string][]string

// 大小类型
manager.GetSizeInBytes(key string) uint
```

### 配置设置方法

```go
// 设置配置值
manager.Set(key string, value interface{})

// 设置默认值
manager.SetDefault(key string, value interface{})

// 检查配置是否存在
manager.IsSet(key string) bool

// 获取所有配置键
manager.AllKeys() []string

// 获取所有配置
manager.AllSettings() map[string]interface{}
```

### 结构体绑定方法

```go
// 解析到结构体
manager.Unmarshal(rawVal interface{}) error

// 解析指定键到结构体
manager.UnmarshalKey(key string, rawVal interface{}) error

// 精确解析（不允许未知字段）
manager.UnmarshalExact(rawVal interface{}) error
```

### 环境变量绑定

```go
// 绑定环境变量
manager.BindEnv(input ...string) error

// 绑定命令行参数
manager.BindPFlag(key string, flag *pflag.Flag) error
manager.BindPFlags(flags *pflag.FlagSet) error
```

## 最佳实践

### 1. 配置结构设计

```go
// 推荐：按功能模块组织配置
type Config struct {
    App      AppConfig      `mapstructure:"app"`
    Server   ServerConfig   `mapstructure:"server"`
    Database DatabaseConfig `mapstructure:"database"`
    Redis    RedisConfig    `mapstructure:"redis"`
    Logger   LoggerConfig   `mapstructure:"logger"`
}
```

### 2. 环境变量命名

```bash
# 推荐：使用统一的前缀和分隔符
APP_SERVER_HOST=localhost
APP_SERVER_PORT=8080
APP_DATABASE_DSN="user:pass@tcp(localhost:3306)/db"
```

### 3. 配置文件组织

```
configs/
├── config.yaml          # 默认配置
├── config.dev.yaml      # 开发环境配置
├── config.prod.yaml     # 生产环境配置
└── config.test.yaml     # 测试环境配置
```

### 4. 敏感信息处理

```go
// 敏感信息通过环境变量传递
manager.BindEnv("database.password", "DB_PASSWORD")
manager.BindEnv("auth.jwt_secret", "JWT_SECRET")
```

### 5. 配置验证

```go
type Config struct {
    Server ServerConfig `mapstructure:"server" validate:"required"`
    Database DatabaseConfig `mapstructure:"database" validate:"required"`
}

func (c *Config) Validate() error {
    validate := validator.New()
    return validate.Struct(c)
}
```

## 错误处理

```go
manager, err := config.Default()
if err != nil {
    log.Fatalf("创建配置管理器失败: %v", err)
}

if err := manager.Load(); err != nil {
    log.Fatalf("加载配置失败: %v", err)
}

// 检查配置是否加载
if !manager.IsLoaded() {
    log.Fatal("配置未加载")
}
```

## 调试

```go
// 启用调试模式
config := config.DefaultConfig()
config.Debug = true

manager, _ := config.NewManager(config)
manager.Load()

// 打印调试信息
manager.Debug()
```

## 注意事项

1. **配置文件优先级**: 环境变量 > 配置文件 > 默认值
2. **键名大小写**: 默认不区分大小写，可通过 `CaseSensitive` 配置
3. **环境变量替换**: 配置键中的 `.` 和 `-` 会被替换为 `_`
4. **配置监听**: 仅在开发环境建议启用，生产环境建议禁用
5. **远程配置**: 需要额外安装对应的驱动包

## 许可证

本包基于 MIT 许可证开源。