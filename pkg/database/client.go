package database

import (
	"context"
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// Client 数据库客户端
type Client struct {
	db     *gorm.DB
	config *Config
}

// NewClient 创建数据库客户端
func NewClient(config *Config) (*Client, error) {
	if config == nil {
		config = DefaultConfig()
	}

	// 验证配置
	if err := validateConfig(config); err != nil {
		return nil, fmt.Errorf("invalid config: %w", err)
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		Logger: NewCustomLogger(config.Log, config.SlowQuery),
	}

	// 连接主库
	db, err := gorm.Open(mysql.Open(config.Master), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to master database: %w", err)
	}

	// 配置读写分离
	if err := configureDBResolver(db, config); err != nil {
		return nil, fmt.Errorf("failed to configure db resolver: %w", err)
	}

	// 配置连接池
	if err := configureConnectionPool(db, config.Pool); err != nil {
		return nil, fmt.Errorf("failed to configure connection pool: %w", err)
	}

	// 安装OpenTelemetry插件
	if config.Tracing.Enabled {
		tracingPlugin := NewTracingPlugin(config.Tracing)
		if err := db.Use(tracingPlugin); err != nil {
			return nil, fmt.Errorf("failed to install tracing plugin: %w", err)
		}
	}

	client := &Client{
		db:     db,
		config: config,
	}

	return client, nil
}

// validateConfig 验证配置
func validateConfig(config *Config) error {
	if config.Master == "" {
		return fmt.Errorf("master DSN is required")
	}
	return nil
}

// configureDBResolver 配置读写分离
func configureDBResolver(db *gorm.DB, config *Config) error {
	if len(config.Slaves) == 0 {
		// 没有从库配置，跳过读写分离设置
		return nil
	}

	// 准备从库连接
	replicas := make([]gorm.Dialector, 0, len(config.Slaves))
	for _, slaveDSN := range config.Slaves {
		replicas = append(replicas, mysql.Open(slaveDSN))
	}

	// 配置dbresolver插件
	resolverConfig := dbresolver.Config{
		// 从库用于读操作
		Replicas: replicas,
		// 读写分离策略
		Policy: dbresolver.RandomPolicy{},
	}

	// 安装dbresolver插件
	if err := db.Use(dbresolver.Register(resolverConfig)); err != nil {
		return fmt.Errorf("failed to register dbresolver: %w", err)
	}

	return nil
}

// configureConnectionPool 配置连接池
func configureConnectionPool(db *gorm.DB, poolConfig PoolConfig) error {
	sqlDB, err := db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	// 设置连接池参数
	sqlDB.SetMaxIdleConns(poolConfig.MaxIdleConns)
	sqlDB.SetMaxOpenConns(poolConfig.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(poolConfig.ConnMaxLifetime)
	sqlDB.SetConnMaxIdleTime(poolConfig.ConnMaxIdleTime)

	return nil
}

// DB 获取GORM数据库实例
func (c *Client) DB() *gorm.DB {
	return c.db
}

// Close 关闭数据库连接
func (c *Client) Close() error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.Close()
}

// Ping 测试数据库连接
func (c *Client) Ping(ctx context.Context) error {
	sqlDB, err := c.db.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}
	return sqlDB.PingContext(ctx)
}

// Stats 获取数据库连接统计信息
func (c *Client) Stats() (map[string]interface{}, error) {
	sqlDB, err := c.db.DB()
	if err != nil {
		return nil, fmt.Errorf("failed to get sql.DB: %w", err)
	}

	stats := sqlDB.Stats()
	return map[string]interface{}{
		"max_open_connections":     stats.MaxOpenConnections,
		"open_connections":         stats.OpenConnections,
		"in_use":                   stats.InUse,
		"idle":                     stats.Idle,
		"wait_count":               stats.WaitCount,
		"wait_duration":            stats.WaitDuration,
		"max_idle_closed":          stats.MaxIdleClosed,
		"max_idle_time_closed":     stats.MaxIdleTimeClosed,
		"max_lifetime_closed":      stats.MaxLifetimeClosed,
	}, nil
}

// Transaction 执行事务
func (c *Client) Transaction(ctx context.Context, fn func(*gorm.DB) error) error {
	return c.db.WithContext(ctx).Transaction(fn)
}

// WithContext 设置上下文
func (c *Client) WithContext(ctx context.Context) *gorm.DB {
	return c.db.WithContext(ctx)
}

// Master 强制使用主库
func (c *Client) Master() *gorm.DB {
	return c.db.Clauses(dbresolver.Write)
}

// Slave 强制使用从库
func (c *Client) Slave() *gorm.DB {
	return c.db.Clauses(dbresolver.Read)
}

// GetConfig 获取配置
func (c *Client) GetConfig() *Config {
	return c.config
}

// SetLogLevel 动态设置日志级别
func (c *Client) SetLogLevel(level string) error {
	var logLevel logger.LogLevel
	switch level {
	case "silent":
		logLevel = logger.Silent
	case "error":
		logLevel = logger.Error
	case "warn":
		logLevel = logger.Warn
	case "info":
		logLevel = logger.Info
	default:
		return fmt.Errorf("invalid log level: %s", level)
	}

	c.db.Logger = c.db.Logger.LogMode(logLevel)
	return nil
}

// HealthCheck 健康检查
func (c *Client) HealthCheck(ctx context.Context) error {
	// 检查主库连接
	if err := c.Ping(ctx); err != nil {
		return fmt.Errorf("master database health check failed: %w", err)
	}

	// 检查从库连接（如果有配置）
	if len(c.config.Slaves) > 0 {
		// 尝试执行一个简单的读操作来验证从库连接
		var count int64
		if err := c.Slave().WithContext(ctx).Raw("SELECT 1").Count(&count).Error; err != nil {
			return fmt.Errorf("slave database health check failed: %w", err)
		}
	}

	return nil
}

// GetSlowQueries 获取慢查询统计（需要数据库支持）
func (c *Client) GetSlowQueries(ctx context.Context, limit int) ([]map[string]interface{}, error) {
	var results []map[string]interface{}
	
	// 查询MySQL慢查询日志表（需要开启慢查询日志）
	query := `
		SELECT 
			start_time,
			user_host,
			query_time,
			lock_time,
			rows_sent,
			rows_examined,
			db,
			sql_text
		FROM mysql.slow_log 
		ORDER BY start_time DESC 
		LIMIT ?
	`
	
	if err := c.db.WithContext(ctx).Raw(query, limit).Scan(&results).Error; err != nil {
		return nil, fmt.Errorf("failed to get slow queries: %w", err)
	}
	
	return results, nil
}