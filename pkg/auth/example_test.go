package auth_test

import (
	"fmt"
	"time"

	"sweet/pkg/auth"
)

func Example() {
	// 创建认证管理器
	authManager, err := auth.NewManager(auth.Config{
		JWTSecret:      "your-jwt-secret",
		TokenExpire:    24 * time.Hour,
		RefreshExpire:  7 * 24 * time.Hour,
		TokenStyle:     auth.TokenStyleJWT,
		LoginModel:     auth.LoginModelMulti,
		EnableRedis:    false, // 使用内存存储
		AutoRenew:      true,
		RenewThreshold: 0.5,
	})
	if err != nil {
		panic(err)
	}

	// 初始化Casbin
	enforcer, err := authManager.InitCasbin(auth.CasbinConfig{
		Model: `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`,
		Adapter: "file",
		DSN:     "casbin_policy.csv",
	})
	if err != nil {
		panic(err)
	}

	// 添加角色和权限
	authManager.AddPermissionForRole("admin", "/api/users", "GET")
	authManager.AddPermissionForRole("admin", "/api/users", "POST")
	authManager.AssignRole(123, "admin")

	// 登录
	token, refreshToken, err := authManager.Login(auth.LoginParams{
		UserID:    123,
		Username:  "admin",
		Device:    auth.DeviceTypeWeb,
		IP:        "127.0.0.1",
		UserAgent: "Mozilla/5.0",
		RememberMe: true,
		ExtraData: map[string]interface{}{
			"role": "admin",
		},
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %s\n", token)
	fmt.Printf("Refresh Token: %s\n", refreshToken)

	// 验证Token
	tokenInfo, err := authManager.ValidateToken(token)
	if err != nil {
		panic(err)
	}
	fmt.Printf("User ID: %d\n", tokenInfo.UserID)
	fmt.Printf("Username: %s\n", tokenInfo.Username)

	// 检查权限
	allowed, err := authManager.CheckPermission(123, "/api/users", "GET")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Has permission: %v\n", allowed)

	// 启用二级认证
	err = authManager.EnableSecondAuth(token, 30*time.Minute)
	if err != nil {
		panic(err)
	}

	// 检查二级认证
	valid, err := authManager.CheckSecondAuth(token)
	if err != nil {
		panic(err)
	}
	fmt.Printf("Second auth valid: %v\n", valid)

	// 获取在线用户
	sessions, err := authManager.GetOnlineUsers()
	if err != nil {
		panic(err)
	}
	fmt.Printf("Online users: %d\n", len(sessions))

	// 注销
	err = authManager.Logout(token)
	if err != nil {
		panic(err)
	}
	fmt.Println("Logged out")

	// 关闭管理器
	authManager.Close()
}

func ExampleManager_Login() {
	// 创建认证管理器
	authManager, _ := auth.NewManager(auth.Config{
		JWTSecret:   "your-jwt-secret",
		TokenExpire: 24 * time.Hour,
		TokenStyle:  auth.TokenStyleJWT,
		LoginModel:  auth.LoginModelMulti,
	})

	// 登录
	token, refreshToken, err := authManager.Login(auth.LoginParams{
		UserID:    123,
		Username:  "admin",
		Device:    auth.DeviceTypeWeb,
		RememberMe: true,
	})
	if err != nil {
		panic(err)
	}
	fmt.Printf("Token: %s\n", token)
	fmt.Printf("Refresh Token: %s\n", refreshToken)
}

func ExampleManager_CheckPermission() {
	// 创建认证管理器
	authManager, _ := auth.NewManager(auth.Config{
		JWTSecret:   "your-jwt-secret",
		TokenExpire: 24 * time.Hour,
	})

	// 初始化Casbin
	authManager.InitCasbin(auth.CasbinConfig{
		Model: `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
`,
		Adapter: "file",
		DSN:     "casbin_policy.csv",
	})

	// 添加角色和权限
	authManager.AddPermissionForRole("admin", "/api/users", "GET")
	authManager.AssignRole(123, "admin")

	// 检查权限
	allowed, err := authManager.CheckPermission(123, "/api/users", "GET")
	if err != nil {
		panic(err)
	}
	fmt.Printf("Has permission: %v\n", allowed)
}

func ExampleManager_KickoutByUserID() {
	// 创建认证管理器
	authManager, _ := auth.NewManager(auth.Config{
		JWTSecret:   "your-jwt-secret",
		TokenExpire: 24 * time.Hour,
	})

	// 添加事件监听器
	authManager.AddListener(func(event auth.AuthEvent) {
		if event.Type == auth.EventKickout {
			fmt.Printf("User %d kicked out\n", event.UserID)
		}
	})

	// 踢人下线
	err := authManager.KickoutByUserID(123)
	if err != nil {
		panic(err)
	}
}
