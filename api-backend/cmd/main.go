// @title ADCMS API
// @version 1.0
// @description ADCMS 多租户内容管理系统 API 文档
// @host localhost:8004
// @BasePath /api
// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
package main

import (
	"adcms/internal/config"
	"adcms/internal/router"
	"adcms/pkg/crontab"
	"adcms/pkg/database"
	"adcms/pkg/logcfg"
	"adcms/pkg/logger"
	"adcms/pkg/storage"
	"fmt"
	"os"
)

func main() {
	configPath := "config.yaml"
	if len(os.Args) > 1 {
		configPath = os.Args[1]
	}

	cfg, err := config.LoadConfig(configPath)
	if err != nil {
		fmt.Printf("Failed to load config: %v\n", err)
		os.Exit(1)
	}

	if err := logger.InitLogger(&cfg.Log); err != nil {
		fmt.Printf("Failed to init logger: %v\n", err)
		os.Exit(1)
	}

	if err := database.InitMySQL(&cfg.MySQL); err != nil {
		logger.Fatalf("Failed to init MySQL: %v", err)
	}
	defer database.CloseMySQL()

	logcfg.Init(database.DB)

	if err := database.AutoMigrate(); err != nil {
		logger.Fatalf("Failed to auto migrate: %v", err)
	}

	if err := database.InitData(); err != nil {
		logger.Warnf("Init data warning: %v", err)
	}

	if err := database.InitRedis(&cfg.Redis); err != nil {
		logger.Fatalf("Failed to init Redis: %v", err)
	}
	defer database.CloseRedis()

	// 初始化文件存储
	initStorage(&cfg.Storage)

	// 启动定时任务
	crontab.Setup()
	defer crontab.Stop()

	r := router.SetupRouter(cfg.Server.Mode)

	addr := fmt.Sprintf(":%d", cfg.Server.Port)
	logger.Infof("Server starting on %s", addr)

	if err := r.Run(addr); err != nil {
		logger.Fatalf("Failed to start server: %v", err)
	}
}

func initStorage(cfg *config.StorageConfig) {
	switch cfg.Type {
	case "minio":
		s, err := storage.NewMinIOStorage(storage.MinIOConfig{
			Endpoint:  cfg.MinIO.Endpoint,
			AccessKey: cfg.MinIO.AccessKey,
			SecretKey: cfg.MinIO.SecretKey,
			Bucket:    cfg.MinIO.Bucket,
			UseSSL:    cfg.MinIO.UseSSL,
			BaseURL:   cfg.MinIO.BaseURL,
		})
		if err != nil {
			logger.Fatalf("Failed to init MinIO storage: %v", err)
		}
		storage.Init(s)
		logger.Infof("Storage: MinIO (%s/%s)", cfg.MinIO.Endpoint, cfg.MinIO.Bucket)
	default:
		s := storage.NewLocalStorage(cfg.Local.BasePath, cfg.Local.BaseURL)
		storage.Init(s)
		logger.Infof("Storage: Local (%s)", cfg.Local.BasePath)
	}
}
