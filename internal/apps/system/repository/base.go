package repository

import (
	"context"
)

// Repository 基础仓储接口
type Repository[T any, P any] interface {
	// Create 创建实体
	Create(ctx context.Context, entity *T) error
	// Update 更新实体
	Update(ctx context.Context, entity *T) error
	// Delete 删除实体
	Delete(ctx context.Context, ids []int64) error
	// FindOne 查询单个实体
	FindOne(ctx context.Context, id int64) (*T, error)
	// FindList 查询实体列表
	FindList(ctx context.Context, params P) (list []*T, total int64, err error)
	// ScanOne 查询单个实体并扫描到目标对象
	ScanOne(ctx context.Context, id int64, dest any) error
	// ScanList 查询实体列表并扫描到目标对象
	ScanList(ctx context.Context, params P, list any, total *int64) error
}

// Pageable 分页参数
type Pageable struct {
	Page int // 页码
	Size int // 每页数量
}

// Sortable 排序参数
type Sortable struct {
	OrderBy  string // 排序字段
	OrderDir string // 排序方向 asc/desc
}
