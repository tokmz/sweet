# 数据库包优化方案

## 当前问题

通过对 `pkg/database` 包的分析，发现以下需要优化的问题：

1. **缺少并发安全保护**：
   - 当前实现中没有全局互斥锁保护全局变量
   - 多个 goroutine 同时初始化或访问数据库可能导致竞态条件

2. **错误信息国际化**：
   - 错误信息使用中文硬编码，不利于国际化
   - 应该使用错误常量，便于后续国际化处理

3. **缺少全局默认数据库实例**：
   - 与 config 包类似，应提供全局默认数据库实例
   - 便于在不需要显式传递 DB 实例的场景使用

4. **初始化检查机制不完善**：
   - 缺少初始化状态检查
   - 没有安全的初始化和关闭机制

5. **文档与实际代码不一致**：
   - README 中的示例与实际代码实现有差异
   - 缺少完整的使用示例

## 优化建议

### 1. 添加并发安全保护

在 `vars.go` 中添加全局互斥锁和初始化标志：

```go
// 全局变量
var (
	// DefaultDB 是默认的数据库连接实例
	DefaultDB *gorm.DB

	// rw 是一个用于保护全局变量的读写锁
	rw sync.RWMutex

	// initialized 标记数据库是否已初始化
	initialized bool
)
```

### 2. 错误信息国际化

在 `vars.go` 中定义标准错误常量：

```go
// 错误定义
var (
	// ErrNotInitialized 表示数据库尚未初始化
	ErrNotInitialized = errors.New("database not initialized")

	// ErrInvalidConfig 表示提供的配置参数无效
	ErrInvalidConfig = errors.New("invalid database configuration")

	// ErrConnectionFailed 表示连接数据库失败
	ErrConnectionFailed = errors.New("failed to connect to database")

	// ErrConfigReadWriteSplit 表示配置读写分离失败
	ErrConfigReadWriteSplit = errors.New("failed to configure read-write splitting")

	// ErrGetSQLDB 表示获取底层sql.DB失败
	ErrGetSQLDB = errors.New("failed to get underlying sql.DB")

	// ErrConfigTracing 表示配置链路追踪失败
	ErrConfigTracing = errors.New("failed to configure tracing")
)
```

### 3. 添加全局默认数据库实例

修改 `Init` 函数，使其同时设置全局默认实例：

```go
// Init 初始化数据库连接并设置为全局默认实例
func Init(config Config) (*gorm.DB, error) {
    // 加锁保护并发安全
    rw.Lock()
    defer rw.Unlock()
    
    // 创建数据库连接
    db, err := initDB(config)
    if err != nil {
        return nil, err
    }
    
    // 设置全局默认实例
    DefaultDB = db
    initialized = true
    
    return db, nil
}

// initDB 内部函数，实际初始化数据库连接
func initDB(config Config) (*gorm.DB, error) {
    // 原有的初始化逻辑
    // ...
}
```

### 4. 实现安全的初始化检查机制

添加获取默认数据库实例的函数：

```go
// GetDB 安全地获取默认数据库实例
func GetDB() (*gorm.DB, error) {
    rw.RLock()
    defer rw.RUnlock()
    
    if !initialized || DefaultDB == nil {
        return nil, ErrNotInitialized
    }
    
    return DefaultDB, nil
}
```

### 5. 更新错误处理

将所有中文错误信息替换为英文错误常量：

```go
// 示例：修改连接主库失败的错误
db, err := gorm.Open(mysql.Open(config.Master), gormConfig)
if err != nil {
    return nil, fmt.Errorf("%w: %s", ErrConnectionFailed, err.Error())
}
```

### 6. 更新文档和示例

更新 README.md，确保示例与实际代码一致，并添加全局默认实例的使用示例。

## 实施计划

1. 修改 `vars.go`，添加全局变量和错误常量
2. 修改 `database.go`，实现并发安全的初始化和访问机制
3. 替换所有中文错误信息为英文错误常量
4. 更新 README.md，确保文档与代码一致
5. 添加单元测试，确保修改不会破坏现有功能

## 预期收益

1. 提高并发安全性，避免竞态条件
2. 便于国际化，提高代码可维护性
3. 简化使用方式，通过全局默认实例减少代码重复
4. 提高代码健壮性，通过初始化检查避免空指针异常
5. 文档与代码保持一致，减少使用者的困惑

## GORM Gen 与 TiDB 集成指南

GORM Gen 是GORM官方的代码生成工具，可用于生成类型安全的数据库访问代码。以下是将GORM Gen与TiDB集成的步骤：

### 1. 安装GORM Gen工具

```bash
go get -u gorm.io/gen
```

### 2. 创建生成器配置文件

在项目的`tools/gentool`目录中创建一个配置文件，例如：

```go
// main.go
package main

import (
	"fmt"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gorm"
)

// TiDB数据库连接配置
const dsn = "user:pass@tcp(tidb.example.com:4000)/dbname?charset=utf8mb4&parseTime=True&loc=Local&tls=true"

func main() {
	// 连接数据库
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Fprintf(os.Stderr, "failed to connect database: %v", err)
		return
	}

	// 初始化生成器
	g := gen.NewGenerator(gen.Config{
		// 输出路径
		OutPath: "../../internal/entity/gen",
		// 输出包名
		OutFile: "gen.go",
		// 是否使用模型包装器
		ModelPkgPath: "../../internal/entity/model",
		// 自定义模型可见性
		FieldNullable: true,
		// 字段支持NULL值
		FieldCoverable: true,
		// 使用单数表名映射（TiDB命名习惯）
		Mode: gen.WithoutContext | gen.WithDefaultQuery,
	})

	// 为TiDB设置方言，TiDB兼容MySQL语法
	g.UseDB(db)

	// 生成所有表的模型
	allModels := g.GenerateAllTable()

	// 也可以为特定表生成模型
	// user := g.GenerateModel("users")
	// product := g.GenerateModel("products")

	// 为模型生成查询API
	g.ApplyBasic(allModels...)

	// 执行生成
	g.Execute()
}
```

### 3. 为TiDB配置生成自定义特性

如果需要支持TiDB的特殊功能如AUTO_RANDOM，可以通过下面的配置：

```go
// 定义支持TiDB AUTO_RANDOM功能的模型
user := g.GenerateModel("users", gen.FieldGORMTag("id", "type:bigint;auto_random"))
```

### 4. 运行代码生成

```bash
cd tools/gentool
go run main.go
```

### 5. 使用生成的代码

```go
import (
	"context"
	
	"your-project/internal/entity/gen"
	"your-project/pkg/database"
)

func GetUserByID(id int64) (user *gen.User, err error) {
	// 获取数据库连接
	db, err := database.GetDB()
	if err != nil {
		return nil, err
	}
	
	// 创建查询
	query := gen.Use(db)
	
	// 使用生成的查询方法
	return query.User.Where(query.User.ID.Eq(id)).First()
}
```

### 6. TiDB特有优化配置

当与TiDB一起使用GORM Gen时，以下是一些推荐的优化配置：

1. **批量操作优化** - TiDB对大批量操作有特殊优化：

```go
// 批量插入时设置合适的批量大小(推荐100-1000)
users := make([]*gen.User, 0, 100)
// ... 填充用户数据
err := query.User.CreateInBatches(users, 100)
```

2. **事务处理** - 使用Sweet项目的重试事务机制：

```go
err := database.Transaction(db, ctx, func(tx *gorm.DB) error {
    q := gen.Use(tx)
    // 使用查询对象进行事务操作
    return q.User.Create(&user)
})
```

3. **使用生成的高级查询API** - 充分利用GORM Gen生成的类型安全查询：

```go
// 复杂查询示例
users, err := query.User.
    Select(query.User.ID, query.User.Name, query.User.CreatedAt).
    Where(query.User.Age.Gt(18)).
    Order(query.User.CreatedAt.Desc()).
    Limit(10).
    Find()
```

4. **设置合理的索引** - TiDB中合理设置索引对性能至关重要：

```sql
-- 为常用查询字段创建索引
CREATE INDEX idx_users_phone ON users (phone);
```

通过以上配置，可以让Sweet项目的database包与TiDB及GORM Gen高效集成，同时利用TiDB的分布式特性实现高可用数据库操作。