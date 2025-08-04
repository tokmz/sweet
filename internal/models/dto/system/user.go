package system

import (
	"time"
)

// CreateUserRequest 创建用户请求
type CreateUserRequest struct {
	Username string  `json:"username" binding:"required,min=3,max=32" comment:"登录用户名"`
	Password string  `json:"password" binding:"required,min=6,max=32" comment:"登录密码"`
	Realname string  `json:"realname" binding:"required,min=2,max=32" comment:"真实姓名"`
	Nickname string  `json:"nickname" binding:"required,min=2,max=32" comment:"昵称"`
	Avatar   *string `json:"avatar" comment:"头像"`
	Email    *string `json:"email" binding:"omitempty,email,max=64" comment:"邮箱"`
	Phone    *string `json:"phone" binding:"omitempty,max=20" comment:"手机号"`
	Status   *int64  `json:"status" binding:"omitempty,oneof=1 2" comment:"状态：1=正常，2=禁用"`
	RoleID   *int64  `json:"role_id" comment:"角色ID"`
	DeptID   *int64  `json:"dept_id" comment:"部门ID"`
	PostID   *int64  `json:"post_id" comment:"岗位ID"`
	Remark   *string `json:"remark" binding:"omitempty,max=255" comment:"备注"`
}

// UpdateUserRequest 更新用户请求
type UpdateUserRequest struct {
	ID       int64   `json:"id" binding:"required" comment:"用户ID"`
	Username string  `json:"username" binding:"required,min=3,max=32" comment:"登录用户名"`
	Realname string  `json:"realname" binding:"required,min=2,max=32" comment:"真实姓名"`
	Nickname string  `json:"nickname" binding:"required,min=2,max=32" comment:"昵称"`
	Avatar   *string `json:"avatar" comment:"头像"`
	Email    *string `json:"email" binding:"omitempty,email,max=64" comment:"邮箱"`
	Phone    *string `json:"phone" binding:"omitempty,max=20" comment:"手机号"`
	Status   *int64  `json:"status" binding:"omitempty,oneof=1 2" comment:"状态：1=正常，2=禁用"`
	RoleID   *int64  `json:"role_id" comment:"角色ID"`
	DeptID   *int64  `json:"dept_id" comment:"部门ID"`
	PostID   *int64  `json:"post_id" comment:"岗位ID"`
	Remark   *string `json:"remark" binding:"omitempty,max=255" comment:"备注"`
}

// ChangePasswordRequest 修改密码请求
type ChangePasswordRequest struct {
	ID          int64  `json:"id" binding:"required" comment:"用户ID"`
	OldPassword string `json:"old_password" binding:"required" comment:"原密码"`
	NewPassword string `json:"new_password" binding:"required,min=6,max=32" comment:"新密码"`
}

// UserQueryRequest 用户查询请求
type UserQueryRequest struct {
	Page     int     `json:"page" form:"page" binding:"omitempty,min=1" comment:"页码"`
	PageSize int     `json:"page_size" form:"page_size" binding:"omitempty,min=1,max=100" comment:"每页数量"`
	Keyword  *string `json:"keyword" form:"keyword" comment:"关键词搜索(用户名、真实姓名、昵称)"`
	Status   *int64  `json:"status" form:"status" binding:"omitempty,oneof=1 2" comment:"状态筛选"`
	RoleID   *int64  `json:"role_id" form:"role_id" comment:"角色筛选"`
	DeptID   *int64  `json:"dept_id" form:"dept_id" comment:"部门筛选"`
}

// UserResponse 用户响应
type UserResponse struct {
	ID        int64      `json:"id" comment:"用户ID"`
	Username  string     `json:"username" comment:"登录用户名"`
	Realname  string     `json:"realname" comment:"真实姓名"`
	Nickname  string     `json:"nickname" comment:"昵称"`
	Avatar    *string    `json:"avatar" comment:"头像"`
	Email     *string    `json:"email" comment:"邮箱"`
	Phone     *string    `json:"phone" comment:"手机号"`
	Status    *int64     `json:"status" comment:"状态：1=正常，2=禁用"`
	RoleID    *int64     `json:"role_id" comment:"角色ID"`
	RoleName  *string    `json:"role_name" comment:"角色名称"`
	DeptID    *int64     `json:"dept_id" comment:"部门ID"`
	DeptName  *string    `json:"dept_name" comment:"部门名称"`
	PostID    *int64     `json:"post_id" comment:"岗位ID"`
	PostName  *string    `json:"post_name" comment:"岗位名称"`
	Remark    *string    `json:"remark" comment:"备注"`
	CreatedAt *time.Time `json:"created_at" comment:"创建时间"`
	UpdatedAt *time.Time `json:"updated_at" comment:"更新时间"`
}

// UserListResponse 用户列表响应
type UserListResponse struct {
	Total int64          `json:"total" comment:"总数"`
	List  []UserResponse `json:"list" comment:"用户列表"`
}