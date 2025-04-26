package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
)

// Init 初始化数据库连接并设置为全局默认实例
func Init(config Config) (*gorm.DB, error) {
	// 参数验证
	if config.Master == "" {
		return nil, ErrInvalidConfig
	}

	// 加锁保护并发安全
	rw.Lock()
	defer rw.Unlock()

	// 创建数据库连接
	db, err := initDB(config)
	if err != nil {
		return nil, err
	}

	// 设置全局默认实例
	DefaultDB = db
	initialized = true

	return db, nil
}

// initDB 内部函数，实际初始化数据库连接
func initDB(config Config) (*gorm.DB, error) {
	// 创建Logger
	gormLogger := createLogger(config)

	// 连接主库
	gormConfig := &gorm.Config{
		Logger: gormLogger,
		NowFunc: func() time.Time {
			return time.Now().Local()
		},
	}

	// 设置查询超时
	if config.QueryTimeout > 0 {
		gormConfig.PrepareStmt = true
	}

	// 处理TiDB特有的配置
	masterDSN := config.Master
	if config.TiDB.Enabled {
		masterDSN = processTiDBDSN(config.Master, config.TiDB)
	}

	// 连接主库
	db, err := gorm.Open(mysql.Open(masterDSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err.Error())
	}

	// 配置读写分离
	if len(config.Slaves) > 0 {
		resolverConfig := dbresolver.Config{
			Sources:  []gorm.Dialector{mysql.Open(masterDSN)},
			Replicas: make([]gorm.Dialector, 0, len(config.Slaves)),
		}

		for _, slave := range config.Slaves {
			slaveDSN := slave
			if config.TiDB.Enabled {
				slaveDSN = processTiDBDSN(slave, config.TiDB)
			}
			resolverConfig.Replicas = append(resolverConfig.Replicas, mysql.Open(slaveDSN))
		}

		// 设置读写分离并配置连接池
		dbResolverPlugin := dbresolver.Register(resolverConfig).
			SetConnMaxLifetime(time.Duration(config.Pool.ConnMaxLifetime) * time.Second).
			SetMaxIdleConns(config.Pool.MaxIdleConns).
			SetMaxOpenConns(config.Pool.MaxOpenConns)

		if err = db.Use(dbResolverPlugin); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrConfigReadWriteSplit, err.Error())
		}
	}

	// 配置连接池
	sqlDB, err := db.DB()
	if err != nil {
		return nil, fmt.Errorf("%w: %s", ErrGetSQLDB, err.Error())
	}

	sqlDB.SetMaxIdleConns(config.Pool.MaxIdleConns)
	sqlDB.SetMaxOpenConns(config.Pool.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Duration(config.Pool.ConnMaxLifetime) * time.Second)

	// 配置链路追踪
	if config.EnableTrace {
		if err := db.Use(&TracingPlugin{
			tracer:             otel.Tracer(tracerName),
			recordSQL:          true,
			recordAffectedRows: true,
		}); err != nil {
			return nil, fmt.Errorf("%w: %s", ErrConfigTracing, err.Error())
		}
	}

	return db, nil
}

// processTiDBDSN 处理TiDB特有的DSN参数
func processTiDBDSN(dsn string, tidbConfig TiDBConfig) string {
	// 如果DSN中已包含tls参数，则不再添加
	if strings.Contains(dsn, "tls=") {
		return dsn
	}

	// 添加TLS配置
	tlsParam := "false"
	if tidbConfig.UseSSL {
		tlsParam = "true"
		// 如果指定了CA证书，可以在这里添加更复杂的TLS配置
	}

	// 检查DSN是否已有参数
	if strings.Contains(dsn, "?") {
		dsn = fmt.Sprintf("%s&tls=%s", dsn, tlsParam)
	} else {
		dsn = fmt.Sprintf("%s?tls=%s", dsn, tlsParam)
	}

	// 添加其他TiDB特有参数
	if tidbConfig.EnableAutoRandom {
		dsn = fmt.Sprintf("%s&allow_auto_random=true", dsn)
	}

	return dsn
}

// 创建GORM日志记录器
func createLogger(config Config) logger.Interface {
	var logLevel = logger.Silent
	if config.EnableLog {
		logLevel = logger.Info
	}

	return logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             time.Duration(config.SlowThreshold) * time.Millisecond,
			LogLevel:                  logLevel,
			IgnoreRecordNotFoundError: true,
			Colorful:                  false,
		},
	)
}

// GetDB 安全地获取默认数据库实例
func GetDB() (*gorm.DB, error) {
	rw.RLock()
	defer rw.RUnlock()

	if !initialized || DefaultDB == nil {
		return nil, ErrNotInitialized
	}

	return DefaultDB, nil
}

// WithContext 设置上下文
func WithContext(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.WithContext(ctx)
}

// MasterDB 强制使用主库
func MasterDB(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.WithContext(ctx).Clauses(dbresolver.Write)
}

// SlaveDB 强制使用从库
func SlaveDB(db *gorm.DB, ctx context.Context) *gorm.DB {
	return db.WithContext(ctx).Clauses(dbresolver.Read)
}

// 实现带重试功能的事务
func Transaction(db *gorm.DB, ctx context.Context, fn func(tx *gorm.DB) error) error {
	return TransactionWithRetry(db, ctx, fn, 0)
}

// TransactionWithRetry 带重试功能的事务
func TransactionWithRetry(db *gorm.DB, ctx context.Context, fn func(tx *gorm.DB) error, retries int) error {
	var err error

	// 如果设置为0，使用默认配置的最大重试次数
	maxRetries := retries
	if maxRetries <= 0 {
		// 从存储在上下文中的配置获取，如果没有则默认为3次
		cfg, ok := ctx.Value("db_config").(Config)
		if ok {
			maxRetries = cfg.MaxRetries
		} else {
			maxRetries = 3
		}
	}

	// 从存储在上下文中的配置获取重试延迟
	var retryDelay time.Duration = 100 * time.Millisecond
	cfg, ok := ctx.Value("db_config").(Config)
	if ok && cfg.RetryDelay > 0 {
		retryDelay = time.Duration(cfg.RetryDelay) * time.Millisecond
	}

	// 获取链路追踪的span
	span := trace.SpanFromContext(ctx)

	for attempt := 0; attempt <= maxRetries; attempt++ {
		// 添加重试信息到span
		if span.IsRecording() {
			span.SetAttributes(
				attribute.Int("db.retry.attempt", attempt),
				attribute.Int("db.retry.max", maxRetries),
			)
		}

		err = db.WithContext(ctx).Transaction(func(tx *gorm.DB) error {
			return fn(tx)
		})

		// 如果成功或者是非临时性错误，不需要重试
		if err == nil || !isRetryableError(err) {
			break
		}

		// 最后一次尝试失败，直接返回错误
		if attempt == maxRetries {
			if span.IsRecording() {
				span.SetAttributes(attribute.Bool("db.retry.exhausted", true))
			}
			break
		}

		// 记录重试信息
		if span.IsRecording() {
			span.AddEvent("db.retry", trace.WithAttributes(
				attribute.String("db.error", err.Error()),
				attribute.Int("db.retry.delay_ms", int(retryDelay.Milliseconds())),
			))
		}

		// 等待一段时间后重试
		time.Sleep(retryDelay)
	}

	return err
}

// isRetryableError 判断错误是否可重试
func isRetryableError(err error) bool {
	if err == nil {
		return false
	}

	// 扩展可重试错误类型，包括TiDB特有的错误
	errMsg := err.Error()
	return errors.Is(err, gorm.ErrInvalidTransaction) ||
		contains(errMsg, "deadlock") ||
		contains(errMsg, "Lock wait timeout") ||
		contains(errMsg, "connection reset") ||
		contains(errMsg, "server has gone away") ||
		contains(errMsg, "too many connections") ||
		// TiDB特有的可重试错误
		contains(errMsg, "Region is unavailable") ||
		contains(errMsg, "TiKV server timeout") ||
		contains(errMsg, "transaction") && contains(errMsg, "abort") ||
		contains(errMsg, "retry statement")
}

// contains 判断字符串是否包含子串
func contains(s, substr string) bool {
	return strings.Contains(s, substr)
}

// TracingPlugin 链路追踪插件
type TracingPlugin struct {
	tracer             trace.Tracer
	recordSQL          bool
	recordAffectedRows bool
}

// Name 返回插件名称
func (tp *TracingPlugin) Name() string {
	return "TracingPlugin"
}

// Initialize 初始化并添加回调
func (tp *TracingPlugin) Initialize(db *gorm.DB) error {
	// 为Create操作注册回调
	err := db.Callback().Create().Before("gorm:create").Register("tracing:before_create", tp.beforeCreate)
	if err != nil {
		return fmt.Errorf("注册Create前回调失败: %w", err)
	}
	err = db.Callback().Create().After("gorm:create").Register("tracing:after_create", tp.afterCreate)
	if err != nil {
		return fmt.Errorf("注册Create后回调失败: %w", err)
	}

	// 为Query操作注册回调
	err = db.Callback().Query().Before("gorm:query").Register("tracing:before_query", tp.beforeQuery)
	if err != nil {
		return fmt.Errorf("注册Query前回调失败: %w", err)
	}
	err = db.Callback().Query().After("gorm:query").Register("tracing:after_query", tp.afterQuery)
	if err != nil {
		return fmt.Errorf("注册Query后回调失败: %w", err)
	}

	// 为Update操作注册回调
	err = db.Callback().Update().Before("gorm:update").Register("tracing:before_update", tp.beforeUpdate)
	if err != nil {
		return fmt.Errorf("注册Update前回调失败: %w", err)
	}
	err = db.Callback().Update().After("gorm:update").Register("tracing:after_update", tp.afterUpdate)
	if err != nil {
		return fmt.Errorf("注册Update后回调失败: %w", err)
	}

	// 为Delete操作注册回调
	err = db.Callback().Delete().Before("gorm:delete").Register("tracing:before_delete", tp.beforeDelete)
	if err != nil {
		return fmt.Errorf("注册Delete前回调失败: %w", err)
	}
	err = db.Callback().Delete().After("gorm:delete").Register("tracing:after_delete", tp.afterDelete)
	if err != nil {
		return fmt.Errorf("注册Delete后回调失败: %w", err)
	}

	// 为Raw操作注册回调
	err = db.Callback().Raw().Before("gorm:raw").Register("tracing:before_raw", tp.beforeRaw)
	if err != nil {
		return fmt.Errorf("注册Raw前回调失败: %w", err)
	}
	err = db.Callback().Raw().After("gorm:raw").Register("tracing:after_raw", tp.afterRaw)
	if err != nil {
		return fmt.Errorf("注册Raw后回调失败: %w", err)
	}

	return nil
}

// startSpan 开始span
func (tp *TracingPlugin) startSpan(ctx context.Context, operation string, db *gorm.DB) (context.Context, trace.Span) {
	if ctx == nil {
		ctx = context.Background()
	}

	opts := []trace.SpanStartOption{
		trace.WithAttributes(
			attribute.String("db.system", "mysql"),
			attribute.String("db.operation", operation),
		),
		trace.WithSpanKind(trace.SpanKindClient),
	}

	// 记录SQL语句
	if tp.recordSQL && db.Statement != nil && db.Statement.SQL.String() != "" {
		opts = append(opts, trace.WithAttributes(attribute.String("db.statement", db.Statement.SQL.String())))
	}

	// 记录表名
	if db.Statement != nil && db.Statement.Table != "" {
		opts = append(opts, trace.WithAttributes(attribute.String("db.table", db.Statement.Table)))
	}

	return tp.tracer.Start(ctx, operation, opts...)
}

// endSpan 结束span
func (tp *TracingPlugin) endSpan(span trace.Span, db *gorm.DB) {
	if db.Error != nil && db.Error != gorm.ErrRecordNotFound {
		span.SetStatus(codes.Error, db.Error.Error())
		span.RecordError(db.Error)
	} else {
		span.SetStatus(codes.Ok, "")
	}

	// 记录影响的行数
	if tp.recordAffectedRows && db.Statement != nil && db.Statement.SQL.String() != "" {
		span.SetAttributes(attribute.Int64("db.rows_affected", db.Statement.RowsAffected))
	}

	span.End()
}

// 各种操作回调
func (tp *TracingPlugin) beforeCreate(db *gorm.DB) {
	ctx, span := tp.startSpan(db.Statement.Context, opCreate, db)
	db.Statement.Context = ctx
	db.Statement.WithContext(ctx)
	db.Statement.Context = context.WithValue(db.Statement.Context, spanKey, span)
}

func (tp *TracingPlugin) afterCreate(db *gorm.DB) {
	span, ok := db.Statement.Context.Value(spanKey).(trace.Span)
	if ok {
		tp.endSpan(span, db)
	}
}

func (tp *TracingPlugin) beforeQuery(db *gorm.DB) {
	ctx, span := tp.startSpan(db.Statement.Context, opQuery, db)
	db.Statement.Context = ctx
	db.Statement.WithContext(ctx)
	db.Statement.Context = context.WithValue(db.Statement.Context, spanKey, span)
}

func (tp *TracingPlugin) afterQuery(db *gorm.DB) {
	span, ok := db.Statement.Context.Value(spanKey).(trace.Span)
	if ok {
		tp.endSpan(span, db)
	}
}

func (tp *TracingPlugin) beforeUpdate(db *gorm.DB) {
	ctx, span := tp.startSpan(db.Statement.Context, opUpdate, db)
	db.Statement.Context = ctx
	db.Statement.WithContext(ctx)
	db.Statement.Context = context.WithValue(db.Statement.Context, spanKey, span)
}

func (tp *TracingPlugin) afterUpdate(db *gorm.DB) {
	span, ok := db.Statement.Context.Value(spanKey).(trace.Span)
	if ok {
		tp.endSpan(span, db)
	}
}

func (tp *TracingPlugin) beforeDelete(db *gorm.DB) {
	ctx, span := tp.startSpan(db.Statement.Context, opDelete, db)
	db.Statement.Context = ctx
	db.Statement.WithContext(ctx)
	db.Statement.Context = context.WithValue(db.Statement.Context, spanKey, span)
}

func (tp *TracingPlugin) afterDelete(db *gorm.DB) {
	span, ok := db.Statement.Context.Value(spanKey).(trace.Span)
	if ok {
		tp.endSpan(span, db)
	}
}

func (tp *TracingPlugin) beforeRaw(db *gorm.DB) {
	ctx, span := tp.startSpan(db.Statement.Context, opRawSQL, db)
	db.Statement.Context = ctx
	db.Statement.WithContext(ctx)
	db.Statement.Context = context.WithValue(db.Statement.Context, spanKey, span)
}

func (tp *TracingPlugin) afterRaw(db *gorm.DB) {
	span, ok := db.Statement.Context.Value(spanKey).(trace.Span)
	if ok {
		tp.endSpan(span, db)
	}
}
