package model

type Site struct {
	TenantBaseModel
	Name     string `gorm:"size:200;not null" json:"name"`
	Type     string `gorm:"size:50" json:"type"`
	URL      string `gorm:"size:500" json:"url"`
	Image    string `gorm:"size:500" json:"image"`
	IsDomain int8   `gorm:"default:0" json:"is_domain"`
	Status   int8   `gorm:"default:1" json:"status"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Remark   string `gorm:"size:500" json:"remark"`
}
