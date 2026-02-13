package model

type Link struct {
	TenantBaseModel
	SiteID   uint   `gorm:"default:0" json:"site_id"`
	CateID   uint   `gorm:"default:0" json:"cate_id"`
	Type     int8   `gorm:"default:1" json:"type"`
	Platform int8   `gorm:"default:1" json:"platform"`
	Form     int8   `gorm:"default:1" json:"form"`
	Name     string `gorm:"size:200;not null" json:"name"`
	Image    string `gorm:"size:500" json:"image"`
	URL      string `gorm:"size:500;not null" json:"url"`
	Logo     string `gorm:"size:500" json:"logo"`
	Desc     string `gorm:"size:500" json:"desc"`
	Sort     int    `gorm:"default:0" json:"sort"`
	Status   int8   `gorm:"default:1" json:"status"`
}
