package database

type Config struct {
	// 主库
	Master string `json:"master"` // dsn user:pass@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local
	// 从库
	Slaves []string `json:"slaves"`
	Pool   ConnPool `json:"pool"`
	// 是否开启日志
	EnableLog bool `json:"enable_log"`
	// 是否开始链路追踪
	EnableTrace bool `json:"enable_trace"`
	// 慢查询阈值（毫秒）
	SlowThreshold int `json:"slow_threshold"`
	// 数据库驱动类型
	DriverType string `json:"driver_type"`
	// 查询超时设置（毫秒），0表示不设置超时
	QueryTimeout int `json:"query_timeout"`
	// 最大重试次数
	MaxRetries int `json:"max_retries"`
	// 重试间隔（毫秒）
	RetryDelay int `json:"retry_delay"`
}

type ConnPool struct {
	MaxIdleConns    int `json:"max_idle_conns"`
	MaxOpenConns    int `json:"max_open_conns"`
	ConnMaxLifetime int `json:"conn_max_lifetime"`
}
