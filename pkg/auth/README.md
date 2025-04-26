# Auth 包

## 功能概述

`auth`包提供了完整的用户认证与权限管理解决方案，适用于需要严格访问控制的Web应用。

### 主要功能

- **登录认证**: 
  - 单端登录：一个用户只能在一个设备上登录，新登录会踢掉旧登录
  - 多端登录：一个用户可以在多个设备上同时登录
  - 同端互斥登录：同一类型设备只能登录一个账号，如手机端只能登录一个账号
  - 七天内免登录：支持长效会话，七天内无需重新登录
  - 令牌刷新：支持访问令牌过期后使用刷新令牌获取新的访问令牌

- **权限认证**: 
  - 基于 Casbin 的权限控制：使用 RBAC 模型进行精细化权限管理
  - 角色认证：基于角色的访问控制
  - 会话二级认证：敏感操作需要二次验证
  - 权限策略动态管理：支持运行时添加、删除、更新权限策略

## 架构设计

### 核心组件

1. **TokenManager**: 负责生成、验证、刷新和撤销令牌
2. **SessionStore**: 会话存储接口，支持多种后端存储（Redis、数据库等）
3. **AuthService**: 认证服务，处理登录、登出等操作
4. **PermissionManager**: 权限管理器，基于 Casbin 实现权限检查
5. **SecondaryAuthVerifier**: 二级认证验证器，用于敏感操作的二次验证
6. **Middleware**: 提供各种认证中间件，用于 Gin 框架集成

### 数据流

1. 用户登录 → 验证凭证 → 生成令牌 → 存储会话 → 返回令牌
2. 请求资源 → 验证令牌 → 检查权限 → 允许/拒绝访问
3. 敏感操作 → 二级认证 → 验证通过 → 允许操作

## 使用方法

### 初始化

```go
import (
    "time"
    
    "github.com/your-org/sweet/pkg/auth"
    "github.com/redis/go-redis/v9"
    "gorm.io/gorm"
)

// 初始化认证服务
func initAuth(redisClient *redis.Client, db *gorm.DB) error {
    // 1. 创建会话存储
    sessionStore := auth.NewRedisSessionStore(redisClient)
    
    // 2. 创建认证配置
    config := auth.Config{
        SecretKey:          "your-secret-key-should-be-long-and-secure", // 生产环境应从配置或环境变量获取
        AccessTokenExpiry:  24 * time.Hour,                              // 访问令牌24小时过期
        RefreshTokenExpiry: 7 * 24 * time.Hour,                          // 刷新令牌7天过期
        LoginMode:          auth.LoginModeSingle,                        // 单端登录模式
    }
    
    // 3. 创建认证服务
    authService := auth.NewAuthService(config, sessionStore)
    
    // 4. 创建权限管理器
    permManager, err := auth.NewPermissionManager("", "mysql", db)
    if err != nil {
        return fmt.Errorf("初始化权限管理器失败: %w", err)
    }
    
    // 5. 创建二级认证验证器
    secondaryAuthVerifier := auth.NewRedisSecondaryAuthVerifier(redisClient, 5*time.Minute)
    
    // 6. 设置全局实例
    auth.SetAuthService(authService)
    auth.SetPermissionManager(permManager)
    auth.SetSecondaryAuthVerifier(secondaryAuthVerifier)
    
    return nil
}
```

### 登录认证流程

```go
import (
    "github.com/gin-gonic/gin"
    "github.com/your-org/sweet/pkg/auth"
    "net/http"
)

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
    authService := auth.GetAuthService()
    
    // 生成令牌
    accessToken, refreshToken, err := authService.Login(user.ID, req.DeviceType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "登录失败"})
        return
    }
    
    // 获取令牌过期时间
    tokenInfo, _ := authService.VerifyToken(accessToken)
    expiresIn := int64(tokenInfo.ExpiresAt.Sub(time.Now()).Seconds())
    
    c.JSON(http.StatusOK, gin.H{
        "code": 200,
        "msg": "登录成功",
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
    userID, _ := auth.GetCurrentUserID(c)
    deviceType, _ := auth.GetCurrentDeviceType(c)
    
    // 登出
    err := auth.GetAuthService().Logout(userID, deviceType)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"code": 500, "msg": "登出失败"})
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
    
    // 刷新令牌
    authService := auth.GetAuthService()
    accessToken, refreshToken, err := authService.RefreshToken(req.RefreshToken)
    if err != nil {
        var msg string
        if err == auth.ErrTokenExpired {
            msg = "刷新令牌已过期"
        } else {
            msg = "刷新令牌无效"
        }
        c.JSON(http.StatusUnauthorized, gin.H{"code": 401, "msg": msg})
        return
    }
    
    // 获取令牌过期时间
    tokenInfo, _ := authService.VerifyToken(accessToken)
    expiresIn := int64(tokenInfo.ExpiresAt.Sub(time.Now()).Seconds())
    
    c.JSON(http.StatusOK, gin.H{
        "code": 200,
        "msg": "刷新成功",
        "data": LoginResponse{
            AccessToken:  accessToken,
            RefreshToken: refreshToken,
            ExpiresIn:    expiresIn,
        },
    })
}
```

### 权限管理

```go
// 添加角色
func addRole(userID uint64, role string) error {
    permManager := auth.GetPermissionManager()
    return permManager.AddRoleForUser(userID, role)
}

// 添加权限策略
func addPolicy(role, obj, act string) error {
    permManager := auth.GetPermissionManager()
    return permManager.AddPolicy(role, obj, act)
}

// 检查用户是否有特定权限
func checkPermission(userID uint64, obj, act string) bool {
    permManager := auth.GetPermissionManager()
    return permManager.CheckPermission(userID, obj, act)
}

// 获取用户角色
func getUserRoles(userID uint64) ([]string, error) {
    permManager := auth.GetPermissionManager()
    return permManager.GetRolesForUser(userID)
}
```

### 二级认证

```go
// 生成二级认证码
func generateSecondaryAuthCode(userID uint64) (string, error) {
    verifier := auth.GetSecondaryAuthVerifier()
    return verifier.GenerateCode(userID)
}

// 验证二级认证码
func verifySecondaryAuthCode(userID uint64, code string) bool {
    verifier := auth.GetSecondaryAuthVerifier()
    return verifier.VerifyCode(userID, code)
}

// 检查用户是否已通过二级认证
func isSecondaryAuthVerified(userID uint64) bool {
    verifier := auth.GetSecondaryAuthVerifier()
    return verifier.IsVerified(userID)
}
```

## 中间件使用

### JWT认证中间件

```go
// 需要JWT认证的路由
authorized := r.Group("/api")
authorized.Use(auth.JWTAuthMiddleware())
{
    authorized.GET("/user/profile", getUserProfile)
    authorized.POST("/logout", logout)
}
```

### 角色认证中间件

```go
// 需要管理员角色的路由
admin := r.Group("/api/admin")
admin.Use(auth.JWTAuthMiddleware(), auth.RoleMiddleware("admin"))
{
    admin.GET("/users", listUsers)
    admin.POST("/settings", updateSettings)
}
```

### 权限认证中间件

```go
// 需要特定权限的路由
products := r.Group("/api/products")
products.Use(auth.JWTAuthMiddleware())
{
    // 需要products:read权限
    products.GET("", auth.PermissionMiddleware("products", "read"), listProducts)
    // 需要products:create权限
    products.POST("", auth.PermissionMiddleware("products", "create"), createProduct)
    // 需要products:update权限
    products.PUT("/:id", auth.PermissionMiddleware("products", "update"), updateProduct)
    // 需要products:delete权限
    products.DELETE("/:id", auth.PermissionMiddleware("products", "delete"), deleteProduct)
}
```

### 二级认证中间件

```go
// 需要二级认证的路由
sensitive := r.Group("/api/sensitive")
sensitive.Use(auth.JWTAuthMiddleware(), auth.SecondaryAuthMiddleware())
{
    sensitive.POST("/payment", processPayment)
    sensitive.PUT("/user/password", changePassword)
}
```

## API参考

### 认证服务 (AuthService)

| 方法 | 描述 |
| --- | --- |
| `Login(userID uint64, deviceType string) (accessToken, refreshToken string, err error)` | 用户登录，生成访问令牌和刷新令牌 |
| `Logout(userID uint64, deviceType string) error` | 用户登出，撤销令牌 |
| `RefreshToken(refreshToken string) (accessToken, newRefreshToken string, err error)` | 使用刷新令牌获取新的访问令牌 |
| `VerifyToken(token string) (*TokenInfo, error)` | 验证令牌有效性，返回令牌信息 |
| `RevokeAllUserTokens(userID uint64) error` | 撤销用户的所有令牌 |

### 权限管理器 (PermissionManager)

| 方法 | 描述 |
| --- | --- |
| `AddRoleForUser(userID uint64, role string) error` | 为用户添加角色 |
| `DeleteRoleForUser(userID uint64, role string) error` | 删除用户的角色 |
| `GetRolesForUser(userID uint64) ([]string, error)` | 获取用户的所有角色 |
| `AddPolicy(role, obj, act string) error` | 添加权限策略 |
| `RemovePolicy(role, obj, act string) error` | 删除权限策略 |
| `CheckPermission(userID uint64, obj, act string) bool` | 检查用户是否有特定权限 |

### 二级认证验证器 (SecondaryAuthVerifier)

| 方法 | 描述 |
| --- | --- |
| `GenerateCode(userID uint64) (string, error)` | 生成二级认证码 |
| `VerifyCode(userID uint64, code string) bool` | 验证二级认证码 |
| `IsVerified(userID uint64) bool` | 检查用户是否已通过二级认证 |

### 上下文辅助函数

| 函数 | 描述 |
| --- | --- |
| `GetCurrentUserID(c *gin.Context) (uint64, bool)` | 从上下文获取当前用户ID |
| `GetCurrentDeviceType(c *gin.Context) (string, bool)` | 从上下文获取当前设备类型 |
| `GetCurrentTokenInfo(c *gin.Context) (*TokenInfo, bool)` | 从上下文获取当前令牌信息 |

## 配置参考

### 认证配置 (Config)

| 字段 | 类型 | 描述 | 默认值 |
| --- | --- | --- | --- |
| `SecretKey` | `string` | 用于签名令牌的密钥 | 必填 |
| `AccessTokenExpiry` | `time.Duration` | 访问令牌的过期时间 | `24 * time.Hour` |
| `RefreshTokenExpiry` | `time.Duration` | 刷新令牌的过期时间 | `7 * 24 * time.Hour` |
| `LoginMode` | `LoginMode` | 登录模式 | `LoginModeSingle` |

### 登录模式 (LoginMode)

| 常量 | 描述 |
| --- | --- |
| `LoginModeSingle` | 单端登录：一个用户只能在一个设备上登录 |
| `LoginModeMulti` | 多端登录：一个用户可以在多个设备上同时登录 |
| `LoginModeMutex` | 同端互斥登录：同一类型设备只能登录一个账号 |

## 错误码参考

| 错误 | 描述 |
| --- | --- |
| `ErrTokenExpired` | 令牌已过期 |
| `ErrTokenInvalid` | 令牌无效 |
| `ErrTokenRevoked` | 令牌已被撤销 |
| `ErrSessionNotFound` | 会话不存在 |
| `ErrPermissionDenied` | 权限不足 |
| `ErrRoleNotFound` | 角色不存在 |
| `ErrSecondaryAuthRequired` | 需要二级认证 |
| `ErrSecondaryAuthCodeInvalid` | 二级认证码无效 |

## 最佳实践

1. **安全存储密钥**：不要在代码中硬编码 `SecretKey`，应从环境变量或安全的配置管理系统获取。

2. **令牌过期时间**：根据应用的安全需求设置合适的令牌过期时间。敏感应用应使用较短的过期时间。

3. **权限粒度**：设计合理的权限粒度，避免过于细粒度导致管理复杂，也避免过于粗粒度导致权限控制不够精确。

4. **二级认证**：对敏感操作（如支付、修改密码）启用二级认证，提高安全性。

5. **错误处理**：妥善处理认证和权限错误，向用户提供清晰但不过于详细的错误信息，避免泄露系统细节。

6. **日志记录**：记录重要的认证和权限事件，如登录尝试、权限变更、敏感操作等，便于审计和问题排查。

7. **定期刷新权限**：在长时间运行的应用中，考虑定期刷新权限缓存，确保权限变更能及时生效。