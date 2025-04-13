package auth

import (
	"fmt"
	"strings"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	fileadapter "github.com/casbin/casbin/v2/persist/file-adapter"
	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
)

// CasbinService Casbin服务
type CasbinService struct {
	enforcer *casbin.Enforcer
	config   CasbinConfig
}

// NewCasbinService 创建Casbin服务
func NewCasbinService(config CasbinConfig) (*CasbinService, error) {
	var enforcer *casbin.Enforcer
	var err error

	// 加载模型
	var m model.Model
	if strings.HasPrefix(config.Model, "file://") {
		// 从文件加载模型
		modelPath := strings.TrimPrefix(config.Model, "file://")
		m, err = model.NewModelFromFile(modelPath)
	} else {
		// 从字符串加载模型
		m, err = model.NewModelFromString(config.Model)
	}
	if err != nil {
		return nil, fmt.Errorf("加载Casbin模型失败: %w", err)
	}

	// 创建适配器
	var adapter interface{}
	switch config.Adapter {
	case "file":
		// 文件适配器
		adapter = fileadapter.NewAdapter(config.DSN)
	case "mysql":
		// MySQL适配器
		db, err := gorm.Open(mysql.Open(config.DSN), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("连接MySQL数据库失败: %w", err)
		}
		
		// 自动迁移表结构
		if config.AutoMigrate {
			if err := db.AutoMigrate(&gormadapter.CasbinRule{}); err != nil {
				return nil, fmt.Errorf("迁移Casbin表结构失败: %w", err)
			}
		}
		
		// 创建适配器
		adapter, err = gormadapter.NewAdapterByDB(db)
		if err != nil {
			return nil, fmt.Errorf("创建MySQL适配器失败: %w", err)
		}
	case "postgres":
		// PostgreSQL适配器
		db, err := gorm.Open(postgres.Open(config.DSN), &gorm.Config{})
		if err != nil {
			return nil, fmt.Errorf("连接PostgreSQL数据库失败: %w", err)
		}
		
		// 自动迁移表结构
		if config.AutoMigrate {
			if err := db.AutoMigrate(&gormadapter.CasbinRule{}); err != nil {
				return nil, fmt.Errorf("迁移Casbin表结构失败: %w", err)
			}
		}
		
		// 创建适配器
		adapter, err = gormadapter.NewAdapterByDB(db)
		if err != nil {
			return nil, fmt.Errorf("创建PostgreSQL适配器失败: %w", err)
		}
	default:
		return nil, fmt.Errorf("不支持的适配器类型: %s", config.Adapter)
	}

	// 创建执行器
	enforcer, err = casbin.NewEnforcer(m, adapter)
	if err != nil {
		return nil, fmt.Errorf("创建Casbin执行器失败: %w", err)
	}

	// 加载策略
	if err := enforcer.LoadPolicy(); err != nil {
		return nil, fmt.Errorf("加载Casbin策略失败: %w", err)
	}

	return &CasbinService{
		enforcer: enforcer,
		config:   config,
	}, nil
}

// CheckPermission 检查权限
func (s *CasbinService) CheckPermission(sub string, obj string, act string) (bool, error) {
	return s.enforcer.Enforce(sub, obj, act)
}

// AddPermissionForUser 为用户添加权限
func (s *CasbinService) AddPermissionForUser(sub string, obj string, act string) error {
	_, err := s.enforcer.AddPolicy(sub, obj, act)
	if err != nil {
		return fmt.Errorf("添加权限失败: %w", err)
	}
	return nil
}

// RemovePermissionForUser 移除用户权限
func (s *CasbinService) RemovePermissionForUser(sub string, obj string, act string) error {
	_, err := s.enforcer.RemovePolicy(sub, obj, act)
	if err != nil {
		return fmt.Errorf("移除权限失败: %w", err)
	}
	return nil
}

// AddRoleForUser 为用户添加角色
func (s *CasbinService) AddRoleForUser(user string, role string) error {
	_, err := s.enforcer.AddGroupingPolicy(user, role)
	if err != nil {
		return fmt.Errorf("添加角色失败: %w", err)
	}
	return nil
}

// RemoveRoleForUser 移除用户角色
func (s *CasbinService) RemoveRoleForUser(user string, role string) error {
	_, err := s.enforcer.RemoveGroupingPolicy(user, role)
	if err != nil {
		return fmt.Errorf("移除角色失败: %w", err)
	}
	return nil
}

// AddPermissionForRole 为角色添加权限
func (s *CasbinService) AddPermissionForRole(role string, obj string, act string) error {
	_, err := s.enforcer.AddPolicy(role, obj, act)
	if err != nil {
		return fmt.Errorf("添加角色权限失败: %w", err)
	}
	return nil
}

// RemovePermissionForRole 移除角色权限
func (s *CasbinService) RemovePermissionForRole(role string, obj string, act string) error {
	_, err := s.enforcer.RemovePolicy(role, obj, act)
	if err != nil {
		return fmt.Errorf("移除角色权限失败: %w", err)
	}
	return nil
}

// GetRolesForUser 获取用户角色
func (s *CasbinService) GetRolesForUser(user string) ([]string, error) {
	return s.enforcer.GetRolesForUser(user)
}

// GetUsersForRole 获取角色用户
func (s *CasbinService) GetUsersForRole(role string) ([]string, error) {
	return s.enforcer.GetUsersForRole(role)
}

// GetPermissionsForUser 获取用户权限
func (s *CasbinService) GetPermissionsForUser(user string) ([][]string, error) {
	return s.enforcer.GetPermissionsForUser(user)
}

// GetPermissionsForRole 获取角色权限
func (s *CasbinService) GetPermissionsForRole(role string) ([][]string, error) {
	return s.enforcer.GetPermissionsForUser(role)
}

// HasRoleForUser 检查用户是否有角色
func (s *CasbinService) HasRoleForUser(user string, role string) (bool, error) {
	return s.enforcer.HasRoleForUser(user, role)
}

// GetAllRoles 获取所有角色
func (s *CasbinService) GetAllRoles() ([]string, error) {
	return s.enforcer.GetAllRoles()
}

// GetAllObjects 获取所有对象
func (s *CasbinService) GetAllObjects() ([]string, error) {
	return s.enforcer.GetAllObjects()
}

// GetAllSubjects 获取所有主体
func (s *CasbinService) GetAllSubjects() ([]string, error) {
	return s.enforcer.GetAllSubjects()
}

// LoadPolicy 加载策略
func (s *CasbinService) LoadPolicy() error {
	return s.enforcer.LoadPolicy()
}

// SavePolicy 保存策略
func (s *CasbinService) SavePolicy() error {
	return s.enforcer.SavePolicy()
}

// GetEnforcer 获取执行器
func (s *CasbinService) GetEnforcer() *casbin.Enforcer {
	return s.enforcer
}
