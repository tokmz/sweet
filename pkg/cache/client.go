package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

// Client Redis客户端接口
type Client interface {
	// Close 关闭连接
	Close() error

	// Do 执行自定义命令
	Do(ctx context.Context, args ...any) (any, error)

	/*
		String操作
	*/
	// Get 获取键值，如果键不存在返回ErrKeyNotExists错误
	Get(ctx context.Context, key string) (string, error)
	// GetWithExists 获取键值，返回值和键是否存在
	GetWithExists(ctx context.Context, key string) (string, bool, error)
	// Set 设置键值
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	// SetNX 不存在时设置键值
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	// MGet 批量获取键值
	MGet(ctx context.Context, keys ...string) ([]KeyValue, error)
	// MSet 批量设置键值
	MSet(ctx context.Context, pairs ...KeyValue) error
	// Del 删除键
	Del(ctx context.Context, keys ...string) (int64, error)
	// Exists 检查键是否存在
	Exists(ctx context.Context, keys ...string) (int64, error)
	// Expire 设置过期时间
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	// TTL 获取剩余生存时间
	TTL(ctx context.Context, key string) (time.Duration, error)

	/*
		Hash操作
	*/
	// HGet 获取哈希表字段值
	HGet(ctx context.Context, key, field string) (string, error)
	// HGetWithExists 获取哈希表字段值，返回值和字段是否存在
	HGetWithExists(ctx context.Context, key, field string) (string, bool, error)
	// HSet 设置哈希表字段值
	HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
	// HGetAll 获取哈希表所有字段和值
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	// HDel 删除哈希表字段
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	// HExists 检查哈希表字段是否存在
	HExists(ctx context.Context, key, field string) (bool, error)
	// HKeys 获取哈希表所有字段
	HKeys(ctx context.Context, key string) ([]string, error)
	// HVals 获取哈希表所有值
	HVals(ctx context.Context, key string) ([]string, error)
	// HLen 获取哈希表字段数量
	HLen(ctx context.Context, key string) (int64, error)

	/*
		List操作
	*/
	// LPush 将值插入到列表头部
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	// RPush 将值插入到列表尾部
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	// LPop 移除并返回列表第一个元素
	LPop(ctx context.Context, key string) (string, error)
	// LPopWithExists 移除并返回列表第一个元素，返回值和列表是否为空
	LPopWithExists(ctx context.Context, key string) (string, bool, error)
	// RPop 移除并返回列表最后一个元素
	RPop(ctx context.Context, key string) (string, error)
	// RPopWithExists 移除并返回列表最后一个元素，返回值和列表是否为空
	RPopWithExists(ctx context.Context, key string) (string, bool, error)
	// LRange 获取列表指定范围内的元素
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// LLen 获取列表长度
	LLen(ctx context.Context, key string) (int64, error)
	// LTrim 对列表进行修剪
	LTrim(ctx context.Context, key string, start, stop int64) error
	// LRem 移除列表中与value相等的元素
	LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error)

	/*
		Set操作
	*/
	// SAdd 向集合添加一个或多个成员
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	// SMembers 返回集合中的所有成员
	SMembers(ctx context.Context, key string) ([]string, error)
	// SRem 移除集合中一个或多个成员
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	// SIsMember 判断成员是否是集合的成员
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	// SCard 获取集合的成员数
	SCard(ctx context.Context, key string) (int64, error)
	// SDiff 返回第一个集合与其他集合之间的差异
	SDiff(ctx context.Context, keys ...string) ([]string, error)
	// SInter 返回所有给定集合的交集
	SInter(ctx context.Context, keys ...string) ([]string, error)
	// SUnion 返回所有给定集合的并集
	SUnion(ctx context.Context, keys ...string) ([]string, error)

	/*
		ZSet操作
	*/
	// ZAdd 向有序集合添加一个或多个成员
	ZAdd(ctx context.Context, key string, members ...*Z) (int64, error)
	// ZRange 返回有序集合中指定范围的成员
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	// ZRangeWithScores 返回有序集合中指定范围的成员和分数
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*Z, error)
	// ZRem 移除有序集合中的一个或多个成员
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	// ZCard 获取有序集合的成员数
	ZCard(ctx context.Context, key string) (int64, error)
	// ZScore 返回有序集合中成员的分数
	ZScore(ctx context.Context, key string, member string) (float64, error)
	// ZRank 返回有序集合中成员的排名
	ZRank(ctx context.Context, key string, member string) (int64, error)
	// ZRevRank 返回有序集合中成员的逆序排名
	ZRevRank(ctx context.Context, key string, member string) (int64, error)

	/*
		事务操作
	*/
	// TxPipeline 创建一个事务管道
	TxPipeline() Pipeliner
	// Watch 监视键值变化
	Watch(ctx context.Context, fn func(tx *Tx) error, keys ...string) error

	/*
		Lua脚本
	*/

	// Eval 执行Lua脚本
	Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error)
	// EvalSha 执行Lua脚本的SHA1校验和
	EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error)
	// ScriptLoad 将脚本加载到脚本缓存
	ScriptLoad(ctx context.Context, script string) (string, error)

	/*
	 发布订阅
	*/
	// Subscribe 订阅频道
	Subscribe(ctx context.Context, channels ...string) *PubSub
	// PSubscribe 订阅模式
	PSubscribe(ctx context.Context, patterns ...string) *PubSub
	// Publish 发布消息到频道
	Publish(ctx context.Context, channel string, message interface{}) (int64, error)
}

// Pipeliner 管道接口
type Pipeliner interface {
	// Exec 执行管道中的所有命令
	Exec(ctx context.Context) error
	// Discard 丢弃管道中的所有命令
	Discard() error
	// Get 获取键值
	Get(ctx context.Context, key string) *StringCmd
	// Set 设置键值
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) *StatusCmd
	// Del 删除键
	Del(ctx context.Context, keys ...string) *IntCmd
	// HGet 获取哈希表字段值
	HGet(ctx context.Context, key, field string) *StringCmd
	// HSet 设置哈希表字段值
	HSet(ctx context.Context, key string, values ...interface{}) *IntCmd
	// ZAdd 向有序集合添加一个或多个成员
	ZAdd(ctx context.Context, key string, members ...*Z) *IntCmd
}

// Tx 事务
type Tx struct {
	Pipeliner
}

// PubSub 发布订阅
type PubSub struct {
	channel <-chan *PubSubMessage
	pubsub  *redis.PubSub
}

// Channel 返回接收消息的通道
func (p *PubSub) Channel() <-chan *PubSubMessage {
	return p.channel
}

// Close 关闭订阅
func (p *PubSub) Close() error {
	if p.pubsub != nil {
		return p.pubsub.Close()
	}
	return nil
}

// StringCmd 字符串命令
type StringCmd struct {
	val string
	err error
}

// Result 返回命令结果
func (c *StringCmd) Result() (string, error) {
	return c.val, c.err
}

// StatusCmd 状态命令
type StatusCmd struct {
	err error
}

// Result 返回命令结果
func (c *StatusCmd) Result() error {
	return c.err
}

// IntCmd 整数命令
type IntCmd struct {
	val int64
	err error
}

// Result 返回命令结果
func (c *IntCmd) Result() (int64, error) {
	return c.val, c.err
}
