package system

import "sweet/internal/models"

// CreateUserReq 创建用户请求
type CreateUserReq struct {
}

// DeleteUserReq 删除用户请求
type DeleteUserReq models.IdsReq

// UpdateUserReq 更新用户请求
type UpdateUserReq struct {
	models.IDReq
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
	models.IDReq
	Username string `json:"username"`
	RoleID   int64  `json:"role_id"`
	DeptID   int64  `json:"dept_id"`
	PostID   int64  `json:"post_id"`
	Status   int64  `json:"status"`
}

type ListUserRes models.PageRes[ListUserItem]

type UserDetailRes struct{}
