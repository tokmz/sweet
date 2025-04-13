// Package auth 提供身份认证与权限控制功能
package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var (
	// ErrTokenNotFound 令牌不存在
	ErrTokenNotFound = errors.New("token not found")
	// ErrTokenInvalid 令牌无效
	ErrTokenInvalid = errors.New("token invalid")
	// ErrTokenExpired 令牌过期
	ErrTokenExpired = errors.New("token expired")
	// ErrNotLogin 未登录
	ErrNotLogin = errors.New("not login")
	// ErrKickedOut 被踢下线
	ErrKickedOut = errors.New("kicked out")
)

const (
	// LoginTypeDefault 默认登录类型，同端互斥
	LoginTypeDefault = "default"
	// LoginTypeSingle 单端登录
	LoginTypeSingle = "single"
	// LoginTypeMulti 多端登录
	LoginTypeMulti = "multi"
	// LoginTypeExclusive 同端互斥登录
	LoginTypeExclusive = "exclusive"
)

// 存储模式
const (
	// StoreMemory 使用内存存储
	StoreMemory = "memory"
	// StoreRedis 使用Redis存储
	StoreRedis = "redis"
)

// JWTClaims 自定义JWT声明
type JWTClaims struct {
	// 标准jwt声明
	jwt.RegisteredClaims
	// 用户ID
	UserID int64 `json:"user_id"`
	// 设备类型
	Device string `json:"device"`
	// 登录类型
	LoginType string `json:"login_type"`
	// 用户角色
	Roles []string `json:"roles,omitempty"`
}

// TokenInfo 令牌信息
type TokenInfo struct {
	// UserID 用户ID
	UserID int64
	// Token 令牌值
	Token string
	// LoginType 登录类型
	LoginType string
	// Device 登录设备
	Device string
	// ExpireAt 过期时间
	ExpireAt time.Time
	// LastActiveTime 最后活跃时间
	LastActiveTime time.Time
	// Roles 用户角色
	Roles []string
}

// TokenStore 令牌存储接口
type TokenStore interface {
	// StoreToken 存储令牌
	StoreToken(ctx context.Context, token TokenInfo) error
	// GetToken 获取令牌信息
	GetToken(ctx context.Context, token string) (*TokenInfo, error)
	// GetUserTokens 获取用户的所有令牌
	GetUserTokens(ctx context.Context, userID int64) ([]TokenInfo, error)
	// GetUserTokensByDevice 获取用户指定设备的所有令牌
	GetUserTokensByDevice(ctx context.Context, userID int64, device string) ([]TokenInfo, error)
	// RemoveToken 移除令牌
	RemoveToken(ctx context.Context, token string) error
	// RemoveUserTokens 移除用户的所有令牌
	RemoveUserTokens(ctx context.Context, userID int64) error
	// RemoveUserTokensByDevice 移除用户指定设备的所有令牌
	RemoveUserTokensByDevice(ctx context.Context, userID int64, device string) error
	// UpdateActiveTime 更新令牌活跃时间
	UpdateActiveTime(ctx context.Context, token string, time time.Time) error
}

// Config 身份认证配置
type Config struct {
	// StoreType 存储类型，默认memory
	StoreType string
	// AccessTokenTimeout 访问令牌有效期，默认30分钟
	AccessTokenTimeout time.Duration
	// RefreshTokenTimeout 刷新令牌有效期，默认7天
	RefreshTokenTimeout time.Duration
	// TokenIssuer 令牌签发者
	TokenIssuer string
	// TokenAudience 令牌接收者
	TokenAudience string
	// JWTSigningKey JWT签名密钥
	JWTSigningKey []byte
	// AllowConcurrentLogin 是否允许同一账号并发登录，默认false
	AllowConcurrentLogin bool
	// IsShare 是否共享token，默认false
	IsShare bool
	// AutoRenew 是否自动续期，默认true
	AutoRenew bool
	// TokenPrefix 令牌前缀，默认"auth:"
	TokenPrefix string
}

// DefaultConfig 默认配置
func DefaultConfig() Config {
	return Config{
		StoreType:            StoreMemory,
		AccessTokenTimeout:   30 * time.Minute,
		RefreshTokenTimeout:  7 * 24 * time.Hour,
		TokenIssuer:          "auth",
		TokenAudience:        "users",
		JWTSigningKey:        []byte("default-secret-key"),
		AllowConcurrentLogin: false,
		IsShare:              false,
		AutoRenew:            true,
		TokenPrefix:          "auth:",
	}
}

// Auth 身份认证实例
type Auth struct {
	// Config 配置
	Config Config
	// Store 存储
	Store TokenStore
	// Enforcer Casbin执行器
	Enforcer CasbinEnforcer
}

// New 创建身份认证实例
func New(config Config) *Auth {
	if config.StoreType == "" || len(config.JWTSigningKey) == 0 {
		config = DefaultConfig()
	}

	var store TokenStore
	if config.StoreType == StoreMemory {
		store = NewMemoryStore()
	} else if config.StoreType == StoreRedis {
		// Redis存储，需要单独实现
		store = NewMemoryStore()
	} else {
		// 默认内存存储
		store = NewMemoryStore()
	}

	// 创建Casbin执行器
	enforcer := NewCasbinEnforcer()

	return &Auth{
		Config:   config,
		Store:    store,
		Enforcer: enforcer,
	}
}

// Login 登录并返回token
func (a *Auth) Login(ctx context.Context, userID int64, device string, loginType string) (string, error) {
	if userID <= 0 {
		return "", errors.New("invalid user id")
	}

	if device == "" {
		device = "unknown"
	}

	if loginType == "" {
		loginType = LoginTypeDefault
	}

	// 根据登录类型处理已有的token
	switch loginType {
	case LoginTypeSingle:
		// 单端登录：一个用户只能有一个token
		err := a.Store.RemoveUserTokens(ctx, userID)
		if err != nil {
			return "", err
		}
	case LoginTypeExclusive:
		// 同端互斥登录：一个用户在一个设备上只能有一个token
		err := a.Store.RemoveUserTokensByDevice(ctx, userID, device)
		if err != nil {
			return "", err
		}
	case LoginTypeMulti:
		// 多端登录：不做处理，允许多个token同时存在
	default:
		// 默认同端互斥
		err := a.Store.RemoveUserTokensByDevice(ctx, userID, device)
		if err != nil {
			return "", err
		}
	}

	// 创建过期时间
	expiresAt := time.Now().Add(a.Config.AccessTokenTimeout)

	// 创建JWT
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    a.Config.TokenIssuer,
			Audience:  jwt.ClaimStrings{a.Config.TokenAudience},
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			NotBefore: jwt.NewNumericDate(time.Now()),
			ID:        fmt.Sprintf("%d-%s-%d", userID, device, time.Now().UnixNano()),
		},
		UserID:    userID,
		Device:    device,
		LoginType: loginType,
		Roles:     []string{}, // 初始化为空角色列表
	}

	// 生成token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	tokenString, err := token.SignedString(a.Config.JWTSigningKey)
	if err != nil {
		return "", err
	}

	// 创建token信息
	tokenInfo := TokenInfo{
		UserID:         userID,
		Token:          tokenString,
		LoginType:      loginType,
		Device:         device,
		ExpireAt:       expiresAt,
		LastActiveTime: time.Now(),
		Roles:          []string{},
	}

	// 存储token
	err = a.Store.StoreToken(ctx, tokenInfo)
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// Logout 注销登录
func (a *Auth) Logout(ctx context.Context, token string) error {
	if token == "" {
		return ErrNotLogin
	}

	return a.Store.RemoveToken(ctx, token)
}

// LogoutByUserID 根据用户ID注销登录
func (a *Auth) LogoutByUserID(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}

	return a.Store.RemoveUserTokens(ctx, userID)
}

// LogoutByUserIDAndDevice 根据用户ID和设备注销登录
func (a *Auth) LogoutByUserIDAndDevice(ctx context.Context, userID int64, device string) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}

	if device == "" {
		device = "unknown"
	}

	return a.Store.RemoveUserTokensByDevice(ctx, userID, device)
}

// ParseToken 解析JWT token
func (a *Auth) ParseToken(tokenString string) (*JWTClaims, error) {
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return a.Config.JWTSigningKey, nil
	})

	if err != nil {
		return nil, err
	}

	if claims, ok := token.Claims.(*JWTClaims); ok && token.Valid {
		return claims, nil
	}

	return nil, ErrTokenInvalid
}

// CheckLogin 检查登录状态
func (a *Auth) CheckLogin(ctx context.Context, token string) (bool, error) {
	if token == "" {
		return false, ErrNotLogin
	}

	// 解析token
	claims, err := a.ParseToken(token)
	if err != nil {
		return false, err
	}
	_ = claims // 使用claims以避免未使用的变量警告

	// 检查token是否存在
	tokenInfo, err := a.Store.GetToken(ctx, token)
	if err != nil {
		return false, err
	}

	// 检查token是否过期
	if time.Now().After(tokenInfo.ExpireAt) {
		return false, ErrTokenExpired
	}

	// 自动续期
	if a.Config.AutoRenew {
		// 如果距离上次活跃时间超过了10分钟，更新活跃时间
		if time.Since(tokenInfo.LastActiveTime) > 10*time.Minute {
			now := time.Now()
			a.Store.UpdateActiveTime(ctx, token, now)

			// 如果剩余有效期不足一半，延长有效期
			halfTimeout := a.Config.AccessTokenTimeout / 2
			if tokenInfo.ExpireAt.Sub(now) < halfTimeout {
				// 创建新的过期时间
				newExpireAt := now.Add(a.Config.AccessTokenTimeout)

				// 这里需要重新创建一个新的TokenInfo，否则无法更新expireAt
				newInfo := *tokenInfo
				newInfo.ExpireAt = newExpireAt
				a.Store.StoreToken(ctx, newInfo)

				// 生成新的JWT token（可选，但建议实现）
				// ...这里需要更新JWT token的实现逻辑
			}
		}
	}

	return true, nil
}

// GetUserID 根据token获取用户ID
func (a *Auth) GetUserID(ctx context.Context, token string) (int64, error) {
	if token == "" {
		return 0, ErrNotLogin
	}

	// 解析token
	claims, err := a.ParseToken(token)
	if err != nil {
		return 0, err
	}

	return claims.UserID, nil
}

// KickOut 踢人下线
func (a *Auth) KickOut(ctx context.Context, token string) error {
	if token == "" {
		return ErrNotLogin
	}

	return a.Store.RemoveToken(ctx, token)
}

// KickOutByUserID 踢指定用户下线
func (a *Auth) KickOutByUserID(ctx context.Context, userID int64) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}

	return a.Store.RemoveUserTokens(ctx, userID)
}

// KickOutByUserIDAndDevice 踢指定用户在指定设备下线
func (a *Auth) KickOutByUserIDAndDevice(ctx context.Context, userID int64, device string) error {
	if userID <= 0 {
		return errors.New("invalid user id")
	}

	if device == "" {
		device = "unknown"
	}

	return a.Store.RemoveUserTokensByDevice(ctx, userID, device)
}

// GetActiveTokens 获取有效的token列表
func (a *Auth) GetActiveTokens(ctx context.Context, userID int64) ([]TokenInfo, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	tokens, err := a.Store.GetUserTokens(ctx, userID)
	if err != nil {
		return nil, err
	}

	// 过滤出未过期的token
	validTokens := make([]TokenInfo, 0, len(tokens))
	for _, token := range tokens {
		if time.Now().Before(token.ExpireAt) {
			validTokens = append(validTokens, token)
		}
	}

	return validTokens, nil
}

// GetActiveTokensByDevice 获取指定设备有效的token列表
func (a *Auth) GetActiveTokensByDevice(ctx context.Context, userID int64, device string) ([]TokenInfo, error) {
	if userID <= 0 {
		return nil, errors.New("invalid user id")
	}

	if device == "" {
		device = "unknown"
	}

	tokens, err := a.Store.GetUserTokensByDevice(ctx, userID, device)
	if err != nil {
		return nil, err
	}

	// 过滤出未过期的token
	validTokens := make([]TokenInfo, 0, len(tokens))
	for _, token := range tokens {
		if time.Now().Before(token.ExpireAt) {
			validTokens = append(validTokens, token)
		}
	}

	return validTokens, nil
}
