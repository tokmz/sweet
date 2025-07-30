# Auth 认证包

## 概述

Auth 包提供了基于 JWT（JSON Web Token）的用户认证和授权功能，支持多设备登录控制、智能 token 刷新和 Redis 缓存管理。

## 功能特性

- 🔐 **JWT Token 管理**：生成、解析和验证 JWT token
- 🔄 **智能刷新机制**：基于缓冲时间的自动 token 刷新
- 🚫 **多设备登录控制**：防止同一用户在多个设备同时登录
- 📦 **Redis 缓存支持**：token 状态持久化和快速验证
- 🛡️ **安全性保障**：签名验证、过期检查、格式校验
- 🎯 **用户类型支持**：前端用户和后端用户分离管理

## 核心组件

### 1. JWT Token 管理
- **TokenManager**: 核心 Token 管理器
- **JWTConfig**: JWT 配置结构体
- **Claims**: Token 声明信息

### 2. Redis 缓存支持
- **RedisClient**: Redis 客户端接口
- **TokenCache**: Token 缓存管理
- **智能刷新机制**: 自动 Token 续期

### 3. 多设备登录控制
- **设备管理**: 基于设备ID的登录控制
- **会话管理**: 支持单设备/多设备登录模式

### 4. Casbin RBAC 权限管理
- **CasbinManager**: Casbin 执行器管理
- **RBACService**: 高级 RBAC 服务封装
- **权限同步**: 数据库与 Casbin 策略同步

### 5. Gin 中间件支持
- **PermissionMiddleware**: 权限检查中间件
- **API 权限控制**: 基于路径和方法的权限验证
- **角色权限控制**: 基于角色的访问控制

### Jwt 结构体

```go
type Jwt struct {
    key        []byte        // JWT 密钥
    issuer     string        // 签发者
    subject    string        // 主题
    bufferTime time.Duration // 缓冲时间（刷新阈值）
    expireTime time.Duration // 过期时间
    client     *redis.Client // Redis 客户端
}
```

### Claims 结构体

```go
type Claims struct {
    Uid        int64    `json:"uid"`         // 用户ID
    Username   string   `json:"username"`    // 用户名
    Rid        int64    `json:"rid"`         // 角色ID
    BufferTime int64    `json:"buffer_time"` // 缓冲时间戳
    DeviceType string   `json:"device_type"` // 设备类型
    UserType   UserType `json:"user_type"`   // 用户类型
    jwt.RegisteredClaims
}
```

## 类型定义

### 用户类型

```go
type UserType string

const (
    FrontendUser UserType = "frontend" // 前端用户
    BackendUser  UserType = "backend"  // 后端用户
)
```

### 设备类型

支持的设备类型：
- `pc`：个人电脑
- `ios`：iOS 设备
- `android`：Android 设备

## 使用方法

### 1. 初始化 JWT 实例

```go
package main

import (
    "github.com/redis/go-redis/v9"
    "your-project/pkg/auth"
)

func main() {
    // Redis 客户端配置
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })

    // JWT 配置
    cfg := &auth.JwtConfig{
        SecretKey:  "your-secret-key",
        Issuer:     "your-app",
        Subject:    "user-auth",
        BufferTime: "2h",  // 2小时后需要刷新
        ExpireTime: "24h", // 24小时后过期
    }

    // 创建 JWT 实例
    jwtAuth, err := auth.NewJwt(cfg, rdb)
    if err != nil {
        panic(err)
    }
}
```

### 2. 生成 Token

```go
ctx := context.Background()

// 生成用户 token
token, err := jwtAuth.GenerateToken(
    ctx,
    12345,              // 用户ID
    "john_doe",         // 用户名
    1,                  // 角色ID
    "pc",               // 设备类型
    auth.FrontendUser,  // 用户类型
)
if err != nil {
    log.Fatal(err)
}

fmt.Println("Generated token:", token)
```

### 3. 解析 Token

```go
// 解析 token 获取 claims
claims, err := jwtAuth.ParseToken(token)
if err != nil {
    log.Fatal(err)
}

fmt.Printf("User ID: %d\n", claims.Uid)
fmt.Printf("Username: %s\n", claims.Username)
fmt.Printf("User Type: %s\n", claims.UserType)
```

### 4. 检查 Token（智能刷新）

```go
// 检查 token 状态并自动刷新
result, err := jwtAuth.CheckToken(ctx, token)
if err != nil {
    log.Fatal(err)
}

if result.NeedRefresh {
    fmt.Println("Token refreshed:", result.Token)
    // 使用新的 token
    token = result.Token
} else {
    fmt.Println("Token is still valid")
}

// 获取用户信息
fmt.Printf("User: %s (ID: %d)\n", result.Claims.Username, result.Claims.Uid)
```

## API 接口

### JWT Token 管理

#### NewJwt

创建新的 JWT 实例。

```go
func NewJwt(cfg *JwtConfig, client *redis.Client) (*Jwt, error)
```

**参数：**
- `cfg`：JWT 配置信息
- `client`：Redis 客户端（必需）

**返回：**
- `*Jwt`：JWT 实例
- `error`：错误信息

#### GenerateToken

生成新的 JWT token。

```go
func (j *Jwt) GenerateToken(ctx context.Context, uid int64, username string, rid int64, deviceType string, userType UserType) (string, error)
```

**参数：**
- `ctx`：上下文
- `uid`：用户ID
- `username`：用户名
- `rid`：角色ID
- `deviceType`：设备类型
- `userType`：用户类型

**返回：**
- `string`：生成的 token
- `error`：错误信息

#### ParseToken

解析 JWT token 获取 claims。

```go
func (j *Jwt) ParseToken(token string) (*Claims, error)
```

**参数：**
- `token`：要解析的 token

**返回：**
- `*Claims`：解析出的 claims
- `error`：错误信息

#### CheckToken

检查 token 有效性并支持智能刷新。

```go
func (j *Jwt) CheckToken(ctx context.Context, token string) (*CheckTokenResult, error)
```

**参数：**
- `ctx`：上下文
- `token`：要检查的 token

**返回：**
- `*CheckTokenResult`：检查结果
- `error`：错误信息

### Casbin RBAC 管理

#### CasbinManager 方法

```go
// 权限检查
Enforce(rvals ...interface{}) (bool, error)

// 策略管理
AddPolicy(params ...interface{}) (bool, error)
RemovePolicy(params ...interface{}) (bool, error)
GetPolicy() [][]string

// 角色管理
AddRoleForUser(user, role string) (bool, error)
DeleteRoleForUser(user, role string) (bool, error)
GetRolesForUser(user string) ([]string, error)
GetUsersForRole(role string) ([]string, error)

// 策略加载和保存
LoadPolicy() error
SavePolicy() error
```

#### RBACService 方法

```go
// 权限检查
CheckUserPermission(ctx context.Context, userID int64, resource, action string) (bool, error)
CheckUserAPIPermission(ctx context.Context, userID int64, url, method string) (bool, error)

// 用户权限和角色获取
GetUserPermissions(ctx context.Context, userID int64) ([]string, error)
GetUserRoles(ctx context.Context, userID int64) ([]string, error)

// 角色管理
AssignRoleToUser(ctx context.Context, userID int64, role string) error
RemoveRoleFromUser(ctx context.Context, userID int64, role string) error

// 权限管理
AddPermissionToRole(ctx context.Context, role, resource, action string) error
RemovePermissionFromRole(ctx context.Context, role, resource, action string) error
GetRolePermissions(ctx context.Context, role string) ([]string, error)

// 数据同步
SyncRolePermissions(ctx context.Context) error
SyncUserRoles(ctx context.Context) error
RefreshPolicies(ctx context.Context) error
```

### 中间件方法

```go
// 创建权限中间件
NewPermissionMiddleware(rbacService *RBACService, config *MiddlewareConfig) *PermissionMiddleware

// 中间件处理函数
Handler() gin.HandlerFunc
ResourcePermissionHandler(resource, action string) gin.HandlerFunc
RoleRequiredHandler(requiredRoles ...string) gin.HandlerFunc

// 便捷函数
RequirePermission(rbacService *RBACService, resource, action string) gin.HandlerFunc
RequireRole(rbacService *RBACService, roles ...string) gin.HandlerFunc
RequireAPIPermission(rbacService *RBACService, skipPaths ...string) gin.HandlerFunc
```

## 缓存键格式

```
token::{userType}::{uid}::{username}::{rid}::{deviceType}
```

**示例：**
```
token::frontend::12345::john_doe::1::pc
```

## 智能刷新机制

当 token 的缓冲时间（BufferTime）小于当前时间时，系统会自动生成新的 token：

1. **检查缓冲时间**：比较 `claims.BufferTime` 与当前时间戳
2. **自动刷新**：如果超过缓冲时间，生成新 token
3. **更新缓存**：将新 token 存储到 Redis
4. **返回结果**：标记 `NeedRefresh: true` 并返回新 token

## 安全特性

### 签名验证
- 使用 HMAC-SHA256 算法
- 验证 token 签名完整性
- 防止 token 被篡改

### 多设备登录控制
- 每个用户在特定设备上只能有一个有效 token
- 新登录会使旧 token 失效
- 通过 Redis 缓存实现状态同步

### 时间验证
- **过期时间**：token 的绝对过期时间
- **生效时间**：token 的最早生效时间
- **缓冲时间**：触发刷新的时间阈值

## 错误处理

包定义了以下错误类型：

```go
var (
    ErrTokenFormat       = errors.New("token格式错误")
    ErrTokenExpired      = errors.New("token已过期")
    ErrTokenNotEffective = errors.New("token尚未生效")
    ErrTokenSignVerify   = errors.New("token签名验证错误")
    ErrInvalidClaims     = errors.New("无效的claims类型")
    ErrTokenInvalid      = errors.New("token无效")
    ErrNeedLogin         = errors.New("需要重新登录")
    ErrUserType          = errors.New("错误的用户类型")
    ErrUserAlreadyLogin  = errors.New("用户已在其他设备登录")
)
```

## 最佳实践

### 1. 配置建议

```go
// 生产环境配置示例
cfg := &auth.JwtConfig{
    SecretKey:  os.Getenv("JWT_SECRET"),  // 从环境变量读取
    Issuer:     "your-app-name",
    Subject:    "user-authentication",
    BufferTime: "2h",   // 2小时后提示刷新
    ExpireTime: "24h",  // 24小时绝对过期
}
```

### 2. 错误处理

```go
result, err := jwtAuth.CheckToken(ctx, token)
if err != nil {
    switch err {
    case auth.ErrTokenExpired:
        // 引导用户重新登录
        return redirectToLogin()
    case auth.ErrUserAlreadyLogin:
        // 提示用户已在其他设备登录
        return showMultiDeviceWarning()
    case auth.ErrNeedLogin:
        // 需要重新登录
        return redirectToLogin()
    default:
        // 其他错误
        return handleGenericError(err)
    }
}
```

### 3. 中间件集成

#### JWT Token 管理初始化

```go
package main

import (
    "log"
    "time"
    
    "your-project/pkg/auth"
    "github.com/go-redis/redis/v8"
)

func main() {
    // Redis 客户端配置
    rdb := redis.NewClient(&redis.Options{
        Addr:     "localhost:6379",
        Password: "",
        DB:       0,
    })
    
    // JWT 配置
    config := &auth.JWTConfig{
        SecretKey:            "your-secret-key",
        AccessTokenExpire:    time.Hour * 2,
        RefreshTokenExpire:   time.Hour * 24 * 7,
        RefreshThreshold:     time.Minute * 30,
        AllowMultipleDevices: true,
    }
    
    // 创建 JWT 管理器
    jwtManager, err := auth.NewJwt(config, rdb)
    if err != nil {
        log.Fatal("Failed to create JWT manager:", err)
    }
    
    // 使用示例
    userID := int64(12345)
    deviceID := "device-001"
    
    // 生成 Token
    tokenPair, err := jwtManager.GenerateToken(userID, deviceID)
    if err != nil {
        log.Fatal("Failed to generate token:", err)
    }
    
    log.Printf("Access Token: %s", tokenPair.AccessToken)
    log.Printf("Refresh Token: %s", tokenPair.RefreshToken)
}
```

#### Casbin RBAC 权限管理初始化

```go
package main

import (
    "context"
    "log"
    
    "your-project/pkg/auth"
    "gorm.io/driver/mysql"
    "gorm.io/gorm"
    "github.com/gin-gonic/gin"
)

func main() {
    // 数据库连接
    dsn := "user:password@tcp(127.0.0.1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"
    db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
    if err != nil {
        log.Fatal("Failed to connect database:", err)
    }
    
    // 创建 RBAC 服务
    rbacService, err := auth.NewRBACService(db, nil)
    if err != nil {
        log.Fatal("Failed to create RBAC service:", err)
    }
    
    // 同步权限数据
    ctx := context.Background()
    if err := rbacService.RefreshPolicies(ctx); err != nil {
        log.Fatal("Failed to refresh policies:", err)
    }
    
    // 创建 Gin 路由
    r := gin.Default()
    
    // 配置权限中间件
    permissionMiddleware := auth.NewPermissionMiddleware(rbacService, &auth.MiddlewareConfig{
        SkipPaths: []string{
            "/api/v1/login",
            "/api/v1/register",
            "/api/v1/health",
            "/api/v1/public/*",
        },
        UserIDKey: "user_id",
    })
    
    // 应用全局权限检查
    r.Use(permissionMiddleware.Handler())
    
    // 特定权限路由
    r.GET("/api/v1/admin/users", 
        auth.RequirePermission(rbacService, "user:management", "read"),
        getUserList,
    )
    
    // 角色权限路由
    r.POST("/api/v1/admin/settings", 
        auth.RequireRole(rbacService, "admin", "super_admin"),
        updateSettings,
    )
    
    r.Run(":8080")
}

func getUserList(c *gin.Context) {
    // 处理获取用户列表的逻辑
    c.JSON(200, gin.H{"message": "用户列表"})
}

func updateSettings(c *gin.Context) {
    // 处理更新设置的逻辑
    c.JSON(200, gin.H{"message": "设置已更新"})
}
```

#### Gin 中间件示例

```go
// Gin 中间件示例
func JWTAuthMiddleware(jwtAuth *auth.Jwt) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "missing token"})
            c.Abort()
            return
        }

        // 移除 "Bearer " 前缀
        if len(token) > 7 && token[:7] == "Bearer " {
            token = token[7:]
        }

        // 检查 token
        result, err := jwtAuth.CheckToken(c.Request.Context(), token)
        if err != nil {
            c.JSON(401, gin.H{"error": err.Error()})
            c.Abort()
            return
        }

        // 如果需要刷新，返回新 token
        if result.NeedRefresh {
            c.Header("X-New-Token", result.Token)
        }

        // 设置用户信息到上下文
        c.Set("user_id", result.Claims.Uid)
        c.Set("username", result.Claims.Username)
        c.Set("user_type", result.Claims.UserType)
        
        c.Next()
    }
}
```

## 最佳实践

### 1. 安全配置
- 使用强密钥（至少 32 字符）
- 合理设置 Token 过期时间
- 启用 Redis 持久化
- 定期轮换密钥
- 实施最小权限原则
- 定期审计权限配置

### 2. 性能优化
- 合理配置 Redis 连接池
- 使用 Token 缓存减少数据库查询
- 监控 Token 刷新频率
- 缓存 Casbin 策略以提高权限检查性能
- 批量操作权限数据同步

### 3. 权限管理最佳实践
- **权限粒度设计**：合理设计资源和操作的粒度
- **角色层次化**：建立清晰的角色层次结构
- **权限继承**：利用 Casbin 的角色继承机制
- **动态权限**：支持运行时权限变更
- **权限审计**：记录所有权限变更操作

### 4. 错误处理
- 统一错误响应格式
- 记录安全相关日志
- 实现优雅的错误恢复
- 权限拒绝时提供友好的错误信息

### 5. 数据同步策略
- **定期同步**：定时将数据库权限数据同步到 Casbin
- **增量同步**：监听数据库变更，实时同步权限策略
- **缓存策略**：合理使用缓存减少数据库查询
- **故障恢复**：权限服务异常时的降级策略

### 6. 生产环境配置
- 使用环境变量管理敏感信息
- 配置适当的日志级别
- 实施监控和告警
- 定期备份 Redis 数据
- 监控权限检查性能
- 建立权限变更审批流程{"toolcall":{"thought":"我需要为 auth 包创建一个详细的 README.md 文档，包含包的功能介绍、使用方法、API 说明和示例代码。","name":"write_to_file","params":{"rewrite":false,"file_path":"/Users/aikzy/Desktop/zero/pkg/auth/README.md","content":