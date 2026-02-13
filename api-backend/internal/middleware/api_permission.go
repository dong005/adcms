package middleware

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"adcms/pkg/utils"
	"context"
	"encoding/json"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const permissionCachePrefix = "user:permissions:"
const permissionCacheTTL = 5 * time.Minute

// GetUserPermissionCodes 获取用户所有权限码（带 Redis 缓存）
// super_admin 返回所有权限码（用于前端按钮级权限控制）
func GetUserPermissionCodes(userID uint) []string {
	ctx := context.Background()
	cacheKey := permissionCachePrefix + formatUint(userID)

	// 尝试从缓存获取
	cached, err := database.RDB.Get(ctx, cacheKey).Result()
	if err == nil && cached != "" {
		var codes []string
		if json.Unmarshal([]byte(cached), &codes) == nil {
			return codes
		}
	}

	var codes []string

	// super_admin 返回所有权限码
	if IsSuperAdmin(userID) {
		var allPerms []model.Permission
		database.DB.Find(&allPerms)
		codes = make([]string, 0, len(allPerms))
		for _, p := range allPerms {
			codes = append(codes, p.Code)
		}
	} else {
		// 从角色权限查询
		var permissions []model.Permission
		database.DB.Raw(`
			SELECT DISTINCT p.* FROM permissions p
			INNER JOIN role_permissions rp ON rp.permission_id = p.id
			INNER JOIN user_roles ur ON ur.role_id = rp.role_id
			WHERE ur.user_id = ?
		`, userID).Scan(&permissions)

		codes = make([]string, 0, len(permissions))
		for _, p := range permissions {
			codes = append(codes, p.Code)
		}

		// TODO: 这里可以扩展用户独立权限表
		// 如果将来需要用户独立权限，可以在这里添加查询逻辑
	}

	// 写入缓存
	if data, err := json.Marshal(codes); err == nil {
		database.RDB.Set(ctx, cacheKey, string(data), permissionCacheTTL)
	}

	return codes
}

// GetUserAPIPermissions 获取用户所有 API 权限（method+path）
func GetUserAPIPermissions(userID uint) []model.Permission {
	var permissions []model.Permission
	database.DB.Raw(`
		SELECT DISTINCT p.* FROM permissions p
		INNER JOIN role_permissions rp ON rp.permission_id = p.id
		INNER JOIN user_roles ur ON ur.role_id = rp.role_id
		WHERE ur.user_id = ? AND p.type = 3 AND p.path != '' AND p.method != ''
	`, userID).Scan(&permissions)
	return permissions
}

// ClearUserPermissionCache 清除用户权限缓存（角色变更时调用）
func ClearUserPermissionCache(userID uint) {
	ctx := context.Background()
	cacheKey := permissionCachePrefix + formatUint(userID)
	database.RDB.Del(ctx, cacheKey)
}

// RequirePermission 基于权限码的中间件，用于单个路由
// 用法: router.PUT("/users/:id", middleware.RequirePermission("user:update"), handler.Update)
func RequirePermission(code string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// super_admin 跳过权限检查
		if IsSuperAdmin(userID) {
			c.Next()
			return
		}

		codes := GetUserPermissionCodes(userID)
		for _, userCode := range codes {
			if userCode == code {
				c.Next()
				return
			}
		}

		utils.Fail(c, 4003, "无操作权限: "+code)
		c.Abort()
	}
}

// RequireAnyPermission 拥有任一权限即可通过
func RequireAnyPermission(codes ...string) gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		if IsSuperAdmin(userID) {
			c.Next()
			return
		}

		userCodes := GetUserPermissionCodes(userID)
		userCodeMap := make(map[string]bool, len(userCodes))
		for _, uc := range userCodes {
			userCodeMap[uc] = true
		}

		for _, code := range codes {
			if userCodeMap[code] {
				c.Next()
				return
			}
		}

		utils.Fail(c, 4003, "无操作权限")
		c.Abort()
	}
}

// APIPermissionCheck 基于 path+method 的自动权限检查中间件
// 挂载到 protected 路由组，自动匹配 permissions 表中的 path+method
func APIPermissionCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		userID := GetUserID(c)
		if userID == 0 {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		// super_admin 跳过
		if IsSuperAdmin(userID) {
			c.Next()
			return
		}

		method := c.Request.Method
		path := c.Request.URL.Path

		// 检查该 path+method 是否在权限表中注册
		var count int64
		database.DB.Model(&model.Permission{}).
			Where("type = 3 AND method = ? AND path != ''", method).
			Count(&count)

		if count == 0 {
			// 该接口未注册权限，默认放行（仅认证即可访问）
			c.Next()
			return
		}

		// 检查是否有匹配的权限
		var matchedPerm model.Permission
		err := database.DB.Where("type = 3 AND method = ?", method).
			Find(&matchedPerm).Error

		// 使用路径模式匹配
		var allAPIPerms []model.Permission
		database.DB.Where("type = 3 AND method = ? AND path != ''", method).Find(&allAPIPerms)

		matched := false
		var matchedCode string
		for _, perm := range allAPIPerms {
			if matchPath(perm.Path, path) {
				matched = true
				matchedCode = perm.Code
				break
			}
		}

		if !matched {
			// 该具体路径未注册权限，放行
			c.Next()
			return
		}

		// 检查用户是否拥有该权限
		codes := GetUserPermissionCodes(userID)
		for _, code := range codes {
			if code == matchedCode {
				c.Next()
				return
			}
		}

		utils.Fail(c, 4003, "无操作权限: "+matchedCode)
		c.Abort()
		_ = err
	}
}

// matchPath 路径模式匹配，支持 :param 通配符
// 例如 /api/users/:id 匹配 /api/users/123
func matchPath(pattern, path string) bool {
	patternParts := strings.Split(strings.Trim(pattern, "/"), "/")
	pathParts := strings.Split(strings.Trim(path, "/"), "/")

	if len(patternParts) != len(pathParts) {
		return false
	}

	for i, part := range patternParts {
		if strings.HasPrefix(part, ":") {
			continue // 通配符，匹配任意值
		}
		if part != pathParts[i] {
			return false
		}
	}
	return true
}

func formatUint(n uint) string {
	if n == 0 {
		return "0"
	}
	var buf [20]byte
	pos := len(buf)
	for n > 0 {
		pos--
		buf[pos] = byte('0' + n%10)
		n /= 10
	}
	return string(buf[pos:])
}
