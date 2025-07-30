package config

import (
	"fmt"
	"log"
	"time"
)

// AppConfig 应用配置结构体示例
type AppConfig struct {
	// 服务器配置
	Server ServerConfig `mapstructure:"server" json:"server" yaml:"server"`
	// 数据库配置
	Database DatabaseConfig `mapstructure:"database" json:"database" yaml:"database"`
	// Redis配置
	Redis RedisConfig `mapstructure:"redis" json:"redis" yaml:"redis"`
	// 日志配置
	Logger LoggerConfig `mapstructure:"logger" json:"logger" yaml:"logger"`
	// 认证配置
	Auth AuthConfig `mapstructure:"auth" json:"auth" yaml:"auth"`
}

// ServerConfig 服务器配置
type ServerConfig struct {
	Host         string        `mapstructure:"host" json:"host" yaml:"host"`
	Port         int           `mapstructure:"port" json:"port" yaml:"port"`
	Mode         string        `mapstructure:"mode" json:"mode" yaml:"mode"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
	IdleTimeout  time.Duration `mapstructure:"idle_timeout" json:"idle_timeout" yaml:"idle_timeout"`
	EnableTLS    bool          `mapstructure:"enable_tls" json:"enable_tls" yaml:"enable_tls"`
	TLSCertFile  string        `mapstructure:"tls_cert_file" json:"tls_cert_file" yaml:"tls_cert_file"`
	TLSKeyFile   string        `mapstructure:"tls_key_file" json:"tls_key_file" yaml:"tls_key_file"`
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Driver          string        `mapstructure:"driver" json:"driver" yaml:"driver"`
	DSN             string        `mapstructure:"dsn" json:"dsn" yaml:"dsn"`
	MaxIdleConns    int           `mapstructure:"max_idle_conns" json:"max_idle_conns" yaml:"max_idle_conns"`
	MaxOpenConns    int           `mapstructure:"max_open_conns" json:"max_open_conns" yaml:"max_open_conns"`
	ConnMaxLifetime time.Duration `mapstructure:"conn_max_lifetime" json:"conn_max_lifetime" yaml:"conn_max_lifetime"`
	SlowThreshold   time.Duration `mapstructure:"slow_threshold" json:"slow_threshold" yaml:"slow_threshold"`
	LogLevel        string        `mapstructure:"log_level" json:"log_level" yaml:"log_level"`
}

// RedisConfig Redis配置
type RedisConfig struct {
	Addr         string        `mapstructure:"addr" json:"addr" yaml:"addr"`
	Password     string        `mapstructure:"password" json:"password" yaml:"password"`
	DB           int           `mapstructure:"db" json:"db" yaml:"db"`
	PoolSize     int           `mapstructure:"pool_size" json:"pool_size" yaml:"pool_size"`
	DialTimeout  time.Duration `mapstructure:"dial_timeout" json:"dial_timeout" yaml:"dial_timeout"`
	ReadTimeout  time.Duration `mapstructure:"read_timeout" json:"read_timeout" yaml:"read_timeout"`
	WriteTimeout time.Duration `mapstructure:"write_timeout" json:"write_timeout" yaml:"write_timeout"`
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level      string `mapstructure:"level" json:"level" yaml:"level"`
	Format     string `mapstructure:"format" json:"format" yaml:"format"`
	Output     string `mapstructure:"output" json:"output" yaml:"output"`
	FilePath   string `mapstructure:"file_path" json:"file_path" yaml:"file_path"`
	MaxSize    int    `mapstructure:"max_size" json:"max_size" yaml:"max_size"`
	MaxBackups int    `mapstructure:"max_backups" json:"max_backups" yaml:"max_backups"`
	MaxAge     int    `mapstructure:"max_age" json:"max_age" yaml:"max_age"`
	Compress   bool   `mapstructure:"compress" json:"compress" yaml:"compress"`
}

// AuthConfig 认证配置
type AuthConfig struct {
	JWTSecret          string        `mapstructure:"jwt_secret" json:"jwt_secret" yaml:"jwt_secret"`
	AccessTokenExpire  time.Duration `mapstructure:"access_token_expire" json:"access_token_expire" yaml:"access_token_expire"`
	RefreshTokenExpire time.Duration `mapstructure:"refresh_token_expire" json:"refresh_token_expire" yaml:"refresh_token_expire"`
	TokenIssuer        string        `mapstructure:"token_issuer" json:"token_issuer" yaml:"token_issuer"`
}

// ExampleBasicUsage 基本使用示例
func ExampleBasicUsage() {
	// 创建默认配置管理器
	manager, err := Default()
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	// 加载配置
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 获取配置值
	serverHost := manager.GetString("server.host")
	serverPort := manager.GetInt("server.port")
	debugMode := manager.GetBool("debug")

	fmt.Printf("Server: %s:%d, Debug: %t\n", serverHost, serverPort, debugMode)
}

// ExampleStructBinding 结构体绑定示例
func ExampleStructBinding() {
	// 创建开发环境配置管理器
	manager, err := Development()
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	// 加载配置
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 设置默认值
	manager.SetDefault("server.host", "localhost")
	manager.SetDefault("server.port", 8080)
	manager.SetDefault("server.mode", "debug")
	manager.SetDefault("database.driver", "mysql")
	manager.SetDefault("redis.addr", "localhost:6379")

	// 解析到结构体
	var appConfig AppConfig
	if err := manager.Unmarshal(&appConfig); err != nil {
		log.Fatalf("Failed to unmarshal config: %v", err)
	}

	fmt.Printf("App Config: %+v\n", appConfig)
}

// ExampleEnvironmentVariables 环境变量示例
func ExampleEnvironmentVariables() {
	// 创建配置管理器
	config := DefaultConfig()
	config.EnvPrefix = "MYAPP"
	config.AutoEnv = true

	manager, err := NewManager(config)
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	// 加载配置
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 绑定特定环境变量
	manager.BindEnv("database.password", "MYAPP_DB_PASSWORD")
	manager.BindEnv("redis.password", "MYAPP_REDIS_PASSWORD")

	// 获取环境变量值
	dbPassword := manager.GetString("database.password")
	redisPassword := manager.GetString("redis.password")

	fmt.Printf("DB Password: %s, Redis Password: %s\n", dbPassword, redisPassword)
}

// ExampleConfigWatch 配置监听示例
func ExampleConfigWatch() {
	// 创建配置管理器并启用监听
	config := DevelopmentConfig()
	config.WatchConfig = true
	config.OnConfigChange = func() {
		fmt.Println("Config file changed!")
		// 重新加载应用配置
		// reloadAppConfig()
	}

	manager, err := NewManager(config)
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	// 加载配置
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Println("Config loaded with file watching enabled")
	// 应用会持续运行并监听配置文件变化
}

// ExampleRemoteConfig 远程配置示例
func ExampleRemoteConfig() {
	// 创建远程配置管理器
	manager, err := WithRemote("etcd", "http://127.0.0.1:2379", "/config/myapp", YAMLFormat)
	if err != nil {
		log.Fatalf("Failed to create remote config manager: %v", err)
	}

	// 加载远程配置
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load remote config: %v", err)
	}

	// 获取配置值
	appName := manager.GetString("app.name")
	version := manager.GetString("app.version")

	fmt.Printf("App: %s, Version: %s\n", appName, version)
}

// ExampleCustomConfig 自定义配置示例
func ExampleCustomConfig() {
	// 创建自定义配置
	config := &Config{
		ConfigName: "myapp",
		ConfigType: YAMLFormat,
		ConfigPaths: []string{
			"./configs",
			"/etc/myapp",
		},
		EnvPrefix:      "MYAPP",
		EnvKeyReplacer: "_",
		AutoEnv:        true,
		WatchConfig:    true,
		Debug:          true,
	}

	// 创建管理器
	manager, err := NewManager(config)
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	// 加载配置
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	// 打印调试信息
	manager.Debug()
}

// ExampleFactoryUsage 工厂模式使用示例
func ExampleFactoryUsage() {
	// 使用工厂创建不同环境的配置管理器
	factory := NewFactory()

	// 根据环境变量创建
	manager, err := factory.CreateWithEnv("APP_ENV")
	if err != nil {
		log.Fatalf("Failed to create config manager: %v", err)
	}

	// 加载配置
	if err := manager.Load(); err != nil {
		log.Fatalf("Failed to load config: %v", err)
	}

	fmt.Printf("Config loaded for environment: %s\n", manager.GetString("env"))
}

// ExamplePanicOnError 错误时panic示例
func ExamplePanicOnError() {
	// 创建并加载配置，失败时panic
	manager := LoadOrPanic(func() (*Manager, error) {
		return WithFile("./configs/app.yaml")
	})

	// 配置已成功加载，可以直接使用
	appName := manager.GetString("app.name")
	fmt.Printf("App Name: %s\n", appName)
}
