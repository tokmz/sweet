package script

import (
	"fmt"
	"strings"

	"gorm.io/driver/mysql"
	"gorm.io/gen"
	"gorm.io/gen/field"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
)

// GenConfig 代码生成器配置
type GenConfig struct {
	DSN             string // 数据库连接串
	OutPath         string // 输出路径
	ModelPkgPath    string // 模型包路径
	WithUnitTest    bool   // 是否生成单元测试
	WithQueryFilter bool   // 是否生成查询过滤器
	TablePrefix     string // 表前缀
	SingularTable   bool   // 是否使用单数表名，默认为true
}

// Generator GORM 生成器
type Generator struct {
	Config *GenConfig
	DB     *gorm.DB
	Gen    *gen.Generator
}

// NewGenerator 创建一个新的生成器实例
func NewGenerator(config *GenConfig) (*Generator, error) {
	if config == nil {
		return nil, fmt.Errorf("config cannot be nil")
	}

	// 设置默认值
	if !config.SingularTable {
		config.SingularTable = true // 默认使用单数表名
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务
		NamingStrategy: schema.NamingStrategy{
			TablePrefix:   config.TablePrefix,   // 设置表前缀
			SingularTable: config.SingularTable, // 使用单数表名
		},
	}

	// 设置命名策略
	namingStrategy := schema.NamingStrategy{
		SingularTable: config.SingularTable, // 使用单数表名
	}

	if config.TablePrefix != "" {
		namingStrategy.TablePrefix = config.TablePrefix // 设置表前缀
	}

	gormConfig.NamingStrategy = namingStrategy

	// 连接数据库
	db, err := gorm.Open(mysql.Open(config.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 初始化生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           config.OutPath,
		ModelPkgPath:      config.ModelPkgPath,
		Mode:              gen.WithDefaultQuery | gen.WithoutContext,
		FieldNullable:     true,  // 可空字段使用指针类型
		FieldCoverable:    true,  // 有默认值的字段使用指针类型
		FieldSignable:     false, // 数字类型保持与数据库一致
		FieldWithIndexTag: false, // 不生成索引标签
		FieldWithTypeTag:  true,  // 生成带类型标签的字段
	})

	// 设置去除表前缀的模型命名策略
	// 定义一个转换为大驼峰命名的函数
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
		if config.TablePrefix != "" && strings.HasPrefix(tableName, config.TablePrefix) {
			name = strings.TrimPrefix(tableName, config.TablePrefix)
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

	return &Generator{
		Config: config,
		DB:     db,
		Gen:    g,
	}, nil
}

// GenerateModel 生成指定表的模型
func (g *Generator) GenerateModel(tableName string, opts ...gen.ModelOpt) error {
	// 生成基本模型文件
	g.Gen.GenerateModel(tableName, opts...)
	return nil
}

// GenerateAllModel 生成数据库中所有表的模型
func (g *Generator) GenerateAllModel() error {
	// 生成所有模型文件
	g.Gen.GenerateAllTable()
	return nil
}

// Execute 执行代码生成
func (g *Generator) Execute() {
	g.Gen.Execute()
}

// SetupModelRelations 设置模型关联关系
func (g *Generator) SetupModelRelations() {
	// 定义软删除字段类型
	softDeleteField := gen.FieldType("deleted_at", "gorm.DeletedAt")

	// 解析表名 - 不需要手动添加前缀，GORM会根据命名策略自动处理
	userTable := "sys_user"
	departmentTable := "sys_department"
	positionTable := "sys_position"
	roleTable := "sys_role"
	menuTable := "sys_menu"
	roleMenuTable := "sys_role_menu"
	userRoleTable := "sys_user_role"

	// 设置用户表的关联
	userOpts := []gen.ModelOpt{
		softDeleteField,
		// 用户与部门的一对一关系 - 使用指针类型
		gen.FieldNew("Department", "*Department", field.Tag{
			"json": "department",
			"gorm": "foreignKey:DeptID;references:ID",
		}),
		// 用户与岗位的一对一关系 - 使用指针类型
		gen.FieldNew("Position", "*Position", field.Tag{
			"json": "position",
			"gorm": "foreignKey:PositionID;references:ID",
		}),
		// 用户与角色的多对多关系
		gen.FieldNew("Roles", "[]*Role", field.Tag{
			"json": "roles",
			"gorm": "many2many:sys_user_role;foreignKey:ID;joinForeignKey:UserID;References:ID;JoinReferences:RoleID",
		}),
	}

	// 设置部门表的关联
	deptOpts := []gen.ModelOpt{
		softDeleteField,
		// 部门与父部门的自引用关系 - 使用指针类型
		gen.FieldNew("Parent", "*Department", field.Tag{
			"json": "parent",
			"gorm": "foreignKey:Pid;references:ID",
		}),
		// 部门与子部门的关系 - 使用指针数组类型
		gen.FieldNew("Children", "[]*Department", field.Tag{
			"json": "children",
			"gorm": "foreignKey:Pid;references:ID",
		}),
	}

	// 设置岗位表的关联
	posOpts := []gen.ModelOpt{
		softDeleteField,
		// 岗位与部门的多对一关系 - 使用指针类型
		gen.FieldNew("Department", "*Department", field.Tag{
			"json": "department",
			"gorm": "foreignKey:DeptID;references:ID",
		}),
	}

	// 设置角色表的关联
	roleOpts := []gen.ModelOpt{
		softDeleteField,
		// 角色与菜单的多对多关系
		gen.FieldNew("Menus", "[]*Menu", field.Tag{
			"json": "menus",
			"gorm": "many2many:sys_role_menu;foreignKey:ID;joinForeignKey:RoleID;References:ID;JoinReferences:MenuID",
		}),
	}

	// 设置菜单表的关联
	menuOpts := []gen.ModelOpt{
		softDeleteField,
		// 菜单与父菜单的自引用关系
		gen.FieldNew("Parent", "*Menu", field.Tag{
			"json": "parent",
			"gorm": "foreignKey:Pid;references:ID",
		}),
		// 菜单与子菜单的关系
		gen.FieldNew("Children", "[]*Menu", field.Tag{
			"json": "children",
			"gorm": "foreignKey:Pid;references:ID",
		}),
	}

	// 设置角色菜单关联表的选项
	var roleMenuOpts []gen.ModelOpt

	// 设置用户角色关联表的选项
	var userRoleOpts []gen.ModelOpt

	// 生成带关联的模型
	user := g.Gen.GenerateModel(userTable, userOpts...)
	dept := g.Gen.GenerateModel(departmentTable, deptOpts...)
	pos := g.Gen.GenerateModel(positionTable, posOpts...)
	role := g.Gen.GenerateModel(roleTable, roleOpts...)
	menu := g.Gen.GenerateModel(menuTable, menuOpts...)
	roleMenu := g.Gen.GenerateModel(roleMenuTable, roleMenuOpts...)
	userRole := g.Gen.GenerateModel(userRoleTable, userRoleOpts...)

	// 应用基本模型
	g.Gen.ApplyBasic(user, dept, pos, role, menu, roleMenu, userRole)
}

// GenerateModelsWithRelations 生成带关联关系的模型
func (g *Generator) GenerateModelsWithRelations() {
	// 设置模型关联关系
	g.SetupModelRelations()

	// 执行代码生成
	g.Execute()
}
