package middleware

import (
	"adcms/pkg/database"
	"adcms/pkg/utils"
	"context"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
)

const (
	ContextUserID   = "user_id"
	ContextTenantID = "tenant_id"
	ContextUsername = "username"
	ContextIsAdmin  = "is_admin"
)

func JWTAuth() gin.HandlerFunc {
	return func(c *gin.Context) {
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			utils.Unauthorized(c, "请先登录")
			c.Abort()
			return
		}

		parts := strings.SplitN(authHeader, " ", 2)
		if !(len(parts) == 2 && parts[0] == "Bearer") {
			utils.Unauthorized(c, "认证格式错误")
			c.Abort()
			return
		}

		token := parts[1]

		blacklisted, err := isTokenBlacklisted(token)
		if err != nil || blacklisted {
			utils.Unauthorized(c, "token已失效")
			c.Abort()
			return
		}

		claims, err := utils.ParseToken(token)
		if err != nil {
			utils.Unauthorized(c, "token无效或已过期")
			c.Abort()
			return
		}

		c.Set(ContextUserID, claims.UserID)
		c.Set(ContextTenantID, claims.TenantID)
		c.Set(ContextUsername, claims.Username)
		c.Set(ContextIsAdmin, claims.IsAdmin)

		c.Next()
	}
}

func isTokenBlacklisted(token string) (bool, error) {
	ctx := context.Background()
	key := "token:blacklist:" + token
	result, err := database.RDB.Exists(ctx, key).Result()
	if err != nil {
		return false, err
	}
	return result > 0, nil
}

func BlacklistToken(token string, expireTime time.Duration) error {
	ctx := context.Background()
	key := "token:blacklist:" + token
	return database.RDB.Set(ctx, key, "1", expireTime).Err()
}

func GetUserID(c *gin.Context) uint {
	userID, exists := c.Get(ContextUserID)
	if !exists {
		return 0
	}
	return userID.(uint)
}

func GetTenantID(c *gin.Context) uint {
	tenantID, exists := c.Get(ContextTenantID)
	if !exists {
		return 0
	}
	return tenantID.(uint)
}

func GetUsername(c *gin.Context) string {
	username, exists := c.Get(ContextUsername)
	if !exists {
		return ""
	}
	return username.(string)
}

func GetIsAdmin(c *gin.Context) int8 {
	isAdmin, exists := c.Get(ContextIsAdmin)
	if !exists {
		return 0
	}
	return isAdmin.(int8)
}
