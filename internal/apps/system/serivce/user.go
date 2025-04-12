package serivce

import (
	"context"
	"sweet/internal/apps/system/types/dto"
	"sweet/internal/apps/system/types/entity"
	"sweet/pkg/resp"
)

var (
	IUser interface {
		// Create 创建用户
		Create(ctx context.Context, req *dto.CreateUserReq) error
		// Update 更新用户
		Update(ctx context.Context, req *dto.UpdateUserReq) error
		// Delete 删除用户
		Delete(ctx context.Context, req *dto.IdsReq) error
		// FindOne 查询用户
		FindOne(ctx context.Context, req *dto.IDReq) (*entity.User, error)
		// FindList 查询用户列表
		// List []*entity.User
		// Total int64
		FindList(ctx context.Context, req *dto.FindListUserReq) (*resp.Page, error)
		// UpdateStatus 批量更新用户状态
		UpdateStatus(ctx context.Context, req *dto.StatusReq) error
		// ResetPassword 重置密码
		ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error
	}

	// IUserAuth 用户认证接口
	IUserAuth interface {
		// AccountLogin 账号登录
		AccountLogin(ctx context.Context, req *dto.AccountLoginReq) (*dto.LoginRes, error)
		// MobileLogin 手机号登录
		MobileLogin(ctx context.Context, req *dto.MobileLoginReq) (*dto.LoginRes, error)
		// EmailLogin 邮箱登录
		EmailLogin(ctx context.Context, req *dto.EmailLoginReq) (*dto.LoginRes, error)
		// Logout 退出登录
		Logout(ctx context.Context, uid int64) error
		// RefreshToken 刷新Token
		RefreshToken(ctx context.Context, uid int64) (*dto.LoginRes, error)
		// Info 获取用户信息
		Info(ctx context.Context, uid int64) (*dto.UserInfoRes, error)
	}
)
