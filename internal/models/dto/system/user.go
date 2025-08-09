package system

import (
	"sweet/internal/models"
	"time"
)

// CreateUserReq 创建用户请求
type CreateUserReq struct {
	Username string  `json:"username" binding:"required,min=5,max=15"`      // 登录用户名 5-15位
	Password string  `json:"password" binding:"required,min=6,max=15"`      // 登录密码 6-15位
	Realname string  `json:"realname" binding:"required"`                   // 真实姓名
	Nickname string  `json:"nickname" binding:"required"`                   // 昵称
	Avatar   *string `json:"avatar"`                                        // 头像
	Email    *string `json:"email"`                                         // 邮箱
	Phone    *string `json:"phone"`                                         // 手机号
	Status   *int64  `json:"status" binding:"required;oneof=1 2;default=1"` // 状态：1=正常，2=禁用
	RoleID   *int64  `json:"role_id" binding:"required"`                    // 角色ID
	DeptID   *int64  `json:"dept_id" binding:"required"`                    // 部门ID
	PostID   *int64  `json:"post_id" binding:"required"`                    // 岗位ID
	Remark   *string `json:"remark"`                                        // 备注
}

// DeleteUserReq 删除用户请求
type DeleteUserReq models.IdsReq

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	models.IDReq
	Password string  `json:"password"` // 登录密码
	Realname string  `json:"realname"` // 真实姓名
	Nickname string  `json:"nickname"` // 昵称
	Avatar   *string `json:"avatar"`   // 头像
	Email    *string `json:"email"`    // 邮箱
	Phone    *string `json:"phone"`    // 手机号
	Status   *int64  `json:"status"`   // 状态：1=正常，2=禁用
	RoleID   *int64  `json:"role_id"`  // 角色ID
	DeptID   *int64  `json:"dept_id"`  // 部门ID
	PostID   *int64  `json:"post_id"`  // 岗位ID
	Remark   *string `json:"remark"`   // 备注
}

// ListUserReq 获取用户列表请求
type ListUserReq struct {
	Username string `json:"username" form:"username"`
	RoleID   int64  `json:"role_id" form:"role_id"`
	DeptID   int64  `json:"dept_id" form:"dept_id"`
	PostID   int64  `json:"post_id" form:"post_id"`
	Status   int64  `json:"status" form:"status"`
	models.SortReq
	models.PageReq
	models.TimeRangeReq
}

// ListUserItem 用户列表项
type ListUserItem struct {
	ID        int64      `json:"id"`         // 管理员ID
	Username  string     `json:"username"`   // 登录用户名
	Realname  string     `json:"realname"`   // 真实姓名
	Nickname  string     `json:"nickname"`   // 昵称
	Avatar    *string    `json:"avatar"`     // 头像
	Status    *int64     `json:"status"`     // 状态：1=正常，2=禁用
	RoleID    int64      `json:"role_id"`    // 角色ID
	RoleName  string     `json:"role_name"`  // 角色名称
	CreatedAt *time.Time `json:"created_at"` // 创建时间
}

type ListUserRes models.PageRes[ListUserItem]

type UserDetailRes struct {
	ID        int64      `json:"id"`         // 管理员ID
	Username  string     `json:"username"`   // 登录用户名
	Password  string     `json:"password"`   // 登录密码
	Salt      string     `json:"salt"`       // 密码盐
	Realname  string     `json:"realname"`   // 真实姓名
	Nickname  string     `json:"nickname"`   // 昵称
	Avatar    *string    `json:"avatar"`     // 头像
	Email     *string    `json:"email"`      // 邮箱
	Phone     *string    `json:"phone"`      // 手机号
	Status    *int64     `json:"status"`     // 状态：1=正常，2=禁用
	RoleID    int64      `json:"role_id"`    // 角色ID
	RoleName  string     `json:"role_name"`  // 角色名称
	DeptID    int64      `json:"dept_id"`    // 部门ID
	DeptName  string     `json:"dept_name"`  // 部门名称
	PostID    int64      `json:"post_id"`    // 岗位ID
	PostName  string     `json:"post_name"`  // 岗位名称
	Remark    *string    `json:"remark"`     // 备注
	CreatedAt *time.Time `json:"created_at"` // 创建时间
	UpdatedAt *time.Time `json:"updated_at"` // 更新时间
}
