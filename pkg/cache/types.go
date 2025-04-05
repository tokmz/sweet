package cache

// Mode Redis运行模式
type Mode string

const (
	// ModeSingle 单机模式
	ModeSingle Mode = "single"
	// ModeCluster 集群模式
	ModeCluster Mode = "cluster"
	// ModeSentinel 哨兵模式
	ModeSentinel Mode = "sentinel"
)

// Config Redis配置
type Config struct {
	// 运行模式：单机、集群、哨兵
	Mode Mode `json:"mode"`

	// 单机模式配置
	Single struct {
		// Redis服务器地址，格式：host:port
		Addr string `json:"addr"`
	} `json:"single"`

	// 集群模式配置
	Cluster struct {
		// Redis集群节点地址列表
		Addrs []string `json:"addrs"`
	} `json:"cluster"`

	// 哨兵模式配置
	Sentinel struct {
		// 主节点名称
		MasterName string `json:"master_name"`
		// 哨兵节点地址列表
		Addrs []string `json:"addrs"`
	} `json:"sentinel"`

	// 通用配置
	// 用户名
	Username string `json:"username"`
	// 密码
	Password string `json:"password"`
	// 默认使用的数据库，仅单机和哨兵模式有效
	DB int `json:"db"`
	// 连接池大小
	PoolSize int `json:"pool_size"`
	// 最小空闲连接数
	MinIdleConns int `json:"min_idle_conns"`
	// 连接最大空闲时间
	IdleTimeout int `json:"idle_timeout"`
	// 连接超时时间(毫秒)
	ConnTimeout int `json:"conn_timeout"`
	// 读取超时时间(毫秒)
	ReadTimeout int `json:"read_timeout"`
	// 写入超时时间(毫秒)
	WriteTimeout int `json:"write_timeout"`
	// 命令执行超时时间(毫秒)
	ExecTimeout int `json:"exec_timeout"`
	// 最大重试次数
	MaxRetries int `json:"max_retries"`
	// 重试间隔(毫秒)
	RetryDelay int `json:"retry_delay"`
	// 最小重试间隔(毫秒)
	MinRetryBackoff int `json:"min_retry_backoff"`
	// 最大重试间隔(毫秒)
	MaxRetryBackoff int `json:"max_retry_backoff"`
	// 是否开启链路追踪
	EnableTrace bool `json:"enable_trace"`
}
