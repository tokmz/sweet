package auth

import (
	"fmt"
	"strconv"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	gormadapter "github.com/casbin/gorm-adapter/v3"
	"gorm.io/gorm"
)

// CasbinPermissionManager Casbin权限管理器实现
type CasbinPermissionManager struct {
	// enforcer Casbin执行器
	enforcer *casbin.Enforcer
}

// NewPermissionManager 创建权限管理器
func NewPermissionManager(modelPath string, adapter string, opts ...interface{}) (*CasbinPermissionManager, error) {
	var a *gormadapter.Adapter
	var err error

	// 根据适配器类型创建适配器
	switch adapter {
	case "mysql":
		// 使用GORM适配器
		if len(opts) > 0 {
			if db, ok := opts[0].(*gorm.DB); ok {
				// 使用现有的GORM连接
				a, err = gormadapter.NewAdapterByDB(db)
			} else {
				// 使用连接字符串
				dsn, _ := opts[0].(string)
				a, err = gormadapter.NewAdapter("mysql", dsn)
			}
		} else {
			return nil, fmt.Errorf("缺少数据库连接参数")
		}
	case "memory":
		// 使用内存适配器
		a, err = gormadapter.NewAdapter("sqlite3", ":memory:")
	default:
		return nil, fmt.Errorf("不支持的适配器类型: %s", adapter)
	}

	if err != nil {
		return nil, err
	}

	// 加载模型
	var m model.Model
	if modelPath != "" {
		// 从文件加载模型
		m, err = model.NewModelFromFile(modelPath)
		if err != nil {
			return nil, err
		}
	} else {
		// 使用默认RBAC模型
		m, err = model.NewModelFromString(`
		[request_definition]
		r = sub, obj, act
		
		[policy_definition]
		p = sub, obj, act
		
		[role_definition]
		g = _, _
		
		[policy_effect]
		e = some(where (p.eft == allow))
		
		[matchers]
		m = g(r.sub, p.sub) && r.obj == p.obj && r.act == p.act
		`)
		if err != nil {
			return nil, err
		}
	}

	// 创建执行器
	enforcer, err := casbin.NewEnforcer(m, a)
	if err != nil {
		return nil, err
	}

	// 加载策略
	err = enforcer.LoadPolicy()
	if err != nil {
		return nil, err
	}

	return &CasbinPermissionManager{
		enforcer: enforcer,
	}, nil
}

// CheckPermission 检查权限
func (m *CasbinPermissionManager) CheckPermission(userID uint64, obj string, act string) bool {
	// 将用户ID转换为字符串
	sub := strconv.FormatUint(userID, 10)

	// 检查权限
	result, err := m.enforcer.Enforce(sub, obj, act)
	if err != nil {
		return false
	}

	return result
}

// AddRoleForUser 为用户添加角色
func (m *CasbinPermissionManager) AddRoleForUser(userID uint64, role string) error {
	// 将用户ID转换为字符串
	sub := strconv.FormatUint(userID, 10)

	// 添加角色
	_, err := m.enforcer.AddRoleForUser(sub, role)
	if err != nil {
		return err
	}

	// 保存策略
	return m.enforcer.SavePolicy()
}

// DeleteRoleForUser 删除用户的角色
func (m *CasbinPermissionManager) DeleteRoleForUser(userID uint64, role string) error {
	// 将用户ID转换为字符串
	sub := strconv.FormatUint(userID, 10)

	// 删除角色
	_, err := m.enforcer.DeleteRoleForUser(sub, role)
	if err != nil {
		return err
	}

	// 保存策略
	return m.enforcer.SavePolicy()
}

// GetRolesForUser 获取用户的所有角色
func (m *CasbinPermissionManager) GetRolesForUser(userID uint64) ([]string, error) {
	// 将用户ID转换为字符串
	sub := strconv.FormatUint(userID, 10)

	// 获取角色
	return m.enforcer.GetRolesForUser(sub)
}

// AddPermissionForRole 为角色添加权限
func (m *CasbinPermissionManager) AddPermissionForRole(role, obj, act string) error {
	// 添加权限
	_, err := m.enforcer.AddPolicy(role, obj, act)
	if err != nil {
		return err
	}

	// 保存策略
	return m.enforcer.SavePolicy()
}

// DeletePermissionForRole 删除角色的权限
func (m *CasbinPermissionManager) DeletePermissionForRole(role, obj, act string) error {
	// 删除权限
	_, err := m.enforcer.RemovePolicy(role, obj, act)
	if err != nil {
		return err
	}

	// 保存策略
	return m.enforcer.SavePolicy()
}

// GetPermissionsForRole 获取角色的所有权限
func (m *CasbinPermissionManager) GetPermissionsForRole(role string) ([][]string, error) {
	// 获取权限
	if m.enforcer == nil {
		return nil, fmt.Errorf("enforcer未初始化")
	}

	// 获取过滤后的策略
	policies, err := m.enforcer.GetFilteredPolicy(0, role)
	if err != nil {
		return nil, err
	}

	// 返回结果
	return policies, nil
}

// GetEnforcer 获取Casbin执行器
func (m *CasbinPermissionManager) GetEnforcer() *casbin.Enforcer {
	return m.enforcer
}
