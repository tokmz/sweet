# 配置包 (pkg/config)

配置包提供了一个简单高效的配置管理解决方案，基于 [Viper](https://github.com/spf13/viper) 库实现。该包支持从多种格式的配置文件中加载配置、热重载、环境变量替换等功能。

## 特性

- 支持多种配置文件格式：YAML、JSON、TOML、HCL、INI 和环境变量
- 允许从多个路径搜索配置文件
- 支持配置热重载（实时监测配置文件变更）
- 自动绑定环境变量
- 简洁的错误处理

## 安装

确保您的项目中已经引入了必要的依赖：

```bash
go get github.com/spf13/viper
go get github.com/fsnotify/fsnotify
```

## 使用方法

### 初始化配置

使用 `InitConfig` 函数初始化配置：

```go
import "nevus/pkg/config"

func initApp() error {
    cfg := &config.Config{
        Name: "config",         // 配置文件名（不含扩展名）
        Type: "yaml",           // 配置文件类型
        Path: []string{         // 搜索路径，按顺序搜索
            "./config",
            "./resource",
            "/etc/nevus",
        },
    }
    
    // 初始化配置
    if err := config.InitConfig(cfg); err != nil {
        return err
    }
    
    return nil
}
```

### 加载配置到结构体

使用 `LoadConfig` 将配置加载到您的结构体中：

```go
type AppConfig struct {
    Server struct {
        Port int    `mapstructure:"port"`
        Host string `mapstructure:"host"`
    } `mapstructure:"server"`
    
    Database struct {
        DSN string `mapstructure:"dsn"`
    } `mapstructure:"database"`
}

func loadAppConfig() (*AppConfig, error) {
    cfg := &config.Config{
        Name: "config",
        Type: "yaml",
        Path: []string{"./config"},
    }
    
    var appConfig AppConfig
    if err := config.LoadConfig(cfg, &appConfig); err != nil {
        return nil, err
    }
    
    return &appConfig, nil
}
```

### 从特定文件加载配置

如果您知道配置文件的确切路径，可以使用 `LoadConfigFile`：

```go
func loadFromFile() (*AppConfig, error) {
    var appConfig AppConfig
    if err := config.LoadConfigFile("./resource/config.yaml", &appConfig); err != nil {
        return nil, err
    }
    
    return &appConfig, nil
}
```

### 监控配置文件变更

使用 `WatchConfig` 监控配置文件变更，并在变更时执行回调：

```go
func watchConfigChanges() error {
    return config.WatchConfig(func() {
        fmt.Println("配置文件已更新")
        
        // 重新加载配置
        var newConfig AppConfig
        if err := config.Viper.Unmarshal(&newConfig); err != nil {
            fmt.Printf("重新加载配置失败: %v\n", err)
            return
        }
        
        // 应用新配置
        applyNewConfig(newConfig)
    })
}
```

## 环境变量替换

配置包支持使用环境变量替换配置值。例如，如果您的配置文件有以下内容：

```yaml
database:
  password: ${DB_PASSWORD}
```

您可以设置环境变量 `DB_PASSWORD`，它将被自动替换。环境变量名称使用 `_` 替换配置键中的 `.`，例如 `database.password` 对应的环境变量是 `DATABASE_PASSWORD`。

## 完整示例

以下是一个完整的示例，展示如何使用配置包：

```go
package main

import (
    "fmt"
    "log"
    
    "nevus/pkg/config"
)

type ServerConfig struct {
    Server struct {
        Name string `mapstructure:"name"`
        Port int    `mapstructure:"port"`
        Host string `mapstructure:"host"`
    } `mapstructure:"server"`
    
    Database struct {
        Driver   string `mapstructure:"driver"`
        Host     string `mapstructure:"host"`
        Port     int    `mapstructure:"port"`
        Username string `mapstructure:"username"`
        Password string `mapstructure:"password"`
        Database string `mapstructure:"database"`
    } `mapstructure:"database"`
}

func main() {
    // 初始化配置
    cfg := &config.Config{
        Name: "config",
        Type: "yaml",
        Path: []string{"./config", "./resource"},
    }
    
    if err := config.InitConfig(cfg); err != nil {
        log.Fatalf("初始化配置失败: %v", err)
    }
    
    // 加载配置到结构体
    var serverConfig ServerConfig
    if err := config.LoadConfig(cfg, &serverConfig); err != nil {
        log.Fatalf("加载配置失败: %v", err)
    }
    
    // 使用配置
    fmt.Printf("服务器名称: %s\n", serverConfig.Server.Name)
    fmt.Printf("端口: %d\n", serverConfig.Server.Port)
    fmt.Printf("数据库连接: %s@%s:%d/%s\n",
        serverConfig.Database.Username,
        serverConfig.Database.Host,
        serverConfig.Database.Port,
        serverConfig.Database.Database,
    )
    
    // 监控配置文件变更
    if err := config.WatchConfig(func() {
        fmt.Println("配置文件已更新")
        
        // 重新加载配置
        var newConfig ServerConfig
        if err := config.Viper.Unmarshal(&newConfig); err != nil {
            fmt.Printf("重新加载配置失败: %v\n", err)
            return
        }
        
        // 应用新配置
        fmt.Printf("新端口: %d\n", newConfig.Server.Port)
    }); err != nil {
        log.Printf("监控配置文件失败: %v", err)
    }
    
    // 保持程序运行
    select {}
}
```

## 错误处理

包中定义了以下错误类型：

- `ErrConfigNotFound` - 配置文件未找到
- `ErrConfigFileType` - 配置文件类型错误
- `ErrInvalidConfigType` - 无效的配置文件类型
- `ErrReadConfig` - 读取配置文件失败
- `ErrConfigValidation` - 配置验证失败

使用这些错误常量可以进行更精确的错误处理：

```go
if err := config.InitConfig(cfg); err != nil {
    if errors.Is(err, config.ErrConfigNotFound) {
        log.Fatal("找不到配置文件，请确认配置文件存在")
    } else if errors.Is(err, config.ErrInvalidConfigType) {
        log.Fatal("配置文件类型无效，请使用支持的格式（yaml、json、toml等）")
    } else {
        log.Fatalf("初始化配置时发生错误: %v", err)
    }
}
```

## 与应用配置系统的集成

本配置包可以与应用程序的内部配置系统无缝集成，例如加载 `internal/types.Config` 类型的配置：

```go
import (
    "nevus/internal/types"
    "nevus/pkg/config"
)

func LoadAppConfig() (*types.Config, error) {
    cfg := &config.Config{
        Name: "config",
        Type: "yaml",
        Path: []string{"./resource"},
    }
    
    var appConfig types.Config
    if err := config.LoadConfig(cfg, &appConfig); err != nil {
        return nil, err
    }
    
    // 全局配置变量
    types.Cfg = &appConfig
    
    return &appConfig, nil
} 