package auth

import (
	"time"
)

// TokenStyle 定义Token的生成风格
type TokenStyle int

const (
	// TokenStyleUUID 使用UUID生成Token
	TokenStyleUUID TokenStyle = iota
	// TokenStyleSimple 使用简单字符串生成Token
	TokenStyleSimple
	// TokenStyleJWT 使用标准JWT格式
	TokenStyleJWT
	// TokenStyleJWTMixed 使用混合JWT格式(JWT+UUID)
	TokenStyleJWTMixed
	// TokenStyleJWTUUID 使用UUID作为JWT的jti
	TokenStyleJWTUUID
	// TokenStyleCustom 使用自定义生成器
	TokenStyleCustom
)

// LoginModel 定义登录模式
type LoginModel int

const (
	// LoginModelSingle 单端登录模式(一个用户只能在一个设备上登录)
	LoginModelSingle LoginModel = iota
	// LoginModelMulti 多端登录模式(一个用户可以在多个设备上同时登录)
	LoginModelMulti
	// LoginModelExclusive 同端互斥登录模式(同类设备互斥，不同类设备可以同时在线)
	LoginModelExclusive
)

// SSOMode 定义单点登录模式
type SSOMode int

const (
	// SSOModeNone 不启用单点登录
	SSOModeNone SSOMode = iota
	// SSOModeLocal 本地单点登录(同域)
	SSOModeLocal
	// SSOModeRedis 基于Redis的单点登录(跨域)
	SSOModeRedis
	// SSOModeJWT 基于JWT的单点登录
	SSOModeJWT
)

// DeviceType 定义设备类型
type DeviceType string

const (
	// DeviceTypeWeb Web端
	DeviceTypeWeb DeviceType = "web"
	// DeviceTypeApp 移动App端
	DeviceTypeApp DeviceType = "app"
	// DeviceTypeDesktop 桌面端
	DeviceTypeDesktop DeviceType = "desktop"
	// DeviceTypeAPI API调用
	DeviceTypeAPI DeviceType = "api"
	// DeviceTypeOther 其他设备
	DeviceTypeOther DeviceType = "other"
)

// Config 认证管理器配置
type Config struct {
	// JWT配置
	JWTSecret     string        // JWT密钥
	TokenExpire   time.Duration // Token过期时间
	RefreshExpire time.Duration // 刷新Token过期时间
	TokenStyle    TokenStyle    // Token风格

	// 登录模式
	LoginModel LoginModel // 登录模式

	// Redis配置
	EnableRedis    bool   // 是否启用Redis
	RedisAddr      string // Redis地址
	RedisPass      string // Redis密码
	RedisDB        int    // Redis数据库
	RedisKeyPrefix string // Redis键前缀

	// 高级配置
	AutoRenew         bool          // 是否自动续签
	RenewThreshold    float64       // 续签阈值(0.0-1.0)
	EnableSecondAuth  bool          // 是否启用二级认证
	SecondAuthExpire  time.Duration // 二级认证过期时间
	EnableConcurrency bool          // 是否启用并发控制
	MaxConcurrency    int           // 最大并发数
}

// CasbinConfig Casbin配置
type CasbinConfig struct {
	Model       string // 模型配置文件路径
	Adapter     string // 适配器类型(file, mysql, postgres等)
	DSN         string // 数据源名称(适配器为数据库时使用)
	TableName   string // 表名(适配器为数据库时使用)
	AutoMigrate bool   // 是否自动迁移表结构
}

// SSOConfig 单点登录配置
type SSOConfig struct {
	Mode      SSOMode // 单点登录模式
	RedisAddr string  // Redis地址
	RedisPass string  // Redis密码
	RedisDB   int     // Redis数据库
	Domain    string  // Cookie域名(同域模式使用)
	Path      string  // Cookie路径(同域模式使用)
}

// LoginParams 登录参数
type LoginParams struct {
	UserID     int64                  // 用户ID
	Username   string                 // 用户名
	Device     DeviceType             // 设备类型
	IP         string                 // IP地址
	UserAgent  string                 // 用户代理
	RememberMe bool                   // 是否记住我
	ExtraData  map[string]interface{} // 额外数据
}

// TokenInfo Token信息
type TokenInfo struct {
	Token        string                 // Token值
	RefreshToken string                 // 刷新Token
	ExpiresAt    int64                  // 过期时间
	UserID       int64                  // 用户ID
	Username     string                 // 用户名
	Device       DeviceType             // 设备类型
	IP           string                 // IP地址
	UserAgent    string                 // 用户代理
	LoginTime    int64                  // 登录时间
	ExtraData    map[string]interface{} // 额外数据
}

// UserSession 用户会话信息
type UserSession struct {
	UserID    int64      // 用户ID
	Username  string     // 用户名
	Device    DeviceType // 设备类型
	IP        string     // IP地址
	UserAgent string     // 用户代理
	LoginTime int64      // 登录时间
	ExpiresAt int64      // 过期时间
	Token     string     // Token值
}

// CustomTokenGenerator 自定义Token生成器
type CustomTokenGenerator func(userID int64, device DeviceType, extraData map[string]interface{}) (string, error)

// AuthEventType 认证事件类型
type AuthEventType int

const (
	// EventLogin 登录事件
	EventLogin AuthEventType = iota
	// EventLogout 注销事件
	EventLogout
	// EventKickout 踢出事件
	EventKickout
	// EventTokenRefresh Token刷新事件
	EventTokenRefresh
	// EventTokenExpired Token过期事件
	EventTokenExpired
	// EventSecondAuthEnabled 二级认证启用事件
	EventSecondAuthEnabled
	// EventSecondAuthDisabled 二级认证禁用事件
	EventSecondAuthDisabled
)

// AuthEvent 认证事件
type AuthEvent struct {
	Type      AuthEventType // 事件类型
	UserID    int64         // 用户ID
	Username  string        // 用户名
	Device    DeviceType    // 设备类型
	IP        string        // IP地址
	Token     string        // Token值
	Timestamp int64         // 事件时间戳
	ExtraData interface{}   // 额外数据
}

// AuthEventListener 认证事件监听器
type AuthEventListener func(event AuthEvent)

// Claims JWT声明
type Claims struct {
	UserID    int64                  `json:"uid"`
	Username  string                 `json:"username"`
	Device    string                 `json:"device"`
	ExtraData map[string]interface{} `json:"extra,omitempty"`
}

// Storage 存储接口
type Storage interface {
	// SaveToken 保存Token信息
	SaveToken(token string, info TokenInfo) error
	// GetToken 获取Token信息
	GetToken(token string) (TokenInfo, error)
	// RemoveToken 移除Token
	RemoveToken(token string) error
	// GetUserTokens 获取用户的所有Token
	GetUserTokens(userID int64) ([]string, error)
	// RemoveUserTokens 移除用户的所有Token
	RemoveUserTokens(userID int64) error
	// GetDeviceToken 获取用户在特定设备上的Token
	GetDeviceToken(userID int64, device DeviceType) (string, error)
	// SaveSecondAuth 保存二级认证信息
	SaveSecondAuth(token string, expireAt int64) error
	// CheckSecondAuth 检查二级认证
	CheckSecondAuth(token string) (bool, error)
	// RemoveSecondAuth 移除二级认证
	RemoveSecondAuth(token string) error
	// GetOnlineUsers 获取在线用户
	GetOnlineUsers() ([]UserSession, error)
	// Close 关闭存储
	Close() error
}
