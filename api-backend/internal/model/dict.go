package model

type DictType struct {
	TenantBaseModel
	Name   string `gorm:"size:100;not null" json:"name"`
	Code   string `gorm:"size:100;not null;uniqueIndex" json:"code"`
	Sort   int    `gorm:"default:0" json:"sort"`
	Status int8   `gorm:"default:1" json:"status"`
	Remark string `gorm:"size:500" json:"remark"`
}

type Dict struct {
	TenantBaseModel
	DictTypeID uint   `gorm:"index;not null" json:"dict_type_id"`
	Name       string `gorm:"size:100;not null" json:"name"`
	Value      string `gorm:"size:200;not null" json:"value"`
	Sort       int    `gorm:"default:0" json:"sort"`
	Status     int8   `gorm:"default:1" json:"status"`
	Remark     string `gorm:"size:500" json:"remark"`
}
