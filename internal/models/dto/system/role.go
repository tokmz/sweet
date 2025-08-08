package system

import (
	"sweet/internal/models"
	"time"
)

// CreateRoleReq 创建角色
type CreateRoleReq struct {
	Name    string  `json:"name"`     // 角色名称
	Code    string  `json:"code"`     // 角色标识
	Sort    int64   `json:"sort"`     // 排序
	IsSuper *int64  `json:"is_super"` // 是否超级管理员：1=是，2否
	Status  *int64  `json:"status"`   // 状态：1=正常，2=禁用
	Remark  *string `json:"remark"`   // 备注
}

// DeleteRoleReq 删除角色
type DeleteRoleReq models.IdsReq

// UpdateRoleReq 更新角色
type UpdateRoleReq struct {
	models.IDReq
	Name    string  `json:"name"`     // 角色名称
	Code    string  `json:"code"`     // 角色标识
	Sort    int64   `json:"sort"`     // 排序
	IsSuper *int64  `json:"is_super"` // 是否超级管理员：1=是，2否
	Status  *int64  `json:"status"`   // 状态：1=正常，2=禁用
	Remark  *string `json:"remark"`   // 备注
}

// RoleListReq 角色列表
type RoleListReq struct {
	Name     string `json:"name"`
	IsSystem *int64 `json:"is_system"` // 是否系统内置：1=是，2否
	IsSuper  *int64 `json:"is_super"`  // 是否超级管理员：1=是，2否
	Status   *int64 `json:"status"`    // 状态：1=正常，2=禁用
	models.TimeRangeReq
}

// RoleListItem 角色列表项
type RoleListItem struct {
	ID        int64      `json:"id"`         // 角色ID
	Name      string     `json:"name"`       // 角色名称
	Code      string     `json:"code"`       // 角色标识
	Sort      int64      `json:"sort"`       // 排序
	IsSystem  *int64     `json:"is_system"`  // 是否系统内置：1=是，2否
	IsSuper   *int64     `json:"is_super"`   // 是否超级管理员：1=是，2否
	Status    *int64     `json:"status"`     // 状态：1=正常，2=禁用
	CreatedAt *time.Time `json:"created_at"` // 创建时间
}

// RoleListRes 角色列表响应
type RoleListRes models.PageRes[RoleListItem]

// RoleDetailRes 角色详情响应
type RoleDetailRes struct {
	ID        int64      `json:"id"`         // 角色ID
	Name      string     `json:"name"`       // 角色名称
	Code      string     `json:"code"`       // 角色标识
	Sort      int64      `json:"sort"`       // 排序
	IsSystem  *int64     `json:"is_system"`  // 是否系统内置：1=是，2否
	IsSuper   *int64     `json:"is_super"`   // 是否超级管理员：1=是，2否
	Status    *int64     `json:"status"`     // 状态：1=正常，2=禁用
	Remark    *string    `json:"remark"`     // 备注
	CreatedAt *time.Time `json:"created_at"` // 创建时间
	UpdatedAt *time.Time `json:"updated_at"` // 更新时间
}

// RoleOptionItem 角色选项响应Item
type RoleOptionItem struct {
	ID   int64  `json:"id"`   // 角色ID
	Name string `json:"name"` // 角色名称
}

// RoleOptionRes 角色选项响应
type RoleOptionRes models.PageRes[RoleOptionItem]

// AssignRoleMenuIdsReq 给角色分配菜单
type AssignRoleMenuIdsReq struct {
	models.IDReq         // 角色ID
	MenuIds      []int64 `json:"menu_ids"` // 菜单ID列表
}

// AssignRoleApiIdsReq 给角色分配Api
type AssignRoleApiIdsReq struct {
	models.IDReq         // 角色ID
	ApiIds       []int64 `json:"api_ids"` // ApiID列表
}
