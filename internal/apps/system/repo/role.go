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
	ErrRoleNotFound   = errs.New(10101, "角色不存在")
	ErrRoleExists     = errs.New(10102, "角色已存在")
	ErrRoleHasUsers   = errs.New(10103, "角色下存在用户，无法删除")
	ErrRoleCodeExists = errs.New(10105, "角色编码已存在")
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
	Status int64  // 状态
	Page   int    // 页码
	Size   int    // 每页数量
}

// roleRepository 角色仓储实现
type roleRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (r *roleRepository) Create(ctx context.Context, role *entity.Role) error {
	return r.q.Transaction(func(tx *query.Query) error {
		do := tx.Role.WithContext(ctx)
		field := tx.Role

		// 检查角色名称或编码是否已存在
		if info, err := do.Where(field.Name.Eq(role.Name)).
			Or(field.Code.Eq(role.Code)).
			First(); err == nil {
			if info.Name == role.Name {
				return ErrRoleExists
			} else if info.Code == role.Code {
				return ErrRoleCodeExists
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询角色失败",
				logger.Err(err),
				logger.String("角色名称", role.Name),
				logger.String("角色编码", role.Code),
			)
			return errs.ErrServer
		}

		// 创建角色
		if err := do.Create(role); err != nil {
			logger.Error(
				"创建角色失败",
				logger.Err(err),
				logger.String("角色名称", role.Name),
				logger.String("角色编码", role.Code),
				logger.Int64("角色状态", utils.SafeInt64(role.Status)),
			)
			return errs.ErrServer
		}
		return nil
	})
}

func (r *roleRepository) Update(ctx context.Context, role *entity.Role) error {
	return r.q.Transaction(func(tx *query.Query) error {
		do := tx.Role.WithContext(ctx)
		field := tx.Role

		// 检查要更新的角色是否存在
		if _, err := do.Where(field.ID.Eq(role.ID)).First(); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRoleNotFound
			}
			logger.Error(
				"查询角色失败",
				logger.Err(err),
				logger.Int64("角色ID", role.ID),
			)
			return errs.ErrServer
		}

		// 检查角色名称是否重复
		if role.Name != "" {
			if _, err := do.Where(field.ID.Neq(role.ID), field.Name.Eq(role.Name)).First(); err == nil {
				return ErrRoleExists
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error(
					"检查角色名称唯一性失败",
					logger.Err(err),
					logger.Int64("角色ID", role.ID),
					logger.String("角色名称", role.Name),
				)
				return errs.ErrServer
			}
		}

		// 检查角色编码是否重复
		if role.Code != "" {
			if _, err := do.Where(field.ID.Neq(role.ID), field.Code.Eq(role.Code)).First(); err == nil {
				return ErrRoleCodeExists
			} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
				logger.Error(
					"检查角色编码唯一性失败",
					logger.Err(err),
					logger.Int64("角色ID", role.ID),
					logger.String("角色编码", role.Code),
				)
				return errs.ErrServer
			}
		}

		// 更新角色信息
		if _, err := do.Where(field.ID.Eq(role.ID)).Updates(role); err != nil {
			logger.Error(
				"更新角色失败",
				logger.Err(err),
				logger.Int64("角色ID", role.ID),
				logger.String("角色名称", role.Name),
				logger.String("角色编码", role.Code),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (r *roleRepository) Delete(ctx context.Context, ids []int64) error {
	return r.q.Transaction(func(tx *query.Query) error {
		// 检查角色是否有关联用户
		for _, id := range ids {
			userIds, err := r.CheckHasUsers(ctx, id)
			if err != nil {
				return err
			}
			if len(userIds) > 0 {
				return ErrRoleHasUsers
			}
		}

		// 删除角色和角色菜单关联
		if result, err := tx.WithContext(ctx).RoleMenu.Where(tx.RoleMenu.RoleID.In(ids...)).Delete(); err != nil {
			logger.Error(
				"删除角色菜单关联失败",
				logger.Err(err),
				logger.Any("角色ID列表", ids),
			)
			return errs.ErrServer
		} else {
			logger.Debug(
				"删除角色菜单关联",
				logger.Int64("关联数量", result.RowsAffected),
				logger.Any("角色ID列表", ids),
			)
		}

		// 删除角色
		if _, err := tx.Role.WithContext(ctx).Where(tx.Role.ID.In(ids...)).Delete(); err != nil {
			logger.Error(
				"删除角色失败",
				logger.Err(err),
				logger.Any("角色ID列表", ids),
				logger.Int("删除数量", len(ids)),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (r *roleRepository) FindOne(ctx context.Context, id int64) (*entity.Role, error) {
	if info, err := r.q.Role.WithContext(ctx).Where(r.q.Role.ID.Eq(id)).First(); err == nil {
		return info, nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, ErrRoleNotFound
	} else {
		logger.Error(
			"查询角色失败",
			logger.Err(err),
			logger.Int64("角色ID", id),
		)
		return nil, errs.ErrServer
	}
}

func (r *roleRepository) FindList(ctx context.Context, params *RoleListParams) (list []*entity.Role, total int64, err error) {
	do := r.q.Role.WithContext(ctx)
	field := r.q.Role

	// 构建查询条件
	if params.Name != "" {
		do = do.Where(field.Name.Like("%" + params.Name + "%"))
	}
	if params.Code != "" {
		do = do.Where(field.Code.Like("%" + params.Code + "%"))
	}
	if params.Status != 0 {
		do = do.Where(field.Status.Eq(params.Status))
	}

	// 默认按照排序字段升序排列
	do = do.Order(field.Sort)

	// 执行分页查询
	if list, total, err = do.FindByPage(params.Page, params.Size); err == nil {
		return list, total, nil
	} else {
		logger.Error(
			"查询角色列表失败",
			logger.Err(err),
			logger.String("角色名称", params.Name),
			logger.String("角色编码", params.Code),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return nil, 0, errs.ErrServer
	}
}

func (r *roleRepository) ScanOne(ctx context.Context, id int64, val any) error {
	if err := r.q.Role.WithContext(ctx).Where(r.q.Role.ID.Eq(id)).Scan(val); err == nil {
		return nil
	} else if errors.Is(err, gorm.ErrRecordNotFound) {
		return ErrRoleNotFound
	} else {
		logger.Error(
			"查询角色失败",
			logger.Err(err),
			logger.Int64("角色ID", id),
		)
		return errs.ErrServer
	}
}

func (r *roleRepository) ScanList(ctx context.Context, params *RoleListParams, list any, total *int64) error {
	do := r.q.Role.WithContext(ctx)
	field := r.q.Role

	// 构建查询条件
	if params.Name != "" {
		do = do.Where(field.Name.Like("%" + params.Name + "%"))
	}
	if params.Code != "" {
		do = do.Where(field.Code.Like("%" + params.Code + "%"))
	}
	if params.Status != 0 {
		do = do.Where(field.Status.Eq(params.Status))
	}

	// 默认按照排序字段升序排列
	do = do.Order(field.Sort)

	// 执行分页查询并扫描到指定结构体
	if count, err := do.ScanByPage(list, params.Page, params.Size); err == nil {
		*total = count
		return nil
	} else {
		logger.Error(
			"查询角色列表失败",
			logger.Err(err),
			logger.String("角色名称", params.Name),
			logger.String("角色编码", params.Code),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return errs.ErrServer
	}
}

func (r *roleRepository) CheckHasUsers(ctx context.Context, roleID int64) ([]int64, error) {
	var userIds []int64
	if err := r.q.User.WithContext(ctx).Where(r.q.User.RoleID.Eq(roleID)).Select(r.q.User.ID).Scan(&userIds); err != nil {
		logger.Error(
			"查询角色下用户失败",
			logger.Err(err),
			logger.Int64("角色ID", roleID),
		)
		return nil, errs.ErrServer
	}
	return userIds, nil
}

func (r *roleRepository) AssignMenus(ctx context.Context, roleID int64, menuIds []int64) error {
	return r.q.Transaction(func(tx *query.Query) error {
		// 检查角色是否存在
		if _, err := tx.Role.WithContext(ctx).Where(tx.Role.ID.Eq(roleID)).First(); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrRoleNotFound
			}
			logger.Error(
				"查询角色失败",
				logger.Err(err),
				logger.Int64("角色ID", roleID),
			)
			return errs.ErrServer
		}

		// 删除原有的角色菜单关联
		if result, err := tx.WithContext(ctx).RoleMenu.Where(tx.RoleMenu.RoleID.Eq(roleID)).Delete(); err != nil {
			logger.Error(
				"删除角色菜单关联失败",
				logger.Err(err),
				logger.Int64("角色ID", roleID),
			)
			return errs.ErrServer
		} else {
			logger.Debug(
				"删除角色菜单关联",
				logger.Int64("关联数量", result.RowsAffected),
				logger.Int64("角色ID", roleID),
			)
		}

		// 添加新的角色菜单关联
		if len(menuIds) > 0 {
			roleMenus := make([]*entity.RoleMenu, 0, len(menuIds))
			for _, menuID := range menuIds {
				roleMenus = append(roleMenus, &entity.RoleMenu{
					RoleID: roleID,
					MenuID: menuID,
				})
			}
			if err := tx.WithContext(ctx).RoleMenu.Create(roleMenus...); err != nil {
				logger.Error(
					"创建角色菜单关联失败",
					logger.Err(err),
					logger.Int64("角色ID", roleID),
					logger.Int("菜单数量", len(menuIds)),
				)
				return errs.ErrServer
			}
		}

		return nil
	})
}

func (r *roleRepository) FindMenusTree(ctx context.Context, roleID int64) ([]*entity.Menu, error) {
	// 查询角色关联的菜单ID
	menuIds, err := r.FindMenusIds(ctx, roleID)
	if err != nil {
		return nil, err
	}

	// 查询菜单信息
	if len(menuIds) == 0 {
		return []*entity.Menu{}, nil
	}

	menus, err := r.q.Menu.WithContext(ctx).Where(r.q.Menu.ID.In(menuIds...)).Find()
	if err != nil {
		logger.Error(
			"查询角色菜单失败",
			logger.Err(err),
			logger.Int64("角色ID", roleID),
			logger.Int("菜单ID数量", len(menuIds)),
		)
		return nil, errs.ErrServer
	}

	// 构建菜单树
	return BuildMenuTree(menus, 0), nil
}

func (r *roleRepository) FindMenusIds(ctx context.Context, roleID int64) ([]int64, error) {
	var menuIds []int64
	if err := r.q.RoleMenu.WithContext(ctx).Where(r.q.RoleMenu.RoleID.Eq(roleID)).Select(r.q.RoleMenu.MenuID).Scan(&menuIds); err != nil {
		logger.Error(
			"查询角色菜单ID失败",
			logger.Err(err),
			logger.Int64("角色ID", roleID),
		)
		return nil, errs.ErrServer
	}
	return menuIds, nil
}

// BuildMenuTree 构建菜单树
func BuildMenuTree(menus []*entity.Menu, parentID int64) []*entity.Menu {
	tree := make([]*entity.Menu, 0)
	for _, menu := range menus {
		if menu.ParentID != nil && *menu.ParentID == parentID {
			children := BuildMenuTree(menus, menu.ID)
			menu.Children = children
			tree = append(tree, menu)
		}
	}
	return tree
}

// NewRoleRepository 创建角色仓储
func NewRoleRepository(db *gorm.DB) RoleRepository {
	return &roleRepository{
		db: db,
		q:  query.Use(db),
	}
}
