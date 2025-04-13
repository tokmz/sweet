package auth

import (
	"sweet/pkg/errs"
)

// 定义错误常量
var (
	// ErrTokenExpired Token已过期
	ErrTokenExpired = errs.New(2001, "token已过期")
	// ErrTokenInvalid Token无效
	ErrTokenInvalid = errs.New(2002, "token无效")
	// ErrTokenNotFound Token不存在
	ErrTokenNotFound = errs.New(2003, "token不存在")
	// ErrUserNotFound 用户不存在
	ErrUserNotFound = errs.New(2004, "用户不存在")
	// ErrPermissionDenied 权限不足
	ErrPermissionDenied = errs.New(2005, "权限不足")
	// ErrLoginFailed 登录失败
	ErrLoginFailed = errs.New(2006, "登录失败")
	// ErrUserDisabled 用户已禁用
	ErrUserDisabled = errs.New(2007, "用户已禁用")
	// ErrUserLocked 用户已锁定
	ErrUserLocked = errs.New(2008, "用户已锁定")
	// ErrSecondAuthRequired 需要二级认证
	ErrSecondAuthRequired = errs.New(2009, "需要二级认证")
	// ErrSecondAuthFailed 二级认证失败
	ErrSecondAuthFailed = errs.New(2010, "二级认证失败")
	// ErrInvalidConfig 配置无效
	ErrInvalidConfig = errs.New(2011, "配置无效")
	// ErrStorageError 存储错误
	ErrStorageError = errs.New(2012, "存储错误")
	// ErrCasbinError Casbin错误
	ErrCasbinError = errs.New(2013, "Casbin错误")
	// ErrConcurrencyLimit 并发限制
	ErrConcurrencyLimit = errs.New(2014, "超出并发限制")
	// ErrNotImplemented 未实现
	ErrNotImplemented = errs.New(2015, "功能未实现")
	// ErrInvalidParams 参数无效
	ErrInvalidParams = errs.New(2016, "参数无效")
	// ErrRefreshTokenInvalid 刷新Token无效
	ErrRefreshTokenInvalid = errs.New(2017, "刷新Token无效")
	// ErrRefreshTokenExpired 刷新Token已过期
	ErrRefreshTokenExpired = errs.New(2018, "刷新Token已过期")
	// ErrDeviceNotAllowed 设备不允许
	ErrDeviceNotAllowed = errs.New(2019, "设备不允许")
	// ErrSSONotEnabled 未启用单点登录
	ErrSSONotEnabled = errs.New(2020, "未启用单点登录")
	// ErrRoleNotExists 角色不存在
	ErrRoleNotExists = errs.New(2021, "角色不存在")
	// ErrRoleExists 角色已存在
	ErrRoleExists = errs.New(2022, "角色已存在")
	// ErrUserRoleExists 用户角色已存在
	ErrUserRoleExists = errs.New(2023, "用户角色已存在")
)

