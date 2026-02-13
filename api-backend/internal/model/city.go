package model

type City struct {
	BaseModel
	PID      uint    `gorm:"column:pid;default:0;index" json:"pid"`
	Level    int8    `gorm:"default:0" json:"level"`
	Name     string  `gorm:"size:100;not null" json:"name"`
	Citycode string  `gorm:"size:20" json:"citycode"`
	Adcode   string  `gorm:"size:20" json:"adcode"`
	PAdcode  string  `gorm:"size:20" json:"p_adcode"`
	Lng      float64 `gorm:"type:decimal(10,6)" json:"lng"`
	Lat      float64 `gorm:"type:decimal(10,6)" json:"lat"`
	Sort     int     `gorm:"default:0" json:"sort"`
}

func (City) TableName() string {
	return "cities"
}
