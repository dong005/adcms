package middleware

import (
	"adcms/internal/model"
	"adcms/pkg/database"
)

// 角色权限等级，数字越小权限越高
const (
	RoleLevelSuperAdmin = 0
	RoleLevelAdmin      = 1
	RoleLevelUser       = 2
	RoleLevelNone       = 99
)

// RoleCodeToLevel 将角色编码转换为权限等级
func RoleCodeToLevel(code string) int {
	switch code {
	case "super_admin":
		return RoleLevelSuperAdmin
	case "admin":
		return RoleLevelAdmin
	case "user":
		return RoleLevelUser
	default:
		return RoleLevelNone
	}
}

// GetUserRoleLevel 获取用户的最高权限等级（数字越小权限越高）
func GetUserRoleLevel(userID uint) int {
	var roles []model.Role
	database.DB.Raw(`
		SELECT r.* FROM roles r
		INNER JOIN user_roles ur ON ur.role_id = r.id
		WHERE ur.user_id = ? AND r.deleted_at IS NULL
	`, userID).Scan(&roles)

	level := RoleLevelNone
	for _, role := range roles {
		l := RoleCodeToLevel(role.Code)
		if l < level {
			level = l
		}
	}
	return level
}

// IsSuperAdmin 判断用户是否为超级管理员
func IsSuperAdmin(userID uint) bool {
	return GetUserRoleLevel(userID) == RoleLevelSuperAdmin
}

// HasHigherOrEqualLevel 判断操作者是否有权操作目标用户
// 返回 true 表示操作者权限 >= 目标用户权限（可以操作）
func HasHigherLevel(operatorID, targetID uint) bool {
	operatorLevel := GetUserRoleLevel(operatorID)
	targetLevel := GetUserRoleLevel(targetID)
	return operatorLevel < targetLevel
}

// HasHigherOrEqualLevel 判断操作者权限是否 >= 目标
func HasHigherOrEqualLevel(operatorID, targetID uint) bool {
	operatorLevel := GetUserRoleLevel(operatorID)
	targetLevel := GetUserRoleLevel(targetID)
	return operatorLevel <= targetLevel
}

// CanOperateRole 判断操作者是否有权操作某个角色
func CanOperateRole(operatorID uint, roleCode string) bool {
	operatorLevel := GetUserRoleLevel(operatorID)
	targetLevel := RoleCodeToLevel(roleCode)
	// 只能操作比自己权限低的角色
	return operatorLevel < targetLevel
}

// CanAssignRoles 判断操作者是否有权分配指定角色列表
func CanAssignRoles(operatorID uint, roleIDs []uint) bool {
	operatorLevel := GetUserRoleLevel(operatorID)
	for _, roleID := range roleIDs {
		var role model.Role
		if err := database.DB.First(&role, roleID).Error; err != nil {
			continue
		}
		targetLevel := RoleCodeToLevel(role.Code)
		// 不能分配比自己权限高或相同的角色
		if targetLevel <= operatorLevel {
			return false
		}
	}
	return true
}
