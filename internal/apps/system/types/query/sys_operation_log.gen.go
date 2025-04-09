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

func newOperationLog(db *gorm.DB, opts ...gen.DOOption) operationLog {
	_operationLog := operationLog{}

	_operationLog.operationLogDo.UseDB(db, opts...)
	_operationLog.operationLogDo.UseModel(&entity.OperationLog{})

	tableName := _operationLog.operationLogDo.TableName()
	_operationLog.ALL = field.NewAsterisk(tableName)
	_operationLog.ID = field.NewInt64(tableName, "id")
	_operationLog.Title = field.NewString(tableName, "title")
	_operationLog.BusinessType = field.NewInt64(tableName, "business_type")
	_operationLog.Method = field.NewString(tableName, "method")
	_operationLog.RequestMethod = field.NewString(tableName, "request_method")
	_operationLog.OperatorType = field.NewInt64(tableName, "operator_type")
	_operationLog.UserID = field.NewInt64(tableName, "user_id")
	_operationLog.Username = field.NewString(tableName, "username")
	_operationLog.URL = field.NewString(tableName, "url")
	_operationLog.IP = field.NewString(tableName, "ip")
	_operationLog.Location = field.NewString(tableName, "location")
	_operationLog.Browser = field.NewString(tableName, "browser")
	_operationLog.Os = field.NewString(tableName, "os")
	_operationLog.Device = field.NewString(tableName, "device")
	_operationLog.Param = field.NewString(tableName, "param")
	_operationLog.Result = field.NewString(tableName, "result")
	_operationLog.Status = field.NewInt64(tableName, "status")
	_operationLog.ErrorStack = field.NewString(tableName, "error_stack")
	_operationLog.CostTime = field.NewInt64(tableName, "cost_time")
	_operationLog.CreatedAt = field.NewTime(tableName, "created_at")

	_operationLog.fillFieldMap()

	return _operationLog
}

// operationLog 操作日志表
type operationLog struct {
	operationLogDo

	ALL           field.Asterisk
	ID            field.Int64  // 日志ID
	Title         field.String // 模块标题
	BusinessType  field.Int64  // 业务类型（1新增 2修改 3删除 4其它）
	Method        field.String // 方法名称
	RequestMethod field.String // 请求方式
	OperatorType  field.Int64  // 操作类别（1后台用户 2手机端用户 3其它）
	UserID        field.Int64  // 操作人ID
	Username      field.String // 操作人名称
	URL           field.String // 请求URL
	IP            field.String // 操作地址
	Location      field.String // 操作地点
	Browser       field.String // 浏览器类型
	Os            field.String // 操作系统
	Device        field.String // 设备
	Param         field.String // 请求参数
	Result        field.String // 返回结果
	Status        field.Int64  // 操作状态（1正常 2异常）
	ErrorStack    field.String // 错误消息
	CostTime      field.Int64  // 消耗时间
	CreatedAt     field.Time   // 创建时间

	fieldMap map[string]field.Expr
}

func (o operationLog) Table(newTableName string) *operationLog {
	o.operationLogDo.UseTable(newTableName)
	return o.updateTableName(newTableName)
}

func (o operationLog) As(alias string) *operationLog {
	o.operationLogDo.DO = *(o.operationLogDo.As(alias).(*gen.DO))
	return o.updateTableName(alias)
}

func (o *operationLog) updateTableName(table string) *operationLog {
	o.ALL = field.NewAsterisk(table)
	o.ID = field.NewInt64(table, "id")
	o.Title = field.NewString(table, "title")
	o.BusinessType = field.NewInt64(table, "business_type")
	o.Method = field.NewString(table, "method")
	o.RequestMethod = field.NewString(table, "request_method")
	o.OperatorType = field.NewInt64(table, "operator_type")
	o.UserID = field.NewInt64(table, "user_id")
	o.Username = field.NewString(table, "username")
	o.URL = field.NewString(table, "url")
	o.IP = field.NewString(table, "ip")
	o.Location = field.NewString(table, "location")
	o.Browser = field.NewString(table, "browser")
	o.Os = field.NewString(table, "os")
	o.Device = field.NewString(table, "device")
	o.Param = field.NewString(table, "param")
	o.Result = field.NewString(table, "result")
	o.Status = field.NewInt64(table, "status")
	o.ErrorStack = field.NewString(table, "error_stack")
	o.CostTime = field.NewInt64(table, "cost_time")
	o.CreatedAt = field.NewTime(table, "created_at")

	o.fillFieldMap()

	return o
}

func (o *operationLog) GetFieldByName(fieldName string) (field.OrderExpr, bool) {
	_f, ok := o.fieldMap[fieldName]
	if !ok || _f == nil {
		return nil, false
	}
	_oe, ok := _f.(field.OrderExpr)
	return _oe, ok
}

func (o *operationLog) fillFieldMap() {
	o.fieldMap = make(map[string]field.Expr, 20)
	o.fieldMap["id"] = o.ID
	o.fieldMap["title"] = o.Title
	o.fieldMap["business_type"] = o.BusinessType
	o.fieldMap["method"] = o.Method
	o.fieldMap["request_method"] = o.RequestMethod
	o.fieldMap["operator_type"] = o.OperatorType
	o.fieldMap["user_id"] = o.UserID
	o.fieldMap["username"] = o.Username
	o.fieldMap["url"] = o.URL
	o.fieldMap["ip"] = o.IP
	o.fieldMap["location"] = o.Location
	o.fieldMap["browser"] = o.Browser
	o.fieldMap["os"] = o.Os
	o.fieldMap["device"] = o.Device
	o.fieldMap["param"] = o.Param
	o.fieldMap["result"] = o.Result
	o.fieldMap["status"] = o.Status
	o.fieldMap["error_stack"] = o.ErrorStack
	o.fieldMap["cost_time"] = o.CostTime
	o.fieldMap["created_at"] = o.CreatedAt
}

func (o operationLog) clone(db *gorm.DB) operationLog {
	o.operationLogDo.ReplaceConnPool(db.Statement.ConnPool)
	return o
}

func (o operationLog) replaceDB(db *gorm.DB) operationLog {
	o.operationLogDo.ReplaceDB(db)
	return o
}

type operationLogDo struct{ gen.DO }

func (o operationLogDo) Debug() *operationLogDo {
	return o.withDO(o.DO.Debug())
}

func (o operationLogDo) WithContext(ctx context.Context) *operationLogDo {
	return o.withDO(o.DO.WithContext(ctx))
}

func (o operationLogDo) ReadDB() *operationLogDo {
	return o.Clauses(dbresolver.Read)
}

func (o operationLogDo) WriteDB() *operationLogDo {
	return o.Clauses(dbresolver.Write)
}

func (o operationLogDo) Session(config *gorm.Session) *operationLogDo {
	return o.withDO(o.DO.Session(config))
}

func (o operationLogDo) Clauses(conds ...clause.Expression) *operationLogDo {
	return o.withDO(o.DO.Clauses(conds...))
}

func (o operationLogDo) Returning(value interface{}, columns ...string) *operationLogDo {
	return o.withDO(o.DO.Returning(value, columns...))
}

func (o operationLogDo) Not(conds ...gen.Condition) *operationLogDo {
	return o.withDO(o.DO.Not(conds...))
}

func (o operationLogDo) Or(conds ...gen.Condition) *operationLogDo {
	return o.withDO(o.DO.Or(conds...))
}

func (o operationLogDo) Select(conds ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.Select(conds...))
}

func (o operationLogDo) Where(conds ...gen.Condition) *operationLogDo {
	return o.withDO(o.DO.Where(conds...))
}

func (o operationLogDo) Order(conds ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.Order(conds...))
}

func (o operationLogDo) Distinct(cols ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.Distinct(cols...))
}

func (o operationLogDo) Omit(cols ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.Omit(cols...))
}

func (o operationLogDo) Join(table schema.Tabler, on ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.Join(table, on...))
}

func (o operationLogDo) LeftJoin(table schema.Tabler, on ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.LeftJoin(table, on...))
}

func (o operationLogDo) RightJoin(table schema.Tabler, on ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.RightJoin(table, on...))
}

func (o operationLogDo) Group(cols ...field.Expr) *operationLogDo {
	return o.withDO(o.DO.Group(cols...))
}

func (o operationLogDo) Having(conds ...gen.Condition) *operationLogDo {
	return o.withDO(o.DO.Having(conds...))
}

func (o operationLogDo) Limit(limit int) *operationLogDo {
	return o.withDO(o.DO.Limit(limit))
}

func (o operationLogDo) Offset(offset int) *operationLogDo {
	return o.withDO(o.DO.Offset(offset))
}

func (o operationLogDo) Scopes(funcs ...func(gen.Dao) gen.Dao) *operationLogDo {
	return o.withDO(o.DO.Scopes(funcs...))
}

func (o operationLogDo) Unscoped() *operationLogDo {
	return o.withDO(o.DO.Unscoped())
}

func (o operationLogDo) Create(values ...*entity.OperationLog) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Create(values)
}

func (o operationLogDo) CreateInBatches(values []*entity.OperationLog, batchSize int) error {
	return o.DO.CreateInBatches(values, batchSize)
}

// Save : !!! underlying implementation is different with GORM
// The method is equivalent to executing the statement: db.Clauses(clause.OnConflict{UpdateAll: true}).Create(values)
func (o operationLogDo) Save(values ...*entity.OperationLog) error {
	if len(values) == 0 {
		return nil
	}
	return o.DO.Save(values)
}

func (o operationLogDo) First() (*entity.OperationLog, error) {
	if result, err := o.DO.First(); err != nil {
		return nil, err
	} else {
		return result.(*entity.OperationLog), nil
	}
}

func (o operationLogDo) Take() (*entity.OperationLog, error) {
	if result, err := o.DO.Take(); err != nil {
		return nil, err
	} else {
		return result.(*entity.OperationLog), nil
	}
}

func (o operationLogDo) Last() (*entity.OperationLog, error) {
	if result, err := o.DO.Last(); err != nil {
		return nil, err
	} else {
		return result.(*entity.OperationLog), nil
	}
}

func (o operationLogDo) Find() ([]*entity.OperationLog, error) {
	result, err := o.DO.Find()
	return result.([]*entity.OperationLog), err
}

func (o operationLogDo) FindInBatch(batchSize int, fc func(tx gen.Dao, batch int) error) (results []*entity.OperationLog, err error) {
	buf := make([]*entity.OperationLog, 0, batchSize)
	err = o.DO.FindInBatches(&buf, batchSize, func(tx gen.Dao, batch int) error {
		defer func() { results = append(results, buf...) }()
		return fc(tx, batch)
	})
	return results, err
}

func (o operationLogDo) FindInBatches(result *[]*entity.OperationLog, batchSize int, fc func(tx gen.Dao, batch int) error) error {
	return o.DO.FindInBatches(result, batchSize, fc)
}

func (o operationLogDo) Attrs(attrs ...field.AssignExpr) *operationLogDo {
	return o.withDO(o.DO.Attrs(attrs...))
}

func (o operationLogDo) Assign(attrs ...field.AssignExpr) *operationLogDo {
	return o.withDO(o.DO.Assign(attrs...))
}

func (o operationLogDo) Joins(fields ...field.RelationField) *operationLogDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Joins(_f))
	}
	return &o
}

func (o operationLogDo) Preload(fields ...field.RelationField) *operationLogDo {
	for _, _f := range fields {
		o = *o.withDO(o.DO.Preload(_f))
	}
	return &o
}

func (o operationLogDo) FirstOrInit() (*entity.OperationLog, error) {
	if result, err := o.DO.FirstOrInit(); err != nil {
		return nil, err
	} else {
		return result.(*entity.OperationLog), nil
	}
}

func (o operationLogDo) FirstOrCreate() (*entity.OperationLog, error) {
	if result, err := o.DO.FirstOrCreate(); err != nil {
		return nil, err
	} else {
		return result.(*entity.OperationLog), nil
	}
}

func (o operationLogDo) FindByPage(offset int, limit int) (result []*entity.OperationLog, count int64, err error) {
	result, err = o.Offset(offset).Limit(limit).Find()
	if err != nil {
		return
	}

	if size := len(result); 0 < limit && 0 < size && size < limit {
		count = int64(size + offset)
		return
	}

	count, err = o.Offset(-1).Limit(-1).Count()
	return
}

func (o operationLogDo) ScanByPage(result interface{}, offset int, limit int) (count int64, err error) {
	count, err = o.Count()
	if err != nil {
		return
	}

	err = o.Offset(offset).Limit(limit).Scan(result)
	return
}

func (o operationLogDo) Scan(result interface{}) (err error) {
	return o.DO.Scan(result)
}

func (o operationLogDo) Delete(models ...*entity.OperationLog) (result gen.ResultInfo, err error) {
	return o.DO.Delete(models)
}

func (o *operationLogDo) withDO(do gen.Dao) *operationLogDo {
	o.DO = *do.(*gen.DO)
	return o
}
