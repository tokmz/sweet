package repo

import (
	"context"
	"errors"
	"sweet/pkg/errs"
	"sweet/pkg/logger"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrPostNotFound   = errs.New(6001, "岗位不存在")
	ErrPostNameExists = errs.New(6002, "岗位名称已存在")
	ErrPostHasUsers   = errs.New(6003, "岗位下存在用户，无法删除")
	ErrInvalidPostID  = errs.New(6004, "无效岗位ID")
	ErrPostCodeExists = errs.New(6005, "岗位编码已存在")
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
	// FindList 查询岗位列表
	FindList(ctx context.Context, params *PostListParams) (list []*entity.Post, total int64, err error)
	// ScanList 查询岗位列表
	ScanList(ctx context.Context, params *PostListParams, list any, total *int64) error
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
	return p.q.Transaction(func(tx *query.Query) error {
		do := tx.Post.WithContext(ctx)
		field := tx.Post

		// 检查岗位名称是否已存在
		if _, err := do.Where(field.Name.Eq(post.Name)).First(); err == nil {
			return ErrPostNameExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询岗位失败",
				logger.Err(err),
				logger.String("岗位名称", post.Name),
			)
			return errs.ErrServer
		}

		// 检查岗位编码是否已存在
		if _, err := do.Where(field.Code.Eq(post.Code)).First(); err == nil {
			return ErrPostCodeExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"查询岗位失败",
				logger.Err(err),
				logger.String("岗位编码", post.Code),
			)
			return errs.ErrServer
		}

		// 创建岗位
		if err := do.Create(post); err != nil {
			logger.Error(
				"创建岗位失败",
				logger.Err(err),
				logger.String("岗位名称", post.Name),
				logger.String("岗位编码", post.Code),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (p *postRepository) Update(ctx context.Context, post *entity.Post) error {
	return p.q.Transaction(func(tx *query.Query) error {
		do := tx.Post.WithContext(ctx)
		field := tx.Post

		// 检查要更新的岗位是否存在
		if _, err := do.Where(field.ID.Eq(post.ID)).First(); err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return ErrPostNotFound
			}
			logger.Error(
				"查询岗位失败",
				logger.Err(err),
				logger.Int64("岗位ID", post.ID),
			)
			return errs.ErrServer
		}

		// 检查岗位名称是否重复
		if _, err := do.Where(field.Name.Eq(post.Name), field.ID.Neq(post.ID)).First(); err == nil {
			return ErrPostNameExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"检查岗位名称唯一性失败",
				logger.Err(err),
				logger.Int64("岗位ID", post.ID),
				logger.String("岗位名称", post.Name),
			)
			return errs.ErrServer
		}

		// 检查岗位编码是否重复
		if _, err := do.Where(field.Code.Eq(post.Code), field.ID.Neq(post.ID)).First(); err == nil {
			return ErrPostCodeExists
		} else if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			logger.Error(
				"检查岗位编码唯一性失败",
				logger.Err(err),
				logger.Int64("岗位ID", post.ID),
				logger.String("岗位编码", post.Code),
			)
			return errs.ErrServer
		}

		// 更新岗位
		if _, err := do.Where(field.ID.Eq(post.ID)).Updates(post); err != nil {
			logger.Error(
				"更新岗位失败",
				logger.Err(err),
				logger.Int64("岗位ID", post.ID),
				logger.String("岗位名称", post.Name),
				logger.String("岗位编码", post.Code),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (p *postRepository) Delete(ctx context.Context, ids []int64) error {
	return p.q.Transaction(func(tx *query.Query) error {
		do := tx.Post.WithContext(ctx)
		userDo := tx.User.WithContext(ctx)

		for _, id := range ids {
			// 检查岗位下是否有用户
			count, err := userDo.Where(tx.User.PostID.Eq(id)).Count()
			if err != nil {
				logger.Error(
					"查询岗位用户数量失败",
					logger.Err(err),
					logger.Int64("岗位ID", id),
				)
				return errs.ErrServer
			}
			if count > 0 {
				return ErrPostHasUsers
			}
		}

		// 删除岗位
		if _, err := do.Where(tx.Post.ID.In(ids...)).Delete(); err != nil {
			logger.Error(
				"删除岗位失败",
				logger.Err(err),
				logger.Any("岗位ID列表", ids),
			)
			return errs.ErrServer
		}

		return nil
	})
}

func (p *postRepository) FindOne(ctx context.Context, id int64) (*entity.Post, error) {
	if id <= 0 {
		return nil, ErrInvalidPostID
	}

	do := p.q.Post.WithContext(ctx)
	field := p.q.Post
	post, err := do.Where(field.ID.Eq(id)).First()
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrPostNotFound
		}
		logger.Error(
			"查询岗位失败",
			logger.Err(err),
			logger.Int64("岗位ID", id),
		)
		return nil, errs.ErrServer
	}
	return post, nil
}

func (p *postRepository) FindList(ctx context.Context, params *PostListParams) (list []*entity.Post, total int64, err error) {
	do := p.q.Post.WithContext(ctx)
	field := p.q.Post

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

	// 默认按照排序字段和ID升序排列
	do = do.Order(field.Sort, field.ID)

	// 计算总数
	total, err = do.Count()
	if err != nil {
		logger.Error(
			"查询岗位总数失败",
			logger.Err(err),
			logger.String("岗位名称", params.Name),
			logger.String("岗位编码", params.Code),
		)
		return nil, 0, errs.ErrServer
	}

	// 分页查询
	if params.Page > 0 && params.Size > 0 {
		// 计算偏移量
		offset := (params.Page - 1) * params.Size
		list, err = do.Offset(offset).Limit(params.Size).Find()
	} else {
		list, err = do.Find()
	}

	if err != nil {
		logger.Error(
			"查询岗位列表失败",
			logger.Err(err),
			logger.String("岗位名称", params.Name),
			logger.String("岗位编码", params.Code),
			logger.Int("页码", params.Page),
			logger.Int("每页数量", params.Size),
		)
		return nil, 0, errs.ErrServer
	}

	return list, total, nil
}

func (p *postRepository) ScanList(ctx context.Context, params *PostListParams, list any, total *int64) error {
	do := p.q.Post.WithContext(ctx)
	field := p.q.Post

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

	// 默认按照排序字段和ID升序排列
	do = do.Order(field.Sort, field.ID)

	// 计算总数
	count, err := do.Count()
	if err != nil {
		logger.Error(
			"查询岗位总数失败",
			logger.Err(err),
			logger.String("岗位名称", params.Name),
			logger.String("岗位编码", params.Code),
		)
		return errs.ErrServer
	}
	*total = count

	// 分页查询
	if params.Page > 0 && params.Size > 0 {
		// 计算偏移量
		offset := (params.Page - 1) * params.Size
		if err := do.Offset(offset).Limit(params.Size).Scan(list); err != nil {
			logger.Error(
				"查询岗位列表失败",
				logger.Err(err),
				logger.String("岗位名称", params.Name),
				logger.String("岗位编码", params.Code),
				logger.Int("页码", params.Page),
				logger.Int("每页数量", params.Size),
			)
			return errs.ErrServer
		}
	} else {
		if err := do.Scan(list); err != nil {
			logger.Error(
				"查询岗位列表失败",
				logger.Err(err),
				logger.String("岗位名称", params.Name),
				logger.String("岗位编码", params.Code),
			)
			return errs.ErrServer
		}
	}

	return nil
}

// NewPostRepository 创建岗位仓储
func NewPostRepository(db *gorm.DB) PostRepository {
	return &postRepository{
		db: db,
		q:  query.Use(db),
	}
}
