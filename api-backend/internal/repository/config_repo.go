package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type ConfigRepository struct {
	db *gorm.DB
}

func NewConfigRepository() *ConfigRepository {
	return &ConfigRepository{db: database.DB}
}

func (r *ConfigRepository) FindAll(tenantID uint) ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := r.db.Where("tenant_id = 0 OR tenant_id = ?", tenantID).Find(&configs).Error
	return configs, err
}

func (r *ConfigRepository) GetByKey(tenantID uint, key string) (*model.SystemConfig, error) {
	var config model.SystemConfig
	err := r.db.Where("(tenant_id = 0 OR tenant_id = ?) AND `key` = ?", tenantID, key).First(&config).Error
	return &config, err
}

func (r *ConfigRepository) Upsert(config *model.SystemConfig) error {
	var existing model.SystemConfig
	err := r.db.Where("tenant_id = ? AND `key` = ?", config.TenantID, config.Key).First(&existing).Error
	if err == nil {
		existing.Value = config.Value
		existing.Description = config.Description
		return r.db.Save(&existing).Error
	}
	return r.db.Create(config).Error
}
