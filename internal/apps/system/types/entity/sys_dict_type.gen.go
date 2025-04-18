// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.
// Code generated by gorm.io/gen. DO NOT EDIT.

package entity

import (
	"time"

	"gorm.io/gorm"
)

const TableNameDictType = "sys_dict_type"

// DictType 字典类型表
type DictType struct {
	ID        int64          `gorm:"column:id;type:bigint;primaryKey;autoIncrement:true;comment:字典类型ID" json:"id"` // 字典类型ID
	Name      string         `gorm:"column:name;type:varchar(64);not null;comment:字典名称" json:"name"`               // 字典名称
	Type      string         `gorm:"column:type;type:varchar(64);not null;comment:字典类型" json:"type"`               // 字典类型
	Status    *int64         `gorm:"column:status;type:tinyint(1);default:1;comment:状态 1-正常 2-禁用" json:"status"`   // 状态 1-正常 2-禁用
	Remark    *string        `gorm:"column:remark;type:varchar(255);comment:备注" json:"remark"`                     // 备注
	CreatedBy *int64         `gorm:"column:created_by;type:bigint;comment:创建人" json:"created_by"`                  // 创建人
	UpdatedBy *int64         `gorm:"column:updated_by;type:bigint;comment:更新人" json:"updated_by"`                  // 更新人
	CreatedAt time.Time      `gorm:"column:created_at;type:datetime;not null;comment:创建时间" json:"created_at"`      // 创建时间
	UpdatedAt *time.Time     `gorm:"column:updated_at;type:datetime;comment:更新时间" json:"updated_at"`               // 更新时间
	DeletedAt gorm.DeletedAt `gorm:"column:deleted_at;type:datetime;comment:删除时间" json:"deleted_at"`               // 删除时间
	Version   *int64         `gorm:"column:version;type:int;comment:乐观锁版本" json:"version"`                         // 乐观锁版本
	DictData  []*DictData    `gorm:"foreignKey:DictType;references:Type" json:"dictData"`
}

// TableName DictType's table name
func (*DictType) TableName() string {
	return TableNameDictType
}
