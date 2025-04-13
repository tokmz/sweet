package repo

import (
	"context"
	"errors"
	"sweet/pkg/errs"
	"sweet/pkg/logger"
	"sweet/pkg/utils"

	"sweet/internal/apps/system/types/dto"
	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"
	"sweet/internal/apps/system/types/vo"

	"slices"

	"gorm.io/gorm"
)

var (
	ErrMenuNotFound         = errs.New(10301, "菜单不存在")
	ErrMenuHasChildren      = errs.New(10302, "存在子菜单，无法删除")
	ErrMenuPermissionExists = errs.New(10303, "权限标识已存在")
	ErrMenuNameExists       = errs.New(10304, "菜单名称已存在")
	ErrMenuPathExists       = errs.New(10305, "路由路径已存在")
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
	// ScanOne 扫描菜单
	ScanOne(ctx context.Context, dest interface{}, id int64) error // 暂时未用
	// 查询Tree结构的Menu 用于前端展示和选择
	Tree(ctx context.Context, params *dto.TreeParams) ([]*entity.Menu, int64, error)
	// 查询RouteTree结构的Menu 用于前端路由
	RouteTree(ctx context.Context, roleId int64) ([]*vo.RouteTree, int64, error)
	// 查询ItemTree结构的Menu 用于角色分配菜单
	ItemTree(ctx context.Context, roleId int64) ([]*vo.ItemTree, int64, error)
	// 查询角色拥有的菜单IDs
	FindMenusIds(ctx context.Context, roleId int64) ([]int64, error)
}

// 暂时未用
// FindMenusTreeParams 路由树参数
// type FindMenusTreeParams struct {
// 	ID   int64
// 	Pids []int64
// }

// menuRepository 菜单仓储实现
type menuRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (m *menuRepository) Create(ctx context.Context, menu *entity.Menu) error {
	return m.q.Transaction(func(tx *query.Query) error {
		do := tx.Menu.WithContext(ctx)
		field := tx.Menu
		//1、检查菜单名，权限标识是否已存在，路径
		if info, err := do.Where(field.Name.Eq(menu.Name)).
			Or(field.Permission.Eq(utils.SafeString(menu.Permission))).
			Or(field.Path.Eq(utils.SafeString(menu.Path))).
			First(); err == nil {
			if info.Name == menu.Name {
				return ErrMenuNameExists
			} else if info.Permission == menu.Permission {
				return ErrMenuPermissionExists
			} else if info.Path == menu.Path {
				return ErrMenuPathExists
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询菜单失败",
				logger.Err(err),
				logger.String("菜单名称", menu.Name),
				logger.String("权限标识", utils.SafeString(menu.Permission)),
				logger.String("路由路径", utils.SafeString(menu.Path)),
			)
			return errs.ErrServer
		}
		if err := do.Create(menu); err != nil {
			logger.Error(
				"创建菜单失败",
				logger.Err(err),
				logger.String("菜单名称", menu.Name),
				logger.Int64("菜单类型", menu.Type),
				logger.Int64("父ID", utils.SafeInt64(menu.ParentID)),
			)
			return errs.ErrServer
		}
		return nil
	})
}

func (m *menuRepository) Update(ctx context.Context, menu *entity.Menu) error {
	return m.q.Transaction(func(tx *query.Query) error {
		do := tx.Menu.WithContext(ctx)
		field := tx.Menu

		// 1. 检查菜单是否存在
		if _, err := do.Where(field.ID.Eq(menu.ID)).First(); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrMenuNotFound
			}
			logger.Error(
				"查询菜单失败",
				logger.Err(err),
				logger.Int64("菜单ID", menu.ID),
			)
			return errs.ErrServer
		}

		// 2. 检查菜单名，权限标识，路径是否与其他菜单冲突
		if info, err := do.Where(field.ID.Neq(menu.ID)).
			Where(field.Name.Eq(menu.Name)).
			Or(field.Permission.Eq(utils.SafeString(menu.Permission))).
			Or(field.Path.Eq(utils.SafeString(menu.Path))).
			First(); err == nil {
			// 存在冲突
			if info.Name == menu.Name {
				return ErrMenuNameExists
			} else if info.Permission != nil && menu.Permission != nil && *info.Permission == *menu.Permission {
				return ErrMenuPermissionExists
			} else if info.Path != nil && menu.Path != nil && *info.Path == *menu.Path {
				return ErrMenuPathExists
			}
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询菜单冲突失败",
				logger.Err(err),
				logger.Int64("菜单ID", menu.ID),
				logger.String("菜单名称", menu.Name),
				logger.String("权限标识", utils.SafeString(menu.Permission)),
				logger.String("路由路径", utils.SafeString(menu.Path)),
			)
			return errs.ErrServer
		}

		// 3. 执行更新
		if _, err := do.Where(field.ID.Eq(menu.ID)).Updates(menu); err != nil {
			logger.Error(
				"更新菜单失败",
				logger.Err(err),
				logger.Int64("菜单ID", menu.ID),
				logger.String("菜单名称", menu.Name),
				logger.Int64("菜单类型", menu.Type),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (m *menuRepository) Delete(ctx context.Context, ids []int64) error {
	return m.q.Transaction(func(tx *query.Query) error {
		do := tx.Menu.WithContext(ctx)
		field := tx.Menu

		// 检查是否有子菜单
		for _, id := range ids {
			children, err := do.Where(field.ParentID.Eq(id)).Count()
			if err != nil {
				logger.Error(
					"查询子菜单失败",
					logger.Err(err),
					logger.Int64("菜单ID", id),
				)
				return errs.ErrServer
			}
			if children > 0 {
				return ErrMenuHasChildren
			}
		}

		// 检查是否存在角色-菜单关联
		rmDo := tx.RoleMenu.WithContext(ctx)
		rmField := tx.RoleMenu

		for _, id := range ids {
			count, err := rmDo.Where(rmField.MenuID.Eq(id)).Count()
			if err != nil {
				logger.Error(
					"查询角色菜单关联失败",
					logger.Err(err),
					logger.Int64("菜单ID", id),
				)
				return errs.ErrServer
			}

			// 如果存在关联，先删除角色-菜单关联
			if count > 0 {
				if _, err = rmDo.Where(rmField.MenuID.Eq(id)).Delete(); err != nil {
					logger.Error(
						"删除角色菜单关联失败",
						logger.Err(err),
						logger.Int64("菜单ID", id),
						logger.Int64("关联数量", count),
					)
					return errs.ErrServer
				}
			}
		}

		// 批量删除菜单
		if _, err := do.Where(field.ID.In(ids...)).Delete(); err != nil {
			logger.Error(
				"删除菜单失败",
				logger.Err(err),
				logger.Any("菜单ID列表", ids),
				logger.Int("删除数量", len(ids)),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (m *menuRepository) FindOne(ctx context.Context, id int64) (*entity.Menu, error) {
	menu, err := m.q.Menu.WithContext(ctx).Where(m.q.Menu.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrMenuNotFound
		}
		logger.Error(
			"查询菜单失败",
			logger.Err(err),
			logger.Int64("菜单ID", id),
		)
		return nil, errs.ErrServer
	}
	return menu, nil
}

func (m *menuRepository) ScanOne(ctx context.Context, dest interface{}, id int64) error {
	if err := m.q.Menu.WithContext(ctx).Where(m.q.Menu.ID.Eq(id)).Scan(dest); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrMenuNotFound
		}
		logger.Error("查询菜单失败", logger.Err(err), logger.Any("id", id))
		return errs.ErrServer
	}
	return nil
}

func (m *menuRepository) Tree(ctx context.Context, params *dto.TreeParams) ([]*entity.Menu, int64, error) {
	// 构建查询条件
	do := m.q.Menu.WithContext(ctx)
	field := m.q.Menu
	var parentMenu *entity.Menu

	// 默认只查询启用的菜单
	do = do.Where(field.Status.Eq(1))

	// 按角色ID查询
	if params.Rid > 0 {
		roleMenuQ := m.q.RoleMenu
		do = do.LeftJoin(roleMenuQ, field.ID.EqCol(roleMenuQ.MenuID)).
			Where(roleMenuQ.RoleID.Eq(params.Rid))
	}

	// 如果指定了父ID或名称，查询所有相关菜单而不限制父ID
	var allMenus []*entity.Menu
	var err error

	if params.Pid > 0 || params.Name != "" {
		// 创建查询对象的副本
		menuQuery := do

		// 按名称模糊查询
		if params.Name != "" {
			menuQuery = menuQuery.Where(field.Name.Like("%" + params.Name + "%"))
		}

		// 获取所有菜单
		allMenus, err = menuQuery.Order(field.Sort).Find()
		if err != nil {
			logger.Error(
				"查询菜单列表失败",
				logger.Err(err),
				logger.Int64("角色ID", params.Rid),
				logger.Int64("父菜单ID", params.Pid),
				logger.String("菜单名称", params.Name),
			)
			return nil, 0, errs.ErrServer
		}

		// 如果指定了父ID，获取父菜单信息
		if params.Pid > 0 {
			parentMenu, err = m.FindOne(ctx, params.Pid)
			if err != nil {
				logger.Error(
					"查询父菜单失败",
					logger.Err(err),
					logger.Int64("父菜单ID", params.Pid),
					logger.Int64("角色ID", params.Rid),
				)
				return nil, 0, err
			}
		}
	} else {
		// 查询顶级菜单
		var zero int64 = 0
		allMenus, err = do.Where(field.ParentID.Eq(zero)).Order(field.Sort).Find()
		if err != nil {
			logger.Error(
				"查询顶级菜单列表失败",
				logger.Err(err),
				logger.Int64("角色ID", params.Rid),
			)
			return nil, 0, errs.ErrServer
		}
	}

	// 如果指定了名称搜索，返回平铺列表
	if params.Name != "" {
		var result []*entity.Menu
		if parentMenu != nil {
			result = append(result, parentMenu)
		}
		result = append(result, allMenus...)
		return result, int64(len(result)), nil
	}

	// 构建菜单树
	menuMap := make(map[int64]*entity.Menu)
	for _, menu := range allMenus {
		menuMap[menu.ID] = menu
		// 初始化子菜单切片
		menu.Children = make([]*entity.Menu, 0)
	}

	// 根据父ID构建菜单树
	var rootMenus []*entity.Menu

	// 如果有指定父ID，则以父菜单为根
	if params.Pid > 0 && parentMenu != nil {
		rootMenus = []*entity.Menu{parentMenu}
		parentMenu.Children = make([]*entity.Menu, 0)

		// 为父菜单添加直接子菜单
		for _, menu := range allMenus {
			if menu.ParentID != nil && *menu.ParentID == params.Pid {
				parentMenu.Children = append(parentMenu.Children, menu)
			}
		}
	} else {
		// 否则构建完整的树结构
		for _, menu := range allMenus {
			if menu.ParentID == nil || *menu.ParentID == 0 {
				// 根菜单
				rootMenus = append(rootMenus, menu)
			} else {
				// 子菜单，添加到父菜单的children中
				if parent, exists := menuMap[*menu.ParentID]; exists {
					parent.Children = append(parent.Children, menu)
				} else {
					// 父菜单不存在，记录异常情况但继续处理
					logger.Warn(
						"菜单的父ID不存在",
						logger.Int64("菜单ID", menu.ID),
						logger.String("菜单名称", menu.Name),
						logger.Int64("父菜单ID", *menu.ParentID),
					)
					// 作为根菜单处理
					rootMenus = append(rootMenus, menu)
				}
			}
		}
	}

	return rootMenus, int64(len(rootMenus)), nil
}

func (m *menuRepository) RouteTree(ctx context.Context, rid int64) ([]*vo.RouteTree, int64, error) {
	// 构建查询条件
	menuQ := m.q.Menu
	roleMenuQ := m.q.RoleMenu
	var menus []*entity.Menu

	// 如果是管理员角色，获取所有菜单
	if rid == 1 { // 假设管理员角色ID为1
		result, err := menuQ.WithContext(ctx).
			Where(menuQ.Status.Eq(1)). // 只获取启用的菜单
			Order(menuQ.Sort).
			Find()
		if err != nil {
			logger.Error(
				"管理员查询菜单失败",
				logger.Err(err),
				logger.Int64("角色ID", rid),
			)
			return nil, 0, errs.ErrServer
		}
		menus = result
	} else {
		// 获取角色关联的菜单
		result, err := menuQ.WithContext(ctx).
			LeftJoin(roleMenuQ, menuQ.ID.EqCol(roleMenuQ.MenuID)).
			Where(roleMenuQ.RoleID.Eq(rid)).
			Where(menuQ.Status.Eq(1)). // 只获取启用的菜单
			Order(menuQ.Sort).
			Find()
		if err != nil {
			logger.Error(
				"查询角色关联菜单失败",
				logger.Err(err),
				logger.Int64("角色ID", rid),
			)
			return nil, 0, errs.ErrServer
		}
		menus = result
	}

	// 构建路由树
	routeTree := buildRouteTree(menus)
	return routeTree, int64(len(routeTree)), nil
}

func (m *menuRepository) ItemTree(ctx context.Context, rid int64) ([]*vo.ItemTree, int64, error) {
	// 查询所有菜单
	menus, err := m.q.Menu.WithContext(ctx).Order(m.q.Menu.Sort).Find()
	if err != nil {
		logger.Error(
			"查询所有菜单失败",
			logger.Err(err),
			logger.Int64("角色ID", rid),
		)
		return nil, 0, errs.ErrServer
	}

	// 查询角色已有菜单
	var roleMenuIDs []int64
	if rid > 0 {
		roleMenus, err := m.q.RoleMenu.WithContext(ctx).Where(m.q.RoleMenu.RoleID.Eq(rid)).Find()
		if err != nil {
			logger.Error(
				"查询角色菜单关联失败",
				logger.Err(err),
				logger.Int64("角色ID", rid),
			)
			return nil, 0, errs.ErrServer
		}

		for _, rm := range roleMenus {
			roleMenuIDs = append(roleMenuIDs, rm.MenuID)
		}
	}

	// 构建ID到菜单的映射
	menuMap := make(map[int64]*entity.Menu)
	for _, menu := range menus {
		menuMap[menu.ID] = menu
	}

	// 构建ItemTree
	itemMap := make(map[int64]*vo.ItemTree)
	var rootItems []*vo.ItemTree

	// 第一步：创建所有节点
	for _, menu := range menus {
		// 检查菜单是否被选中
		checked := slices.Contains(roleMenuIDs, menu.ID)

		item := &vo.ItemTree{
			ID:       menu.ID,
			Name:     menu.Name,
			Checked:  checked,
			Children: make([]*vo.ItemTree, 0),
		}
		itemMap[menu.ID] = item
	}

	// 第二步：构建树形结构
	for _, menu := range menus {
		item := itemMap[menu.ID]
		if utils.SafeInt64(menu.ParentID) == 0 {
			// 根节点
			rootItems = append(rootItems, item)
		} else {
			// 子节点
			if parent, ok := itemMap[utils.SafeInt64(menu.ParentID)]; ok {
				parent.Children = append(parent.Children, item)
			} else {
				logger.Warn(
					"菜单的父ID不存在",
					logger.Int64("菜单ID", menu.ID),
					logger.String("菜单名称", menu.Name),
					logger.Int64("父菜单ID", utils.SafeInt64(menu.ParentID)),
				)
				// 作为根菜单处理
				rootItems = append(rootItems, item)
			}
		}
	}

	return rootItems, int64(len(rootItems)), nil
}

// FindMenusIds 查询角色拥有的菜单IDs
func (m *menuRepository) FindMenusIds(ctx context.Context, roleId int64) ([]int64, error) {
	if roleId <= 0 {
		return []int64{}, nil
	}

	// 查询角色-菜单关联表，获取角色拥有的所有菜单ID
	roleMenus, err := m.q.RoleMenu.WithContext(ctx).Where(m.q.RoleMenu.RoleID.Eq(roleId)).Find()
	if err != nil {
		return nil, err
	}

	// 提取菜单IDs
	menuIDs := make([]int64, 0, len(roleMenus))
	for _, rm := range roleMenus {
		menuIDs = append(menuIDs, rm.MenuID)
	}

	return menuIDs, nil
}

// NewMenuRepository 创建菜单仓储
func NewMenuRepository(db *gorm.DB) MenuRepository {
	return &menuRepository{
		db: db,
		q:  query.Use(db),
	}
}

// buildRouteTree 构建路由树
func buildRouteTree(menus []*entity.Menu) []*vo.RouteTree {
	// 创建菜单ID到路由树的映射
	menuMap := make(map[int64]*vo.RouteTree)
	// 存储根节点
	var rootNodes []*vo.RouteTree

	// 第一遍遍历：创建所有节点
	for _, menu := range menus {
		route := &vo.RouteTree{
			ID:        menu.ID,
			ParentID:  utils.SafeInt64(menu.ParentID),
			Name:      menu.Name,
			Path:      utils.SafeString(menu.Path),
			Component: utils.SafeString(menu.Component),
			Redirect:  utils.SafeString(menu.Redirect),
			Meta:      buildRouteMeta(menu),
			Children:  make([]*vo.RouteTree, 0),
		}
		menuMap[menu.ID] = route
	}

	// 第二遍遍历：构建树形结构
	for _, menu := range menus {
		route := menuMap[menu.ID]
		if utils.SafeInt64(menu.ParentID) == 0 {
			// 根节点
			rootNodes = append(rootNodes, route)
		} else {
			// 子节点
			if parent, ok := menuMap[utils.SafeInt64(menu.ParentID)]; ok {
				parent.Children = append(parent.Children, route)
			}
		}
	}

	return rootNodes
}

// buildRouteMeta 构建路由元信息
func buildRouteMeta(menu *entity.Menu) vo.RouteMeta {
	return vo.RouteMeta{
		Title:        menu.Name,
		Icon:         utils.SafeString(menu.Icon),
		Hidden:       utils.SafeBoolFromInt64(menu.Hidden),
		KeepAlive:    utils.SafeBoolFromInt64(menu.KeepAlive),
		AlwaysShow:   utils.SafeBoolFromInt64(menu.AlwaysShow),
		Target:       utils.SafeString(menu.Target),
		ActiveMenu:   utils.SafeString(menu.ActiveMenu),
		Breadcrumb:   utils.SafeBoolFromInt64(menu.Breadcrumb),
		Affix:        utils.SafeBoolFromInt64(menu.Affix),
		FrameSrc:     utils.SafeString(menu.FrameSrc),
		FrameLoading: utils.SafeBoolFromInt64(menu.FrameLoading),
		Transition:   utils.SafeString(menu.Transition),
		Permission:   utils.SafeString(menu.Permission),
	}
}
