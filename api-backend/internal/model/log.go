package model

import "time"

type OperationLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TenantID  uint      `gorm:"index" json:"tenant_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Module    string    `gorm:"size:50" json:"module"`
	Action    string    `gorm:"size:50" json:"action"`
	Method    string    `gorm:"size:10" json:"method"`
	Path      string    `gorm:"size:255" json:"path"`
	Params    string    `gorm:"type:text" json:"params"`
	Response  string    `gorm:"type:text" json:"response"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Duration  int64     `json:"duration"`
	CreatedAt time.Time `json:"created_at"`
}

func (OperationLog) TableName() string {
	return "operation_logs"
}

type LoginLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TenantID  uint      `gorm:"index" json:"tenant_id"`
	UserID    uint      `gorm:"index" json:"user_id"`
	Username  string    `gorm:"size:50" json:"username"`
	IP        string    `gorm:"size:45" json:"ip"`
	UserAgent string    `gorm:"size:500" json:"user_agent"`
	Status    int8      `json:"status"`
	Message   string    `gorm:"size:255" json:"message"`
	CreatedAt time.Time `json:"created_at"`
}

func (LoginLog) TableName() string {
	return "login_logs"
}

type EmailLog struct {
	ID        uint      `gorm:"primarykey" json:"id"`
	TenantID  uint      `gorm:"index" json:"tenant_id"`
	To        string    `gorm:"size:255" json:"to"`
	Subject   string    `gorm:"size:500" json:"subject"`
	Status    int8      `json:"status"` // 1=成功 0=失败
	Error     string    `gorm:"size:500" json:"error"`
	CreatedAt time.Time `json:"created_at"`
}

func (EmailLog) TableName() string {
	return "email_logs"
}

type SmsLog struct {
	ID         uint      `gorm:"primarykey" json:"id"`
	TenantID   uint      `gorm:"index" json:"tenant_id"`
	Phone      string    `gorm:"size:20" json:"phone"`
	TemplateID string    `gorm:"size:50" json:"template_id"`
	Params     string    `gorm:"size:500" json:"params"`
	Status     int8      `json:"status"` // 1=成功 0=失败
	Error      string    `gorm:"size:500" json:"error"`
	CreatedAt  time.Time `json:"created_at"`
}

func (SmsLog) TableName() string {
	return "sms_logs"
}
