package scripts

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

	// 设置命名策略
	namingStrategy := schema.NamingStrategy{
		SingularTable: config.SingularTable, // 使用单数表名
	}

	if config.TablePrefix != "" {
		namingStrategy.TablePrefix = config.TablePrefix // 设置表前缀
	}

	// 创建GORM配置
	gormConfig := &gorm.Config{
		SkipDefaultTransaction: true, // 跳过默认事务
		NamingStrategy:         namingStrategy,
	}

	// 连接数据库
	db, err := gorm.Open(mysql.Open(config.DSN), gormConfig)
	if err != nil {
		return nil, fmt.Errorf("failed to connect database: %v", err)
	}

	// 初始化生成器
	g := gen.NewGenerator(gen.Config{
		OutPath:           config.OutPath,
		ModelPkgPath:      config.ModelPkgPath,
		Mode:              gen.WithDefaultQuery | gen.WithoutContext | gen.WithQueryInterface,
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

	// 解析表名 - 根据SQL文件中的实际表名
	userTable := "sw_sys_user"
	roleTable := "sw_sys_role"
	menuTable := "sw_sys_menu"
	menuConfigTable := "sw_sys_menu_config"
	roleMenuTable := "sw_sys_role_menu"
	apiGroupTable := "sw_sys_api_group"
	apiTable := "sw_sys_api"
	roleApiTable := "sw_sys_role_api"
	deptTable := "sw_sys_dept"
	postTable := "sw_sys_post"
	loginLogTable := "sw_sys_login_log"
	operationLogTable := "sw_sys_operation_log"
	// 文件管理相关表
	fileTable := "sw_sys_file"
	fileChunkTable := "sw_sys_file_chunk"
	ossConfigTable := "sw_sys_oss_config"
	fileRefTable := "sw_sys_file_ref"

	// 生成带关联的模型
	user := g.Gen.GenerateModel(userTable)
	role := g.Gen.GenerateModel(roleTable)
	menu := g.Gen.GenerateModel(menuTable)
	apiGroup := g.Gen.GenerateModel(apiGroupTable)
	api := g.Gen.GenerateModel(apiTable)
	dept := g.Gen.GenerateModel(deptTable)
	post := g.Gen.GenerateModel(postTable)
	// 文件管理相关模型
	file := g.Gen.GenerateModel(fileTable)
	fileChunk := g.Gen.GenerateModel(fileChunkTable)
	ossConfig := g.Gen.GenerateModel(ossConfigTable)
	fileRef := g.Gen.GenerateModel(fileRefTable)

	// 设置系统用户表的关联
	userOpts := []gen.ModelOpt{
		softDeleteField,
		// 用户与角色的多对一关系
		gen.FieldRelate(field.BelongsTo, "Role", role, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"RoleID"},
				"references": {"ID"},
			},
			JSONTag: "role",
		}),
		// 用户与部门的多对一关系
		gen.FieldRelate(field.BelongsTo, "Dept", dept, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"DeptID"},
				"references": {"ID"},
			},
			JSONTag: "dept",
		}),
		// 用户与岗位的多对一关系
		gen.FieldRelate(field.BelongsTo, "Post", post, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"PostID"},
				"references": {"ID"},
			},
			JSONTag: "post",
		}),
	}

	// 设置角色表的关联
	roleOpts := []gen.ModelOpt{
		softDeleteField,
	}

	// 设置菜单表的关联
	menuOpts := []gen.ModelOpt{
		softDeleteField,
		// 菜单与父菜单的自引用关系
		gen.FieldRelate(field.BelongsTo, "Parent", menu, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"ParentID"},
				"references": {"ID"},
			},
			JSONTag: "parent",
		}),
		// 菜单与子菜单的关系
		gen.FieldRelate(field.HasMany, "Children", menu, &field.RelateConfig{
			RelateSlicePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"ParentID"},
				"references": {"ID"},
			},
			JSONTag: "children",
		}),
	}

	// 设置角色菜单关联表的选项
	roleMenuOpts := []gen.ModelOpt{
		// 角色菜单关联表与角色的多对一关系
		gen.FieldRelate(field.BelongsTo, "Role", role, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"RoleID"},
				"references": {"ID"},
			},
			JSONTag: "role",
		}),
		// 角色菜单关联表与菜单的多对一关系
		gen.FieldRelate(field.BelongsTo, "Menu", menu, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"MenuID"},
				"references": {"ID"},
			},
			JSONTag: "menu",
		}),
	}

	// 设置菜单配置表的选项
	menuConfigOpts := []gen.ModelOpt{
		// 菜单配置与菜单的一对一关系
		gen.FieldRelate(field.BelongsTo, "Menu", menu, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"MenuID"},
				"references": {"ID"},
			},
			JSONTag: "menu",
		}),
	}

	// 设置API分组表的选项
	apiGroupOpts := []gen.ModelOpt{
		softDeleteField,
	}

	// 设置API表的选项
	apiOpts := []gen.ModelOpt{
		softDeleteField,
		// API与API分组的多对一关系
		gen.FieldRelate(field.BelongsTo, "ApiGroup", apiGroup, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"Group"},
				"references": {"Code"},
			},
			JSONTag: "api_group",
		}),
	}

	// 设置角色API关联表的选项
	roleApiOpts := []gen.ModelOpt{
		// 角色API关联表与角色的多对一关系
		gen.FieldRelate(field.BelongsTo, "Role", role, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"RoleID"},
				"references": {"ID"},
			},
			JSONTag: "role",
		}),
		// 角色API关联表与API的多对一关系
		gen.FieldRelate(field.BelongsTo, "Api", api, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"ApiID"},
				"references": {"ID"},
			},
			JSONTag: "api",
		}),
	}

	// 设置登录日志表的选项
	loginLogOpts := []gen.ModelOpt{
		// 登录日志与用户的多对一关系
		gen.FieldRelate(field.BelongsTo, "User", user, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"UserID"},
				"references": {"ID"},
			},
			JSONTag: "user",
		}),
	}

	// 设置部门表的选项
	deptOpts := []gen.ModelOpt{
		softDeleteField,
		// 部门与父部门的自引用关系
		gen.FieldRelate(field.BelongsTo, "Parent", dept, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"ParentID"},
				"references": {"ID"},
			},
			JSONTag: "parent",
		}),
		// 部门与子部门的关系
		gen.FieldRelate(field.HasMany, "Children", dept, &field.RelateConfig{
			RelateSlicePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"ParentID"},
				"references": {"ID"},
			},
			JSONTag: "children",
		}),
	}

	// 设置岗位表的选项
	postOpts := []gen.ModelOpt{
		softDeleteField,
		// 岗位与部门的多对一关系
		gen.FieldRelate(field.BelongsTo, "Dept", dept, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"DeptID"},
				"references": {"ID"},
			},
			JSONTag: "dept",
		}),
	}

	// 设置操作日志表的选项
	operationLogOpts := []gen.ModelOpt{
		// 操作日志与用户的多对一关系
		gen.FieldRelate(field.BelongsTo, "User", user, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"UserID"},
				"references": {"ID"},
			},
			JSONTag: "user",
		}),
	}

	// 设置文件表的选项
	fileOpts := []gen.ModelOpt{
		softDeleteField,
		// 文件与上传用户的多对一关系
		gen.FieldRelate(field.BelongsTo, "UploadUser", user, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"UploadUserID"},
				"references": {"ID"},
			},
			JSONTag: "upload_user",
		}),
		// 文件与文件引用的一对多关系
		gen.FieldRelate(field.HasMany, "FileRefs", fileRef, &field.RelateConfig{
			GORMTag: map[string][]string{
				"foreignKey": {"FileID"},
				"references": {"ID"},
			},
			JSONTag: "file_refs",
		}),
	}

	// 设置文件分片表的选项
	fileChunkOpts := []gen.ModelOpt{
		// 文件分片与上传用户的多对一关系
		gen.FieldRelate(field.BelongsTo, "UploadUser", user, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"UploadUserID"},
				"references": {"ID"},
			},
			JSONTag: "upload_user",
		}),
	}

	// 设置OSS配置表的选项
	ossConfigOpts := []gen.ModelOpt{
		softDeleteField,
	}

	// 设置文件引用表的选项
	fileRefOpts := []gen.ModelOpt{
		// 文件引用与文件的多对一关系
		gen.FieldRelate(field.BelongsTo, "File", file, &field.RelateConfig{
			RelatePointer: true,
			GORMTag: map[string][]string{
				"foreignKey": {"FileID"},
				"references": {"ID"},
			},
			JSONTag: "file",
		}),
	}

	// 重新生成带关联的模型
	user = g.Gen.GenerateModel(userTable, userOpts...)
	role = g.Gen.GenerateModel(roleTable, roleOpts...)
	menu = g.Gen.GenerateModel(menuTable, menuOpts...)
	menuConfig := g.Gen.GenerateModel(menuConfigTable, menuConfigOpts...)
	roleMenu := g.Gen.GenerateModel(roleMenuTable, roleMenuOpts...)
	apiGroup = g.Gen.GenerateModel(apiGroupTable, apiGroupOpts...)
	api = g.Gen.GenerateModel(apiTable, apiOpts...)
	roleApi := g.Gen.GenerateModel(roleApiTable, roleApiOpts...)
	dept = g.Gen.GenerateModel(deptTable, deptOpts...)
	post = g.Gen.GenerateModel(postTable, postOpts...)
	loginLog := g.Gen.GenerateModel(loginLogTable, loginLogOpts...)
	operationLog := g.Gen.GenerateModel(operationLogTable, operationLogOpts...)
	// 文件管理相关模型
	file = g.Gen.GenerateModel(fileTable, fileOpts...)
	fileChunk = g.Gen.GenerateModel(fileChunkTable, fileChunkOpts...)
	ossConfig = g.Gen.GenerateModel(ossConfigTable, ossConfigOpts...)
	fileRef = g.Gen.GenerateModel(fileRefTable, fileRefOpts...)

	// 应用基本模型
	g.Gen.ApplyBasic(user, role, menu, menuConfig, roleMenu, apiGroup, api, roleApi, dept, post, loginLog, operationLog, file, fileChunk, ossConfig, fileRef)
}

// GenerateModelsWithRelations 生成带关联关系的模型
func (g *Generator) GenerateModelsWithRelations() {
	// 设置模型关联关系
	g.SetupModelRelations()

	// 执行代码生成
	g.Execute()
}

// GenerateSystemModels 生成系统模块的所有模型（便捷方法）
func GenerateSystemModels(dsn, outPath, modelPkgPath string) error {
	// 创建生成器配置
	config := &GenConfig{
		DSN:             dsn,
		OutPath:         outPath,
		ModelPkgPath:    modelPkgPath,
		WithUnitTest:    false,
		WithQueryFilter: true,
		TablePrefix:     "sw_",
		SingularTable:   true,
	}

	// 创建生成器实例
	generator, err := NewGenerator(config)
	if err != nil {
		return fmt.Errorf("failed to create generator: %v", err)
	}

	// 生成带关联关系的模型
	generator.GenerateModelsWithRelations()

	return nil
}

// GenerateAllSystemTables 生成系统所有表的模型（不包含关联关系）
func GenerateAllSystemTables(dsn, outPath, modelPkgPath string) error {
	// 创建生成器配置
	config := &GenConfig{
		DSN:             dsn,
		OutPath:         outPath,
		ModelPkgPath:    modelPkgPath,
		WithUnitTest:    false,
		WithQueryFilter: true,
		TablePrefix:     "sw_",
		SingularTable:   true,
	}

	// 创建生成器实例
	generator, err := NewGenerator(config)
	if err != nil {
		return fmt.Errorf("failed to create generator: %v", err)
	}

	// 生成所有表的模型
	err = generator.GenerateAllModel()
	if err != nil {
		return fmt.Errorf("failed to generate models: %v", err)
	}

	// 执行代码生成
	generator.Execute()

	return nil
}
