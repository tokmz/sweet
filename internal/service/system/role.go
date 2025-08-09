package system

import (
	"context"
	"errors"
	"sweet/internal/global"
	"sweet/internal/models"
	systemDTO "sweet/internal/models/dto/system"
	"sweet/internal/models/entity"
	"sweet/internal/models/query"
	"sweet/pkg/errs"
	"time"

	"go.uber.org/zap"
	"gorm.io/gorm"
)

type RoleService struct{}

func (s *RoleService) CreateRole(ctx context.Context, req *systemDTO.CreateRoleReq) error {
	return global.Query.Transaction(func(tx *query.Query) error {
		dao := tx.SysRole
		if exist, err := dao.WithContext(ctx).Where(dao.Code.Eq(req.Code)).Or(dao.Name.Eq(req.Name)).First(); err == nil {
			// 判断是哪个字段存在
			if exist.Code == req.Code {
				return errs.ErrRoleCodeExists
			}
			return errs.ErrRoleNameExists
		} else if !errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error(
				"查询角色失败",
				zap.String("code", req.Code),
				zap.String("name", req.Name),
				zap.Error(err),
			)
			return errs.ErrServer
		}
		roleEntity := entity.SysRole{
			Name:    req.Name,
			Code:    req.Code,
			Sort:    req.Sort,
			IsSuper: req.IsSuper,
			Status:  req.Status,
			Remark:  req.Remark,
		}
		if err := dao.WithContext(ctx).Create(&roleEntity); err != nil {
			global.Logger.Error(
				"创建角色失败",
				zap.Any("req", req),
				zap.Error(err),
			)
			return errs.ErrServer
		}
		return nil
	})
}

func (s *RoleService) DeleteRole(ctx context.Context, req *systemDTO.DeleteRoleReq) error {
	dao := global.Query.SysRole
	if _, err := dao.WithContext(ctx).Where(dao.ID.In(req.Ids...)).Delete(); err != nil {
		global.Logger.Error(
			"删除角色失败",
			zap.Int64s("ids", req.Ids),
			zap.Error(err),
		)
		return errs.ErrServer
	}
	return nil
}

func (s *RoleService) UpdateRole(ctx context.Context, req *systemDTO.UpdateRoleReq) error {
	return global.Query.Transaction(func(tx *query.Query) error {
		dao := tx.SysRole

		// 检查角色是否存在
		existingRole, err := dao.WithContext(ctx).Where(dao.ID.Eq(req.ID)).First()
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				global.Logger.Error(
					"角色不存在",
					zap.Int64("id", req.ID),
					zap.Error(err),
				)
				return errs.ErrRoleNotFound
			}
			global.Logger.Error(
				"查询角色失败",
				zap.Int64("id", req.ID),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		// 检查是否为系统内置角色
		if existingRole.IsSystem != nil && *existingRole.IsSystem == 1 {
			global.Logger.Error(
				"系统内置角色不允许修改",
				zap.Int64("id", req.ID),
				zap.String("name", existingRole.Name),
			)
			return errs.ErrSystemRoleCannotModify
		}

		// 检查角色名称是否重复（排除自己）
		if req.Name != "" && req.Name != existingRole.Name {
			count, err := dao.WithContext(ctx).Where(
				dao.Name.Eq(req.Name),
				dao.ID.Neq(req.ID),
			).Count()
			if err != nil {
				global.Logger.Error(
					"检查角色名称重复失败",
					zap.String("name", req.Name),
					zap.Error(err),
				)
				return errs.ErrServer
			}
			if count > 0 {
				global.Logger.Error(
					"角色名称已存在",
					zap.String("name", req.Name),
				)
				return errs.ErrRoleNameExists
			}
		}

		// 检查角色标识是否重复（排除自己）
		if req.Code != "" && req.Code != existingRole.Code {
			count, err := dao.WithContext(ctx).Where(
				dao.Code.Eq(req.Code),
				dao.ID.Neq(req.ID),
			).Count()
			if err != nil {
				global.Logger.Error(
					"检查角色标识重复失败",
					zap.String("code", req.Code),
					zap.Error(err),
				)
				return errs.ErrServer
			}
			if count > 0 {
				global.Logger.Error(
					"角色标识已存在",
					zap.String("code", req.Code),
				)
				return errs.ErrRoleCodeExists
			}
		}

		// 构建更新数据
		updateData := make(map[string]interface{})
		if req.Name != "" {
			updateData["name"] = req.Name
		}
		if req.Code != "" {
			updateData["code"] = req.Code
		}
		if req.Sort != 0 {
			updateData["sort"] = req.Sort
		}
		if req.IsSuper != nil {
			updateData["is_super"] = *req.IsSuper
		}
		if req.Status != nil {
			updateData["status"] = *req.Status
		}
		if req.Remark != nil {
			updateData["remark"] = *req.Remark
		}
		updateData["updated_at"] = time.Now()

		// 执行更新
		_, err = dao.WithContext(ctx).Where(dao.ID.Eq(req.ID)).Updates(updateData)
		if err != nil {
			global.Logger.Error(
				"更新角色失败",
				zap.Int64("id", req.ID),
				zap.Any("updateData", updateData),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (s *RoleService) ListRole(ctx context.Context, req *systemDTO.RoleListReq) (*systemDTO.RoleListRes, error) {
	dao := global.Query.SysRole
	query := dao.WithContext(ctx)

	// 条件查询
	if req.Name != "" {
		query = query.Where(dao.Name.Like("%" + req.Name + "%"))
	}
	if req.IsSystem != nil {
		query = query.Where(dao.IsSystem.Eq(*req.IsSystem))
	}
	if req.IsSuper != nil {
		query = query.Where(dao.IsSuper.Eq(*req.IsSuper))
	}
	if req.Status != nil {
		query = query.Where(dao.Status.Eq(*req.Status))
	}
	// 时间范围查询
	if req.StartTime > 0 {
		query = query.Where(dao.CreatedAt.Gte(time.Unix(req.StartTime, 0)))
	}
	if req.EndTime > 0 {
		query = query.Where(dao.CreatedAt.Lte(time.Unix(req.EndTime, 0)))
	}

	// 排序
	if req.Field != "" {
		if req.Order == "asc" {
			switch req.Field {
			case "sort":
				query = query.Order(dao.Sort)
			case "created_at":
				query = query.Order(dao.CreatedAt)
			case "name":
				query = query.Order(dao.Name)
			default:
				query = query.Order(dao.Sort, dao.CreatedAt.Desc())
			}
		} else {
			switch req.Field {
			case "sort":
				query = query.Order(dao.Sort.Desc())
			case "created_at":
				query = query.Order(dao.CreatedAt.Desc())
			case "name":
				query = query.Order(dao.Name.Desc())
			default:
				query = query.Order(dao.Sort, dao.CreatedAt.Desc())
			}
		}
	} else {
		// 默认排序：Sort升序，CreatedAt降序
		query = query.Order(dao.Sort, dao.CreatedAt.Desc())
	}

	// 分页参数验证
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.Size <= 0 {
		req.Size = 10
	}
	if req.Size > 100 {
		req.Size = 100
	}

	// 使用FindByPage方法一次性完成分页查询和计数，提高性能
	offset := (req.Page - 1) * req.Size
	roles, total, err := query.FindByPage(offset, req.Size)
	if err != nil {
		global.Logger.Error(
			"查询角色列表失败",
			zap.Any("req", req),
			zap.Error(err),
		)
		return nil, errs.ErrServer
	}

	// 转换为DTO
	list := make([]*systemDTO.RoleListItem, 0, len(roles))
	for _, role := range roles {
		list = append(list, &systemDTO.RoleListItem{
			ID:        role.ID,
			Name:      role.Name,
			Code:      role.Code,
			Sort:      role.Sort,
			IsSystem:  role.IsSystem,
			IsSuper:   role.IsSuper,
			Status:    role.Status,
			CreatedAt: role.CreatedAt,
		})
	}

	return &systemDTO.RoleListRes{
		List:  list,
		Total: total,
	}, nil
}

func (s *RoleService) GetRoleDetail(ctx context.Context, req *models.IDReq) (*systemDTO.RoleDetailRes, error) {
	dao := global.Query.SysRole

	// 使用Scan方法查询角色详情
	var detail systemDTO.RoleDetailRes
	err := dao.WithContext(ctx).Where(dao.ID.Eq(req.ID)).Scan(&detail)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			global.Logger.Error(
				"角色不存在",
				zap.Int64("id", req.ID),
				zap.Error(err),
			)
			return nil, errs.ErrRoleNotFound
		}
		global.Logger.Error(
			"查询角色详情失败",
			zap.Int64("id", req.ID),
			zap.Error(err),
		)
		return nil, errs.ErrServer
	}

	return &detail, nil
}

func (s *RoleService) RoleOptions(ctx context.Context) (*systemDTO.RoleOptionRes, error) {
	dao := global.Query.SysRole
	query := dao.WithContext(ctx)

	// 只查询正常状态的角色
	query = query.Where(dao.Status.Eq(1))

	// 按排序和创建时间排序
	query = query.Order(dao.Sort, dao.CreatedAt.Desc())

	// 使用Find方法查询所有角色选项
	roles, err := query.Find()
	if err != nil {
		global.Logger.Error(
			"查询角色选项失败",
			zap.Error(err),
		)
		return nil, errs.ErrServer
	}

	// 转换为DTO
	list := make([]*systemDTO.RoleOptionItem, 0, len(roles))
	for _, role := range roles {
		list = append(list, &systemDTO.RoleOptionItem{
			ID:   role.ID,
			Name: role.Name,
		})
	}

	return &systemDTO.RoleOptionRes{
		List:  list,
		Total: int64(len(roles)),
	}, nil
}

func (s *RoleService) RoleMenuIds(ctx context.Context, req *models.IDReq) (*systemDTO.RoleMenuIdsRes, error) {
	dao := global.Query.SysRoleMenu
	query := dao.WithContext(ctx)

	// 根据角色ID查询角色菜单关联记录
	roleMenus, err := query.Where(dao.RoleID.Eq(req.ID)).Find()
	if err != nil {
		global.Logger.Error(
			"查询角色菜单关联失败",
			zap.Int64("role_id", req.ID),
			zap.Error(err),
		)
		return nil, errs.ErrServer
	}

	// 提取菜单ID列表
	menuIds := make([]int64, 0, len(roleMenus))
	for _, roleMenu := range roleMenus {
		menuIds = append(menuIds, roleMenu.MenuID)
	}

	return &systemDTO.RoleMenuIdsRes{
		Ids: menuIds,
	}, nil
}

func (s *RoleService) AssignRoleMenuIds(ctx context.Context, req *systemDTO.AssignRoleMenuIdsReq) error {
	// 如果没有新的菜单ID，需要确认是否为有意的清空操作
	if len(req.MenuIds) == 0 && !req.ConfirmClear {
		global.Logger.Warn(
			"尝试清空角色菜单权限但未确认",
			zap.Int64("role_id", req.ID),
		)
		return errs.ErrRoleMenuIdsEmpty
	}

	return global.Query.Transaction(func(tx *query.Query) error {
		dao := tx.SysRoleMenu

		// 先删除该角色的所有现有菜单关联
		if _, err := dao.WithContext(ctx).Where(dao.RoleID.Eq(req.ID)).Delete(); err != nil {
			global.Logger.Error(
				"删除角色菜单关联失败",
				zap.Int64("role_id", req.ID),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		// 如果没有新的菜单ID，直接返回（清空权限）
		if len(req.MenuIds) == 0 {
			global.Logger.Info(
				"角色菜单权限已清空",
				zap.Int64("role_id", req.ID),
			)
			return nil
		}

		// 批量插入新的菜单关联
		roleMenus := make([]*entity.SysRoleMenu, 0, len(req.MenuIds))
		for _, menuId := range req.MenuIds {
			roleMenus = append(roleMenus, &entity.SysRoleMenu{
				RoleID: req.ID,
				MenuID: menuId,
			})
		}

		if err := dao.WithContext(ctx).CreateInBatches(roleMenus, 100); err != nil {
			global.Logger.Error(
				"批量创建角色菜单关联失败",
				zap.Int64("role_id", req.ID),
				zap.Any("menu_ids", req.MenuIds),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		global.Logger.Info(
			"角色菜单分配成功",
			zap.Int64("role_id", req.ID),
			zap.Int("menu_count", len(req.MenuIds)),
		)
		return nil
	})
}

func (s *RoleService) RoleApiIds(ctx context.Context, req *models.IDReq) (*systemDTO.RoleApiIdsRes, error) {
	dao := global.Query.SysRoleApi
	query := dao.WithContext(ctx)

	// 根据角色ID查询角色API关联记录
	roleApis, err := query.Where(dao.RoleID.Eq(req.ID)).Find()
	if err != nil {
		global.Logger.Error(
			"查询角色API关联失败",
			zap.Int64("role_id", req.ID),
			zap.Error(err),
		)
		return nil, errs.ErrServer
	}

	// 提取API ID列表
	apiIds := make([]int64, 0, len(roleApis))
	for _, roleApi := range roleApis {
		apiIds = append(apiIds, roleApi.APIID)
	}

	return &systemDTO.RoleApiIdsRes{
		Ids: apiIds,
	}, nil
}

func (s *RoleService) AssignRoleApiIds(ctx context.Context, req *systemDTO.AssignRoleApiIdsReq) error {
	// 如果没有新的API ID，需要确认是否为有意的清空操作
	if len(req.ApiIds) == 0 && (req.ConfirmClear == nil || !*req.ConfirmClear) {
		global.Logger.Warn(
			"尝试清空角色API权限但未确认",
			zap.Int64("role_id", req.ID),
		)
		return errs.ErrParams
	}

	return global.Query.Transaction(func(tx *query.Query) error {
		dao := tx.SysRoleApi

		// 先删除该角色的所有现有API关联
		if _, err := dao.WithContext(ctx).Where(dao.RoleID.Eq(req.ID)).Delete(); err != nil {
			global.Logger.Error(
				"删除角色API关联失败",
				zap.Int64("role_id", req.ID),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		// 如果没有新的API ID，直接返回（清空权限）
		if len(req.ApiIds) == 0 {
			global.Logger.Info(
				"角色API权限已清空",
				zap.Int64("role_id", req.ID),
			)
			return nil
		}

		// 批量插入新的API关联
		roleApis := make([]*entity.SysRoleApi, 0, len(req.ApiIds))
		for _, apiId := range req.ApiIds {
			roleApis = append(roleApis, &entity.SysRoleApi{
				RoleID: req.ID,
				APIID:  apiId,
			})
		}

		if err := dao.WithContext(ctx).CreateInBatches(roleApis, 100); err != nil {
			global.Logger.Error(
				"批量创建角色API关联失败",
				zap.Int64("role_id", req.ID),
				zap.Any("api_ids", req.ApiIds),
				zap.Error(err),
			)
			return errs.ErrServer
		}

		global.Logger.Info(
			"角色API分配成功",
			zap.Int64("role_id", req.ID),
			zap.Int("api_count", len(req.ApiIds)),
		)
		return nil
	})
}

func NewRoleService() IRoleService {
	return &RoleService{}
}
