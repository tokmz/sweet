package auth

import (
	"context"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// DefaultAuthService 默认认证服务实现
type DefaultAuthService struct {
	// config 配置
	config Config
	// tokenManager 令牌管理器
	tokenManager TokenManager
	// sessionStore 会话存储
	sessionStore SessionStore
}

// NewAuthService 创建认证服务
func NewAuthService(config Config, sessionStore SessionStore) *DefaultAuthService {
	// 创建令牌管理器
	tokenManager := NewJWTTokenManager(
		config.SecretKey,
		config.AccessTokenExpiry,
		config.RefreshTokenExpiry,
		sessionStore,
	)

	return &DefaultAuthService{
		config:       config,
		tokenManager: tokenManager,
		sessionStore: sessionStore,
	}
}

// Login 登录
func (s *DefaultAuthService) Login(userID uint64, deviceType string) (accessToken string, refreshToken string, err error) {
	// 生成访问令牌
	accessToken, err = s.tokenManager.GenerateToken(userID, deviceType, TokenTypeAccess)
	if err != nil {
		return "", "", err
	}

	// 生成刷新令牌
	refreshToken, err = s.tokenManager.GenerateToken(userID, deviceType, TokenTypeRefresh)
	if err != nil {
		return "", "", err
	}

	// 创建会话信息
	now := time.Now()
	session := &SessionInfo{
		UserID:       userID,
		DeviceType:   deviceType,
		LoginTime:    now.Unix(),
		ExpireTime:   now.Add(s.config.RefreshTokenExpiry).Unix(),
		RefreshToken: refreshToken,
		// IP和UserAgent应该从请求中获取，这里暂时留空
		IP:        "",
		UserAgent: "",
	}

	// 根据登录模式处理会话
	switch s.config.LoginMode {
	case LoginModeSingle:
		// 单端登录：清除用户所有会话
		var sessions []*SessionInfo
		sessions, err = s.sessionStore.GetUserSessions(userID)
		if err == nil {
			for _, oldSession := range sessions {
				_ = s.sessionStore.RemoveSession(oldSession.UserID, oldSession.DeviceType)
			}
		}
		// 保存新会话
		err = s.sessionStore.SaveSession(session)
		if err != nil {
			return "", "", err
		}

	case LoginModeMulti:
		// 多端登录：直接保存会话
		err = s.sessionStore.SaveSession(session)
		if err != nil {
			return "", "", err
		}

	case LoginModeMutex:
		// 同端互斥登录：清除同设备类型的所有会话
		sessions, err := s.sessionStore.GetDeviceTypeSessions(deviceType)
		if err == nil {
			for _, oldSession := range sessions {
				_ = s.sessionStore.RemoveSession(oldSession.UserID, oldSession.DeviceType)
			}
		}
		// 保存新会话
		err = s.sessionStore.SaveSession(session)
		if err != nil {
			return "", "", err
		}

	default:
		return "", "", ErrInvalidLoginMode
	}

	return accessToken, refreshToken, nil
}

// Logout 登出
func (s *DefaultAuthService) Logout(userID uint64, deviceType string) error {
	// 撤销令牌
	return s.tokenManager.RevokeToken(userID, deviceType)
}

// RefreshToken 刷新令牌
func (s *DefaultAuthService) RefreshToken(refreshToken string) (accessToken string, newRefreshToken string, err error) {
	// 刷新令牌
	return s.tokenManager.RefreshToken(refreshToken)
}

// VerifyToken 验证令牌
func (s *DefaultAuthService) VerifyToken(token string) (*TokenInfo, error) {
	// 解析令牌
	return s.tokenManager.ParseToken(token)
}

// GetUserSessions 获取用户所有会话
func (s *DefaultAuthService) GetUserSessions(userID uint64) ([]*SessionInfo, error) {
	// 获取用户所有会话
	return s.sessionStore.GetUserSessions(userID)
}

// GetConfig 获取配置
func (s *DefaultAuthService) GetConfig() Config {
	return s.config
}

// RedisSecondaryAuthVerifier Redis二级认证验证器实现
type RedisSecondaryAuthVerifier struct {
	// client Redis客户端
	client *redis.Client
	// keyPrefix 键前缀
	keyPrefix string
	// codeExpiry 验证码过期时间
	codeExpiry time.Duration
}

// NewRedisSecondaryAuthVerifier 创建Redis二级认证验证器
func NewRedisSecondaryAuthVerifier(client *redis.Client, codeExpiry time.Duration) *RedisSecondaryAuthVerifier {
	return &RedisSecondaryAuthVerifier{
		client:     client,
		keyPrefix:  "sweet:auth:secondary:",
		codeExpiry: codeExpiry,
	}
}

// codeKey 生成验证码键
func (v *RedisSecondaryAuthVerifier) codeKey(userID uint64) string {
	return fmt.Sprintf("%scode:%d", v.keyPrefix, userID)
}

// verifiedKey 生成已验证键
func (v *RedisSecondaryAuthVerifier) verifiedKey(userID uint64) string {
	return fmt.Sprintf("%sverified:%d", v.keyPrefix, userID)
}

// GenerateCode 生成验证码
func (v *RedisSecondaryAuthVerifier) GenerateCode(userID uint64) (string, error) {
	// 生成6位随机验证码
	code := fmt.Sprintf("%06d", time.Now().UnixNano()%1000000)

	ctx := context.Background()
	// 保存验证码
	key := v.codeKey(userID)
	err := v.client.Set(ctx, key, code, v.codeExpiry).Err()
	if err != nil {
		return "", err
	}

	return code, nil
}

// VerifyCode 验证验证码
func (v *RedisSecondaryAuthVerifier) VerifyCode(userID uint64, code string) bool {
	ctx := context.Background()
	// 获取验证码
	key := v.codeKey(userID)
	savedCode, err := v.client.Get(ctx, key).Result()
	if err != nil {
		return false
	}

	// 验证码匹配
	if savedCode == code {
		// 删除验证码
		_ = v.client.Del(ctx, key).Err()
		return true
	}

	return false
}

// SetVerified 设置已验证状态
func (v *RedisSecondaryAuthVerifier) SetVerified(userID uint64, duration time.Duration) error {
	ctx := context.Background()
	// 设置已验证状态
	key := v.verifiedKey(userID)
	return v.client.Set(ctx, key, "1", duration).Err()
}

// IsVerified 检查是否已验证
func (v *RedisSecondaryAuthVerifier) IsVerified(userID uint64) bool {
	ctx := context.Background()
	// 检查已验证状态
	key := v.verifiedKey(userID)
	result, err := v.client.Exists(ctx, key).Result()
	if err != nil {
		return false
	}

	return result > 0
}
