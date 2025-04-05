package database

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
