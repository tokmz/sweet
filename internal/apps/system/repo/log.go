package repo

import (
	"context"
	"sweet/pkg/errs"
	"time"

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
	// Create 创建日志
	Create(ctx context.Context, log *entity.OperationLog) error
	// Delete 删除日志
	Delete(ctx context.Context, ids []int64) error
	// DeleteByDate 删除指定时间段的日志
	DeleteByDate(ctx context.Context, start, end time.Time) error
	// Clear 清空日志
	Clear(ctx context.Context) error
	// FindOne 查询日志
	FindOne(ctx context.Context, id int64) (*entity.OperationLog, error)
	// FindList 查询日志列表
	FindList(ctx context.Context, params *OperationLogListParams) (list []*entity.OperationLog, total int64, err error)
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

// logRepository 日志仓储实现
type logRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (l *logRepository) Create(ctx context.Context, log *entity.OperationLog) error {
	return l.q.OperationLog.WithContext(ctx).Create(log)
}

func (l *logRepository) Delete(ctx context.Context, ids []int64) error {
	if len(ids) == 0 {
		return nil
	}
	_, err := l.q.OperationLog.WithContext(ctx).Where(l.q.OperationLog.ID.In(ids...)).Delete()
	return err
}

func (l *logRepository) DeleteByDate(ctx context.Context, start, end time.Time) error {
	if start.After(end) {
		start, end = end, start // 确保开始时间不晚于结束时间
	}

	_, err := l.q.OperationLog.WithContext(ctx).
		Where(l.q.OperationLog.CreatedAt.Between(start, end)).
		Delete()
	return err
}

func (l *logRepository) Clear(ctx context.Context) error {
	// 使用原生SQL执行truncate操作，效率更高
	return l.db.WithContext(ctx).Exec("TRUNCATE TABLE " + entity.TableNameOperationLog).Error
}

func (l *logRepository) FindOne(ctx context.Context, id int64) (*entity.OperationLog, error) {
	if id <= 0 {
		return nil, ErrInvalidLogID
	}

	log, err := l.q.OperationLog.WithContext(ctx).Where(l.q.OperationLog.ID.Eq(id)).First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, ErrLogNotFound
		}
		return nil, err
	}
	return log, nil
}

func (l *logRepository) FindList(ctx context.Context, params *OperationLogListParams) (list []*entity.OperationLog, total int64, err error) {
	// 构建查询条件
	query := l.q.OperationLog.WithContext(ctx)

	// 根据参数进行筛选
	if params.Title != "" {
		query = query.Where(l.q.OperationLog.Title.Like("%" + params.Title + "%"))
	}

	if params.BusinessType != nil {
		query = query.Where(l.q.OperationLog.BusinessType.Eq(*params.BusinessType))
	}

	if params.OperatorType != nil {
		query = query.Where(l.q.OperationLog.OperatorType.Eq(*params.OperatorType))
	}

	if params.Status != nil {
		query = query.Where(l.q.OperationLog.Status.Eq(*params.Status))
	}

	if params.Username != "" {
		query = query.Where(l.q.OperationLog.Username.Like("%" + params.Username + "%"))
	}

	if params.RequestMethod != "" {
		query = query.Where(l.q.OperationLog.RequestMethod.Eq(params.RequestMethod))
	}

	if params.Method != "" {
		query = query.Where(l.q.OperationLog.Method.Like("%" + params.Method + "%"))
	}

	// 处理时间范围查询
	if params.Start != "" && params.End != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", params.Start)
		if err == nil {
			endTime, err := time.Parse("2006-01-02 15:04:05", params.End)
			if err == nil {
				query = query.Where(l.q.OperationLog.CreatedAt.Between(startTime, endTime))
			}
		}
	} else if params.Start != "" {
		startTime, err := time.Parse("2006-01-02 15:04:05", params.Start)
		if err == nil {
			query = query.Where(l.q.OperationLog.CreatedAt.Gte(startTime))
		}
	} else if params.End != "" {
		endTime, err := time.Parse("2006-01-02 15:04:05", params.End)
		if err == nil {
			query = query.Where(l.q.OperationLog.CreatedAt.Lte(endTime))
		}
	}

	// 获取总数
	count, err := query.Count()
	if err != nil {
		return nil, 0, err
	}

	// 分页并排序获取数据
	if params.Page <= 0 {
		params.Page = 1
	}
	if params.Size <= 0 {
		params.Size = 10
	}

	offset := (params.Page - 1) * params.Size
	result, err := query.
		Order(l.q.OperationLog.CreatedAt.Desc()).
		Offset(offset).
		Limit(params.Size).
		Find()

	return result, count, err
}

// NewLogRepository 创建日志仓储
func NewLogRepository(db *gorm.DB) LogRepository {
	return &logRepository{
		db: db,
		q:  query.Use(db),
	}
}
