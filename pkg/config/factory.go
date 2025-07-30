package config

import (
	"fmt"
	"os"
	"path/filepath"
)

// Factory 配置管理器工厂
type Factory struct{}

// NewFactory 创建新的工厂实例
func NewFactory() *Factory {
	return &Factory{}
}

// CreateDefault 创建默认配置管理器
func (f *Factory) CreateDefault() (*Manager, error) {
	return NewManager(DefaultConfig())
}

// CreateDevelopment 创建开发环境配置管理器
func (f *Factory) CreateDevelopment() (*Manager, error) {
	return NewManager(DevelopmentConfig())
}

// CreateProduction 创建生产环境配置管理器
func (f *Factory) CreateProduction() (*Manager, error) {
	return NewManager(ProductionConfig())
}

// CreateWithFile 使用指定配置文件创建管理器
func (f *Factory) CreateWithFile(configPath string) (*Manager, error) {
	if configPath == "" {
		return nil, fmt.Errorf("config path cannot be empty")
	}

	// 检查文件是否存在
	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("config file does not exist: %s", configPath)
	}

	// 解析文件路径
	dir := filepath.Dir(configPath)
	filename := filepath.Base(configPath)
	ext := filepath.Ext(filename)
	name := filename[:len(filename)-len(ext)]

	// 确定配置格式
	var format Format
	switch ext {
	case ".yaml", ".yml":
		format = YAMLFormat
	case ".json":
		format = JSONFormat
	case ".toml":
		format = TOMLFormat
	case ".properties":
		format = PropertiesFormat
	case ".ini":
		format = INIFormat
	default:
		return nil, fmt.Errorf("unsupported config file format: %s", ext)
	}

	config := &Config{
		ConfigName:     name,
		ConfigType:     format,
		ConfigPaths:    []string{dir},
		EnvPrefix:      "APP",
		EnvKeyReplacer: "_",
		AutoEnv:        true,
		AllowEmptyEnv:  false,
		WatchConfig:    false,
		CaseSensitive:  false,
		Debug:          false,
	}

	return NewManager(config)
}

// CreateWithEnv 根据环境变量创建管理器
func (f *Factory) CreateWithEnv(envKey string) (*Manager, error) {
	if envKey == "" {
		envKey = "APP_ENV"
	}

	env := os.Getenv(envKey)
	switch env {
	case "development", "dev":
		return f.CreateDevelopment()
	case "production", "prod":
		return f.CreateProduction()
	case "test", "testing":
		return f.CreateTesting()
	default:
		return f.CreateDefault()
	}
}

// CreateTesting 创建测试环境配置管理器
func (f *Factory) CreateTesting() (*Manager, error) {
	config := DefaultConfig()
	config.ConfigPaths = []string{
		".",
		"./testdata",
		"./configs",
	}
	config.WatchConfig = false
	config.Debug = true
	return NewManager(config)
}

// CreateWithRemote 创建远程配置管理器
func (f *Factory) CreateWithRemote(provider, endpoint, configPath string, configType Format) (*Manager, error) {
	if provider == "" {
		return nil, fmt.Errorf("remote provider cannot be empty")
	}
	if endpoint == "" {
		return nil, fmt.Errorf("remote endpoint cannot be empty")
	}
	if configPath == "" {
		return nil, fmt.Errorf("remote config path cannot be empty")
	}

	config := DefaultConfig()
	config.RemoteProvider = provider
	config.RemoteEndpoint = endpoint
	config.RemoteConfigPath = configPath
	config.RemoteConfigType = configType

	return NewManager(config)
}

// CreateCustom 创建自定义配置管理器
func (f *Factory) CreateCustom(config *Config) (*Manager, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}
	return NewManager(config)
}

// 全局工厂实例
var defaultFactory = NewFactory()

// 便捷函数

// Default 创建默认配置管理器
func Default() (*Manager, error) {
	return defaultFactory.CreateDefault()
}

// Development 创建开发环境配置管理器
func Development() (*Manager, error) {
	return defaultFactory.CreateDevelopment()
}

// Production 创建生产环境配置管理器
func Production() (*Manager, error) {
	return defaultFactory.CreateProduction()
}

// WithFile 使用指定配置文件创建管理器
func WithFile(configPath string) (*Manager, error) {
	return defaultFactory.CreateWithFile(configPath)
}

// WithEnv 根据环境变量创建管理器
func WithEnv(envKey string) (*Manager, error) {
	return defaultFactory.CreateWithEnv(envKey)
}

// Testing 创建测试环境配置管理器
func Testing() (*Manager, error) {
	return defaultFactory.CreateTesting()
}

// WithRemote 创建远程配置管理器
func WithRemote(provider, endpoint, configPath string, configType Format) (*Manager, error) {
	return defaultFactory.CreateWithRemote(provider, endpoint, configPath, configType)
}

// Custom 创建自定义配置管理器
func Custom(config *Config) (*Manager, error) {
	return defaultFactory.CreateCustom(config)
}

// MustLoad 加载配置管理器，失败时panic
func MustLoad(manager *Manager) *Manager {
	if err := manager.Load(); err != nil {
		panic(fmt.Sprintf("failed to load config: %v", err))
	}
	return manager
}

// LoadOrPanic 创建并加载配置管理器，失败时panic
func LoadOrPanic(createFunc func() (*Manager, error)) *Manager {
	manager, err := createFunc()
	if err != nil {
		panic(fmt.Sprintf("failed to create config manager: %v", err))
	}
	return MustLoad(manager)
}
