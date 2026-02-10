package database

import (
	"adcms/internal/config"
	"fmt"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitMySQL(cfg *config.MySQLConfig) error {
	var err error

	gormConfig := &gorm.Config{
		Logger: logger.Default.LogMode(logger.Info),
	}

	DB, err = gorm.Open(mysql.Open(cfg.DSN()), gormConfig)
	if err != nil {
		return fmt.Errorf("failed to connect to mysql: %w", err)
	}

	sqlDB, err := DB.DB()
	if err != nil {
		return fmt.Errorf("failed to get sql.DB: %w", err)
	}

	sqlDB.SetMaxIdleConns(cfg.MaxIdleConns)
	sqlDB.SetMaxOpenConns(cfg.MaxOpenConns)
	sqlDB.SetConnMaxLifetime(time.Hour)

	return nil
}

func CloseMySQL() error {
	if DB != nil {
		sqlDB, err := DB.DB()
		if err != nil {
			return err
		}
		return sqlDB.Close()
	}
	return nil
}
