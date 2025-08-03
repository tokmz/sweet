package global

import (
	"sweet/internal/models/query"
	"sweet/pkg/cache"
	"sweet/pkg/database"
	"sweet/pkg/logger"
)

var (
	// DBClient 数据库客户端
	DBClient *database.Client
	// Query 数据库查询
	Query *query.Query
	// CacheClient 缓存客户端
	CacheClient *cache.Client
	// Logger 日志客户端
	Logger logger.Logger
)
