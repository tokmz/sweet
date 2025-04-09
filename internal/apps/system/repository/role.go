package repository

import (
	"context"
	"sweet/pkg/errs"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrRoleNotFound   = errs.New(4000, "角色不存在")
	ErrRoleExists     = errs.New(4001, "角色已存在")
	ErrRoleHasUsers   = errs.New(4002, "角色下存在用户，无法删除")
	ErrInvalidRoleID  = errs.New(4003, "无效角色ID")
	ErrRoleCodeExists = errs.New(4004, "角色编码已存在")
)

type RoleRepository interface {
	// Create 创建角色
	Create(ctx context.Context, role *entity.Role) error
	// Update 更新角色
	Update(ctx context.Context, role *entity.Role) error
	// Delete 删除角色
	Delete(ctx context.Context, ids []int64) error
	// FindOne 查询角色
	FindOne(ctx context.Context, id int64) (*entity.Role, error)
	// FindList 查询角色列表
	FindList(ctx context.Context, params *RoleListParams) (list []*entity.Role, total int64, err error)
	// ScanOne 查询角色
	ScanOne(ctx context.Context, id int64, val any) error
	// ScanList 查询角色列表
	ScanList(ctx context.Context, params *RoleListParams, list any, total *int64) error
	// CheckHasUsers 检查角色是否有用户绑定
	// 返回用户ID列表
	CheckHasUsers(ctx context.Context, roleID int64) ([]int64, error)

	// AssignMenus 给角色分配菜单
	AssignMenus(ctx context.Context, roleID int64, menuIds []int64) error
	// FindMenusTree 查询角色关联的菜单树
	FindMenusTree(ctx context.Context, roleID int64) ([]*entity.Menu, error)
	// FindMenusIds 查询角色关联菜单ids
	FindMenusIds(ctx context.Context, roleID int64) ([]int64, error)
}

// RoleListParams 角色列表查询参数
type RoleListParams struct {
	Name   string // 角色名称
	Code   string // 角色编码
	Status *int64 // 状态
	Page   int    // 页码
	Size   int    // 每页数量
}

// roleRepository 角色仓储实现
type roleRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) Delete(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) FindOne(ctx context.Context, id int64) (*entity.Role, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) FindList(ctx context.Context, params *RoleListParams) (list []*entity.Role, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) ScanOne(ctx context.Context, id int64, val any) error {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) ScanList(ctx context.Context, params *RoleListParams, list any, total *int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) CheckHasUsers(ctx context.Context, roleID int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) AssignMenus(ctx context.Context, roleID int64, menuIds []int64) error {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) FindMenusTree(ctx context.Context, roleID int64) ([]*entity.Menu, error) {
	//TODO implement me
	panic("implement me")
}

func (r *roleRepository) FindMenusIds(ctx context.Context, roleID int64) ([]int64, error) {
	//TODO implement me
	panic("implement me")
}

// NewRoleRepository 创建角色仓储
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
		q:  query.Use(db),
	}
}
