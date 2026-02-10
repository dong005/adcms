package model

type Tag struct {
	TenantBaseModel
	Name string `gorm:"size:50;not null" json:"name"`
	Slug string `gorm:"size:100" json:"slug"`
}

func (Tag) TableName() string {
	return "tags"
}
