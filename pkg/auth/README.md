# Auth 认证与授权框架

Auth 是一个基于 JWT 和 Casbin 的认证与授权框架，提供了登录认证、多端登录控制、权限管理等功能。

## 功能特性

- **登录认证**：使用 JWT 实现安全的登录验证功能
- **多端登录**：支持多种登录模式
  - 单端登录：一个用户只能有一个有效 token
  - 多端登录：一个用户可以拥有多个有效 token
  - 同端互斥登录：一个用户在同一类型的设备上只能有一个有效 token
- **存储模式**：支持多种 token 存储方式
  - 内存存储：适用于单机部署
  - Redis 存储：适用于分布式部署（可扩展）
- **Token 特性**：
  - 基于 JWT：标准的 JWT 实现，支持签名验证
  - 自动续期：根据用户活跃度自动延长 token 有效期
  - 踢人下线：支持指定用户或设备强制下线
- **权限控制**：
  - 基于 Casbin：使用强大的 Casbin 库实现 RBAC 权限模型
  - 角色管理：支持角色分配、角色撤销等功能
  - 权限管理：细粒度的权限控制，支持为用户和角色分配权限

## 快速入门

### 安装

```bash
go get github.com/your-org/auth
```

### 基本使用

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/your-org/auth"
)

func main() {
	// 创建自定义配置
	config := auth.Config{
		StoreType:            auth.StoreMemory,
		AccessTokenTimeout:   30 * time.Minute,
		RefreshTokenTimeout:  7 * 24 * time.Hour,
		TokenIssuer:          "myapp",
		TokenAudience:        "users",
		JWTSigningKey:        []byte("your-secret-key"),  // 生产环境请使用安全的密钥管理
		AllowConcurrentLogin: false,
		AutoRenew:            true,
		TokenPrefix:          "auth:",
	}
	
	// 创建认证实例
	a := auth.New(config)
	
	// 登录
	ctx := context.Background()
	userID := int64(1001)
	token, err := a.Login(ctx, userID, "web", auth.LoginTypeDefault)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("登录成功，token:", token)
	
	// 检查登录状态
	isLogin, err := a.CheckLogin(ctx, token)
	if err != nil {
		panic(err)
	}
	fmt.Println("是否已登录:", isLogin)
	
	// 获取用户ID
	id, err := a.GetUserID(ctx, token)
	if err != nil {
		panic(err)
	}
	fmt.Println("用户ID:", id)
	
	// 为用户分配角色
	err = a.AssignRoleToUser(ctx, userID, "admin")
	if err != nil {
		panic(err)
	}
	
	// 为角色添加权限
	err = a.AddPermissionForRole(ctx, "admin", "user", "create")
	if err != nil {
		panic(err)
	}
	
	// 检查用户是否拥有权限
	hasPermission, err := a.HasPermission(ctx, token, "user", "create")
	if err != nil {
		panic(err)
	}
	fmt.Println("是否有创建用户权限:", hasPermission)
	
	// 退出登录
	err = a.Logout(ctx, token)
	if err != nil {
		panic(err)
	}
	fmt.Println("已退出登录")
}
```

## 登录模式示例

### 单端登录

```go
// 单端登录：一个用户只能在一个设备上登录，新登录会踢掉旧登录
token, err := a.Login(ctx, userID, "web", auth.LoginTypeSingle)
```

### 多端登录

```go
// 多端登录：允许一个用户在多个设备上同时登录
token, err := a.Login(ctx, userID, "web", auth.LoginTypeMulti)
```

### 同端互斥登录

```go
// 同端互斥登录：一个用户在同一类型设备上只能登录一次，新登录会踢掉旧登录
token, err := a.Login(ctx, userID, "web", auth.LoginTypeExclusive)
```

## 踢人下线

```go
// 根据 token 踢人下线
err := a.KickOut(ctx, token)

// 根据用户 ID 踢人下线
err := a.KickOutByUserID(ctx, userID)

// 根据用户 ID 和设备类型踢人下线
err := a.KickOutByUserIDAndDevice(ctx, userID, "web")
```

## 获取在线用户

```go
// 获取指定用户的所有有效 token
tokens, err := a.GetActiveTokens(ctx, userID)

// 获取指定用户在特定设备上的所有有效 token
tokens, err := a.GetActiveTokensByDevice(ctx, userID, "web")
```

## 权限管理

### 角色管理

```go
// 为用户分配角色
err := a.AssignRoleToUser(ctx, userID, "admin")

// 从用户撤销角色
err := a.RevokeRoleFromUser(ctx, userID, "admin")

// 获取用户的所有角色
roles, err := a.GetRoles(ctx, token)

// 获取指定用户的所有角色
roles, err := a.GetUserRoles(ctx, userID)

// 检查用户是否拥有指定角色
hasRole, err := a.HasRole(ctx, token, "admin")

// 强制用户拥有指定角色，否则返回错误
err := a.EnforceRole(ctx, token, "admin")
```

### 权限管理

```go
// 为角色添加权限
err := a.AddPermissionForRole(ctx, "admin", "user", "create")

// 从角色移除权限
err := a.RemovePermissionFromRole(ctx, "admin", "user", "delete")

// 为用户添加直接权限（不通过角色）
err := a.AddPermissionForUser(ctx, userID, "article", "edit")

// 从用户移除直接权限
err := a.RemovePermissionFromUser(ctx, userID, "article", "edit")

// 获取用户的所有权限
permissions, err := a.GetPermissions(ctx, token)

// 获取指定用户的所有权限
permissions, err := a.GetUserPermissions(ctx, userID)

// 检查用户是否拥有指定权限
hasPermission, err := a.HasPermission(ctx, token, "user", "create")

// 强制用户拥有指定权限，否则返回错误
err := a.EnforcePermission(ctx, token, "user", "create")
```

## Casbin 权限模型

框架使用 Casbin 的 RBAC 模型，默认规则如下：

```
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")
```

这允许：
- 用户可以继承角色拥有的权限
- 通配符匹配：`role:admin, *, *` 表示管理员角色拥有所有对象的所有操作权限

## Redis 存储示例

```go
package main

import (
	"context"
	
	"github.com/redis/go-redis/v9"
	"github.com/your-org/auth"
)

func main() {
	// 创建 Redis 客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	
	// 创建配置
	config := auth.DefaultConfig()
	config.StoreType = auth.StoreRedis
	config.RedisClient = redisClient
	config.JWTSigningKey = []byte("your-secret-key")
	
	// 创建认证实例
	a := auth.New(config)
	
	// 现在可以使用 a 进行认证和授权操作
	// ...
}
```

## 许可证

MIT License 