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
)

// 默认超时和重试配置
const (
	defaultConnTimeout     = 5000 // 5秒
	defaultReadTimeout     = 3000 // 3秒
	defaultWriteTimeout    = 3000 // 3秒
	defaultMaxRetries      = 3
	defaultRetryDelay      = 100 // 100毫秒
	defaultMinRetryBackoff = 8   // 8毫秒
	defaultMaxRetryBackoff = 512 // 512毫秒
	defaultPoolSize        = 10
	defaultMinIdleConns    = 5
	defaultIdleTimeout     = 300 // 5分钟
)

// 错误信息
const (
	ErrInvalidMode    = "无效的Redis模式"
	ErrEmptyAddrs     = "未提供Redis地址"
	ErrEmptyMasterSet = "未提供哨兵主节点名称"
)
