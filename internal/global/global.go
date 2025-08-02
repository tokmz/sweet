package global

import (
	"sweet/pkg/cache"
	"sweet/pkg/database"
	"sweet/pkg/logger"
)

var (
	// DBClient 数据库客户端
	DBClient *database.Client
	// CacheClient 缓存客户端
	CacheClient *cache.Client
	// Logger 日志客户端
	Logger *logger.Logger
)
