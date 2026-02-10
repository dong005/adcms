package model

type Category struct {
	TenantBaseModel
	ParentID    uint   `gorm:"default:0" json:"parent_id"`
	Name        string `gorm:"size:50;not null" json:"name"`
	Slug        string `gorm:"size:100" json:"slug"`
	Description string `gorm:"size:255" json:"description"`
	Sort        int    `gorm:"default:0" json:"sort"`
	Status      int8   `gorm:"default:1" json:"status"`
}

func (Category) TableName() string {
	return "categories"
}
