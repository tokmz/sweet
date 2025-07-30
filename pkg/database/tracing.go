package database

import (
	"context"
	"fmt"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
	"gorm.io/gorm"
)

const (
	// 插件名称
	pluginName = "otel:tracing"
	// Span名称前缀
	spanPrefix = "gorm"
)

// TracingPlugin OpenTelemetry追踪插件
type TracingPlugin struct {
	config TracingConfig
	tracer trace.Tracer
}

// NewTracingPlugin 创建追踪插件
func NewTracingPlugin(config TracingConfig) *TracingPlugin {
	tracer := otel.Tracer(config.ServiceName)
	return &TracingPlugin{
		config: config,
		tracer: tracer,
	}
}

// Name 返回插件名称
func (p *TracingPlugin) Name() string {
	return pluginName
}

// Initialize 初始化插件
func (p *TracingPlugin) Initialize(db *gorm.DB) error {
	if !p.config.Enabled {
		return nil
	}

	// 注册回调函数
	if err := p.registerCallbacks(db); err != nil {
		return fmt.Errorf("failed to register tracing callbacks: %w", err)
	}

	return nil
}

// registerCallbacks 注册回调函数
func (p *TracingPlugin) registerCallbacks(db *gorm.DB) error {
	// Create操作
	if err := db.Callback().Create().Before("gorm:create").Register("otel:before_create", p.beforeCreate); err != nil {
		return err
	}
	if err := db.Callback().Create().After("gorm:create").Register("otel:after_create", p.afterCreate); err != nil {
		return err
	}

	// Query操作
	if err := db.Callback().Query().Before("gorm:query").Register("otel:before_query", p.beforeQuery); err != nil {
		return err
	}
	if err := db.Callback().Query().After("gorm:query").Register("otel:after_query", p.afterQuery); err != nil {
		return err
	}

	// Update操作
	if err := db.Callback().Update().Before("gorm:update").Register("otel:before_update", p.beforeUpdate); err != nil {
		return err
	}
	if err := db.Callback().Update().After("gorm:update").Register("otel:after_update", p.afterUpdate); err != nil {
		return err
	}

	// Delete操作
	if err := db.Callback().Delete().Before("gorm:delete").Register("otel:before_delete", p.beforeDelete); err != nil {
		return err
	}
	if err := db.Callback().Delete().After("gorm:delete").Register("otel:after_delete", p.afterDelete); err != nil {
		return err
	}

	// Row操作
	if err := db.Callback().Row().Before("gorm:row").Register("otel:before_row", p.beforeRow); err != nil {
		return err
	}
	if err := db.Callback().Row().After("gorm:row").Register("otel:after_row", p.afterRow); err != nil {
		return err
	}

	// Raw操作
	if err := db.Callback().Raw().Before("gorm:raw").Register("otel:before_raw", p.beforeRaw); err != nil {
		return err
	}
	if err := db.Callback().Raw().After("gorm:raw").Register("otel:after_raw", p.afterRaw); err != nil {
		return err
	}

	return nil
}

// beforeCreate Create操作前回调
func (p *TracingPlugin) beforeCreate(db *gorm.DB) {
	p.before(db, "create")
}

// afterCreate Create操作后回调
func (p *TracingPlugin) afterCreate(db *gorm.DB) {
	p.after(db)
}

// beforeQuery Query操作前回调
func (p *TracingPlugin) beforeQuery(db *gorm.DB) {
	p.before(db, "query")
}

// afterQuery Query操作后回调
func (p *TracingPlugin) afterQuery(db *gorm.DB) {
	p.after(db)
}

// beforeUpdate Update操作前回调
func (p *TracingPlugin) beforeUpdate(db *gorm.DB) {
	p.before(db, "update")
}

// afterUpdate Update操作后回调
func (p *TracingPlugin) afterUpdate(db *gorm.DB) {
	p.after(db)
}

// beforeDelete Delete操作前回调
func (p *TracingPlugin) beforeDelete(db *gorm.DB) {
	p.before(db, "delete")
}

// afterDelete Delete操作后回调
func (p *TracingPlugin) afterDelete(db *gorm.DB) {
	p.after(db)
}

// beforeRow Row操作前回调
func (p *TracingPlugin) beforeRow(db *gorm.DB) {
	p.before(db, "row")
}

// afterRow Row操作后回调
func (p *TracingPlugin) afterRow(db *gorm.DB) {
	p.after(db)
}

// beforeRaw Raw操作前回调
func (p *TracingPlugin) beforeRaw(db *gorm.DB) {
	p.before(db, "raw")
}

// afterRaw Raw操作后回调
func (p *TracingPlugin) afterRaw(db *gorm.DB) {
	p.after(db)
}

// before 通用前置回调
func (p *TracingPlugin) before(db *gorm.DB, operation string) {
	ctx := db.Statement.Context
	if ctx == nil {
		ctx = context.Background()
	}

	spanName := fmt.Sprintf("%s.%s", spanPrefix, operation)
	ctx, span := p.tracer.Start(ctx, spanName)

	// 设置基本属性
	span.SetAttributes(
		attribute.String("db.system", "mysql"),
		attribute.String("db.operation", operation),
	)

	// 设置表名
	if db.Statement.Table != "" {
		span.SetAttributes(attribute.String("db.sql.table", db.Statement.Table))
	}

	// 存储span和开始时间
	db.Set("otel:span", span)
	db.Set("otel:start_time", time.Now())
	db.Statement.Context = ctx
}

// after 通用后置回调
func (p *TracingPlugin) after(db *gorm.DB) {
	spanValue, exists := db.Get("otel:span")
	if !exists {
		return
	}

	span, ok := spanValue.(trace.Span)
	if !ok {
		return
	}
	defer span.End()

	// 获取开始时间
	startTimeValue, _ := db.Get("otel:start_time")
	startTime, _ := startTimeValue.(time.Time)
	duration := time.Since(startTime)

	// 设置执行时间
	span.SetAttributes(attribute.Int64("db.duration_ms", duration.Milliseconds()))

	// 设置影响行数
	if db.RowsAffected >= 0 {
		span.SetAttributes(attribute.Int64("db.rows_affected", db.RowsAffected))
	}

	// 记录SQL语句
	if p.config.RecordSQL && db.Statement.SQL.String() != "" {
		span.SetAttributes(attribute.String("db.statement", db.Statement.SQL.String()))
	}

	// 记录查询参数
	if p.config.RecordParams && len(db.Statement.Vars) > 0 {
		params := make([]string, len(db.Statement.Vars))
		for i, v := range db.Statement.Vars {
			params[i] = fmt.Sprintf("%v", v)
		}
		span.SetAttributes(attribute.StringSlice("db.statement.params", params))
	}

	// 处理错误
	if db.Error != nil {
		span.RecordError(db.Error)
		span.SetStatus(codes.Error, db.Error.Error())
		span.SetAttributes(attribute.Bool("db.error", true))
	} else {
		span.SetStatus(codes.Ok, "")
		span.SetAttributes(attribute.Bool("db.error", false))
	}
}