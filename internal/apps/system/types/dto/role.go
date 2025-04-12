package dto

import "sweet/internal/common"

// CreateRoleReq 创建角色请求
type CreateRoleReq struct {
	Name      string  `json:"name"`       // 角色名称
	Code      string  `json:"code"`       // 角色编码
	Sort      *int    `json:"sort"`       // 排序
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	DataScope *int64  `json:"data_scope"` // 数据范围 1-全部 2-自定义 3-本部门 4-本部门及子部门 5-仅本人
	Remark    *string `json:"remark"`     // 备注
	CreatedBy *int64  `json:"created_by"` // 创建人
}

// UpdateRoleReq 更新角色请求
type UpdateRoleReq struct {
	common.IDReq
	Name      string  `json:"name"`       // 角色名称
	Code      string  `json:"code"`       // 角色编码
	Sort      *int    `json:"sort"`       // 排序
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	DataScope *int64  `json:"data_scope"` // 数据范围 1-全部 2-自定义 3-本部门 4-本部门及子部门 5-仅本人
	Remark    *string `json:"remark"`     // 备注
	UpdatedBy *int64  `json:"updated_by"` // 更新人
}

// FindListRoleReq 查询角色列表请求
type FindListRoleReq struct {
	Name      string `json:"name"`       // 角色名称
	Code      string `json:"code"`       // 角色编码
	Status    int64  `json:"status"`     // 状态
	DataScope int64  `json:"data_scope"` // 数据范围
	common.PageReq
}

// AssignMenusReq 给角色分配菜单请求
type AssignMenusReq struct {
	common.IDReq         // 角色ID
	MenuIds      []int64 `json:"menu_ids"` // 菜单ID列表
}

// UpdateRoleStatusReq 更新角色状态请求
type UpdateRoleStatusReq struct {
	common.StatusReq
}
