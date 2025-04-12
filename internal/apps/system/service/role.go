package service

import (
	"context"
	"sweet/internal/apps/system/repo"
	"sweet/internal/apps/system/types/dto"
	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/vo"
	"sweet/internal/common"
	"sweet/pkg/logger"
)

type RoleService interface {
	// Create 创建角色
	Create(ctx context.Context, req *dto.CreateRoleReq) error
	// Update 更新角色
	Update(ctx context.Context, req *dto.UpdateRoleReq) error
	// Delete 删除角色
	Delete(ctx context.Context, req *common.IdsReq) error
	// FindOne 查询角色
	FindOne(ctx context.Context, req *common.IDReq) (*entity.Role, error)
	// FindList 查询角色列表
	FindList(ctx context.Context, req *dto.FindListRoleReq) (*common.Page, error)
	// ListItem 查询角色列表 返回全部可用角色，用于创建用户和更新用户角色信息使用
	ListItem(ctx context.Context) (*common.Page, error)
	// UpdateStatus 批量更新角色状态
	UpdateStatus(ctx context.Context, req *common.StatusReq) error

	// AssignMenus 给角色分配菜单
	AssignMenus(ctx context.Context, req *dto.AssignMenusReq) error
	// FindMenusIds 查询角色关联菜单ids
	FindMenusIds(ctx context.Context, req *common.IDReq) ([]int64, error)
}

// roleService 角色服务实现
type roleService struct {
	roleRepo repo.RoleRepository
}

// Create 创建角色
func (r *roleService) Create(ctx context.Context, req *dto.CreateRoleReq) error {
	logger.Info("创建角色",
		logger.String("角色名称", req.Name),
		logger.String("角色编码", req.Code))

	// 构建角色实体
	role := &entity.Role{
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		DataScope: req.DataScope,
		Remark:    req.Remark,
		CreatedBy: req.CreatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		role.Sort = &sort
	}

	// 调用仓储层创建角色
	return r.roleRepo.Create(ctx, role)
}

// Update 更新角色
func (r *roleService) Update(ctx context.Context, req *dto.UpdateRoleReq) error {
	logger.Info("更新角色",
		logger.Int64("角色ID", req.ID),
		logger.String("角色名称", req.Name),
		logger.String("角色编码", req.Code))

	// 构建角色实体
	role := &entity.Role{
		ID:        req.ID,
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		DataScope: req.DataScope,
		Remark:    req.Remark,
		UpdatedBy: req.UpdatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		role.Sort = &sort
	}

	// 调用仓储层更新角色
	return r.roleRepo.Update(ctx, role)
}

// Delete 删除角色
func (r *roleService) Delete(ctx context.Context, req *common.IdsReq) error {
	logger.Info("删除角色",
		logger.Any("角色ID列表", req.Ids),
		logger.Int("删除数量", len(req.Ids)))

	// 调用仓储层删除角色
	return r.roleRepo.Delete(ctx, req.Ids)
}

// FindOne 查询角色
func (r *roleService) FindOne(ctx context.Context, req *common.IDReq) (*entity.Role, error) {
	logger.Info("查询角色", logger.Int64("角色ID", req.ID))

	// 调用仓储层查询角色
	return r.roleRepo.FindOne(ctx, req.ID)
}

// FindList 查询角色列表
func (r *roleService) FindList(ctx context.Context, req *dto.FindListRoleReq) (*common.Page, error) {
	logger.Info("查询角色列表",
		logger.String("角色名称", req.Name),
		logger.String("角色编码", req.Code),
		logger.Int("页码", req.Page),
		logger.Int("每页数量", req.Size))

	// 构建查询参数
	params := &repo.RoleListParams{
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
		Page:   req.Page,
		Size:   req.Size,
	}

	// 调用仓储层查询角色列表
	list, total, err := r.roleRepo.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建分页响应
	return common.NewPage(list, total), nil
}

// ListItem 查询角色列表 返回全部可用角色
func (r *roleService) ListItem(ctx context.Context) (*common.Page, error) {
	logger.Info("查询可用角色列表")

	// 构建查询参数，只查询状态正常的角色
	params := &repo.RoleListParams{
		Status: 1, // 正常状态的角色
		Page:   1,
		Size:   1000, // 设置较大值返回全部数据
	}

	// 调用仓储层查询角色列表
	var list []*vo.RoleItemRes
	total := int64(0)
	err := r.roleRepo.ScanList(ctx, params, &list, &total)
	if err != nil {
		return nil, err
	}

	// 构建分页响应
	return common.NewPage(list, total), nil
}

// UpdateStatus 批量更新角色状态
func (r *roleService) UpdateStatus(ctx context.Context, req *common.StatusReq) error {
	logger.Info("批量更新角色状态",
		logger.Any("角色ID列表", req.Ids),
		logger.Int64("状态", req.Status),
		logger.Int("数量", len(req.Ids)))

	// 遍历处理每个角色
	for _, id := range req.Ids {
		// 先查询角色
		role, err := r.roleRepo.FindOne(ctx, id)
		if err != nil {
			logger.Error("查询角色失败",
				logger.Err(err),
				logger.Int64("角色ID", id))
			continue
		}

		// 更新状态
		status := req.Status
		role.Status = &status

		// 调用仓储层更新角色
		if err = r.roleRepo.Update(ctx, role); err != nil {
			logger.Error("更新角色状态失败",
				logger.Err(err),
				logger.Int64("角色ID", id),
				logger.Int64("状态", req.Status))
			return err
		}
	}

	return nil
}

// AssignMenus 给角色分配菜单
func (r *roleService) AssignMenus(ctx context.Context, req *dto.AssignMenusReq) error {
	logger.Info("给角色分配菜单",
		logger.Int64("角色ID", req.ID),
		logger.Any("菜单ID列表", req.MenuIds),
		logger.Int("菜单数量", len(req.MenuIds)))

	// 调用仓储层给角色分配菜单
	return r.roleRepo.AssignMenus(ctx, req.ID, req.MenuIds)
}

// FindMenusIds 查询角色关联菜单ids
func (r *roleService) FindMenusIds(ctx context.Context, req *common.IDReq) ([]int64, error) {
	logger.Info("查询角色关联菜单ids", logger.Int64("角色ID", req.ID))

	// 调用仓储层查询角色关联菜单ids
	return r.roleRepo.FindMenusIds(ctx, req.ID)
}

// NewRoleService 创建角色服务
func NewRoleService(roleRepo repo.RoleRepository) RoleService {
	return &roleService{
		roleRepo: roleRepo,
	}
}
