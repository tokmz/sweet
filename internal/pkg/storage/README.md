# Storage 存储包

一个功能强大、高度可扩展的Go语言文件存储包，支持多种存储后端和高级文件处理功能。

## 🚀 核心特性

### 1. 多存储后端支持（策略模式）
- **本地存储**：支持本地文件系统存储
- **七牛云存储**：集成七牛云对象存储服务
- **阿里云OSS**：支持阿里云对象存储服务
- **腾讯云COS**：集成腾讯云对象存储服务
- **可扩展架构**：基于策略模式，轻松添加新的存储后端

### 2. 大文件处理能力
- **并发分片上传**：自动将大文件分片并发上传，提升传输效率
- **秒传功能**：基于文件MD5校验，实现文件秒传
- **断点续传**：支持网络中断后的断点续传功能
- **智能分片**：根据文件大小和网络状况自动调整分片策略

### 3. 安全校验链（责任链模式）
- **文件类型校验**：支持白名单/黑名单模式的文件类型检查
- **文件大小限制**：可配置的文件大小上限检查
- **病毒扫描**：集成病毒扫描引擎，确保文件安全
- **内容检测**：支持敏感内容检测和过滤
- **可扩展校验**：基于责任链模式，轻松添加自定义校验规则

### 4. 事件钩子系统
- **上传前钩子**：文件上传前的预处理和校验
- **上传中钩子**：上传过程中的进度监控和状态更新
- **上传完成钩子**：上传完成后的后处理操作
- **清理钩子**：文件删除和清理时的回调处理
- **错误钩子**：异常情况的处理和通知

### 5. HTTP Handler 支持
- **RESTful API**：提供标准的文件上传、下载、删除接口
- **进度监控**：实时上传/下载进度反馈
- **批量操作**：支持批量文件上传和管理
- **预签名URL**：生成临时访问链接
- **中间件集成**：与Gin等Web框架无缝集成

## 📁 项目结构

```
storage/
├── README.md              # 项目文档
├── interface.go           # 核心接口定义
├── factory.go             # 存储工厂类
├── manager.go             # 存储管理器
├── qiniu.go              # 七牛云存储实现
├── tencent.go            # 腾讯云COS实现
├── alioss.go             # 阿里云OSS实现（待实现）
├── local.go              # 本地存储实现（待实现）
├── validator/            # 校验器目录
│   ├── chain.go          # 责任链模式实现
│   ├── file_type.go      # 文件类型校验
│   ├── file_size.go      # 文件大小校验
│   └── virus_scan.go     # 病毒扫描校验
├── uploader/             # 上传器目录
│   ├── chunked.go        # 分片上传实现
│   ├── resumable.go      # 断点续传实现
│   └── concurrent.go     # 并发上传实现
├── hooks/                # 事件钩子目录
│   ├── interface.go      # 钩子接口定义
│   ├── manager.go        # 钩子管理器
│   └── builtin.go        # 内置钩子实现
└── handler/              # HTTP处理器目录
    ├── upload.go         # 上传处理器
    ├── download.go       # 下载处理器
    ├── manage.go         # 管理处理器
    └── middleware.go     # 中间件
```

## 🔧 快速开始

### 基本使用

```go
package main

import (
    "context"
    "log"
    "github.com/your-org/sweet/internal/pkg/storage"
)

func main() {
    // 创建存储管理器
    manager := storage.NewStorageManager()
    
    // 注册七牛云存储
    qiniuConfig := &storage.QiniuConfig{
        AccessKey: "your-access-key",
        SecretKey: "your-secret-key",
        Bucket:    "your-bucket",
        Domain:    "your-domain",
    }
    
    err := manager.RegisterStorage(storage.Qiniu, qiniuConfig)
    if err != nil {
        log.Fatal(err)
    }
    
    // 设置默认存储
    manager.SetDefaultStorage(storage.Qiniu)
    
    // 上传文件
    file, _ := os.Open("example.jpg")
    defer file.Close()
    
    result, err := manager.UploadToDefault(
        context.Background(),
        "images/example.jpg",
        file,
        1024*1024, // 1MB
        &storage.UploadOptions{
            ContentType: "image/jpeg",
        },
    )
    
    if err != nil {
        log.Fatal(err)
    }
    
    log.Printf("文件上传成功: %s", result.URL)
}
```

### 配置校验链

```go
// 创建校验链
validatorChain := validator.NewChain()

// 添加文件类型校验
validatorChain.Add(&validator.FileTypeValidator{
    AllowedTypes: []string{".jpg", ".png", ".pdf"},
})

// 添加文件大小校验
validatorChain.Add(&validator.FileSizeValidator{
    MaxSize: 10 * 1024 * 1024, // 10MB
})

// 添加病毒扫描
validatorChain.Add(&validator.VirusScanValidator{
    ScanEngine: "clamav",
})

// 在上传时使用校验链
uploadOptions := &storage.UploadOptions{
    Validator: validatorChain,
}
```

### 配置事件钩子

```go
// 创建钩子管理器
hookManager := hooks.NewManager()

// 注册上传前钩子
hookManager.RegisterBeforeUpload(func(ctx context.Context, key string, size int64) error {
    log.Printf("准备上传文件: %s, 大小: %d", key, size)
    return nil
})

// 注册上传完成钩子
hookManager.RegisterAfterUpload(func(ctx context.Context, result *storage.UploadResult) error {
    log.Printf("文件上传完成: %s", result.URL)
    // 可以在这里发送通知、更新数据库等
    return nil
})

// 在存储管理器中使用钩子
manager.SetHookManager(hookManager)
```

### HTTP Handler 集成

```go
package main

import (
    "github.com/gin-gonic/gin"
    "github.com/your-org/sweet/internal/pkg/storage"
    "github.com/your-org/sweet/internal/pkg/storage/handler"
)

func main() {
    // 创建存储管理器
    manager := storage.NewStorageManager()
    // ... 配置存储后端
    
    // 创建HTTP处理器
    uploadHandler := handler.NewUploadHandler(manager)
    downloadHandler := handler.NewDownloadHandler(manager)
    
    // 创建Gin路由
    r := gin.Default()
    
    // 注册路由
    api := r.Group("/api/v1")
    {
        api.POST("/upload", uploadHandler.Upload)
        api.POST("/upload/chunked", uploadHandler.ChunkedUpload)
        api.GET("/download/:key", downloadHandler.Download)
        api.DELETE("/files/:key", uploadHandler.Delete)
        api.GET("/files/:key/info", uploadHandler.GetFileInfo)
    }
    
    r.Run(":8080")
}
```

## 🔧 配置说明

### 存储后端配置

#### 七牛云配置
```go
type QiniuConfig struct {
    AccessKey string `json:"access_key" yaml:"access_key"`
    SecretKey string `json:"secret_key" yaml:"secret_key"`
    Bucket    string `json:"bucket" yaml:"bucket"`
    Domain    string `json:"domain" yaml:"domain"`
    UseHTTPS  bool   `json:"use_https" yaml:"use_https"`
    Zone      string `json:"zone" yaml:"zone"` // 存储区域
}
```

#### 腾讯云COS配置
```go
type TencentConfig struct {
    SecretID  string `json:"secret_id" yaml:"secret_id"`
    SecretKey string `json:"secret_key" yaml:"secret_key"`
    Region    string `json:"region" yaml:"region"`
    Bucket    string `json:"bucket" yaml:"bucket"`
    Domain    string `json:"domain" yaml:"domain"`
}
```

### 分片上传配置
```go
type ChunkedUploadConfig struct {
    ChunkSize    int64 `json:"chunk_size"`    // 分片大小，默认5MB
    MaxChunks    int   `json:"max_chunks"`    // 最大分片数
    Concurrency  int   `json:"concurrency"`   // 并发数，默认3
    RetryTimes   int   `json:"retry_times"`   // 重试次数
    EnableResume bool  `json:"enable_resume"` // 是否启用断点续传
}
```

### 存储配置文件

#### YAML配置示例
```yaml
# storage.yaml
storage:
  # 默认存储类型
  default: "qiniu"
  
  # 文件上传限制
  upload:
    # 全局文件大小限制 (字节)
    max_file_size: 104857600  # 100MB
    
    # 允许的文件类型 (白名单模式)
    allowed_types:
      - ".jpg"
      - ".jpeg"
      - ".png"
      - ".gif"
      - ".pdf"
      - ".doc"
      - ".docx"
      - ".xls"
      - ".xlsx"
      - ".ppt"
      - ".pptx"
      - ".txt"
      - ".zip"
      - ".rar"
    
    # 禁止的文件类型 (黑名单模式)
    forbidden_types:
      - ".exe"
      - ".bat"
      - ".cmd"
      - ".scr"
      - ".vbs"
      - ".js"
      - ".jar"
    
    # 按文件类型的大小限制
    type_size_limits:
      image:  # 图片类型
        types: [".jpg", ".jpeg", ".png", ".gif"]
        max_size: 10485760  # 10MB
      document:  # 文档类型
        types: [".pdf", ".doc", ".docx", ".xls", ".xlsx", ".ppt", ".pptx"]
        max_size: 52428800  # 50MB
      archive:  # 压缩包类型
        types: [".zip", ".rar", ".7z"]
        max_size: 104857600  # 100MB
    
    # 分片上传配置
    chunked:
      enabled: true
      chunk_size: 5242880  # 5MB
      max_chunks: 1000
      concurrency: 3
      retry_times: 3
      enable_resume: true
      # 大于此大小的文件自动启用分片上传
      auto_chunk_threshold: 20971520  # 20MB
  
  # 安全校验配置
  security:
    # 病毒扫描
    virus_scan:
      enabled: true
      engine: "clamav"
      timeout: 30  # 扫描超时时间(秒)
      async: true   # 异步扫描
    
    # 内容检测
    content_detection:
      enabled: true
      # 敏感词检测
      sensitive_words: true
      # 图片内容检测
      image_content: true
  
  # 存储后端配置
  backends:
    # 七牛云配置
    qiniu:
      access_key: "${QINIU_ACCESS_KEY}"
      secret_key: "${QINIU_SECRET_KEY}"
      bucket: "${QINIU_BUCKET}"
      domain: "${QINIU_DOMAIN}"
      use_https: true
      zone: "z0"  # 华东
    
    # 腾讯云COS配置
    tencent:
      secret_id: "${TENCENT_SECRET_ID}"
      secret_key: "${TENCENT_SECRET_KEY}"
      region: "ap-beijing"
      bucket: "${TENCENT_BUCKET}"
      domain: "${TENCENT_DOMAIN}"
    
    # 阿里云OSS配置
    alioss:
      access_key_id: "${ALIOSS_ACCESS_KEY_ID}"
      access_key_secret: "${ALIOSS_ACCESS_KEY_SECRET}"
      endpoint: "oss-cn-hangzhou.aliyuncs.com"
      bucket: "${ALIOSS_BUCKET}"
      domain: "${ALIOSS_DOMAIN}"
    
    # 本地存储配置
    local:
      root_path: "./uploads"
      url_prefix: "/static/uploads"
      create_date_dir: true  # 按日期创建目录
  
  # 缓存配置
  cache:
    # 文件元数据缓存
    metadata:
      enabled: true
      ttl: 3600  # 缓存时间(秒)
      max_size: 10000  # 最大缓存条目数
    
    # 预签名URL缓存
    presigned_url:
      enabled: true
      ttl: 1800  # 缓存时间(秒)
      max_size: 5000
  
  # 事件钩子配置
  hooks:
    # 上传前钩子
    before_upload:
      - "validate_user_quota"  # 验证用户配额
      - "check_duplicate"      # 检查重复文件
    
    # 上传完成钩子
    after_upload:
      - "update_user_quota"    # 更新用户配额
      - "send_notification"    # 发送通知
      - "generate_thumbnail"   # 生成缩略图
    
    # 删除钩子
    before_delete:
      - "check_permission"     # 检查删除权限
    
    after_delete:
      - "cleanup_cache"        # 清理缓存
      - "update_statistics"    # 更新统计信息
  
  # 监控配置
  monitoring:
    # 指标收集
    metrics:
      enabled: true
      interval: 60  # 收集间隔(秒)
    
    # 日志配置
    logging:
      level: "info"
      format: "json"
      # 敏感信息脱敏
      mask_sensitive: true
```

#### JSON配置示例
```json
{
  "storage": {
    "default": "qiniu",
    "upload": {
      "max_file_size": 104857600,
      "allowed_types": [".jpg", ".jpeg", ".png", ".pdf"],
      "forbidden_types": [".exe", ".bat", ".cmd"],
      "type_size_limits": {
        "image": {
          "types": [".jpg", ".jpeg", ".png", ".gif"],
          "max_size": 10485760
        }
      },
      "chunked": {
        "enabled": true,
        "chunk_size": 5242880,
        "concurrency": 3,
        "enable_resume": true
      }
    },
    "backends": {
      "qiniu": {
        "access_key": "your-access-key",
        "secret_key": "your-secret-key",
        "bucket": "your-bucket",
        "domain": "your-domain",
        "use_https": true
      }
    }
  }
}
```

#### Go配置结构体
```go
// StorageConfig 存储配置
type StorageConfig struct {
    Default string                    `json:"default" yaml:"default"`
    Upload  UploadConfig             `json:"upload" yaml:"upload"`
    Security SecurityConfig          `json:"security" yaml:"security"`
    Backends map[string]interface{}  `json:"backends" yaml:"backends"`
    Cache   CacheConfig             `json:"cache" yaml:"cache"`
    Hooks   HooksConfig             `json:"hooks" yaml:"hooks"`
    Monitoring MonitoringConfig      `json:"monitoring" yaml:"monitoring"`
}

// UploadConfig 上传配置
type UploadConfig struct {
    MaxFileSize      int64                          `json:"max_file_size" yaml:"max_file_size"`
    AllowedTypes     []string                       `json:"allowed_types" yaml:"allowed_types"`
    ForbiddenTypes   []string                       `json:"forbidden_types" yaml:"forbidden_types"`
    TypeSizeLimits   map[string]TypeSizeLimit       `json:"type_size_limits" yaml:"type_size_limits"`
    Chunked         ChunkedUploadConfig            `json:"chunked" yaml:"chunked"`
}

// TypeSizeLimit 按类型的大小限制
type TypeSizeLimit struct {
    Types   []string `json:"types" yaml:"types"`
    MaxSize int64    `json:"max_size" yaml:"max_size"`
}

// SecurityConfig 安全配置
type SecurityConfig struct {
    VirusScan        VirusScanConfig        `json:"virus_scan" yaml:"virus_scan"`
    ContentDetection ContentDetectionConfig `json:"content_detection" yaml:"content_detection"`
}

// VirusScanConfig 病毒扫描配置
type VirusScanConfig struct {
    Enabled bool   `json:"enabled" yaml:"enabled"`
    Engine  string `json:"engine" yaml:"engine"`
    Timeout int    `json:"timeout" yaml:"timeout"`
    Async   bool   `json:"async" yaml:"async"`
}

// ContentDetectionConfig 内容检测配置
type ContentDetectionConfig struct {
    Enabled        bool `json:"enabled" yaml:"enabled"`
    SensitiveWords bool `json:"sensitive_words" yaml:"sensitive_words"`
    ImageContent   bool `json:"image_content" yaml:"image_content"`
}

// CacheConfig 缓存配置
type CacheConfig struct {
    Metadata     CacheItemConfig `json:"metadata" yaml:"metadata"`
    PresignedURL CacheItemConfig `json:"presigned_url" yaml:"presigned_url"`
}

// CacheItemConfig 缓存项配置
type CacheItemConfig struct {
    Enabled bool `json:"enabled" yaml:"enabled"`
    TTL     int  `json:"ttl" yaml:"ttl"`
    MaxSize int  `json:"max_size" yaml:"max_size"`
}

// HooksConfig 钩子配置
type HooksConfig struct {
    BeforeUpload []string `json:"before_upload" yaml:"before_upload"`
    AfterUpload  []string `json:"after_upload" yaml:"after_upload"`
    BeforeDelete []string `json:"before_delete" yaml:"before_delete"`
    AfterDelete  []string `json:"after_delete" yaml:"after_delete"`
}

// MonitoringConfig 监控配置
type MonitoringConfig struct {
    Metrics MetricsConfig `json:"metrics" yaml:"metrics"`
    Logging LoggingConfig `json:"logging" yaml:"logging"`
}

// MetricsConfig 指标配置
type MetricsConfig struct {
    Enabled  bool `json:"enabled" yaml:"enabled"`
    Interval int  `json:"interval" yaml:"interval"`
}

// LoggingConfig 日志配置
type LoggingConfig struct {
    Level         string `json:"level" yaml:"level"`
    Format        string `json:"format" yaml:"format"`
    MaskSensitive bool   `json:"mask_sensitive" yaml:"mask_sensitive"`
}
```

### 配置文件加载

```go
package main

import (
    "fmt"
    "log"
    "gopkg.in/yaml.v3"
    "os"
    "github.com/your-org/sweet/internal/pkg/storage"
)

// LoadStorageConfig 加载存储配置
func LoadStorageConfig(configPath string) (*StorageConfig, error) {
    data, err := os.ReadFile(configPath)
    if err != nil {
        return nil, fmt.Errorf("读取配置文件失败: %v", err)
    }
    
    var config StorageConfig
    if err := yaml.Unmarshal(data, &config); err != nil {
        return nil, fmt.Errorf("解析配置文件失败: %v", err)
    }
    
    // 验证配置
    if err := validateConfig(&config); err != nil {
        return nil, fmt.Errorf("配置验证失败: %v", err)
    }
    
    return &config, nil
}

// validateConfig 验证配置
func validateConfig(config *StorageConfig) error {
    if config.Default == "" {
        return fmt.Errorf("必须指定默认存储类型")
    }
    
    if config.Upload.MaxFileSize <= 0 {
        return fmt.Errorf("最大文件大小必须大于0")
    }
    
    // 验证存储后端配置
    if _, exists := config.Backends[config.Default]; !exists {
        return fmt.Errorf("默认存储类型 %s 的配置不存在", config.Default)
    }
    
    return nil
}

// InitStorageWithConfig 使用配置初始化存储
func InitStorageWithConfig(configPath string) (*storage.StorageManager, error) {
    config, err := LoadStorageConfig(configPath)
    if err != nil {
        return nil, err
    }
    
    manager := storage.NewStorageManager()
    
    // 注册存储后端
    for backendType, backendConfig := range config.Backends {
        storageType, err := storage.ParseStorageType(backendType)
        if err != nil {
            log.Printf("跳过未知存储类型: %s", backendType)
            continue
        }
        
        // 根据类型创建具体配置
        var storageConfig storage.Config
        switch storageType {
        case storage.Qiniu:
            storageConfig = createQiniuConfig(backendConfig)
        case storage.TencentCOS:
            storageConfig = createTencentConfig(backendConfig)
        // ... 其他存储类型
        }
        
        if err := manager.RegisterStorage(storageType, storageConfig); err != nil {
            return nil, fmt.Errorf("注册存储 %s 失败: %v", backendType, err)
        }
    }
    
    // 设置默认存储
    defaultType, _ := storage.ParseStorageType(config.Default)
    if err := manager.SetDefaultStorage(defaultType); err != nil {
        return nil, fmt.Errorf("设置默认存储失败: %v", err)
    }
    
    return manager, nil
}

func main() {
    // 从配置文件初始化存储
    manager, err := InitStorageWithConfig("./config/storage.yaml")
    if err != nil {
        log.Fatal(err)
    }
    
    // 使用存储管理器
    // ...
}
```

## 🛡️ 安全特性

### 文件类型校验
- 支持基于文件扩展名的白名单/黑名单过滤
- 支持基于MIME类型的内容检测
- 支持基于文件头的真实类型检测

### 文件大小限制
- 支持全局文件大小限制
- 支持按文件类型的差异化大小限制
- 支持动态调整大小限制

### 病毒扫描
- 集成ClamAV病毒扫描引擎
- 支持自定义病毒扫描服务
- 异步扫描，不影响上传性能

## 📊 性能优化

### 并发上传
- 自动检测文件大小，智能选择上传策略
- 大文件自动启用分片并发上传
- 可配置的并发数和分片大小

### 缓存策略
- 文件元数据缓存，减少重复查询
- 预签名URL缓存，提升访问速度
- 智能缓存失效机制

### 连接池
- HTTP客户端连接池复用
- 数据库连接池优化
- 合理的超时和重试机制

## 🔍 监控和日志

### 指标监控
- 上传/下载成功率
- 平均响应时间
- 文件大小分布
- 存储后端健康状态

### 日志记录
- 结构化日志输出
- 可配置的日志级别
- 敏感信息脱敏
- 链路追踪支持

## 🧪 测试

```bash
# 运行单元测试
go test ./...

# 运行基准测试
go test -bench=. ./...

# 生成测试覆盖率报告
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out
```

## 📝 许可证

MIT License

## 🤝 贡献

欢迎提交Issue和Pull Request来帮助改进这个项目。

## 📞 支持

如有问题，请通过以下方式联系：
- 提交GitHub Issue
- 发送邮件至：support@example.com

---

**注意**：本文档描述的是存储包的完整功能规划，部分功能可能仍在开发中。请查看具体的实现文件了解当前可用的功能。