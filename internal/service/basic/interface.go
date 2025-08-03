package basic

import (
	"context"
	"sweet/internal/models"
	basicDto "sweet/internal/models/dto/basic"
)

type IBasicService interface {
	// LoginLog 获取登录日志服务
	LoginLog() ILoginLogService
	// OperationLog 获取操作日志服务
	OperationLog() IOperationLogService
}

// ILoginLogService 登录日志服务接口
type ILoginLogService interface {
	// CreateLoginLog 创建登录日志
	CreateLoginLog(ctx context.Context, req *basicDto.CreateLoginLogReq) error
	// DeleteLoginLog 删除登录日志
	DeleteLoginLog(ctx context.Context, req *basicDto.DeleteLoginLogReq) error
	// ClearAllLoginLog 清空所有登录日志
	ClearAllLoginLog(ctx context.Context) error
	// ClearLoginLog 清空自己的登录日志
	ClearLoginLog(ctx context.Context, uid int64) error
	// ListLoginLog 获取登录日志列表
	ListLoginLog(ctx context.Context, req *basicDto.ListLoginLogReq) (*basicDto.ListLoginLogRes, error)
	// GetLoginLog 获取登录日志详情
	GetLoginLog(ctx context.Context, req *models.IDReq) (*basicDto.LoginLogDetailRes, error)
}

// IOperationLogService 操作日志服务接口
type IOperationLogService interface {
	// CreateOperationLog 创建操作日志
	CreateOperationLog(ctx context.Context, req *basicDto.CreateOperationLogReq) error
	// DeleteOperationLog 删除操作日志
	DeleteOperationLog(ctx context.Context, req *basicDto.DeleteOperationLogReq) error
	// ClearAllOperationLog 清空所有操作日志
	ClearAllOperationLog(ctx context.Context) error
	// ClearOperationLog 清空自己的操作日志
	ClearOperationLog(ctx context.Context, uid int64) error
	// ListOperationLog 获取操作日志列表
	ListOperationLog(ctx context.Context, req *basicDto.ListOperationLogReq) (*basicDto.ListOperationLogRes, error)
	// GetOperationLog 获取操作日志详情
	GetOperationLog(ctx context.Context, req *models.IDReq) (*basicDto.OperationLogDetailRes, error)
}
