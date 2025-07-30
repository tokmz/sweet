package database

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"gorm.io/gorm/utils"
)

// CustomLogger 自定义日志记录器
type CustomLogger struct {
	logger.Config
	slowQueryConfig SlowQueryConfig
	logger          *log.Logger
	// 日志格式字符串
	infoStr, warnStr, errStr            string
	traceStr, traceErrStr, traceWarnStr string
}

// NewCustomLogger 创建自定义日志记录器
func NewCustomLogger(config LogConfig, slowQueryConfig SlowQueryConfig) logger.Interface {
	customLogger := &CustomLogger{
		Config: logger.Config{
			SlowThreshold:             slowQueryConfig.Threshold,
			LogLevel:                  config.Level,
			IgnoreRecordNotFoundError: config.IgnoreRecordNotFoundError,
			ParameterizedQueries:      config.ParameterizedQueries,
			Colorful:                  config.Colorful,
		},
		slowQueryConfig: slowQueryConfig,
		logger:          log.New(os.Stdout, "\r\n", log.LstdFlags),
	}
	
	// 初始化日志格式字符串
	customLogger.initLogStrings()
	return customLogger
}

// LogMode 设置日志模式
func (l *CustomLogger) LogMode(level logger.LogLevel) logger.Interface {
	newLogger := *l
	newLogger.LogLevel = level
	return &newLogger
}

// Info 输出信息日志
func (l *CustomLogger) Info(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Info {
		l.logger.Printf(l.infoStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Warn 输出警告日志
func (l *CustomLogger) Warn(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Warn {
		l.logger.Printf(l.warnStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Error 输出错误日志
func (l *CustomLogger) Error(ctx context.Context, msg string, data ...interface{}) {
	if l.LogLevel >= logger.Error {
		l.logger.Printf(l.errStr+msg, append([]interface{}{utils.FileWithLineNum()}, data...)...)
	}
}

// Trace 输出SQL追踪日志
func (l *CustomLogger) Trace(ctx context.Context, begin time.Time, fc func() (sql string, rowsAffected int64), err error) {
	if l.LogLevel <= logger.Silent {
		return
	}

	elapsed := time.Since(begin)
	sql, rows := fc()

	// 处理慢查询
	if l.slowQueryConfig.Enabled && elapsed >= l.slowQueryConfig.Threshold {
		l.logSlowQuery(elapsed, sql, rows, err)
	}

	switch {
	case err != nil && l.LogLevel >= logger.Error && (!errors.Is(err, gorm.ErrRecordNotFound) || !l.IgnoreRecordNotFoundError):
		l.logger.Printf(l.traceErrStr, utils.FileWithLineNum(), err, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case elapsed > l.SlowThreshold && l.SlowThreshold != 0 && l.LogLevel >= logger.Warn:
		l.logger.Printf(l.traceWarnStr, utils.FileWithLineNum(), float64(elapsed.Nanoseconds())/1e6, rows, sql)
	case l.LogLevel == logger.Info:
		l.logger.Printf(l.traceStr, float64(elapsed.Nanoseconds())/1e6, rows, sql)
	}
}

// logSlowQuery 记录慢查询日志
func (l *CustomLogger) logSlowQuery(elapsed time.Duration, sql string, rows int64, err error) {
	slowQueryLog := fmt.Sprintf("[SLOW QUERY] [%.3fms] [rows:%d] %s", 
		float64(elapsed.Nanoseconds())/1e6, rows, sql)
	
	if err != nil {
		slowQueryLog += fmt.Sprintf(" [ERROR: %v]", err)
	}
	
	l.logger.Printf("\033[31m%s\033[0m\n%s\n", utils.FileWithLineNum(), slowQueryLog)
}

// 日志格式字符串
var (
	infoStr      = "\033[32m[info] \033[0m"
	warnStr      = "\033[33m[warn] \033[0m"
	errStr       = "\033[31m[error] \033[0m"
	traceStr     = "\033[35m[%.3fms] \033[34m[rows:%d]\033[0m %s"
	traceWarnStr = "\033[33m%s\n\033[33m[%.3fms] \033[34m[rows:%d]\033[0m %s"
	traceErrStr  = "\033[31m%s\n[ERROR: %v] \033[31m[%.3fms] \033[34m[rows:%d]\033[0m %s\033[0m"
)

// initLogStrings 初始化日志格式字符串
func (l *CustomLogger) initLogStrings() {
	if l.Colorful {
		l.infoStr = infoStr
		l.warnStr = warnStr
		l.errStr = errStr
		l.traceStr = traceStr
		l.traceWarnStr = traceWarnStr
		l.traceErrStr = traceErrStr
	} else {
		l.infoStr = "[info] "
		l.warnStr = "[warn] "
		l.errStr = "[error] "
		l.traceStr = "[%.3fms] [rows:%d] %s"
		l.traceWarnStr = "%s\n[%.3fms] [rows:%d] %s"
		l.traceErrStr = "%s\n[ERROR: %v] [%.3fms] [rows:%d] %s"
	}
}