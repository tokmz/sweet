# SaToken Go版

SaToken Go版是一个简化版的类似于Java中的[SaToken](https://sa-token.cc/doc.html)的身份认证与权限控制框架，提供了登录认证、单端登录、多端登录、同端互斥登录、七天内免登录、踢人下线和权限管理等功能。

## 功能特性

- **登录认证**：提供基本的登录验证功能
- **多端登录**：支持多种登录模式
  - 单端登录：一个用户只能有一个有效token
  - 多端登录：一个用户可以拥有多个有效token
  - 同端互斥登录：一个用户在同一类型的设备上只能有一个有效token
- **存储模式**：支持多种token存储方式
  - 内存存储：适用于单机部署
  - Redis存储：适用于分布式部署
- **Token特性**：
  - 自动续期：根据用户活跃度自动延长token有效期
  - 踢人下线：支持指定用户或设备强制下线
  - 七天内免登录：支持通过配置长时间token
- **权限控制**：
  - 基于角色的权限控制：RBAC权限模型
  - 权限验证：支持细粒度的权限验证
  - 通配符匹配：支持权限通配符匹配

## 快速入门

### 安装

```bash
go get github.com/your-org/satoken
```

### 基本使用

```go
package main

import (
	"context"
	"fmt"
	"time"
	
	"github.com/your-org/satoken"
)

func main() {
	// 创建SaToken实例，使用默认配置
	sa := satoken.New(satoken.DefaultConfig())
	
	// 登录
	ctx := context.Background()
	userID := int64(1001)
	token, err := sa.Login(ctx, userID, "web", satoken.LoginTypeDefault)
	if err != nil {
		panic(err)
	}
	
	fmt.Println("登录成功，token:", token)
	
	// 检查登录状态
	isLogin, err := sa.CheckLogin(ctx, token)
	if err != nil {
		panic(err)
	}
	fmt.Println("是否已登录:", isLogin)
	
	// 获取用户ID
	id, err := sa.GetUserID(ctx, token)
	if err != nil {
		panic(err)
	}
	fmt.Println("用户ID:", id)
	
	// 退出登录
	err = sa.Logout(ctx, token)
	if err != nil {
		panic(err)
	}
	fmt.Println("已退出登录")
}
```

### 配置详解

```go
// 创建自定义配置
config := satoken.Config{
	StoreType:           satoken.StoreMemory,  // 存储类型
	AccessTokenTimeout:  30 * time.Minute,     // 访问令牌有效期
	RefreshTokenTimeout: 7 * 24 * time.Hour,   // 刷新令牌有效期
	TokenStyle:          "random",             // 令牌风格
	AllowConcurrentLogin: false,               // 是否允许同账号并发登录
	IsShare:             false,                // 是否共享token
	AutoRenew:           true,                 // 是否自动续期
	TokenPrefix:         "satoken:",           // 令牌前缀
}

// 创建SaToken实例
sa := satoken.New(config)
```

### 使用Redis存储

```go
package main

import (
	"context"
	"fmt"
	
	"github.com/redis/go-redis/v9"
	"github.com/your-org/satoken"
)

func main() {
	// 创建Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	
	// 创建配置
	config := satoken.DefaultConfig()
	config.StoreType = satoken.StoreRedis
	config.RedisClient = redisClient
	
	// 创建SaToken实例
	sa := satoken.New(config)
	
	// 现在可以使用sa进行登录认证等操作
	// ...
}
```

## 登录模式示例

### 单端登录

```go
// 单端登录：一个用户只能在一个设备上登录，新登录会踢掉旧登录
token, err := sa.Login(ctx, userID, "web", satoken.LoginTypeSingle)
```

### 多端登录

```go
// 多端登录：允许一个用户在多个设备上同时登录
token, err := sa.Login(ctx, userID, "web", satoken.LoginTypeMulti)
```

### 同端互斥登录

```go
// 同端互斥登录：一个用户在同一类型设备上只能登录一次，新登录会踢掉旧登录
token, err := sa.Login(ctx, userID, "web", satoken.LoginTypeExclusive)
```

## 踢人下线

```go
// 根据token踢人下线
err := sa.KickOut(ctx, token)

// 根据用户ID踢人下线
err := sa.KickOutByUserID(ctx, userID)

// 根据用户ID和设备类型踢人下线
err := sa.KickOutByUserIDAndDevice(ctx, userID, "web")
```

## 获取在线用户

```go
// 获取指定用户的所有有效token
tokens, err := sa.GetActiveTokens(ctx, userID)

// 获取指定用户在特定设备上的所有有效token
tokens, err := sa.GetActiveTokensByDevice(ctx, userID, "web")
```

## 权限管理

### 权限分配

```go
// 为用户分配权限
err := sa.AssignPermissionsToUser(ctx, userID, "system:user:list", "system:user:create")

// 为用户分配角色
err := sa.AssignRolesToUser(ctx, userID, "admin", "editor")

// 为角色分配权限
err := sa.AssignPermissionsToRole(ctx, "editor", "content:article:list", "content:article:edit")
```

### 权限验证

```go
// 检查用户是否拥有指定权限
hasPermission, err := sa.HasPermission(ctx, token, "system:user:create")

// 检查用户是否拥有指定角色
hasRole, err := sa.HasRole(ctx, token, "admin")

// 强制用户拥有指定权限，否则返回错误
err := sa.EnforcePermission(ctx, token, "system:user:create")

// 强制用户拥有指定角色，否则返回错误
err := sa.EnforceRole(ctx, token, "admin")
```

### 获取用户权限信息

```go
// 获取用户所有权限（包括角色权限）
permissions, err := sa.GetAllPermissions(ctx, userID)

// 获取用户所有角色
roles, err := sa.GetAllRoles(ctx, userID)
```

### 通配符权限匹配

SaToken支持使用通配符进行权限匹配，例如：

```go
// 为角色分配包含通配符的权限
sa.AssignPermissionsToRole(ctx, "admin", "system:*") // 系统所有权限

// 验证特定权限时，会自动匹配通配符
// 如果用户有system:*权限，下面的验证会通过
hasPermission, _ := sa.HasPermission(ctx, token, "system:user:create")
```

通配符规则：
- `*` 表示匹配任何内容
- `system:*` 匹配所有以system:开头的权限
- `*:view` 匹配所有以view结尾的权限
- `system:user:*` 匹配所有system:user:下的权限

## 许可证

MIT License 