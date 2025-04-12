package repo

import (
	"context"
	"sweet/pkg/errs"

	"sweet/internal/apps/system/types/entity"
	"sweet/internal/apps/system/types/query"

	"gorm.io/gorm"
)

var (
	ErrDictTypeNotFound   = errs.New(7000, "字典类型不存在")
	ErrDictTypeExists     = errs.New(7001, "字典类型已存在")
	ErrDictTypeHasData    = errs.New(7002, "字典类型下存在字典数据，无法删除")
	ErrInvalidDictTypeID  = errs.New(7003, "无效字典类型ID")
	ErrDictTypeCodeExists = errs.New(7004, "字典类型编码已存在")
	ErrDictDataNotFound   = errs.New(7010, "字典数据不存在")
	ErrDictDataExists     = errs.New(7011, "字典数据已存在")
	ErrInvalidDictDataID  = errs.New(7012, "无效字典数据ID")
)

// DictRepository 字典仓储接口
type DictRepository interface {
	/* 字典类型 */

	// CreateType 创建字典类型
	CreateType(ctx context.Context, dictType *entity.DictType) error
	// UpdateType 更新字典类型
	UpdateType(ctx context.Context, dictType *entity.DictType) error
	// DeleteType 删除字典类型
	DeleteType(ctx context.Context, ids []int64) error
	// FindType 查询字典类型
	FindType(ctx context.Context, id int64) (*entity.DictType, error)
	// FindTypeByType 通过类型查询字典类型
	FindTypeByType(ctx context.Context, dictType string) (*entity.DictType, error)
	// FindTypeList 查询字典类型列表
	FindTypeList(ctx context.Context, params *DictTypeListParams) (list []*entity.DictType, total int64, err error)
	// CheckTypeHasData 检查字典类型下是否有字典数据
	CheckTypeHasData(ctx context.Context, dictType string) (bool, error)

	/* 字典数据 */

	// CreateData 创建字典数据
	CreateData(ctx context.Context, dictData *entity.DictData) error
	// UpdateData 更新字典数据
	UpdateData(ctx context.Context, dictData *entity.DictData) error
	// DeleteData 删除字典数据
	DeleteData(ctx context.Context, ids []int64) error
	// FindData 查询字典数据
	FindData(ctx context.Context, id int64) (*entity.DictData, error)
	// FindDataList 查询字典数据列表
	FindDataList(ctx context.Context, params *DictDataListParams) (list []*entity.DictData, total int64, err error)
	// FindDataByType 通过字典类型查询字典数据列表
	FindDataByType(ctx context.Context, dictType string) ([]*entity.DictData, error)
	// FindDataByTypeAndValue 通过字典类型和字典值查询字典数据
	FindDataByTypeAndValue(ctx context.Context, dictType string, dictValue string) (*entity.DictData, error)
}

// DictTypeListParams 字典类型列表查询参数
type DictTypeListParams struct {
	Name   string // 字典名称
	Type   string // 字典类型
	Status *int64 // 状态
	Page   int    // 页码
	Size   int    // 每页数量
}

// DictDataListParams 字典数据列表查询参数
type DictDataListParams struct {
	DictType  string // 字典类型
	DictLabel string // 字典标签
	Status    *int64 // 状态
	Page      int    // 页码
	Size      int    // 每页数量
}

// dictRepository 字典仓储实现
type dictRepository struct {
	db *gorm.DB
	q  *query.Query
}

/* 字典类型 */

func (d *dictRepository) CreateType(ctx context.Context, dictType *entity.DictType) error {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) UpdateType(ctx context.Context, dictType *entity.DictType) error {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) DeleteType(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) FindType(ctx context.Context, id int64) (*entity.DictType, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) FindTypeByType(ctx context.Context, dictType string) (*entity.DictType, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) FindTypeList(ctx context.Context, params *DictTypeListParams) (list []*entity.DictType, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) CheckTypeHasData(ctx context.Context, dictType string) (bool, error) {
	//TODO implement me
	panic("implement me")
}

/* 字典数据 */

func (d *dictRepository) CreateData(ctx context.Context, dictData *entity.DictData) error {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) UpdateData(ctx context.Context, dictData *entity.DictData) error {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) DeleteData(ctx context.Context, ids []int64) error {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) FindData(ctx context.Context, id int64) (*entity.DictData, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) FindDataList(ctx context.Context, params *DictDataListParams) (list []*entity.DictData, total int64, err error) {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) FindDataByType(ctx context.Context, dictType string) ([]*entity.DictData, error) {
	//TODO implement me
	panic("implement me")
}

func (d *dictRepository) FindDataByTypeAndValue(ctx context.Context, dictType string, dictValue string) (*entity.DictData, error) {
	//TODO implement me
	panic("implement me")
}

// NewDictRepository 创建字典仓储
func NewDictRepository(db *gorm.DB) DictRepository {
	return &dictRepository{
		db: db,
		q:  query.Use(db),
	}
}
