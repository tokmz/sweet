package config

import (
	"fmt"
	"path/filepath"
	"strings"
	"time"
)

// Format 配置文件格式
type Format string

const (
	// YAMLFormat YAML格式
	YAMLFormat Format = "yaml"
	// JSONFormat JSON格式
	JSONFormat Format = "json"
	// TOMLFormat TOML格式
	TOMLFormat Format = "toml"
	// PropertiesFormat Properties格式
	PropertiesFormat Format = "properties"
	// INIFormat INI格式
	INIFormat Format = "ini"
)

// Config Viper配置管理器配置
type Config struct {
	// 基本配置
	ConfigName  string   // 配置文件名(不含扩展名)
	ConfigType  Format   // 配置文件格式
	ConfigPaths []string // 配置文件搜索路径列表

	// 环境变量配置
	EnvPrefix      string // 环境变量前缀
	EnvKeyReplacer string // 环境变量键替换符(默认为"_")
	AutoEnv        bool   // 是否自动读取环境变量
	AllowEmptyEnv  bool   // 是否允许空环境变量

	// 远程配置
	RemoteProvider   string // 远程配置提供者: etcd, consul, firestore
	RemoteEndpoint   string // 远程配置端点
	RemoteConfigPath string // 远程配置路径
	RemoteConfigType Format // 远程配置格式

	// 监听配置
	WatchConfig    bool          // 是否监听配置文件变化
	WatchInterval  time.Duration // 监听间隔(仅对某些提供者有效)
	OnConfigChange func()        // 配置变化回调函数

	// 高级配置
	CaseSensitive bool // 是否区分大小写
	Debug         bool // 是否启用调试模式
}

// DefaultConfig 返回默认配置
func DefaultConfig() *Config {
	return &Config{
		ConfigName: "config",
		ConfigType: YAMLFormat,
		ConfigPaths: []string{
			".",
			"./configs",
			"./config",
			"/etc/app",
			"$HOME/.app",
		},
		EnvPrefix:      "APP",
		EnvKeyReplacer: "_",
		AutoEnv:        true,
		AllowEmptyEnv:  false,
		WatchConfig:    false,
		WatchInterval:  time.Second * 5,
		CaseSensitive:  false,
		Debug:          false,
	}
}

// DevelopmentConfig 返回开发环境配置
func DevelopmentConfig() *Config {
	config := DefaultConfig()
	config.ConfigPaths = []string{
		".",
		"./configs",
		"./config",
	}
	config.WatchConfig = true
	config.Debug = true
	return config
}

// ProductionConfig 返回生产环境配置
func ProductionConfig() *Config {
	config := DefaultConfig()
	config.ConfigPaths = []string{
		"/etc/app",
		"./configs",
	}
	config.WatchConfig = false
	config.Debug = false
	return config
}

// Validate 验证配置
func (c *Config) Validate() error {
	if c.ConfigName == "" {
		return fmt.Errorf("config name cannot be empty")
	}

	if len(c.ConfigPaths) == 0 {
		return fmt.Errorf("config paths cannot be empty")
	}

	validFormats := []Format{YAMLFormat, JSONFormat, TOMLFormat, PropertiesFormat, INIFormat}
	validFormat := false
	for _, format := range validFormats {
		if c.ConfigType == format {
			validFormat = true
			break
		}
	}
	if !validFormat {
		return fmt.Errorf("unsupported config type: %s", c.ConfigType)
	}

	if c.RemoteProvider != "" {
		validProviders := []string{"etcd", "consul", "firestore"}
		validProvider := false
		for _, provider := range validProviders {
			if c.RemoteProvider == provider {
				validProvider = true
				break
			}
		}
		if !validProvider {
			return fmt.Errorf("unsupported remote provider: %s", c.RemoteProvider)
		}

		if c.RemoteEndpoint == "" {
			return fmt.Errorf("remote endpoint cannot be empty when using remote provider")
		}

		if c.RemoteConfigPath == "" {
			return fmt.Errorf("remote config path cannot be empty when using remote provider")
		}
	}

	return nil
}

// GetConfigFileName 获取完整的配置文件名
func (c *Config) GetConfigFileName() string {
	return fmt.Sprintf("%s.%s", c.ConfigName, c.ConfigType)
}

// GetEnvKey 获取环境变量键名
func (c *Config) GetEnvKey(key string) string {
	if c.EnvPrefix == "" {
		return strings.ToUpper(key)
	}
	replacer := c.EnvKeyReplacer
	if replacer == "" {
		replacer = "_"
	}
	envKey := strings.ReplaceAll(key, ".", replacer)
	envKey = strings.ReplaceAll(envKey, "-", replacer)
	return strings.ToUpper(fmt.Sprintf("%s%s%s", c.EnvPrefix, replacer, envKey))
}

// GetSearchPaths 获取配置文件搜索路径
func (c *Config) GetSearchPaths() []string {
	paths := make([]string, 0, len(c.ConfigPaths))
	for _, path := range c.ConfigPaths {
		// 展开环境变量
		expandedPath := filepath.Clean(path)
		paths = append(paths, expandedPath)
	}
	return paths
}
