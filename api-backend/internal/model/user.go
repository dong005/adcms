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
	// 新增字段
	IsAdmin     int8       `gorm:"default:0" json:"is_admin"`      // 0=普通用户 1=管理员(租户) 2=超级管理员
	Company     string     `gorm:"size:200" json:"company"`        // 租户公司名称
	Domain      string     `gorm:"size:200" json:"domain"`         // 租户绑定域名
	ExpireTime  *time.Time `json:"expire_time"`                    // 租户到期时间
	MaxUsers    uint       `gorm:"default:0" json:"max_users"`     // 最大用户数，0=不限
	LoginCount  uint       `gorm:"default:0" json:"login_count"`   // 登录次数
	Remark      string     `gorm:"size:500" json:"remark"`         // 备注
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

type UserMenu struct {
	UserID    uint      `gorm:"primaryKey"`
	MenuID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

func (UserMenu) TableName() string {
	return "user_menus"
}

func (UserRole) TableName() string {
	return "user_roles"
}
