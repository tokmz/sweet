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

func newMenu(db *gorm.DB, opts ...gen.DOOption) menu {
	_menu := menu{}

	_menu.menuDo.UseDB(db, opts...)
	_menu.menuDo.UseModel(&entity.Menu{})

	tableName := _menu.menuDo.TableName()
	_menu.ALL = field.NewAsterisk(tableName)
	_menu.ID = field.NewInt64(tableName, "id")
	_menu.ParentID = field.NewInt64(tableName, "parent_id")
	_menu.Name = field.NewString(tableName, "name")
	_menu.Permission = field.NewString(tableName, "permission")
	_menu.Type = field.NewInt64(tableName, "type")
	_menu.Path = field.NewString(tableName, "path")
	_menu.Component = field.NewString(tableName, "component")
	_menu.Redirect = field.NewString(tableName, "redirect")
	_menu.Icon = field.NewString(tableName, "icon")
	_menu.Sort = field.NewInt64(tableName, "sort")
	_menu.Hidden = field.NewInt64(tableName, "hidden")
	_menu.Status = field.NewInt64(tableName, "status")
	_menu.AlwaysShow = field.NewInt64(tableName, "always_show")
	_menu.KeepAlive = field.NewInt64(tableName, "keep_alive")
	_menu.Target = field.NewString(tableName, "target")
	_menu.Title = field.NewString(tableName, "title")
	_menu.ActiveMenu = field.NewString(tableName, "active_menu")
	_menu.Breadcrumb = field.NewInt64(tableName, "breadcrumb")
	_menu.Affix = field.NewInt64(tableName, "affix")
	_menu.FrameSrc = field.NewString(tableName, "frame_src")
	_menu.FrameLoading = field.NewInt64(tableName, "frame_loading")
	_menu.Transition = field.NewString(tableName, "transition")
	_menu.Remark = field.NewString(tableName, "remark")
	_menu.CreatedBy = field.NewInt64(tableName, "created_by")
	_menu.UpdatedBy = field.NewInt64(tableName, "updated_by")
	_menu.CreatedAt = field.NewTime(tableName, "created_at")
	_menu.UpdatedAt = field.NewTime(tableName, "updated_at")
	_menu.DeletedAt = field.NewField(tableName, "deleted_at")
	_menu.Version = field.NewInt64(tableName, "version")

	_menu.fillFieldMap()

	return _menu
}

// menu 菜单权限表
type menu struct {
	menuDo

	ALL          field.Asterisk
	ID           field.Int64  // 菜单ID
	ParentID     field.Int64  // 父菜单ID
	Name         field.String // 菜单名称
	Permission   field.String // 权限标识
	Type         field.Int64  // 类型 1-目录 2-菜单 3-按钮
	Path         field.String // 路由地址
	Component    field.String // 组件路径
	Redirect     field.String // 重定向地址
	Icon         field.String // 菜单图标
	Sort         field.Int64  // 排序
	Hidden       field.Int64  // 是否隐藏 1-是 2-否
	Status       field.Int64  // 状态 1-正常 2-禁用
	AlwaysShow   field.Int64  // 是否总是显示 1-是 2-否
	KeepAlive    field.Int64  // 是否缓存 1-是 2-否
	Target       field.String // 打开方式 _self _blank
	Title        field.String // 菜单标题
	ActiveMenu   field.String // 激活菜单
	Breadcrumb   field.Int64  // 是否显示面包屑 1-是 2-否
	Affix        field.Int64  // 是否固定 1-是 2-否
	FrameSrc     field.String // iframe地址
	FrameLoading field.Int64  // iframe加载状态 1-是 2-否
	Transition   field.String // 过渡动画
	Remark       field.String // 备注
	CreatedBy    field.Int64  // 创建人
	UpdatedBy    field.Int64  // 更新人
	CreatedAt    field.Time   // 创建时间
	UpdatedAt    field.Time   // 更新时间
	DeletedAt    field.Field  // 删除时间
	Version      field.Int64  // 乐观锁版本

	fieldMap map[string]field.Expr
}

func (m menu) Table(newTableName string) *menu {
	m.menuDo.UseTable(newTableName)
	return m.updateTableName(newTableName)
}

func (m menu) As(alias string) *menu {
	m.menuDo.DO = *(m.menuDo.As(alias).(*gen.DO))
	return m.updateTableName(alias)
}

func (m *menu) updateTableName(table string) *menu {
	m.ALL = field.NewAsterisk(table)
	m.ID = field.NewInt64(table, "id")
	m.ParentID = field.NewInt64(table, "parent_id")
	m.Name = field.NewString(table, "name")
	m.Permission = field.NewString(table, "permission")
	m.Type = field.NewInt64(table, "type")
	m.Path = field.NewString(table, "path")
	m.Component = field.NewString(table, "component")
	m.Redirect = field.NewString(table, "redirect")
	m.Icon = field.NewString(table, "icon")
	m.Sort = field.NewInt64(table, "sort")
	m.Hidden = field.NewInt64(table, "hidden")
	m.Status = field.NewInt64(table, "status")
	m.AlwaysShow = field.NewInt64(table, "always_show")
	m.KeepAlive = field.NewInt64(table, "keep_alive")
	m.Target = field.NewString(table, "target")
	m.Title = field.NewString(table, "title")
	m.ActiveMenu = field.NewString(table, "active_menu")
	m.Breadcrumb = field.NewInt64(table, "breadcrumb")
	m.Affix = field.NewInt64(table, "affix")
	m.FrameSrc = field.NewString(table, "frame_src")
	m.FrameLoading = field.NewInt64(table, "frame_loading")
	m.Transition = field.NewString(table, "transition")
	m.Remark = field.NewString(table, "remark")
	m.CreatedBy = field.NewInt64(table, "created_by")
	m.UpdatedBy = field.NewInt64(table, "updated_by")
	m.CreatedAt = field.NewTime(table, "created_at")
	m.UpdatedAt = field.NewTime(table, "updated_at")
	m.DeletedAt = field.NewField(table, "deleted_at")
	m.Version = field.NewInt64(table, "version")

	m.fillFieldMap()

	return m
}

func (m *menu) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := m.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (m *menu) fillFieldMap() {
	m.fieldMap = make(map[string]field.Expr, 32)
	m.fieldMap["id"] = m.ID
	m.fieldMap["parent_id"] = m.ParentID
	m.fieldMap["name"] = m.Name
	m.fieldMap["permission"] = m.Permission
	m.fieldMap["type"] = m.Type
	m.fieldMap["path"] = m.Path
	m.fieldMap["component"] = m.Component
	m.fieldMap["redirect"] = m.Redirect
	m.fieldMap["icon"] = m.Icon
	m.fieldMap["sort"] = m.Sort
	m.fieldMap["hidden"] = m.Hidden
	m.fieldMap["status"] = m.Status
	m.fieldMap["always_show"] = m.AlwaysShow
	m.fieldMap["keep_alive"] = m.KeepAlive
	m.fieldMap["target"] = m.Target
	m.fieldMap["title"] = m.Title
	m.fieldMap["active_menu"] = m.ActiveMenu
	m.fieldMap["breadcrumb"] = m.Breadcrumb
	m.fieldMap["affix"] = m.Affix
	m.fieldMap["frame_src"] = m.FrameSrc
	m.fieldMap["frame_loading"] = m.FrameLoading
	m.fieldMap["transition"] = m.Transition
	m.fieldMap["remark"] = m.Remark
	m.fieldMap["created_by"] = m.CreatedBy
	m.fieldMap["updated_by"] = m.UpdatedBy
	m.fieldMap["created_at"] = m.CreatedAt
	m.fieldMap["updated_at"] = m.UpdatedAt
	m.fieldMap["deleted_at"] = m.DeletedAt
	m.fieldMap["version"] = m.Version

}

func (m menu) clone(db *gorm.DB) menu {
	m.menuDo.ReplaceConnPool(db.Statement.ConnPool)
	return m
}

func (m menu) replaceDB(db *gorm.DB) menu {
	m.menuDo.ReplaceDB(db)
	return m
}

type menuDo struct{ gen.DO }

func (m menuDo) Debug() *menuDo {
	return m.withDO(m.DO.Debug())
}

func (m menuDo) WithContext(ctx context.Context) *menuDo {
	return m.withDO(m.DO.WithContext(ctx))
}

func (m menuDo) ReadDB() *menuDo {
	return m.Clauses(dbresolver.Read)
}

func (m menuDo) WriteDB() *menuDo {
	return m.Clauses(dbresolver.Write)
}

func (m menuDo) Session(config *gorm.Session) *menuDo {
	return m.withDO(m.DO.Session(config))
}

func (m menuDo) Clauses(conds ...clause.Expression) *menuDo {
	return m.withDO(m.DO.Clauses(conds...))
}

func (m menuDo) Returning(value interface{}, columns ...string) *menuDo {
	return m.withDO(m.DO.Returning(value, columns...))
}

func (m menuDo) Not(conds ...gen.Condition) *menuDo {
	return m.withDO(m.DO.Not(conds...))
}

func (m menuDo) Or(conds ...gen.Condition) *menuDo {
	return m.withDO(m.DO.Or(conds...))
}

func (m menuDo) Select(conds ...field.Expr) *menuDo {
	return m.withDO(m.DO.Select(conds...))
}

func (m menuDo) Where(conds ...gen.Condition) *menuDo {
	return m.withDO(m.DO.Where(conds...))
}

func (m menuDo) Order(conds ...field.Expr) *menuDo {
	return m.withDO(m.DO.Order(conds...))
}

func (m menuDo) Distinct(cols ...field.Expr) *menuDo {
	return m.withDO(m.DO.Distinct(cols...))
}

func (m menuDo) Omit(cols ...field.Expr) *menuDo {
	return m.withDO(m.DO.Omit(cols...))
}

func (m menuDo) Join(table schema.Tabler, on ...field.Expr) *menuDo {
	return m.withDO(m.DO.Join(table, on...))
}

func (m menuDo) LeftJoin(table schema.Tabler, on ...field.Expr) *menuDo {
	return m.withDO(m.DO.LeftJoin(table, on...))
}

func (m menuDo) RightJoin(table schema.Tabler, on ...field.Expr) *menuDo {
	return m.withDO(m.DO.RightJoin(table, on...))
}

func (m menuDo) Group(cols ...field.Expr) *menuDo {
	return m.withDO(m.DO.Group(cols...))
}

func (m menuDo) Having(conds ...gen.Condition) *menuDo {
	return m.withDO(m.DO.Having(conds...))
}

func (m menuDo) Limit(limit int) *menuDo {
	return m.withDO(m.DO.Limit(limit))
}

func (m menuDo) Offset(offset int) *menuDo {
	return m.withDO(m.DO.Offset(offset))
}

func (m menuDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *menuDo {
	return m.withDO(m.DO.Scopes(funcs...))
}

func (m menuDo) Unscoped() *menuDo {
	return m.withDO(m.DO.Unscoped())
}

func (m menuDo) Create(values ...*entity.Menu) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Create(values)
}

func (m menuDo) CreateInBatches(values []*entity.Menu, batchSize int) error {
	return m.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (m menuDo) Save(values ...*entity.Menu) error {
	if len(values) == 0 {
		return nil
	}
	return m.DO.Save(values)
}

func (m menuDo) First() (*entity.Menu, error) {
	if result, err := m.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Menu), nil
	}
}

func (m menuDo) Take() (*entity.Menu, error) {
	if result, err := m.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Menu), nil
	}
}

func (m menuDo) Last() (*entity.Menu, error) {
	if result, err := m.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Menu), nil
	}
}

func (m menuDo) Find() ([]*entity.Menu, error) {
	result, err := m.DO.Find()
	return result.([]*entity.Menu), err
}

func (m menuDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.Menu, err error) {
	buf := make([]*entity.Menu, 0, batchSize)
	err = m.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (m menuDo) FindInBatches(result *[]*entity.Menu, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return m.DO.FindInBatches(result, batchSize, fc)
}

func (m menuDo) Attrs(attrs ...field.AssignExpr) *menuDo {
	return m.withDO(m.DO.Attrs(attrs...))
}

func (m menuDo) Assign(attrs ...field.AssignExpr) *menuDo {
	return m.withDO(m.DO.Assign(attrs...))
}

func (m menuDo) Joins(fields ...field.RelationField) *menuDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Joins(_f))
	}
	return &m
}

func (m menuDo) Preload(fields ...field.RelationField) *menuDo {
	for _, _f := range fields {
		m = *m.withDO(m.DO.Preload(_f))
	}
	return &m
}

func (m menuDo) FirstOrInit() (*entity.Menu, error) {
	if result, err := m.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Menu), nil
	}
}

func (m menuDo) FirstOrCreate() (*entity.Menu, error) {
	if result, err := m.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.Menu), nil
	}
}

func (m menuDo) FindByPage(offset int, limit int) (result []*entity.Menu, count int64, err error) {
	result, err = m.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = m.Offset(-1).Limit(-1).Count()
	return
}

func (m menuDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = m.Count()
	if err != nil {
		return
	}

	err = m.Offset(offset).Limit(limit).Scan(result)
	return
}

func (m menuDo) Scan(result interface{}) (err error) {
	return m.DO.Scan(result)
}

func (m menuDo) Delete(models ...*entity.Menu) (result gen.ResultInfo, err error) {
	return m.DO.Delete(models)
}

func (m *menuDo) withDO(do gen.Dao) *menuDo {
	m.DO = *do.(*gen.DO)
	return m
}
