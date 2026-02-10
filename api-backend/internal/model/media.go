package model

type Media struct {
	TenantBaseModel
	UserID   uint   `gorm:"index" json:"user_id"`
	Name     string `gorm:"size:255;not null" json:"name"`
	Path     string `gorm:"size:255;not null" json:"path"`
	Type     string `gorm:"size:50" json:"type"`
	Size     int64  `json:"size"`
	MimeType string `gorm:"size:100" json:"mime_type"`
}

func (Media) TableName() string {
	return "media"
}
