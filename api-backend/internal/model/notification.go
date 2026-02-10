package model

import "time"

// Notification 站内消息表
type Notification struct {
	TenantBaseModel
	SenderID   uint       `gorm:"default:0" json:"sender_id"`   // 0=系统消息
	ReceiverID uint       `gorm:"index;not null" json:"receiver_id"`
	Title      string     `gorm:"size:200;not null" json:"title"`
	Content    string     `gorm:"type:text" json:"content"`
	Type       string     `gorm:"size:50;default:system" json:"type"` // system/task/message
	IsRead     int8       `gorm:"default:0" json:"is_read"`          // 0=未读 1=已读
	ReadAt     *time.Time `json:"read_at"`
	Extra      string     `gorm:"type:text" json:"extra"`      // JSON扩展字段
	ReplyToID  uint       `gorm:"default:0" json:"reply_to_id"` // 回复哪条消息，0=原始消息
	SenderName string     `gorm:"-" json:"sender_name"`         // 发送者名称（非DB字段）
}

func (Notification) TableName() string {
	return "notifications"
}
