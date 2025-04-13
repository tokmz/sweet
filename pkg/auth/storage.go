package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"time"

	"github.com/redis/go-redis/v9"
)

// MemoryStorage 内存存储实现
type MemoryStorage struct {
	tokens      map[string]TokenInfo
	secondAuth  map[string]int64
	tokensMutex sync.RWMutex
	authMutex   sync.RWMutex
}

// NewMemoryStorage 创建内存存储
func NewMemoryStorage() *MemoryStorage {
	return &MemoryStorage{
		tokens:     make(map[string]TokenInfo),
		secondAuth: make(map[string]int64),
	}
}

// SaveToken 保存Token信息
func (s *MemoryStorage) SaveToken(token string, info TokenInfo) error {
	s.tokensMutex.Lock()
	defer s.tokensMutex.Unlock()
	s.tokens[token] = info
	return nil
}

// GetToken 获取Token信息
func (s *MemoryStorage) GetToken(token string) (TokenInfo, error) {
	s.tokensMutex.RLock()
	defer s.tokensMutex.RUnlock()
	info, ok := s.tokens[token]
	if !ok {
		return TokenInfo{}, ErrTokenNotFound
	}
	return info, nil
}

// RemoveToken 移除Token
func (s *MemoryStorage) RemoveToken(token string) error {
	s.tokensMutex.Lock()
	defer s.tokensMutex.Unlock()
	delete(s.tokens, token)
	return nil
}

// GetUserTokens 获取用户的所有Token
func (s *MemoryStorage) GetUserTokens(userID int64) ([]string, error) {
	s.tokensMutex.RLock()
	defer s.tokensMutex.RUnlock()
	var tokens []string
	for token, info := range s.tokens {
		if info.UserID == userID {
			tokens = append(tokens, token)
		}
	}
	return tokens, nil
}

// RemoveUserTokens 移除用户的所有Token
func (s *MemoryStorage) RemoveUserTokens(userID int64) error {
	s.tokensMutex.Lock()
	defer s.tokensMutex.Unlock()
	for token, info := range s.tokens {
		if info.UserID == userID {
			delete(s.tokens, token)
		}
	}
	return nil
}

// GetDeviceToken 获取用户在特定设备上的Token
func (s *MemoryStorage) GetDeviceToken(userID int64, device DeviceType) (string, error) {
	s.tokensMutex.RLock()
	defer s.tokensMutex.RUnlock()
	for token, info := range s.tokens {
		if info.UserID == userID && info.Device == device {
			return token, nil
		}
	}
	return "", ErrTokenNotFound
}

// SaveSecondAuth 保存二级认证信息
func (s *MemoryStorage) SaveSecondAuth(token string, expireAt int64) error {
	s.authMutex.Lock()
	defer s.authMutex.Unlock()
	s.secondAuth[token] = expireAt
	return nil
}

// CheckSecondAuth 检查二级认证
func (s *MemoryStorage) CheckSecondAuth(token string) (bool, error) {
	s.authMutex.RLock()
	defer s.authMutex.RUnlock()
	expireAt, ok := s.secondAuth[token]
	if !ok {
		return false, nil
	}
	now := time.Now().Unix()
	if now > expireAt {
		return false, nil
	}
	return true, nil
}

// RemoveSecondAuth 移除二级认证
func (s *MemoryStorage) RemoveSecondAuth(token string) error {
	s.authMutex.Lock()
	defer s.authMutex.Unlock()
	delete(s.secondAuth, token)
	return nil
}

// GetOnlineUsers 获取在线用户
func (s *MemoryStorage) GetOnlineUsers() ([]UserSession, error) {
	s.tokensMutex.RLock()
	defer s.tokensMutex.RUnlock()
	var sessions []UserSession
	for token, info := range s.tokens {
		sessions = append(sessions, UserSession{
			UserID:    info.UserID,
			Username:  info.Username,
			Device:    info.Device,
			IP:        info.IP,
			UserAgent: info.UserAgent,
			LoginTime: info.LoginTime,
			ExpiresAt: info.ExpiresAt,
			Token:     token,
		})
	}
	return sessions, nil
}

// Close 关闭存储
func (s *MemoryStorage) Close() error {
	return nil
}

// RedisStorage Redis存储实现
type RedisStorage struct {
	client        *redis.Client
	keyPrefix     string
	tokenExpire   time.Duration
	refreshExpire time.Duration
	ownClient     bool // 是否自己创建的客户端
}

// NewRedisStorage 创建Redis存储（通过连接参数）
func NewRedisStorage(addr, password string, db int, keyPrefix string, tokenExpire, refreshExpire time.Duration) (*RedisStorage, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: password,
		DB:       db,
	})

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{
		client:        client,
		keyPrefix:     keyPrefix,
		tokenExpire:   tokenExpire,
		refreshExpire: refreshExpire,
		ownClient:     true,
	}, nil
}

// NewRedisStorageWithClient 创建Redis存储（通过已有客户端）
func NewRedisStorageWithClient(client *redis.Client, keyPrefix string, tokenExpire, refreshExpire time.Duration) (*RedisStorage, error) {
	if client == nil {
		return nil, errs.New(2050, "Redis客户端不能为空")
	}

	// 测试连接
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := client.Ping(ctx).Err(); err != nil {
		return nil, err
	}

	return &RedisStorage{
		client:        client,
		keyPrefix:     keyPrefix,
		tokenExpire:   tokenExpire,
		refreshExpire: refreshExpire,
		ownClient:     false,
	}, nil
}

// tokenKey 生成Token键
func (s *RedisStorage) tokenKey(token string) string {
	return s.keyPrefix + "token:" + token
}

// userTokensKey 生成用户Token键
func (s *RedisStorage) userTokensKey(userID int64) string {
	return fmt.Sprintf("%suser:%d:tokens", s.keyPrefix, userID)
}

// deviceTokenKey 生成设备Token键
func (s *RedisStorage) deviceTokenKey(userID int64, device DeviceType) string {
	return fmt.Sprintf("%suser:%d:device:%s", s.keyPrefix, userID, device)
}

// secondAuthKey 生成二级认证键
func (s *RedisStorage) secondAuthKey(token string) string {
	return s.keyPrefix + "second_auth:" + token
}

// SaveToken 保存Token信息
func (s *RedisStorage) SaveToken(token string, info TokenInfo) error {
	ctx := context.Background()

	// 序列化Token信息
	data, err := json.Marshal(info)
	if err != nil {
		return err
	}

	// 保存Token信息
	tokenKey := s.tokenKey(token)
	if err := s.client.Set(ctx, tokenKey, data, s.tokenExpire).Err(); err != nil {
		return err
	}

	// 保存用户Token关联
	userTokensKey := s.userTokensKey(info.UserID)
	if err := s.client.SAdd(ctx, userTokensKey, token).Err(); err != nil {
		return err
	}
	s.client.Expire(ctx, userTokensKey, s.refreshExpire)

	// 保存设备Token关联
	deviceTokenKey := s.deviceTokenKey(info.UserID, info.Device)
	if err := s.client.Set(ctx, deviceTokenKey, token, s.tokenExpire).Err(); err != nil {
		return err
	}

	return nil
}

// GetToken 获取Token信息
func (s *RedisStorage) GetToken(token string) (TokenInfo, error) {
	ctx := context.Background()

	// 获取Token信息
	tokenKey := s.tokenKey(token)
	data, err := s.client.Get(ctx, tokenKey).Bytes()
	if err != nil {
		if err == redis.Nil {
			return TokenInfo{}, ErrTokenNotFound
		}
		return TokenInfo{}, err
	}

	// 反序列化Token信息
	var info TokenInfo
	if err := json.Unmarshal(data, &info); err != nil {
		return TokenInfo{}, err
	}

	return info, nil
}

// RemoveToken 移除Token
func (s *RedisStorage) RemoveToken(token string) error {
	ctx := context.Background()

	// 获取Token信息
	info, err := s.GetToken(token)
	if err != nil {
		return err
	}

	// 移除Token信息
	tokenKey := s.tokenKey(token)
	if err := s.client.Del(ctx, tokenKey).Err(); err != nil {
		return err
	}

	// 移除用户Token关联
	userTokensKey := s.userTokensKey(info.UserID)
	if err := s.client.SRem(ctx, userTokensKey, token).Err(); err != nil {
		return err
	}

	// 移除设备Token关联
	deviceTokenKey := s.deviceTokenKey(info.UserID, info.Device)
	if err := s.client.Del(ctx, deviceTokenKey).Err(); err != nil {
		return err
	}

	// 移除二级认证
	secondAuthKey := s.secondAuthKey(token)
	s.client.Del(ctx, secondAuthKey)

	return nil
}

// GetUserTokens 获取用户的所有Token
func (s *RedisStorage) GetUserTokens(userID int64) ([]string, error) {
	ctx := context.Background()

	// 获取用户Token关联
	userTokensKey := s.userTokensKey(userID)
	tokens, err := s.client.SMembers(ctx, userTokensKey).Result()
	if err != nil {
		if err == redis.Nil {
			return []string{}, nil
		}
		return nil, err
	}

	return tokens, nil
}

// RemoveUserTokens 移除用户的所有Token
func (s *RedisStorage) RemoveUserTokens(userID int64) error {
	ctx := context.Background()

	// 获取用户Token关联
	userTokensKey := s.userTokensKey(userID)
	tokens, err := s.client.SMembers(ctx, userTokensKey).Result()
	if err != nil {
		if err == redis.Nil {
			return nil
		}
		return err
	}

	// 移除所有Token
	for _, token := range tokens {
		if err := s.RemoveToken(token); err != nil {
			return err
		}
	}

	// 移除用户Token关联
	if err := s.client.Del(ctx, userTokensKey).Err(); err != nil {
		return err
	}

	return nil
}

// GetDeviceToken 获取用户在特定设备上的Token
func (s *RedisStorage) GetDeviceToken(userID int64, device DeviceType) (string, error) {
	ctx := context.Background()

	// 获取设备Token关联
	deviceTokenKey := s.deviceTokenKey(userID, device)
	token, err := s.client.Get(ctx, deviceTokenKey).Result()
	if err != nil {
		if err == redis.Nil {
			return "", ErrTokenNotFound
		}
		return "", err
	}

	return token, nil
}

// SaveSecondAuth 保存二级认证信息
func (s *RedisStorage) SaveSecondAuth(token string, expireAt int64) error {
	ctx := context.Background()

	// 保存二级认证信息
	secondAuthKey := s.secondAuthKey(token)
	expireDuration := time.Until(time.Unix(expireAt, 0))
	if expireDuration <= 0 {
		return ErrInvalidParams
	}

	return s.client.Set(ctx, secondAuthKey, "1", expireDuration).Err()
}

// CheckSecondAuth 检查二级认证
func (s *RedisStorage) CheckSecondAuth(token string) (bool, error) {
	ctx := context.Background()

	// 检查二级认证信息
	secondAuthKey := s.secondAuthKey(token)
	exists, err := s.client.Exists(ctx, secondAuthKey).Result()
	if err != nil {
		return false, err
	}

	return exists > 0, nil
}

// RemoveSecondAuth 移除二级认证
func (s *RedisStorage) RemoveSecondAuth(token string) error {
	ctx := context.Background()

	// 移除二级认证信息
	secondAuthKey := s.secondAuthKey(token)
	return s.client.Del(ctx, secondAuthKey).Err()
}

// GetOnlineUsers 获取在线用户
func (s *RedisStorage) GetOnlineUsers() ([]UserSession, error) {
	ctx := context.Background()

	// 获取所有Token键
	pattern := s.tokenKey("*")
	var cursor uint64
	var sessions []UserSession

	for {
		var keys []string
		var err error
		keys, cursor, err = s.client.Scan(ctx, cursor, pattern, 100).Result()
		if err != nil {
			return nil, err
		}

		// 获取Token信息
		for _, key := range keys {
			token := key[len(s.tokenKey("")):]
			info, err := s.GetToken(token)
			if err != nil {
				continue
			}

			sessions = append(sessions, UserSession{
				UserID:    info.UserID,
				Username:  info.Username,
				Device:    info.Device,
				IP:        info.IP,
				UserAgent: info.UserAgent,
				LoginTime: info.LoginTime,
				ExpiresAt: info.ExpiresAt,
				Token:     token,
			})
		}

		if cursor == 0 {
			break
		}
	}

	return sessions, nil
}

// Close 关闭存储
func (s *RedisStorage) Close() error {
	// 只关闭自己创建的客户端
	if s.ownClient {
		return s.client.Close()
	}
	return nil
}
