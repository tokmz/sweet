package repository

import (
	"context"
	"sweet/pkg/errs"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrLogNotFound  = errs.New(8000, "日志不存在")
	ErrInvalidLogID = errs.New(8001, "无效日志ID")
)

// LogRepository 日志仓储接口
type LogRepository interface {
	/* 操作日志 */

	// CreateOperationLog 创建操作日志
	CreateOperationLog(ctx context.Context, log *entity.OperationLog) error
	// DeleteOperationLog 删除操作日志
	DeleteOperationLog(ctx context.Context, ids []int64) error
	// DeleteOperationLogByDate 删除指定日期前的操作日志
	DeleteOperationLogByDate(ctx context.Context, date string) error
	// ClearOperationLog 清空操作日志
	ClearOperationLog(ctx context.Context) error
	// FindOperationLog 查询操作日志
	FindOperationLog(ctx context.Context, id int64) (*entity.OperationLog, error)
	// FindOperationLogList 查询操作日志列表
	FindOperationLogList(ctx context.Context, params *OperationLogListParams) (list []*entity.OperationLog, total int64, err error)

	/* 登录日志 */

	// CreateLoginLog 创建登录日志
	CreateLoginLog(ctx context.Context, log *entity.LoginLog) error
	// DeleteLoginLog 删除登录日志
	DeleteLoginLog(ctx context.Context, ids []int64) error
	// DeleteLoginLogByDate 删除指定日期前的登录日志
	DeleteLoginLogByDate(ctx context.Context, date string) error
	// ClearLoginLog 清空登录日志
	ClearLoginLog(ctx context.Context) error
	// FindLoginLog 查询登录日志
	FindLoginLog(ctx context.Context, id int64) (*entity.LoginLog, error)
	// FindLoginLogList 查询登录日志列表
	FindLoginLogList(ctx context.Context, params *LoginLogListParams) (list []*entity.LoginLog, total int64, err error)
}

// OperationLogListParams 操作日志列表查询参数
type OperationLogListParams struct {
	Title         string // 模块标题
	BusinessType  *int64 // 业务类型
	OperatorType  *int64 // 操作类别
	Status        *int64 // 操作状态
	Username      string // 操作人员
	RequestMethod string // 请求方式
	Method        string // 方法名称
	Start         string // 开始时间
	End           string // 结束时间
	Page          int    // 页码
	Size          int    // 每页数量
}

// LoginLogListParams 登录日志列表查询参数
type LoginLogListParams struct {
	Username string // 用户名
	Status   *int64 // 状态
	Ip       string // IP地址
	Location string // 登录地点
	Browser  string // 浏览器
	Os       string // 操作系统
	Device   string // 设备
	Start    string // 开始时间
	End      string // 结束时间
	Page     int    // 页码
	Size     int    // 每页数量
}

// logRepository 日志仓储实现
type logRepository struct {
	db *gorm.DB
	q  *query.Query
}

/* 操作日志 */

func (l *logRepository) CreateOperationLog(ctx context.Context, log *entity.OperationLog) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) DeleteOperationLog(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) DeleteOperationLogByDate(ctx context.Context, date string) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) ClearOperationLog(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) FindOperationLog(ctx context.Context, id int64) (*entity.OperationLog, error) {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) FindOperationLogList(ctx context.Context, params *OperationLogListParams) (list []*entity.OperationLog, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

/* 登录日志 */

func (l *logRepository) CreateLoginLog(ctx context.Context, log *entity.LoginLog) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) DeleteLoginLog(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) DeleteLoginLogByDate(ctx context.Context, date string) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) ClearLoginLog(ctx context.Context) error {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) FindLoginLog(ctx context.Context, id int64) (*entity.LoginLog, error) {
	//TODO implement me
	panic("implement me")
}

func (l *logRepository) FindLoginLogList(ctx context.Context, params *LoginLogListParams) (list []*entity.LoginLog, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

// NewLogRepository 创建日志仓储
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{
		db: db,
		q:  query.Use(db),
	}
}
