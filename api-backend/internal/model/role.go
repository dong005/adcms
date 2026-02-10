package model

import "time"

type Role struct {
	TenantBaseModel
	Name        string       `gorm:"size:50;not null" json:"name"`
	Code        string       `gorm:"size:50;not null" json:"code"`
	Description string       `gorm:"size:255" json:"description"`
	Status      int8         `gorm:"default:1" json:"status"`
	Sort        int          `gorm:"default:0" json:"sort"`
	DataScope   int8         `gorm:"default:1" json:"data_scope"` // 1=全部 2=本部门及下级 3=本部门 4=仅自己
	Permissions []Permission `gorm:"many2many:role_permissions;" json:"permissions,omitempty"`
	Menus       []Menu       `gorm:"many2many:role_menus;" json:"menus,omitempty"`
}

func (Role) TableName() string {
	return "roles"
}

type RolePermission struct {
	RoleID       uint      `gorm:"primaryKey"`
	PermissionID uint      `gorm:"primaryKey"`
	CreatedAt    time.Time `json:"created_at"`
}

func (RolePermission) TableName() string {
	return "role_permissions"
}

type RoleMenu struct {
	RoleID    uint      `gorm:"primaryKey"`
	MenuID    uint      `gorm:"primaryKey"`
	CreatedAt time.Time `json:"created_at"`
}

func (RoleMenu) TableName() string {
	return "role_menus"
}
