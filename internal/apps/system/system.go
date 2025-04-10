package system

import (
	"fmt"
	"sweet/internal/apps/system/types"
	"sweet/internal/apps/system/types/query"
	"sweet/pkg/cache"
	"sweet/pkg/config"
	"sweet/pkg/database"

	"gorm.io/gorm"
)

// DB 全局数据库实例
var DB *gorm.DB

// CacheClient 全局缓存客户端
var CacheClient cache.Client

func Init() {
	if err := config.InitConfig(&config.Config{
		Name: "config",
		Type: "yaml",
		Path: []string{"/Users/aikzy/Desktop/go/sweet/internal/apps/system"},
	}); err != nil {
		panic("初始化system模块配置失败: " + err.Error())
	}

	cfg := types.Config{}
	if err := config.LoadConfig(&cfg); err != nil {
		panic("加载system模块配置失败: " + err.Error())
	} else {
		fmt.Printf("system模块配置加载成功\n")
	}

	fmt.Printf("正在启动%s模块\n介绍：%s\n版本:%s\n", cfg.App.Name, cfg.App.Intro, cfg.App.Version)

	dbCfg := cfg.Database
	var err error
	DB, err = database.Init(database.Config{
		Master: dbCfg.Master,
		Slaves: dbCfg.Slaves,
		Pool: database.ConnPool{
			MaxIdleConns:    dbCfg.Pool.MaxIdleConns,    // 最大空闲连接数
			MaxOpenConns:    dbCfg.Pool.MaxOpenConns,    // 最大打开连接数
			ConnMaxLifetime: dbCfg.Pool.ConnMaxLifetime, // 连接最大生命周期
		},
		EnableLog:     dbCfg.EnableLog,     // 启用SQL日志
		EnableTrace:   dbCfg.EnableTrace,   // 启用链路追踪
		SlowThreshold: dbCfg.SlowThreshold, // 慢查询阈值
		DriverType:    dbCfg.DriverType,    // 数据库驱动类型
		QueryTimeout:  dbCfg.QueryTimeout,  // 查询超时
		MaxRetries:    dbCfg.MaxRetries,    // 最大重试次数
		RetryDelay:    dbCfg.RetryDelay,    // 重试延迟（毫秒）
	})
	if err != nil {
		panic(cfg.App.Name + "模块连接数据库失败：" + err.Error())
	} else {
		fmt.Printf("%s模块连接数据库成功\n", cfg.App.Name)
	}

	// 设置query层
	query.SetDefault(DB)

	// 初始化缓存客户端
	cacheCfg := cfg.Cache
	CacheClient, err = cache.NewClient(cache.Config{
		Mode: cache.Mode(cacheCfg.Mode),
		Single: cache.SingleConfig{
			Addr: cacheCfg.Single.Addr,
		},
		Cluster: cache.ClusterConfig{
			Addrs: cacheCfg.Cluster.Addrs,
		},
		Sentinel: cache.SentinelConfig{
			MasterName: cacheCfg.Sentinel.MasterName,
			Addrs:      cacheCfg.Sentinel.Addrs,
		},
		Username:        cacheCfg.Username,
		Password:        cacheCfg.Password,
		DB:              cacheCfg.DB,
		PoolSize:        cacheCfg.PoolSize,
		MinIdleConns:    cacheCfg.MinIdleConns,
		IdleTimeout:     cacheCfg.IdleTimeout,
		ConnTimeout:     cacheCfg.ConnTimeout,
		ReadTimeout:     cacheCfg.ReadTimeout,
		WriteTimeout:    cacheCfg.WriteTimeout,
		ExecTimeout:     cacheCfg.ExecTimeout,
		MaxRetries:      cacheCfg.MaxRetries,
		RetryDelay:      cacheCfg.RetryDelay,
		MinRetryBackoff: cacheCfg.MinRetryBackoff,
		MaxRetryBackoff: cacheCfg.MaxRetryBackoff,
		EnableTrace:     cacheCfg.EnableTrace,
		EnableReadWrite: cacheCfg.EnableReadWrite,
	})
	if err != nil {
		panic("system模块连接缓存服务失败：" + err.Error())
	} else {
		fmt.Printf("%s模块连接缓存服务成功\n", cfg.App.Name)
	}

	fmt.Printf("%s模块启动成功\n", cfg.App.Name)
}
