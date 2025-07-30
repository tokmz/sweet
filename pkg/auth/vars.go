package auth

import (
	"errors"

	"github.com/golang-jwt/jwt/v5"
)

// UserType 用户类型
type UserType string

const (
	// FrontendUser 前端用户
	FrontendUser UserType = "frontend"
	// BackendUser 后端用户
	BackendUser UserType = "backend"
	// TokenCache 令牌缓存
	TokenCache = "token::%s::%d::%s::%d::%s" // 用户类型 Token类型 用户ID 用户名 角色ID 设备类型
)

type Claims struct {
	// 用户ID
	Uid int64 `json:"uid"`
	// 用户名
	Username string `json:"username"`
	// 角色ID
	Rid int64 `json:"rid"`
	// 缓冲时间
	BufferTime int64 `json:"buffer_time"`
	// 设备类型（pc,ios,android）
	DeviceType string `json:"device_type"`
	// 用户类型（frontend,backend）
	UserType UserType `json:"user_type"`
	jwt.RegisteredClaims
}

// jwt 相关错误处理
var (
	// ErrTokenFormat 格式错误
	ErrTokenFormat = errors.New("token格式错误")
	// ErrTokenExpired 过期
	ErrTokenExpired = errors.New("token已过期")
	// ErrTokenNotEffective token还未生效
	ErrTokenNotEffective = errors.New("token尚未生效")
	// ErrTokenSignVerify 签名验证错误
	ErrTokenSignVerify = errors.New("token签名验证错误")
	// ErrInvalidClaims 无效的claims类型
	ErrInvalidClaims = errors.New("无效的claims类型")
	// ErrTokenInvalid token无效
	ErrTokenInvalid = errors.New("token无效")
	// ErrNeedLogin 需要重新登录
	ErrNeedLogin = errors.New("需要重新登录")
	// ErrUserType 错误的用户类型
	ErrUserType = errors.New("错误的用户类型")
	// ErrUserAlreadyLogin 用户已在其他设备登录
	ErrUserAlreadyLogin = errors.New("用户已在其他设备登录")
)
