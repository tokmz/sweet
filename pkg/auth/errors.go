package auth

import (
	"errors"
	"fmt"
)

// AuthError 认证错误类型
type AuthError struct {
	// 错误类型
	Type string
	// 错误消息
	Message string
	// 错误上下文
	Context map[string]interface{}
	// 原始错误
	Wrapped error
}

// Error 实现error接口
func (e *AuthError) Error() string {
	if e.Wrapped != nil {
		return fmt.Sprintf("%s: %s (%v)", e.Type, e.Message, e.Wrapped)
	}
	return fmt.Sprintf("%s: %s", e.Type, e.Message)
}

// Unwrap 实现errors.Unwrap接口
func (e *AuthError) Unwrap() error {
	return e.Wrapped
}

// WithContext 添加上下文信息
func (e *AuthError) WithContext(key string, value interface{}) *AuthError {
	if e.Context == nil {
		e.Context = make(map[string]interface{})
	}
	e.Context[key] = value
	return e
}

// Wrap 包裹原始错误
func (e *AuthError) Wrap(err error) *AuthError {
	e.Wrapped = err
	return e
}

// 错误类型常量
const (
	ErrTypeInvalidToken          = "INVALID_TOKEN"
	ErrTypeTokenExpired          = "TOKEN_EXPIRED"
	ErrTypeUserNotFound          = "USER_NOT_FOUND"
	ErrTypeInvalidCredentials    = "INVALID_CREDENTIALS"
	ErrTypeSessionNotFound       = "SESSION_NOT_FOUND"
	ErrTypePermissionDenied      = "PERMISSION_DENIED"
	ErrTypeInvalidRefreshToken   = "INVALID_REFRESH_TOKEN"
	ErrTypeRefreshTokenExpired   = "REFRESH_TOKEN_EXPIRED"
	ErrTypeServiceNotInitialized = "SERVICE_NOT_INITIALIZED"
	ErrTypeInvalidLoginMode      = "INVALID_LOGIN_MODE"
	ErrTypeSecondaryAuthRequired = "SECONDARY_AUTH_REQUIRED"
	ErrTypeInvalidVerificationCode = "INVALID_VERIFICATION_CODE"
)

// 错误实例
var (
	// ErrInvalidToken 无效的令牌
	ErrInvalidToken = &AuthError{
		Type:    ErrTypeInvalidToken,
		Message: "无效的令牌，请检查令牌格式或重新登录",
	}

	// ErrTokenExpired 令牌已过期
	ErrTokenExpired = &AuthError{
		Type:    ErrTypeTokenExpired,
		Message: "令牌已过期，请刷新令牌或重新登录",
	}

	// ErrUserNotFound 用户不存在
	ErrUserNotFound = &AuthError{
		Type:    ErrTypeUserNotFound,
		Message: "用户不存在，请检查用户ID或联系管理员",
	}

	// ErrInvalidCredentials 无效的凭证
	ErrInvalidCredentials = &AuthError{
		Type:    ErrTypeInvalidCredentials,
		Message: "无效的凭证，请检查用户名和密码",
	}

	// ErrSessionNotFound 会话不存在
	ErrSessionNotFound = &AuthError{
		Type:    ErrTypeSessionNotFound,
		Message: "会话不存在或已过期，请重新登录",
	}

	// ErrPermissionDenied 权限被拒绝
	ErrPermissionDenied = &AuthError{
		Type:    ErrTypePermissionDenied,
		Message: "权限被拒绝，您没有权限执行此操作",
	}

	// ErrInvalidRefreshToken 无效的刷新令牌
	ErrInvalidRefreshToken = &AuthError{
		Type:    ErrTypeInvalidRefreshToken,
		Message: "无效的刷新令牌，请重新登录",
	}

	// ErrRefreshTokenExpired 刷新令牌已过期
	ErrRefreshTokenExpired = &AuthError{
		Type:    ErrTypeRefreshTokenExpired,
		Message: "刷新令牌已过期，请重新登录",
	}

	// ErrServiceNotInitialized 服务未初始化
	ErrServiceNotInitialized = &AuthError{
		Type:    ErrTypeServiceNotInitialized,
		Message: "认证服务未初始化，请先初始化服务",
	}

	// ErrInvalidLoginMode 无效的登录模式
	ErrInvalidLoginMode = &AuthError{
		Type:    ErrTypeInvalidLoginMode,
		Message: "无效的登录模式，请检查配置",
	}

	// ErrSecondaryAuthRequired 需要二级认证
	ErrSecondaryAuthRequired = &AuthError{
		Type:    ErrTypeSecondaryAuthRequired,
		Message: "需要二级认证，请完成验证",
	}

	// ErrInvalidVerificationCode 无效的验证码
	ErrInvalidVerificationCode = &AuthError{
		Type:    ErrTypeInvalidVerificationCode,
		Message: "无效的验证码，请重新获取验证码",
	}
)

// IsAuthError 检查是否为特定类型的认证错误
func IsAuthError(err error, errType string) bool {
	var authErr *AuthError
	if errors.As(err, &authErr) {
		return authErr.Type == errType
	}
	return false
}

// NewAuthError 创建新的认证错误
func NewAuthError(errType, message string) *AuthError {
	return &AuthError{
		Type:    errType,
		Message: message,
		Context: make(map[string]interface{}),
	}
}
