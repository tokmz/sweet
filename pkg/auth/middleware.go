package auth

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

// JWTAuthMiddleware 创建JWT认证中间件
func JWTAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取认证服务
		authService, err := GetAuthService()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "认证服务未初始化",
			})
			c.Abort()
			return
		}

		// 从请求头获取令牌
		auth := c.GetHeader("Authorization")
		if auth == "" {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未提供认证令牌",
			})
			c.Abort()
			return
		}

		// 解析令牌
		parts := strings.SplitN(auth, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "认证格式无效",
			})
			c.Abort()
			return
		}

		// 验证令牌
		tokenInfo, err := authService.VerifyToken(parts[1])
		if err != nil {
			var msg string
			if err == ErrTokenExpired {
				msg = "认证令牌已过期"
			} else {
				msg = "认证令牌无效"
			}
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  msg,
			})
			c.Abort()
			return
		}

		// 检查令牌类型
		if tokenInfo.TokenType != TokenTypeAccess {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "无效的令牌类型",
			})
			c.Abort()
			return
		}

		// 将用户信息存储到上下文
		c.Set(string(ContextKeyUserID), tokenInfo.UserID)
		c.Set(string(ContextKeyDeviceType), tokenInfo.DeviceType)
		c.Set(string(ContextKeyTokenInfo), tokenInfo)

		c.Next()
	}
}

// PermissionMiddleware 创建权限认证中间件
func PermissionMiddleware(obj string, act string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取权限管理器
		permManager, err := GetPermissionManager()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "权限管理器未初始化",
			})
			c.Abort()
			return
		}

		// 从上下文获取用户ID
		userID, exists := c.Get(string(ContextKeyUserID))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未认证的请求",
			})
			c.Abort()
			return
		}

		// 检查权限
		if !permManager.CheckPermission(userID.(uint64), obj, act) {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// RoleMiddleware 创建角色认证中间件
func RoleMiddleware(role string) gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取权限管理器
		permManager, err := GetPermissionManager()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "权限管理器未初始化",
			})
			c.Abort()
			return
		}

		// 从上下文获取用户ID
		userID, exists := c.Get(string(ContextKeyUserID))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未认证的请求",
			})
			c.Abort()
			return
		}

		// 获取用户角色
		roles, err := permManager.GetRolesForUser(userID.(uint64))
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "获取角色信息失败",
			})
			c.Abort()
			return
		}

		// 检查角色
		hasRole := false
		for _, r := range roles {
			if r == role {
				hasRole = true
				break
			}
		}

		if !hasRole {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "角色权限不足",
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// SecondaryAuthMiddleware 创建二级认证中间件
func SecondaryAuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取二级认证验证器
		verifier, err := GetSecondaryAuthVerifier()
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"code": 500,
				"msg":  "二级认证验证器未初始化",
			})
			c.Abort()
			return
		}

		// 从上下文获取用户ID
		userID, exists := c.Get(string(ContextKeyUserID))
		if !exists {
			c.JSON(http.StatusUnauthorized, gin.H{
				"code": 401,
				"msg":  "未认证的请求",
			})
			c.Abort()
			return
		}

		// 检查是否已通过二级认证
		if !verifier.IsVerified(userID.(uint64)) {
			c.JSON(http.StatusForbidden, gin.H{
				"code": 403,
				"msg":  "需要二级认证",
				"data": gin.H{
					"require_secondary_auth": true,
				},
			})
			c.Abort()
			return
		}

		c.Next()
	}
}

// GetCurrentUserID 从上下文获取当前用户ID
func GetCurrentUserID(c *gin.Context) (uint64, bool) {
	userID, exists := c.Get(string(ContextKeyUserID))
	if !exists {
		return 0, false
	}
	return userID.(uint64), true
}

// GetCurrentDeviceType 从上下文获取当前设备类型
func GetCurrentDeviceType(c *gin.Context) (string, bool) {
	deviceType, exists := c.Get(string(ContextKeyDeviceType))
	if !exists {
		return "", false
	}
	return deviceType.(string), true
}

// GetCurrentTokenInfo 从上下文获取当前令牌信息
func GetCurrentTokenInfo(c *gin.Context) (*TokenInfo, bool) {
	tokenInfo, exists := c.Get(string(ContextKeyTokenInfo))
	if !exists {
		return nil, false
	}
	return tokenInfo.(*TokenInfo), true
}
