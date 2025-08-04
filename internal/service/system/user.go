package system

import (
	"context"
	"fmt"
	"strings"
	"time"

	"gorm.io/gorm"

	"sweet/internal/models/entity"
	dto "sweet/internal/models/dto/system"
	"sweet/pkg/crypto"
	"sweet/pkg/errs"
)

// UserService 用户服务实现
type UserService struct {
	db *gorm.DB
}

// NewUserService 创建用户服务实例
func NewUserService(db *gorm.DB) IUserService {
	return &UserService{
		db: db,
	}
}

// CreateUser 创建用户
func (s *UserService) CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error) {
	// 检查用户名是否已存在
	exists, err := s.CheckUserExists(ctx, req.Username)
	if err != nil {
		return nil, err
	}
	if exists {
		return nil, errs.NewBusinessError("用户名已存在")
	}

	// 生成密码盐和哈希
	salt := crypto.GenerateSalt()
	hashedPassword := crypto.HashPassword(req.Password, salt)

	// 设置默认状态
	status := int64(1)
	if req.Status != nil {
		status = *req.Status
	}

	// 创建用户实体
	user := &entity.SysUser{
		Username: req.Username,
		Password: hashedPassword,
		Salt:     salt,
		Realname: req.Realname,
		Nickname: req.Nickname,
		Avatar:   req.Avatar,
		Email:    req.Email,
		Phone:    req.Phone,
		Status:   &status,
		RoleID:   req.RoleID,
		DeptID:   req.DeptID,
		PostID:   req.PostID,
		Remark:   req.Remark,
	}

	// 保存到数据库
	if err := s.db.WithContext(ctx).Create(user).Error; err != nil {
		return nil, fmt.Errorf("创建用户失败: %w", err)
	}

	// 返回用户信息
	return s.GetUserByID(ctx, user.ID)
}

// UpdateUser 更新用户
func (s *UserService) UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error) {
	// 检查用户是否存在
	var user entity.SysUser
	if err := s.db.WithContext(ctx).First(&user, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewBusinessError("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	// 检查用户名是否被其他用户使用
	if user.Username != req.Username {
		exists, err := s.CheckUserExists(ctx, req.Username, req.ID)
		if err != nil {
			return nil, err
		}
		if exists {
			return nil, errs.NewBusinessError("用户名已被其他用户使用")
		}
	}

	// 更新用户信息
	updates := map[string]interface{}{
		"username": req.Username,
		"realname": req.Realname,
		"nickname": req.Nickname,
		"avatar":   req.Avatar,
		"email":    req.Email,
		"phone":    req.Phone,
		"status":   req.Status,
		"role_id":  req.RoleID,
		"dept_id":  req.DeptID,
		"post_id":  req.PostID,
		"remark":   req.Remark,
	}

	if err := s.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		return nil, fmt.Errorf("更新用户失败: %w", err)
	}

	// 返回更新后的用户信息
	return s.GetUserByID(ctx, req.ID)
}

// DeleteUser 删除用户
func (s *UserService) DeleteUser(ctx context.Context, id int64) error {
	// 检查用户是否存在
	var user entity.SysUser
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errs.NewBusinessError("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 软删除用户
	if err := s.db.WithContext(ctx).Delete(&user).Error; err != nil {
		return fmt.Errorf("删除用户失败: %w", err)
	}

	return nil
}

// GetUserByID 根据ID获取用户
func (s *UserService) GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error) {
	var user entity.SysUser
	query := s.db.WithContext(ctx).
		Preload("Role").
		Preload("Dept").
		Preload("Post")

	if err := query.First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewBusinessError("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return s.convertToUserResponse(&user), nil
}

// GetUserByUsername 根据用户名获取用户
func (s *UserService) GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error) {
	var user entity.SysUser
	query := s.db.WithContext(ctx).
		Preload("Role").
		Preload("Dept").
		Preload("Post")

	if err := query.Where("username = ?", username).First(&user).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errs.NewBusinessError("用户不存在")
		}
		return nil, fmt.Errorf("查询用户失败: %w", err)
	}

	return s.convertToUserResponse(&user), nil
}

// GetUserList 获取用户列表
func (s *UserService) GetUserList(ctx context.Context, req *dto.UserQueryRequest) (*dto.UserListResponse, error) {
	// 设置默认分页参数
	page := 1
	pageSize := 10
	if req.Page > 0 {
		page = req.Page
	}
	if req.PageSize > 0 {
		pageSize = req.PageSize
	}

	// 构建查询条件
	query := s.db.WithContext(ctx).Model(&entity.SysUser{}).
		Preload("Role").
		Preload("Dept").
		Preload("Post")

	// 关键词搜索
	if req.Keyword != nil && *req.Keyword != "" {
		keyword := "%" + *req.Keyword + "%"
		query = query.Where("username LIKE ? OR realname LIKE ? OR nickname LIKE ?", keyword, keyword, keyword)
	}

	// 状态筛选
	if req.Status != nil {
		query = query.Where("status = ?", *req.Status)
	}

	// 角色筛选
	if req.RoleID != nil {
		query = query.Where("role_id = ?", *req.RoleID)
	}

	// 部门筛选
	if req.DeptID != nil {
		query = query.Where("dept_id = ?", *req.DeptID)
	}

	// 查询总数
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, fmt.Errorf("查询用户总数失败: %w", err)
	}

	// 分页查询
	var users []entity.SysUser
	offset := (page - 1) * pageSize
	if err := query.Offset(offset).Limit(pageSize).Find(&users).Error; err != nil {
		return nil, fmt.Errorf("查询用户列表失败: %w", err)
	}

	// 转换为响应格式
	list := make([]dto.UserResponse, len(users))
	for i, user := range users {
		list[i] = *s.convertToUserResponse(&user)
	}

	return &dto.UserListResponse{
		Total: total,
		List:  list,
	}, nil
}

// ChangePassword 修改密码
func (s *UserService) ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error {
	// 查询用户
	var user entity.SysUser
	if err := s.db.WithContext(ctx).First(&user, req.ID).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errs.NewBusinessError("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 验证原密码
	if !crypto.VerifyPassword(req.OldPassword, user.Salt, user.Password) {
		return errs.NewBusinessError("原密码错误")
	}

	// 生成新密码哈希
	newSalt := crypto.GenerateSalt()
	newHashedPassword := crypto.HashPassword(req.NewPassword, newSalt)

	// 更新密码
	updates := map[string]interface{}{
		"password": newHashedPassword,
		"salt":     newSalt,
	}

	if err := s.db.WithContext(ctx).Model(&user).Updates(updates).Error; err != nil {
		return fmt.Errorf("更新密码失败: %w", err)
	}

	return nil
}

// UpdateUserStatus 更新用户状态
func (s *UserService) UpdateUserStatus(ctx context.Context, id int64, status int64) error {
	// 检查用户是否存在
	var user entity.SysUser
	if err := s.db.WithContext(ctx).First(&user, id).Error; err != nil {
		if err == gorm.ErrRecordNotFound {
			return errs.NewBusinessError("用户不存在")
		}
		return fmt.Errorf("查询用户失败: %w", err)
	}

	// 更新状态
	if err := s.db.WithContext(ctx).Model(&user).Update("status", status).Error; err != nil {
		return fmt.Errorf("更新用户状态失败: %w", err)
	}

	return nil
}

// CheckUserExists 检查用户是否存在
func (s *UserService) CheckUserExists(ctx context.Context, username string, excludeID ...int64) (bool, error) {
	query := s.db.WithContext(ctx).Model(&entity.SysUser{}).Where("username = ?", username)

	// 排除指定ID
	if len(excludeID) > 0 && excludeID[0] > 0 {
		query = query.Where("id != ?", excludeID[0])
	}

	var count int64
	if err := query.Count(&count).Error; err != nil {
		return false, fmt.Errorf("检查用户是否存在失败: %w", err)
	}

	return count > 0, nil
}

// convertToUserResponse 转换为用户响应格式
func (s *UserService) convertToUserResponse(user *entity.SysUser) *dto.UserResponse {
	resp := &dto.UserResponse{
		ID:        user.ID,
		Username:  user.Username,
		Realname:  user.Realname,
		Nickname:  user.Nickname,
		Avatar:    user.Avatar,
		Email:     user.Email,
		Phone:     user.Phone,
		Status:    user.Status,
		RoleID:    user.RoleID,
		DeptID:    user.DeptID,
		PostID:    user.PostID,
		Remark:    user.Remark,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}

	// 设置关联信息
	if user.Role != nil {
		roleName := user.Role.Name
		resp.RoleName = &roleName
	}

	if user.Dept != nil {
		deptName := user.Dept.Name
		resp.DeptName = &deptName
	}

	if user.Post != nil {
		postName := user.Post.Name
		resp.PostName = &postName
	}

	return resp
}