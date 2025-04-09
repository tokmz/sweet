package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
)

// Do 执行自定义命令
func (c *clusterClient) Do(ctx context.Context, args ...interface{}) (interface{}, error) {
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
func (c *clusterClient) Get(ctx context.Context, key string) (string, error) {
	ctx, span := c.startTrace(ctx, opGet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	val, err := c.client.Get(ctx, key).Result()
	c.endTrace(span, err)
	return val, c.handleCommandError("Get", err)
}

// GetWithExists 获取键值，返回值和键是否存在
func (c *clusterClient) GetWithExists(ctx context.Context, key string) (string, bool, error) {
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
func (c *clusterClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, span := c.startTrace(ctx, opSet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	err := c.client.Set(ctx, key, value, expiration).Err()
	c.endTrace(span, err)
	return c.handleCommandError("Set", err)
}

// SetNX 不存在时设置键值
func (c *clusterClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	ctx, span := c.startTrace(ctx, opSet, attribute.String("key", key))
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	ok, err := c.client.SetNX(ctx, key, value, expiration).Result()
	c.endTrace(span, err)
	return ok, c.handleCommandError("SetNX", err)
}

// Del 删除键
func (c *clusterClient) Del(ctx context.Context, keys ...string) (int64, error) {
	ctx, span := c.startTrace(ctx, opDel)
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	n, err := c.client.Del(ctx, keys...).Result()
	c.endTrace(span, err)
	return n, c.handleCommandError("Del", err)
}

// Exists 检查键是否存在
func (c *clusterClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	ctx, span := c.startTrace(ctx, "redis.Exists")
	defer c.endTrace(span, nil)

	ctx, cancel := c.withTimeout(ctx)
	defer cancel()

	n, err := c.client.Exists(ctx, keys...).Result()
	c.endTrace(span, err)
	return n, c.handleCommandError("Exists", err)
}

// 为了满足接口而添加的存根方法
func (c *clusterClient) MGet(ctx context.Context, keys ...string) ([]KeyValue, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) MSet(ctx context.Context, pairs ...KeyValue) error {
	return ErrCommandFailed
}

func (c *clusterClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return false, ErrCommandFailed
}

func (c *clusterClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) HGet(ctx context.Context, key, field string) (string, error) {
	return "", ErrCommandFailed
}

func (c *clusterClient) HGetWithExists(ctx context.Context, key, field string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *clusterClient) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return false, ErrCommandFailed
}

func (c *clusterClient) HKeys(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) HVals(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) HLen(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) LPop(ctx context.Context, key string) (string, error) {
	return "", ErrCommandFailed
}

func (c *clusterClient) LPopWithExists(ctx context.Context, key string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *clusterClient) RPop(ctx context.Context, key string) (string, error) {
	return "", ErrCommandFailed
}

func (c *clusterClient) RPopWithExists(ctx context.Context, key string) (string, bool, error) {
	return "", false, ErrCommandFailed
}

func (c *clusterClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) LLen(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) LTrim(ctx context.Context, key string, start, stop int64) error {
	return ErrCommandFailed
}

func (c *clusterClient) LRem(ctx context.Context, key string, count int64, value interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return false, ErrCommandFailed
}

func (c *clusterClient) SCard(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) SDiff(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) SInter(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) SUnion(ctx context.Context, keys ...string) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) ZAdd(ctx context.Context, key string, members ...*Z) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*Z, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) ZCard(ctx context.Context, key string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) ZScore(ctx context.Context, key string, member string) (float64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) ZRank(ctx context.Context, key string, member string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) ZRevRank(ctx context.Context, key string, member string) (int64, error) {
	return 0, ErrCommandFailed
}

func (c *clusterClient) TxPipeline() Pipeliner {
	return nil
}

func (c *clusterClient) Watch(ctx context.Context, fn func(tx *Tx) error, keys ...string) error {
	return ErrCommandFailed
}

func (c *clusterClient) Eval(ctx context.Context, script string, keys []string, args ...interface{}) (interface{}, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) EvalSha(ctx context.Context, sha1 string, keys []string, args ...interface{}) (interface{}, error) {
	return nil, ErrCommandFailed
}

func (c *clusterClient) ScriptLoad(ctx context.Context, script string) (string, error) {
	return "", ErrCommandFailed
}

func (c *clusterClient) Subscribe(ctx context.Context, channels ...string) *PubSub {
	return &PubSub{channel: nil}
}

func (c *clusterClient) PSubscribe(ctx context.Context, patterns ...string) *PubSub {
	return &PubSub{channel: nil}
}

func (c *clusterClient) Publish(ctx context.Context, channel string, message interface{}) (int64, error) {
	return 0, ErrCommandFailed
}
