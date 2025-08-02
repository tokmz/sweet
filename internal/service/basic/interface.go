package basic

import (
	"context"
	"sweet/internal/models"
	basicDto "sweet/internal/models/dto/basic"
)

// ILoginLogService 登录日志服务接口
type ILoginLogService interface {
	// CreateLoginLog 创建登录日志
	CreateLoginLog(ctx context.Context, req *basicDto.CreateLoginLogReq) error
	// DeleteLoginLog 删除登录日志
	DeleteLoginLog(ctx context.Context, req *basicDto.DeleteLoginLogReq) error
	// ClearLoginLog 清空自己的登录日志
	ClearLoginLog(ctx context.Context, uid int64) error
	// ListLoginLog 获取登录日志列表
	ListLoginLog(ctx context.Context, req *basicDto.ListLoginLogReq) (*basicDto.ListLoginLogRes, error)
	// GetLoginLog 获取登录日志详情
	GetLoginLog(ctx context.Context, req *models.IDReq) (*basicDto.LoginLogDetailRes, error)
}
