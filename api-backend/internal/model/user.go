package model

import "time"

type User struct {
	TenantBaseModel
	Username    string     `gorm:"size:50;index;not null" json:"username"`
	Password    string     `gorm:"size:255;not null" json:"-"`
	Email       string     `gorm:"size:100;index" json:"email"`
	Phone       string     `gorm:"size:20" json:"phone"`
	Nickname    string     `gorm:"size:50" json:"nickname"`
	Avatar      string     `gorm:"size:255" json:"avatar"`
	DepartmentID uint       `gorm:"default:0;index" json:"department_id"`
	Status       int8       `gorm:"default:1" json:"status"`
	TOTPEnabled  int8       `gorm:"default:0" json:"totp_enabled"`
	TOTPSecret  string     `gorm:"size:32" json:"-"`
	EmailNotify int8       `gorm:"default:1" json:"email_notify"` // 1=接收系统邮件通知 0=不接收
	LastLoginAt *time.Time `json:"last_login_at"`
	LastLoginIP string     `gorm:"size:45" json:"last_login_ip"`
	Roles       []Role     `gorm:"many2many:user_roles;" json:"roles,omitempty"`
}

func (User) TableName() string {
	return "users"
}

type UserRole struct {
	UserID    uint      `gorm:"primaryKey"`
	RoleID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserRole) TableName() string {
	return "user_roles"
}
