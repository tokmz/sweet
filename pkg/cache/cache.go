package cache

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/trace"
)

// 基础客户端
type baseClient struct {
	enableTrace bool
	tracer      trace.Tracer
	execTimeout int
}

// 开始链路追踪
func (c *baseClient) startTrace(ctx context.Context, op string, attrs ...attribute.KeyValue) (context.Context, trace.Span) {
	if !c.enableTrace {
		return ctx, nil
	}
	return c.tracer.Start(ctx, op, trace.WithAttributes(attrs...))
}

// 结束链路追踪
func (c *baseClient) endTrace(span trace.Span, err error) {
	if span == nil {
		return
	}
	if err != nil {
		span.RecordError(err)
	}
	span.End()
}

// 处理超时
func (c *baseClient) withTimeout(ctx context.Context) (context.Context, context.CancelFunc) {
	if c.execTimeout > 0 {
		return context.WithTimeout(ctx, time.Duration(c.execTimeout)*time.Millisecond)
	}
	return ctx, func() {}
}

// 处理redis.Nil错误
func (c *baseClient) handleNilError(err error) error {
	if errors.Is(err, redis.Nil) {
		return ErrKeyNotExists
	}
	return err
}

// 处理命令错误
func (c *baseClient) handleCommandError(cmd string, err error) error {
	if err == nil {
		return nil
	}
	if errors.Is(err, redis.Nil) {
		return ErrKeyNotExists
	}
	return &CommandError{
		Command: cmd,
		Err:     err,
	}
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
		return nil, ErrInvalidMode
	}
}

// 验证配置
func validateConfig(config Config) error {
	switch config.Mode {
	case ModeSingle:
		if config.Single.Addr == "" {
			return ErrEmptyAddrs
		}
	case ModeCluster:
		if len(config.Cluster.Addrs) == 0 {
			return ErrEmptyAddrs
		}
	case ModeSentinel:
		if len(config.Sentinel.Addrs) == 0 {
			return ErrEmptyAddrs
		}
		if config.Sentinel.MasterName == "" {
			return ErrEmptyMasterSet
		}
	default:
		return ErrInvalidMode
	}
	return nil
}

// 设置默认配置
func setDefaultConfig(config *Config) {
	// 连接超时
	if config.ConnTimeout <= 0 {
		config.ConnTimeout = defaultConnTimeout
	}
	// 读取超时
	if config.ReadTimeout <= 0 {
		config.ReadTimeout = defaultReadTimeout
	}
	// 写入超时
	if config.WriteTimeout <= 0 {
		config.WriteTimeout = defaultWriteTimeout
	}
	// 最大重试次数
	if config.MaxRetries <= 0 {
		config.MaxRetries = defaultMaxRetries
	}
	// 重试间隔
	if config.RetryDelay <= 0 {
		config.RetryDelay = defaultRetryDelay
	}
	// 最小重试间隔
	if config.MinRetryBackoff <= 0 {
		config.MinRetryBackoff = defaultMinRetryBackoff
	}
	// 最大重试间隔
	if config.MaxRetryBackoff <= 0 {
		config.MaxRetryBackoff = defaultMaxRetryBackoff
	}
	// 连接池大小
	if config.PoolSize <= 0 {
		config.PoolSize = defaultPoolSize
	}
	// 最小空闲连接数
	if config.MinIdleConns <= 0 {
		config.MinIdleConns = defaultMinIdleConns
	}
	// 连接最大空闲时间
	if config.IdleTimeout <= 0 {
		config.IdleTimeout = defaultIdleTimeout
	}
}

// 单机客户端实现
type singleClient struct {
	client *redis.Client
	baseClient
}

// Close 关闭连接
func (c *singleClient) Close() error {
	return c.client.Close()
}

// 集群客户端实现
type clusterClient struct {
	client *redis.ClusterClient
	baseClient
}

// Close 关闭连接
func (c *clusterClient) Close() error {
	return c.client.Close()
}

// 哨兵客户端实现
type sentinelClient struct {
	client *redis.Client
	baseClient
}

// Close 关闭连接
func (c *sentinelClient) Close() error {
	return c.client.Close()
}

// 创建单机客户端
func newSingleClient(config Config) (Client, error) {
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
		return nil, &ConnectionError{
			Addr: config.Single.Addr,
			Err:  err,
		}
	}

	// 创建客户端实例
	c := &singleClient{
		client: client,
		baseClient: baseClient{
			enableTrace: config.EnableTrace,
			execTimeout: config.ExecTimeout,
		},
	}

	// 初始化链路追踪
	if config.EnableTrace {
		c.tracer = otel.Tracer(tracerName)
	}

	return c, nil
}

// 创建集群客户端
func newClusterClient(config Config) (Client, error) {
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
		return nil, &ConnectionError{
			Addr: fmt.Sprintf("%v", config.Cluster.Addrs),
			Err:  err,
		}
	}

	// 创建客户端实例
	c := &clusterClient{
		client: client,
		baseClient: baseClient{
			enableTrace: config.EnableTrace,
			execTimeout: config.ExecTimeout,
		},
	}

	// 初始化链路追踪
	if config.EnableTrace {
		c.tracer = otel.Tracer(tracerName)
	}

	return c, nil
}

// 创建哨兵客户端
func newSentinelClient(config Config) (Client, error) {
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
		return nil, &ConnectionError{
			Addr: fmt.Sprintf("%v (master: %s)", config.Sentinel.Addrs, config.Sentinel.MasterName),
			Err:  err,
		}
	}

	// 创建客户端实例
	c := &sentinelClient{
		client: client,
		baseClient: baseClient{
			enableTrace: config.EnableTrace,
			execTimeout: config.ExecTimeout,
		},
	}

	// 初始化链路追踪
	if config.EnableTrace {
		c.tracer = otel.Tracer(tracerName)
	}

	return c, nil
}
