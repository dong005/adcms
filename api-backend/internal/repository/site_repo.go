package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type SiteRepository struct {
	db *gorm.DB
}

func NewSiteRepository() *SiteRepository {
	return &SiteRepository{db: database.DB}
}

func (r *SiteRepository) List(tenantID uint, keyword string) ([]model.Site, error) {
	var sites []model.Site
	query := r.db.Model(&model.Site{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("sort ASC, id DESC").Find(&sites).Error
	return sites, err
}

func (r *SiteRepository) FindByID(id uint) (*model.Site, error) {
	var s model.Site
	err := r.db.First(&s, id).Error
	return &s, err
}

func (r *SiteRepository) Create(s *model.Site) error {
	return r.db.Create(s).Error
}

func (r *SiteRepository) Update(s *model.Site) error {
	return r.db.Save(s).Error
}

func (r *SiteRepository) Delete(id uint) error {
	return r.db.Delete(&model.Site{}, id).Error
}
