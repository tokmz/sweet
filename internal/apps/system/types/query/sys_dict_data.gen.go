// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package query

import (
	"context"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"gorm.io/gorm/schema"

	"gorm.io/gen"
	"gorm.io/gen/field"

	"gorm.io/plugin/dbresolver"

	"sweet/internal/apps/system/types/entity"
)

func newDictData(db *gorm.DB, opts ...gen.DOOption) dictData {
	_dictData := dictData{}

	_dictData.dictDataDo.UseDB(db, opts...)
	_dictData.dictDataDo.UseModel(&entity.DictData{})

	tableName := _dictData.dictDataDo.TableName()
	_dictData.ALL = field.NewAsterisk(tableName)
	_dictData.ID = field.NewInt64(tableName, "id")
	_dictData.DictType = field.NewString(tableName, "dict_type")
	_dictData.DictLabel = field.NewString(tableName, "dict_label")
	_dictData.DictValue = field.NewString(tableName, "dict_value")
	_dictData.DictSort = field.NewInt64(tableName, "dict_sort")
	_dictData.CSSClass = field.NewString(tableName, "css_class")
	_dictData.ListClass = field.NewString(tableName, "list_class")
	_dictData.IsDefault = field.NewInt64(tableName, "is_default")
	_dictData.Status = field.NewInt64(tableName, "status")
	_dictData.Remark = field.NewString(tableName, "remark")
	_dictData.CreatedBy = field.NewInt64(tableName, "created_by")
	_dictData.UpdatedBy = field.NewInt64(tableName, "updated_by")
	_dictData.CreatedAt = field.NewTime(tableName, "created_at")
	_dictData.UpdatedAt = field.NewTime(tableName, "updated_at")
	_dictData.DeletedAt = field.NewField(tableName, "deleted_at")
	_dictData.Version = field.NewInt64(tableName, "version")

	_dictData.fillFieldMap()

	return _dictData
}

// dictData 字典数据表
type dictData struct {
	dictDataDo

	ALL       field.Asterisk
	ID        field.Int64  // 字典数据ID
	DictType  field.String // 字典类型
	DictLabel field.String // 字典标签
	DictValue field.String // 字典值
	DictSort  field.Int64  // 排序
	CSSClass  field.String // 样式属性
	ListClass field.String // 表格回显样式
	IsDefault field.Int64  // 是否默认 1-是 2-否
	Status    field.Int64  // 状态 1-正常 2-禁用
	Remark    field.String // 备注
	CreatedBy field.Int64  // 创建人
	UpdatedBy field.Int64  // 更新人
	CreatedAt field.Time   // 创建时间
	UpdatedAt field.Time   // 更新时间
	DeletedAt field.Field  // 删除时间
	Version   field.Int64  // 乐观锁版本

	fieldMap map[string]field.Expr
}

func (d dictData) Table(newTableName string) *dictData {
	d.dictDataDo.UseTable(newTableName)
	return d.updateTableName(newTableName)
}

func (d dictData) As(alias string) *dictData {
	d.dictDataDo.DO = *(d.dictDataDo.As(alias).(*gen.DO))
	return d.updateTableName(alias)
}

func (d *dictData) updateTableName(table string) *dictData {
	d.ALL = field.NewAsterisk(table)
	d.ID = field.NewInt64(table, "id")
	d.DictType = field.NewString(table, "dict_type")
	d.DictLabel = field.NewString(table, "dict_label")
	d.DictValue = field.NewString(table, "dict_value")
	d.DictSort = field.NewInt64(table, "dict_sort")
	d.CSSClass = field.NewString(table, "css_class")
	d.ListClass = field.NewString(table, "list_class")
	d.IsDefault = field.NewInt64(table, "is_default")
	d.Status = field.NewInt64(table, "status")
	d.Remark = field.NewString(table, "remark")
	d.CreatedBy = field.NewInt64(table, "created_by")
	d.UpdatedBy = field.NewInt64(table, "updated_by")
	d.CreatedAt = field.NewTime(table, "created_at")
	d.UpdatedAt = field.NewTime(table, "updated_at")
	d.DeletedAt = field.NewField(table, "deleted_at")
	d.Version = field.NewInt64(table, "version")

	d.fillFieldMap()

	return d
}

func (d *dictData) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := d.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (d *dictData) fillFieldMap() {
	d.fieldMap = make(map[string]field.Expr, 17)
	d.fieldMap["id"] = d.ID
	d.fieldMap["dict_type"] = d.DictType
	d.fieldMap["dict_label"] = d.DictLabel
	d.fieldMap["dict_value"] = d.DictValue
	d.fieldMap["dict_sort"] = d.DictSort
	d.fieldMap["css_class"] = d.CSSClass
	d.fieldMap["list_class"] = d.ListClass
	d.fieldMap["is_default"] = d.IsDefault
	d.fieldMap["status"] = d.Status
	d.fieldMap["remark"] = d.Remark
	d.fieldMap["created_by"] = d.CreatedBy
	d.fieldMap["updated_by"] = d.UpdatedBy
	d.fieldMap["created_at"] = d.CreatedAt
	d.fieldMap["updated_at"] = d.UpdatedAt
	d.fieldMap["deleted_at"] = d.DeletedAt
	d.fieldMap["version"] = d.Version

}

func (d dictData) clone(db *gorm.DB) dictData {
	d.dictDataDo.ReplaceConnPool(db.Statement.ConnPool)
	return d
}

func (d dictData) replaceDB(db *gorm.DB) dictData {
	d.dictDataDo.ReplaceDB(db)
	return d
}

type dictDataDo struct{ gen.DO }

func (d dictDataDo) Debug() *dictDataDo {
	return d.withDO(d.DO.Debug())
}

func (d dictDataDo) WithContext(ctx context.Context) *dictDataDo {
	return d.withDO(d.DO.WithContext(ctx))
}

func (d dictDataDo) ReadDB() *dictDataDo {
	return d.Clauses(dbresolver.Read)
}

func (d dictDataDo) WriteDB() *dictDataDo {
	return d.Clauses(dbresolver.Write)
}

func (d dictDataDo) Session(config *gorm.Session) *dictDataDo {
	return d.withDO(d.DO.Session(config))
}

func (d dictDataDo) Clauses(conds ...clause.Expression) *dictDataDo {
	return d.withDO(d.DO.Clauses(conds...))
}

func (d dictDataDo) Returning(value interface{}, columns ...string) *dictDataDo {
	return d.withDO(d.DO.Returning(value, columns...))
}

func (d dictDataDo) Not(conds ...gen.Condition) *dictDataDo {
	return d.withDO(d.DO.Not(conds...))
}

func (d dictDataDo) Or(conds ...gen.Condition) *dictDataDo {
	return d.withDO(d.DO.Or(conds...))
}

func (d dictDataDo) Select(conds ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.Select(conds...))
}

func (d dictDataDo) Where(conds ...gen.Condition) *dictDataDo {
	return d.withDO(d.DO.Where(conds...))
}

func (d dictDataDo) Order(conds ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.Order(conds...))
}

func (d dictDataDo) Distinct(cols ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.Distinct(cols...))
}

func (d dictDataDo) Omit(cols ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.Omit(cols...))
}

func (d dictDataDo) Join(table schema.Tabler, on ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.Join(table, on...))
}

func (d dictDataDo) LeftJoin(table schema.Tabler, on ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.LeftJoin(table, on...))
}

func (d dictDataDo) RightJoin(table schema.Tabler, on ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.RightJoin(table, on...))
}

func (d dictDataDo) Group(cols ...field.Expr) *dictDataDo {
	return d.withDO(d.DO.Group(cols...))
}

func (d dictDataDo) Having(conds ...gen.Condition) *dictDataDo {
	return d.withDO(d.DO.Having(conds...))
}

func (d dictDataDo) Limit(limit int) *dictDataDo {
	return d.withDO(d.DO.Limit(limit))
}

func (d dictDataDo) Offset(offset int) *dictDataDo {
	return d.withDO(d.DO.Offset(offset))
}

func (d dictDataDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *dictDataDo {
	return d.withDO(d.DO.Scopes(funcs...))
}

func (d dictDataDo) Unscoped() *dictDataDo {
	return d.withDO(d.DO.Unscoped())
}

func (d dictDataDo) Create(values ...*entity.DictData) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Create(values)
}

func (d dictDataDo) CreateInBatches(values []*entity.DictData, batchSize int) error {
	return d.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (d dictDataDo) Save(values ...*entity.DictData) error {
	if len(values) == 0 {
		return nil
	}
	return d.DO.Save(values)
}

func (d dictDataDo) First() (*entity.DictData, error) {
	if result, err := d.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.DictData), nil
	}
}

func (d dictDataDo) Take() (*entity.DictData, error) {
	if result, err := d.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.DictData), nil
	}
}

func (d dictDataDo) Last() (*entity.DictData, error) {
	if result, err := d.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.DictData), nil
	}
}

func (d dictDataDo) Find() ([]*entity.DictData, error) {
	result, err := d.DO.Find()
	return result.([]*entity.DictData), err
}

func (d dictDataDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.DictData, err error) {
	buf := make([]*entity.DictData, 0, batchSize)
	err = d.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (d dictDataDo) FindInBatches(result *[]*entity.DictData, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return d.DO.FindInBatches(result, batchSize, fc)
}

func (d dictDataDo) Attrs(attrs ...field.AssignExpr) *dictDataDo {
	return d.withDO(d.DO.Attrs(attrs...))
}

func (d dictDataDo) Assign(attrs ...field.AssignExpr) *dictDataDo {
	return d.withDO(d.DO.Assign(attrs...))
}

func (d dictDataDo) Joins(fields ...field.RelationField) *dictDataDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Joins(_f))
	}
	return &d
}

func (d dictDataDo) Preload(fields ...field.RelationField) *dictDataDo {
	for _, _f := range fields {
		d = *d.withDO(d.DO.Preload(_f))
	}
	return &d
}

func (d dictDataDo) FirstOrInit() (*entity.DictData, error) {
	if result, err := d.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.DictData), nil
	}
}

func (d dictDataDo) FirstOrCreate() (*entity.DictData, error) {
	if result, err := d.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.DictData), nil
	}
}

func (d dictDataDo) FindByPage(offset int, limit int) (result []*entity.DictData, count int64, err error) {
	result, err = d.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = d.Offset(-1).Limit(-1).Count()
	return
}

func (d dictDataDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = d.Count()
	if err != nil {
		return
	}

	err = d.Offset(offset).Limit(limit).Scan(result)
	return
}

func (d dictDataDo) Scan(result interface{}) (err error) {
	return d.DO.Scan(result)
}

func (d dictDataDo) Delete(models ...*entity.DictData) (result gen.ResultInfo, err error) {
	return d.DO.Delete(models)
}

func (d *dictDataDo) withDO(do gen.Dao) *dictDataDo {
	d.DO = *do.(*gen.DO)
	return d
}
