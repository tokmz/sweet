package database

import (
	"context"
	"fmt"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// User 示例用户模型
type User struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	Name      string    `gorm:"size:100;not null" json:"name"`
	Email     string    `gorm:"size:100;uniqueIndex" json:"email"`
	Age       int       `json:"age"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

// ExampleUsage 展示数据库包的使用方法
func ExampleUsage() {
	ctx := context.Background()

	// 1. 创建数据库配置
	config := &Config{
		Master: "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
		Slaves: []string{
			"user:password@tcp(slave1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			"user:password@tcp(slave2:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Pool: PoolConfig{
			MaxIdleConns:    10,
			MaxOpenConns:    100,
			ConnMaxLifetime: time.Hour,
			ConnMaxIdleTime: time.Minute * 30,
		},
		Log: LogConfig{
			Level:                     logger.Info,
			Colorful:                  true,
			IgnoreRecordNotFoundError: false,
			ParameterizedQueries:      false,
		},
		SlowQuery: SlowQueryConfig{
			Enabled:   true,
			Threshold: time.Millisecond * 200,
		},
		Tracing: TracingConfig{
			Enabled:      true,
			ServiceName:  "user-service",
			RecordSQL:    true,
			RecordParams: false,
		},
	}

	// 2. 创建数据库客户端
	client, err := NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create database client: %v", err)
	}
	defer client.Close()

	// 3. 测试连接
	if err = client.Ping(ctx); err != nil {
		log.Fatalf("Failed to ping database: %v", err)
	}
	fmt.Println("Database connection successful!")

	// 4. 健康检查
	if err = client.HealthCheck(ctx); err != nil {
		log.Printf("Health check failed: %v", err)
	} else {
		fmt.Println("Database health check passed!")
	}

	// 5. 获取连接统计信息
	stats, err := client.Stats()
	if err != nil {
		log.Printf("Failed to get stats: %v", err)
	} else {
		fmt.Printf("Database stats: %+v\n", stats)
	}

	// 6. 自动迁移表结构
	if err = client.DB().AutoMigrate(&User{}); err != nil {
		log.Fatalf("Failed to migrate: %v", err)
	}

	// 7. 写操作示例（使用主库）
	fmt.Println("\n=== 写操作示例 ===")
	user := &User{
		Name:  "张三",
		Email: "zhangsan@example.com",
		Age:   25,
	}

	// 创建用户（自动使用主库）
	if err = client.DB().WithContext(ctx).Create(user).Error; err != nil {
		log.Printf("Failed to create user: %v", err)
	} else {
		fmt.Printf("Created user: %+v\n", user)
	}

	// 强制使用主库进行更新
	if err = client.Master().WithContext(ctx).Model(user).Update("age", 26).Error; err != nil {
		log.Printf("Failed to update user: %v", err)
	} else {
		fmt.Println("Updated user age to 26")
	}

	// 8. 读操作示例（使用从库）
	fmt.Println("\n=== 读操作示例 ===")
	var users []User

	// 查询所有用户（自动使用从库）
	if err = client.DB().WithContext(ctx).Find(&users).Error; err != nil {
		log.Printf("Failed to find users: %v", err)
	} else {
		fmt.Printf("Found %d users\n", len(users))
	}

	// 强制使用从库进行查询
	var user2 User
	if err = client.Slave().WithContext(ctx).First(&user2, user.ID).Error; err != nil {
		log.Printf("Failed to find user: %v", err)
	} else {
		fmt.Printf("Found user from slave: %+v\n", user2)
	}

	// 9. 事务示例
	fmt.Println("\n=== 事务示例 ===")
	err = client.Transaction(ctx, func(tx *gorm.DB) error {
		// 在事务中创建多个用户
		users := []User{
			{Name: "李四", Email: "lisi@example.com", Age: 30},
			{Name: "王五", Email: "wangwu@example.com", Age: 28},
		}

		for _, u := range users {
			if err = tx.Create(&u).Error; err != nil {
				return err
			}
		}

		fmt.Println("Transaction completed successfully")
		return nil
	})

	if err != nil {
		log.Printf("Transaction failed: %v", err)
	}

	// 10. 慢查询示例
	fmt.Println("\n=== 慢查询示例 ===")
	// 执行一个可能较慢的查询
	var count int64
	if err = client.DB().WithContext(ctx).Model(&User{}).Where("age > ?", 20).Count(&count).Error; err != nil {
		log.Printf("Failed to count users: %v", err)
	} else {
		fmt.Printf("Users count: %d\n", count)
	}

	// 11. 动态设置日志级别
	fmt.Println("\n=== 动态日志级别示例 ===")
	if err = client.SetLogLevel("warn"); err != nil {
		log.Printf("Failed to set log level: %v", err)
	} else {
		fmt.Println("Log level set to warn")
	}

	// 12. 原生SQL示例
	fmt.Println("\n=== 原生SQL示例 ===")
	var result struct {
		TotalUsers int     `json:"total_users"`
		AvgAge     float64 `json:"avg_age"`
	}

	if err = client.DB().WithContext(ctx).Raw("SELECT COUNT(*) as total_users, AVG(age) as avg_age FROM users").Scan(&result).Error; err != nil {
		log.Printf("Failed to execute raw SQL: %v", err)
	} else {
		fmt.Printf("Raw SQL result: %+v\n", result)
	}

	// 13. 获取慢查询日志（如果数据库支持）
	fmt.Println("\n=== 慢查询日志示例 ===")
	slowQueries, err := client.GetSlowQueries(ctx, 10)
	if err != nil {
		log.Printf("Failed to get slow queries: %v", err)
	} else {
		fmt.Printf("Found %d slow queries\n", len(slowQueries))
		for i, query := range slowQueries {
			fmt.Printf("Slow query %d: %+v\n", i+1, query)
		}
	}
}

// ExampleMinimalUsage 最小化使用示例
func ExampleMinimalUsage() {
	// 使用默认配置
	config := DefaultConfig()
	config.Master = "user:password@tcp(localhost:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local"

	client, err := NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create database client: %v", err)
	}
	defer client.Close()

	// 简单的数据库操作
	ctx := context.Background()
	var users []User
	if err := client.DB().WithContext(ctx).Find(&users).Error; err != nil {
		log.Printf("Failed to find users: %v", err)
	}

	fmt.Printf("Found %d users\n", len(users))
}

// ExampleAdvancedConfiguration 高级配置示例
func ExampleAdvancedConfiguration() {
	config := &Config{
		Master: "user:password@tcp(master:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
		Slaves: []string{
			"user:password@tcp(slave1:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
			"user:password@tcp(slave2:3306)/dbname?charset=utf8mb4&parseTime=True&loc=Local",
		},
		Pool: PoolConfig{
			MaxIdleConns:    20,
			MaxOpenConns:    200,
			ConnMaxLifetime: time.Hour * 2,
			ConnMaxIdleTime: time.Hour,
		},
		Log: LogConfig{
			Level:                     logger.Warn, // 只记录警告和错误
			Colorful:                  false,       // 生产环境关闭颜色
			IgnoreRecordNotFoundError: true,        // 忽略记录未找到错误
			ParameterizedQueries:      true,        // 使用参数化查询
		},
		SlowQuery: SlowQueryConfig{
			Enabled:   true,
			Threshold: time.Millisecond * 100, // 100ms慢查询阈值
		},
		Tracing: TracingConfig{
			Enabled:      true,
			ServiceName:  "production-service",
			RecordSQL:    false, // 生产环境可能不记录SQL
			RecordParams: false, // 生产环境不记录参数
		},
	}

	client, err := NewClient(config)
	if err != nil {
		log.Fatalf("Failed to create database client: %v", err)
	}
	defer client.Close()

	fmt.Println("Advanced configuration database client created successfully")
}
