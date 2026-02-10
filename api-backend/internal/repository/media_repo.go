package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type MediaRepository struct {
	db *gorm.DB
}

func NewMediaRepository() *MediaRepository {
	return &MediaRepository{db: database.DB}
}

func (r *MediaRepository) Create(media *model.Media) error {
	return r.db.Create(media).Error
}

func (r *MediaRepository) Delete(id uint) error {
	return r.db.Delete(&model.Media{}, id).Error
}

func (r *MediaRepository) FindByID(id uint) (*model.Media, error) {
	var media model.Media
	err := r.db.First(&media, id).Error
	return &media, err
}

func (r *MediaRepository) List(tenantID uint, page, pageSize int, mediaType string) ([]model.Media, int64, error) {
	var medias []model.Media
	var total int64

	query := r.db.Model(&model.Media{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if mediaType != "" {
		query = query.Where("type = ?", mediaType)
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Offset(offset).Limit(pageSize).Order("id DESC").Find(&medias).Error
	return medias, total, err
}
