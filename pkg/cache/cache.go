package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/trace"
)

// Client Redis客户端接口
type Client interface {
	// 关闭连接
	Close() error
	// 执行自定义命令
	Do(ctx context.Context, args ...interface{}) (interface{}, error)

	// String操作
	Get(ctx context.Context, key string) (string, error)
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	SetNX(ctx context.Context, key string, value interface{}, expiration time.Duration) (bool, error)
	Del(ctx context.Context, keys ...string) (int64, error)
	Exists(ctx context.Context, keys ...string) (int64, error)
	Expire(ctx context.Context, key string, expiration time.Duration) (bool, error)
	TTL(ctx context.Context, key string) (time.Duration, error)

	// Hash操作
	HGet(ctx context.Context, key, field string) (string, error)
	HSet(ctx context.Context, key string, values ...interface{}) (int64, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HDel(ctx context.Context, key string, fields ...string) (int64, error)
	HExists(ctx context.Context, key, field string) (bool, error)

	// List操作
	LPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	RPush(ctx context.Context, key string, values ...interface{}) (int64, error)
	LPop(ctx context.Context, key string) (string, error)
	RPop(ctx context.Context, key string) (string, error)
	LRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	LLen(ctx context.Context, key string) (int64, error)

	// Set操作
	SAdd(ctx context.Context, key string, members ...interface{}) (int64, error)
	SMembers(ctx context.Context, key string) ([]string, error)
	SRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	SIsMember(ctx context.Context, key string, member interface{}) (bool, error)
	SCard(ctx context.Context, key string) (int64, error)

	// ZSet操作
	ZAdd(ctx context.Context, key string, members ...*Z) (int64, error)
	ZRange(ctx context.Context, key string, start, stop int64) ([]string, error)
	ZRangeWithScores(ctx context.Context, key string, start, stop int64) ([]*Z, error)
	ZRem(ctx context.Context, key string, members ...interface{}) (int64, error)
	ZCard(ctx context.Context, key string) (int64, error)
}

// Z 有序集合成员
type Z struct {
	Score  float64
	Member interface{}
}

// 基础客户端
type baseClient struct {
	enableTrace bool
	tracer      trace.Tracer
}

// 单机模式客户端
type singleClient struct {
	baseClient
	client *redis.Client
}

// 集群模式客户端
type clusterClient struct {
	baseClient
	client *redis.ClusterClient
}

// 哨兵模式客户端
type sentinelClient struct {
	baseClient
	client *redis.Client
}

// NewClient 创建Redis客户端
func NewClient(config Config) (Client, error) {
	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, err
	}

	// 设置默认值
	setDefaultConfig(&config)

	// 根据模式创建对应客户端
	switch config.Mode {
	case ModeSingle:
		return newSingleClient(config)
	case ModeCluster:
		return newClusterClient(config)
	case ModeSentinel:
		return newSentinelClient(config)
	default:
		return nil, errors.New(ErrInvalidMode)
	}
}

// 验证配置
func validateConfig(config Config) error {
	switch config.Mode {
	case ModeSingle:
		if config.Single.Addr == "" {
			return errors.New(ErrEmptyAddrs)
		}
	case ModeCluster:
		if len(config.Cluster.Addrs) == 0 {
			return errors.New(ErrEmptyAddrs)
		}
	case ModeSentinel:
		if len(config.Sentinel.Addrs) == 0 {
			return errors.New(ErrEmptyAddrs)
		}
		if config.Sentinel.MasterName == "" {
			return errors.New(ErrEmptyMasterSet)
		}
	default:
		return errors.New(ErrInvalidMode)
	}
	return nil
}

// 设置默认值
func setDefaultConfig(config *Config) {
	// 连接池配置
	if config.PoolSize == 0 {
		config.PoolSize = defaultPoolSize
	}
	if config.MinIdleConns == 0 {
		config.MinIdleConns = defaultMinIdleConns
	}
	if config.IdleTimeout == 0 {
		config.IdleTimeout = defaultIdleTimeout
	}

	// 超时设置
	if config.ConnTimeout == 0 {
		config.ConnTimeout = defaultConnTimeout
	}
	if config.ReadTimeout == 0 {
		config.ReadTimeout = defaultReadTimeout
	}
	if config.WriteTimeout == 0 {
		config.WriteTimeout = defaultWriteTimeout
	}

	// 重试设置
	if config.MaxRetries == 0 {
		config.MaxRetries = defaultMaxRetries
	}
	if config.RetryDelay == 0 {
		config.RetryDelay = defaultRetryDelay
	}
	if config.MinRetryBackoff == 0 {
		config.MinRetryBackoff = defaultMinRetryBackoff
	}
	if config.MaxRetryBackoff == 0 {
		config.MaxRetryBackoff = defaultMaxRetryBackoff
	}
}

// 创建单机客户端
func newSingleClient(config Config) (*singleClient, error) {
	// 创建连接选项
	opt := &redis.Options{
		Addr:            config.Single.Addr,
		Username:        config.Username,
		Password:        config.Password,
		DB:              config.DB,
		MaxRetries:      config.MaxRetries,
		MinRetryBackoff: time.Duration(config.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff: time.Duration(config.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:     time.Duration(config.ConnTimeout) * time.Millisecond,
		ReadTimeout:     time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout:    time.Duration(config.WriteTimeout) * time.Millisecond,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		MaxIdleConns:    config.PoolSize,
		ConnMaxIdleTime: time.Duration(config.IdleTimeout) * time.Second,
	}

	// 创建客户端
	client := redis.NewClient(opt)

	// 检查连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接Redis失败: %w", err)
	}

	// 创建客户端实例
	c := &singleClient{
		client: client,
		baseClient: baseClient{
			enableTrace: config.EnableTrace,
		},
	}

	// 初始化链路追踪
	if config.EnableTrace {
		c.tracer = otel.Tracer(tracerName)
	}

	return c, nil
}

// 创建集群客户端
func newClusterClient(config Config) (*clusterClient, error) {
	// 创建连接选项
	opt := &redis.ClusterOptions{
		Addrs:           config.Cluster.Addrs,
		Username:        config.Username,
		Password:        config.Password,
		MaxRetries:      config.MaxRetries,
		MinRetryBackoff: time.Duration(config.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff: time.Duration(config.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:     time.Duration(config.ConnTimeout) * time.Millisecond,
		ReadTimeout:     time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout:    time.Duration(config.WriteTimeout) * time.Millisecond,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		MaxIdleConns:    config.PoolSize,
		ConnMaxIdleTime: time.Duration(config.IdleTimeout) * time.Second,
	}

	// 创建客户端
	client := redis.NewClusterClient(opt)

	// 检查连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接Redis集群失败: %w", err)
	}

	// 创建客户端实例
	c := &clusterClient{
		client: client,
		baseClient: baseClient{
			enableTrace: config.EnableTrace,
		},
	}

	// 初始化链路追踪
	if config.EnableTrace {
		c.tracer = otel.Tracer(tracerName)
	}

	return c, nil
}

// 创建哨兵客户端
func newSentinelClient(config Config) (*sentinelClient, error) {
	// 创建连接选项
	opt := &redis.FailoverOptions{
		MasterName:      config.Sentinel.MasterName,
		SentinelAddrs:   config.Sentinel.Addrs,
		Username:        config.Username,
		Password:        config.Password,
		DB:              config.DB,
		MaxRetries:      config.MaxRetries,
		MinRetryBackoff: time.Duration(config.MinRetryBackoff) * time.Millisecond,
		MaxRetryBackoff: time.Duration(config.MaxRetryBackoff) * time.Millisecond,
		DialTimeout:     time.Duration(config.ConnTimeout) * time.Millisecond,
		ReadTimeout:     time.Duration(config.ReadTimeout) * time.Millisecond,
		WriteTimeout:    time.Duration(config.WriteTimeout) * time.Millisecond,
		PoolSize:        config.PoolSize,
		MinIdleConns:    config.MinIdleConns,
		MaxIdleConns:    config.PoolSize,
		ConnMaxIdleTime: time.Duration(config.IdleTimeout) * time.Second,
	}

	// 创建客户端
	client := redis.NewFailoverClient(opt)

	// 检查连接
	ctx := context.Background()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("连接Redis哨兵失败: %w", err)
	}

	// 创建客户端实例
	c := &sentinelClient{
		client: client,
		baseClient: baseClient{
			enableTrace: config.EnableTrace,
		},
	}

	// 初始化链路追踪
	if config.EnableTrace {
		c.tracer = otel.Tracer(tracerName)
	}

	return c, nil
}
