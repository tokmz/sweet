package service

import (
	"context"
	"sweet/internal/apps/system/repo"
	"sweet/internal/apps/system/types/dto"
	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/vo"
	"sweet/internal/common"
	"sweet/pkg/logger"
	"sweet/pkg/utils"
)

type MenuService interface {
	// Create 创建菜单
	Create(ctx context.Context, req *dto.CreateMenuReq) error
	// Update 更新菜单
	Update(ctx context.Context, req *dto.UpdateMenuReq) error
	// Delete 删除菜单
	Delete(ctx context.Context, req *common.IdsReq) error
	// FindOne 查询菜单
	FindOne(ctx context.Context, req *common.IDReq) (*entity.Menu, error)
	// Tree 查询菜单树
	Tree(ctx context.Context, req *dto.FindMenuTreeReq) ([]*entity.Menu, int64, error)
	// RouteTree 查询路由树 根据角色ID
	RouteTree(ctx context.Context, req *dto.RouteTreeReq) ([]*vo.RouteTree, int64, error)
	// ItemTree 查询Item菜单树 用于分配给角色
	ItemTree(ctx context.Context, req *dto.RouteTreeReq) ([]*vo.ItemTree, int64, error)
}

// menuService 菜单服务实现
type menuService struct {
	menuRepo repo.MenuRepository
}

// Create 创建菜单
func (m *menuService) Create(ctx context.Context, req *dto.CreateMenuReq) error {
	logger.Info("创建菜单",
		logger.String("菜单名称", req.Name),
		logger.Int64("菜单类型", req.Type),
		logger.Int64("父ID", utils.SafeInt64(req.ParentID)))

	// 构建菜单实体
	menu := &entity.Menu{
		Name:         req.Name,
		Type:         req.Type,
		ParentID:     req.ParentID,
		Permission:   req.Permission,
		Path:         req.Path,
		Component:    req.Component,
		Redirect:     req.Redirect,
		Icon:         req.Icon,
		Hidden:       req.Hidden,
		Status:       req.Status,
		AlwaysShow:   req.AlwaysShow,
		KeepAlive:    req.KeepAlive,
		Target:       req.Target,
		Title:        req.Title,
		ActiveMenu:   req.ActiveMenu,
		Breadcrumb:   req.Breadcrumb,
		Affix:        req.Affix,
		FrameSrc:     req.FrameSrc,
		FrameLoading: req.FrameLoading,
		Transition:   req.Transition,
		Remark:       req.Remark,
		CreatedBy:    req.CreatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		menu.Sort = &sort
	}

	// 调用仓储层创建菜单
	return m.menuRepo.Create(ctx, menu)
}

// Update 更新菜单
func (m *menuService) Update(ctx context.Context, req *dto.UpdateMenuReq) error {
	logger.Info("更新菜单",
		logger.Int64("菜单ID", req.ID),
		logger.String("菜单名称", req.Name),
		logger.Int64("菜单类型", req.Type),
		logger.Int64("父ID", utils.SafeInt64(req.ParentID)))

	// 构建菜单实体
	menu := &entity.Menu{
		ID:           req.ID,
		Name:         req.Name,
		Type:         req.Type,
		ParentID:     req.ParentID,
		Permission:   req.Permission,
		Path:         req.Path,
		Component:    req.Component,
		Redirect:     req.Redirect,
		Icon:         req.Icon,
		Hidden:       req.Hidden,
		Status:       req.Status,
		AlwaysShow:   req.AlwaysShow,
		KeepAlive:    req.KeepAlive,
		Target:       req.Target,
		Title:        req.Title,
		ActiveMenu:   req.ActiveMenu,
		Breadcrumb:   req.Breadcrumb,
		Affix:        req.Affix,
		FrameSrc:     req.FrameSrc,
		FrameLoading: req.FrameLoading,
		Transition:   req.Transition,
		Remark:       req.Remark,
		UpdatedBy:    req.UpdatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		menu.Sort = &sort
	}

	// 调用仓储层更新菜单
	return m.menuRepo.Update(ctx, menu)
}

// Delete 删除菜单
func (m *menuService) Delete(ctx context.Context, req *common.IdsReq) error {
	logger.Info("删除菜单",
		logger.Any("菜单ID列表", req.Ids),
		logger.Int("删除数量", len(req.Ids)))

	// 调用仓储层删除菜单
	return m.menuRepo.Delete(ctx, req.Ids)
}

// FindOne 查询菜单
func (m *menuService) FindOne(ctx context.Context, req *common.IDReq) (*entity.Menu, error) {
	logger.Info("查询菜单", logger.Int64("菜单ID", req.ID))

	// 调用仓储层查询菜单
	return m.menuRepo.FindOne(ctx, req.ID)
}

// Tree 查询菜单树
func (m *menuService) Tree(ctx context.Context, req *dto.FindMenuTreeReq) ([]*entity.Menu, int64, error) {
	logger.Info("查询菜单树",
		logger.Int64("父菜单ID", utils.SafeInt64(req.ParentID)),
		logger.String("菜单名称", utils.SafeString(req.Name)),
		logger.Int64("角色ID", utils.SafeInt64(req.RoleID)))

	// 构建查询参数
	params := &dto.TreeParams{
		Pid: utils.SafeInt64(req.ParentID),
		Rid: utils.SafeInt64(req.RoleID),
	}

	if req.Name != nil {
		params.Name = *req.Name
	}

	// 调用仓储层查询菜单树
	return m.menuRepo.Tree(ctx, params)
}

// RouteTree 查询路由树
func (m *menuService) RouteTree(ctx context.Context, req *dto.RouteTreeReq) ([]*vo.RouteTree, int64, error) {
	logger.Info("查询路由树", logger.Int64("角色ID", req.RoleID))

	// 调用仓储层获取路由树
	return m.menuRepo.RouteTree(ctx, req.RoleID)
}

// ItemTree 查询Item菜单树 用于分配给角色
func (m *menuService) ItemTree(ctx context.Context, req *dto.RouteTreeReq) ([]*vo.ItemTree, int64, error) {
	logger.Info("查询Item菜单树", logger.Int64("角色ID", req.RoleID))

	// 调用仓储层获取ItemTree
	return m.menuRepo.ItemTree(ctx, req.RoleID)
}

// NewMenuService 创建菜单服务
func NewMenuService(menuRepo repo.MenuRepository) MenuService {
	return &menuService{
		menuRepo: menuRepo,
	}
}
