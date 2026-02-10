package middleware

import (
	"adcms/pkg/database"
	"adcms/pkg/utils"
	"context"
	"fmt"
	"time"

	"github.com/gin-gonic/gin"
)

// RateLimit IP 级接口限流中间件（滑动窗口）
// maxRequests: 窗口内最大请求数, window: 时间窗口
func RateLimit(maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:%s:%s", c.Request.URL.Path, ip)

		current, _ := database.RDB.Incr(ctx, key).Result()
		if current == 1 {
			database.RDB.Expire(ctx, key, window)
		}

		if current > maxRequests {
			utils.Fail(c, 4029, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// UserRateLimit 用户级接口限流
func UserRateLimit(maxRequests int64, window time.Duration) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		userID := GetUserID(c)
		if userID == 0 {
			c.Next()
			return
		}

		key := fmt.Sprintf("ratelimit:user:%d:%s", userID, c.Request.URL.Path)
		current, _ := database.RDB.Incr(ctx, key).Result()
		if current == 1 {
			database.RDB.Expire(ctx, key, window)
		}

		if current > maxRequests {
			utils.Fail(c, 4029, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}

// LoginFailLimit 登录失败次数限制
// maxAttempts: 最大失败次数, lockDuration: 锁定时长
const (
	LoginMaxAttempts = 5
	LoginLockDuration = 15 * time.Minute
	loginFailPrefix  = "login:fail:"
	loginLockPrefix  = "login:lock:"
)

// IsLoginLocked 检查账号是否被锁定
func IsLoginLocked(username string) (bool, int64) {
	ctx := context.Background()
	lockKey := loginLockPrefix + username

	ttl, err := database.RDB.TTL(ctx, lockKey).Result()
	if err != nil || ttl <= 0 {
		return false, 0
	}

	return true, int64(ttl.Seconds())
}

// RecordLoginFail 记录登录失败
// 返回: 剩余尝试次数, 是否被锁定
func RecordLoginFail(username string) (int64, bool) {
	ctx := context.Background()
	failKey := loginFailPrefix + username

	count, _ := database.RDB.Incr(ctx, failKey).Result()
	if count == 1 {
		database.RDB.Expire(ctx, failKey, LoginLockDuration)
	}

	if count >= int64(LoginMaxAttempts) {
		// 锁定账号
		lockKey := loginLockPrefix + username
		database.RDB.Set(ctx, lockKey, "1", LoginLockDuration)
		database.RDB.Del(ctx, failKey)
		return 0, true
	}

	remaining := int64(LoginMaxAttempts) - count
	return remaining, false
}

// ClearLoginFail 登录成功后清除失败记录
func ClearLoginFail(username string) {
	ctx := context.Background()
	database.RDB.Del(ctx, loginFailPrefix+username)
	database.RDB.Del(ctx, loginLockPrefix+username)
}

// GlobalRateLimit 全局 API 限流（每个 IP 每分钟最多 N 次）
func GlobalRateLimit(maxPerMinute int64) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := context.Background()
		ip := c.ClientIP()
		key := fmt.Sprintf("ratelimit:global:%s", ip)

		current, _ := database.RDB.Incr(ctx, key).Result()
		if current == 1 {
			database.RDB.Expire(ctx, key, time.Minute)
		}

		if current > maxPerMinute {
			utils.Fail(c, 4029, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}

		c.Next()
	}
}
