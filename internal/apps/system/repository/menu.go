package repository

import (
	"context"
	"sweet/pkg/errs"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrMenuNotFound         = errs.New(5000, "菜单不存在")
	ErrMenuExists           = errs.New(5001, "菜单已存在")
	ErrMenuHasChildren      = errs.New(5002, "存在子菜单，无法删除")
	ErrInvalidMenuID        = errs.New(5003, "无效菜单ID")
	ErrMenuPermissionExists = errs.New(5004, "权限标识已存在")
)

type MenuRepository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *entity.Menu) error
	// Update 更新菜单
	Update(ctx context.Context, menu *entity.Menu) error
	// Delete 删除菜单
	Delete(ctx context.Context, ids []int64) error
	// FindOne 查询菜单
	FindOne(ctx context.Context, id int64) (*entity.Menu, error)
	// ScanOne 查询菜单
	ScanOne(ctx context.Context, id int64, val any) error
	// Tree 查询菜单树
	Tree(ctx context.Context, params *TreeParams) ([]*entity.Menu, int64, error)
	// RouteTree 查询路由树 根据角色ID
	RouteTree(ctx context.Context, rid int64) ([]*RouteTree, int64, error)
	// ItemTree 查询Item菜单树 用于分配给角色
	ItemTree(ctx context.Context, rid int64) ([]*ItemTree, int64, error)
}

// TreeParams 菜单列表查询参数
type TreeParams struct {
	Rid  int64  // 角色ID 获取角色的菜单树
	Pid  int64  // 父菜单ID 获取此ID和子菜单
	Name string // 菜单名称 模糊查询 同时获取子菜单内容
}

type RouteTree struct {
}

type ItemTree struct {
}

// menuRepository 菜单仓储实现
type menuRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (m *menuRepository) Create(ctx context.Context, menu *entity.Menu) error {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) Update(ctx context.Context, menu *entity.Menu) error {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) Delete(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindOne(ctx context.Context, id int64) (*entity.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) ScanOne(ctx context.Context, id int64, val any) error {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) Tree(ctx context.Context, params *TreeParams) ([]*entity.Menu, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) RouteTree(ctx context.Context, rid int64) ([]*RouteTree, int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) ItemTree(ctx context.Context, rid int64) ([]*ItemTree, int64, error) {
	//TODO implement me
	panic("implement me")
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{
		db: db,
		q:  query.Use(db),
	}
}
