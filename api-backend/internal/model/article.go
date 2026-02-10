package model

import "time"

type Article struct {
	TenantBaseModel
	CategoryID  uint       `gorm:"index" json:"category_id"`
	UserID      uint       `gorm:"index" json:"user_id"`
	Title       string     `gorm:"size:255;not null" json:"title"`
	Slug        string     `gorm:"size:255" json:"slug"`
	Summary     string     `gorm:"size:500" json:"summary"`
	Content     string     `gorm:"type:longtext" json:"content"`
	Cover       string     `gorm:"size:255" json:"cover"`
	Status      int8       `gorm:"default:0" json:"status"`
	ViewCount   int        `gorm:"default:0" json:"view_count"`
	PublishedAt *time.Time `json:"published_at"`
	Tags        []Tag      `gorm:"many2many:article_tags;" json:"tags,omitempty"`
}

func (Article) TableName() string {
	return "articles"
}

type ArticleTag struct {
	ArticleID uint `gorm:"primaryKey"`
	TagID     uint `gorm:"primaryKey"`
}

func (ArticleTag) TableName() string {
	return "article_tags"
}
