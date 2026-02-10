package logcfg

import (
	"sync"
	"time"

	"gorm.io/gorm"
)

var (
	db             *gorm.DB
	logConfigCache = make(map[string]string)
	logConfigMu    sync.RWMutex
	logConfigTime  time.Time
	cacheTTL       = time.Minute
)

// Init 初始化，传入 DB 实例避免循环导入
func Init(database *gorm.DB) {
	db = database
}

// IsLogEnabled 检查某类日志是否启用，带1分钟内存缓存
func IsLogEnabled(key string) bool {
	if db == nil {
		return true
	}

	logConfigMu.RLock()
	if time.Since(logConfigTime) < cacheTTL {
		val, ok := logConfigCache[key]
		logConfigMu.RUnlock()
		if ok {
			return val != "0"
		}
		return true
	}
	logConfigMu.RUnlock()

	logConfigMu.Lock()
	defer logConfigMu.Unlock()

	// 双重检查
	if time.Since(logConfigTime) < cacheTTL {
		val, ok := logConfigCache[key]
		if ok {
			return val != "0"
		}
		return true
	}

	type kv struct {
		Key   string
		Value string
	}
	var configs []kv
	db.Table("system_configs").Select("`key`, value").Where("`key` LIKE 'log_%_enabled'").Find(&configs)

	newCache := make(map[string]string)
	for _, c := range configs {
		newCache[c.Key] = c.Value
	}
	logConfigCache = newCache
	logConfigTime = time.Now()

	val, ok := newCache[key]
	if ok {
		return val != "0"
	}
	return true
}

// ClearCache 清除日志配置缓存（配置更新后调用）
func ClearCache() {
	logConfigMu.Lock()
	defer logConfigMu.Unlock()
	logConfigTime = time.Time{}
}
