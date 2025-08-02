package basic

import (
	"context"
	"fmt"
	"sweet/internal/global"
	"sweet/internal/models"
	"sweet/internal/models/entity"
	"sweet/internal/models/query"
	basicDto "sweet/internal/models/dto/basic"
	"time"

	"gorm.io/gorm"
)

type LoginLogService struct{}

// CreateLoginLog 创建登录日志
func (s *LoginLogService) CreateLoginLog(ctx context.Context, req *basicDto.CreateLoginLogReq) error {
	if req == nil {
		return fmt.Errorf("请求参数不能为空")
	}

	// 参数验证
	if req.Username == "" {
		return fmt.Errorf("用户名不能为空")
	}
	if req.IP == "" {
		return fmt.Errorf("IP地址不能为空")
	}

	// 构建登录日志实体
	loginLog := &entity.SysLoginLog{
		UserID:     req.UserID,
		Username:   req.Username,
		LoginType:  req.LoginType,
		ClientType: req.ClientType,
		IP:         req.IP,
		Location:   req.Location,
		UserAgent:  req.UserAgent,
		DeviceInfo: req.DeviceInfo,
		Browser:    req.Browser,
		Os:         req.Os,
		Status:     req.Status,
		FailReason: req.FailReason,
	}

	// 使用事务创建登录日志
	q := query.Use(global.DBClient.DB())
	err := q.Transaction(func(tx *query.Query) error {
		return tx.SysLoginLog.WithContext(ctx).Create(loginLog)
	})

	if err != nil {
		return fmt.Errorf("创建登录日志失败: %w", err)
	}

	return nil
}

// DeleteLoginLog 删除登录日志
func (s *LoginLogService) DeleteLoginLog(ctx context.Context, req *basicDto.DeleteLoginLogReq) error {
	if req == nil {
		return fmt.Errorf("请求参数不能为空")
	}

	q := query.Use(global.DBClient.DB())

	// 使用事务删除登录日志
	err := q.Transaction(func(tx *query.Query) error {
		loginLogQuery := tx.SysLoginLog.WithContext(ctx)

		// 如果指定了用户ID，则删除该用户的所有登录日志
		if req.Uid > 0 {
			_, err := loginLogQuery.Where(tx.SysLoginLog.UserID.Eq(req.Uid)).Delete()
			return err
		}

		// 如果指定了日志ID列表，则删除指定的登录日志
		if len(req.Ids) > 0 {
			_, err := loginLogQuery.Where(tx.SysLoginLog.ID.In(req.Ids...)).Delete()
			return err
		}

		return fmt.Errorf("请指定要删除的用户ID或日志ID")
	})

	if err != nil {
		return fmt.Errorf("删除登录日志失败: %w", err)
	}

	return nil
}

// ClearLoginLog 清空指定用户的登录日志
func (s *LoginLogService) ClearLoginLog(ctx context.Context, uid int64) error {
	if uid <= 0 {
		return fmt.Errorf("用户ID无效")
	}

	q := query.Use(global.DBClient.DB())

	// 使用事务清空用户登录日志
	err := q.Transaction(func(tx *query.Query) error {
		_, err := tx.SysLoginLog.WithContext(ctx).Where(tx.SysLoginLog.UserID.Eq(uid)).Delete()
		return err
	})

	if err != nil {
		return fmt.Errorf("清空用户登录日志失败: %w", err)
	}

	return nil
}

// ListLoginLog 获取登录日志列表
func (s *LoginLogService) ListLoginLog(ctx context.Context, req *basicDto.ListLoginLogReq) (*basicDto.ListLoginLogRes, error) {
	if req == nil {
		return nil, fmt.Errorf("请求参数不能为空")
	}

	// 设置默认分页参数
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	q := query.Use(global.DBClient.DB())
	loginLogQuery := q.SysLoginLog.WithContext(ctx)

	// 构建查询条件
	if req.Uid > 0 {
		loginLogQuery = loginLogQuery.Where(q.SysLoginLog.UserID.Eq(req.Uid))
	}
	if req.LoginType > 0 {
		loginLogQuery = loginLogQuery.Where(q.SysLoginLog.LoginType.Eq(req.LoginType))
	}
	if req.ClientType > 0 {
		loginLogQuery = loginLogQuery.Where(q.SysLoginLog.ClientType.Eq(req.ClientType))
	}
	if req.Status > 0 {
		loginLogQuery = loginLogQuery.Where(q.SysLoginLog.Status.Eq(req.Status))
	}

	// 时间范围查询
	if req.StartTime > 0 {
		startTime := time.Unix(req.StartTime, 0)
		loginLogQuery = loginLogQuery.Where(q.SysLoginLog.CreatedAt.Gte(startTime))
	}
	if req.EndTime > 0 {
		endTime := time.Unix(req.EndTime, 0)
		loginLogQuery = loginLogQuery.Where(q.SysLoginLog.CreatedAt.Lte(endTime))
	}

	// 排序
	if req.Field != "" && req.Order != "" {
		switch req.Field {
		case "id":
			if req.Order == "asc" {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.ID)
			} else {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.ID.Desc())
			}
		case "username":
			if req.Order == "asc" {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.Username)
			} else {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.Username.Desc())
			}
		case "login_type":
			if req.Order == "asc" {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.LoginType)
			} else {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.LoginType.Desc())
			}
		case "status":
			if req.Order == "asc" {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.Status)
			} else {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.Status.Desc())
			}
		default:
			if req.Order == "asc" {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.CreatedAt)
			} else {
				loginLogQuery = loginLogQuery.Order(q.SysLoginLog.CreatedAt.Desc())
			}
		}
	} else {
		// 默认按创建时间倒序排列
		loginLogQuery = loginLogQuery.Order(q.SysLoginLog.CreatedAt.Desc())
	}

	// 分页查询
	offset := (req.Page - 1) * req.Size
	logs, total, err := loginLogQuery.FindByPage(offset, req.Size)
	if err != nil {
		return nil, fmt.Errorf("查询登录日志列表失败: %w", err)
	}

	// 转换为响应格式
	var items []*basicDto.ListLoginLogItem
	for _, log := range logs {
		item := &basicDto.ListLoginLogItem{
			ID:         log.ID,
			UserID:     log.UserID,
			Username:   log.Username,
			LoginType:  log.LoginType,
			ClientType: log.ClientType,
			Browser:    log.Browser,
			Os:         log.Os,
			Status:     log.Status,
			CreatedAt:  log.CreatedAt,
		}
		items = append(items, item)
	}

	return &basicDto.ListLoginLogRes{
		Total: total,
		List:  items,
	}, nil
}

// GetLoginLog 获取登录日志详情
func (s *LoginLogService) GetLoginLog(ctx context.Context, req *models.IDReq) (*basicDto.LoginLogDetailRes, error) {
	if req == nil || req.ID <= 0 {
		return nil, fmt.Errorf("日志ID无效")
	}

	q := query.Use(global.DBClient.DB())

	// 查询登录日志详情，包含用户信息
	log, err := q.SysLoginLog.WithContext(ctx).
		Preload(q.SysLoginLog.User).
		Where(q.SysLoginLog.ID.Eq(req.ID)).
		First()
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, fmt.Errorf("登录日志不存在")
		}
		return nil, fmt.Errorf("查询登录日志详情失败: %w", err)
	}

	// 构建响应数据
	detail := &basicDto.LoginLogDetailRes{
		ID:         log.ID,
		UserID:     log.UserID,
		Username:   log.Username,
		LoginType:  log.LoginType,
		ClientType: log.ClientType,
		IP:         log.IP,
		Location:   log.Location,
		UserAgent:  log.UserAgent,
		DeviceInfo: log.DeviceInfo,
		Browser:    log.Browser,
		Os:         log.Os,
		Status:     log.Status,
		FailReason: log.FailReason,
		CreatedAt:  log.CreatedAt,
	}

	// 如果有关联的用户信息，则填充用户详情
	if log.User != nil {
		detail.Realname = log.User.Realname
		detail.Nickname = log.User.Nickname
		detail.Avatar = log.User.Avatar
		detail.Email = log.User.Email
		detail.Phone = log.User.Phone
	}

	return detail, nil
}

// NewLoginLogService 创建登录日志服务实例
func NewLoginLogService() *LoginLogService {
	return &LoginLogService{}
}
