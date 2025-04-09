package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
)

// Do 执行自定义命令
func (c *singleClient) Do(ctx context.Context, args ...interface{}) (interface{}, error) {
	ctx, span := c.startTrace(ctx, opDo)
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	cmd := c.client.Do(ctx, args...)
	err := cmd.Err()
	c.endTrace(span, err)
	return cmd.Val(), c.handleCommandError("Do", err)
}

// Get 获取键值
func (c *singleClient) Get(ctx context.Context, key string) (string, error) {
	ctx, span := c.startTrace(ctx, opGet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	val, err := c.client.Get(ctx, key).Result()
	c.endTrace(span, err)
	return val, c.handleCommandError("Get", err)
}

// GetWithExists 获取键值，返回值和键是否存在
func (c *singleClient) GetWithExists(ctx context.Context, key string) (string, bool, error) {
	ctx, span := c.startTrace(ctx, opGet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	val, err := c.client.Get(ctx, key).Result()
	if err == redis.Nil {
		c.endTrace(span, nil)
		return "", false, nil
	}
	if err != nil {
		c.endTrace(span, err)
		return "", false, &CommandError{Command: "Get", Err: err}
	}
	c.endTrace(span, nil)
	return val, true, nil
}

// Set 设置键值
func (c *singleClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, span := c.startTrace(ctx, opSet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	err := c.client.Set(ctx, key, value, expiration).Err()
	c.endTrace(span, err)
	return c.handleCommandError("Set", err)
}

// SetNX 不存在时设置键值
func (c *singleClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ctx, span := c.startTrace(ctx, opSet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	ok, err := c.client.SetNX(ctx, key, value, expiration).Result()
	c.endTrace(span, err)
	return ok, c.handleCommandError("SetNX", err)
}

// Del 删除键
func (c *singleClient) Del(ctx context.Context, keys ...string) (int64, error) {
	ctx, span := c.startTrace(ctx, opDel)
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	n, err := c.client.Del(ctx, keys...).Result()
	c.endTrace(span, err)
	return n, c.handleCommandError("Del", err)
}

// Exists 检查键是否存在
func (c *singleClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	ctx, span := c.startTrace(ctx, "redis.Exists")
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	n, err := c.client.Exists(ctx, keys...).Result()
	c.endTrace(span, err)
	return n, c.handleCommandError("Exists", err)
}

// 以下是为了满足接口而添加的存根方法，实际项目中应该完整实现

// MGet 批量获取键值
func (c *singleClient) MGet(ctx context.Context, keys ...string) ([]KeyValue, error) {
	return nil, ErrCommandFailed
}

// MSet 批量设置键值
func (c *singleClient) MSet(ctx context.Context, pairs ...KeyValue) error {
	return ErrCommandFailed
}

// Expire 设置过期时间
func (c *singleClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return false, ErrCommandFailed
}

// TTL 获取剩余生存时间
func (c *singleClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, ErrCommandFailed
}

// HGet 获取哈希表字段值
func (c *singleClient) HGet(ctx context.Context, key, field string) (string, error) {
	return "", ErrCommandFailed
}

// HGetWithExists 获取哈希表字段值，返回值和字段是否存在
func (c *singleClient) HGetWithExists(ctx context.Context, key, field string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

// HSet 设置哈希表字段值
func (c *singleClient) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

// HGetAll 获取哈希表所有字段和值
func (c *singleClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return nil, ErrCommandFailed
}

// HDel 删除哈希表字段
func (c *singleClient) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return 0, ErrCommandFailed
}

// HExists 检查哈希表字段是否存在
func (c *singleClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return false, ErrCommandFailed
}

// HKeys 获取哈希表所有字段
func (c *singleClient) HKeys(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

// HVals 获取哈希表所有值
func (c *singleClient) HVals(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

// HLen 获取哈希表字段数量
func (c *singleClient) HLen(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

// 其他接口方法的存根实现
// 列表相关操作
func (c *singleClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) LPop(ctx context.Context, key string) (string, error) {
	return "", ErrCommandFailed
}

func (c *singleClient) LPopWithExists(ctx context.Context, key string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *singleClient) RPop(ctx context.Context, key string) (string, error) {
	return "", ErrCommandFailed
}

func (c *singleClient) RPopWithExists(ctx context.Context, key string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *singleClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) LLen(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	return ErrCommandFailed
}

func (c *singleClient) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

// 集合相关操作
func (c *singleClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, ErrCommandFailed
}

func (c *singleClient) SCard(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) SInter(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

// 有序集合相关操作
func (c *singleClient) ZAdd(ctx context.Context, key string, members ...*Z) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*Z, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) ZCard(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) ZRank(ctx context.Context, key string, member string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *singleClient) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return 0, ErrCommandFailed
}

// 事务相关操作
func (c *singleClient) TxPipeline() Pipeliner {
	return nil
}

func (c *singleClient) Watch(ctx context.Context, fn func(tx *Tx) error, keys ...string) error {
	return ErrCommandFailed
}

// Lua脚本相关操作
func (c *singleClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return nil, ErrCommandFailed
}

func (c *singleClient) ScriptLoad(ctx context.Context, script string) (string, error) {
	return "", ErrCommandFailed
}

// 发布订阅相关操作
func (c *singleClient) Subscribe(ctx context.Context, channels ...string) *PubSub {
	return &PubSub{channel: nil}
}

func (c *singleClient) PSubscribe(ctx context.Context, patterns ...string) *PubSub {
	return &PubSub{channel: nil}
}

func (c *singleClient) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	return 0, ErrCommandFailed
}
