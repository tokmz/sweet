# Sweet 链路追踪包

## 简介

链路追踪包是Sweet社交电商分销系统的核心组件之一，基于OpenTelemetry标准实现，提供了全面的分布式链路追踪功能。通过该包，可以追踪系统中各个服务、组件之间的调用关系，帮助开发者快速定位性能瓶颈和故障点。

## 主要特性

- **多种采样策略**：支持固定比例采样、父级采样、全采样、不采样和动态采样
- **安全增强**：支持TLS加密传输，敏感信息自动过滤和掩码处理
- **性能优化**：批处理机制、压缩传输、内存优化等提升性能
- **多种导出器**：支持HTTP、gRPC和标准输出等多种导出方式
- **中间件集成**：内置Gin、GORM、gRPC、Redis、RabbitMQ等中间件的链路追踪支持
- **丰富的工具函数**：提供了大量辅助函数，简化链路追踪的使用

## 配置选项

### 基础配置

```go
type Config struct {
	// 是否启用链路追踪
	Enabled bool `json:"enabled"`
	// 链路追踪服务地址，如：http://jaeger:14268/api/traces
	Endpoint string `json:"endpoint"`
	// 链路追踪采样率，范围 0.0-1.0
	SamplingRate float64 `json:"sampling_rate"`
	// 服务名称
	ServiceName string `json:"service_name"`
	// 环境名称
	Environment string `json:"environment"`
	// 服务实例ID
	InstanceID string `json:"instance_id"`
	// 额外的服务标签
	ServiceTags map[string]string `json:"service_tags"`
	// 是否记录SQL语句
	RecordSQL bool `json:"record_sql"`
	// 是否记录HTTP请求和响应头
	RecordHTTPHeaders bool `json:"record_http_headers"`
	// 是否记录HTTP请求和响应体
	RecordHTTPBody bool `json:"record_http_body"`
	// 是否记录错误堆栈
	RecordErrorStack bool `json:"record_error_stack"`
	// 最大属性长度限制
	MaxAttributeValueLength int `json:"max_attribute_value_length"`
	// 安全配置
	Security SecurityConfig `json:"security"`
	// 性能配置
	Performance PerformanceConfig `json:"performance"`
	// 采样策略
	Sampling SamplingConfig `json:"sampling"`
}
```

### 安全配置

```go
type SecurityConfig struct {
	// 是否启用TLS
	EnableTLS bool `json:"enable_tls"`
	// 证书文件路径
	CertFile string `json:"cert_file"`
	// 密钥文件路径
	KeyFile string `json:"key_file"`
	// CA证书文件路径
	CACertFile string `json:"ca_cert_file"`
	// 是否跳过证书验证（仅用于开发环境）
	InsecureSkipVerify bool `json:"insecure_skip_verify"`
	// 敏感信息过滤规则
	SensitiveKeys []string `json:"sensitive_keys"`
	// 是否启用敏感信息过滤
	EnableSensitiveFilter bool `json:"enable_sensitive_filter"`
}
```

### 性能配置

```go
type PerformanceConfig struct {
	// 批处理大小
	BatchSize int `json:"batch_size"`
	// 批处理发送间隔（毫秒）
	BatchIntervalMs int `json:"batch_interval_ms"`
	// 最大导出并发数
	MaxExportBatchSize int `json:"max_export_batch_size"`
	// 是否启用压缩
	EnableCompression bool `json:"enable_compression"`
	// 是否启用内存优化
	EnableMemoryOptimization bool `json:"enable_memory_optimization"`
	// 缓存大小
	CacheSize int `json:"cache_size"`
}
```

### 采样配置

```go
type SamplingConfig struct {
	// 采样类型: ratio(比例采样), parent(父级采样), always_on(全采样), always_off(不采样), dynamic(动态采样)
	Type string `json:"type"`
	// 动态采样参数
	Dynamic DynamicSamplingConfig `json:"dynamic"`
}

type DynamicSamplingConfig struct {
	// 最大QPS阈值，超过此值降低采样率
	MaxQPS int `json:"max_qps"`
	// 最小采样率
	MinRate float64 `json:"min_rate"`
	// 目标采样率
	TargetSamplesPerSecond int `json:"target_samples_per_second"`
}
```

## 使用示例

### 初始化链路追踪

```go
package main

import (
	"context"
	"log"

	"sweet/pkg/tracing"
)

func main() {
	// 创建配置
	cfg := tracing.DefaultConfig()
	cfg.ServiceName = "my-service"
	cfg.Endpoint = "http://jaeger:14268/api/traces"
	cfg.SamplingRate = 0.1
	cfg.Security.EnableTLS = true
	cfg.Security.CACertFile = "/path/to/ca.crt"
	cfg.Sampling.Type = "dynamic"

	// 初始化链路追踪
	if err := tracing.Init(cfg); err != nil {
		log.Fatalf("Failed to initialize tracing: %v", err)
	}
	defer tracing.Close(context.Background())

	// 应用程序逻辑...
}
```

### 创建和使用Span

```go
func HandleRequest(ctx context.Context, req Request) (Response, error) {
	// 创建新的span
	ctx, span := tracing.StartSpan(ctx, "HandleRequest")
	defer tracing.End(ctx)

	// 添加属性
	tracing.SetAttributes(ctx,
		attribute.String("request.id", req.ID),
		attribute.String("request.type", req.Type),
	)

	// 处理请求
	result, err := processRequest(ctx, req)

	// 记录错误（如果有）
	if err != nil {
		tracing.RecordError(ctx, err)
	}

	return result, err
}
```

### 使用中间件

#### Gin中间件

```go
func SetupRouter() *gin.Engine {
	r := gin.New()

	// 添加链路追踪中间件
	r.Use(middleware.GinTracer())

	// 路由配置...

	return r
}
```

#### GORM中间件

```go
func SetupDB(dsn string) (*gorm.DB, error) {
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	// 添加链路追踪插件
	tracer := middleware.NewGormTracingPlugin("my-service")
	if err := db.Use(tracer); err != nil {
		return nil, err
	}

	return db, nil
}
```

#### Redis中间件

```go
func SetupRedis(addr string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: addr,
	})

	// 添加链路追踪
	client = middleware.WrapRedisClient(client, "my-service")

	return client
}
```

#### RabbitMQ中间件

```go
func SetupRabbitMQ(conn *amqp.Connection) (*middleware.WrapChannel, error) {
	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	// 添加链路追踪
	wrappedCh := middleware.NewWrappedChannel(ch, "my-service")

	return wrappedCh, nil
}
```

## 最佳实践

1. **合理设置采样率**：生产环境建议使用动态采样，根据流量自动调整采样率
2. **启用TLS加密**：生产环境必须启用TLS加密，保护链路追踪数据安全
3. **过滤敏感信息**：确保敏感信息不会被记录到链路追踪系统中
4. **设置合理的批处理参数**：根据系统负载调整批处理大小和间隔，平衡性能和资源消耗
5. **使用有意义的Span名称**：Span名称应该能够清晰表达操作的含义
6. **记录关键属性**：为Span添加有助于问题排查的关键属性，但避免过多无用信息
7. **正确结束Span**：确保所有创建的Span都被正确结束，避免资源泄漏

## 故障排查

1. **链路追踪数据不可见**：检查采样率设置，确保不是设置为0或过低
2. **TLS连接失败**：检查证书路径和权限，确保证书有效且未过期
3. **性能问题**：调整批处理参数，减小批处理间隔或增大批处理大小
4. **内存占用过高**：检查是否有未正确结束的Span，或减小批处理缓冲区大小
5. **敏感信息泄露**：检查敏感信息过滤配置，添加更多的敏感关键词

## 扩展开发

链路追踪包设计为可扩展的，可以通过以下方式进行扩展：

1. **添加新的中间件**：实现特定组件的链路追踪支持
2. **自定义采样策略**：实现自定义的采样策略
3. **扩展敏感信息过滤**：添加更多的敏感信息识别和掩码规则
4. **添加新的导出器**：支持更多的链路追踪后端系统

## 版本历史

- **v1.0.0**：基础链路追踪功能，支持HTTP和gRPC导出
- **v1.1.0**：添加Gin和GORM中间件支持
- **v1.2.0**：添加gRPC中间件支持
- **v2.0.0**：重构核心架构，添加安全性和性能优化
- **v2.1.0**：添加Redis和RabbitMQ中间件支持，增强敏感信息过滤