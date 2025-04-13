package auth_test

import (
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
	"sweet/pkg/auth"
)

func ExampleWithRedisClient() {
	// 方式1：使用已有的Redis客户端
	redisClient := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})

	authManager, err := auth.NewManager(auth.Config{
		JWTSecret:   "your-jwt-secret",
		TokenExpire: 24 * time.Hour,
	}, auth.WithRedisClient(redisClient, "auth:"))
	if err != nil {
		panic(err)
	}

	// 使用authManager...
	fmt.Println("Auth manager created with existing Redis client")

	// 关闭管理器（不会关闭传入的Redis客户端）
	authManager.Close()

	// 关闭Redis客户端（需要单独关闭）
	redisClient.Close()
	// Output: Auth manager created with existing Redis client
}

func ExampleWithRedisStorage() {
	// 方式2：通过连接参数使用Redis
	authManager, err := auth.NewManager(auth.Config{
		JWTSecret:   "your-jwt-secret",
		TokenExpire: 24 * time.Hour,
	}, auth.WithRedisStorage("localhost:6379", "", 0, "auth:"))
	if err != nil {
		fmt.Println("Failed to create auth manager:", err)
		return
	}

	// 使用authManager...
	fmt.Println("Auth manager created with Redis storage")

	// 关闭管理器（会自动关闭Redis客户端）
	authManager.Close()
	// Output: Auth manager created with Redis storage
}

func ExampleWithMemoryStorage() {
	// 方式3：使用内存存储（默认）
	authManager, err := auth.NewManager(auth.Config{
		JWTSecret:   "your-jwt-secret",
		TokenExpire: 24 * time.Hour,
	})
	if err != nil {
		panic(err)
	}

	// 使用authManager...
	fmt.Println("Auth manager created with default memory storage")

	// 或者显式指定
	authManager2, err := auth.NewManager(auth.Config{
		JWTSecret:   "your-jwt-secret",
		TokenExpire: 24 * time.Hour,
	}, auth.WithMemoryStorage())
	if err != nil {
		panic(err)
	}

	// 使用authManager2...
	fmt.Println("Auth manager created with explicit memory storage")

	// 关闭管理器
	authManager.Close()
	authManager2.Close()
	// Output:
	// Auth manager created with default memory storage
	// Auth manager created with explicit memory storage
}
