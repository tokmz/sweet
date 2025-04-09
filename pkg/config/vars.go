package config

import (
	"errors"
	"sync"

	"github.com/spf13/viper"
)

// 全局变量
var (
	// Viper 是Viper库的全局实例，提供对配置的访问
	Viper *viper.Viper

	// rw 是一个用于保护Viper实例的读写锁
	rw sync.RWMutex
)

// SupportedConfigTypes 包含所有支持的配置文件类型
//
// 当前支持的类型为：
// - yaml: YAML格式
// - json: JSON格式
// - toml: TOML格式
// - hcl: HCL格式
// - ini: INI格式
// - env: 环境变量格式
var SupportedConfigTypes = []string{"yaml", "json", "toml", "hcl", "ini", "env"}

// 错误定义
var (
	// ErrConfigNotFound 表示找不到配置文件
	ErrConfigNotFound = errors.New("配置文件未找到")

	// ErrConfigFileType 表示配置文件类型错误（例如，配置文件没有扩展名）
	ErrConfigFileType = errors.New("配置文件类型错误")

	// ErrInvalidConfigType 表示指定的配置文件类型不受支持
	ErrInvalidConfigType = errors.New("无效的配置文件类型")

	// ErrReadConfig 表示读取配置文件时发生错误
	ErrReadConfig = errors.New("读取配置文件失败")

	// ErrConfigValidation 表示配置验证失败（例如，解析配置到结构体时）
	ErrConfigValidation = errors.New("配置验证失败")
)
