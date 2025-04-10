package repository

import (
	"context"
	"sweet/pkg/errs"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrDeptNotFound    = errs.New(3000, "部门不存在")
	ErrDeptExists      = errs.New(3001, "部门已存在")
	ErrDeptHasChildren = errs.New(3002, "存在子部门，无法删除")
	ErrDeptHasUsers    = errs.New(3003, "部门下存在用户，无法删除")
	ErrInvalidDeptID   = errs.New(3004, "无效部门ID")
	ErrDeptCodeExists  = errs.New(3005, "部门编码已存在")
)

type DeptRepository interface {
	// Create 创建部门
	Create(ctx context.Context, dept *entity.Dept) error
	// Update 更新部门
	Update(ctx context.Context, dept *entity.Dept) error
	// Delete 删除部门
	Delete(ctx context.Context, ids []int64) error
	// FindOne 查询部门
	FindOne(ctx context.Context, id int64) (*entity.Dept, error)
	// FindList 查询部门列表
	FindList(ctx context.Context, params *DeptListParams) ([]*entity.Dept, error)
	// ScanOne 查询部门
	ScanOne(ctx context.Context, id int64, val any) error
	// ScanList 查询部门列表
	ScanList(ctx context.Context, params *DeptListParams, list any) error
	// SubDept 查询子部门
	SubDept(ctx context.Context, id int64) (*entity.Dept, error)
	// Tree 查询部门树
	Tree(ctx context.Context, params *DeptListParams) ([]*entity.Dept, error)

	// DeptUser 获取部门下用户，如果有子部门，递归查询子部门下用户
	DeptUser(ctx context.Context, params DeptUserParams) ([]*entity.User, error)
}

// DeptListParams 部门列表查询参数
type DeptListParams struct {
	Name      string // 部门名称
	Code      string // 部门编码
	Status    *int64 // 状态
	ParentID  *int64 // 父部门ID
	ExcludeID *int64 // 排除ID
	Page      int    // 页码
	Size      int    // 每页数量
}

type DeptUserParams struct {
	ID   int64
	Page int
	Size int
}

// deptRepository 部门仓储实现
type deptRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (d *deptRepository) Create(ctx context.Context, dept *entity.Dept) error {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) Update(ctx context.Context, dept *entity.Dept) error {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) Delete(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) FindOne(ctx context.Context, id int64) (*entity.Dept, error) {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) FindList(ctx context.Context, params *DeptListParams) ([]*entity.Dept, error) {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) ScanOne(ctx context.Context, id int64, val any) error {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) ScanList(ctx context.Context, params *DeptListParams, list any) error {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) SubDept(ctx context.Context, id int64) (*entity.Dept, error) {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) Tree(ctx context.Context, params *DeptListParams) ([]*entity.Dept, error) {
	//TODO implement me
	panic("implement me")
}

func (d *deptRepository) DeptUser(ctx context.Context, params DeptUserParams) ([]*entity.User, error) {
	//TODO implement me
	panic("implement me")
}

// NewDeptRepository 创建部门仓储
func NewDeptRepository(db *gorm.DB) DeptRepository {
	return &deptRepository{
		db: db,
		q:  query.Use(db),
	}
}
