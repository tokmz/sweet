package basic

import (
	"sweet/internal/models"
	"time"
)

// CreateLoginLogReq 创建登录日志请求
type CreateLoginLogReq struct {
	UserID     *int64  `json:"user_id"`     // 用户ID
	Username   string  `json:"username"`    // 登录用户名
	LoginType  *int64  `json:"login_type"`  // 登录类型（1账号密码 2手机验证码 3邮箱验证码 4第三方登录 5微信 6QQ 7支付宝）
	ClientType *int64  `json:"client_type"` // 客户端类型（1Web 2移动端 3小程序 4API 5管理后台）
	IP         string  `json:"ip"`          // 登录IP
	Location   *string `json:"location"`    // IP归属地
	UserAgent  *string `json:"user_agent"`  // 用户代理
	DeviceInfo *string `json:"device_info"` // 设备信息
	Browser    *string `json:"browser"`     // 浏览器
	Os         *string `json:"os"`          // 操作系统
	Status     *int64  `json:"status"`      // 登录状态（1成功 2失败 3异常）
	FailReason *string `json:"fail_reason"` // 失败原因
}

// DeleteLoginLogReq 删除登录日志请求
type DeleteLoginLogReq struct {
	Uid           int64 `json:"uid"` // 用户ID （可选 根据用户ID删除）
	models.IdsReq       // 删除登录日志ID列表
}

// ListLoginLogReq 获取登录日志列表请求
type ListLoginLogReq struct {
	Uid                 int64               `json:"uid"`         // 用户ID （可选 根据用户ID查询）
	LoginType           int64               `json:"login_type"`  // 登录类型 （可选 根据登录类型查询）
	ClientType          int64               `json:"client_type"` // 客户端类型 （可选 根据客户端类型查询）
	Status              int64               `json:"status"`      // 登录状态 （可选 根据登录状态查询）
	models.PageReq      `json:"page"`       // 分页参数
	models.SortReq      `json:"sort"`       // 排序参数
	models.TimeRangeReq `json:"time_range"` // 时间范围参数
}

// ListLoginLogItem 获取登录日志列表项
type ListLoginLogItem struct {
	ID         int64      `json:"id"`          // 日志ID
	UserID     *int64     `json:"user_id"`     // 用户ID
	Username   string     `json:"username"`    // 登录用户名
	LoginType  *int64     `json:"login_type"`  // 登录类型（1账号密码 2手机验证码 3邮箱验证码 4第三方登录 5微信 6QQ 7支付宝）
	ClientType *int64     `json:"client_type"` // 客户端类型（1Web 2移动端 3小程序 4API 5管理后台）
	Browser    *string    `json:"browser"`     // 浏览器
	Os         *string    `json:"os"`          // 操作系统
	Status     *int64     `json:"status"`      // 登录状态（1成功 2失败 3异常）
	CreatedAt  *time.Time `json:"created_at"`  // 创建时间
}

// ListLoginLogRes 获取登录日志列表响应
type ListLoginLogRes models.PageRes[ListLoginLogItem]

// LoginLogDetailRes 获取登录日志详情响应
type LoginLogDetailRes struct {
	ID         int64      `json:"id"`          // 日志ID
	UserID     *int64     `json:"user_id"`     // 用户ID
	Username   string     `json:"username"`    // 登录用户名
	Realname   string     `json:"realname"`    // 真实姓名
	Nickname   string     `json:"nickname"`    // 昵称
	Avatar     *string    `json:"avatar"`      // 头像
	Email      *string    `json:"email"`       // 邮箱
	Phone      *string    `json:"phone"`       // 手机号
	LoginType  *int64     `json:"login_type"`  // 登录类型（1账号密码 2手机验证码 3邮箱验证码 4第三方登录 5微信 6QQ 7支付宝）
	ClientType *int64     `json:"client_type"` // 客户端类型（1Web 2移动端 3小程序 4API 5管理后台）
	IP         string     `json:"ip"`          // 登录IP
	Location   *string    `json:"location"`    // IP归属地
	UserAgent  *string    `json:"user_agent"`  // 用户代理
	DeviceInfo *string    `json:"device_info"` // 设备信息
	Browser    *string    `json:"browser"`     // 浏览器
	Os         *string    `json:"os"`          // 操作系统
	Status     *int64     `json:"status"`      // 登录状态（1成功 2失败 3异常）
	FailReason *string    `json:"fail_reason"` // 失败原因
	CreatedAt  *time.Time `json:"created_at"`  // 创建时间
}
