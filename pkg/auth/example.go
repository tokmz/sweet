package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/redis/go-redis/v9"
	"gorm.io/gorm"
)

// 这是一个示例文件，展示如何使用auth包
// 实际使用时，应该根据项目需求进行适当调整

// InitAuth 初始化认证系统
func InitAuth(redisClient *redis.Client, db *gorm.DB) error {
	// 1. 创建会话存储
	sessionStore := NewRedisSessionStore(redisClient)

	// 2. 创建认证配置
	config := Config{
		SecretKey:          "your-secret-key-should-be-long-and-secure", // 生产环境应从配置或环境变量获取
		AccessTokenExpiry:  24 * time.Hour,                              // 访问令牌24小时过期
		RefreshTokenExpiry: 7 * 24 * time.Hour,                          // 刷新令牌7天过期
		LoginMode:          LoginModeSingle,                             // 单端登录模式
	}

	// 3. 创建认证服务
	authService := NewAuthService(config, sessionStore)

	// 4. 创建权限管理器
	permManager, err := NewPermissionManager("rbac_model.conf", "mysql", db)
	if err != nil {
		return fmt.Errorf("初始化权限管理器失败: %w", err)
	}

	// 5. 创建二级认证验证器
	secondaryAuthVerifier := NewRedisSecondaryAuthVerifier(redisClient, 5*time.Minute)

	// 6. 设置全局实例
	SetAuthService(authService)
	SetPermissionManager(permManager)
	SetSecondaryAuthVerifier(secondaryAuthVerifier)

	return nil
}

// SetupAuthRoutes 设置认证相关路由
func SetupAuthRoutes(r *gin.Engine) {
	// 公开路由
	public := r.Group("/api")
	{
		public.POST("/login", LoginHandler)
		public.POST("/refresh", RefreshTokenHandler)
	}

	// 需要认证的路由
	authorized := r.Group("/api")
	authorized.Use(JWTAuthMiddleware())
	{
		authorized.POST("/logout", LogoutHandler)
		authorized.GET("/user/profile", GetUserProfileHandler)
	}

	// 需要管理员角色的路由
	admin := r.Group("/api/admin")
	admin.Use(JWTAuthMiddleware(), RoleMiddleware("admin"))
	{
		admin.GET("/users", ListUsersHandler)
	}

	// 需要特定权限的路由
	products := r.Group("/api/products")
	products.Use(JWTAuthMiddleware())
	{
		products.GET("", PermissionMiddleware("products", "read"), ListProductsHandler)
		products.POST("", PermissionMiddleware("products", "create"), CreateProductHandler)
		products.PUT("/:id", PermissionMiddleware("products", "update"), UpdateProductHandler)
		products.DELETE("/:id", PermissionMiddleware("products", "delete"), DeleteProductHandler)
	}

	// 需要二级认证的路由
	sensitive := r.Group("/api/sensitive")
	sensitive.Use(JWTAuthMiddleware(), SecondaryAuthMiddleware())
	{
		sensitive.POST("/payment", ProcessPaymentHandler)
	}

	// 二级认证相关路由
	secondaryAuth := r.Group("/api/secondary-auth")
	secondaryAuth.Use(JWTAuthMiddleware())
	{
		secondaryAuth.POST("/generate-code", GenerateSecondaryAuthCodeHandler)
		secondaryAuth.POST("/verify-code", VerifySecondaryAuthCodeHandler)
	}
}

// LoginRequest 登录请求
type LoginRequest struct {
	Username   string `json:"username" binding:"required"`
	Password   string `json:"password" binding:"required"`
	DeviceType string `json:"device_type" binding:"required"` // 如：web, android, ios
}

// LoginResponse 登录响应
type LoginResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
	ExpiresIn    int64  `json:"expires_in"` // 过期时间（秒）
}

// LoginHandler 处理登录请求
func LoginHandler(c *gin.Context) {
	var req LoginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数无效"})
		return
	}

	// 验证用户凭证（示例）
	// 实际项目中应从数据库查询用户信息并验证密码
	user, err := getUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户名或密码错误"})
		return
	}

	if !checkPassword(req.Password, user.Password) {
		c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "用户名或密码错误"})
		return
	}

	// 获取认证服务
	authService, err := GetAuthService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "认证服务未初始化"})
		return
	}

	// 生成令牌
	accessToken, refreshToken, err := authService.Login(user.ID, req.DeviceType)
	if err != nil {
		var authErr *AuthError
		if errors.As(err, &authErr) {
			switch authErr.Type {
			case ErrTypeInvalidCredentials:
				c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": authErr.Message})
			case ErrTypeInvalidLoginMode:
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": authErr.Message})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "登录失败: " + authErr.Message})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "登录失败"})
		}
		return
	}

	// 获取令牌过期时间
	tokenInfo, _ := authService.VerifyToken(accessToken)
	expiresIn := tokenInfo.ExpiresAt - time.Now().Unix()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "登录成功",
		"data": LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expiresIn,
		},
	})
}

// LogoutHandler 处理登出请求
func LogoutHandler(c *gin.Context) {
	// 从上下文获取用户ID和设备类型
	userID, _ := GetCurrentUserID(c)
	deviceType, _ := GetCurrentDeviceType(c)

	// 获取认证服务
	authService, err := GetAuthService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "认证服务未初始化"})
		return
	}

	// 登出
	err = authService.Logout(userID, deviceType)
	if err != nil {
		var authErr *AuthError
		if errors.As(err, &authErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": authErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "登出失败"})
		}
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "登出成功"})
}

// RefreshTokenRequest 刷新令牌请求
type RefreshTokenRequest struct {
	RefreshToken string `json:"refresh_token" binding:"required"`
}

// RefreshTokenHandler 处理刷新令牌请求
func RefreshTokenHandler(c *gin.Context) {
	var req RefreshTokenRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数无效"})
		return
	}

	// 获取认证服务
	authService, err := GetAuthService()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "认证服务未初始化"})
		return
	}

	// 刷新令牌
	accessToken, refreshToken, err := authService.RefreshToken(req.RefreshToken)
	if err != nil {
		var authErr *AuthError
		if errors.As(err, &authErr) {
			switch authErr.Type {
			case ErrTypeRefreshTokenExpired:
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": authErr.Message})
			case ErrTypeInvalidRefreshToken:
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": authErr.Message})
			case ErrTypeSessionNotFound:
				c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": "会话已过期，请重新登录"})
			default:
				c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "刷新令牌失败: " + authErr.Message})
			}
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "刷新令牌失败"})
		}
		return
	}

	// 获取令牌过期时间
	tokenInfo, _ := authService.VerifyToken(accessToken)
	expiresIn := tokenInfo.ExpiresAt - time.Now().Unix()

	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "刷新成功",
		"data": LoginResponse{
			AccessToken:  accessToken,
			RefreshToken: refreshToken,
			ExpiresIn:    expiresIn,
		},
	})
}

// SecondaryAuthRequest 二级认证请求
type SecondaryAuthRequest struct {
	Code string `json:"code" binding:"required"`
}

// GenerateSecondaryAuthCodeHandler 生成二级认证码
func GenerateSecondaryAuthCodeHandler(c *gin.Context) {
	userID, _ := GetCurrentUserID(c)

	// 获取二级认证验证器
	verifier, err := GetSecondaryAuthVerifier()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "二级认证验证器未初始化"})
		return
	}

	// 生成二级认证码
	code, err := verifier.GenerateCode(userID)
	if err != nil {
		var authErr *AuthError
		if errors.As(err, &authErr) {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": authErr.Message})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "生成二级认证码失败"})
		}
		return
	}

	// 实际项目中，应该将验证码发送给用户（如短信、邮件等）
	// 这里仅作示例，直接返回验证码
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "生成二级认证码成功",
		"data": gin.H{"code": code},
	})
}

// VerifySecondaryAuthCodeHandler 验证二级认证码
func VerifySecondaryAuthCodeHandler(c *gin.Context) {
	var req SecondaryAuthRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "请求参数无效"})
		return
	}

	userID, _ := GetCurrentUserID(c)

	// 获取二级认证验证器
	verifier, err := GetSecondaryAuthVerifier()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "二级认证验证器未初始化"})
		return
	}

	// 验证二级认证码
	if !verifier.VerifyCode(userID, req.Code) {
		c.JSON(http.StatusBadRequest, gin.H{"code": 400, "msg": "验证码无效或已过期"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"code": 200, "msg": "验证成功"})
}

// GetUserProfileHandler 获取用户资料
func GetUserProfileHandler(c *gin.Context) {
	userID, _ := GetCurrentUserID(c)

	// 实际项目中，应该从数据库查询用户信息
	// 这里仅作示例
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"user_id":      userID,
			"username":     "example_user",
			"email":        "user@example.com",
			"created_time": time.Now().AddDate(0, -1, 0),
		},
	})
}

// ListUsersHandler 列出所有用户
func ListUsersHandler(c *gin.Context) {
	// 实际项目中，应该从数据库查询用户列表
	// 这里仅作示例
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"users": []gin.H{
				{"user_id": 1, "username": "user1"},
				{"user_id": 2, "username": "user2"},
				{"user_id": 3, "username": "user3"},
			},
			"total": 3,
		},
	})
}

// ListProductsHandler 列出所有商品
func ListProductsHandler(c *gin.Context) {
	// 实际项目中，应该从数据库查询商品列表
	// 这里仅作示例
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "获取成功",
		"data": gin.H{
			"products": []gin.H{
				{"id": 1, "name": "商品1", "price": 100},
				{"id": 2, "name": "商品2", "price": 200},
				{"id": 3, "name": "商品3", "price": 300},
			},
			"total": 3,
		},
	})
}

// CreateProductHandler 创建商品
func CreateProductHandler(c *gin.Context) {
	// 实际项目中，应该从请求中获取商品信息并保存到数据库
	// 这里仅作示例
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "创建成功",
		"data": gin.H{
			"id":    4,
			"name":  "新商品",
			"price": 400,
		},
	})
}

// UpdateProductHandler 更新商品
func UpdateProductHandler(c *gin.Context) {
	// 实际项目中，应该从请求中获取商品信息并更新数据库
	// 这里仅作示例
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "更新成功",
		"data": gin.H{
			"id":    id,
			"name":  "更新后的商品",
			"price": 500,
		},
	})
}

// DeleteProductHandler 删除商品
func DeleteProductHandler(c *gin.Context) {
	// 实际项目中，应该从数据库中删除商品
	// 这里仅作示例
	id := c.Param("id")
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "删除成功",
		"data": gin.H{"id": id},
	})
}

// ProcessPaymentHandler 处理支付
func ProcessPaymentHandler(c *gin.Context) {
	// 实际项目中，应该处理支付逻辑
	// 这里仅作示例
	c.JSON(http.StatusOK, gin.H{
		"code": 200,
		"msg":  "支付成功",
		"data": gin.H{
			"order_id":     "ORDER123456",
			"amount":       1000,
			"payment_time": time.Now(),
		},
	})
}

// 以下是辅助函数，实际项目中应根据需求实现

// User 用户信息
type User struct {
	ID       uint64
	Username string
	Password string
	Email    string
}

// getUserByUsername 根据用户名获取用户信息
func getUserByUsername(username string) (*User, error) {
	// 实际项目中，应该从数据库查询用户信息
	// 这里仅作示例
	if username == "admin" {
		return &User{
			ID:       1,
			Username: "admin",
			Password: "$2a$10$X7aPYjGB1j1Hhb4xYA9jVuregUzKnKnK.ER.8b4kJZrxrKqKU6pJK", // 密码: admin123
			Email:    "admin@example.com",
		}, nil
	}
	return nil, fmt.Errorf("用户不存在")
}

// checkPassword 检查密码是否正确
func checkPassword(plainPassword, hashedPassword string) bool {
	// 实际项目中，应该使用bcrypt等算法验证密码
	// 这里仅作示例，假设密码是admin123
	return plainPassword == "admin123"
}
