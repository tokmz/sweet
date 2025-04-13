package vo

import "time"

// PostDetailRes 岗位详情响应
type PostDetailRes struct {
	ID        int64      `json:"id"`         // 岗位ID
	Name      string     `json:"name"`       // 岗位名称
	Code      string     `json:"code"`       // 岗位编码
	Sort      *int64     `json:"sort"`       // 排序
	Status    *int64     `json:"status"`     // 状态 1-正常 2-禁用
	Remark    *string    `json:"remark"`     // 备注
	CreatedBy *int64     `json:"created_by"` // 创建人
	UpdatedBy *int64     `json:"updated_by"` // 更新人
	CreatedAt time.Time  `json:"created_at"` // 创建时间
	UpdatedAt *time.Time `json:"updated_at"` // 更新时间
}

// PostListRes 岗位列表响应
type PostListRes struct {
	ID        int64     `json:"id"`         // 岗位ID
	Name      string    `json:"name"`       // 岗位名称
	Code      string    `json:"code"`       // 岗位编码
	Sort      *int64    `json:"sort"`       // 排序
	Status    *int64    `json:"status"`     // 状态 1-正常 2-禁用
	Remark    *string   `json:"remark"`     // 备注
	CreatedAt time.Time `json:"created_at"` // 创建时间
}

// PostItemRes 岗位列表项响应，用于选择器等简化场景
type PostItemRes struct {
	ID   int64  `json:"id"`   // 岗位ID
	Name string `json:"name"` // 岗位名称
}
