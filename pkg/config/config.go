package config

import (
	"fmt"
	"path/filepath"
	"strings"

	"github.com/fsnotify/fsnotify"
	"github.com/spf13/viper"
)

// InitConfig 初始化Viper配置系统
//
// 该函数接收一个Config结构体指针作为参数，该结构体包含配置文件的名称、类型和搜索路径。
// 函数会验证配置类型是否有效，创建一个新的Viper实例，设置配置文件名和类型，
// 添加配置搜索路径，并尝试读取配置文件。
//
// 参数:
//   - cfg: 配置参数，包含名称、类型和搜索路径
//
// 返回:
//   - error: 如果初始化过程中发生错误，则返回相应的错误；否则返回nil
func InitConfig(cfg *Config) error {
	rw.Lock()
	defer rw.Unlock()

	// Validate config file type
	if !isValidConfigType(cfg.Type) {
		return ErrInvalidConfigType
	}

	// Create a new Viper instance
	v := viper.New()
	v.SetConfigName(cfg.Name)
	v.SetConfigType(cfg.Type)

	// Add config paths
	for _, path := range cfg.Path {
		v.AddConfigPath(path)
	}

	// Set defaults, environment variables, etc. if needed
	v.AutomaticEnv()
	v.SetEnvKeyReplacer(strings.NewReplacer(".", "_"))

	// Read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrConfigNotFound
		}
		return fmt.Errorf("%w: %s", ErrReadConfig, err.Error())
	}

	// Store the Viper instance
	Viper = v
	return nil
}

// LoadConfig 将配置加载到提供的结构体中
//
// 该函数首先检查Viper实例是否已初始化，如果没有，则调用InitConfig进行初始化。
// 然后将配置解析到目标结构体中。结构体字段应使用mapstructure标签指定对应的配置键。
//
// 参数:
//   - cfg: 配置参数，包含名称、类型和搜索路径
//   - target: 指向要填充配置值的结构体的指针
//
// 返回:
//   - error: 如果加载过程中发生错误，则返回相应的错误；否则返回nil
func LoadConfig(target any) error {
	// Initialize config if needed
	if Viper == nil {
		return fmt.Errorf("viper is not initialized")
	}

	rw.RLock()
	defer rw.RUnlock()

	// Unmarshal config into target struct
	if err := Viper.Unmarshal(target); err != nil {
		return fmt.Errorf("%w: %s", ErrConfigValidation, err.Error())
	}

	return nil
}

// LoadConfigFile 从特定文件加载配置到提供的结构体中
//
// 该函数接收一个文件路径和一个目标结构体指针。它会创建一个新的Viper实例，
// 设置配置文件，验证配置类型是否有效，然后尝试读取配置文件并解析到目标结构体中。
//
// 参数:
//   - filePath: 配置文件的完整路径
//   - target: 指向要填充配置值的结构体的指针
//
// 返回:
//   - error: 如果加载过程中发生错误，则返回相应的错误；否则返回nil
func LoadConfigFile(filePath string, target interface{}) error {
	rw.Lock()
	defer rw.Unlock()

	// Create a new Viper instance
	v := viper.New()

	// Set config file
	ext := filepath.Ext(filePath)
	if ext == "" {
		return ErrConfigFileType
	}

	// Remove dot from extension (.yaml -> yaml)
	configType := ext[1:]
	if !isValidConfigType(configType) {
		return ErrInvalidConfigType
	}

	v.SetConfigFile(filePath)

	// Read config
	if err := v.ReadInConfig(); err != nil {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			return ErrConfigNotFound
		}
		return fmt.Errorf("%w: %s", ErrReadConfig, err.Error())
	}

	// Unmarshal config into target struct
	if err := v.Unmarshal(target); err != nil {
		return fmt.Errorf("%w: %s", ErrConfigValidation, err.Error())
	}

	// Store the Viper instance
	Viper = v
	return nil
}

// isValidConfigType 检查提供的配置类型是否受支持
//
// 参数:
//   - configType: 要检查的配置类型字符串（例如"yaml"、"json"等）
//
// 返回:
//   - bool: 如果配置类型受支持则返回true，否则返回false
func isValidConfigType(configType string) bool {
	for _, supportedType := range SupportedConfigTypes {
		if configType == supportedType {
			return true
		}
	}
	return false
}

// WatchConfig 设置配置文件变更监视器
//
// 该函数设置一个监视器来监控配置文件的变更。当配置文件发生变更时，
// 会调用提供的回调函数。这对于实现配置热重载非常有用。
//
// 参数:
//   - callback: 当配置文件变更时要调用的函数
//
// 返回:
//   - error: 如果设置监视器过程中发生错误，则返回相应的错误；否则返回nil
func WatchConfig(callback func()) error {
	if Viper == nil {
		return fmt.Errorf("viper is not initialized")
	}

	rw.RLock()
	defer rw.RUnlock()

	Viper.WatchConfig()
	Viper.OnConfigChange(func(e fsnotify.Event) {
		if callback != nil {
			callback()
		}
	})

	return nil
}
