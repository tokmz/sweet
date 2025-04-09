package config

// Config 定义配置结构
//
// Config结构体用于指定配置文件的位置和格式信息。
// 它包含配置文件名称（不含扩展名）、配置文件类型（如yaml、json等）
// 以及配置文件的搜索路径列表。
type Config struct {
	Name string   // 配置文件名（不含扩展名）
	Type string   // 配置文件类型（yaml、json、toml等）
	Path []string // 配置文件搜索路径
}
