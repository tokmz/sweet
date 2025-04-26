package auth

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

// RedisSessionStore Redis会话存储实现
type RedisSessionStore struct {
	// client Redis客户端
	client *redis.Client
	// keyPrefix 键前缀
	keyPrefix string
}

// NewRedisSessionStore 创建Redis会话存储
func NewRedisSessionStore(client *redis.Client) *RedisSessionStore {
	return &RedisSessionStore{
		client:    client,
		keyPrefix: "sweet:auth:session:",
	}
}

// userSessionKey 生成用户会话键
func (s *RedisSessionStore) userSessionKey(userID uint64, deviceType string) string {
	return fmt.Sprintf("%suser:%d:%s", s.keyPrefix, userID, deviceType)
}

// deviceTypeKey 生成设备类型键
func (s *RedisSessionStore) deviceTypeKey(deviceType string) string {
	return fmt.Sprintf("%sdevice:%s", s.keyPrefix, deviceType)
}

// SaveSession 保存会话
func (s *RedisSessionStore) SaveSession(session *SessionInfo) error {
	// 序列化会话
	data, err := json.Marshal(session)
	if err != nil {
		return err
	}

	ctx := context.Background()
	// 计算过期时间
	expiration := time.Until(time.Unix(session.ExpireTime, 0))

	// 保存用户会话
	key := s.userSessionKey(session.UserID, session.DeviceType)
	err = s.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		return err
	}

	// 将用户ID添加到设备类型集合
	deviceKey := s.deviceTypeKey(session.DeviceType)
	err = s.client.SAdd(ctx, deviceKey, session.UserID).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetSession 获取会话
func (s *RedisSessionStore) GetSession(userID uint64, deviceType string) (*SessionInfo, error) {
	ctx := context.Background()
	// 获取用户会话
	key := s.userSessionKey(userID, deviceType)
	data, err := s.client.Get(ctx, key).Bytes()
	if err != nil {
		if err == redis.Nil {
			return nil, ErrSessionNotFound
		}
		return nil, err
	}

	// 反序列化会话
	var session SessionInfo
	err = json.Unmarshal(data, &session)
	if err != nil {
		return nil, err
	}

	return &session, nil
}

// RemoveSession 删除会话
func (s *RedisSessionStore) RemoveSession(userID uint64, deviceType string) error {
	ctx := context.Background()
	// 删除用户会话
	key := s.userSessionKey(userID, deviceType)
	err := s.client.Del(ctx, key).Err()
	if err != nil {
		return err
	}

	// 从设备类型集合中移除用户ID
	deviceKey := s.deviceTypeKey(deviceType)
	err = s.client.SRem(ctx, deviceKey, userID).Err()
	if err != nil {
		return err
	}

	return nil
}

// GetUserSessions 获取用户所有会话
func (s *RedisSessionStore) GetUserSessions(userID uint64) ([]*SessionInfo, error) {
	ctx := context.Background()
	// 获取所有设备类型
	deviceTypes := []string{DeviceTypeWeb, DeviceTypeMobile, DeviceTypeTablet, DeviceTypeDesktop, DeviceTypeAPI}

	// 获取所有会话
	sessions := make([]*SessionInfo, 0)
	for _, deviceType := range deviceTypes {
		key := s.userSessionKey(userID, deviceType)
		data, err := s.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, err
		}

		// 反序列化会话
		var session SessionInfo
		err = json.Unmarshal(data, &session)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// GetDeviceTypeSessions 获取指定设备类型的所有会话
func (s *RedisSessionStore) GetDeviceTypeSessions(deviceType string) ([]*SessionInfo, error) {
	ctx := context.Background()
	// 获取设备类型集合中的所有用户ID
	deviceKey := s.deviceTypeKey(deviceType)
	userIDs, err := s.client.SMembers(ctx, deviceKey).Result()
	if err != nil {
		return nil, err
	}

	// 获取所有会话
	sessions := make([]*SessionInfo, 0, len(userIDs))
	for _, userIDStr := range userIDs {
		var userID uint64
		_, err := fmt.Sscanf(userIDStr, "%d", &userID)
		if err != nil {
			continue
		}

		key := s.userSessionKey(userID, deviceType)
		data, err := s.client.Get(ctx, key).Bytes()
		if err != nil {
			if err == redis.Nil {
				continue
			}
			return nil, err
		}

		// 反序列化会话
		var session SessionInfo
		err = json.Unmarshal(data, &session)
		if err != nil {
			return nil, err
		}

		sessions = append(sessions, &session)
	}

	return sessions, nil
}

// CleanExpiredSessions 清理过期会话
func (s *RedisSessionStore) CleanExpiredSessions() error {
	// Redis会自动清理过期的键，无需手动实现
	return nil
}
