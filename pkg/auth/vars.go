package auth

import "sync"

var (
	// 全局认证服务实例
	authService AuthService
	// 全局权限管理器实例
	permissionManager PermissionManager
	// 全局二级认证验证器实例
	secondaryAuthVerifier SecondaryAuthVerifier
	// 互斥锁，保护全局变量
	authMutex sync.RWMutex
)

// SetAuthService 设置全局认证服务实例
func SetAuthService(service AuthService) {
	authMutex.Lock()
	defer authMutex.Unlock()
	authService = service
}

// GetAuthService 获取全局认证服务实例
func GetAuthService() (AuthService, error) {
	authMutex.RLock()
	defer authMutex.RUnlock()
	if authService == nil {
		return nil, ErrServiceNotInitialized
	}
	return authService, nil
}

// MustGetAuthService 获取全局认证服务实例，如果未初始化则panic
func MustGetAuthService() AuthService {
	service, err := GetAuthService()
	if err != nil {
		panic(err)
	}
	return service
}

// SetPermissionManager 设置全局权限管理器实例
func SetPermissionManager(manager PermissionManager) {
	authMutex.Lock()
	defer authMutex.Unlock()
	permissionManager = manager
}

// GetPermissionManager 获取全局权限管理器实例
func GetPermissionManager() (PermissionManager, error) {
	authMutex.RLock()
	defer authMutex.RUnlock()
	if permissionManager == nil {
		return nil, ErrServiceNotInitialized
	}
	return permissionManager, nil
}

// MustGetPermissionManager 获取全局权限管理器实例，如果未初始化则panic
func MustGetPermissionManager() PermissionManager {
	manager, err := GetPermissionManager()
	if err != nil {
		panic(err)
	}
	return manager
}

// SetSecondaryAuthVerifier 设置全局二级认证验证器实例
func SetSecondaryAuthVerifier(verifier SecondaryAuthVerifier) {
	authMutex.Lock()
	defer authMutex.Unlock()
	secondaryAuthVerifier = verifier
}

// GetSecondaryAuthVerifier 获取全局二级认证验证器实例
func GetSecondaryAuthVerifier() (SecondaryAuthVerifier, error) {
	authMutex.RLock()
	defer authMutex.RUnlock()
	if secondaryAuthVerifier == nil {
		return nil, ErrServiceNotInitialized
	}
	return secondaryAuthVerifier, nil
}

// MustGetSecondaryAuthVerifier 获取全局二级认证验证器实例，如果未初始化则panic
func MustGetSecondaryAuthVerifier() SecondaryAuthVerifier {
	verifier, err := GetSecondaryAuthVerifier()
	if err != nil {
		panic(err)
	}
	return verifier
}

// 设备类型常量
const (
	// DeviceTypeWeb Web端
	DeviceTypeWeb = "web"
	// DeviceTypeMobile 移动端
	DeviceTypeMobile = "mobile"
	// DeviceTypeTablet 平板端
	DeviceTypeTablet = "tablet"
	// DeviceTypeDesktop 桌面端
	DeviceTypeDesktop = "desktop"
	// DeviceTypeAPI API调用
	DeviceTypeAPI = "api"
)

// 令牌类型常量
const (
	// TokenTypeAccess 访问令牌
	TokenTypeAccess = "access"
	// TokenTypeRefresh 刷新令牌
	TokenTypeRefresh = "refresh"
)
