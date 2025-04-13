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

// PostService 岗位服务接口
type PostService interface {
	// Create 创建岗位
	Create(ctx context.Context, req *dto.CreatePostReq) error
	// Update 更新岗位
	Update(ctx context.Context, req *dto.UpdatePostReq) error
	// Delete 删除岗位
	Delete(ctx context.Context, req *common.IdsReq) error
	// FindOne 查询岗位
	FindOne(ctx context.Context, req *common.IDReq) (*vo.PostDetailRes, error)
	// FindList 查询岗位列表
	FindList(ctx context.Context, req *dto.FindListPostReq) (*common.Page, error)
	// ListItem 查询岗位列表项，用于下拉选择
	ListItem(ctx context.Context) ([]*vo.PostItemRes, error)
	// UpdateStatus 批量更新岗位状态
	UpdateStatus(ctx context.Context, req *common.StatusReq) error
}

// postService 岗位服务实现
type postService struct {
	postRepo repo.PostRepository
}

// Create 创建岗位
func (p *postService) Create(ctx context.Context, req *dto.CreatePostReq) error {
	logger.Info("创建岗位",
		logger.String("岗位名称", req.Name),
		logger.String("岗位编码", req.Code))

	// 构建岗位实体
	post := &entity.Post{
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		Remark:    req.Remark,
		CreatedBy: req.CreatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		post.Sort = &sort
	}

	// 调用仓储层创建岗位
	return p.postRepo.Create(ctx, post)
}

// Update 更新岗位
func (p *postService) Update(ctx context.Context, req *dto.UpdatePostReq) error {
	logger.Info("更新岗位",
		logger.Int64("岗位ID", req.ID),
		logger.String("岗位名称", req.Name),
		logger.String("岗位编码", req.Code))

	// 构建岗位实体
	post := &entity.Post{
		ID:        req.ID,
		Name:      req.Name,
		Code:      req.Code,
		Status:    req.Status,
		Remark:    req.Remark,
		UpdatedBy: req.UpdatedBy,
	}

	// 设置排序值
	if req.Sort != nil {
		sort := int64(*req.Sort)
		post.Sort = &sort
	}

	// 调用仓储层更新岗位
	return p.postRepo.Update(ctx, post)
}

// Delete 删除岗位
func (p *postService) Delete(ctx context.Context, req *common.IdsReq) error {
	logger.Info("删除岗位",
		logger.Any("岗位ID列表", req.Ids),
		logger.Int("删除数量", len(req.Ids)))

	// 调用仓储层删除岗位
	return p.postRepo.Delete(ctx, req.Ids)
}

// FindOne 查询岗位
func (p *postService) FindOne(ctx context.Context, req *common.IDReq) (*vo.PostDetailRes, error) {
	logger.Info("查询岗位", logger.Int64("岗位ID", req.ID))

	// 调用仓储层查询岗位
	post, err := p.postRepo.FindOne(ctx, req.ID)
	if err != nil {
		return nil, err
	}

	// 转换为VO
	res := &vo.PostDetailRes{
		ID:        post.ID,
		Name:      post.Name,
		Code:      post.Code,
		Sort:      post.Sort,
		Status:    post.Status,
		Remark:    post.Remark,
		CreatedBy: post.CreatedBy,
		UpdatedBy: post.UpdatedBy,
		CreatedAt: post.CreatedAt,
		UpdatedAt: post.UpdatedAt,
	}

	return res, nil
}

// FindList 查询岗位列表
func (p *postService) FindList(ctx context.Context, req *dto.FindListPostReq) (*common.Page, error) {
	logger.Info("查询岗位列表",
		logger.String("岗位名称", req.Name),
		logger.String("岗位编码", req.Code),
		logger.Int("页码", req.Page),
		logger.Int("每页数量", req.Size))

	// 构建查询参数
	params := &repo.PostListParams{
		Name:   req.Name,
		Code:   req.Code,
		Status: req.Status,
		Page:   req.Page,
		Size:   req.Size,
	}

	// 调用仓储层查询岗位列表
	list, total, err := p.postRepo.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建VO列表
	voList := make([]*vo.PostListRes, 0, len(list))
	for _, post := range list {
		voList = append(voList, &vo.PostListRes{
			ID:        post.ID,
			Name:      post.Name,
			Code:      post.Code,
			Sort:      post.Sort,
			Status:    post.Status,
			Remark:    post.Remark,
			CreatedAt: post.CreatedAt,
		})
	}

	// 构建分页响应
	return common.NewPage(voList, total), nil
}

// ListItem 查询岗位列表项，用于下拉选择
func (p *postService) ListItem(ctx context.Context) ([]*vo.PostItemRes, error) {
	logger.Info("查询岗位列表项")

	// 构建查询参数，只查询状态正常的岗位
	status := int64(1) // 正常状态的岗位
	params := &repo.PostListParams{
		Status: &status,
		Page:   1,
		Size:   1000, // 设置较大值返回全部数据
	}

	// 调用仓储层查询岗位列表
	list, _, err := p.postRepo.FindList(ctx, params)
	if err != nil {
		return nil, err
	}

	// 构建VO列表
	voList := make([]*vo.PostItemRes, 0, len(list))
	for _, post := range list {
		voList = append(voList, &vo.PostItemRes{
			ID:   post.ID,
			Name: post.Name,
		})
	}

	return voList, nil
}

// UpdateStatus 批量更新岗位状态
func (p *postService) UpdateStatus(ctx context.Context, req *common.StatusReq) error {
	logger.Info("批量更新岗位状态",
		logger.Any("岗位ID列表", req.Ids),
		logger.Int64("状态", req.Status),
		logger.Int("数量", len(req.Ids)))

	// 遍历处理每个岗位
	for _, id := range req.Ids {
		// 先查询岗位
		post, err := p.postRepo.FindOne(ctx, id)
		if err != nil {
			logger.Error("查询岗位失败",
				logger.Err(err),
				logger.Int64("岗位ID", id))
			continue
		}

		// 更新状态
		status := req.Status
		post.Status = &status

		// 调用仓储层更新岗位
		if err = p.postRepo.Update(ctx, post); err != nil {
			logger.Error("更新岗位状态失败",
				logger.Err(err),
				logger.Int64("岗位ID", id),
				logger.Int64("状态", req.Status))
			return err
		}
	}

	return nil
}

// NewPostService 创建岗位服务
func NewPostService(postRepo repo.PostRepository) PostService {
	return &postService{
		postRepo: postRepo,
	}
}
