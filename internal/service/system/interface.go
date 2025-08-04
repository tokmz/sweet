package system

import (
	"context"
	
	dto "sweet/internal/models/dto/system"
)

// ISystemService 系统服务接口
type ISystemService interface {
	// User 用户服务接口
	User() IUserService
	// Auth 认证服务接口
	Auth() IAuthService
	// Role 角色服务接口
	Role() IRoleService
	// Menu 菜单服务接口
	Menu() IMenuService
	// Api API服务接口
	Api() IApiService
	// Dept 部门服务接口
	Dept() IDeptService
	// Post 岗位服务接口
	Post() IPostService
}

// IUserService 用户服务接口
type IUserService interface {
	// CreateUser 创建用户
	CreateUser(ctx context.Context, req *dto.CreateUserRequest) (*dto.UserResponse, error)
	
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, req *dto.UpdateUserRequest) (*dto.UserResponse, error)
	
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, id int64) error
	
	// GetUserByID 根据ID获取用户
	GetUserByID(ctx context.Context, id int64) (*dto.UserResponse, error)
	
	// GetUserByUsername 根据用户名获取用户
	GetUserByUsername(ctx context.Context, username string) (*dto.UserResponse, error)
	
	// GetUserList 获取用户列表
	GetUserList(ctx context.Context, req *dto.UserQueryRequest) (*dto.UserListResponse, error)
	
	// ChangePassword 修改密码
	ChangePassword(ctx context.Context, req *dto.ChangePasswordRequest) error
	
	// UpdateUserStatus 更新用户状态
	UpdateUserStatus(ctx context.Context, id int64, status int64) error
	
	// CheckUserExists 检查用户是否存在
	CheckUserExists(ctx context.Context, username string, excludeID ...int64) (bool, error)
}

// IAuthService 认证服务接口
type IAuthService interface {
}

// IRoleService 角色服务接口
type IRoleService interface {
}

// IMenuService 菜单服务接口
type IMenuService interface {
}

// IApiService API服务接口
type IApiService interface {
}

// IDeptService 部门服务接口
type IDeptService interface {
}

// IPostService 岗位服务接口
type IPostService interface {
}
