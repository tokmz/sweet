package vo

import "time"

// DeptTreeRes 部门树响应
type DeptTreeRes struct {
	ID        int64          `json:"id"`         // 部门ID
	ParentID  *int64         `json:"parent_id"`  // 父部门ID
	Name      string         `json:"name"`       // 部门名称
	Code      *string        `json:"code"`       // 部门编码
	Ancestors *string        `json:"ancestors"`  // 祖级列表
	Leader    *string        `json:"leader"`     // 负责人
	Phone     *string        `json:"phone"`      // 联系电话
	Email     *string        `json:"email"`      // 邮箱
	Sort      *int64         `json:"sort"`       // 排序
	Status    *int64         `json:"status"`     // 状态 1-正常 2-禁用
	CreatedAt time.Time      `json:"created_at"` // 创建时间
	Children  []*DeptTreeRes `json:"children"`   // 子部门
}

// DeptDetailRes 部门详情响应
type DeptDetailRes struct {
	ID        int64      `json:"id"`         // 部门ID
	ParentID  *int64     `json:"parent_id"`  // 父部门ID
	Name      string     `json:"name"`       // 部门名称
	Code      *string    `json:"code"`       // 部门编码
	Ancestors *string    `json:"ancestors"`  // 祖级列表
	Leader    *string    `json:"leader"`     // 负责人
	Phone     *string    `json:"phone"`      // 联系电话
	Email     *string    `json:"email"`      // 邮箱
	Sort      *int64     `json:"sort"`       // 排序
	Status    *int64     `json:"status"`     // 状态 1-正常 2-禁用
	CreatedBy *int64     `json:"created_by"` // 创建人
	UpdatedBy *int64     `json:"updated_by"` // 更新人
	CreatedAt time.Time  `json:"created_at"` // 创建时间
	UpdatedAt *time.Time `json:"updated_at"` // 更新时间
}

// DeptItemRes 部门列表项响应，用于选择器等简化场景
type DeptItemRes struct {
	ID   int64  `json:"id"`   // 部门ID
	Name string `json:"name"` // 部门名称
}
