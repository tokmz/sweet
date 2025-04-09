package cache

// 链路追踪相关常量
const (
	tracerName = "sweet/pkg/cache"
	opDo       = "redis.Do"
	opGet      = "redis.Get"
	opSet      = "redis.Set"
	opDel      = "redis.Del"
	opPipeline = "redis.Pipeline"
	opWatch    = "redis.Watch"
	opEval     = "redis.Eval"
	opPub      = "redis.Publish"
	opSub      = "redis.Subscribe"
)

// 默认超时和重试配置
const (
	defaultConnTimeout     = 5000 // 5秒
	defaultReadTimeout     = 3000 // 3秒
	defaultWriteTimeout    = 3000 // 3秒
	defaultExecTimeout     = 0    // 不设置超时
	defaultMaxRetries      = 3
	defaultRetryDelay      = 100 // 100毫秒
	defaultMinRetryBackoff = 8   // 8毫秒
	defaultMaxRetryBackoff = 512 // 512毫秒
	defaultPoolSize        = 10
	defaultMinIdleConns    = 5
	defaultIdleTimeout     = 300 // 5分钟
)

// 默认配置
var defaultConfig = Config{
	ConnTimeout:     defaultConnTimeout,
	ReadTimeout:     defaultReadTimeout,
	WriteTimeout:    defaultWriteTimeout,
	ExecTimeout:     defaultExecTimeout,
	MaxRetries:      defaultMaxRetries,
	RetryDelay:      defaultRetryDelay,
	MinRetryBackoff: defaultMinRetryBackoff,
	MaxRetryBackoff: defaultMaxRetryBackoff,
	PoolSize:        defaultPoolSize,
	MinIdleConns:    defaultMinIdleConns,
	IdleTimeout:     defaultIdleTimeout,
	EnableTrace:     false,
	EnableReadWrite: false,
}
