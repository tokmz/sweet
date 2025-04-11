package repository

import (
	"context"
	"sweet/pkg/errs"
	"sweet/pkg/utils"

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

// RouteTree 路由树结构
type RouteTree struct {
	ID        int64        `json:"id"`        // 菜单ID
	ParentID  int64        `json:"parentId"`  // 父菜单ID
	Name      string       `json:"name"`      // 路由名称
	Path      string       `json:"path"`      // 路由地址
	Component string       `json:"component"` // 组件路径
	Redirect  string       `json:"redirect"`  // 重定向地址
	Meta      RouteMeta    `json:"meta"`      // 路由元信息
	Children  []*RouteTree `json:"children"`  // 子路由
}

// RouteMeta 路由元信息
type RouteMeta struct {
	Title        string `json:"title"`        // 菜单标题
	Icon         string `json:"icon"`         // 菜单图标
	Hidden       bool   `json:"hidden"`       // 是否隐藏
	KeepAlive    bool   `json:"keepAlive"`    // 是否缓存
	AlwaysShow   bool   `json:"alwaysShow"`   // 是否总是显示
	Target       string `json:"target"`       // 打开方式
	ActiveMenu   string `json:"activeMenu"`   // 激活菜单
	Breadcrumb   bool   `json:"breadcrumb"`   // 是否显示面包屑
	Affix        bool   `json:"affix"`        // 是否固定
	FrameSrc     string `json:"frameSrc"`     // iframe地址
	FrameLoading bool   `json:"frameLoading"` // iframe加载状态
	Transition   string `json:"transition"`   // 过渡动画
	Permission   string `json:"permission"`   // 权限标识
}

type ItemTree struct {
	ID       string      `json:"id"`
	Name     string      `json:"name"`
	Children []*ItemTree `json:"children"`
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
	// 判断角色ID是否有效
	if rid <= 0 {
		return nil, 0, ErrInvalidMenuID
	}

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
			return nil, 0, err
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
			return nil, 0, err
		}
		menus = result
	}

	// 构建路由树
	routeTree := buildRouteTree(menus)
	return routeTree, int64(len(routeTree)), nil
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

// buildRouteTree 构建路由树
func buildRouteTree(menus []*entity.Menu) []*RouteTree {
	// 创建菜单ID到路由树的映射
	menuMap := make(map[int64]*RouteTree)
	// 存储根节点
	var rootNodes []*RouteTree

	// 第一遍遍历：创建所有节点
	for _, menu := range menus {
		route := &RouteTree{
			ID:        menu.ID,
			ParentID:  utils.SafeInt64(menu.ParentID),
			Name:      menu.Name,
			Path:      utils.SafeString(menu.Path),
			Component: utils.SafeString(menu.Component),
			Redirect:  utils.SafeString(menu.Redirect),
			Meta:      buildRouteMeta(menu),
			Children:  make([]*RouteTree, 0),
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
func buildRouteMeta(menu *entity.Menu) RouteMeta {
	return RouteMeta{
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
