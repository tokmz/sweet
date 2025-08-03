package basic

import (
	"context"
	"errors"
	"sweet/internal/global"
	"sweet/internal/models"
	basicDto "sweet/internal/models/dto/basic"
	"sweet/internal/models/entity"
	"sweet/pkg/errs"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

// OperationService 操作日志服务实现
type OperationService struct{}

// CreateOperationLog 创建操作日志
func (s *OperationService) CreateOperationLog(ctx context.Context, req *basicDto.CreateOperationLogReq) error {
	if err := global.Query.SysOperationLog.WithContext(ctx).Create(&entity.SysOperationLog{
		UserID:        req.UserID,
		Username:      req.Username,
		Module:        req.Module,
		Operation:     req.Operation,
		Method:        req.Method,
		URL:           req.URL,
		IP:            req.IP,
		Location:      req.Location,
		UserAgent:     req.UserAgent,
		RequestParams: req.Params,
		ResponseData:  req.Result,
		Status:        req.Status,
		ErrorMsg:      req.ErrorMsg,
		CostTime:      req.Duration,
	}); err != nil {
		global.Logger.Error(
			"创建操作日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

// DeleteOperationLog 删除操作日志
func (s *OperationService) DeleteOperationLog(ctx context.Context, req *basicDto.DeleteOperationLogReq) error {
	dao := global.Query.SysOperationLog
	do := dao.WithContext(ctx)
	if req.Uid != 0 {
		do = do.Where(dao.UserID.Eq(req.Uid))
	}

	if _, err := do.Where(dao.ID.In(req.Ids...)).Delete(); err != nil {
		global.Logger.Error(
			"删除操作日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

// ClearAllOperationLog 清空所有操作日志
func (s *OperationService) ClearAllOperationLog(ctx context.Context) error {
	if _, err := global.Query.SysOperationLog.WithContext(ctx).Delete(); err != nil {
		global.Logger.Error(
			"清空所有操作日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

// ClearOperationLog 清空指定用户的操作日志
func (s *OperationService) ClearOperationLog(ctx context.Context, uid int64) error {
	if _, err := global.Query.SysOperationLog.WithContext(ctx).Where(global.Query.SysOperationLog.UserID.Eq(uid)).Delete(); err != nil {
		global.Logger.Error(
			"清空操作日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

// ListOperationLog 获取操作日志列表
func (s *OperationService) ListOperationLog(ctx context.Context, req *basicDto.ListOperationLogReq) (*basicDto.ListOperationLogRes, error) {
	dao := global.Query.SysOperationLog
	do := dao.WithContext(ctx)

	// 用户ID条件
	if req.Uid != 0 {
		do = do.Where(dao.UserID.Eq(req.Uid))
	}

	// 操作模块条件
	if req.Module != "" {
		do = do.Where(dao.Module.Eq(req.Module))
	}

	// 操作类型条件
	if req.Operation != "" {
		do = do.Where(dao.Operation.Eq(req.Operation))
	}

	// HTTP方法条件
	if req.Method != "" {
		do = do.Where(dao.Method.Eq(req.Method))
	}

	// 操作状态条件
	if req.Status != 0 {
		do = do.Where(dao.Status.Eq(req.Status))
	}

	// 时间范围条件
	if req.StartTime != 0 {
		startTime := time.Unix(req.StartTime, 0)
		do = do.Where(dao.CreatedAt.Gte(startTime))
	}
	if req.EndTime != 0 {
		endTime := time.Unix(req.EndTime, 0)
		do = do.Where(dao.CreatedAt.Lte(endTime))
	}

	// 排序处理
	if req.Field != "" {
		switch req.Field {
		case "id":
			if req.Order == "asc" {
				do = do.Order(dao.ID)
			} else {
				do = do.Order(dao.ID.Desc())
			}
		case "user_id":
			if req.Order == "asc" {
				do = do.Order(dao.UserID)
			} else {
				do = do.Order(dao.UserID.Desc())
			}
		case "module":
			if req.Order == "asc" {
				do = do.Order(dao.Module)
			} else {
				do = do.Order(dao.Module.Desc())
			}
		case "operation":
			if req.Order == "asc" {
				do = do.Order(dao.Operation)
			} else {
				do = do.Order(dao.Operation.Desc())
			}
		case "method":
			if req.Order == "asc" {
				do = do.Order(dao.Method)
			} else {
				do = do.Order(dao.Method.Desc())
			}
		case "status":
			if req.Order == "asc" {
				do = do.Order(dao.Status)
			} else {
				do = do.Order(dao.Status.Desc())
			}
		case "created_at":
			if req.Order == "asc" {
				do = do.Order(dao.CreatedAt)
			} else {
				do = do.Order(dao.CreatedAt.Desc())
			}
		}
	} else {
		// 默认按创建时间倒序
		do = do.Order(dao.CreatedAt.Desc())
	}

	// 分页参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	// 使用FindByPage方法一次性完成分页查询和计数，提高性能
	offset := (req.Page - 1) * req.Size
	operationLogs, total, err := do.FindByPage(offset, req.Size)
	if err != nil {
		global.Logger.Error("查询操作日志列表失败", zap.Error(err))
		return nil, errs.NewError(1000, "查询操作日志列表失败")
	}

	// 转换为响应格式
	items := make([]*basicDto.ListOperationLogItem, 0, len(operationLogs))
	for _, log := range operationLogs {
		item := &basicDto.ListOperationLogItem{
			ID:        log.ID,
			UserID:    log.UserID,
			Username:  log.Username,
			Module:    log.Module,
			Operation: log.Operation,
			Method:    log.Method,
			URL:       log.URL,
			IP:        log.IP,
			Status:    log.Status,
			Duration:  log.CostTime,
			CreatedAt: log.CreatedAt,
		}
		items = append(items, item)
	}

	// 构建响应
	res := &basicDto.ListOperationLogRes{
		List:  items,
		Total: total,
	}

	return res, nil
}

// GetOperationLog 获取操作日志详情
func (s *OperationService) GetOperationLog(ctx context.Context, req *models.IDReq) (*basicDto.OperationLogDetailRes, error) {
	// 查询操作日志并预加载用户信息
	dao := global.Query.SysOperationLog
	do := dao.WithContext(ctx)

	// 预加载用户信息
	operationLog, err := do.Preload(dao.User).Where(dao.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("查询操作日志失败, id: %d", zap.Int64("id", req.ID))
			return nil, errs.ErrNotFound
		}
		global.Logger.Error("查询操作日志失败", zap.Error(err))
		return nil, errs.ErrServer
	}

	// 构建响应数据
	res := &basicDto.OperationLogDetailRes{
		ID:        operationLog.ID,
		UserID:    operationLog.UserID,
		Username:  "",
		Realname:  "",
		Nickname:  "",
		Avatar:    nil,
		Email:     nil,
		Phone:     nil,
		Module:    operationLog.Module,
		Operation: operationLog.Operation,
		Method:    operationLog.Method,
		URL:       operationLog.URL,
		IP:        operationLog.IP,
		Location:  operationLog.Location,
		UserAgent: operationLog.UserAgent,
		Params:    operationLog.RequestParams,
		Result:    operationLog.ResponseData,
		Status:    operationLog.Status,
		ErrorMsg:  operationLog.ErrorMsg,
		Duration:  operationLog.CostTime,
		CreatedAt: operationLog.CreatedAt,
	}

	// 如果有关联的用户信息，填充用户详细信息
	if operationLog.User != nil {
		res.Username = operationLog.User.Username
		res.Realname = operationLog.User.Realname
		res.Nickname = operationLog.User.Nickname
		res.Avatar = operationLog.User.Avatar
		res.Email = operationLog.User.Email
		res.Phone = operationLog.User.Phone
	} else if operationLog.Username != nil {
		// 如果没有关联用户但有用户名，使用存储的用户名
		res.Username = *operationLog.Username
	}

	return res, nil
}

// NewOperationLogService 创建操作日志服务实例
func NewOperationLogService() IOperationLogService {
	return &OperationService{}
}
