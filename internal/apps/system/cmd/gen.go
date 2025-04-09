package cmd

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

type GenConfig struct {
	DSN             string // 数据库连接串
	OutPath         string // 输出路径
	ModelPkgPath    string // 模型包路径
	WithUnitTest    bool   // 是否生成单元测试
	WithQueryFilter bool   // 是否生成查询过滤器
	TablePrefix     string // 表前缀
	SingularTable   bool   // 是否使用单数表名，默认为true
}

type Gen struct {
	cfg *GenConfig
	DB  *gorm.DB
	Gen *gen.Generator
}

func InitGen(cfg *GenConfig) (*Gen, error) {
	if cfg == nil {
		return nil, fmt.Errorf("cfg is nil")
	}

	gormCfg := &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   cfg.TablePrefix,   // 设置表前缀
			SingularTable: cfg.SingularTable, // 使用单数表名
		},
	}

	// 设置命名策略
	namingStrategy := schema.NamingStrategy{
		SingularTable: cfg.SingularTable, // 使用单数表名
	}

	if cfg.TablePrefix != "" {
		namingStrategy.TablePrefix = cfg.TablePrefix // 设置表前缀
	}

	gormCfg.NamingStrategy = namingStrategy

	db, err := gorm.Open(mysql.Open(cfg.DSN), gormCfg)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	g := gen.NewGenerator(gen.Config{
		OutPath:           cfg.OutPath,
		ModelPkgPath:      cfg.ModelPkgPath,
		Mode:              gen.WithDefaultQuery | gen.WithoutContext,
		FieldNullable:     true,  // 可空字段使用指针类型
		FieldCoverable:    true,  // 有默认值的字段使用指针类型
		FieldSignable:     false, // 数字类型保持与数据库一致
		FieldWithIndexTag: false, // 不生成索引标签
		FieldWithTypeTag:  true,  // 生成带类型标签的字段
	})

	toCamelCase := func(s string) string {
		// 将下划线分隔的单词转换为大驼峰命名
		parts := strings.Split(s, "_")
		for i, part := range parts {
			if len(part) > 0 {
				parts[i] = strings.ToUpper(part[:1]) + part[1:]
			}
		}
		return strings.Join(parts, "")
	}

	// 应用模型命名策略
	g.WithModelNameStrategy(func(tableName string) string {
		// 移除表前缀
		name := tableName
		if cfg.TablePrefix != "" && strings.HasPrefix(tableName, cfg.TablePrefix) {
			name = strings.TrimPrefix(tableName, cfg.TablePrefix)
		}
		// 转换为大驼峰命名
		return toCamelCase(name)
	})

	// 使用数据库连接
	g.UseDB(db)

	// 自定义字段的数据类型映射
	dataMap := map[string]func(columnType gorm.ColumnType) (dataType string){
		"tinyint":   func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"smallint":  func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"mediumint": func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"bigint":    func(columnType gorm.ColumnType) (dataType string) { return "int64" },
		"int":       func(columnType gorm.ColumnType) (dataType string) { return "int64" },
	}
	g.WithDataTypeMap(dataMap)
	return &Gen{
		cfg: cfg,
		DB:  db,
		Gen: g,
	}, nil
}

func (g *Gen) SetupModelRelations() {
	// 定义软删除字段类型
	softDeleteField := gen.FieldType("deleted_at", "gorm.DeletedAt")

	// 定义表名
	userTable := "sys_user"
	roleTable := "sys_role"
	deptTable := "sys_dept"
	postTable := "sys_post"
	menuTable := "sys_menu"
	roleMenuTable := "sys_role_menu"
	dictTypeTable := "sys_dict_type"
	dictDataTable := "sys_dict_data"
	loginLogTable := "sys_login_log"
	operationLogTable := "sys_operation_log"

	// 用户表关联配置
	userOpts := []gen.ModelOpt{
		softDeleteField,
		gen.FieldNew("Role", "*Role", field.Tag{
			"json": "role",
			"gorm": "foreignKey:RoleID;references:ID",
		}),
		gen.FieldNew("Dept", "*Dept", field.Tag{
			"json": "dept",
			"gorm": "foreignKey:DeptID;references:ID",
		}),
		gen.FieldNew("Post", "*Post", field.Tag{
			"json": "post",
			"gorm": "foreignKey:PostID;references:ID",
		}),
	}

	// 部门表关联配置
	deptOpts := []gen.ModelOpt{
		softDeleteField,
		gen.FieldNew("Parent", "*Dept", field.Tag{
			"json": "parent",
			"gorm": "foreignKey:Pid;references:ID",
		}),
		gen.FieldNew("Children", "[]*Dept", field.Tag{
			"json": "children",
			"gorm": "foreignKey:Pid;references:ID",
		}),
		gen.FieldNew("Users", "[]*User", field.Tag{
			"json": "users",
			"gorm": "foreignKey:DeptID;references:ID",
		}),
	}

	// 岗位表关联配置
	postOpts := []gen.ModelOpt{
		softDeleteField,
		gen.FieldNew("Users", "[]*User", field.Tag{
			"json": "users",
			"gorm": "foreignKey:PostID;references:ID",
		}),
	}

	// 角色表关联配置
	roleOpts := []gen.ModelOpt{
		softDeleteField,
		gen.FieldNew("Users", "[]*User", field.Tag{
			"json": "users",
			"gorm": "foreignKey:RoleID;references:ID",
		}),
		gen.FieldNew("Menus", "[]*Menu", field.Tag{
			"json": "menus",
			"gorm": "many2many:sys_role_menu;foreignKey:ID;joinForeignKey:RoleID;References:ID;JoinReferences:MenuID",
		}),
	}

	// 菜单表关联配置
	menuOpts := []gen.ModelOpt{
		softDeleteField,
		gen.FieldNew("Parent", "*Menu", field.Tag{
			"json": "parent",
			"gorm": "foreignKey:ParentID;references:ID",
		}),
		gen.FieldNew("Children", "[]*Menu", field.Tag{
			"json": "children",
			"gorm": "foreignKey:ParentID;references:ID",
		}),
		gen.FieldNew("Roles", "[]*Role", field.Tag{
			"json": "roles",
			"gorm": "many2many:sys_role_menu;foreignKey:ID;joinForeignKey:MenuID;References:ID;JoinReferences:RoleID",
		}),
	}

	// 角色菜单关联表配置
	roleMenuOpts := []gen.ModelOpt{
		gen.FieldNew("Role", "*Role", field.Tag{
			"json": "role",
			"gorm": "foreignKey:RoleID;references:ID",
		}),
		gen.FieldNew("Menu", "*Menu", field.Tag{
			"json": "menu",
			"gorm": "foreignKey:MenuID;references:ID",
		}),
	}

	// 字典类型表关联配置
	dictTypeOpts := []gen.ModelOpt{
		softDeleteField,
		gen.FieldNew("DictData", "[]*DictData", field.Tag{
			"json": "dictData",
			"gorm": "foreignKey:DictType;references:Type",
		}),
	}

	// 字典数据表关联配置
	dictDataOpts := []gen.ModelOpt{
		softDeleteField,
		gen.FieldNew("Type", "*DictType", field.Tag{
			"json": "dictType",
			"gorm": "foreignKey:DictType;references:Type",
		}),
	}

	// 登录日志表配置
	var loginLogOpts []gen.ModelOpt

	// 操作日志表配置
	var operationLogOpts []gen.ModelOpt

	// 生成模型
	user := g.Gen.GenerateModel(userTable, userOpts...)
	dept := g.Gen.GenerateModel(deptTable, deptOpts...)
	post := g.Gen.GenerateModel(postTable, postOpts...)
	role := g.Gen.GenerateModel(roleTable, roleOpts...)
	menu := g.Gen.GenerateModel(menuTable, menuOpts...)
	roleMenu := g.Gen.GenerateModel(roleMenuTable, roleMenuOpts...)
	dictType := g.Gen.GenerateModel(dictTypeTable, dictTypeOpts...)
	dictData := g.Gen.GenerateModel(dictDataTable, dictDataOpts...)
	loginLog := g.Gen.GenerateModel(loginLogTable, loginLogOpts...)
	operationLog := g.Gen.GenerateModel(operationLogTable, operationLogOpts...)

	// 应用所有模型
	g.Gen.ApplyBasic(user, dept, post, role, menu, roleMenu, dictType, dictData, loginLog, operationLog)

	g.Gen.Execute()
}
