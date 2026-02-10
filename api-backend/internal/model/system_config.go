package model

type SystemConfig struct {
	BaseModel
	TenantID    uint   `gorm:"index;default:0" json:"tenant_id"`
	Key         string `gorm:"size:100;not null" json:"key"`
	Value       string `gorm:"type:text" json:"value"`
	Description string `gorm:"size:255" json:"description"`
}

func (SystemConfig) TableName() string {
	return "system_configs"
}
