package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type LinkRepository struct {
	db *gorm.DB
}

func NewLinkRepository() *LinkRepository {
	return &LinkRepository{db: database.DB}
}

func (r *LinkRepository) List(tenantID uint, keyword string) ([]model.Link, error) {
	var links []model.Link
	query := r.db.Model(&model.Link{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ?", "%"+keyword+"%")
	}
	err := query.Order("sort ASC, id DESC").Find(&links).Error
	return links, err
}

func (r *LinkRepository) FindByID(id uint) (*model.Link, error) {
	var l model.Link
	err := r.db.First(&l, id).Error
	return &l, err
}

func (r *LinkRepository) Create(l *model.Link) error {
	return r.db.Create(l).Error
}

func (r *LinkRepository) Update(l *model.Link) error {
	return r.db.Save(l).Error
}

func (r *LinkRepository) Delete(id uint) error {
	return r.db.Delete(&model.Link{}, id).Error
}
