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

type LoginLogService struct{}

func (s *LoginLogService) CreateLoginLog(ctx context.Context, req *basicDto.CreateLoginLogReq) error {
	if err := global.Query.SysLoginLog.WithContext(ctx).Create(&entity.SysLoginLog{
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
	}); err != nil {
		global.Logger.Error(
			"创建登录日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

func (s *LoginLogService) DeleteLoginLog(ctx context.Context, req *basicDto.DeleteLoginLogReq) error {
	dao := global.Query.SysLoginLog
	do := dao.WithContext(ctx)
	if req.Uid != 0 {
		do = do.Where(dao.UserID.Eq(req.Uid))
	}

	if _, err := do.Where(dao.ID.In(req.Ids...)).Delete(); err != nil {
		global.Logger.Error(
			"删除登录日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

func (s *LoginLogService) ClearLoginLog(ctx context.Context, uid int64) error {
	if _, err := global.Query.SysLoginLog.WithContext(ctx).Where(global.Query.SysLoginLog.UserID.Eq(uid)).Delete(); err != nil {
		global.Logger.Error(
			"清空登录日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

func (s *LoginLogService) ClearAllLoginLog(ctx context.Context) error {
	if _, err := global.Query.SysLoginLog.WithContext(ctx).Delete(); err != nil {
		global.Logger.Error(
			"清空所有登录日志失败",
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

func (s *LoginLogService) ListLoginLog(ctx context.Context, req *basicDto.ListLoginLogReq) (*basicDto.ListLoginLogRes, error) {
	dao := global.Query.SysLoginLog
	do := dao.WithContext(ctx)

	// 用户ID条件
	if req.Uid != 0 {
		do = do.Where(dao.UserID.Eq(req.Uid))
	}

	// 登录类型条件
	if req.LoginType != 0 {
		do = do.Where(dao.LoginType.Eq(req.LoginType))
	}

	// 客户端类型条件
	if req.ClientType != 0 {
		do = do.Where(dao.ClientType.Eq(req.ClientType))
	}

	// 登录状态条件
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
		case "login_type":
			if req.Order == "asc" {
				do = do.Order(dao.LoginType)
			} else {
				do = do.Order(dao.LoginType.Desc())
			}
		case "client_type":
			if req.Order == "asc" {
				do = do.Order(dao.ClientType)
			} else {
				do = do.Order(dao.ClientType.Desc())
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
	loginLogs, total, err := do.FindByPage(offset, req.Size)
	if err != nil {
		global.Logger.Error("查询登录日志列表失败", zap.Error(err))
		return nil, errs.NewError(1000, "查询登录日志列表失败")
	}

	// 转换为响应格式
	items := make([]*basicDto.ListLoginLogItem, 0, len(loginLogs))
	for _, log := range loginLogs {
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

	// 构建响应
	res := &basicDto.ListLoginLogRes{
		List:  items,
		Total: total,
	}

	return res, nil
}

func (s *LoginLogService) GetLoginLog(ctx context.Context, req *models.IDReq) (*basicDto.LoginLogDetailRes, error) {
	// 查询登录日志并预加载用户信息
	dao := global.Query.SysLoginLog
	do := dao.WithContext(ctx)

	// 预加载用户信息
	loginLog, err := do.Preload(dao.User).Where(dao.ID.Eq(req.ID)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error("查询登录日志失败, id: %d", zap.Int64("id", req.ID))
			return nil, errs.ErrNotFound
		}
		global.Logger.Error("查询登录日志失败", zap.Error(err))
		return nil, errs.ErrServer
	}

	// 构建响应数据
	res := &basicDto.LoginLogDetailRes{
		ID:         loginLog.ID,
		UserID:     loginLog.UserID,
		Username:   loginLog.Username,
		LoginType:  loginLog.LoginType,
		ClientType: loginLog.ClientType,
		IP:         loginLog.IP,
		Location:   loginLog.Location,
		UserAgent:  loginLog.UserAgent,
		DeviceInfo: loginLog.DeviceInfo,
		Browser:    loginLog.Browser,
		Os:         loginLog.Os,
		Status:     loginLog.Status,
		FailReason: loginLog.FailReason,
		CreatedAt:  loginLog.CreatedAt,
	}

	// 如果有关联的用户信息，填充用户详细信息
	if loginLog.User != nil {
		res.Realname = loginLog.User.Realname
		res.Nickname = loginLog.User.Nickname
		res.Avatar = loginLog.User.Avatar
		res.Email = loginLog.User.Email
		res.Phone = loginLog.User.Phone
	}

	return res, nil
}

func NewLoginLogService() ILoginLogService {
	return &LoginLogService{}
}
