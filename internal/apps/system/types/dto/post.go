package dto

import "sweet/internal/common"

// CreatePostReq 创建岗位请求
type CreatePostReq struct {
	Name      string  `json:"name"`       // 岗位名称
	Code      string  `json:"code"`       // 岗位编码
	Sort      *int    `json:"sort"`       // 排序
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	Remark    *string `json:"remark"`     // 备注
	CreatedBy *int64  `json:"created_by"` // 创建人
}

// UpdatePostReq 更新岗位请求
type UpdatePostReq struct {
	common.IDReq
	Name      string  `json:"name"`       // 岗位名称
	Code      string  `json:"code"`       // 岗位编码
	Sort      *int    `json:"sort"`       // 排序
	Status    *int64  `json:"status"`     // 状态 1-正常 2-禁用
	Remark    *string `json:"remark"`     // 备注
	UpdatedBy *int64  `json:"updated_by"` // 更新人
}

// FindListPostReq 查询岗位列表请求
type FindListPostReq struct {
	Name   string `json:"name"`   // 岗位名称
	Code   string `json:"code"`   // 岗位编码
	Status *int64 `json:"status"` // 状态
	common.PageReq
}
