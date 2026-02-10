package model

type Permission struct {
	BaseModel
	Name        string `gorm:"size:50;not null" json:"name"`
	Code        string `gorm:"size:100;uniqueIndex;not null" json:"code"`
	Type        int8   `gorm:"default:1" json:"type"`
	ParentID    uint   `gorm:"default:0" json:"parent_id"`
	Path        string `gorm:"size:255" json:"path"`
	Method      string `gorm:"size:10" json:"method"`
	Description string `gorm:"size:255" json:"description"`
}

func (Permission) TableName() string {
	return "permissions"
}
