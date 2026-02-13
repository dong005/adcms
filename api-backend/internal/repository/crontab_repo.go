package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type CrontabRepository struct {
	db *gorm.DB
}

func NewCrontabRepository() *CrontabRepository {
	return &CrontabRepository{db: database.DB}
}

func (r *CrontabRepository) List(tenantID uint) ([]model.Crontab, error) {
	var crontabs []model.Crontab
	query := r.db.Model(&model.Crontab{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Order("id DESC").Find(&crontabs).Error
	return crontabs, err
}

func (r *CrontabRepository) FindByID(id uint) (*model.Crontab, error) {
	var c model.Crontab
	err := r.db.First(&c, id).Error
	return &c, err
}

func (r *CrontabRepository) Create(c *model.Crontab) error {
	return r.db.Create(c).Error
}

func (r *CrontabRepository) Update(c *model.Crontab) error {
	return r.db.Save(c).Error
}

func (r *CrontabRepository) Delete(id uint) error {
	return r.db.Delete(&model.Crontab{}, id).Error
}
