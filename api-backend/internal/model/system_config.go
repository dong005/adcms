package model

type ConfigGroup struct {
	BaseModel
	Name string `gorm:"size:100;not null" json:"name"`
	Sort int    `gorm:"default:0" json:"sort"`
}

func (ConfigGroup) TableName() string {
	return "config_groups"
}

type SystemConfig struct {
	BaseModel
	TenantID    uint   `gorm:"index;default:0" json:"tenant_id"`
	GroupID     uint   `gorm:"default:0" json:"group_id"`
	Key         string `gorm:"size:100;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Type        string `gorm:"size:20;default:text" json:"type"`
	Options     string `gorm:"type:text" json:"options"`
	Description string `gorm:"size:255" json:"description"`
	Sort        int    `gorm:"default:0" json:"sort"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}

type ConfigWeb struct {
	TenantBaseModel
	Name  string `gorm:"size:100;not null" json:"name"`
	Code  string `gorm:"size:100;not null" json:"code"`
	Value string `gorm:"type:text" json:"value"`
	Sort  int    `gorm:"default:0" json:"sort"`
}

func (ConfigWeb) TableName() string {
	return "config_webs"
}
