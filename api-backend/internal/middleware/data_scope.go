package middleware

import (
	"adcms/internal/model"
	"adcms/internal/repository"
	"adcms/pkg/database"

	"github.com/gin-gonic/gin"
)

const ContextDataScope = "data_scope"
const ContextDeptIDs = "data_dept_ids"

// DataScopeFilter 数据权限过滤中间件
// 根据用户角色的 data_scope 设置，自动注入可访问的部门ID列表
func DataScopeFilter() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		tenantID := GetTenantID(c)

		if userID == 0 {
			c.Next()
			return
		}

		// super_admin 不限制
		if IsSuperAdmin(userID) {
			c.Set(ContextDataScope, int8(model.DataScopeAll))
			c.Next()
			return
		}

		// 获取用户角色中最宽的数据权限范围
		var roles []model.Role
		database.DB.Raw(`
			SELECT r.* FROM roles r
			INNER JOIN user_roles ur ON ur.role_id = r.id
			WHERE ur.user_id = ? AND r.deleted_at IS NULL
		`, userID).Scan(&roles)

		bestScope := int8(model.DataScopeSelf) // 默认仅自己
		for _, role := range roles {
			if role.DataScope < bestScope {
				bestScope = role.DataScope
			}
		}

		c.Set(ContextDataScope, bestScope)

		switch bestScope {
		case model.DataScopeAll:
			// 全部数据，不限制
		case model.DataScopeDeptTree:
			// 本部门及下级
			var user model.User
			database.DB.First(&user, userID)
			if user.DepartmentID > 0 {
				deptRepo := repository.NewDepartmentRepository()
				deptIDs := deptRepo.GetChildDeptIDs(tenantID, user.DepartmentID)
				c.Set(ContextDeptIDs, deptIDs)
			}
		case model.DataScopeDept:
			// 仅本部门
			var user model.User
			database.DB.First(&user, userID)
			if user.DepartmentID > 0 {
				c.Set(ContextDeptIDs, []uint{user.DepartmentID})
			}
		case model.DataScopeSelf:
			// 仅自己，不设置部门ID
		}

		c.Next()
	}
}

// GetDataScope 获取当前请求的数据权限范围
func GetDataScope(c *gin.Context) int8 {
	scope, exists := c.Get(ContextDataScope)
	if !exists {
		return model.DataScopeSelf
	}
	return scope.(int8)
}

// GetDataDeptIDs 获取当前请求可访问的部门ID列表
func GetDataDeptIDs(c *gin.Context) []uint {
	ids, exists := c.Get(ContextDeptIDs)
	if !exists {
		return nil
	}
	return ids.([]uint)
}

// ApplyDataScope 将数据权限应用到查询（通用辅助函数）
// 调用方在 handler 中使用：db = middleware.ApplyDataScope(c, db, "user_id", "department_id")
func ApplyDataScope(c *gin.Context, db *interface{}) {
	// 此函数由各 handler 根据需要自行调用
	// 参见 GetDataScope() 和 GetDataDeptIDs()
}
