package auth

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"

	"github.com/redis/go-redis/v9"
)

type Jwt struct {
	key        []byte // jwt 秘钥
	issuer     string // 签发
	subject    string
	bufferTime time.Duration // 时间戳
	expireTime time.Duration // 过期时间
	client     *redis.Client
}

type JwtConfig struct {
	SecretKey  string // 密钥
	Issuer     string // 签发者
	Subject    string
	BufferTime string // 刷新阈值时间 1s、1h、1d
	ExpireTime string // 令牌过期时间 1s、1h、1d
}

var (
	localJwt *Jwt
)

func NewJwt(cfg *JwtConfig, client *redis.Client) error {
	if client == nil {
		return errors.New("redis client is nil")
	}

	// 处理过期和缓冲时间
	bufferTime, err := time.ParseDuration(cfg.BufferTime)
	if err != nil {
		return fmt.Errorf("解析缓冲时间失败: %w", err)
	}

	expireTime, err := time.ParseDuration(cfg.ExpireTime)
	if err != nil {
		return fmt.Errorf("解析过期时间失败: %w", err)
	}

	localJwt = &Jwt{
		key:        []byte(cfg.SecretKey),
		issuer:     cfg.Issuer,
		subject:    cfg.Subject,
		bufferTime: bufferTime,
		expireTime: expireTime,
		client:     client,
	}
	return nil
}

// GenerateToken 生成Token
func GenerateToken(ctx context.Context, uid int64, username string, rid int64, deviceType string, userType UserType) (string, error) {
	bufferTime := time.Now().Add(localJwt.bufferTime)
	expire := time.Now().Add(localJwt.expireTime)
	claims := &Claims{
		Uid:        uid,
		Username:   username,
		Rid:        rid,
		DeviceType: deviceType,
		UserType:   userType,
		BufferTime: bufferTime.Unix(),
		RegisteredClaims: jwt.RegisteredClaims{
			Issuer:    localJwt.issuer,
			Subject:   localJwt.subject,
			ExpiresAt: jwt.NewNumericDate(expire),
			NotBefore: jwt.NewNumericDate(time.Now()),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
		},
	}

	if token, err := jwt.NewWithClaims(jwt.SigningMethodHS256, claims).SignedString(localJwt.key); err != nil {
		return "", fmt.Errorf("生成token失败: %w", err)
	} else {
		key := fmt.Sprintf(TokenCache, userType, uid, username, rid, deviceType)
		if err = localJwt.client.Set(ctx, key, token, localJwt.expireTime).Err(); err != nil {
			return "", fmt.Errorf("缓存token失败: %w", err)
		}
		return token, nil
	}
}

// ParseToken 解析Token
func ParseToken(token string) (*Claims, error) {
	// 检查token是否为空
	if token == "" {
		return nil, errors.New("token不能为空")
	}

	claimsType, err := jwt.ParseWithClaims(token, &Claims{}, func(token *jwt.Token) (interface{}, error) {
		// 验证签名方法
		if method, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("不支持的签名方法: %v", token.Header["alg"])
		} else if method != jwt.SigningMethodHS256 {
			return nil, fmt.Errorf("期望HS256签名方法，实际: %v", method.Alg())
		}
		return localJwt.key, nil
	})

	// 处理解析错误
	if err != nil {
		switch {
		case errors.Is(err, jwt.ErrTokenMalformed):
			return nil, ErrTokenFormat
		case errors.Is(err, jwt.ErrTokenExpired):
			return nil, ErrTokenExpired
		case errors.Is(err, jwt.ErrTokenNotValidYet):
			return nil, ErrTokenNotEffective
		case errors.Is(err, jwt.ErrTokenSignatureInvalid):
			return nil, ErrTokenSignVerify
		case errors.Is(err, jwt.ErrTokenInvalidClaims):
			return nil, ErrInvalidClaims
		default:
			return nil, fmt.Errorf("解析token失败: %w", err)
		}
	}

	// 检查token是否有效
	if !claimsType.Valid {
		return nil, ErrTokenInvalid
	}

	// 类型断言获取Claims
	claims, ok := claimsType.Claims.(*Claims)
	if !ok {
		return nil, ErrInvalidClaims
	}

	return claims, nil
}

type CheckTokenResult struct {
	NeedRefresh bool
	Token       string
	Claims      *Claims
}

func CheckToken(ctx context.Context, token string) (*CheckTokenResult, error) {
	if claims, err := ParseToken(token); err != nil {
		return nil, err
	} else {
		var result string
		key := fmt.Sprintf(TokenCache, claims.UserType, claims.Uid, claims.Username, claims.Rid, claims.DeviceType)
		if result, err = localJwt.client.Get(ctx, key).Result(); err != nil {
			if errors.Is(err, redis.Nil) {
				// 需要重新登录
				return nil, ErrNeedLogin
			}
			return nil, fmt.Errorf("缓存查询失败: %w", err)
		}

		if result != token {
			return nil, ErrUserAlreadyLogin
		}

		if claims.BufferTime < time.Now().Unix() {
			// 需要刷新token
			newToken, err := GenerateToken(ctx, claims.Uid, claims.Username, claims.Rid, claims.DeviceType, claims.UserType)
			if err != nil {
				return nil, fmt.Errorf("刷新token失败: %w", err)
			}
			return &CheckTokenResult{
				NeedRefresh: true,
				Token:       newToken,
				Claims:      claims,
			}, nil
		}

		return &CheckTokenResult{
			NeedRefresh: false,
			Token:       token,
			Claims:      claims,
		}, nil
	}
}
