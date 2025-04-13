package auth

import (
	"context"
	"fmt"
)

// HasPermission 检查用户是否拥有指定权限
func (a *Auth) HasPermission(ctx context.Context, token string, obj string, act string) (bool, error) {
	// 获取用户ID
	userID, err := a.GetUserID(ctx, token)
	if err != nil {
		return false, err
	}

	// 使用Casbin检查权限
	return a.Enforcer.HasPermission(userID, obj, act)
}

// HasRole 检查用户是否拥有指定角色
func (a *Auth) HasRole(ctx context.Context, token string, role string) (bool, error) {
	// 获取用户ID
	userID, err := a.GetUserID(ctx, token)
	if err != nil {
		return false, err
	}

	// 使用Casbin检查角色
	return a.Enforcer.HasRole(userID, role)
}

// EnforcePermission 强制用户拥有指定权限
func (a *Auth) EnforcePermission(ctx context.Context, token string, obj string, act string) error {
	// 获取用户ID
	userID, err := a.GetUserID(ctx, token)
	if err != nil {
		return err
	}

	// 使用Casbin强制检查权限
	return a.Enforcer.EnforcePermission(userID, obj, act)
}

// EnforceRole 强制用户拥有指定角色
func (a *Auth) EnforceRole(ctx context.Context, token string, role string) error {
	// 获取用户ID
	userID, err := a.GetUserID(ctx, token)
	if err != nil {
		return err
	}

	// 使用Casbin强制检查角色
	return a.Enforcer.EnforceRole(userID, role)
}

// AssignRoleToUser 为用户分配角色
func (a *Auth) AssignRoleToUser(ctx context.Context, userID int64, role string) error {
	// 使用Casbin分配角色
	return a.Enforcer.AssignRole(userID, role)
}

// RevokeRoleFromUser 从用户撤销角色
func (a *Auth) RevokeRoleFromUser(ctx context.Context, userID int64, role string) error {
	// 使用Casbin撤销角色
	return a.Enforcer.RevokeRole(userID, role)
}

// AddPermissionForRole 为角色添加权限
func (a *Auth) AddPermissionForRole(ctx context.Context, role string, obj string, act string) error {
	// 使用Casbin添加权限
	return a.Enforcer.AddPermissionForRole(role, obj, act)
}

// RemovePermissionFromRole 从角色移除权限
func (a *Auth) RemovePermissionFromRole(ctx context.Context, role string, obj string, act string) error {
	// 使用Casbin移除权限
	return a.Enforcer.RemovePermissionForRole(role, obj, act)
}

// AddPermissionForUser 为用户添加权限
func (a *Auth) AddPermissionForUser(ctx context.Context, userID int64, obj string, act string) error {
	// 使用Casbin添加权限
	return a.Enforcer.AddPermissionForUser(userID, obj, act)
}

// RemovePermissionFromUser 从用户移除权限
func (a *Auth) RemovePermissionFromUser(ctx context.Context, userID int64, obj string, act string) error {
	// 使用Casbin移除权限
	return a.Enforcer.RemovePermissionForUser(userID, obj, act)
}

// GetRoles 获取用户的所有角色
func (a *Auth) GetRoles(ctx context.Context, token string) ([]string, error) {
	// 获取用户ID
	userID, err := a.GetUserID(ctx, token)
	if err != nil {
		return nil, err
	}

	// 使用Casbin获取角色
	return a.Enforcer.GetRolesForUser(userID)
}

// GetPermissions 获取用户的所有权限
func (a *Auth) GetPermissions(ctx context.Context, token string) ([]string, error) {
	// 获取用户ID
	userID, err := a.GetUserID(ctx, token)
	if err != nil {
		return nil, err
	}

	// 使用Casbin获取权限
	perms, err := a.Enforcer.GetPermissionsForUser(userID)
	if err != nil {
		return nil, err
	}

	// 转换权限格式为字符串数组
	result := make([]string, 0, len(perms))
	for _, perm := range perms {
		if len(perm) >= 3 {
			// 格式化为 "object:action" 形式
			result = append(result, fmt.Sprintf("%s:%s", perm[1], perm[2]))
		}
	}

	return result, nil
}

// GetUserRoles 获取指定用户的所有角色
func (a *Auth) GetUserRoles(ctx context.Context, userID int64) ([]string, error) {
	// 使用Casbin获取角色
	return a.Enforcer.GetRolesForUser(userID)
}

// GetUserPermissions 获取指定用户的所有权限
func (a *Auth) GetUserPermissions(ctx context.Context, userID int64) ([]string, error) {
	// 使用Casbin获取权限
	perms, err := a.Enforcer.GetPermissionsForUser(userID)
	if err != nil {
		return nil, err
	}

	// 转换权限格式为字符串数组
	result := make([]string, 0, len(perms))
	for _, perm := range perms {
		if len(perm) >= 3 {
			// 格式化为 "object:action" 形式
			result = append(result, fmt.Sprintf("%s:%s", perm[1], perm[2]))
		}
	}

	return result, nil
}
