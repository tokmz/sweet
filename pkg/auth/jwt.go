package auth

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

// JWTClaims JWT声明结构
type JWTClaims struct {
	jwt.RegisteredClaims
	UserID    int64                  `json:"uid"`
	Username  string                 `json:"username"`
	Device    string                 `json:"device"`
	ExtraData map[string]interface{} `json:"extra,omitempty"`
}

// JWTService JWT服务
type JWTService struct {
	secret        string
	tokenExpire   time.Duration
	refreshExpire time.Duration
}

// NewJWTService 创建JWT服务
func NewJWTService(secret string, tokenExpire, refreshExpire time.Duration) *JWTService {
	return &JWTService{
		secret:        secret,
		tokenExpire:   tokenExpire,
		refreshExpire: refreshExpire,
	}
}

// GenerateToken 生成JWT Token
func (s *JWTService) GenerateToken(userID int64, username string, device string, extraData map[string]interface{}) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.tokenExpire)
	
	// 创建JWT声明
	claims := JWTClaims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expiresAt),
			IssuedAt:  jwt.NewNumericDate(now),
			NotBefore: jwt.NewNumericDate(now),
			Issuer:    "sweet-auth",
			Subject:   username,
			ID:        uuid.New().String(),
		},
		UserID:    userID,
		Username:  username,
		Device:    device,
		ExtraData: extraData,
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// GenerateRefreshToken 生成刷新Token
func (s *JWTService) GenerateRefreshToken(userID int64, username string, device string) (string, error) {
	now := time.Now()
	expiresAt := now.Add(s.refreshExpire)
	
	// 创建JWT声明
	claims := jwt.RegisteredClaims{
		ExpiresAt: jwt.NewNumericDate(expiresAt),
		IssuedAt:  jwt.NewNumericDate(now),
		NotBefore: jwt.NewNumericDate(now),
		Issuer:    "sweet-auth-refresh",
		Subject:   username,
		ID:        uuid.New().String(),
	}

	// 创建Token
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	// 签名Token
	tokenString, err := token.SignedString([]byte(s.secret))
	if err != nil {
		return "", err
	}

	return tokenString, nil
}

// ParseToken 解析JWT Token
func (s *JWTService) ParseToken(tokenString string) (*JWTClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &JWTClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrTokenExpired
		}
		return nil, ErrTokenInvalid
	}

	// 验证Token有效性
	if !token.Valid {
		return nil, ErrTokenInvalid
	}

	// 获取声明
	claims, ok := token.Claims.(*JWTClaims)
	if !ok {
		return nil, ErrTokenInvalid
	}

	return claims, nil
}

// ParseRefreshToken 解析刷新Token
func (s *JWTService) ParseRefreshToken(tokenString string) (*jwt.RegisteredClaims, error) {
	// 解析Token
	token, err := jwt.ParseWithClaims(tokenString, &jwt.RegisteredClaims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名算法
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.secret), nil
	})

	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired) {
			return nil, ErrRefreshTokenExpired
		}
		return nil, ErrRefreshTokenInvalid
	}

	// 验证Token有效性
	if !token.Valid {
		return nil, ErrRefreshTokenInvalid
	}

	// 获取声明
	claims, ok := token.Claims.(*jwt.RegisteredClaims)
	if !ok {
		return nil, ErrRefreshTokenInvalid
	}

	return claims, nil
}

// GenerateTokenByStyle 根据指定风格生成Token
func (s *JWTService) GenerateTokenByStyle(style TokenStyle, userID int64, username string, device string, extraData map[string]interface{}) (string, error) {
	switch style {
	case TokenStyleUUID:
		// 使用UUID生成Token
		return uuid.New().String(), nil
	case TokenStyleSimple:
		// 使用简单字符串生成Token
		return generateSimpleToken(userID, device), nil
	case TokenStyleJWT:
		// 使用标准JWT格式
		return s.GenerateToken(userID, username, device, extraData)
	case TokenStyleJWTMixed:
		// 使用混合JWT格式(JWT+UUID)
		jwtToken, err := s.GenerateToken(userID, username, device, extraData)
		if err != nil {
			return "", err
		}
		return jwtToken + "." + uuid.New().String(), nil
	case TokenStyleJWTUUID:
		// 使用UUID作为JWT的jti
		return s.GenerateToken(userID, username, device, extraData)
	default:
		return "", ErrNotImplemented
	}
}

// 生成简单Token
func generateSimpleToken(userID int64, device string) string {
	// 简单格式: 时间戳-用户ID-随机字符串
	timestamp := time.Now().UnixNano()
	randomStr := uuid.New().String()
	return time.Now().Format("20060102150405") + "-" + device + "-" + randomStr
}
