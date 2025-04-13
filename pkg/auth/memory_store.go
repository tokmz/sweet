package auth

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"
)

// MemoryStore 内存存储实现
type MemoryStore struct {
	// tokenMap 存储token到token信息的映射
	tokenMap map[string]TokenInfo
	// userTokenMap 存储用户ID到该用户所有token的映射
	userTokenMap map[int64]map[string]string
	// 用户设备token map 用户ID:设备名->token列表
	userDeviceTokenMap map[string]map[string]string
	mutex              sync.RWMutex
}

// NewMemoryStore 创建内存存储
func NewMemoryStore() *MemoryStore {
	return &MemoryStore{
		tokenMap:           make(map[string]TokenInfo),
		userTokenMap:       make(map[int64]map[string]string),
		userDeviceTokenMap: make(map[string]map[string]string),
	}
}

// StoreToken 存储令牌
func (m *MemoryStore) StoreToken(ctx context.Context, token TokenInfo) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 存储token信息
	m.tokenMap[token.Token] = token

	// 存储用户token映射
	if _, ok := m.userTokenMap[token.UserID]; !ok {
		m.userTokenMap[token.UserID] = make(map[string]string)
	}
	m.userTokenMap[token.UserID][token.Token] = token.Token

	// 存储用户设备token映射
	key := fmt.Sprintf("%d:%s", token.UserID, token.Device)
	if _, ok := m.userDeviceTokenMap[key]; !ok {
		m.userDeviceTokenMap[key] = make(map[string]string)
	}
	m.userDeviceTokenMap[key][token.Token] = token.Token

	return nil
}

// GetToken 获取令牌信息
func (m *MemoryStore) GetToken(ctx context.Context, token string) (*TokenInfo, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	info, ok := m.tokenMap[token]
	if !ok {
		return nil, ErrTokenNotFound
	}

	// 检查令牌是否过期
	if time.Now().After(info.ExpireAt) {
		return nil, ErrTokenExpired
	}

	return &info, nil
}

// GetUserTokens 获取用户的所有令牌
func (m *MemoryStore) GetUserTokens(ctx context.Context, userID int64) ([]TokenInfo, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	tokenMap, ok := m.userTokenMap[userID]
	if !ok {
		return []TokenInfo{}, nil
	}

	tokens := make([]TokenInfo, 0, len(tokenMap))
	for token := range tokenMap {
		info, ok := m.tokenMap[token]
		if ok && time.Now().Before(info.ExpireAt) {
			tokens = append(tokens, info)
		}
	}

	return tokens, nil
}

// GetUserTokensByDevice 获取用户指定设备的所有令牌
func (m *MemoryStore) GetUserTokensByDevice(ctx context.Context, userID int64, device string) ([]TokenInfo, error) {
	m.mutex.RLock()
	defer m.mutex.RUnlock()

	key := fmt.Sprintf("%d:%s", userID, device)
	tokenMap, ok := m.userDeviceTokenMap[key]
	if !ok {
		return []TokenInfo{}, nil
	}

	tokens := make([]TokenInfo, 0, len(tokenMap))
	for token := range tokenMap {
		info, ok := m.tokenMap[token]
		if ok && time.Now().Before(info.ExpireAt) {
			tokens = append(tokens, info)
		}
	}

	return tokens, nil
}

// RemoveToken 移除令牌
func (m *MemoryStore) RemoveToken(ctx context.Context, token string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	info, ok := m.tokenMap[token]
	if !ok {
		return nil // 令牌不存在，无需删除
	}

	// 从tokenMap中删除
	delete(m.tokenMap, token)

	// 从userTokenMap中删除
	if userTokens, ok := m.userTokenMap[info.UserID]; ok {
		delete(userTokens, token)
	}

	// 从userDeviceTokenMap中删除
	key := fmt.Sprintf("%d:%s", info.UserID, info.Device)
	if deviceTokens, ok := m.userDeviceTokenMap[key]; ok {
		delete(deviceTokens, token)
	}

	return nil
}

// RemoveUserTokens 移除用户的所有令牌
func (m *MemoryStore) RemoveUserTokens(ctx context.Context, userID int64) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 获取用户所有token
	tokenMap, ok := m.userTokenMap[userID]
	if !ok {
		return nil // 用户无token，无需删除
	}

	// 从tokenMap中删除所有token
	for token := range tokenMap {
		delete(m.tokenMap, token)
	}

	// 删除用户token映射
	delete(m.userTokenMap, userID)

	// 删除用户设备token映射
	for key := range m.userDeviceTokenMap {
		if fmt.Sprintf("%d:", userID) == key[:len(strconv.FormatInt(userID, 10))+1] {
			delete(m.userDeviceTokenMap, key)
		}
	}

	return nil
}

// RemoveUserTokensByDevice 移除用户指定设备的所有令牌
func (m *MemoryStore) RemoveUserTokensByDevice(ctx context.Context, userID int64, device string) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	// 获取用户设备token
	key := fmt.Sprintf("%d:%s", userID, device)
	deviceTokens, ok := m.userDeviceTokenMap[key]
	if !ok {
		return nil // 设备无token，无需删除
	}

	// 从tokenMap和userTokenMap中删除相关token
	for token := range deviceTokens {
		delete(m.tokenMap, token)
		if userTokens, ok := m.userTokenMap[userID]; ok {
			delete(userTokens, token)
		}
	}

	// 删除设备token映射
	delete(m.userDeviceTokenMap, key)

	return nil
}

// UpdateActiveTime 更新令牌活跃时间
func (m *MemoryStore) UpdateActiveTime(ctx context.Context, token string, t time.Time) error {
	m.mutex.Lock()
	defer m.mutex.Unlock()

	info, ok := m.tokenMap[token]
	if !ok {
		return ErrTokenNotFound
	}

	info.LastActiveTime = t
	m.tokenMap[token] = info

	return nil
}
