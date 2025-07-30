package config

import (
	"os"
	"path/filepath"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestConfig_Validate(t *testing.T) {
	tests := []struct {
		name    string
		config  *Config
		wantErr bool
		errMsg  string
	}{
		{
			name:    "valid default config",
			config:  DefaultConfig(),
			wantErr: false,
		},
		{
			name: "empty config name",
			config: &Config{
				ConfigName:  "",
				ConfigType:  YAMLFormat,
				ConfigPaths: []string{"."},
			},
			wantErr: true,
			errMsg:  "config name cannot be empty",
		},
		{
			name: "empty config paths",
			config: &Config{
				ConfigName:  "test",
				ConfigType:  YAMLFormat,
				ConfigPaths: []string{},
			},
			wantErr: true,
			errMsg:  "config paths cannot be empty",
		},
		{
			name: "invalid config type",
			config: &Config{
				ConfigName:  "test",
				ConfigType:  Format("invalid"),
				ConfigPaths: []string{"."},
			},
			wantErr: true,
			errMsg:  "unsupported config type: invalid",
		},
		{
			name: "invalid remote provider",
			config: &Config{
				ConfigName:       "test",
				ConfigType:       YAMLFormat,
				ConfigPaths:      []string{"."},
				RemoteProvider:   "invalid",
				RemoteEndpoint:   "http://localhost",
				RemoteConfigPath: "/config",
			},
			wantErr: true,
			errMsg:  "unsupported remote provider: invalid",
		},
		{
			name: "missing remote endpoint",
			config: &Config{
				ConfigName:       "test",
				ConfigType:       YAMLFormat,
				ConfigPaths:      []string{"."},
				RemoteProvider:   "etcd",
				RemoteEndpoint:   "",
				RemoteConfigPath: "/config",
			},
			wantErr: true,
			errMsg:  "remote endpoint cannot be empty when using remote provider",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.config.Validate()
			if tt.wantErr {
				require.Error(t, err)
				assert.Contains(t, err.Error(), tt.errMsg)
			} else {
				require.NoError(t, err)
			}
		})
	}
}

func TestConfig_GetConfigFileName(t *testing.T) {
	config := &Config{
		ConfigName: "app",
		ConfigType: YAMLFormat,
	}

	filename := config.GetConfigFileName()
	assert.Equal(t, "app.yaml", filename)
}

func TestConfig_GetEnvKey(t *testing.T) {
	tests := []struct {
		name     string
		config   *Config
		key      string
		expected string
	}{
		{
			name: "with prefix",
			config: &Config{
				EnvPrefix:      "APP",
				EnvKeyReplacer: "_",
			},
			key:      "server.host",
			expected: "APP_SERVER_HOST",
		},
		{
			name: "without prefix",
			config: &Config{
				EnvPrefix:      "",
				EnvKeyReplacer: "_",
			},
			key:      "server.host",
			expected: "SERVER.HOST",
		},
		{
			name: "with dash",
			config: &Config{
				EnvPrefix:      "APP",
				EnvKeyReplacer: "_",
			},
			key:      "log-level",
			expected: "APP_LOG_LEVEL",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := tt.config.GetEnvKey(tt.key)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestConfig_GetSearchPaths(t *testing.T) {
	config := &Config{
		ConfigPaths: []string{
			".",
			"./configs",
			"/etc/app",
		},
	}

	paths := config.GetSearchPaths()
	expected := []string{".", "configs", "/etc/app"}
	assert.Equal(t, expected, paths)
}

func TestDefaultConfigs(t *testing.T) {
	t.Run("DefaultConfig", func(t *testing.T) {
		config := DefaultConfig()
		assert.Equal(t, "config", config.ConfigName)
		assert.Equal(t, YAMLFormat, config.ConfigType)
		assert.Equal(t, "APP", config.EnvPrefix)
		assert.True(t, config.AutoEnv)
		assert.False(t, config.WatchConfig)
		assert.False(t, config.Debug)
	})

	t.Run("DevelopmentConfig", func(t *testing.T) {
		config := DevelopmentConfig()
		assert.True(t, config.WatchConfig)
		assert.True(t, config.Debug)
		assert.Contains(t, config.ConfigPaths, ".")
		assert.Contains(t, config.ConfigPaths, "./configs")
	})

	t.Run("ProductionConfig", func(t *testing.T) {
		config := ProductionConfig()
		assert.False(t, config.WatchConfig)
		assert.False(t, config.Debug)
		assert.Contains(t, config.ConfigPaths, "/etc/app")
	})
}

func TestManager_Creation(t *testing.T) {
	t.Run("NewManager with valid config", func(t *testing.T) {
		config := DefaultConfig()
		manager, err := NewManager(config)
		require.NoError(t, err)
		assert.NotNil(t, manager)
		assert.Equal(t, config, manager.GetConfig())
		assert.False(t, manager.IsLoaded())
	})

	t.Run("NewManager with nil config", func(t *testing.T) {
		manager, err := NewManager(nil)
		require.NoError(t, err)
		assert.NotNil(t, manager)
		assert.NotNil(t, manager.GetConfig())
	})

	t.Run("NewManager with invalid config", func(t *testing.T) {
		config := &Config{
			ConfigName: "", // invalid
		}
		manager, err := NewManager(config)
		require.Error(t, err)
		assert.Nil(t, manager)
	})
}

func TestManager_BasicOperations(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "test.yaml")
	configContent := `
server:
  host: localhost
  port: 8080
  debug: true
database:
  driver: mysql
  max_connections: 100
features:
  - auth
  - logging
  - monitoring
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 创建配置管理器
	config := &Config{
		ConfigName:  "test",
		ConfigType:  YAMLFormat,
		ConfigPaths: []string{tempDir},
		AutoEnv:     false,
		Debug:       true,
	}

	manager, err := NewManager(config)
	require.NoError(t, err)

	// 加载配置
	err = manager.Load()
	require.NoError(t, err)
	assert.True(t, manager.IsLoaded())

	// 测试基本类型获取
	assert.Equal(t, "localhost", manager.GetString("server.host"))
	assert.Equal(t, 8080, manager.GetInt("server.port"))
	assert.True(t, manager.GetBool("server.debug"))
	assert.Equal(t, "mysql", manager.GetString("database.driver"))
	assert.Equal(t, 100, manager.GetInt("database.max_connections"))

	// 测试切片获取
	features := manager.GetStringSlice("features")
	assert.Equal(t, []string{"auth", "logging", "monitoring"}, features)

	// 测试设置和默认值
	manager.Set("runtime.version", "1.0.0")
	assert.Equal(t, "1.0.0", manager.GetString("runtime.version"))

	manager.SetDefault("cache.ttl", "5m")
	assert.Equal(t, "5m", manager.GetString("cache.ttl"))

	// 测试IsSet
	assert.True(t, manager.IsSet("server.host"))
	assert.False(t, manager.IsSet("nonexistent.key"))

	// 测试AllKeys
	keys := manager.AllKeys()
	assert.Contains(t, keys, "server.host")
	assert.Contains(t, keys, "database.driver")
}

func TestManager_StructBinding(t *testing.T) {
	type ServerConfig struct {
		Host  string `mapstructure:"host"`
		Port  int    `mapstructure:"port"`
		Debug bool   `mapstructure:"debug"`
	}

	type DatabaseConfig struct {
		Driver         string `mapstructure:"driver"`
		MaxConnections int    `mapstructure:"max_connections"`
	}

	type AppConfig struct {
		Server   ServerConfig   `mapstructure:"server"`
		Database DatabaseConfig `mapstructure:"database"`
		Features []string       `mapstructure:"features"`
	}

	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "app.yaml")
	configContent := `
server:
  host: 0.0.0.0
  port: 9090
  debug: false
database:
  driver: postgres
  max_connections: 50
features:
  - api
  - web
`
	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 创建配置管理器
	config := &Config{
		ConfigName:  "app",
		ConfigType:  YAMLFormat,
		ConfigPaths: []string{tempDir},
		AutoEnv:     false,
	}

	manager, err := NewManager(config)
	require.NoError(t, err)

	err = manager.Load()
	require.NoError(t, err)

	// 测试完整结构体解析
	var appConfig AppConfig
	err = manager.Unmarshal(&appConfig)
	require.NoError(t, err)

	assert.Equal(t, "0.0.0.0", appConfig.Server.Host)
	assert.Equal(t, 9090, appConfig.Server.Port)
	assert.False(t, appConfig.Server.Debug)
	assert.Equal(t, "postgres", appConfig.Database.Driver)
	assert.Equal(t, 50, appConfig.Database.MaxConnections)
	assert.Equal(t, []string{"api", "web"}, appConfig.Features)

	// 测试部分结构体解析
	var serverConfig ServerConfig
	err = manager.UnmarshalKey("server", &serverConfig)
	require.NoError(t, err)

	assert.Equal(t, "0.0.0.0", serverConfig.Host)
	assert.Equal(t, 9090, serverConfig.Port)
	assert.False(t, serverConfig.Debug)
}

func TestManager_EnvironmentVariables(t *testing.T) {
	// 设置环境变量
	os.Setenv("TEST_SERVER_HOST", "env-host")
	os.Setenv("TEST_SERVER_PORT", "3000")
	os.Setenv("TEST_DEBUG", "true")
	defer func() {
		os.Unsetenv("TEST_SERVER_HOST")
		os.Unsetenv("TEST_SERVER_PORT")
		os.Unsetenv("TEST_DEBUG")
	}()

	// 创建配置管理器
	config := &Config{
		ConfigName:     "nonexistent", // 故意使用不存在的配置文件
		ConfigType:     YAMLFormat,
		ConfigPaths:    []string{"."},
		EnvPrefix:      "TEST",
		EnvKeyReplacer: "_",
		AutoEnv:        true,
	}

	manager, err := NewManager(config)
	require.NoError(t, err)

	err = manager.Load() // 即使配置文件不存在也应该成功
	require.NoError(t, err)

	// 测试环境变量读取
	assert.Equal(t, "env-host", manager.GetString("server.host"))
	assert.Equal(t, 3000, manager.GetInt("server.port"))
	assert.True(t, manager.GetBool("debug"))
}

func TestManager_TypeConversions(t *testing.T) {
	manager, err := NewManager(DefaultConfig())
	require.NoError(t, err)

	err = manager.Load()
	require.NoError(t, err)

	// 设置各种类型的值
	manager.Set("test.string", "hello")
	manager.Set("test.int", 42)
	manager.Set("test.int32", int32(32))
	manager.Set("test.int64", int64(64))
	manager.Set("test.uint", uint(100))
	manager.Set("test.uint32", uint32(132))
	manager.Set("test.uint64", uint64(164))
	manager.Set("test.float64", 3.14)
	manager.Set("test.bool", true)
	manager.Set("test.duration", "5m")
	manager.Set("test.time", "2023-01-01T00:00:00Z")
	manager.Set("test.size", "1MB")
	manager.Set("test.int_slice", []int{1, 2, 3})
	manager.Set("test.string_slice", []string{"a", "b", "c"})
	manager.Set("test.string_map", map[string]interface{}{"key": "value"})

	// 测试类型转换
	assert.Equal(t, "hello", manager.GetString("test.string"))
	assert.Equal(t, 42, manager.GetInt("test.int"))
	assert.Equal(t, int32(32), manager.GetInt32("test.int32"))
	assert.Equal(t, int64(64), manager.GetInt64("test.int64"))
	assert.Equal(t, uint(100), manager.GetUint("test.uint"))
	assert.Equal(t, uint32(132), manager.GetUint32("test.uint32"))
	assert.Equal(t, uint64(164), manager.GetUint64("test.uint64"))
	assert.Equal(t, 3.14, manager.GetFloat64("test.float64"))
	assert.True(t, manager.GetBool("test.bool"))
	assert.Equal(t, 5*time.Minute, manager.GetDuration("test.duration"))
	assert.Equal(t, uint(1024*1024), manager.GetSizeInBytes("test.size"))
	assert.Equal(t, []int{1, 2, 3}, manager.GetIntSlice("test.int_slice"))
	assert.Equal(t, []string{"a", "b", "c"}, manager.GetStringSlice("test.string_slice"))

	stringMap := manager.GetStringMap("test.string_map")
	assert.Equal(t, "value", stringMap["key"])
}

func TestManager_ConfigFileUsed(t *testing.T) {
	// 创建临时配置文件
	tempDir := t.TempDir()
	configFile := filepath.Join(tempDir, "config.yaml")
	configContent := "app:\n  name: test\n"

	err := os.WriteFile(configFile, []byte(configContent), 0644)
	require.NoError(t, err)

	// 创建配置管理器
	config := &Config{
		ConfigName:  "config",
		ConfigType:  YAMLFormat,
		ConfigPaths: []string{tempDir},
		AutoEnv:     false,
	}

	manager, err := NewManager(config)
	require.NoError(t, err)

	err = manager.Load()
	require.NoError(t, err)

	// 检查使用的配置文件路径
	usedFile := manager.ConfigFileUsed()
	assert.Equal(t, configFile, usedFile)
}

func TestManager_Debug(t *testing.T) {
	// 创建启用调试的配置管理器
	config := DefaultConfig()
	config.Debug = true

	manager, err := NewManager(config)
	require.NoError(t, err)

	err = manager.Load()
	require.NoError(t, err)

	// 调用Debug方法（主要测试不会panic）
	manager.Debug()

	// 测试禁用调试的情况
	config.Debug = false
	manager2, err := NewManager(config)
	require.NoError(t, err)

	err = manager2.Load()
	require.NoError(t, err)

	manager2.Debug() // 应该不输出任何内容
}
