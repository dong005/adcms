package model

import "time"

type Crontab struct {
	TenantBaseModel
	Name       string     `gorm:"size:200;not null" json:"name"`
	Type       int8       `gorm:"default:0" json:"type"`
	Content    string     `gorm:"type:text" json:"content"`
	Expression string     `gorm:"size:100;not null" json:"expression"`
	Command    string     `gorm:"size:500;not null" json:"command"`
	Maximums   int        `gorm:"default:0" json:"maximums"`
	Executes   int        `gorm:"default:0" json:"executes"`
	Status     int8       `gorm:"default:0" json:"status"`
	Sort       int        `gorm:"default:0" json:"sort"`
	Remark     string     `gorm:"size:500" json:"remark"`
	LastRunAt  *time.Time `json:"last_run_at"`
	NextRunAt  *time.Time `json:"next_run_at"`
}
