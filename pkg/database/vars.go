package database

import (
	"errors"
	"sync"

	"gorm.io/gorm"
)

// 全局变量
var (
	// DefaultDB 是默认的数据库连接实例
	DefaultDB *gorm.DB

	// rw 是一个用于保护全局变量的读写锁
	rw sync.RWMutex

	// initialized 标记数据库是否已初始化
	initialized bool
)

// 链路追踪相关常量
const (
	tracerName = "sweet/pkg/database"
	opCreate   = "gorm.Create"
	opQuery    = "gorm.Query"
	opUpdate   = "gorm.Update"
	opDelete   = "gorm.Delete"
	opRawSQL   = "gorm.RawSQL"
)

// 上下文键类型定义
type contextKey string

// 上下文键常量
const spanKey contextKey = "db_span"

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
