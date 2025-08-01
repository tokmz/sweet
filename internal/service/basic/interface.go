package basic

import (
	"context"
	"sweet/internal/models/dto/basic"
	"sweet/internal/models/entity"
)

// LoginLogService 登录日志服务接口
type LoginLogService interface {
	// CreateLoginLog 创建登录日志
	CreateLoginLog(ctx context.Context, req *basic.LoginLogCreateReq) (*basic.LoginLogResp, error)

	// UpdateLoginLog 更新登录日志（主要用于更新登录持续时间、退出类型等）
	UpdateLoginLog(ctx context.Context, req *basic.LoginLogUpdateReq) error

	// GetLoginLogByID 根据ID获取登录日志
	GetLoginLogByID(ctx context.Context, id int64) (*basic.LoginLogResp, error)

	// GetLoginLogList 获取登录日志列表
	GetLoginLogList(ctx context.Context, req *basic.LoginLogQueryReq) (*basic.LoginLogListResp, error)

	// DeleteLoginLog 删除登录日志（软删除）
	DeleteLoginLog(ctx context.Context, id int64) error

	// BatchDeleteLoginLog 批量删除登录日志
	BatchDeleteLoginLog(ctx context.Context, ids []int64) error

	// GetLoginLogStats 获取登录日志统计信息
	GetLoginLogStats(ctx context.Context, userID *int64) (*basic.LoginLogStatsResp, error)

	// GetUserLoginHistory 获取用户登录历史
	GetUserLoginHistory(ctx context.Context, userID int64, limit int) ([]*basic.LoginLogResp, error)

	// GetOnlineUsers 获取在线用户列表
	GetOnlineUsers(ctx context.Context) ([]*basic.LoginLogResp, error)

	// ForceLogoutUser 强制用户下线
	ForceLogoutUser(ctx context.Context, sessionID string) error

	// CleanExpiredLogs 清理过期日志
	CleanExpiredLogs(ctx context.Context, days int) error

	// GetRiskLoginLogs 获取高风险登录日志
	GetRiskLoginLogs(ctx context.Context, riskLevel int64, limit int) ([]*basic.LoginLogResp, error)

	// RecordLoginSuccess 记录登录成功
	RecordLoginSuccess(ctx context.Context, userID int64, username, ip, userAgent string, loginType, clientType int64) (*entity.SysLoginLog, error)

	// RecordLoginFail 记录登录失败
	RecordLoginFail(ctx context.Context, username, ip, userAgent, failReason string, loginType, clientType int64) (*entity.SysLoginLog, error)

	// RecordLogout 记录用户退出
	RecordLogout(ctx context.Context, sessionID string, logoutType int64) error

	// UpdateLoginDuration 更新登录持续时间
	UpdateLoginDuration(ctx context.Context, sessionID string, duration int64) error

	// CheckSuspiciousLogin 检查可疑登录
	CheckSuspiciousLogin(ctx context.Context, userID int64, ip string) (bool, error)

	// GetLoginTrend 获取登录趋势数据
	GetLoginTrend(ctx context.Context, days int) (map[string]int64, error)

	// ExportLoginLogs 导出登录日志
	ExportLoginLogs(ctx context.Context, req *basic.LoginLogQueryReq) ([]byte, error)
}