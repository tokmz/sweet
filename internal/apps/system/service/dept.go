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

// DeptService 部门服务接口
type DeptService interface {
	// Create 创建部门
	Create(ctx context.Context, req *dto.CreateDeptReq) error
	// Update 更新部门
	Update(ctx context.Context, req *dto.UpdateDeptReq) error
	// Delete 删除部门
	Delete(ctx context.Context, req *common.IdsReq) error
	// FindOne 查询部门
	FindOne(ctx context.Context, req *common.IDReq) (*vo.DeptDetailRes, error)
	// FindList 查询部门列表
	FindList(ctx context.Context, req *dto.FindListDeptReq) (*common.Page, error)
	// Tree 查询部门树
	Tree(ctx context.Context, req *dto.DeptTreeReq) ([]*vo.DeptTreeRes, error)
	// SubDept 查询子部门
	SubDept(ctx context.Context, req *dto.SubDeptReq) (*vo.DeptTreeRes, error)
	// ListItem 查询部门列表项，用于下拉选择
	ListItem(ctx context.Context) ([]*vo.DeptItemRes, error)
}

// deptService 部门服务实现
type deptService struct {
	deptRepo repo.DeptRepository
}

// Create 创建部门
func (d *deptService) Create(ctx context.Context, req *dto.CreateDeptReq) error {
	logger.Info("创建部门",
		logger.String("部门名称", req.Name),
		logger.String("部门编码", utils.SafeString(req.Code)),
		logger.Int64("父部门ID", utils.SafeInt64(req.ParentID)))

	// 构建部门实体
	dept := &entity.Dept{
		Name:      req.Name,
		Code:      req.Code,
		Pid:       req.ParentID,
		Ancestors: req.Ancestors,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
		CreatedBy: req.CreatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		dept.Sort = &sort
	}

	// 调用仓储层创建部门
	return d.deptRepo.Create(ctx, dept)
}

// Update 更新部门
func (d *deptService) Update(ctx context.Context, req *dto.UpdateDeptReq) error {
	logger.Info("更新部门",
		logger.Int64("部门ID", req.ID),
		logger.String("部门名称", req.Name),
		logger.String("部门编码", utils.SafeString(req.Code)),
		logger.Int64("父部门ID", utils.SafeInt64(req.ParentID)))

	// 构建部门实体
	dept := &entity.Dept{
		ID:        req.ID,
		Name:      req.Name,
		Code:      req.Code,
		Pid:       req.ParentID,
		Ancestors: req.Ancestors,
		Leader:    req.Leader,
		Phone:     req.Phone,
		Email:     req.Email,
		Status:    req.Status,
		UpdatedBy: req.UpdatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		dept.Sort = &sort
	}

	// 调用仓储层更新部门
	return d.deptRepo.Update(ctx, dept)
}

// Delete 删除部门
func (d *deptService) Delete(ctx context.Context, req *common.IdsReq) error {
	logger.Info("删除部门",
		logger.Any("部门ID列表", req.Ids),
		logger.Int("删除数量", len(req.Ids)))

	// 调用仓储层删除部门
	return d.deptRepo.Delete(ctx, req.Ids)
}

// FindOne 查询部门
func (d *deptService) FindOne(ctx context.Context, req *common.IDReq) (*vo.DeptDetailRes, error) {
	logger.Info("查询部门", logger.Int64("部门ID", req.ID))

	// 调用仓储层查询部门
	dept, err := d.deptRepo.FindOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 转换为VO
	res := &vo.DeptDetailRes{
		ID:        dept.ID,
		ParentID:  dept.Pid,
		Name:      dept.Name,
		Code:      dept.Code,
		Ancestors: dept.Ancestors,
		Leader:    dept.Leader,
		Phone:     dept.Phone,
		Email:     dept.Email,
		Sort:      dept.Sort,
		Status:    dept.Status,
		CreatedBy: dept.CreatedBy,
		UpdatedBy: dept.UpdatedBy,
		CreatedAt: dept.CreatedAt,
		UpdatedAt: dept.UpdatedAt,
	}

	return res, nil
}

// FindList 查询部门列表
func (d *deptService) FindList(ctx context.Context, req *dto.FindListDeptReq) (*common.Page, error) {
	logger.Info("查询部门列表",
		logger.String("部门名称", req.Name),
		logger.String("部门编码", req.Code),
		logger.Int64("父部门ID", utils.SafeInt64(req.ParentID)),
		logger.Int("页码", req.Page),
		logger.Int("每页数量", req.Size))

	// 构建查询参数
	params := &repo.DeptListParams{
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		ParentID:  req.ParentID,
		ExcludeID: req.ExcludeID,
		Page:      req.Page,
		Size:      req.Size,
	}

	// 调用仓储层查询部门列表
	depts, err := d.deptRepo.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建VO列表
	list := make([]*vo.DeptDetailRes, 0, len(depts))
	for _, dept := range depts {
		list = append(list, &vo.DeptDetailRes{
			ID:        dept.ID,
			ParentID:  dept.Pid,
			Name:      dept.Name,
			Code:      dept.Code,
			Ancestors: dept.Ancestors,
			Leader:    dept.Leader,
			Phone:     dept.Phone,
			Email:     dept.Email,
			Sort:      dept.Sort,
			Status:    dept.Status,
			CreatedBy: dept.CreatedBy,
			UpdatedBy: dept.UpdatedBy,
			CreatedAt: dept.CreatedAt,
			UpdatedAt: dept.UpdatedAt,
		})
	}

	// 简单处理，因为FindList返回的是全部列表，没有分页信息
	total := int64(len(list))

	// 构建分页响应
	return common.NewPage(list, total), nil
}

// Tree 查询部门树
func (d *deptService) Tree(ctx context.Context, req *dto.DeptTreeReq) ([]*vo.DeptTreeRes, error) {
	logger.Info("查询部门树",
		logger.String("部门名称", req.Name),
		logger.String("部门编码", req.Code),
		logger.Int64("父部门ID", utils.SafeInt64(req.ParentID)))

	// 构建查询参数
	params := &repo.DeptListParams{
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		ParentID:  req.ParentID,
		ExcludeID: req.ExcludeID,
	}

	// 调用仓储层查询部门树
	depts, err := d.deptRepo.Tree(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建VO树
	return buildDeptTreeVO(depts), nil
}

// SubDept 查询子部门
func (d *deptService) SubDept(ctx context.Context, req *dto.SubDeptReq) (*vo.DeptTreeRes, error) {
	logger.Info("查询子部门", logger.Int64("部门ID", req.ID))

	// 调用仓储层查询子部门
	dept, err := d.deptRepo.SubDept(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 转换为VO
	return convertDeptToTreeVO(dept), nil
}

// ListItem 查询部门列表项，用于下拉选择
func (d *deptService) ListItem(ctx context.Context) ([]*vo.DeptItemRes, error) {
	logger.Info("查询部门列表项")

	// 构建查询参数，只查询状态正常的部门
	status := int64(1) // 正常状态的部门
	params := &repo.DeptListParams{
		Status: &status,
	}

	// 调用仓储层查询部门列表
	depts, err := d.deptRepo.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建VO列表
	list := make([]*vo.DeptItemRes, 0, len(depts))
	for _, dept := range depts {
		list = append(list, &vo.DeptItemRes{
			ID:   dept.ID,
			Name: dept.Name,
		})
	}

	return list, nil
}

// buildDeptTreeVO 构建部门树VO
func buildDeptTreeVO(depts []*entity.Dept) []*vo.DeptTreeRes {
	voList := make([]*vo.DeptTreeRes, 0, len(depts))
	for _, dept := range depts {
		deptVO := convertDeptToTreeVO(dept)
		voList = append(voList, deptVO)
	}
	return voList
}

// convertDeptToTreeVO 将部门实体转换为树形VO
func convertDeptToTreeVO(dept *entity.Dept) *vo.DeptTreeRes {
	if dept == nil {
		return nil
	}

	res := &vo.DeptTreeRes{
		ID:        dept.ID,
		ParentID:  dept.Pid,
		Name:      dept.Name,
		Code:      dept.Code,
		Ancestors: dept.Ancestors,
		Leader:    dept.Leader,
		Phone:     dept.Phone,
		Email:     dept.Email,
		Sort:      dept.Sort,
		Status:    dept.Status,
		CreatedAt: dept.CreatedAt,
	}

	// 递归处理子部门
	if len(dept.Children) > 0 {
		children := make([]*vo.DeptTreeRes, 0, len(dept.Children))
		for _, child := range dept.Children {
			children = append(children, convertDeptToTreeVO(child))
		}
		res.Children = children
	}

	return res
}

// NewDeptService 创建部门服务
func NewDeptService(deptRepo repo.DeptRepository) DeptService {
	return &deptService{
		deptRepo: deptRepo,
	}
}
