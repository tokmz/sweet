package auth

import (
	"errors"
	"fmt"
	"sync"

	"github.com/casbin/casbin/v2"
	"github.com/casbin/casbin/v2/model"
	"github.com/casbin/casbin/v2/persist"
)

var (
	// ErrNoPermission 无权限
	ErrNoPermission = errors.New("no permission")
	// ErrNoRole 无角色
	ErrNoRole = errors.New("no role")
)

// CasbinEnforcer 是Casbin权限检查接口
type CasbinEnforcer interface {
	// HasPermission 检查用户是否有访问某个资源的权限
	HasPermission(userID int64, obj string, act string) (bool, error)
	// HasRole 检查用户是否拥有某个角色
	HasRole(userID int64, role string) (bool, error)
	// AssignRole 为用户分配角色
	AssignRole(userID int64, role string) error
	// RevokeRole 撤销用户的角色
	RevokeRole(userID int64, role string) error
	// AddPermissionForRole 为角色添加权限
	AddPermissionForRole(role string, obj string, act string) error
	// RemovePermissionForRole 从角色中移除权限
	RemovePermissionForRole(role string, obj string, act string) error
	// AddPermissionForUser 为用户添加权限
	AddPermissionForUser(userID int64, obj string, act string) error
	// RemovePermissionForUser 移除用户权限
	RemovePermissionForUser(userID int64, obj string, act string) error
	// GetRolesForUser 获取用户的所有角色
	GetRolesForUser(userID int64) ([]string, error)
	// GetPermissionsForUser 获取用户的所有权限
	GetPermissionsForUser(userID int64) ([][]string, error)
	// EnforceRole 强制用户拥有角色
	EnforceRole(userID int64, role string) error
	// EnforcePermission 强制用户拥有权限
	EnforcePermission(userID int64, obj string, act string) error
}

// DefaultCasbinModel 默认的Casbin RBAC模型
const DefaultCasbinModel = `
[request_definition]
r = sub, obj, act

[policy_definition]
p = sub, obj, act

[role_definition]
g = _, _

[policy_effect]
e = some(where (p.eft == allow))

[matchers]
m = g(r.sub, p.sub) && (r.obj == p.obj || p.obj == "*") && (r.act == p.act || p.act == "*")
`

// CasbinAdapter 是Casbin适配器接口
type CasbinAdapter interface {
	persist.Adapter
}

// MemoryCasbinAdapter 是基于内存的Casbin适配器
type MemoryCasbinAdapter struct {
	policies [][]string
	mutex    sync.RWMutex
}

// NewMemoryCasbinAdapter 创建内存Casbin适配器
func NewMemoryCasbinAdapter() *MemoryCasbinAdapter {
	return &MemoryCasbinAdapter{
		policies: make([][]string, 0),
	}
}

// LoadPolicy 从内存中加载策略
func (a *MemoryCasbinAdapter) LoadPolicy(model model.Model) error {
	a.mutex.RLock()
	defer a.mutex.RUnlock()

	for _, p := range a.policies {
		if len(p) < 3 {
			continue
		}

		lineText := "p, " + p[0] + ", " + p[1] + ", " + p[2]
		persist.LoadPolicyLine(lineText, model)
	}

	return nil
}

// SavePolicy 保存策略到内存
func (a *MemoryCasbinAdapter) SavePolicy(model model.Model) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.policies = make([][]string, 0)

	for _, ast := range model["p"] {
		for _, rule := range ast.Policy {
			if len(rule) < 3 {
				continue
			}
			a.policies = append(a.policies, rule)
		}
	}

	for _, ast := range model["g"] {
		for _, rule := range ast.Policy {
			if len(rule) < 2 {
				continue
			}
			a.policies = append(a.policies, rule)
		}
	}

	return nil
}

// AddPolicy 添加策略
func (a *MemoryCasbinAdapter) AddPolicy(sec string, ptype string, rule []string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	a.policies = append(a.policies, rule)
	return nil
}

// RemovePolicy 移除策略
func (a *MemoryCasbinAdapter) RemovePolicy(sec string, ptype string, rule []string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	for i, r := range a.policies {
		if arrayEquals(r, rule) {
			a.policies = append(a.policies[:i], a.policies[i+1:]...)
			return nil
		}
	}

	return nil
}

// RemoveFilteredPolicy 根据条件移除策略
func (a *MemoryCasbinAdapter) RemoveFilteredPolicy(sec string, ptype string, fieldIndex int, fieldValues ...string) error {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	var temp [][]string
	for _, r := range a.policies {
		matched := true
		for i, v := range fieldValues {
			if v != "" && fieldIndex+i < len(r) && r[fieldIndex+i] != v {
				matched = false
				break
			}
		}
		if !matched {
			temp = append(temp, r)
		}
	}
	a.policies = temp

	return nil
}

// DefaultCasbinEnforcer 是基于Casbin的权限验证实现
type DefaultCasbinEnforcer struct {
	enforcer *casbin.Enforcer
}

// NewCasbinEnforcer 创建Casbin执行器
func NewCasbinEnforcer() CasbinEnforcer {
	// 创建Casbin模型
	m, _ := model.NewModelFromString(DefaultCasbinModel)

	// 创建内存适配器
	adapter := NewMemoryCasbinAdapter()

	// 创建执行器
	e, _ := casbin.NewEnforcer(m, adapter)

	// 加载策略
	e.LoadPolicy()

	return &DefaultCasbinEnforcer{
		enforcer: e,
	}
}

// HasPermission 检查用户是否有访问某个资源的权限
func (e *DefaultCasbinEnforcer) HasPermission(userID int64, obj string, act string) (bool, error) {
	sub := fmt.Sprintf("user:%d", userID)
	return e.enforcer.Enforce(sub, obj, act)
}

// HasRole 检查用户是否拥有某个角色
func (e *DefaultCasbinEnforcer) HasRole(userID int64, role string) (bool, error) {
	sub := fmt.Sprintf("user:%d", userID)
	return e.enforcer.HasGroupingPolicy(sub, "role:"+role)
}

// AssignRole 为用户分配角色
func (e *DefaultCasbinEnforcer) AssignRole(userID int64, role string) error {
	sub := fmt.Sprintf("user:%d", userID)
	_, err := e.enforcer.AddGroupingPolicy(sub, "role:"+role)
	return err
}

// RevokeRole 撤销用户的角色
func (e *DefaultCasbinEnforcer) RevokeRole(userID int64, role string) error {
	sub := fmt.Sprintf("user:%d", userID)
	_, err := e.enforcer.RemoveGroupingPolicy(sub, "role:"+role)
	return err
}

// AddPermissionForRole 为角色添加权限
func (e *DefaultCasbinEnforcer) AddPermissionForRole(role string, obj string, act string) error {
	_, err := e.enforcer.AddPolicy("role:"+role, obj, act)
	return err
}

// RemovePermissionForRole 从角色中移除权限
func (e *DefaultCasbinEnforcer) RemovePermissionForRole(role string, obj string, act string) error {
	_, err := e.enforcer.RemovePolicy("role:"+role, obj, act)
	return err
}

// AddPermissionForUser 为用户添加权限
func (e *DefaultCasbinEnforcer) AddPermissionForUser(userID int64, obj string, act string) error {
	sub := fmt.Sprintf("user:%d", userID)
	_, err := e.enforcer.AddPolicy(sub, obj, act)
	return err
}

// RemovePermissionForUser 移除用户权限
func (e *DefaultCasbinEnforcer) RemovePermissionForUser(userID int64, obj string, act string) error {
	sub := fmt.Sprintf("user:%d", userID)
	_, err := e.enforcer.RemovePolicy(sub, obj, act)
	return err
}

// GetRolesForUser 获取用户的所有角色
func (e *DefaultCasbinEnforcer) GetRolesForUser(userID int64) ([]string, error) {
	sub := fmt.Sprintf("user:%d", userID)
	roles, err := e.enforcer.GetRolesForUser(sub)
	if err != nil {
		return nil, err
	}

	// 去掉角色前缀"role:"
	result := make([]string, 0, len(roles))
	for _, role := range roles {
		if len(role) > 5 && role[:5] == "role:" {
			result = append(result, role[5:])
		}
	}

	return result, nil
}

// GetPermissionsForUser 获取用户的所有权限
func (e *DefaultCasbinEnforcer) GetPermissionsForUser(userID int64) ([][]string, error) {
	sub := fmt.Sprintf("user:%d", userID)
	perms, err := e.enforcer.GetPermissionsForUser(sub)
	if err != nil {
		return nil, err
	}
	return perms, nil
}

// EnforceRole 强制用户拥有角色
func (e *DefaultCasbinEnforcer) EnforceRole(userID int64, role string) error {
	has, err := e.HasRole(userID, role)
	if err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("%w: %s", ErrNoRole, role)
	}

	return nil
}

// EnforcePermission 强制用户拥有权限
func (e *DefaultCasbinEnforcer) EnforcePermission(userID int64, obj string, act string) error {
	has, err := e.HasPermission(userID, obj, act)
	if err != nil {
		return err
	}

	if !has {
		return fmt.Errorf("%w: %s %s", ErrNoPermission, obj, act)
	}

	return nil
}

// 工具函数，检查两个字符串数组是否相等
func arrayEquals(a []string, b []string) bool {
	if len(a) != len(b) {
		return false
	}

	for i, v := range a {
		if v != b[i] {
			return false
		}
	}

	return true
}
