package repository

import (
	"context"
	"sweet/pkg/errs"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound   = errs.New(6000, "岗位不存在")
	ErrPostExists     = errs.New(6001, "岗位已存在")
	ErrPostHasUsers   = errs.New(6002, "岗位下存在用户，无法删除")
	ErrInvalidPostID  = errs.New(6003, "无效岗位ID")
	ErrPostCodeExists = errs.New(6004, "岗位编码已存在")
)

// PostRepository 岗位仓储接口
type PostRepository interface {
	// Create 创建岗位
	Create(ctx context.Context, post *entity.Post) error
	// Update 更新岗位
	Update(ctx context.Context, post *entity.Post) error
	// Delete 删除岗位
	Delete(ctx context.Context, ids []int64) error
	// FindOne 查询岗位
	FindOne(ctx context.Context, id int64) (*entity.Post, error)
	// FindByCode 通过编码查询岗位
	FindByCode(ctx context.Context, code string) (*entity.Post, error)
	// FindList 查询岗位列表
	FindList(ctx context.Context, params *PostListParams) (list []*entity.Post, total int64, err error)
	// CheckHasUsers 检查岗位下是否有用户
	CheckHasUsers(ctx context.Context, postID int64) (bool, error)
}

// PostListParams 岗位列表查询参数
type PostListParams struct {
	Name   string // 岗位名称
	Code   string // 岗位编码
	Status *int64 // 状态
	Page   int    // 页码
	Size   int    // 每页数量
}

// postRepository 岗位仓储实现
type postRepository struct {
	db *gorm.DB
	q  *query.Query
}

func (p *postRepository) Create(ctx context.Context, post *entity.Post) error {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) Update(ctx context.Context, post *entity.Post) error {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) Delete(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) FindOne(ctx context.Context, id int64) (*entity.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) FindByCode(ctx context.Context, code string) (*entity.Post, error) {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) FindList(ctx context.Context, params *PostListParams) (list []*entity.Post, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (p *postRepository) CheckHasUsers(ctx context.Context, postID int64) (bool, error) {
	//TODO implement me
	panic("implement me")
}

// NewPostRepository 创建岗位仓储
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
		q:  query.Use(db),
	}
}
