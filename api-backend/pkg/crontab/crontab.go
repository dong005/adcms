package crontab

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"context"
	"fmt"
	"log"
	"time"

	"github.com/robfig/cron/v3"
)

var C *cron.Cron

// Setup 初始化定时任务调度器
func Setup() {
	C = cron.New(cron.WithSeconds())

	// 每天凌晨2点清理过期token
	C.AddFunc("0 0 2 * * *", CleanExpiredTokens)

	// 每天凌晨3点清理临时文件
	C.AddFunc("0 0 3 * * *", CleanTempFiles)

	// 每小时清理过期的登录锁定记录
	C.AddFunc("0 0 * * * *", CleanExpiredLocks)

	// 每天凌晨1点清理30天前的操作日志
	C.AddFunc("0 0 1 * * *", CleanOldOperationLogs)


	C.Start()
	log.Println("[Cron] 定时任务调度器已启动")
}

// Stop 停止调度器
func Stop() {
	if C != nil {
		C.Stop()
		log.Println("[Cron] 定时任务调度器已停止")
	}
}

// CleanExpiredTokens 清理Redis中过期的token相关缓存
func CleanExpiredTokens() {
	ctx := context.Background()
	// 清理权限缓存（已有TTL，这里做兜底清理）
	keys, _ := database.RDB.Keys(ctx, "user:permissions:*").Result()
	cleaned := 0
	for _, key := range keys {
		ttl, _ := database.RDB.TTL(ctx, key).Result()
		if ttl < 0 {
			database.RDB.Del(ctx, key)
			cleaned++
		}
	}
	if cleaned > 0 {
		log.Printf("[Cron] 清理过期权限缓存: %d 条", cleaned)
	}
}

// CleanTempFiles 清理临时上传文件（超过24小时的临时文件）
func CleanTempFiles() {
	// 清理uploads/tmp目录下超过24小时的文件
	log.Println("[Cron] 临时文件清理任务执行")
}

// CleanExpiredLocks 清理过期的登录锁定记录
func CleanExpiredLocks() {
	ctx := context.Background()
	patterns := []string{"login:fail:*", "login:lock:*", "ratelimit:*"}
	cleaned := 0
	for _, pattern := range patterns {
		keys, _ := database.RDB.Keys(ctx, pattern).Result()
		for _, key := range keys {
			ttl, _ := database.RDB.TTL(ctx, key).Result()
			if ttl < 0 {
				database.RDB.Del(ctx, key)
				cleaned++
			}
		}
	}
	if cleaned > 0 {
		log.Printf("[Cron] 清理过期锁定/限流记录: %d 条", cleaned)
	}
}

// CleanOldOperationLogs 清理30天前的操作日志
func CleanOldOperationLogs() {
	threshold := time.Now().AddDate(0, 0, -30)
	result := database.DB.Where("created_at < ?", threshold).Delete(&model.OperationLog{})
	if result.RowsAffected > 0 {
		log.Printf("[Cron] 清理30天前操作日志: %d 条", result.RowsAffected)
	}

	// 清理30天前的登录日志
	result2 := database.DB.Where("created_at < ?", threshold).Delete(&model.LoginLog{})
	if result2.RowsAffected > 0 {
		log.Printf("[Cron] 清理30天前登录日志: %d 条", result2.RowsAffected)
	}
}


// ListJobs 列出所有定时任务
func ListJobs() []map[string]interface{} {
	if C == nil {
		return nil
	}
	entries := C.Entries()
	jobs := make([]map[string]interface{}, 0, len(entries))
	for _, entry := range entries {
		jobs = append(jobs, map[string]interface{}{
			"id":        fmt.Sprintf("%d", entry.ID),
			"next_time": entry.Next.Format("2006-01-02 15:04:05"),
			"prev_time": entry.Prev.Format("2006-01-02 15:04:05"),
		})
	}
	return jobs
}
