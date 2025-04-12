package service

import (
	"context"
	"sweet/internal/apps/system/repo"
	"sweet/internal/apps/system/types/dto"
	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/vo"
	"sweet/internal/common"
	"sweet/pkg/crypto"
	"sweet/pkg/logger"
	"sweet/pkg/utils"
)

// 将接口变量定义改为类型定义
// UserService 用户服务接口
type UserService interface {
	// Create 创建用户
	Create(ctx context.Context, req *dto.CreateUserReq) error
	// Update 更新用户
	Update(ctx context.Context, req *dto.UpdateUserReq) error
	// Delete 删除用户
	Delete(ctx context.Context, req *common.IdsReq) error
	// FindOne 查询用户
	FindOne(ctx context.Context, req *common.IDReq) (*entity.User, error)
	// FindList 查询用户列表
	FindList(ctx context.Context, req *dto.FindListUserReq) (*common.Page, error)
	// UpdateStatus 批量更新用户状态
	UpdateStatus(ctx context.Context, req *common.StatusReq) error
	// ResetPassword 重置密码
	ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error
}

// IUserAuthService 用户认证接口
type UserAuthService interface {
	// AccountLogin 账号登录
	AccountLogin(ctx context.Context, req *dto.AccountLoginReq) (*vo.UserInfoRes, error)
	// MobileLogin 手机号登录
	MobileLogin(ctx context.Context, req *dto.MobileLoginReq) (*vo.UserInfoRes, error)
	// EmailLogin 邮箱登录
	EmailLogin(ctx context.Context, req *dto.EmailLoginReq) (*vo.UserInfoRes, error)
	// Logout 退出登录
	Logout(ctx context.Context, uid int64) error
	// RefreshToken 刷新Token
	RefreshToken(ctx context.Context, uid int64) (*vo.UserInfoRes, error)
	// Info 获取用户信息
	Info(ctx context.Context, uid int64) (*vo.UserInfoRes, error)
}

// userService 用户服务实现
type userService struct {
	userRepo repo.UserRepository
}

// Create 创建用户
func (s *userService) Create(ctx context.Context, req *dto.CreateUserReq) error {
	logger.Info("创建用户",
		logger.String("用户名", req.Username),
		logger.String("真实姓名", utils.SafeString(req.RealName)))

	// 生成密码盐值
	salt := crypto.Hash.Salt()

	// 构建用户实体
	user := &entity.User{
		Username:  req.Username,
		Password:  crypto.Hash.Md5Salt(req.Password, salt), // 加盐哈希后的密码
		Salt:      &salt,                                   // 保存盐值
		RealName:  req.RealName,
		NickName:  req.NickName,
		Avatar:    req.Avatar,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Gender:    req.Gender,
		DeptID:    req.DeptID,
		PostID:    req.PostID,
		RoleID:    req.RoleID,
		Remark:    req.Remark,
		Status:    req.Status,
		CreatedBy: req.CreatedBy,
	}

	// 调用仓储层创建用户
	return s.userRepo.Create(ctx, user)
}

// Update 更新用户
func (s *userService) Update(ctx context.Context, req *dto.UpdateUserReq) error {
	logger.Info("更新用户",
		logger.Int64("用户ID", req.ID),
		logger.String("用户名", req.Username))

	// 构建用户实体
	user := &entity.User{
		ID:        req.ID,
		Username:  req.Username,
		RealName:  req.RealName,
		NickName:  req.NickName,
		Avatar:    req.Avatar,
		Email:     req.Email,
		Mobile:    req.Mobile,
		Gender:    req.Gender,
		DeptID:    req.DeptID,
		PostID:    req.PostID,
		RoleID:    req.RoleID,
		Remark:    req.Remark,
		Status:    req.Status,
		UpdatedBy: req.UpdatedBy,
	}

	// 调用仓储层更新用户
	return s.userRepo.Update(ctx, user)
}

// Delete 删除用户
func (s *userService) Delete(ctx context.Context, req *common.IdsReq) error {
	logger.Info("删除用户",
		logger.Any("用户ID列表", req.Ids),
		logger.Int("删除数量", len(req.Ids)))

	// 调用仓储层删除用户
	return s.userRepo.Delete(ctx, req.Ids)
}

// FindOne 查询用户
func (s *userService) FindOne(ctx context.Context, req *common.IDReq) (*entity.User, error) {
	logger.Info("查询用户", logger.Int64("用户ID", req.ID))

	// 调用仓储层查询用户
	return s.userRepo.FindOne(ctx, req.ID)
}

// FindList 查询用户列表
func (s *userService) FindList(ctx context.Context, req *dto.FindListUserReq) (*common.Page, error) {
	logger.Info("查询用户列表",
		logger.String("用户名", req.Username),
		logger.String("真实姓名", req.RealName),
		logger.Int("页码", req.Page),
		logger.Int("每页数量", req.Size))

	// 构建查询参数
	params := &repo.UserListParams{
		Username: req.Username,
		RealName: req.RealName,
		Status:   req.Status,
		DeptID:   req.DeptID,
		RoleID:   req.RoleID,
		PostID:   req.PostID,
		Page:     req.Page,
		Size:     req.Size,
	}

	// 调用仓储层查询用户列表
	list, total, err := s.userRepo.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建分页响应
	return common.NewPage(list, total), nil
}

// UpdateStatus 批量更新用户状态
func (s *userService) UpdateStatus(ctx context.Context, req *common.StatusReq) error {
	logger.Info("批量更新用户状态",
		logger.Any("用户ID列表", req.Ids),
		logger.Int64("状态", req.Status),
		logger.Int("数量", len(req.Ids)))

	// 直接调用仓储层的批量更新方法
	return s.userRepo.UpdateStatus(ctx, req.Ids, req.Status)
}

// ResetPassword 重置密码
func (s *userService) ResetPassword(ctx context.Context, req *dto.ResetPasswordReq) error {
	logger.Info("重置用户密码", logger.Int64("用户ID", req.ID))

	// 先查询用户
	user, err := s.userRepo.FindOne(ctx, req.ID)
	if err != nil {
		return err
	}

	// 设置新密码，注意：实际应用中应该对密码进行加密处理
	user.Password = crypto.Hash.Md5Salt(req.Password, utils.SafeString(user.Salt))

	// 调用仓储层更新用户
	return s.userRepo.Update(ctx, user)
}

// NewUserService 创建用户服务
func NewUserService(userRepo repo.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}
