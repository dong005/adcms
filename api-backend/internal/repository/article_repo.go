package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type ArticleRepository struct {
	db *gorm.DB
}

func NewArticleRepository() *ArticleRepository {
	return &ArticleRepository{db: database.DB}
}

func (r *ArticleRepository) Create(article *model.Article) error {
	return r.db.Create(article).Error
}

func (r *ArticleRepository) Update(article *model.Article) error {
	return r.db.Save(article).Error
}

func (r *ArticleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Article{}, id).Error
}

func (r *ArticleRepository) FindByID(id uint) (*model.Article, error) {
	var article model.Article
	err := r.db.Preload("Tags").First(&article, id).Error
	return &article, err
}

func (r *ArticleRepository) List(tenantID uint, page, pageSize int, categoryID uint, status int8, keyword string) ([]model.Article, int64, error) {
	var articles []model.Article
	var total int64

	query := r.db.Model(&model.Article{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if categoryID > 0 {
		query = query.Where("category_id = ?", categoryID)
	}
	if status >= 0 {
		query = query.Where("status = ?", status)
	}
	if keyword != "" {
		query = query.Where("title LIKE ?", "%"+keyword+"%")
	}

	query.Count(&total)

	offset := (page - 1) * pageSize
	err := query.Preload("Tags").Offset(offset).Limit(pageSize).Order("id DESC").Find(&articles).Error
	return articles, total, err
}

func (r *ArticleRepository) UpdateStatus(id uint, status int8) error {
	updates := map[string]interface{}{"status": status}
	if status == 1 {
		updates["published_at"] = gorm.Expr("NOW()")
	}
	return r.db.Model(&model.Article{}).Where("id = ?", id).Updates(updates).Error
}

func (r *ArticleRepository) AssignTags(articleID uint, tagIDs []uint) error {
	r.db.Where("article_id = ?", articleID).Delete(&model.ArticleTag{})
	for _, tagID := range tagIDs {
		r.db.Create(&model.ArticleTag{ArticleID: articleID, TagID: tagID})
	}
	return nil
}
