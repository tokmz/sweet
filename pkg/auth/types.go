package auth

import (
	"time"

	"github.com/casbin/casbin/v2"
)

// LoginMode 登录模式
type LoginMode int

const (
	// LoginModeSingle 单端登录：一个用户只能在一个设备上登录，新登录会踢掉旧登录
	LoginModeSingle LoginMode = iota + 1
	// LoginModeMulti 多端登录：一个用户可以在多个设备上同时登录
	LoginModeMulti
	// LoginModeMutex 同端互斥登录：同一类型设备只能登录一个账号
	LoginModeMutex
)

// Config 认证配置
type Config struct {
	// SecretKey 密钥
	SecretKey string
	// AccessTokenExpiry 访问令牌过期时间
	AccessTokenExpiry time.Duration
	// RefreshTokenExpiry 刷新令牌过期时间
	RefreshTokenExpiry time.Duration
	// LoginMode 登录模式
	LoginMode LoginMode
}

// TokenInfo 令牌信息
type TokenInfo struct {
	// UserID 用户ID
	UserID uint64 `json:"user_id"`
	// DeviceType 设备类型
	DeviceType string `json:"device_type"`
	// ExpiresAt 过期时间
	ExpiresAt int64 `json:"expires_at"`
	// IssuedAt 签发时间
	IssuedAt int64 `json:"issued_at"`
	// TokenType 令牌类型
	TokenType string `json:"token_type"`
}

// SessionInfo 会话信息
type SessionInfo struct {
	// UserID 用户ID
	UserID uint64 `json:"user_id"`
	// DeviceType 设备类型
	DeviceType string `json:"device_type"`
	// LoginTime 登录时间
	LoginTime int64 `json:"login_time"`
	// ExpireTime 过期时间
	ExpireTime int64 `json:"expire_time"`
	// RefreshToken 刷新令牌
	RefreshToken string `json:"refresh_token"`
	// IP 登录IP
	IP string `json:"ip"`
	// UserAgent 用户代理
	UserAgent string `json:"user_agent"`
}

// SessionStore 会话存储接口
type SessionStore interface {
	// SaveSession 保存会话
	SaveSession(session *SessionInfo) error
	// GetSession 获取会话
	GetSession(userID uint64, deviceType string) (*SessionInfo, error)
	// RemoveSession 删除会话
	RemoveSession(userID uint64, deviceType string) error
	// GetUserSessions 获取用户所有会话
	GetUserSessions(userID uint64) ([]*SessionInfo, error)
	// GetDeviceTypeSessions 获取指定设备类型的所有会话
	GetDeviceTypeSessions(deviceType string) ([]*SessionInfo, error)
	// CleanExpiredSessions 清理过期会话
	CleanExpiredSessions() error
}

// TokenManager 令牌管理器接口
type TokenManager interface {
	// GenerateToken 生成令牌
	GenerateToken(userID uint64, deviceType string, tokenType string) (string, error)
	// ParseToken 解析令牌
	ParseToken(token string) (*TokenInfo, error)
	// RefreshToken 刷新令牌
	RefreshToken(refreshToken string) (string, string, error)
	// RevokeToken 撤销令牌
	RevokeToken(userID uint64, deviceType string) error
}

// AuthService 认证服务接口
type AuthService interface {
	// Login 登录
	Login(userID uint64, deviceType string) (accessToken string, refreshToken string, err error)
	// Logout 登出
	Logout(userID uint64, deviceType string) error
	// RefreshToken 刷新令牌
	RefreshToken(refreshToken string) (accessToken string, newRefreshToken string, err error)
	// VerifyToken 验证令牌
	VerifyToken(token string) (*TokenInfo, error)
	// GetUserSessions 获取用户所有会话
	GetUserSessions(userID uint64) ([]*SessionInfo, error)
	// GetConfig 获取配置
	GetConfig() Config
}

// PermissionManager 权限管理器接口
type PermissionManager interface {
	// CheckPermission 检查权限
	CheckPermission(userID uint64, obj string, act string) bool
	// AddRoleForUser 为用户添加角色
	AddRoleForUser(userID uint64, role string) error
	// DeleteRoleForUser 删除用户的角色
	DeleteRoleForUser(userID uint64, role string) error
	// GetRolesForUser 获取用户的所有角色
	GetRolesForUser(userID uint64) ([]string, error)
	// AddPermissionForRole 为角色添加权限
	AddPermissionForRole(role, obj, act string) error
	// DeletePermissionForRole 删除角色的权限
	DeletePermissionForRole(role, obj, act string) error
	// GetPermissionsForRole 获取角色的所有权限
	GetPermissionsForRole(role string) ([][]string, error)
	// GetEnforcer 获取Casbin执行器
	GetEnforcer() *casbin.Enforcer
}

// SecondaryAuthVerifier 二级认证验证器接口
type SecondaryAuthVerifier interface {
	// GenerateCode 生成验证码
	GenerateCode(userID uint64) (string, error)
	// VerifyCode 验证验证码
	VerifyCode(userID uint64, code string) bool
	// SetVerified 设置已验证状态
	SetVerified(userID uint64, duration time.Duration) error
	// IsVerified 检查是否已验证
	IsVerified(userID uint64) bool
}

// ContextKey 上下文键类型
type ContextKey string

const (
	// ContextKeyUserID 用户ID上下文键
	ContextKeyUserID ContextKey = "user_id"
	// ContextKeyDeviceType 设备类型上下文键
	ContextKeyDeviceType ContextKey = "device_type"
	// ContextKeyTokenInfo 令牌信息上下文键
	ContextKeyTokenInfo ContextKey = "token_info"
)
