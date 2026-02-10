package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type CategoryRepository struct {
	db *gorm.DB
}

func NewCategoryRepository() *CategoryRepository {
	return &CategoryRepository{db: database.DB}
}

func (r *CategoryRepository) Create(category *model.Category) error {
	return r.db.Create(category).Error
}

func (r *CategoryRepository) Update(category *model.Category) error {
	return r.db.Save(category).Error
}

func (r *CategoryRepository) Delete(id uint) error {
	return r.db.Delete(&model.Category{}, id).Error
}

func (r *CategoryRepository) FindByID(id uint) (*model.Category, error) {
	var category model.Category
	err := r.db.First(&category, id).Error
	return &category, err
}

func (r *CategoryRepository) FindAll(tenantID uint) ([]model.Category, error) {
	var categories []model.Category
	query := r.db.Order("sort ASC, id ASC")
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Find(&categories).Error
	return categories, err
}

func (r *CategoryRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Category{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}
