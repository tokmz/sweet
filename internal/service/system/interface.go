package system

import (
	"context"
	"sweet/internal/models"
	systemDTO "sweet/internal/models/dto/system"
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
	CreateUser(ctx context.Context, req *systemDTO.CreateUserReq) error
	// DeleteUser 删除用户
	DeleteUser(ctx context.Context, req *systemDTO.DeleteUserReq) error
	// UpdateUser 更新用户
	UpdateUser(ctx context.Context, req *systemDTO.UpdateUserReq) error
	// ListUser 获取用户列表
	ListUser(ctx context.Context, req *systemDTO.ListUserReq) (*systemDTO.ListUserRes, error)
	// GetUserDetail 获取用户详情
	GetUserDetail(ctx context.Context, req *models.IDReq) (*systemDTO.UserDetailRes, error)
}

// IAuthService 认证服务接口
type IAuthService interface{}

// IRoleService 角色服务接口
type IRoleService interface {
	// 创建角色
	CreateRole(ctx context.Context, req *systemDTO.CreateRoleReq) error
	// 删除角色
	DeleteRole(ctx context.Context, req *systemDTO.DeleteRoleReq) error
	// 更新角色
	UpdateRole(ctx context.Context, req *systemDTO.UpdateRoleReq) error
	// 获取角色列表
	ListRole(ctx context.Context, req *systemDTO.RoleListReq) (*systemDTO.RoleListRes, error)
	// 获取角色详情
	GetRoleDetail(ctx context.Context, req *models.IDReq) (*systemDTO.RoleDetailRes, error)
	// 获取角色选项
	RoleOptions(ctx context.Context) (*systemDTO.RoleOptionRes, error)
	// 获取角色菜单Ids
	RoleMenuIds(ctx context.Context, req *models.IDReq) (*models.IdsReq, error)
	// 分配角色菜单
	AssignRoleMenuIds(ctx context.Context, req *systemDTO.AssignRoleMenuIdsReq) error
	// 获取角色ApiIds
	RoleApiIds(ctx context.Context, req *models.IDReq) (*models.IdsReq, error)
	// 分配角色ApiIds
	AssignRoleApiIds(ctx context.Context, req *systemDTO.AssignRoleApiIdsReq) error
}

// IMenuService 菜单服务接口
type IMenuService interface {
	// 创建菜单
	CreateMenu(ctx context.Context, req *systemDTO.CreateMenuReq) error
	// 删除菜单
	DeleteMenu(ctx context.Context, req *systemDTO.DeleteMenuReq) error
	// 更新菜单
	UpdateMenu(ctx context.Context, req *systemDTO.UpdateMenuReq) error
	// 菜单树
	MenuTree(ctx context.Context, req *models.IDReq) (*systemDTO.MenuTreeRes, error)
	// 菜单选项
	MenuOptions(ctx context.Context) (*systemDTO.MenuOptionsRes, error)
	// 菜单详情 - 包含按钮列表
	MenuDetail(ctx context.Context, req *models.IDReq) (*systemDTO.MenuDetailRes, error)
	// 创建按钮
	CreateButton(ctx context.Context, req *systemDTO.CreateButtonReq) error
	// 删除按钮
	DeleteButton(ctx context.Context, req *systemDTO.DeleteButtonReq) error
	// 更新按钮
	UpdateButton(ctx context.Context, req *systemDTO.UpdateButtonReq) error
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
