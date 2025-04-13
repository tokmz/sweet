# Sweet Auth 认证与授权包

Sweet Auth 是一个功能丰富的认证与授权包，集成了 JWT 和 Casbin，提供了完整的用户认证、权限控制和会话管理功能。

## 主要功能

### 登录认证
- **多种登录模式**：单端登录、多端登录、同端互斥登录
- **记住我**：七天内免登录
- **Token管理**：生成、验证、刷新、撤销

### 权限认证
- **基于角色的访问控制**：使用 Casbin 实现 RBAC
- **权限粒度控制**：支持到 API 级别的权限控制
- **动态权限**：支持运行时更新权限规则
- **通配符权限**：支持使用通配符进行权限匹配

### 会话管理
- **踢人下线**：根据账号ID或Token值强制用户下线
- **会话监控**：实时监控在线用户
- **会话存储**：支持内存存储和Redis持久化

### 单点登录(SSO)
- **多种SSO模式**：同域、跨域、基于Redis
- **单点注销**：一处注销，全局生效

### 高级特性
- **二级认证**：敏感操作需要二次验证
- **全局侦听器**：登录、注销等关键操作的AOP处理
- **参数签名**：API调用签名校验，防篡改，防重放
- **自动续签**：Token自动续期

## 快速开始

### 初始化

```go
// 创建认证管理器
authManager, err := auth.NewManager(auth.Config{
    JWTSecret:      "your-jwt-secret",
    TokenExpire:    24 * time.Hour,
    RefreshExpire:  7 * 24 * time.Hour,
    TokenStyle:     auth.TokenStyleUUID,
    LoginModel:     auth.LoginModelMulti,
    EnableRedis:    true,
    RedisKeyPrefix: "sweet:auth:",
})
if err != nil {
    panic(err)
}

// 初始化Casbin
enforcer, err := authManager.InitCasbin(auth.CasbinConfig{
    Model:        "config/rbac_model.conf",
    Adapter:      "mysql",
    DSN:          "user:pass@tcp(127.0.0.1:3306)/dbname",
    TableName:    "casbin_rule",
    AutoMigrate:  true,
})
if err != nil {
    panic(err)
}
```

### 登录认证

```go
// 登录并获取Token
token, refreshToken, err := authManager.Login(auth.LoginParams{
    UserID:    123,
    Username:  "admin",
    Device:    auth.DeviceTypeWeb,
    RememberMe: true,
    ExtraData: map[string]interface{}{
        "role": "admin",
    },
})
if err != nil {
    // 处理登录失败
}

// 验证Token
info, err := authManager.ValidateToken(token)
if err != nil {
    // Token无效
}

// 注销
err = authManager.Logout(token)
if err != nil {
    // 处理注销失败
}
```

### 权限控制

```go
// 检查权限
allowed, err := authManager.CheckPermission(123, "/api/users", "GET")
if err != nil || !allowed {
    // 无权限访问
}

// 为用户分配角色
err = authManager.AssignRole(123, "admin")
if err != nil {
    // 角色分配失败
}

// 为角色添加权限
err = authManager.AddPermissionForRole("admin", "/api/users", "GET")
if err != nil {
    // 权限添加失败
}
```

### 会话管理

```go
// 获取在线用户
onlineUsers, err := authManager.GetOnlineUsers()

// 踢人下线
err = authManager.KickoutByUserID(123)
if err != nil {
    // 踢人失败
}

// 踢出指定Token
err = authManager.KickoutByToken(token)
if err != nil {
    // 踢人失败
}
```

## 权限管理详解

### RBAC模型说明

Sweet Auth 采用基于角色的访问控制(RBAC)模型，核心概念如下：

1. **用户(User)**：系统中的用户实体，通常用用户ID表示
2. **角色(Role)**：一组权限的集合，如管理员、操作员、普通用户等
3. **对象(Object)**：被访问的资源，如API路径、功能模块等
4. **操作(Action)**：对资源的操作，如GET、POST、DELETE等
5. **策略(Policy)**：定义了角色对资源的操作权限

RBAC模型示例：
- 用户123被分配了角色"admin"
- 角色"admin"拥有对"/api/users"进行"GET"操作的权限
- 所以用户123可以对"/api/users"进行"GET"操作

### 权限管理示例

#### 角色管理

```go
// 为用户分配角色
err := authManager.AssignRole(123, "admin")

// 为用户分配多个角色
err = authManager.AssignRole(123, "editor")
err = authManager.AssignRole(123, "viewer")

// 移除用户角色
err = authManager.RemoveRole(123, "editor")

// 检查用户是否有特定角色
hasRole, err := authManager.HasRole(123, "admin")

// 获取用户所有角色
roles, err := authManager.GetRolesForUser(123)
// 返回: ["admin", "viewer"]
```

#### 权限分配

```go
// 为角色添加权限 (角色, 资源, 操作)
err := authManager.AddPermissionForRole("admin", "/api/users", "GET")
err = authManager.AddPermissionForRole("admin", "/api/users", "POST")
err = authManager.AddPermissionForRole("editor", "/api/posts", "*")
err = authManager.AddPermissionForRole("viewer", "/api/posts", "GET")

// 直接获取用户权限
permissions, err := authManager.GetPermissionsForUser(123)
// 返回: [["/api/users", "GET"], ["/api/users", "POST"], ...]
```

#### 权限验证

```go
// 检查用户是否有权限执行操作
allowed, err := authManager.CheckPermission(123, "/api/users", "POST")
if allowed {
    // 允许访问
} else {
    // 拒绝访问
}

// 检查用户是否有权限访问通配符资源
allowed, err = authManager.CheckPermission(123, "/api/posts/123", "GET")
// 如果用户有 "/api/posts/*" 的权限，此处也会返回 true
```

### 通配符权限

Sweet Auth 支持使用通配符进行灵活的权限匹配，通配符规则如下：

1. **精确匹配**：完全一致的路径和操作
   - 例如: `/api/users` 与 `/api/users` 匹配
  
2. **`*` 匹配单级路径**：表示单个路径段的任意值
   - 例如: `/api/*/edit` 匹配 `/api/users/edit`, `/api/posts/edit`
   
3. **`**` 匹配多级路径**：表示任意多个路径段
   - 例如: `/api/**` 匹配 `/api/users`, `/api/posts/123/comments` 等所有以 `/api/` 开头的路径
   
4. **操作通配符**：使用 `*` 表示所有操作
   - 例如: 权限 `/api/users`, `*` 匹配该资源的所有操作（GET, POST, PUT, DELETE等）

通配符权限示例：

```go
// 添加通配符权限
err := authManager.AddPermissionForRole("admin", "/api/**", "*")         // 管理员可访问所有API
err = authManager.AddPermissionForRole("editor", "/api/posts/*", "*")    // 编辑者可完全管理任何单个文章
err = authManager.AddPermissionForRole("viewer", "/api/posts/**", "GET") // 访客只能查看所有文章资源

// 权限检查
allowed, _ := authManager.CheckPermission(123, "/api/users/456", "GET")    // 管理员可访问
allowed, _ = authManager.CheckPermission(456, "/api/posts/789", "PUT")     // 编辑者可更新文章
allowed, _ = authManager.CheckPermission(789, "/api/posts/123/comments", "GET") // 访客可查看文章评论
```

## 高级用法

### 二级认证

```go
// 开启二级认证
err = authManager.EnableSecondAuth(token, 30*time.Minute)
if err != nil {
    // 开启失败
}

// 验证二级认证
valid, err := authManager.CheckSecondAuth(token)
if err != nil || !valid {
    // 二级认证无效
}
```

### 单点登录

```go
// 配置SSO
authManager.ConfigSSO(auth.SSOConfig{
    Mode:      auth.SSOModeCrossRedis,
    RedisAddr: "127.0.0.1:6379",
    RedisDB:   0,
    RedisPass: "",
})

// 单点注销
err = authManager.SSOLogout(userID)
if err != nil {
    // 单点注销失败
}
```

### 临时Token

```go
// 创建临时Token
tempToken, err := authManager.CreateTempToken(123, 5*time.Minute, map[string]interface{}{
    "action": "reset-password",
})
if err != nil {
    // 创建失败
}

// 验证临时Token
claims, err := authManager.ValidateTempToken(tempToken)
if err != nil {
    // 临时Token无效
}
```

### 事件监听

```go
// 添加事件监听器
authManager.AddListener(func(event auth.AuthEvent) {
    switch event.Type {
    case auth.EventLogin:
        fmt.Printf("用户登录: %d, IP: %s\n", event.UserID, event.IP)
    case auth.EventLogout:
        fmt.Printf("用户注销: %d\n", event.UserID)
    case auth.EventKickout:
        fmt.Printf("用户被踢下线: %d\n", event.UserID)
    }
})
```

## 配置参考

### 认证管理器配置

```go
type Config struct {
    // JWT配置
    JWTSecret      string        // JWT密钥
    TokenExpire    time.Duration // Token过期时间
    RefreshExpire  time.Duration // 刷新Token过期时间
    TokenStyle     TokenStyle    // Token风格
    
    // 登录模式
    LoginModel     LoginModel    // 登录模式
    
    // Redis配置
    EnableRedis    bool          // 是否启用Redis
    RedisAddr      string        // Redis地址
    RedisPass      string        // Redis密码
    RedisDB        int           // Redis数据库
    RedisKeyPrefix string        // Redis键前缀
    
    // 高级配置
    AutoRenew      bool          // 是否自动续签
    RenewThreshold float64       // 续签阈值(0.0-1.0)
    EnableSecondAuth bool        // 是否启用二级认证
    SecondAuthExpire time.Duration // 二级认证过期时间
    EnableConcurrency bool        // 是否启用并发控制
    MaxConcurrency    int         // 最大并发数
}
```

### Casbin配置

```go
type CasbinConfig struct {
    Model        string // 模型配置文件路径
    Adapter      string // 适配器类型(file, mysql, postgres等)
    DSN          string // 数据源名称(适配器为数据库时使用)
    TableName    string // 表名(适配器为数据库时使用)
    AutoMigrate  bool   // 是否自动迁移表结构
}
```

## 错误处理

包中定义了一系列错误常量，可以用于错误处理：

```go
var (
    ErrTokenExpired       = errors.New("token已过期")
    ErrTokenInvalid       = errors.New("token无效")
    ErrUserNotFound       = errors.New("用户不存在")
    ErrPermissionDenied   = errors.New("权限不足")
    ErrLoginFailed        = errors.New("登录失败")
    ErrUserDisabled       = errors.New("用户已禁用")
    ErrUserLocked         = errors.New("用户已锁定")
    ErrSecondAuthRequired = errors.New("需要二级认证")
    ErrCasbinError        = errors.New("Casbin错误")
    ErrRoleNotExists      = errors.New("角色不存在")
    ErrRoleExists         = errors.New("角色已存在")
)
```

## 扩展与集成

Sweet Auth 设计为可扩展的，您可以：

1. **自定义Token生成**：实现自己的Token生成逻辑
2. **自定义存储方式**：实现Storage接口支持其他存储
3. **集成到Web框架**：与Gin、Echo等框架集成
4. **自定义认证事件处理**：通过监听器处理各种认证事件

```go
// 示例：集成到Gin框架
func AuthMiddleware(authManager *auth.Manager) gin.HandlerFunc {
    return func(c *gin.Context) {
        token := c.GetHeader("Authorization")
        if token == "" {
            c.JSON(401, gin.H{"error": "未授权"})
            c.Abort()
            return
        }
        
        // 验证Token
        info, err := authManager.ValidateToken(token)
        if err != nil {
            c.JSON(401, gin.H{"error": err.Error()})
            c.Abort()
            return
        }
        
        // 权限验证
        allowed, _ := authManager.CheckPermission(
            info.UserID, 
            c.Request.URL.Path, 
            c.Request.Method,
        )
        if !allowed {
            c.JSON(403, gin.H{"error": "权限不足"})
            c.Abort()
            return
        }
        
        // 设置用户信息
        c.Set("userID", info.UserID)
        c.Set("username", info.Username)
        
        c.Next()
    }
}
```

## 性能与安全建议

1. **生产环境配置**
   - 使用强密钥作为JWTSecret
   - 适当设置TokenExpire，平衡安全性和用户体验
   - 生产环境推荐使用Redis存储

2. **安全最佳实践**
   - 使用HTTPS传输Token
   - 实现请求签名机制
   - 定期审计权限配置
   - 对敏感操作启用二级认证

3. **性能优化**
   - 合理设置Token缓存
   - 避免过于复杂的权限规则
   - 监控并发连接数

## 后续开发计划

1. 完善单点登录(SSO)功能
2. 添加更多存储适配器
3. 增强权限管理UI组件
4. 提供更多权限模型支持
