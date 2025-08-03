package basic

import (
	"sweet/internal/models"
	"time"
)

// CreateOperationLogReq 创建操作日志请求
type CreateOperationLogReq struct {
	UserID    *int64  `json:"user_id"`    // 操作用户ID
	Username  *string `json:"username"`   // 用户账号
	Module    string  `json:"module"`     // 操作模块
	Operation string  `json:"operation"`  // 操作类型
	Method    string  `json:"method"`     // HTTP方法
	URL       string  `json:"url"`        // 请求URL
	IP        string  `json:"ip"`         // 操作IP
	Location  *string `json:"location"`   // IP归属地
	UserAgent *string `json:"user_agent"` // 用户代理
	Params    *string `json:"params"`     // 请求参数
	Result    *string `json:"result"`     // 操作结果
	Status    *int64  `json:"status"`     // 操作状态（1成功 2失败）
	ErrorMsg  *string `json:"error_msg"`  // 错误信息
	Duration  *int64  `json:"duration"`   // 执行时长（毫秒）
}

// DeleteOperationLogReq 删除操作日志请求
type DeleteOperationLogReq struct {
	Uid           int64 `json:"uid"` // 用户ID （可选 根据用户ID删除）
	models.IdsReq       // 删除操作日志ID列表
}

// ListOperationLogReq 获取操作日志列表请求
type ListOperationLogReq struct {
	Uid       int64  `json:"uid"`       // 用户ID （可选 根据用户ID查询）
	Module    string `json:"module"`    // 操作模块 （可选 根据模块查询）
	Operation string `json:"operation"` // 操作类型 （可选 根据操作类型查询）
	Method    string `json:"method"`    // HTTP方法 （可选 根据方法查询）
	Status    int64  `json:"status"`    // 操作状态 （可选 根据状态查询）
	models.TimeRangeReq
	models.PageReq // 分页参数
	models.SortReq // 排序参数
}

// ListOperationLogItem 获取操作日志列表项
type ListOperationLogItem struct {
	ID        int64      `json:"id"`         // 日志ID
	UserID    *int64     `json:"user_id"`    // 操作用户ID
	Username  *string    `json:"username"`   // 用户账号
	Module    string     `json:"module"`     // 操作模块
	Operation string     `json:"operation"`  // 操作类型
	Method    string     `json:"method"`     // HTTP方法
	URL       string     `json:"url"`        // 请求URL
	IP        string     `json:"ip"`         // 操作IP
	Status    *int64     `json:"status"`     // 操作状态（1成功 2失败）
	Duration  *int64     `json:"duration"`   // 执行时长（毫秒）
	CreatedAt *time.Time `json:"created_at"` // 创建时间
}

// ListOperationLogRes 获取操作日志列表响应
type ListOperationLogRes models.PageRes[ListOperationLogItem]

// OperationLogDetailRes 获取操作日志详情响应
type OperationLogDetailRes struct {
	ID        int64      `json:"id"`         // 日志ID
	UserID    *int64     `json:"user_id"`    // 操作用户ID
	Username  string     `json:"username"`   // 操作用户名
	Realname  string     `json:"realname"`   // 真实姓名
	Nickname  string     `json:"nickname"`   // 昵称
	Avatar    *string    `json:"avatar"`     // 头像
	Email     *string    `json:"email"`      // 邮箱
	Phone     *string    `json:"phone"`      // 手机号
	Module    string     `json:"module"`     // 操作模块
	Operation string     `json:"operation"`  // 操作类型
	Method    string     `json:"method"`     // HTTP方法
	URL       string     `json:"url"`        // 请求URL
	IP        string     `json:"ip"`         // 操作IP
	Location  *string    `json:"location"`   // IP归属地
	UserAgent *string    `json:"user_agent"` // 用户代理
	Params    *string    `json:"params"`     // 请求参数
	Result    *string    `json:"result"`     // 操作结果
	Status    *int64     `json:"status"`     // 操作状态（1成功 2失败）
	ErrorMsg  *string    `json:"error_msg"`  // 错误信息
	Duration  *int64     `json:"duration"`   // 执行时长（毫秒）
	CreatedAt *time.Time `json:"created_at"` // 创建时间
}
