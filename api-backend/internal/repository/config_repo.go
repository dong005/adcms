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
	err := r.db.Where("tenant_id = 0 OR tenant_id = ?", tenantID).Order("sort ASC, id ASC").Find(&configs).Error
	return configs, err
}

func (r *ConfigRepository) FindByGroup(tenantID uint, groupID uint) ([]model.SystemConfig, error) {
	var configs []model.SystemConfig
	err := r.db.Where("(tenant_id = 0 OR tenant_id = ?) AND group_id = ?", tenantID, groupID).Order("sort ASC, id ASC").Find(&configs).Error
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

// ========== ConfigGroup ==========

func (r *ConfigRepository) ListGroups() ([]model.ConfigGroup, error) {
	var groups []model.ConfigGroup
	err := r.db.Order("sort ASC, id ASC").Find(&groups).Error
	return groups, err
}

func (r *ConfigRepository) FindGroupByID(id uint) (*model.ConfigGroup, error) {
	var g model.ConfigGroup
	err := r.db.First(&g, id).Error
	return &g, err
}

func (r *ConfigRepository) CreateGroup(g *model.ConfigGroup) error {
	return r.db.Create(g).Error
}

func (r *ConfigRepository) UpdateGroup(g *model.ConfigGroup) error {
	return r.db.Save(g).Error
}

func (r *ConfigRepository) DeleteGroup(id uint) error {
	return r.db.Delete(&model.ConfigGroup{}, id).Error
}

// ========== ConfigWeb ==========

func (r *ConfigRepository) ListWebs(tenantID uint) ([]model.ConfigWeb, error) {
	var webs []model.ConfigWeb
	query := r.db.Model(&model.ConfigWeb{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Order("sort ASC, id ASC").Find(&webs).Error
	return webs, err
}

func (r *ConfigRepository) FindWebByID(id uint) (*model.ConfigWeb, error) {
	var w model.ConfigWeb
	err := r.db.First(&w, id).Error
	return &w, err
}

func (r *ConfigRepository) UpsertWeb(web *model.ConfigWeb) error {
	var existing model.ConfigWeb
	err := r.db.Where("tenant_id = ? AND code = ?", web.TenantID, web.Code).First(&existing).Error
	if err == nil {
		existing.Name = web.Name
		existing.Value = web.Value
		existing.Sort = web.Sort
		return r.db.Save(&existing).Error
	}
	return r.db.Create(web).Error
}

func (r *ConfigRepository) DeleteWeb(id uint) error {
	return r.db.Delete(&model.ConfigWeb{}, id).Error
}
