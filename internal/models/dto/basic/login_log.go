package basic

import (
	"time"
)

// LoginLogCreateReq 创建登录日志请求
type LoginLogCreateReq struct {
	UserID        *int64  `json:"user_id" validate:"omitempty,min=1" comment:"用户ID"`
	Username      string  `json:"username" validate:"required,max=64" comment:"登录用户名"`
	LoginType     *int64  `json:"login_type" validate:"omitempty,oneof=1 2 3 4 5 6 7" comment:"登录类型（1账号密码 2手机验证码 3邮箱验证码 4第三方登录 5微信 6QQ 7支付宝）"`
	ClientType    *int64  `json:"client_type" validate:"omitempty,oneof=1 2 3 4 5" comment:"客户端类型（1Web 2移动端 3小程序 4API 5管理后台）"`
	IP            string  `json:"ip" validate:"required,ip" comment:"登录IP"`
	Location      *string `json:"location" validate:"omitempty,max=64" comment:"IP归属地"`
	UserAgent     *string `json:"user_agent" validate:"omitempty" comment:"用户代理"`
	DeviceInfo    *string `json:"device_info" validate:"omitempty,max=64" comment:"设备信息"`
	Browser       *string `json:"browser" validate:"omitempty,max=100" comment:"浏览器"`
	Os            *string `json:"os" validate:"omitempty,max=100" comment:"操作系统"`
	Status        *int64  `json:"status" validate:"omitempty,oneof=1 2 3" comment:"登录状态（1成功 2失败 3异常）"`
	FailReason    *string `json:"fail_reason" validate:"omitempty,max=64" comment:"失败原因"`
	SessionID     *string `json:"session_id" validate:"omitempty,max=128" comment:"会话ID"`
	RiskLevel     *int64  `json:"risk_level" validate:"omitempty,oneof=1 2 3" comment:"风险等级（1低风险 2中风险 3高风险）"`
}

// LoginLogUpdateReq 更新登录日志请求
type LoginLogUpdateReq struct {
	ID            int64   `json:"id" validate:"required,min=1" comment:"日志ID"`
	LoginDuration *int64  `json:"login_duration" validate:"omitempty,min=0" comment:"登录持续时间（秒）"`
	LogoutType    *int64  `json:"logout_type" validate:"omitempty,oneof=1 2 3" comment:"退出类型（1主动退出 2超时退出 3强制退出）"`
	Status        *int64  `json:"status" validate:"omitempty,oneof=1 2 3" comment:"登录状态（1成功 2失败 3异常）"`
	FailReason    *string `json:"fail_reason" validate:"omitempty,max=64" comment:"失败原因"`
}

// LoginLogQueryReq 查询登录日志请求
type LoginLogQueryReq struct {
	UserID     *int64     `json:"user_id" form:"user_id" validate:"omitempty,min=1" comment:"用户ID"`
	Username   *string    `json:"username" form:"username" validate:"omitempty,max=64" comment:"登录用户名"`
	LoginType  *int64     `json:"login_type" form:"login_type" validate:"omitempty,oneof=1 2 3 4 5 6 7" comment:"登录类型"`
	ClientType *int64     `json:"client_type" form:"client_type" validate:"omitempty,oneof=1 2 3 4 5" comment:"客户端类型"`
	IP         *string    `json:"ip" form:"ip" validate:"omitempty,ip" comment:"登录IP"`
	Status     *int64     `json:"status" form:"status" validate:"omitempty,oneof=1 2 3" comment:"登录状态"`
	RiskLevel  *int64     `json:"risk_level" form:"risk_level" validate:"omitempty,oneof=1 2 3" comment:"风险等级"`
	StartTime  *time.Time `json:"start_time" form:"start_time" validate:"omitempty" comment:"开始时间"`
	EndTime    *time.Time `json:"end_time" form:"end_time" validate:"omitempty" comment:"结束时间"`
	Page       int        `json:"page" form:"page" validate:"min=1" comment:"页码"`
	PageSize   int        `json:"page_size" form:"page_size" validate:"min=1,max=100" comment:"每页数量"`
}

// LoginLogResp 登录日志响应
type LoginLogResp struct {
	ID            int64      `json:"id" comment:"日志ID"`
	UserID        *int64     `json:"user_id" comment:"用户ID"`
	Username      string     `json:"username" comment:"登录用户名"`
	LoginType     *int64     `json:"login_type" comment:"登录类型"`
	LoginTypeName string     `json:"login_type_name" comment:"登录类型名称"`
	ClientType    *int64     `json:"client_type" comment:"客户端类型"`
	ClientTypeName string    `json:"client_type_name" comment:"客户端类型名称"`
	IP            string     `json:"ip" comment:"登录IP"`
	Location      *string    `json:"location" comment:"IP归属地"`
	UserAgent     *string    `json:"user_agent" comment:"用户代理"`
	DeviceInfo    *string    `json:"device_info" comment:"设备信息"`
	Browser       *string    `json:"browser" comment:"浏览器"`
	Os            *string    `json:"os" comment:"操作系统"`
	Status        *int64     `json:"status" comment:"登录状态"`
	StatusName    string     `json:"status_name" comment:"登录状态名称"`
	FailReason    *string    `json:"fail_reason" comment:"失败原因"`
	SessionID     *string    `json:"session_id" comment:"会话ID"`
	LoginDuration *int64     `json:"login_duration" comment:"登录持续时间（秒）"`
	LogoutType    *int64     `json:"logout_type" comment:"退出类型"`
	LogoutTypeName string    `json:"logout_type_name" comment:"退出类型名称"`
	RiskLevel     *int64     `json:"risk_level" comment:"风险等级"`
	RiskLevelName string     `json:"risk_level_name" comment:"风险等级名称"`
	CreatedAt     *time.Time `json:"created_at" comment:"创建时间"`
	UpdatedAt     *time.Time `json:"updated_at" comment:"更新时间"`
	// 关联用户信息
	User *UserSimpleResp `json:"user,omitempty" comment:"用户信息"`
}

// UserSimpleResp 用户简单信息响应
type UserSimpleResp struct {
	ID       int64   `json:"id" comment:"用户ID"`
	Username string  `json:"username" comment:"用户名"`
	Realname string  `json:"realname" comment:"真实姓名"`
	Nickname string  `json:"nickname" comment:"昵称"`
	Avatar   *string `json:"avatar" comment:"头像"`
}

// LoginLogListResp 登录日志列表响应
type LoginLogListResp struct {
	List  []*LoginLogResp `json:"list" comment:"日志列表"`
	Total int64           `json:"total" comment:"总数"`
	Page  int             `json:"page" comment:"当前页"`
	Size  int             `json:"size" comment:"每页数量"`
}

// LoginLogStatsResp 登录日志统计响应
type LoginLogStatsResp struct {
	TotalCount    int64 `json:"total_count" comment:"总登录次数"`
	SuccessCount  int64 `json:"success_count" comment:"成功登录次数"`
	FailCount     int64 `json:"fail_count" comment:"失败登录次数"`
	TodayCount    int64 `json:"today_count" comment:"今日登录次数"`
	OnlineCount   int64 `json:"online_count" comment:"在线用户数"`
	HighRiskCount int64 `json:"high_risk_count" comment:"高风险登录次数"`
}