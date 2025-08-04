package storage

import (
	"context"
	"io"
	"net/http"
	"time"
)

// StorageType 存储类型
type StorageType string

const (
	StorageTypeLocal   StorageType = "local"   // 本地存储
	StorageTypeQiniu   StorageType = "qiniu"   // 七牛云存储
	StorageTypeTencent StorageType = "tencent" // 腾讯云存储
	StorageTypeAliOSS  StorageType = "alioss"  // 阿里云OSS存储
)

// UploadOptions 上传选项
type UploadOptions struct {
	FileName      string            // 文件名
	ContentType   string            // 内容类型
	ContentLength int64             // 内容长度
	Metadata      map[string]string // 元数据
	Expires       time.Duration     // 过期时间
	IsPublic      bool              // 是否公开访问
	PartSize      int64             // 分片大小，用于分片上传
	Threads       int               // 并发上传的线程数
}

// UploadResult 上传结果
type UploadResult struct {
	URL          string            // 访问URL
	ETag         string            // ETag
	Size         int64             // 文件大小
	ContentType  string            // 内容类型
	LastModified time.Time         // 最后修改时间
	Metadata     map[string]string // 元数据
}

// DownloadOptions 下载选项
type DownloadOptions struct {
	Range      string // 范围，如 "bytes=0-1023"
	IfModified string // 如果修改则下载
}

// FileInfo 文件信息
type FileInfo struct {
	Name         string            // 文件名
	Size         int64             // 文件大小
	LastModified time.Time         // 最后修改时间
	ETag         string            // ETag
	ContentType  string            // 内容类型
	IsDir        bool              // 是否是目录
	Metadata     map[string]string // 元数据
}

// Storage 存储接口
type Storage interface {
	// Upload 上传文件
	// reader: 文件内容读取器
	// key: 存储的键名/路径
	// options: 上传选项
	Upload(ctx context.Context, reader io.Reader, key string, options *UploadOptions) (*UploadResult, error)

	// Download 下载文件
	// key: 存储的键名/路径
	// writer: 写入下载内容的写入器
	// options: 下载选项
	Download(ctx context.Context, key string, writer io.Writer, options *DownloadOptions) error

	// Delete 删除文件
	// key: 存储的键名/路径
	Delete(ctx context.Context, key string) error

	// Exists 检查文件是否存在
	// key: 存储的键名/路径
	Exists(ctx context.Context, key string) (bool, error)

	// GetFileInfo 获取文件信息
	// key: 存储的键名/路径
	GetFileInfo(ctx context.Context, key string) (*FileInfo, error)

	// GeneratePresignedURL 生成预签名URL
	// key: 存储的键名/路径
	// expiration: 过期时间
	// isDownload: 是否是下载链接
	GeneratePresignedURL(ctx context.Context, key string, expiration time.Duration, isDownload bool) (string, error)

	// Copy 复制文件
	// srcKey: 源文件键名/路径
	// destKey: 目标文件键名/路径
	Copy(ctx context.Context, srcKey, destKey string) error

	// Move 移动文件
	// srcKey: 源文件键名/路径
	// destKey: 目标文件键名/路径
	Move(ctx context.Context, srcKey, destKey string) error

	// List 列出文件
	// prefix: 前缀
	// delimiter: 分隔符
	// marker: 标记
	// limit: 限制数量
	List(ctx context.Context, prefix, delimiter, marker string, limit int) ([]*FileInfo, string, error)

	// GetStorageType 获取存储类型
	GetStorageType() StorageType

	// Close 关闭客户端
	Close() error
}

// Config 配置接口
type Config interface {
	// Validate 验证配置
	Validate() error
	// GetStorageType 获取存储类型
	GetStorageType() StorageType
}

// MultipartUpload 分片上传接口
type MultipartUpload interface {
	// InitiateMultipartUpload 初始化分片上传
	InitiateMultipartUpload(ctx context.Context, key string, options *UploadOptions) (string, error)

	// UploadPart 上传分片
	// uploadID: 上传ID
	// partNumber: 分片号（从1开始）
	// reader: 分片内容读取器
	UploadPart(ctx context.Context, uploadID string, partNumber int, reader io.Reader) (string, error)

	// CompleteMultipartUpload 完成分片上传
	// uploadID: 上传ID
	// parts: 分片信息列表
	CompleteMultipartUpload(ctx context.Context, uploadID string, parts []PartInfo) (*UploadResult, error)

	// AbortMultipartUpload 取消分片上传
	// uploadID: 上传ID
	AbortMultipartUpload(ctx context.Context, uploadID string) error

	// ListParts 列出已上传的分片
	// uploadID: 上传ID
	ListParts(ctx context.Context, uploadID string) ([]PartInfo, error)
}

// PartInfo 分片信息
type PartInfo struct {
	PartNumber int       // 分片号
	ETag       string    // 分片ETag
	Size       int64     // 分片大小
	UploadTime time.Time // 上传时间
}

// EventType 事件类型
type EventType string

const (
	EventTypeBeforeUpload EventType = "before_upload" // 上传前
	EventTypeAfterUpload  EventType = "after_upload"  // 上传后
	EventTypeBeforeDelete EventType = "before_delete" // 删除前
	EventTypeAfterDelete  EventType = "after_delete"  // 删除后
	EventTypeUploadError  EventType = "upload_error"  // 上传错误
	EventTypeDeleteError  EventType = "delete_error"  // 删除错误
)

// Event 事件
type Event struct {
	Type      EventType              // 事件类型
	Key       string                 // 文件键名
	Timestamp time.Time              // 事件时间
	Data      map[string]interface{} // 事件数据
	Error     error                  // 错误信息（如果有）
}

// EventHandler 事件处理器
type EventHandler interface {
	// Handle 处理事件
	Handle(ctx context.Context, event *Event) error
}

// SecurityValidator 安全验证器接口
type SecurityValidator interface {
	// ValidateFileType 验证文件类型
	ValidateFileType(fileName string, contentType string) error

	// ValidateFileSize 验证文件大小
	ValidateFileSize(size int64) error

	// ScanVirus 病毒扫描
	ScanVirus(ctx context.Context, reader io.Reader) error

	// ValidateContent 内容检测
	ValidateContent(ctx context.Context, reader io.Reader, contentType string) error
}

// StorageManager 存储管理器接口
type StorageManager interface {
	// RegisterStorage 注册存储实例
	RegisterStorage(storageType StorageType, storage Storage) error

	// GetStorage 获取存储实例
	GetStorage(storageType StorageType) (Storage, error)

	// SetDefaultStorage 设置默认存储
	SetDefaultStorage(storageType StorageType) error

	// GetDefaultStorage 获取默认存储
	GetDefaultStorage() (Storage, error)

	// ListStorageTypes 列出所有注册的存储类型
	ListStorageTypes() []StorageType

	// RemoveStorage 移除存储实例
	RemoveStorage(storageType StorageType) error

	// Close 关闭所有存储实例
	Close() error
}

// Cache 缓存接口
type Cache interface {
	// Get 获取缓存
	Get(ctx context.Context, key string) ([]byte, error)

	// Set 设置缓存
	Set(ctx context.Context, key string, value []byte, expiration time.Duration) error

	// Delete 删除缓存
	Delete(ctx context.Context, key string) error

	// Exists 检查缓存是否存在
	Exists(ctx context.Context, key string) (bool, error)

	// Clear 清空缓存
	Clear(ctx context.Context) error
}

// HTTPHandler HTTP处理器接口
type HTTPHandler interface {
	// HandleUpload 处理文件上传
	HandleUpload(w http.ResponseWriter, r *http.Request)

	// HandleDownload 处理文件下载
	HandleDownload(w http.ResponseWriter, r *http.Request)

	// HandleDelete 处理文件删除
	HandleDelete(w http.ResponseWriter, r *http.Request)

	// HandleFileInfo 处理获取文件信息
	HandleFileInfo(w http.ResponseWriter, r *http.Request)

	// HandleList 处理文件列表
	HandleList(w http.ResponseWriter, r *http.Request)
}

// Statistics 统计信息
type Statistics struct {
	TotalFiles    int64     // 总文件数
	TotalSize     int64     // 总大小（字节）
	UploadCount   int64     // 上传次数
	DownloadCount int64     // 下载次数
	DeleteCount   int64     // 删除次数
	ErrorCount    int64     // 错误次数
	LastUpload    time.Time // 最后上传时间
	LastDownload  time.Time // 最后下载时间
	AverageSize   float64   // 平均文件大小
	StorageUsage  float64   // 存储使用率（0-1）
}

// Monitor 监控接口
type Monitor interface {
	// GetStatistics 获取统计信息
	GetStatistics(ctx context.Context) (*Statistics, error)

	// RecordUpload 记录上传
	RecordUpload(ctx context.Context, size int64, duration time.Duration) error

	// RecordDownload 记录下载
	RecordDownload(ctx context.Context, size int64, duration time.Duration) error

	// RecordError 记录错误
	RecordError(ctx context.Context, operation string, err error) error

	// GetHealthStatus 获取健康状态
	GetHealthStatus(ctx context.Context) (bool, error)
}

// Logger 日志接口
type Logger interface {
	// Debug 调试日志
	Debug(msg string, fields ...interface{})

	// Info 信息日志
	Info(msg string, fields ...interface{})

	// Warn 警告日志
	Warn(msg string, fields ...interface{})

	// Error 错误日志
	Error(msg string, fields ...interface{})

	// Fatal 致命错误日志
	Fatal(msg string, fields ...interface{})
}

// ProgressCallback 进度回调函数
type ProgressCallback func(uploaded, total int64)

// RetryPolicy 重试策略
type RetryPolicy struct {
	MaxRetries    int           // 最大重试次数
	InitialDelay  time.Duration // 初始延迟
	MaxDelay      time.Duration // 最大延迟
	BackoffFactor float64       // 退避因子
}

// CompressionType 压缩类型
type CompressionType string

const (
	CompressionTypeNone CompressionType = "none" // 不压缩
	CompressionTypeGzip CompressionType = "gzip" // Gzip压缩
	CompressionTypeLZ4  CompressionType = "lz4"  // LZ4压缩
)

// EncryptionType 加密类型
type EncryptionType string

const (
	EncryptionTypeNone   EncryptionType = "none"   // 不加密
	EncryptionTypeAES256 EncryptionType = "aes256" // AES256加密
)

// StorageOptions 存储选项
type StorageOptions struct {
	Compression CompressionType  // 压缩类型
	Encryption  EncryptionType   // 加密类型
	RetryPolicy *RetryPolicy     // 重试策略
	Timeout     time.Duration    // 超时时间
	Progress    ProgressCallback // 进度回调
}
