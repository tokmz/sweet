package vo

type LoginRes struct {
	Token  string `json:"token"`
	Expire int64  `json:"expire"`
}

type UserInfoRes struct {
	ID          int64    `json:"id"`          // 用户ID
	Username    string   `json:"username"`    // 用户名
	RealName    string   `json:"real_name"`   // 真实姓名
	NickName    string   `json:"nick_name"`   // 昵称
	Avatar      string   `json:"avatar"`      // 头像地址
	Email       string   `json:"email"`       // 邮箱
	Mobile      string   `json:"mobile"`      // 手机号
	Gender      int64    `json:"gender"`      // 性别 1-男 2-女 3-未知
	DeptID      int64    `json:"dept_id"`     // 部门ID
	PostID      int64    `json:"post_id"`     // 岗位ID
	RoleID      int64    `json:"role_id"`     // 角色ID
	RoleName    string   `json:"role_name"`   // 角色名称
	Permissions []string `json:"permissions"` // 权限列表
	Token       string   `json:"token"`       // 访问令牌
	ExpireAt    int64    `json:"expire_at"`   // 过期时间
	CreatedAt   int64    `json:"created_at"`  // 创建时间
}
