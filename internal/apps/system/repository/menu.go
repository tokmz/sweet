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

// MenuRepository 菜单仓储接口
type MenuRepository interface {
	// Create 创建菜单
	Create(ctx context.Context, menu *entity.Menu) error
	// Update 更新菜单
	Update(ctx context.Context, menu *entity.Menu) error
	// Delete 删除菜单
	Delete(ctx context.Context, ids []int64) error
	// FindOne 查询菜单
	FindOne(ctx context.Context, id int64) (*entity.Menu, error)
	// FindByPermission 通过权限标识查询菜单
	FindByPermission(ctx context.Context, permission string) (*entity.Menu, error)
	// FindChildren 查询子菜单
	FindChildren(ctx context.Context, parentID int64) ([]*entity.Menu, error)
	// FindChildrenIDs 查询子菜单ID列表
	FindChildrenIDs(ctx context.Context, parentID int64) ([]int64, error)
	// FindList 查询菜单列表
	FindList(ctx context.Context, params *MenuListParams) ([]*entity.Menu, error)
	// FindTree 查询菜单树
	FindTree(ctx context.Context, params *MenuListParams) ([]*entity.MenuTree, error)
	// FindUserMenus 查询用户菜单列表
	FindUserMenus(ctx context.Context, userID int64) ([]*entity.Menu, error)
	// FindUserMenuTree 查询用户菜单树
	FindUserMenuTree(ctx context.Context, userID int64) ([]*entity.MenuTree, error)
	// FindUserPermissions 查询用户权限列表
	FindUserPermissions(ctx context.Context, userID int64) ([]string, error)
}

// MenuListParams 菜单列表查询参数
type MenuListParams struct {
	Name       string // 菜单名称
	Permission string // 权限标识
	Status     *int64 // 状态
	Type       *int64 // 菜单类型
	ParentID   *int64 // 父菜单ID
	ExcludeID  *int64 // 排除ID
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

func (m *menuRepository) FindByPermission(ctx context.Context, permission string) (*entity.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindChildren(ctx context.Context, parentID int64) ([]*entity.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindChildrenIDs(ctx context.Context, parentID int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindList(ctx context.Context, params *MenuListParams) ([]*entity.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindTree(ctx context.Context, params *MenuListParams) ([]*entity.MenuTree, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindUserMenus(ctx context.Context, userID int64) ([]*entity.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindUserMenuTree(ctx context.Context, userID int64) ([]*entity.MenuTree, error) {
	//TODO implement me
	panic("implement me")
}

func (m *menuRepository) FindUserPermissions(ctx context.Context, userID int64) ([]string, error) {
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
