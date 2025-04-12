package dto

import "sweet/internal/common"

type CreateUserReq struct {
	Username  string  `json:"username"`   // 用户名
	Password  string  `json:"password"`   // 密码
	RealName  *string `json:"real_name"`  // 真实姓名
	NickName  *string `json:"nick_name"`  // 昵称
	Avatar    *string `json:"avatar"`     // 头像地址
	Email     *string `json:"email"`      // 邮箱
	Mobile    *string `json:"mobile"`     // 手机号
	Gender    *int64  `json:"gender"`     // 性别 1-男 2-女 3-未知
	DeptID    *int64  `json:"dept_id"`    // 部门ID
	PostID    *int64  `json:"post_id"`    // 岗位ID
	RoleID    *int64  `json:"role_id"`    // 角色ID
	Remark    *string `json:"remark"`     // 备注
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	CreatedBy *int64  `json:"created_by"` // 创建人
}

type UpdateUserReq struct {
	common.IDReq
	Username  string  `json:"username"`   // 用户名
	RealName  *string `json:"real_name"`  // 真实姓名
	NickName  *string `json:"nick_name"`  // 昵称
	Avatar    *string `json:"avatar"`     // 头像地址
	Email     *string `json:"email"`      // 邮箱
	Mobile    *string `json:"mobile"`     // 手机号
	Gender    *int64  `json:"gender"`     // 性别 1-男 2-女 3-未知
	DeptID    *int64  `json:"dept_id"`    // 部门ID
	PostID    *int64  `json:"post_id"`    // 岗位ID
	RoleID    *int64  `json:"role_id"`    // 角色ID
	Remark    *string `json:"remark"`     // 备注
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	UpdatedBy *int64  `json:"updated_by"` // 更新人
}

type FindListUserReq struct {
	Username string `json:"username"`  // 用户名
	RealName string `json:"real_name"` // 真实姓名
	Status   int64  `json:"status"`    // 状态
	Gender   int64  `json:"gender"`    // 性别 1-男 2-女 3-未知
	DeptID   int64  `json:"dept_id"`   // 部门ID
	PostID   int64  `json:"post_id"`   // 岗位ID
	RoleID   int64  `json:"role_id"`   // 角色ID
	common.PageReq
}

type ResetPasswordReq struct {
	common.IDReq
	Password string `json:"password"` // 密码
}

type AccountLoginReq struct {
	Username string `json:"username"` // 用户名
	Password string `json:"password"` // 密码
}

type EmailLoginReq struct {
	Email    string `json:"email"`    // 邮箱
	Code     string `json:"code"`     // 验证码
	Password string `json:"password"` // 密码
}

type MobileLoginReq struct {
	Mobile   string `json:"mobile"`   // 手机号
	Code     string `json:"code"`     // 验证码
	Password string `json:"password"` // 密码
}
