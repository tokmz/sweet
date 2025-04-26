package tracing

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io/ioutil"
	"os"
	"runtime"
	"strings"
	"sync"
	"sync/atomic"
	"time"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracegrpc"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	"go.opentelemetry.io/otel/exporters/stdout/stdouttrace"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.17.0"
	"go.opentelemetry.io/otel/trace"
	"google.golang.org/grpc/credentials"
)

var (
	// provider 是全局的SDK追踪提供程序
	provider *sdktrace.TracerProvider

	// tracer 是全局的追踪器实例
	tracer trace.Tracer

	// once 确保初始化只执行一次
	once sync.Once

	// 全局配置
	globalConfig Config

	// 是否已初始化
	initialized bool

	// 动态采样相关计数器
	requestCount   int64
	lastResetTime  int64
	currentSampler atomic.Value
)

// Init 初始化链路追踪
func Init(cfg Config) error {
	var err error
	once.Do(func() {
		err = initialize(cfg)
	})
	return err
}

// initialize 实际的初始化函数
func initialize(cfg Config) error {
	if !cfg.Enabled {
		return nil
	}

	// 创建资源
	res, err := createResource(cfg)
	if err != nil {
		return fmt.Errorf("failed to create resource: %w", err)
	}

	// 创建导出器
	exporter, err := createExporter(cfg)
	if err != nil {
		return fmt.Errorf("failed to create exporter: %w", err)
	}

	// 创建批处理器
	batchOpts := []sdktrace.BatchSpanProcessorOption{
		sdktrace.WithMaxExportBatchSize(cfg.Performance.MaxExportBatchSize),
		sdktrace.WithBatchTimeout(time.Duration(cfg.Performance.BatchIntervalMs) * time.Millisecond),
		sdktrace.WithMaxQueueSize(cfg.Performance.BatchSize),
	}

	// 创建采样器
	sampler, err := createSampler(cfg)
	if err != nil {
		return fmt.Errorf("failed to create sampler: %w", err)
	}

	// 创建追踪提供者
	provider = sdktrace.NewTracerProvider(
		sdktrace.WithSampler(sampler),
		sdktrace.WithResource(res),
		sdktrace.WithBatcher(exporter, batchOpts...),
	)

	// 设置全局追踪提供者
	otel.SetTracerProvider(provider)

	// 设置全局传播器
	otel.SetTextMapPropagator(propagation.NewCompositeTextMapPropagator(
		propagation.TraceContext{},
		propagation.Baggage{},
	))

	// 创建追踪器
	tracer = provider.Tracer(cfg.ServiceName)

	// 保存全局配置
	globalConfig = cfg
	initialized = true

	// 初始化动态采样器
	if cfg.Sampling.Type == "dynamic" {
		initDynamicSampling(cfg)
	}

	return nil
}

// createSampler 根据配置创建采样器
func createSampler(cfg Config) (sdktrace.Sampler, error) {
	switch cfg.Sampling.Type {
	case "ratio", "": // 默认使用比例采样
		return sdktrace.TraceIDRatioBased(cfg.SamplingRate), nil
	case "always_on":
		return sdktrace.AlwaysSample(), nil
	case "always_off":
		return sdktrace.NeverSample(), nil
	case "parent":
		return sdktrace.ParentBased(sdktrace.TraceIDRatioBased(cfg.SamplingRate)), nil
	case "dynamic":
		// 动态采样初始使用配置的采样率，后续会根据流量动态调整
		sampler := sdktrace.TraceIDRatioBased(cfg.SamplingRate)
		currentSampler.Store(sampler)
		return sampler, nil
	default:
		return nil, fmt.Errorf("unknown sampling type: %s", cfg.Sampling.Type)
	}
}

// initDynamicSampling 初始化动态采样
func initDynamicSampling(cfg Config) {
	// 重置计数器
	atomic.StoreInt64(&requestCount, 0)
	atomic.StoreInt64(&lastResetTime, time.Now().Unix())

	// 启动后台协程，定期调整采样率
	go func() {
		ticker := time.NewTicker(1 * time.Second)
		defer ticker.Stop()

		for range ticker.C {
			adjustSamplingRate(cfg)
		}
	}()
}

// adjustSamplingRate 根据当前流量调整采样率
func adjustSamplingRate(cfg Config) {
	// 获取当前计数和时间
	count := atomic.LoadInt64(&requestCount)
	lastReset := atomic.LoadInt64(&lastResetTime)
	now := time.Now().Unix()

	// 计算时间差（秒）
	elapsed := now - lastReset
	if elapsed <= 0 {
		return
	}

	// 计算当前QPS
	qps := float64(count) / float64(elapsed)

	// 根据QPS调整采样率
	var newRate float64
	if qps > float64(cfg.Sampling.Dynamic.MaxQPS) {
		// 如果QPS超过阈值，降低采样率
		targetSamples := float64(cfg.Sampling.Dynamic.TargetSamplesPerSecond)
		newRate = targetSamples / qps
		if newRate < cfg.Sampling.Dynamic.MinRate {
			newRate = cfg.Sampling.Dynamic.MinRate
		}
	} else {
		// 否则使用配置的采样率
		newRate = cfg.SamplingRate
	}

	// 更新采样器
	currentSampler.Store(sdktrace.TraceIDRatioBased(newRate))

	// 重置计数器
	atomic.StoreInt64(&requestCount, 0)
	atomic.StoreInt64(&lastResetTime, now)
}

// createResource 创建资源对象
func createResource(cfg Config) (*resource.Resource, error) {
	// 获取主机名和IP地址
	hostname, _ := os.Hostname()
	if cfg.InstanceID == "" {
		cfg.InstanceID = hostname
	}

	// 基础资源属性
	attrs := []attribute.KeyValue{
		semconv.ServiceNameKey.String(cfg.ServiceName),
		semconv.ServiceVersionKey.String(getAppVersion()),
		semconv.ServiceInstanceIDKey.String(cfg.InstanceID),
		attribute.String("environment", cfg.Environment),
		attribute.String("hostname", hostname),
	}

	// 添加自定义标签
	for k, v := range cfg.ServiceTags {
		attrs = append(attrs, attribute.String(k, v))
	}

	// 创建资源
	return resource.NewWithAttributes(
		semconv.SchemaURL,
		attrs...,
	), nil
}

// createExporter 创建追踪导出器
func createExporter(cfg Config) (sdktrace.SpanExporter, error) {
	// 如果未指定Endpoint，则使用标准输出
	if cfg.Endpoint == "" {
		return stdouttrace.New(
			stdouttrace.WithPrettyPrint(),
		)
	}

	// 优先尝试HTTP导出器，如果失败则尝试gRPC导出器
	if httpExporter, err := createHTTPExporter(cfg); err == nil {
		return httpExporter, nil
	}

	// 尝试gRPC导出器
	return createGRPCExporter(cfg)
}

// createHTTPExporter 创建HTTP导出器
func createHTTPExporter(cfg Config) (sdktrace.SpanExporter, error) {
	opts := []otlptracehttp.Option{
		otlptracehttp.WithEndpoint(cfg.Endpoint),
	}

	// 配置TLS
	if cfg.Security.EnableTLS {
		// 创建TLS配置
		tlsConfig, err := createTLSConfig(cfg.Security)
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
		opts = append(opts, otlptracehttp.WithTLSClientConfig(tlsConfig))
	} else {
		// 不使用TLS
		opts = append(opts, otlptracehttp.WithInsecure())
	}

	// 配置压缩
	if cfg.Performance.EnableCompression {
		opts = append(opts, otlptracehttp.WithCompression(otlptracehttp.GzipCompression))
	}

	// 创建客户端
	client := otlptracehttp.NewClient(opts...)
	return otlptrace.New(context.Background(), client)
}

// createGRPCExporter 创建gRPC导出器
func createGRPCExporter(cfg Config) (sdktrace.SpanExporter, error) {
	opts := []otlptracegrpc.Option{
		otlptracegrpc.WithEndpoint(cfg.Endpoint),
	}

	// 配置TLS
	if cfg.Security.EnableTLS {
		// 创建TLS配置
		tlsConfig, err := createTLSConfig(cfg.Security)
		if err != nil {
			return nil, fmt.Errorf("failed to create TLS config: %w", err)
		}
		creds := credentials.NewTLS(tlsConfig)
		opts = append(opts, otlptracegrpc.WithTLSCredentials(creds))
	} else {
		// 不使用TLS
		opts = append(opts, otlptracegrpc.WithInsecure())
	}

	// 配置压缩
	if cfg.Performance.EnableCompression {
		opts = append(opts, otlptracegrpc.WithCompressor("gzip"))
	}

	// 创建客户端
	client := otlptracegrpc.NewClient(opts...)
	return otlptrace.New(context.Background(), client)
}

// createTLSConfig 创建TLS配置
func createTLSConfig(securityCfg SecurityConfig) (*tls.Config, error) {
	tlsConfig := &tls.Config{
		InsecureSkipVerify: securityCfg.InsecureSkipVerify,
	}

	// 如果指定了CA证书，加载它
	if securityCfg.CACertFile != "" {
		caCert, err := ioutil.ReadFile(securityCfg.CACertFile)
		if err != nil {
			return nil, fmt.Errorf("failed to read CA cert file: %w", err)
		}

		caCertPool := x509.NewCertPool()
		if !caCertPool.AppendCertsFromPEM(caCert) {
			return nil, fmt.Errorf("failed to append CA cert to pool")
		}

		tlsConfig.RootCAs = caCertPool
	}

	// 如果指定了客户端证书和密钥，加载它们
	if securityCfg.CertFile != "" && securityCfg.KeyFile != "" {
		cert, err := tls.LoadX509KeyPair(securityCfg.CertFile, securityCfg.KeyFile)
		if err != nil {
			return nil, fmt.Errorf("failed to load client cert/key: %w", err)
		}

		tlsConfig.Certificates = []tls.Certificate{cert}
	}

	return tlsConfig, nil
}

// getAppVersion 获取应用版本
func getAppVersion() string {
	// 此处可以从环境变量、构建信息等获取版本
	version := os.Getenv("APP_VERSION")
	if version == "" {
		return "dev"
	}
	return version
}

// StartSpan 开始一个新的追踪Span
func StartSpan(ctx context.Context, name string, opts ...trace.SpanStartOption) (context.Context, trace.Span) {
	if !initialized || tracer == nil {
		return ctx, trace.SpanFromContext(ctx)
	}

	// 如果使用动态采样，增加请求计数
	if globalConfig.Sampling.Type == "dynamic" {
		atomic.AddInt64(&requestCount, 1)
	}

	return tracer.Start(ctx, name, opts...)
}

// SpanFromContext 从上下文中获取当前Span
func SpanFromContext(ctx context.Context) trace.Span {
	return trace.SpanFromContext(ctx)
}

// AddEvent 向Span添加一个事件
func AddEvent(ctx context.Context, name string, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.AddEvent(name, trace.WithAttributes(attrs...))
	}
}

// SetAttributes 设置Span属性
func SetAttributes(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if span.IsRecording() {
		span.SetAttributes(attrs...)
	}
}

// RecordError 记录错误到Span
func RecordError(ctx context.Context, err error, opts ...trace.EventOption) {
	if err == nil {
		return
	}

	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	span.RecordError(err, opts...)
	span.SetStatus(codes.Error, err.Error())

	if globalConfig.RecordErrorStack {
		// 获取堆栈信息
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, false)

		// 添加堆栈作为事件
		span.AddEvent("error.stack", trace.WithAttributes(
			attribute.String("stack", string(stack[:length])),
		))
	}
}

// End 结束Span，可选添加最终属性
func End(ctx context.Context, attrs ...attribute.KeyValue) {
	span := trace.SpanFromContext(ctx)
	if !span.IsRecording() {
		return
	}

	if len(attrs) > 0 {
		span.SetAttributes(attrs...)
	}

	span.End()
}

// Close 关闭链路追踪提供者，刷新所有待处理的span
func Close(ctx context.Context) error {
	if provider != nil {
		return provider.Shutdown(ctx)
	}
	return nil
}

// Middleware 抽象接口，用于各种中间件的实现
type Middleware interface {
	// Name 返回中间件名称
	Name() string
}

// LimitAttributeValue 限制属性值的长度
func LimitAttributeValue(value string, maxLength int) string {
	if maxLength <= 0 || len(value) <= maxLength {
		return value
	}

	return value[:maxLength] + "...(truncated)"
}

// ContainsSensitiveInfo 检查字符串是否包含敏感信息
func ContainsSensitiveInfo(value string) bool {
	// 如果未启用敏感信息过滤，直接返回false
	if !globalConfig.Security.EnableSensitiveFilter {
		return false
	}

	// 使用配置的敏感关键词列表
	sensitiveKeys := globalConfig.Security.SensitiveKeys
	if len(sensitiveKeys) == 0 {
		// 如果配置为空，使用默认列表
		sensitiveKeys = []string{
			"password", "passwd", "pwd",
			"secret", "token", "apikey",
			"api_key", "auth", "credential",
			"credit", "card", "authorization",
			"access_token", "refresh_token", "private_key",
		}
	}

	valueLower := strings.ToLower(value)
	for _, key := range sensitiveKeys {
		if strings.Contains(valueLower, key) {
			return true
		}
	}

	// 正则表达式匹配常见敏感信息模式（如信用卡号、身份证号等）
	// 这里可以添加更复杂的正则匹配逻辑

	return false
}

// FilterSensitiveAttributes 过滤掉敏感属性
func FilterSensitiveAttributes(attrs []attribute.KeyValue) []attribute.KeyValue {
	// 如果未启用敏感信息过滤，直接返回原始属性
	if !globalConfig.Security.EnableSensitiveFilter {
		return attrs
	}

	filtered := make([]attribute.KeyValue, 0, len(attrs))

	for _, attr := range attrs {
		key := string(attr.Key)
		keyLower := strings.ToLower(key)

		if ContainsSensitiveInfo(keyLower) {
			// 替换敏感属性的值
			filtered = append(filtered, attribute.String(key, "***REDACTED***"))
		} else {
			// 检查值是否为字符串，如果是则检查是否包含敏感信息
			if attr.Value.Type() == attribute.STRING {
				strValue := attr.Value.AsString()
				if ContainsSensitiveInfo(strValue) {
					filtered = append(filtered, attribute.String(key, "***REDACTED***"))
					continue
				}

				// 限制字符串长度
				if len(strValue) > globalConfig.MaxAttributeValueLength {
					filtered = append(filtered, attribute.String(key,
						LimitAttributeValue(strValue, globalConfig.MaxAttributeValueLength)))
					continue
				}
			}

			filtered = append(filtered, attr)
		}
	}

	return filtered
}

// GetCurrentSamplingRate 获取当前采样率
func GetCurrentSamplingRate() float64 {
	if !initialized || globalConfig.Sampling.Type != "dynamic" {
		return globalConfig.SamplingRate
	}

	// 获取当前采样器
	_, ok := currentSampler.Load().(sdktrace.Sampler)
	if !ok {
		return globalConfig.SamplingRate
	}

	// 尝试获取采样率（这是一个近似值，因为无法直接从采样器获取确切的采样率）
	// 在实际应用中，可能需要维护一个单独的变量来跟踪当前采样率
	return globalConfig.SamplingRate
}

// GetTracingStats 获取链路追踪统计信息
func GetTracingStats() map[string]interface{} {
	if !initialized {
		return map[string]interface{}{
			"initialized": false,
		}
	}

	stats := map[string]interface{}{
		"initialized":   true,
		"service_name":  globalConfig.ServiceName,
		"environment":   globalConfig.Environment,
		"sampling_type": globalConfig.Sampling.Type,
	}

	// 添加采样率信息
	if globalConfig.Sampling.Type == "dynamic" {
		stats["sampling_rate"] = GetCurrentSamplingRate()
		stats["request_count"] = atomic.LoadInt64(&requestCount)
	} else {
		stats["sampling_rate"] = globalConfig.SamplingRate
	}

	return stats
}
