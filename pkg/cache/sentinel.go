package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel/attribute"
)

// Close 关闭连接
func (c *sentinelClient) Close() error {
	return c.client.Close()
}

// Do 执行自定义命令
func (c *sentinelClient) Do(ctx context.Context, args ...interface{}) (interface{}, error) {
	ctx, span := c.startTrace(ctx, opDo)
	defer c.endTrace(ctx, span, nil)

	cmd := c.client.Do(ctx, args...)
	return cmd.Result()
}

// Get 获取键值
func (c *sentinelClient) Get(ctx context.Context, key string) (string, error) {
	ctx, span := c.startTrace(ctx, opGet, attribute.String("key", key))
	defer func() {
		if span != nil {
			span.End()
		}
	}()

	val, err := c.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			if span != nil {
				span.SetAttributes(attribute.Bool("exists", false))
			}
			return "", nil
		}
		c.endTrace(ctx, span, err)
		return "", err
	}

	if span != nil {
		span.SetAttributes(attribute.Bool("exists", true))
	}
	return val, nil
}

// Set 设置键值
func (c *sentinelClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	ctx, span := c.startTrace(ctx, opSet,
		attribute.String("key", key),
		attribute.Int64("expiration_ms", expiration.Milliseconds()))
	defer c.endTrace(ctx, span, nil)

	return c.client.Set(ctx, key, value, expiration).Err()
}

// SetNX 不存在时设置键值
func (c *sentinelClient) SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error) {
	return c.client.SetNX(ctx, key, value, expiration).Result()
}

// Del 删除键
func (c *sentinelClient) Del(ctx context.Context, keys ...string) (int64, error) {
	ctx, span := c.startTrace(ctx, opDel, attribute.StringSlice("keys", keys))
	defer c.endTrace(ctx, span, nil)

	return c.client.Del(ctx, keys...).Result()
}

// Exists 检查键是否存在
func (c *sentinelClient) Exists(ctx context.Context, keys ...string) (int64, error) {
	return c.client.Exists(ctx, keys...).Result()
}

// Expire 设置过期时间
func (c *sentinelClient) Expire(ctx context.Context, key string, expiration time.Duration) (bool, error) {
	return c.client.Expire(ctx, key, expiration).Result()
}

// TTL 获取剩余生存时间
func (c *sentinelClient) TTL(ctx context.Context, key string) (time.Duration, error) {
	return c.client.TTL(ctx, key).Result()
}

// HGet 获取哈希表字段值
func (c *sentinelClient) HGet(ctx context.Context, key, field string) (string, error) {
	val, err := c.client.HGet(ctx, key, field).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// HSet 设置哈希表字段值
func (c *sentinelClient) HSet(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return c.client.HSet(ctx, key, values...).Result()
}

// HGetAll 获取哈希表所有字段和值
func (c *sentinelClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	return c.client.HGetAll(ctx, key).Result()
}

// HDel 删除哈希表字段
func (c *sentinelClient) HDel(ctx context.Context, key string, fields ...string) (int64, error) {
	return c.client.HDel(ctx, key, fields...).Result()
}

// HExists 检查哈希表字段是否存在
func (c *sentinelClient) HExists(ctx context.Context, key, field string) (bool, error) {
	return c.client.HExists(ctx, key, field).Result()
}

// LPush 将值插入到列表头部
func (c *sentinelClient) LPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return c.client.LPush(ctx, key, values...).Result()
}

// RPush 将值插入到列表尾部
func (c *sentinelClient) RPush(ctx context.Context, key string, values ...interface{}) (int64, error) {
	return c.client.RPush(ctx, key, values...).Result()
}

// LPop 移除并返回列表第一个元素
func (c *sentinelClient) LPop(ctx context.Context, key string) (string, error) {
	val, err := c.client.LPop(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// RPop 移除并返回列表最后一个元素
func (c *sentinelClient) RPop(ctx context.Context, key string) (string, error) {
	val, err := c.client.RPop(ctx, key).Result()
	if err == redis.Nil {
		return "", nil
	}
	return val, err
}

// LRange 获取列表指定范围内的元素
func (c *sentinelClient) LRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.LRange(ctx, key, start, stop).Result()
}

// LLen 获取列表长度
func (c *sentinelClient) LLen(ctx context.Context, key string) (int64, error) {
	return c.client.LLen(ctx, key).Result()
}

// SAdd 向集合添加一个或多个成员
func (c *sentinelClient) SAdd(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.client.SAdd(ctx, key, members...).Result()
}

// SMembers 返回集合中的所有成员
func (c *sentinelClient) SMembers(ctx context.Context, key string) ([]string, error) {
	return c.client.SMembers(ctx, key).Result()
}

// SRem 移除集合中一个或多个成员
func (c *sentinelClient) SRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.client.SRem(ctx, key, members...).Result()
}

// SIsMember 判断成员是否是集合的成员
func (c *sentinelClient) SIsMember(ctx context.Context, key string, member interface{}) (bool, error) {
	return c.client.SIsMember(ctx, key, member).Result()
}

// SCard 获取集合的成员数
func (c *sentinelClient) SCard(ctx context.Context, key string) (int64, error) {
	return c.client.SCard(ctx, key).Result()
}

// ZAdd 向有序集合添加一个或多个成员
func (c *sentinelClient) ZAdd(ctx context.Context, key string, members ...*Z) (int64, error) {
	// 转换成redis.Z类型
	redisMembers := make([]redis.Z, 0, len(members))
	for _, m := range members {
		redisMembers = append(redisMembers, redis.Z{
			Score:  m.Score,
			Member: m.Member,
		})
	}
	return c.client.ZAdd(ctx, key, redisMembers...).Result()
}

// ZRange 返回有序集合中指定范围的成员
func (c *sentinelClient) ZRange(ctx context.Context, key string, start, stop int64) ([]string, error) {
	return c.client.ZRange(ctx, key, start, stop).Result()
}

// ZRangeWithScores 返回有序集合中指定范围的成员和分数
func (c *sentinelClient) ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*Z, error) {
	result, err := c.client.ZRangeWithScores(ctx, key, start, stop).Result()
	if err != nil {
		return nil, err
	}

	members := make([]*Z, 0, len(result))
	for _, z := range result {
		members = append(members, &Z{
			Score:  z.Score,
			Member: z.Member,
		})
	}
	return members, nil
}

// ZRem 移除有序集合中的一个或多个成员
func (c *sentinelClient) ZRem(ctx context.Context, key string, members ...interface{}) (int64, error) {
	return c.client.ZRem(ctx, key, members...).Result()
}

// ZCard 获取有序集合的成员数
func (c *sentinelClient) ZCard(ctx context.Context, key string) (int64, error) {
	return c.client.ZCard(ctx, key).Result()
}
