package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type TagRepository struct {
	db *gorm.DB
}

func NewTagRepository() *TagRepository {
	return &TagRepository{db: database.DB}
}

func (r *TagRepository) Create(tag *model.Tag) error {
	return r.db.Create(tag).Error
}

func (r *TagRepository) Update(tag *model.Tag) error {
	return r.db.Save(tag).Error
}

func (r *TagRepository) Delete(id uint) error {
	return r.db.Delete(&model.Tag{}, id).Error
}

func (r *TagRepository) FindByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := r.db.First(&tag, id).Error
	return &tag, err
}

func (r *TagRepository) List(tenantID uint) ([]model.Tag, error) {
	var tags []model.Tag
	query := r.db.Order("id DESC")
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Find(&tags).Error
	return tags, err
}
