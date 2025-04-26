package tracing

// Config 链路追踪配置
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

// SecurityConfig 安全相关配置
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

// PerformanceConfig 性能相关配置
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

// SamplingConfig 采样策略配置
type SamplingConfig struct {
	// 采样类型: ratio(比例采样), parent(父级采样), always_on(全采样), always_off(不采样), dynamic(动态采样)
	Type string `json:"type"`
	// 动态采样参数
	Dynamic DynamicSamplingConfig `json:"dynamic"`
}

// DynamicSamplingConfig 动态采样配置
type DynamicSamplingConfig struct {
	// 最大QPS阈值，超过此值降低采样率
	MaxQPS int `json:"max_qps"`
	// 最小采样率
	MinRate float64 `json:"min_rate"`
	// 目标采样率
	TargetSamplesPerSecond int `json:"target_samples_per_second"`
}

// DefaultConfig 返回默认配置
func DefaultConfig() Config {
	return Config{
		Enabled:                 true,
		SamplingRate:            0.1,
		ServiceName:             "sweet-service",
		Environment:             "development",
		RecordSQL:               true,
		RecordHTTPHeaders:       true,
		RecordHTTPBody:          false,
		RecordErrorStack:        true,
		MaxAttributeValueLength: 2048,
		Security: SecurityConfig{
			EnableTLS:             false,
			InsecureSkipVerify:    false,
			EnableSensitiveFilter: true,
			SensitiveKeys: []string{
				"password", "passwd", "pwd", "secret", "token", "apikey",
				"api_key", "auth", "credential", "credit", "card", "authorization",
				"access_token", "refresh_token", "private_key", "session",
			},
		},
		Performance: PerformanceConfig{
			BatchSize:                512,
			BatchIntervalMs:          1000,
			MaxExportBatchSize:       512,
			EnableCompression:        true,
			EnableMemoryOptimization: true,
			CacheSize:                1000,
		},
		Sampling: SamplingConfig{
			Type: "ratio",
			Dynamic: DynamicSamplingConfig{
				MaxQPS:                 1000,
				MinRate:                0.01,
				TargetSamplesPerSecond: 100,
			},
		},
	}
}
