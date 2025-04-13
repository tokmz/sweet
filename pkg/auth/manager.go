package auth

import (
	"context"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/casbin/casbin/v2"
)

// Manager 认证管理器
type Manager struct {
	config         Config
	jwtService     *JWTService
	storage        Storage
	casbinService  *CasbinService
	tokenGenerator CustomTokenGenerator
	listeners      []AuthEventListener
	listenersMutex sync.RWMutex
}

// ManagerOption 管理器选项
type ManagerOption func(*Manager) error

// WithRedisClient 使用已有的Redis客户端
func WithRedisClient(client *redis.Client, keyPrefix string) ManagerOption {
	return func(m *Manager) error {
		storage, err := NewRedisStorageWithClient(client, keyPrefix, m.config.TokenExpire, m.config.RefreshExpire)
		if err != nil {
			return fmt.Errorf("创建Redis存储失败: %w", err)
		}
		m.storage = storage
		m.config.EnableRedis = true
		return nil
	}
}

// WithRedisStorage 使用Redis存储（通过连接参数）
func WithRedisStorage(addr, password string, db int, keyPrefix string) ManagerOption {
	return func(m *Manager) error {
		storage, err := NewRedisStorage(addr, password, db, keyPrefix, m.config.TokenExpire, m.config.RefreshExpire)
		if err != nil {
			return fmt.Errorf("创建Redis存储失败: %w", err)
		}
		m.storage = storage
		m.config.EnableRedis = true
		return nil
	}
}

// WithMemoryStorage 使用内存存储
func WithMemoryStorage() ManagerOption {
	return func(m *Manager) error {
		m.storage = NewMemoryStorage()
		m.config.EnableRedis = false
		return nil
	}
}

// NewManager 创建认证管理器
func NewManager(config Config, options ...ManagerOption) (*Manager, error) {
	// 验证配置
	if config.JWTSecret == "" {
		return nil, ErrInvalidConfig
	}
	if config.TokenExpire <= 0 {
		config.TokenExpire = 24 * time.Hour
	}
	if config.RefreshExpire <= 0 {
		config.RefreshExpire = 7 * 24 * time.Hour
	}
	if config.RenewThreshold <= 0 || config.RenewThreshold >= 1 {
		config.RenewThreshold = 0.5
	}
	if config.MaxConcurrency <= 0 {
		config.MaxConcurrency = 100
	}

	// 创建JWT服务
	jwtService := NewJWTService(config.JWTSecret, config.TokenExpire, config.RefreshExpire)

	// 创建管理器
	manager := &Manager{
		config:     config,
		jwtService: jwtService,
		listeners:  make([]AuthEventListener, 0),
	}

	// 默认使用内存存储
	manager.storage = NewMemoryStorage()

	// 如果配置了Redis但没有提供选项，使用配置创建Redis存储
	if config.EnableRedis && len(options) == 0 {
		storage, err := NewRedisStorage(
			config.RedisAddr,
			config.RedisPass,
			config.RedisDB,
			config.RedisKeyPrefix,
			config.TokenExpire,
			config.RefreshExpire,
		)
		if err != nil {
			return nil, fmt.Errorf("创建Redis存储失败: %w", err)
		}
		manager.storage = storage
	}

	// 应用选项
	for _, option := range options {
		if err := option(manager); err != nil {
			// 关闭已创建的存储
			if manager.storage != nil {
				manager.storage.Close()
			}
			return nil, err
		}
	}

	return manager, nil
}

// InitCasbin 初始化Casbin
func (m *Manager) InitCasbin(config CasbinConfig) (*casbin.Enforcer, error) {
	casbinService, err := NewCasbinService(config)
	if err != nil {
		return nil, err
	}
	m.casbinService = casbinService
	return casbinService.GetEnforcer(), nil
}

// SetTokenGenerator 设置自定义Token生成器
func (m *Manager) SetTokenGenerator(generator CustomTokenGenerator) {
	m.tokenGenerator = generator
}

// AddListener 添加事件监听器
func (m *Manager) AddListener(listener AuthEventListener) {
	m.listenersMutex.Lock()
	defer m.listenersMutex.Unlock()
	m.listeners = append(m.listeners, listener)
}

// RemoveListener 移除事件监听器
func (m *Manager) RemoveListener(listener AuthEventListener) {
	m.listenersMutex.Lock()
	defer m.listenersMutex.Unlock()
	for i, l := range m.listeners {
		if fmt.Sprintf("%p", l) == fmt.Sprintf("%p", listener) {
			m.listeners = append(m.listeners[:i], m.listeners[i+1:]...)
			break
		}
	}
}

// notifyListeners 通知事件监听器
func (m *Manager) notifyListeners(event AuthEvent) {
	m.listenersMutex.RLock()
	listeners := make([]AuthEventListener, len(m.listeners))
	copy(listeners, m.listeners)
	m.listenersMutex.RUnlock()

	for _, listener := range listeners {
		go listener(event)
	}
}

// Login 登录
func (m *Manager) Login(params LoginParams) (string, string, error) {
	// 验证参数
	if params.UserID <= 0 {
		return "", "", ErrInvalidParams
	}

	// 检查登录模式
	if m.config.LoginModel == LoginModelSingle {
		// 单端登录模式，踢出用户所有设备
		if err := m.storage.RemoveUserTokens(params.UserID); err != nil {
			return "", "", fmt.Errorf("移除用户Token失败: %w", err)
		}
	} else if m.config.LoginModel == LoginModelExclusive {
		// 同端互斥登录模式，踢出用户同类设备
		token, err := m.storage.GetDeviceToken(params.UserID, params.Device)
		if err == nil && token != "" {
			if err := m.storage.RemoveToken(token); err != nil {
				return "", "", fmt.Errorf("移除设备Token失败: %w", err)
			}

			// 通知事件
			m.notifyListeners(AuthEvent{
				Type:      EventKickout,
				UserID:    params.UserID,
				Username:  params.Username,
				Device:    params.Device,
				IP:        params.IP,
				Token:     token,
				Timestamp: time.Now().Unix(),
			})
		}
	}

	// 生成Token
	var token string
	var err error
	if m.config.TokenStyle == TokenStyleCustom && m.tokenGenerator != nil {
		// 使用自定义生成器
		token, err = m.tokenGenerator(params.UserID, params.Device, params.ExtraData)
	} else {
		// 使用内置生成器
		token, err = m.jwtService.GenerateTokenByStyle(m.config.TokenStyle, params.UserID, params.Username, string(params.Device), params.ExtraData)
	}
	if err != nil {
		return "", "", fmt.Errorf("生成Token失败: %w", err)
	}

	// 生成刷新Token
	refreshToken, err := m.jwtService.GenerateRefreshToken(params.UserID, params.Username, string(params.Device))
	if err != nil {
		return "", "", fmt.Errorf("生成刷新Token失败: %w", err)
	}

	// 计算过期时间
	var expiresAt int64
	if params.RememberMe {
		expiresAt = time.Now().Add(m.config.RefreshExpire).Unix()
	} else {
		expiresAt = time.Now().Add(m.config.TokenExpire).Unix()
	}

	// 保存Token信息
	tokenInfo := TokenInfo{
		Token:        token,
		RefreshToken: refreshToken,
		ExpiresAt:    expiresAt,
		UserID:       params.UserID,
		Username:     params.Username,
		Device:       params.Device,
		IP:           params.IP,
		UserAgent:    params.UserAgent,
		LoginTime:    time.Now().Unix(),
		ExtraData:    params.ExtraData,
	}
	if err := m.storage.SaveToken(token, tokenInfo); err != nil {
		return "", "", fmt.Errorf("保存Token失败: %w", err)
	}

	// 通知事件
	m.notifyListeners(AuthEvent{
		Type:      EventLogin,
		UserID:    params.UserID,
		Username:  params.Username,
		Device:    params.Device,
		IP:        params.IP,
		Token:     token,
		Timestamp: time.Now().Unix(),
	})

	return token, refreshToken, nil
}

// Logout 注销
func (m *Manager) Logout(token string) error {
	// 获取Token信息
	info, err := m.storage.GetToken(token)
	if err != nil {
		return err
	}

	// 移除Token
	if err := m.storage.RemoveToken(token); err != nil {
		return err
	}

	// 通知事件
	m.notifyListeners(AuthEvent{
		Type:      EventLogout,
		UserID:    info.UserID,
		Username:  info.Username,
		Device:    info.Device,
		IP:        info.IP,
		Token:     token,
		Timestamp: time.Now().Unix(),
	})

	return nil
}

// ValidateToken 验证Token
func (m *Manager) ValidateToken(token string) (*TokenInfo, error) {
	// 从存储中获取Token信息
	info, err := m.storage.GetToken(token)
	if err != nil {
		return nil, err
	}

	// 检查Token是否过期
	now := time.Now().Unix()
	if now > info.ExpiresAt {
		// 移除过期Token
		m.storage.RemoveToken(token)

		// 通知事件
		m.notifyListeners(AuthEvent{
			Type:      EventTokenExpired,
			UserID:    info.UserID,
			Username:  info.Username,
			Device:    info.Device,
			Token:     token,
			Timestamp: now,
		})

		return nil, ErrTokenExpired
	}

	// 如果启用了自动续签，检查是否需要续签
	if m.config.AutoRenew {
		tokenLifetime := info.ExpiresAt - info.LoginTime
		elapsedTime := now - info.LoginTime
		if float64(elapsedTime)/float64(tokenLifetime) > m.config.RenewThreshold {
			// 需要续签
			newExpiresAt := now + int64(m.config.TokenExpire.Seconds())
			info.ExpiresAt = newExpiresAt
			m.storage.SaveToken(token, info)

			// 通知事件
			m.notifyListeners(AuthEvent{
				Type:      EventTokenRefresh,
				UserID:    info.UserID,
				Username:  info.Username,
				Device:    info.Device,
				Token:     token,
				Timestamp: now,
			})
		}
	}

	return &info, nil
}

// RefreshToken 刷新Token
func (m *Manager) RefreshToken(refreshToken string) (string, string, error) {
	// 解析刷新Token
	claims, err := m.jwtService.ParseRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}

	// 获取用户信息
	username := claims.Subject
	userID, err := strconv.ParseInt(claims.ID, 10, 64)
	if err != nil {
		return "", "", ErrRefreshTokenInvalid
	}

	// 生成新Token
	newToken, err := m.jwtService.GenerateTokenByStyle(m.config.TokenStyle, userID, username, "", nil)
	if err != nil {
		return "", "", fmt.Errorf("生成Token失败: %w", err)
	}

	// 生成新刷新Token
	newRefreshToken, err := m.jwtService.GenerateRefreshToken(userID, username, "")
	if err != nil {
		return "", "", fmt.Errorf("生成刷新Token失败: %w", err)
	}

	// 获取用户Token
	tokens, err := m.storage.GetUserTokens(userID)
	if err != nil {
		return "", "", err
	}

	// 查找匹配的Token
	var oldToken string
	var tokenInfo TokenInfo
	for _, token := range tokens {
		info, err := m.storage.GetToken(token)
		if err != nil {
			continue
		}
		if info.RefreshToken == refreshToken {
			oldToken = token
			tokenInfo = info
			break
		}
	}

	if oldToken == "" {
		return "", "", ErrRefreshTokenInvalid
	}

	// 移除旧Token
	if err := m.storage.RemoveToken(oldToken); err != nil {
		return "", "", err
	}

	// 更新Token信息
	tokenInfo.Token = newToken
	tokenInfo.RefreshToken = newRefreshToken
	tokenInfo.ExpiresAt = time.Now().Add(m.config.TokenExpire).Unix()

	// 保存新Token
	if err := m.storage.SaveToken(newToken, tokenInfo); err != nil {
		return "", "", fmt.Errorf("保存Token失败: %w", err)
	}

	// 通知事件
	m.notifyListeners(AuthEvent{
		Type:      EventTokenRefresh,
		UserID:    tokenInfo.UserID,
		Username:  tokenInfo.Username,
		Device:    tokenInfo.Device,
		Token:     newToken,
		Timestamp: time.Now().Unix(),
	})

	return newToken, newRefreshToken, nil
}

// KickoutByUserID 根据用户ID踢人下线
func (m *Manager) KickoutByUserID(userID int64) error {
	// 获取用户Token
	tokens, err := m.storage.GetUserTokens(userID)
	if err != nil {
		return err
	}

	// 获取用户信息
	var username string
	if len(tokens) > 0 {
		info, err := m.storage.GetToken(tokens[0])
		if err == nil {
			username = info.Username
		}
	}

	// 移除所有Token
	if err := m.storage.RemoveUserTokens(userID); err != nil {
		return err
	}

	// 通知事件
	m.notifyListeners(AuthEvent{
		Type:      EventKickout,
		UserID:    userID,
		Username:  username,
		Timestamp: time.Now().Unix(),
		ExtraData: tokens,
	})

	return nil
}

// KickoutByToken 根据Token踢人下线
func (m *Manager) KickoutByToken(token string) error {
	// 获取Token信息
	info, err := m.storage.GetToken(token)
	if err != nil {
		return err
	}

	// 移除Token
	if err := m.storage.RemoveToken(token); err != nil {
		return err
	}

	// 通知事件
	m.notifyListeners(AuthEvent{
		Type:      EventKickout,
		UserID:    info.UserID,
		Username:  info.Username,
		Device:    info.Device,
		IP:        info.IP,
		Token:     token,
		Timestamp: time.Now().Unix(),
	})

	return nil
}

// EnableSecondAuth 启用二级认证
func (m *Manager) EnableSecondAuth(token string, duration time.Duration) error {
	// 验证Token
	info, err := m.ValidateToken(token)
	if err != nil {
		return err
	}

	// 计算过期时间
	expireAt := time.Now().Add(duration).Unix()

	// 保存二级认证信息
	if err := m.storage.SaveSecondAuth(token, expireAt); err != nil {
		return err
	}

	// 通知事件
	m.notifyListeners(AuthEvent{
		Type:      EventSecondAuthEnabled,
		UserID:    info.UserID,
		Username:  info.Username,
		Device:    info.Device,
		Token:     token,
		Timestamp: time.Now().Unix(),
	})

	return nil
}

// CheckSecondAuth 检查二级认证
func (m *Manager) CheckSecondAuth(token string) (bool, error) {
	// 验证Token
	_, err := m.ValidateToken(token)
	if err != nil {
		return false, err
	}

	// 检查二级认证
	return m.storage.CheckSecondAuth(token)
}

// DisableSecondAuth 禁用二级认证
func (m *Manager) DisableSecondAuth(token string) error {
	// 验证Token
	info, err := m.ValidateToken(token)
	if err != nil {
		return err
	}

	// 移除二级认证
	if err := m.storage.RemoveSecondAuth(token); err != nil {
		return err
	}

	// 通知事件
	m.notifyListeners(AuthEvent{
		Type:      EventSecondAuthDisabled,
		UserID:    info.UserID,
		Username:  info.Username,
		Device:    info.Device,
		Token:     token,
		Timestamp: time.Now().Unix(),
	})

	return nil
}

// GetOnlineUsers 获取在线用户
func (m *Manager) GetOnlineUsers() ([]UserSession, error) {
	return m.storage.GetOnlineUsers()
}

// CreateTempToken 创建临时Token
func (m *Manager) CreateTempToken(userID int64, duration time.Duration, extraData map[string]interface{}) (string, error) {
	// 生成临时Token
	token, err := m.jwtService.GenerateToken(userID, "", "temp", extraData)
	if err != nil {
		return "", fmt.Errorf("生成临时Token失败: %w", err)
	}

	// 保存临时Token
	tokenInfo := TokenInfo{
		Token:     token,
		ExpiresAt: time.Now().Add(duration).Unix(),
		UserID:    userID,
		Device:    DeviceTypeOther,
		LoginTime: time.Now().Unix(),
		ExtraData: extraData,
	}
	if err := m.storage.SaveToken(token, tokenInfo); err != nil {
		return "", fmt.Errorf("保存临时Token失败: %w", err)
	}

	return token, nil
}

// ValidateTempToken 验证临时Token
func (m *Manager) ValidateTempToken(token string) (*TokenInfo, error) {
	// 验证Token
	info, err := m.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	// 检查是否为临时Token
	if info.Device != DeviceTypeOther {
		return nil, ErrTokenInvalid
	}

	return info, nil
}

// CheckPermission 检查权限
func (m *Manager) CheckPermission(userID int64, obj string, act string) (bool, error) {
	if m.casbinService == nil {
		return false, ErrCasbinError
	}
	return m.casbinService.CheckPermission(strconv.FormatInt(userID, 10), obj, act)
}

// AssignRole 为用户分配角色
func (m *Manager) AssignRole(userID int64, role string) error {
	if m.casbinService == nil {
		return ErrCasbinError
	}
	return m.casbinService.AddRoleForUser(strconv.FormatInt(userID, 10), role)
}

// RemoveRole 移除用户角色
func (m *Manager) RemoveRole(userID int64, role string) error {
	if m.casbinService == nil {
		return ErrCasbinError
	}
	return m.casbinService.RemoveRoleForUser(strconv.FormatInt(userID, 10), role)
}

// AddPermissionForRole 为角色添加权限
func (m *Manager) AddPermissionForRole(role string, obj string, act string) error {
	if m.casbinService == nil {
		return ErrCasbinError
	}
	return m.casbinService.AddPermissionForRole(role, obj, act)
}

// GetRolesForUser 获取用户角色
func (m *Manager) GetRolesForUser(userID int64) ([]string, error) {
	if m.casbinService == nil {
		return nil, ErrCasbinError
	}
	return m.casbinService.GetRolesForUser(strconv.FormatInt(userID, 10))
}

// GetPermissionsForUser 获取用户权限
func (m *Manager) GetPermissionsForUser(userID int64) ([][]string, error) {
	if m.casbinService == nil {
		return nil, ErrCasbinError
	}
	return m.casbinService.GetPermissionsForUser(strconv.FormatInt(userID, 10))
}

// HasRole 检查用户是否有角色
func (m *Manager) HasRole(userID int64, role string) (bool, error) {
	if m.casbinService == nil {
		return false, ErrCasbinError
	}
	return m.casbinService.HasRoleForUser(strconv.FormatInt(userID, 10), role)
}

// ConfigSSO 配置单点登录
func (m *Manager) ConfigSSO(config SSOConfig) error {
	// 单点登录功能需要Redis支持
	if !m.config.EnableRedis {
		return ErrSSONotEnabled
	}

	// TODO: 实现单点登录配置
	return ErrNotImplemented
}

// SSOLogin 单点登录
func (m *Manager) SSOLogin(params LoginParams) (string, string, error) {
	// TODO: 实现单点登录
	return "", "", ErrNotImplemented
}

// SSOLogout 单点注销
func (m *Manager) SSOLogout(userID int64) error {
	// TODO: 实现单点注销
	return ErrNotImplemented
}

// Close 关闭管理器
func (m *Manager) Close() error {
	return m.storage.Close()
}
