package repo

import (
	"context"
	"errors"
	"sweet/pkg/errs"
	"sweet/pkg/logger"
	"sweet/pkg/utils"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrDeptNotFound    = errs.New(3001, "部门不存在")
	ErrDeptNameExists  = errs.New(3002, "部门名称已存在")
	ErrDeptHasChildren = errs.New(3003, "存在子部门，无法删除")
	ErrDeptHasUsers    = errs.New(3004, "部门下存在用户，无法删除")
	ErrInvalidDeptID   = errs.New(3005, "无效部门ID")
	ErrDeptCodeExists  = errs.New(3006, "部门编码已存在")
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
	return d.q.Transaction(func(tx *query.Query) error {
		do := tx.Dept.WithContext(ctx)
		field := tx.Dept

		// 检查部门名称是否已存在（同级部门下名称不能重复）
		var parentID int64 = 0
		if dept.Pid != nil {
			parentID = *dept.Pid
		}

		if _, err := do.Where(field.Name.Eq(dept.Name), field.Pid.Eq(parentID)).First(); err == nil {
			return ErrDeptNameExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询部门失败",
				logger.Err(err),
				logger.String("部门名称", dept.Name),
				logger.Int64("父部门ID", parentID),
			)
			return errs.ErrServer
		}

		// 检查部门编码是否已存在
		if dept.Code != nil && *dept.Code != "" {
			if _, err := do.Where(field.Code.Eq(*dept.Code)).First(); err == nil {
				return ErrDeptCodeExists
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error(
					"查询部门失败",
					logger.Err(err),
					logger.String("部门编码", *dept.Code),
				)
				return errs.ErrServer
			}
		}

		// 构建祖级列表
		if dept.Pid != nil && *dept.Pid > 0 {
			parent, err := d.FindOne(ctx, *dept.Pid)
			if err != nil {
				return err
			}

			// 设置祖级列表
			ancestors := ""
			if parent.Ancestors != nil {
				ancestors = *parent.Ancestors + "," + utils.ToString(*dept.Pid)
			} else {
				ancestors = utils.ToString(*dept.Pid)
			}
			dept.Ancestors = &ancestors
		}

		// 创建部门
		if err := do.Create(dept); err != nil {
			logger.Error(
				"创建部门失败",
				logger.Err(err),
				logger.String("部门名称", dept.Name),
				logger.String("部门编码", utils.SafeString(dept.Code)),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (d *deptRepository) Update(ctx context.Context, dept *entity.Dept) error {
	return d.q.Transaction(func(tx *query.Query) error {
		do := tx.Dept.WithContext(ctx)
		field := tx.Dept

		// 检查要更新的部门是否存在
		oldDept, err := d.FindOne(ctx, dept.ID)
		if err != nil {
			return err
		}

		// 检查名称是否重复（同级部门下名称不能重复）
		var parentID int64 = 0
		if dept.Pid != nil {
			parentID = *dept.Pid
		}

		if _, err := do.Where(field.Name.Eq(dept.Name), field.Pid.Eq(parentID), field.ID.Neq(dept.ID)).First(); err == nil {
			return ErrDeptNameExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询部门失败",
				logger.Err(err),
				logger.String("部门名称", dept.Name),
				logger.Int64("父部门ID", parentID),
			)
			return errs.ErrServer
		}

		// 检查编码是否重复
		if dept.Code != nil && *dept.Code != "" {
			if _, err := do.Where(field.Code.Eq(*dept.Code), field.ID.Neq(dept.ID)).First(); err == nil {
				return ErrDeptCodeExists
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error(
					"查询部门失败",
					logger.Err(err),
					logger.String("部门编码", *dept.Code),
				)
				return errs.ErrServer
			}
		}

		// 构建祖级列表
		if dept.Pid != nil && *dept.Pid > 0 && (oldDept.Pid == nil || *oldDept.Pid != *dept.Pid) {
			// 不能将部门的父级设置为自己或下级
			if *dept.Pid == dept.ID {
				return errs.New(3007, "不能将部门的父级设置为自己")
			}

			// 检查是否设置为自己的下级
			children, err := d.FindList(ctx, &DeptListParams{ParentID: &dept.ID})
			if err != nil {
				return err
			}

			childrenIDs := make([]int64, 0, len(children))
			for _, child := range children {
				childrenIDs = append(childrenIDs, child.ID)
			}

			if utils.Contains(childrenIDs, *dept.Pid) {
				return errs.New(3008, "不能将部门的父级设置为其下级")
			}

			// 设置新的祖级列表
			parent, err := d.FindOne(ctx, *dept.Pid)
			if err != nil {
				return err
			}

			ancestors := ""
			if parent.Ancestors != nil {
				ancestors = *parent.Ancestors + "," + utils.ToString(*dept.Pid)
			} else {
				ancestors = utils.ToString(*dept.Pid)
			}
			dept.Ancestors = &ancestors

			// 更新所有下级部门的祖级列表
			if len(childrenIDs) > 0 {
				for _, childID := range childrenIDs {
					child, err := d.FindOne(ctx, childID)
					if err != nil {
						continue
					}

					childAncestors := ancestors + "," + utils.ToString(dept.ID)
					child.Ancestors = &childAncestors

					if _, err := do.Where(field.ID.Eq(child.ID)).Updates(child); err != nil {
						logger.Error(
							"更新子部门祖级列表失败",
							logger.Err(err),
							logger.Int64("子部门ID", child.ID),
						)
					}
				}
			}
		}

		// 更新部门
		if _, err := do.Where(field.ID.Eq(dept.ID)).Updates(dept); err != nil {
			logger.Error(
				"更新部门失败",
				logger.Err(err),
				logger.Int64("部门ID", dept.ID),
				logger.String("部门名称", dept.Name),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (d *deptRepository) Delete(ctx context.Context, ids []int64) error {
	return d.q.Transaction(func(tx *query.Query) error {
		do := tx.Dept.WithContext(ctx)
		userDo := tx.User.WithContext(ctx)

		for _, id := range ids {
			// 检查是否有子部门
			children, err := d.FindList(ctx, &DeptListParams{ParentID: &id})
			if err != nil {
				return err
			}
			if len(children) > 0 {
				return ErrDeptHasChildren
			}

			// 检查部门下是否有用户
			count, err := userDo.Where(tx.User.DeptID.Eq(id)).Count()
			if err != nil {
				logger.Error(
					"查询部门用户数量失败",
					logger.Err(err),
					logger.Int64("部门ID", id),
				)
				return errs.ErrServer
			}
			if count > 0 {
				return ErrDeptHasUsers
			}
		}

		// 删除部门
		if _, err := do.Where(tx.Dept.ID.In(ids...)).Delete(); err != nil {
			logger.Error(
				"删除部门失败",
				logger.Err(err),
				logger.Any("部门ID列表", ids),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (d *deptRepository) FindOne(ctx context.Context, id int64) (*entity.Dept, error) {
	if id <= 0 {
		return nil, ErrInvalidDeptID
	}

	do := d.q.Dept.WithContext(ctx)
	field := d.q.Dept
	dept, err := do.Where(field.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrDeptNotFound
		}
		logger.Error(
			"查询部门失败",
			logger.Err(err),
			logger.Int64("部门ID", id),
		)
		return nil, errs.ErrServer
	}
	return dept, nil
}

func (d *deptRepository) FindList(ctx context.Context, params *DeptListParams) ([]*entity.Dept, error) {
	do := d.q.Dept.WithContext(ctx)
	field := d.q.Dept

	// 构建查询条件
	if params.Name != "" {
		do = do.Where(field.Name.Like("%" + params.Name + "%"))
	}
	if params.Code != "" {
		do = do.Where(field.Code.Like("%" + params.Code + "%"))
	}
	if params.Status != nil {
		do = do.Where(field.Status.Eq(*params.Status))
	}
	if params.ParentID != nil {
		do = do.Where(field.Pid.Eq(*params.ParentID))
	}
	if params.ExcludeID != nil {
		do = do.Where(field.ID.Neq(*params.ExcludeID))
	}

	// 默认按照排序字段和ID升序排列
	do = do.Order(field.Sort, field.ID)

	// 分页查询
	var list []*entity.Dept
	var err error
	if params.Page > 0 && params.Size > 0 {
		// 计算偏移量
		offset := (params.Page - 1) * params.Size
		list, err = do.Offset(offset).Limit(params.Size).Find()
	} else {
		list, err = do.Find()
	}

	if err != nil {
		logger.Error(
			"查询部门列表失败",
			logger.Err(err),
			logger.String("部门名称", params.Name),
			logger.String("部门编码", params.Code),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return nil, errs.ErrServer
	}

	return list, nil
}

func (d *deptRepository) ScanOne(ctx context.Context, id int64, val any) error {
	if id <= 0 {
		return ErrInvalidDeptID
	}

	do := d.q.Dept.WithContext(ctx)
	field := d.q.Dept
	if err := do.Where(field.ID.Eq(id)).Scan(val); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrDeptNotFound
		}
		logger.Error(
			"查询部门失败",
			logger.Err(err),
			logger.Int64("部门ID", id),
		)
		return errs.ErrServer
	}
	return nil
}

func (d *deptRepository) ScanList(ctx context.Context, params *DeptListParams, list any) error {
	do := d.q.Dept.WithContext(ctx)
	field := d.q.Dept

	// 构建查询条件
	if params.Name != "" {
		do = do.Where(field.Name.Like("%" + params.Name + "%"))
	}
	if params.Code != "" {
		do = do.Where(field.Code.Like("%" + params.Code + "%"))
	}
	if params.Status != nil {
		do = do.Where(field.Status.Eq(*params.Status))
	}
	if params.ParentID != nil {
		do = do.Where(field.Pid.Eq(*params.ParentID))
	}
	if params.ExcludeID != nil {
		do = do.Where(field.ID.Neq(*params.ExcludeID))
	}

	// 默认按照排序字段和ID升序排列
	do = do.Order(field.Sort, field.ID)

	// 分页查询
	if params.Page > 0 && params.Size > 0 {
		// 计算偏移量
		offset := (params.Page - 1) * params.Size
		if err := do.Offset(offset).Limit(params.Size).Scan(list); err != nil {
			logger.Error(
				"查询部门列表失败",
				logger.Err(err),
				logger.String("部门名称", params.Name),
				logger.String("部门编码", params.Code),
				logger.Int("页码", params.Page),
				logger.Int("每页数量", params.Size),
			)
			return errs.ErrServer
		}
	} else {
		if err := do.Scan(list); err != nil {
			logger.Error(
				"查询部门列表失败",
				logger.Err(err),
				logger.String("部门名称", params.Name),
				logger.String("部门编码", params.Code),
			)
			return errs.ErrServer
		}
	}

	return nil
}

// SubDept 查询子部门
func (d *deptRepository) SubDept(ctx context.Context, id int64) (*entity.Dept, error) {
	if id <= 0 {
		return nil, ErrInvalidDeptID
	}

	dept, err := d.FindOne(ctx, id)
	if err != nil {
		return nil, err
	}

	params := &DeptListParams{
		ParentID: &id,
	}
	children, err := d.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	dept.Children = children
	return dept, nil
}

// Tree 查询部门树
func (d *deptRepository) Tree(ctx context.Context, params *DeptListParams) ([]*entity.Dept, error) {
	// 查询所有部门
	allDepts, err := d.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建部门树
	return buildDeptTree(allDepts, 0), nil
}

// buildDeptTree 构建部门树
func buildDeptTree(depts []*entity.Dept, parentID int64) []*entity.Dept {
	tree := make([]*entity.Dept, 0)
	for _, dept := range depts {
		if dept.Pid != nil && *dept.Pid == parentID {
			children := buildDeptTree(depts, dept.ID)
			dept.Children = children
			tree = append(tree, dept)
		}
	}
	return tree
}

// NewDeptRepository 创建部门仓储
func NewDeptRepository(db *gorm.DB) DeptRepository {
	return &deptRepository{
		db: db,
		q:  query.Use(db),
	}
}
