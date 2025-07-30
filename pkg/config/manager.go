package config

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
)

// Manager Viper配置管理器
type Manager struct {
	config *Config
	viper  *viper.Viper
	mu     sync.RWMutex
	loaded bool
}

// NewManager 创建新的配置管理器
func NewManager(config *Config) (*Manager, error) {
	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	v := viper.New()

	m := &Manager{
		config: config,
		viper:  v,
	}

	return m, nil
}

// NewManagerWithViper 使用现有的viper实例创建管理器
func NewManagerWithViper(v *viper.Viper, config *Config) (*Manager, error) {
	if v == nil {
		return nil, fmt.Errorf("viper instance cannot be nil")
	}

	if config == nil {
		config = DefaultConfig()
	}

	if err := config.Validate(); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	m := &Manager{
		config: config,
		viper:  v,
	}

	return m, nil
}

// Load 加载配置
func (m *Manager) Load() error {
	m.mu.Lock()
	defer m.mu.Unlock()

	// 设置配置文件信息
	m.viper.SetConfigName(m.config.ConfigName)
	m.viper.SetConfigType(string(m.config.ConfigType))

	// 添加搜索路径
	for _, path := range m.config.GetSearchPaths() {
		m.viper.AddConfigPath(path)
	}

	// 配置环境变量
	if m.config.AutoEnv {
		m.viper.AutomaticEnv()
	}

	if m.config.EnvPrefix != "" {
		m.viper.SetEnvPrefix(m.config.EnvPrefix)
	}

	if m.config.EnvKeyReplacer != "" {
		replacer := strings.NewReplacer(".", m.config.EnvKeyReplacer, "-", m.config.EnvKeyReplacer)
		m.viper.SetEnvKeyReplacer(replacer)
	}

	m.viper.AllowEmptyEnv(m.config.AllowEmptyEnv)

	// 加载远程配置
	if m.config.RemoteProvider != "" {
		if err := m.loadRemoteConfig(); err != nil {
			return fmt.Errorf("failed to load remote config: %w", err)
		}
	} else {
		// 加载本地配置文件
		if err := m.viper.ReadInConfig(); err != nil {
			if _, ok := err.(viper.ConfigFileNotFoundError); ok {
				if m.config.Debug {
					log.Printf("Config file not found: %v", err)
				}
				// 配置文件不存在时不返回错误，允许仅使用环境变量
			} else {
				return fmt.Errorf("failed to read config file: %w", err)
			}
		} else {
			if m.config.Debug {
				log.Printf("Using config file: %s", m.viper.ConfigFileUsed())
			}
		}
	}

	// 启用配置监听
	if m.config.WatchConfig {
		m.viper.WatchConfig()
		if m.config.OnConfigChange != nil {
			m.viper.OnConfigChange(func(e fsnotify.Event) {
				if m.config.Debug {
					log.Printf("Config file changed: %s", e.Name)
				}
				m.config.OnConfigChange()
			})
		}
	}

	m.loaded = true
	return nil
}

// loadRemoteConfig 加载远程配置
func (m *Manager) loadRemoteConfig() error {
	err := m.viper.AddRemoteProvider(m.config.RemoteProvider, m.config.RemoteEndpoint, m.config.RemoteConfigPath)
	if err != nil {
		return fmt.Errorf("failed to add remote provider: %w", err)
	}

	m.viper.SetConfigType(string(m.config.RemoteConfigType))

	err = m.viper.ReadRemoteConfig()
	if err != nil {
		return fmt.Errorf("failed to read remote config: %w", err)
	}

	// 启用远程配置监听
	if m.config.WatchConfig {
		go func() {
			for {
				time.Sleep(m.config.WatchInterval)
				err := m.viper.WatchRemoteConfig()
				if err != nil {
					if m.config.Debug {
						log.Printf("Failed to watch remote config: %v", err)
					}
					continue
				}
				if m.config.OnConfigChange != nil {
					m.config.OnConfigChange()
				}
			}
		}()
	}

	return nil
}

// IsLoaded 检查配置是否已加载
func (m *Manager) IsLoaded() bool {
	m.mu.RLock()
	defer m.mu.RUnlock()
	return m.loaded
}

// GetViper 获取底层的viper实例
func (m *Manager) GetViper() *viper.Viper {
	return m.viper
}

// GetConfig 获取管理器配置
func (m *Manager) GetConfig() *Config {
	return m.config
}

// Get 获取配置值
func (m *Manager) Get(key string) interface{} {
	return m.viper.Get(key)
}

// GetString 获取字符串配置值
func (m *Manager) GetString(key string) string {
	return m.viper.GetString(key)
}

// GetBool 获取布尔配置值
func (m *Manager) GetBool(key string) bool {
	return m.viper.GetBool(key)
}

// GetInt 获取整数配置值
func (m *Manager) GetInt(key string) int {
	return m.viper.GetInt(key)
}

// GetInt32 获取32位整数配置值
func (m *Manager) GetInt32(key string) int32 {
	return m.viper.GetInt32(key)
}

// GetInt64 获取64位整数配置值
func (m *Manager) GetInt64(key string) int64 {
	return m.viper.GetInt64(key)
}

// GetUint 获取无符号整数配置值
func (m *Manager) GetUint(key string) uint {
	return m.viper.GetUint(key)
}

// GetUint32 获取32位无符号整数配置值
func (m *Manager) GetUint32(key string) uint32 {
	return m.viper.GetUint32(key)
}

// GetUint64 获取64位无符号整数配置值
func (m *Manager) GetUint64(key string) uint64 {
	return m.viper.GetUint64(key)
}

// GetFloat64 获取浮点数配置值
func (m *Manager) GetFloat64(key string) float64 {
	return m.viper.GetFloat64(key)
}

// GetTime 获取时间配置值
func (m *Manager) GetTime(key string) time.Time {
	return m.viper.GetTime(key)
}

// GetDuration 获取时间间隔配置值
func (m *Manager) GetDuration(key string) time.Duration {
	return m.viper.GetDuration(key)
}

// GetIntSlice 获取整数切片配置值
func (m *Manager) GetIntSlice(key string) []int {
	return m.viper.GetIntSlice(key)
}

// GetStringSlice 获取字符串切片配置值
func (m *Manager) GetStringSlice(key string) []string {
	return m.viper.GetStringSlice(key)
}

// GetStringMap 获取字符串映射配置值
func (m *Manager) GetStringMap(key string) map[string]interface{} {
	return m.viper.GetStringMap(key)
}

// GetStringMapString 获取字符串到字符串映射配置值
func (m *Manager) GetStringMapString(key string) map[string]string {
	return m.viper.GetStringMapString(key)
}

// GetStringMapStringSlice 获取字符串到字符串切片映射配置值
func (m *Manager) GetStringMapStringSlice(key string) map[string][]string {
	return m.viper.GetStringMapStringSlice(key)
}

// GetSizeInBytes 获取字节大小配置值
func (m *Manager) GetSizeInBytes(key string) uint {
	return m.viper.GetSizeInBytes(key)
}

// Set 设置配置值
func (m *Manager) Set(key string, value interface{}) {
	m.viper.Set(key, value)
}

// SetDefault 设置默认配置值
func (m *Manager) SetDefault(key string, value interface{}) {
	m.viper.SetDefault(key, value)
}

// IsSet 检查配置键是否已设置
func (m *Manager) IsSet(key string) bool {
	return m.viper.IsSet(key)
}

// AllKeys 获取所有配置键
func (m *Manager) AllKeys() []string {
	return m.viper.AllKeys()
}

// AllSettings 获取所有配置设置
func (m *Manager) AllSettings() map[string]interface{} {
	return m.viper.AllSettings()
}

// Unmarshal 将配置解析到结构体
func (m *Manager) Unmarshal(rawVal interface{}) error {
	return m.viper.Unmarshal(rawVal)
}

// UnmarshalKey 将指定键的配置解析到结构体
func (m *Manager) UnmarshalKey(key string, rawVal interface{}) error {
	return m.viper.UnmarshalKey(key, rawVal)
}

// UnmarshalExact 精确解析配置到结构体(不允许未知字段)
func (m *Manager) UnmarshalExact(rawVal interface{}) error {
	return m.viper.UnmarshalExact(rawVal)
}

// BindEnv 绑定环境变量
func (m *Manager) BindEnv(input ...string) error {
	return m.viper.BindEnv(input...)
}

// BindPFlag 绑定pflags
func (m *Manager) BindPFlag(key string, flag *pflag.Flag) error {
	return m.viper.BindPFlag(key, flag)
}

// BindPFlags 绑定pflags集合
func (m *Manager) BindPFlags(flags *pflag.FlagSet) error {
	return m.viper.BindPFlags(flags)
}

// WriteConfig 写入配置文件
func (m *Manager) WriteConfig() error {
	return m.viper.WriteConfig()
}

// SafeWriteConfig 安全写入配置文件(文件不存在时才写入)
func (m *Manager) SafeWriteConfig() error {
	return m.viper.SafeWriteConfig()
}

// WriteConfigAs 写入配置到指定文件
func (m *Manager) WriteConfigAs(filename string) error {
	return m.viper.WriteConfigAs(filename)
}

// SafeWriteConfigAs 安全写入配置到指定文件
func (m *Manager) SafeWriteConfigAs(filename string) error {
	return m.viper.SafeWriteConfigAs(filename)
}

// ConfigFileUsed 获取使用的配置文件路径
func (m *Manager) ConfigFileUsed() string {
	return m.viper.ConfigFileUsed()
}

// Debug 打印调试信息
func (m *Manager) Debug() {
	if !m.config.Debug {
		return
	}

	log.Println("=== Config Manager Debug Info ===")
	log.Printf("Config Name: %s", m.config.ConfigName)
	log.Printf("Config Type: %s", m.config.ConfigType)
	log.Printf("Config Paths: %v", m.config.ConfigPaths)
	log.Printf("Env Prefix: %s", m.config.EnvPrefix)
	log.Printf("Auto Env: %t", m.config.AutoEnv)
	log.Printf("Watch Config: %t", m.config.WatchConfig)
	log.Printf("Config File Used: %s", m.ConfigFileUsed())
	log.Printf("All Keys: %v", m.AllKeys())
	log.Println("=== End Debug Info ===")
}
