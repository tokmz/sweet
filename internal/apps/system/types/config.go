package types

// Config 系统配置
type Config struct {
	App      AppConfig      `mapstructure:"app"`
	Database DatabaseConfig `mapstructure:"database"`
	Cache    CacheConfig    `mapstructure:"cache"`
}

// AppConfig 应用配置
type AppConfig struct {
	Name    string `mapstructure:"name"`    // 应用名称
	Version string `mapstructure:"version"` // 应用版本
	Intro   string `mapstructure:"intro"`   // 模块简介
}

// DatabaseConfig 数据库配置
type DatabaseConfig struct {
	Master        string             `mapstructure:"master"`         // 主库连接DSN
	Slaves        []string           `mapstructure:"slaves"`         // 从库连接DSN数组
	Pool          DatabasePoolConfig `mapstructure:"pool"`           // 连接池配置
	EnableLog     bool               `mapstructure:"enable_log"`     // 是否开启日志
	EnableTrace   bool               `mapstructure:"enable_trace"`   // 是否开启链路追踪
	SlowThreshold int                `mapstructure:"slow_threshold"` // 慢查询阈值（毫秒）
	DriverType    string             `mapstructure:"driver_type"`    // 数据库驱动类型
	QueryTimeout  int                `mapstructure:"query_timeout"`  // 查询超时设置（毫秒）
	MaxRetries    int                `mapstructure:"max_retries"`    // 最大重试次数
	RetryDelay    int                `mapstructure:"retry_delay"`    // 重试间隔（毫秒）
}

// DatabasePoolConfig 数据库连接池配置
type DatabasePoolConfig struct {
	MaxIdleConns    int `mapstructure:"max_idle_conns"`    // 最大空闲连接数
	MaxOpenConns    int `mapstructure:"max_open_conns"`    // 最大打开连接数
	ConnMaxLifetime int `mapstructure:"conn_max_lifetime"` // 连接最大生命周期（秒）
}

// CacheConfig 缓存配置
type CacheConfig struct {
	Mode            string         `mapstructure:"mode"`              // 运行模式：single（单机）, cluster（集群）, sentinel（哨兵）
	Single          SingleConfig   `mapstructure:"single"`            // 单机模式配置
	Cluster         ClusterConfig  `mapstructure:"cluster"`           // 集群模式配置
	Sentinel        SentinelConfig `mapstructure:"sentinel"`          // 哨兵模式配置
	Username        string         `mapstructure:"username"`          // Redis用户名
	Password        string         `mapstructure:"password"`          // Redis密码
	DB              int            `mapstructure:"db"`                // 默认使用的数据库
	PoolSize        int            `mapstructure:"pool_size"`         // 连接池大小
	MinIdleConns    int            `mapstructure:"min_idle_conns"`    // 最小空闲连接数
	IdleTimeout     int            `mapstructure:"idle_timeout"`      // 连接最大空闲时间(秒)
	ConnTimeout     int            `mapstructure:"conn_timeout"`      // 连接超时时间(毫秒)
	ReadTimeout     int            `mapstructure:"read_timeout"`      // 读取超时时间(毫秒)
	WriteTimeout    int            `mapstructure:"write_timeout"`     // 写入超时时间(毫秒)
	ExecTimeout     int            `mapstructure:"exec_timeout"`      // 命令执行超时时间(毫秒)
	MaxRetries      int            `mapstructure:"max_retries"`       // 最大重试次数
	RetryDelay      int            `mapstructure:"retry_delay"`       // 重试间隔(毫秒)
	MinRetryBackoff int            `mapstructure:"min_retry_backoff"` // 最小重试间隔(毫秒)
	MaxRetryBackoff int            `mapstructure:"max_retry_backoff"` // 最大重试间隔(毫秒)
	EnableTrace     bool           `mapstructure:"enable_trace"`      // 是否开启链路追踪
	EnableReadWrite bool           `mapstructure:"enable_read_write"` // 是否启用读写分离
}

// SingleConfig 单机模式配置
type SingleConfig struct {
	Addr string `mapstructure:"addr"` // Redis地址
}

// ClusterConfig 集群模式配置
type ClusterConfig struct {
	Addrs []string `mapstructure:"addrs"` // Redis集群地址
}

// SentinelConfig 哨兵模式配置
type SentinelConfig struct {
	MasterName string   `mapstructure:"master_name"` // 主节点名称
	Addrs      []string `mapstructure:"addrs"`       // 哨兵地址
}

// LoggerConfig 日志配置
type LoggerConfig struct {
	Level            string `mapstructure:"level"`             // 日志级别: debug, info, warn, error, fatal
	EnableConsole    bool   `mapstructure:"enable_console"`    // 是否输出到控制台
	EnableFile       bool   `mapstructure:"enable_file"`       // 是否输出到文件
	Filename         string `mapstructure:"filename"`          // 日志文件路径
	MaxSize          int    `mapstructure:"max_size"`          // 单个日志文件最大尺寸，单位MB
	MaxBackups       int    `mapstructure:"max_backups"`       // 保留的旧日志文件最大数量
	MaxAge           int    `mapstructure:"max_age"`           // 保留旧日志文件的最大天数
	Compress         bool   `mapstructure:"compress"`          // 是否压缩旧日志文件
	DisableTimestamp bool   `mapstructure:"disable_timestamp"` // 是否禁用时间戳
	DisableCaller    bool   `mapstructure:"disable_caller"`    // 是否禁用调用者信息
	DisableTrace     bool   `mapstructure:"disable_trace"`     // 是否禁用堆栈跟踪
	Development      bool   `mapstructure:"development"`       // 是否为开发模式
}

// JWTConfig JWT配置
type JWTConfig struct {
	Secret string `mapstructure:"secret"` // JWT密钥
	Expire int    `mapstructure:"expire"` // 过期时间（秒）
	Issuer string `mapstructure:"issuer"` // 签发者
}

// CORSConfig 跨域配置
type CORSConfig struct {
	Enable           bool     `mapstructure:"enable"`            // 是否启用跨域
	AllowOrigins     []string `mapstructure:"allow_origins"`     // 允许的源
	AllowMethods     []string `mapstructure:"allow_methods"`     // 允许的方法
	AllowHeaders     []string `mapstructure:"allow_headers"`     // 允许的头部
	ExposeHeaders    []string `mapstructure:"expose_headers"`    // 暴露的头部
	AllowCredentials bool     `mapstructure:"allow_credentials"` // 是否允许凭证
	MaxAge           int      `mapstructure:"max_age"`           // 预检请求缓存时间（秒）
}

//
//// LoadConfig 加载配置
//func LoadConfig() (*Config, error) {
//	var cfg Config
//	// 初始化配置
//	configParams := &config.Config{
//		Name: "config",
//		Type: "yaml",
//		Path: []string{"./internal/apps/system"},
//	}
//
//	// 初始化Viper
//	if err := config.InitConfig(configParams); err != nil {
//		return nil, err
//	}
//
//	// 加载配置到结构体
//	if err := config.LoadConfig(&cfg); err != nil {
//		return nil, err
//	}
//
//	return &cfg, nil
//}
