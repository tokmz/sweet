package config

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestFactory_CreateDefault(t *testing.T) {
	factory := NewFactory()
	manager, err := factory.CreateDefault()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	config := manager.GetConfig()
	assert.Equal(t, "config", config.ConfigName)
	assert.Equal(t, YAMLFormat, config.ConfigType)
	assert.Equal(t, "APP", config.EnvPrefix)
	assert.False(t, config.WatchConfig)
	assert.False(t, config.Debug)
}

func TestFactory_CreateDevelopment(t *testing.T) {
	factory := NewFactory()
	manager, err := factory.CreateDevelopment()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	config := manager.GetConfig()
	assert.True(t, config.WatchConfig)
	assert.True(t, config.Debug)
	assert.Contains(t, config.ConfigPaths, ".")
	assert.Contains(t, config.ConfigPaths, "./configs")
}

func TestFactory_CreateProduction(t *testing.T) {
	factory := NewFactory()
	manager, err := factory.CreateProduction()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	config := manager.GetConfig()
	assert.False(t, config.WatchConfig)
	assert.False(t, config.Debug)
	assert.Contains(t, config.ConfigPaths, "/etc/app")
}

func TestFactory_CreateTesting(t *testing.T) {
	factory := NewFactory()
	manager, err := factory.CreateTesting()
	require.NoError(t, err)
	assert.NotNil(t, manager)

	config := manager.GetConfig()
	assert.False(t, config.WatchConfig)
	assert.True(t, config.Debug)
	assert.Contains(t, config.ConfigPaths, ".")
	assert.Contains(t, config.ConfigPaths, "./testdata")
	assert.Contains(t, config.ConfigPaths, "./configs")
}

func TestFactory_CreateWithFile(t *testing.T) {
	factory := NewFactory()

	t.Run("valid yaml file", func(t *testing.T) {
		// 创建临时配置文件
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "app.yaml")
		configContent := "app:\n  name: test\n"
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		manager, err := factory.CreateWithFile(configFile)
		require.NoError(t, err)
		assert.NotNil(t, manager)

		config := manager.GetConfig()
		assert.Equal(t, "app", config.ConfigName)
		assert.Equal(t, YAMLFormat, config.ConfigType)
		assert.Contains(t, config.ConfigPaths, tempDir)
	})

	t.Run("valid json file", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.json")
		configContent := `{"app": {"name": "test"}}`
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		manager, err := factory.CreateWithFile(configFile)
		require.NoError(t, err)
		assert.NotNil(t, manager)

		config := manager.GetConfig()
		assert.Equal(t, "config", config.ConfigName)
		assert.Equal(t, JSONFormat, config.ConfigType)
	})

	t.Run("valid toml file", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "settings.toml")
		configContent := "[app]\nname = \"test\"\n"
		err := os.WriteFile(configFile, []byte(configContent), 0644)
		require.NoError(t, err)

		manager, err := factory.CreateWithFile(configFile)
		require.NoError(t, err)
		assert.NotNil(t, manager)

		config := manager.GetConfig()
		assert.Equal(t, "settings", config.ConfigName)
		assert.Equal(t, TOMLFormat, config.ConfigType)
	})

	t.Run("empty config path", func(t *testing.T) {
		manager, err := factory.CreateWithFile("")
		require.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "config path cannot be empty")
	})

	t.Run("nonexistent file", func(t *testing.T) {
		manager, err := factory.CreateWithFile("/nonexistent/config.yaml")
		require.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "config file does not exist")
	})

	t.Run("unsupported format", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.xml")
		err := os.WriteFile(configFile, []byte("<config></config>"), 0644)
		require.NoError(t, err)

		manager, err := factory.CreateWithFile(configFile)
		require.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "unsupported config file format")
	})
}

func TestFactory_CreateWithEnv(t *testing.T) {
	factory := NewFactory()

	tests := []struct {
		name       string
		envValue   string
		expectedFn func() (*Manager, error)
	}{
		{
			name:       "development",
			envValue:   "development",
			expectedFn: factory.CreateDevelopment,
		},
		{
			name:       "dev",
			envValue:   "dev",
			expectedFn: factory.CreateDevelopment,
		},
		{
			name:       "production",
			envValue:   "production",
			expectedFn: factory.CreateProduction,
		},
		{
			name:       "prod",
			envValue:   "prod",
			expectedFn: factory.CreateProduction,
		},
		{
			name:       "test",
			envValue:   "test",
			expectedFn: factory.CreateTesting,
		},
		{
			name:       "testing",
			envValue:   "testing",
			expectedFn: factory.CreateTesting,
		},
		{
			name:       "default",
			envValue:   "unknown",
			expectedFn: factory.CreateDefault,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// 设置环境变量
			os.Setenv("TEST_ENV", tt.envValue)
			defer os.Unsetenv("TEST_ENV")

			manager, err := factory.CreateWithEnv("TEST_ENV")
			require.NoError(t, err)
			assert.NotNil(t, manager)

			// 创建期望的管理器进行比较
			expectedManager, err := tt.expectedFn()
			require.NoError(t, err)

			// 比较配置
			actualConfig := manager.GetConfig()
			expectedConfig := expectedManager.GetConfig()
			assert.Equal(t, expectedConfig.WatchConfig, actualConfig.WatchConfig)
			assert.Equal(t, expectedConfig.Debug, actualConfig.Debug)
		})
	}

	t.Run("empty env key uses default", func(t *testing.T) {
		os.Setenv("APP_ENV", "production")
		defer os.Unsetenv("APP_ENV")

		manager, err := factory.CreateWithEnv("")
		require.NoError(t, err)
		assert.NotNil(t, manager)

		// 应该使用 APP_ENV 环境变量
		config := manager.GetConfig()
		assert.False(t, config.WatchConfig) // production config
		assert.False(t, config.Debug)       // production config
	})
}

func TestFactory_CreateWithRemote(t *testing.T) {
	factory := NewFactory()

	t.Run("valid remote config", func(t *testing.T) {
		manager, err := factory.CreateWithRemote(
			"etcd",
			"http://127.0.0.1:2379",
			"/config/app",
			YAMLFormat,
		)
		require.NoError(t, err)
		assert.NotNil(t, manager)

		config := manager.GetConfig()
		assert.Equal(t, "etcd", config.RemoteProvider)
		assert.Equal(t, "http://127.0.0.1:2379", config.RemoteEndpoint)
		assert.Equal(t, "/config/app", config.RemoteConfigPath)
		assert.Equal(t, YAMLFormat, config.RemoteConfigType)
	})

	t.Run("empty provider", func(t *testing.T) {
		manager, err := factory.CreateWithRemote(
			"",
			"http://127.0.0.1:2379",
			"/config/app",
			YAMLFormat,
		)
		require.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "remote provider cannot be empty")
	})

	t.Run("empty endpoint", func(t *testing.T) {
		manager, err := factory.CreateWithRemote(
			"etcd",
			"",
			"/config/app",
			YAMLFormat,
		)
		require.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "remote endpoint cannot be empty")
	})

	t.Run("empty config path", func(t *testing.T) {
		manager, err := factory.CreateWithRemote(
			"etcd",
			"http://127.0.0.1:2379",
			"",
			YAMLFormat,
		)
		require.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "remote config path cannot be empty")
	})
}

func TestFactory_CreateCustom(t *testing.T) {
	factory := NewFactory()

	t.Run("valid custom config", func(t *testing.T) {
		customConfig := &Config{
			ConfigName:  "custom",
			ConfigType:  JSONFormat,
			ConfigPaths: []string{"/custom/path"},
			EnvPrefix:   "CUSTOM",
			AutoEnv:     true,
			Debug:       true,
		}

		manager, err := factory.CreateCustom(customConfig)
		require.NoError(t, err)
		assert.NotNil(t, manager)
		assert.Equal(t, customConfig, manager.GetConfig())
	})

	t.Run("nil config", func(t *testing.T) {
		manager, err := factory.CreateCustom(nil)
		require.Error(t, err)
		assert.Nil(t, manager)
		assert.Contains(t, err.Error(), "config cannot be nil")
	})
}

// 测试便捷函数
func TestConvenienceFunctions(t *testing.T) {
	t.Run("Default", func(t *testing.T) {
		manager, err := Default()
		require.NoError(t, err)
		assert.NotNil(t, manager)
	})

	t.Run("Development", func(t *testing.T) {
		manager, err := Development()
		require.NoError(t, err)
		assert.NotNil(t, manager)
		assert.True(t, manager.GetConfig().Debug)
	})

	t.Run("Production", func(t *testing.T) {
		manager, err := Production()
		require.NoError(t, err)
		assert.NotNil(t, manager)
		assert.False(t, manager.GetConfig().Debug)
	})

	t.Run("Testing", func(t *testing.T) {
		manager, err := Testing()
		require.NoError(t, err)
		assert.NotNil(t, manager)
		assert.Contains(t, manager.GetConfig().ConfigPaths, "./testdata")
	})

	t.Run("WithFile", func(t *testing.T) {
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "test.yaml")
		err := os.WriteFile(configFile, []byte("test: value"), 0644)
		require.NoError(t, err)

		manager, err := WithFile(configFile)
		require.NoError(t, err)
		assert.NotNil(t, manager)
	})

	t.Run("WithEnv", func(t *testing.T) {
		os.Setenv("TEST_ENV_VAR", "development")
		defer os.Unsetenv("TEST_ENV_VAR")

		manager, err := WithEnv("TEST_ENV_VAR")
		require.NoError(t, err)
		assert.NotNil(t, manager)
	})

	t.Run("WithRemote", func(t *testing.T) {
		manager, err := WithRemote("etcd", "http://localhost:2379", "/config", YAMLFormat)
		require.NoError(t, err)
		assert.NotNil(t, manager)
	})

	t.Run("Custom", func(t *testing.T) {
		config := DefaultConfig()
		manager, err := Custom(config)
		require.NoError(t, err)
		assert.NotNil(t, manager)
	})
}

func TestMustLoad(t *testing.T) {
	t.Run("successful load", func(t *testing.T) {
		// 创建临时配置文件
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.yaml")
		err := os.WriteFile(configFile, []byte("app:\n  name: test"), 0644)
		require.NoError(t, err)

		manager, err := WithFile(configFile)
		require.NoError(t, err)

		// MustLoad 应该成功
		loadedManager := MustLoad(manager)
		assert.NotNil(t, loadedManager)
		assert.True(t, loadedManager.IsLoaded())
	})

	t.Run("failed load should panic", func(t *testing.T) {
		// 创建一个会失败的配置管理器
		config := &Config{
			ConfigName:     "nonexistent",
			ConfigType:     YAMLFormat,
			ConfigPaths:    []string{"/nonexistent/path"},
			AutoEnv:        false,
			RemoteProvider: "invalid", // 这会导致验证失败
		}

		manager, _ := NewManager(config)

		// MustLoad 应该 panic
		assert.Panics(t, func() {
			MustLoad(manager)
		})
	})
}

func TestLoadOrPanic(t *testing.T) {
	t.Run("successful create and load", func(t *testing.T) {
		// 创建临时配置文件
		tempDir := t.TempDir()
		configFile := filepath.Join(tempDir, "config.yaml")
		err := os.WriteFile(configFile, []byte("app:\n  name: test"), 0644)
		require.NoError(t, err)

		// LoadOrPanic 应该成功
		manager := LoadOrPanic(func() (*Manager, error) {
			return WithFile(configFile)
		})

		assert.NotNil(t, manager)
		assert.True(t, manager.IsLoaded())
	})

	t.Run("failed create should panic", func(t *testing.T) {
		// LoadOrPanic 应该 panic
		assert.Panics(t, func() {
			LoadOrPanic(func() (*Manager, error) {
				return WithFile("/nonexistent/config.yaml")
			})
		})
	})

	t.Run("failed load should panic", func(t *testing.T) {
		// LoadOrPanic 应该 panic
		assert.Panics(t, func() {
			LoadOrPanic(func() (*Manager, error) {
				// 创建一个会在加载时失败的管理器
				config := &Config{
					ConfigName:       "config",
					ConfigType:       YAMLFormat,
					ConfigPaths:      []string{"/nonexistent"},
					AutoEnv:          false,
					RemoteProvider:   "etcd",
					RemoteEndpoint:   "invalid-endpoint",
					RemoteConfigPath: "/config",
				}
				return NewManager(config)
			})
		})
	})
}
