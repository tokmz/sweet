package dto

import "sweet/internal/common"

// CreateDeptReq 创建部门请求
type CreateDeptReq struct {
	Name      string  `json:"name"`       // 部门名称
	Code      *string `json:"code"`       // 部门编码
	ParentID  *int64  `json:"parent_id"`  // 父部门ID
	Ancestors *string `json:"ancestors"`  // 祖级列表
	Leader    *string `json:"leader"`     // 负责人
	Phone     *string `json:"phone"`      // 联系电话
	Email     *string `json:"email"`      // 邮箱
	Sort      *int    `json:"sort"`       // 排序
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	CreatedBy *int64  `json:"created_by"` // 创建人
}

// UpdateDeptReq 更新部门请求
type UpdateDeptReq struct {
	common.IDReq
	Name      string  `json:"name"`       // 部门名称
	Code      *string `json:"code"`       // 部门编码
	ParentID  *int64  `json:"parent_id"`  // 父部门ID
	Ancestors *string `json:"ancestors"`  // 祖级列表
	Leader    *string `json:"leader"`     // 负责人
	Phone     *string `json:"phone"`      // 联系电话
	Email     *string `json:"email"`      // 邮箱
	Sort      *int    `json:"sort"`       // 排序
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	UpdatedBy *int64  `json:"updated_by"` // 更新人
}

// FindListDeptReq 查询部门列表请求
type FindListDeptReq struct {
	Name      string `json:"name"`       // 部门名称
	Code      string `json:"code"`       // 部门编码
	Status    *int64 `json:"status"`     // 状态
	ParentID  *int64 `json:"parent_id"`  // 父部门ID
	ExcludeID *int64 `json:"exclude_id"` // 排除ID
	common.PageReq
}

// DeptTreeReq 查询部门树请求
type DeptTreeReq struct {
	Name      string `json:"name"`       // 部门名称
	Code      string `json:"code"`       // 部门编码
	Status    *int64 `json:"status"`     // 状态
	ParentID  *int64 `json:"parent_id"`  // 父部门ID
	ExcludeID *int64 `json:"exclude_id"` // 排除ID
}

// SubDeptReq 查询子部门请求
type SubDeptReq struct {
	common.IDReq // 部门ID
}
