package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

// CustomClaims 自定义JWT声明结构体
type CustomClaims struct {
	// UserID 用户ID
	UserID uint64 `json:"user_id"`
	// DeviceType 设备类型
	DeviceType string `json:"device_type"`
	// TokenType 令牌类型
	TokenType string `json:"token_type"`
	// IssuedAt 签发时间
	IssuedAt int64 `json:"issued_at"`
	// ExpiresAt 过期时间
	ExpiresAt int64 `json:"expires_at"`
	jwt.RegisteredClaims
}

// JWTTokenManager JWT令牌管理器
type JWTTokenManager struct {
	// secretKey 密钥
	secretKey string
	// accessExpiry 访问令牌过期时间
	accessExpiry time.Duration
	// refreshExpiry 刷新令牌过期时间
	refreshExpiry time.Duration
	// sessionStore 会话存储
	sessionStore SessionStore
}

// NewJWTTokenManager 创建JWT令牌管理器
func NewJWTTokenManager(secretKey string, accessExpiry, refreshExpiry time.Duration, sessionStore SessionStore) *JWTTokenManager {
	return &JWTTokenManager{
		secretKey:     secretKey,
		accessExpiry:  accessExpiry,
		refreshExpiry: refreshExpiry,
		sessionStore:  sessionStore,
	}
}

// GenerateToken 生成令牌
func (m *JWTTokenManager) GenerateToken(userID uint64, deviceType string, tokenType string) (string, error) {
	// 参数验证
	if userID == 0 {
		return "", NewAuthError(ErrTypeInvalidCredentials, "用户ID不能为空").WithContext("user_id", userID)
	}
	if deviceType == "" {
		return "", NewAuthError(ErrTypeInvalidCredentials, "设备类型不能为空").WithContext("device_type", deviceType)
	}

	// 确定过期时间
	var expiry time.Duration
	if tokenType == TokenTypeAccess {
		expiry = m.accessExpiry
	} else if tokenType == TokenTypeRefresh {
		expiry = m.refreshExpiry
	} else {
		return "", NewAuthError(ErrTypeInvalidToken, "未知的令牌类型").WithContext("token_type", tokenType)
	}

	// 创建JWT声明
	now := time.Now()
	claims := &CustomClaims{
		UserID:     userID,
		DeviceType: deviceType,
		TokenType:  tokenType,
		IssuedAt:   now.Unix(),
		ExpiresAt:  now.Add(expiry).Unix(),
	}

	// 创建令牌
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名令牌
	tokenString, err := token.SignedString([]byte(m.secretKey))
	if err != nil {
		return "", NewAuthError(ErrTypeInvalidToken, "令牌签名失败").Wrap(err)
	}

	return tokenString, nil
}

// ParseToken 解析令牌
func (m *JWTTokenManager) ParseToken(tokenString string) (*TokenInfo, error) {
	// 解析令牌
	var customClaims CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &customClaims, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, NewAuthError(ErrTypeInvalidToken, "意外的签名方法").WithContext("alg", token.Header["alg"])
		}
		return []byte(m.secretKey), nil
	})

	if err != nil {
		// 检查是否为过期错误
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired.WithContext("token", tokenString)
		}
		return nil, ErrInvalidToken.WithContext("token", tokenString).Wrap(err)
	}

	// 验证令牌
	if !token.Valid {
		return nil, ErrInvalidToken.WithContext("token", tokenString)
	}

	// 检查是否过期
	if time.Now().Unix() > customClaims.ExpiresAt {
		return nil, ErrTokenExpired.WithContext("token", tokenString).WithContext("expires_at", customClaims.ExpiresAt)
	}

	// 返回令牌信息
	return &TokenInfo{
		UserID:     customClaims.UserID,
		DeviceType: customClaims.DeviceType,
		TokenType:  customClaims.TokenType,
		ExpiresAt:  customClaims.ExpiresAt,
		IssuedAt:   customClaims.IssuedAt,
	}, nil
}

// RefreshToken 刷新令牌
func (m *JWTTokenManager) RefreshToken(refreshToken string) (string, string, error) {
	// 解析刷新令牌
	tokenInfo, err := m.ParseToken(refreshToken)
	if err != nil {
		var authErr *AuthError
		if errors.As(err, &authErr) && authErr.Type == ErrTypeTokenExpired {
			return "", "", ErrRefreshTokenExpired.WithContext("refresh_token", refreshToken)
		}
		return "", "", ErrInvalidRefreshToken.WithContext("refresh_token", refreshToken).Wrap(err)
	}

	// 验证令牌类型
	if tokenInfo.TokenType != TokenTypeRefresh {
		return "", "", ErrInvalidRefreshToken.WithContext("token_type", tokenInfo.TokenType)
	}

	// 从会话存储中获取会话
	session, err := m.sessionStore.GetSession(tokenInfo.UserID, tokenInfo.DeviceType)
	if err != nil {
		if err == ErrSessionNotFound {
			return "", "", ErrInvalidRefreshToken.WithContext("user_id", tokenInfo.UserID).WithContext("device_type", tokenInfo.DeviceType)
		}
		return "", "", NewAuthError(ErrTypeSessionNotFound, "获取会话失败").WithContext("user_id", tokenInfo.UserID).WithContext("device_type", tokenInfo.DeviceType).Wrap(err)
	}

	// 验证刷新令牌是否匹配
	if session.RefreshToken != refreshToken {
		return "", "", ErrInvalidRefreshToken.WithContext("user_id", tokenInfo.UserID).WithContext("device_type", tokenInfo.DeviceType)
	}

	// 生成新的访问令牌
	newAccessToken, err := m.GenerateToken(tokenInfo.UserID, tokenInfo.DeviceType, TokenTypeAccess)
	if err != nil {
		return "", "", err
	}

	// 生成新的刷新令牌
	newRefreshToken, err := m.GenerateToken(tokenInfo.UserID, tokenInfo.DeviceType, TokenTypeRefresh)
	if err != nil {
		return "", "", err
	}

	// 更新会话
	now := time.Now()
	session.RefreshToken = newRefreshToken
	session.ExpireTime = now.Add(m.refreshExpiry).Unix()

	// 保存会话
	err = m.sessionStore.SaveSession(session)
	if err != nil {
		return "", "", err
	}

	return newAccessToken, newRefreshToken, nil
}

// RevokeToken 撤销令牌
func (m *JWTTokenManager) RevokeToken(userID uint64, deviceType string) error {
	// 从会话存储中删除会话
	return m.sessionStore.RemoveSession(userID, deviceType)
}
