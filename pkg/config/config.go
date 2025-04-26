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
// 该函数是线程安全的，可以在多个goroutine中并发调用。
// 如果已经初始化过，会先关闭之前的实例，然后创建新的实例。
//
// 参数:
//   - cfg: 配置参数，包含名称、类型和搜索路径
//
// 返回:
//   - error: 如果初始化过程中发生错误，则返回相应的错误；否则返回nil
func InitConfig(cfg *Config) error {
	if cfg == nil {
		return ErrInvalidConfig
	}

	rw.Lock()
	defer rw.Unlock()

	// 如果已经初始化，先关闭之前的资源（如有必要）
	if initialized && Viper != nil {
		// Viper没有Close方法，但如果将来需要清理资源，可以在这里添加
		// 目前只需要重置Viper实例
	}

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
	initialized = true
	return nil
}

// LoadConfig 将配置加载到提供的结构体中
//
// 该函数首先检查Viper实例是否已初始化，如果没有，则返回错误。
// 然后将配置解析到目标结构体中。结构体字段应使用mapstructure标签指定对应的配置键。
// 该函数是线程安全的，可以在多个goroutine中并发调用。
//
// 参数:
//   - target: 指向要填充配置值的结构体的指针
//
// 返回:
//   - error: 如果加载过程中发生错误，则返回相应的错误；否则返回nil
func LoadConfig(target any) error {
	if target == nil {
		return ErrInvalidTarget
	}

	rw.RLock()
	defer rw.RUnlock()

	// 检查是否已初始化
	if !initialized || Viper == nil {
		return ErrNotInitialized
	}

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
// 该函数是线程安全的，可以在多个goroutine中并发调用。
// 如果已经初始化过，会先关闭之前的实例，然后创建新的实例。
//
// 参数:
//   - filePath: 配置文件的完整路径
//   - target: 指向要填充配置值的结构体的指针
//
// 返回:
//   - error: 如果加载过程中发生错误，则返回相应的错误；否则返回nil
func LoadConfigFile(filePath string, target interface{}) error {
	if filePath == "" {
		return ErrInvalidFilePath
	}

	if target == nil {
		return ErrInvalidTarget
	}

	rw.Lock()
	defer rw.Unlock()

	// 如果已经初始化，先关闭之前的资源（如有必要）
	if initialized && Viper != nil {
		// Viper没有Close方法，但如果将来需要清理资源，可以在这里添加
	}

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
	initialized = true
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

// Close 关闭配置系统，释放资源
//
// 该函数用于在应用程序结束时清理配置系统资源。
// 目前Viper没有需要显式关闭的资源，但保留此函数以便将来扩展。
// 该函数是线程安全的，可以在多个goroutine中并发调用。
//
// 返回:
//   - error: 如果关闭过程中发生错误，则返回相应的错误；否则返回nil
func Close() error {
	rw.Lock()
	defer rw.Unlock()

	if !initialized {
		return nil
	}

	// 目前Viper没有需要显式关闭的资源
	// 如果将来需要清理资源，可以在这里添加

	Viper = nil
	initialized = false
	return nil
}

// IsInitialized 检查配置系统是否已初始化
//
// 该函数返回配置系统的初始化状态。
// 该函数是线程安全的，可以在多个goroutine中并发调用。
//
// 返回:
//   - bool: 如果配置系统已初始化则返回true，否则返回false
func IsInitialized() bool {
	rw.RLock()
	defer rw.RUnlock()

	return initialized && Viper != nil
}

// WatchConfig 设置配置文件变更监视器
//
// 该函数设置一个监视器来监控配置文件的变更。当配置文件发生变更时，
// 会调用提供的回调函数。这对于实现配置热重载非常有用。
// 该函数是线程安全的，可以在多个goroutine中并发调用。
//
// 参数:
//   - callback: 当配置文件变更时要调用的函数
//
// 返回:
//   - error: 如果设置监视器过程中发生错误，则返回相应的错误；否则返回nil
func WatchConfig(callback func()) error {
	rw.RLock()
	defer rw.RUnlock()

	if !initialized || Viper == nil {
		return ErrNotInitialized
	}

	Viper.WatchConfig()
	Viper.OnConfigChange(func(e fsnotify.Event) {
		if callback != nil {
			callback()
		}
	})

	return nil
}
