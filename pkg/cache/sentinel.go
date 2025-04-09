package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
)

// Do 执行自定义命令
func (c *sentinelClient) Do(ctx context.Context, args ...interface{}) (interface{}, error) {
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
func (c *sentinelClient) Get(ctx context.Context, key string) (string, error) {
	ctx, span := c.startTrace(ctx, opGet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	val, err := c.client.Get(ctx, key).Result()
	c.endTrace(span, err)
	return val, c.handleCommandError("Get", err)
}

// GetWithExists 获取键值，返回值和键是否存在
func (c *sentinelClient) GetWithExists(ctx context.Context, key string) (string, bool, error) {
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
func (c *sentinelClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, span := c.startTrace(ctx, opSet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	err := c.client.Set(ctx, key, value, expiration).Err()
	c.endTrace(span, err)
	return c.handleCommandError("Set", err)
}

// SetNX 不存在时设置键值
func (c *sentinelClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ctx, span := c.startTrace(ctx, opSet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	ok, err := c.client.SetNX(ctx, key, value, expiration).Result()
	c.endTrace(span, err)
	return ok, c.handleCommandError("SetNX", err)
}

// Del 删除键
func (c *sentinelClient) Del(ctx context.Context, keys ...string) (int64, error) {
	ctx, span := c.startTrace(ctx, opDel)
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	n, err := c.client.Del(ctx, keys...).Result()
	c.endTrace(span, err)
	return n, c.handleCommandError("Del", err)
}

// Exists 检查键是否存在
func (c *sentinelClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	ctx, span := c.startTrace(ctx, "redis.Exists")
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	n, err := c.client.Exists(ctx, keys...).Result()
	c.endTrace(span, err)
	return n, c.handleCommandError("Exists", err)
}

// 为了满足接口而添加的存根方法
func (c *sentinelClient) MGet(ctx context.Context, keys ...string) ([]KeyValue, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) MSet(ctx context.Context, pairs ...KeyValue) error {
	return ErrCommandFailed
}

func (c *sentinelClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return false, ErrCommandFailed
}

func (c *sentinelClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) HGet(ctx context.Context, key, field string) (string, error) {
	return "", ErrCommandFailed
}

func (c *sentinelClient) HGetWithExists(ctx context.Context, key, field string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *sentinelClient) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return false, ErrCommandFailed
}

func (c *sentinelClient) HKeys(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) HVals(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) HLen(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) LPop(ctx context.Context, key string) (string, error) {
	return "", ErrCommandFailed
}

func (c *sentinelClient) LPopWithExists(ctx context.Context, key string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *sentinelClient) RPop(ctx context.Context, key string) (string, error) {
	return "", ErrCommandFailed
}

func (c *sentinelClient) RPopWithExists(ctx context.Context, key string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *sentinelClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) LLen(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	return ErrCommandFailed
}

func (c *sentinelClient) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, ErrCommandFailed
}

func (c *sentinelClient) SCard(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) SInter(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) ZAdd(ctx context.Context, key string, members ...*Z) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*Z, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) ZCard(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) ZRank(ctx context.Context, key string, member string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *sentinelClient) TxPipeline() Pipeliner {
	return nil
}

func (c *sentinelClient) Watch(ctx context.Context, fn func(tx *Tx) error, keys ...string) error {
	return ErrCommandFailed
}

func (c *sentinelClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return nil, ErrCommandFailed
}

func (c *sentinelClient) ScriptLoad(ctx context.Context, script string) (string, error) {
	return "", ErrCommandFailed
}

func (c *sentinelClient) Subscribe(ctx context.Context, channels ...string) *PubSub {
	return &PubSub{channel: nil}
}

func (c *sentinelClient) PSubscribe(ctx context.Context, patterns ...string) *PubSub {
	return &PubSub{channel: nil}
}

func (c *sentinelClient) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	return 0, ErrCommandFailed
}
