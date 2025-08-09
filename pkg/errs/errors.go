package errs

// common error
var (
	ErrServer        = NewError(1000, "服务异常")
	ErrParams        = NewError(1001, "参数异常")
	ErrAuthorization = NewError(1002, "认证异常")
	ErrNotFound      = NewError(1003, "资源不存在")
	ErrConflict      = NewError(1004, "资源冲突")
)

// system user error
var (
	ErrLoginFromOther = NewError(1005, "其他设备登录")
	ErrUserNotFound   = NewError(1006, "账号不存在")
	ErrUserExists     = NewError(1007, "该用户名已被注册")
	ErrPhoneNotValid  = NewError(1008, "手机号格式错误")
	ErrPhoneNotFound  = NewError(1009, "该手机号不存在")
	ErrPhoneExists    = NewError(1010, "该手机号已被绑定")
	ErrEmailNotValid  = NewError(1011, "邮箱格式错误")
	ErrEmailNotFound  = NewError(1012, "邮箱不存在")
	ErrEmailExists    = NewError(1013, "该邮箱已被绑定")
	ErrPassword       = NewError(1014, "密码错误")
)

// system role error
var (
	ErrRoleNotFound           = NewError(1020, "角色不存在")
	ErrRoleCodeExists         = NewError(1021, "角色编码已存在")
	ErrRoleNameExists         = NewError(1022, "角色名称已存在")
	ErrRoleInUse              = NewError(1023, "角色正在使用中，无法删除")
	ErrSystemRoleCannotModify = NewError(1024, "系统内置角色不允许修改")
	ErrRoleMenuIdsEmpty       = NewError(1025, "角色菜单ID列表不能为空")
)
