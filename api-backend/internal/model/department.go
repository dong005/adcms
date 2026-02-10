package model

// Department 部门表
type Department struct {
	TenantBaseModel
	ParentID    uint         `gorm:"default:0" json:"parent_id"`
	Name        string       `gorm:"size:100;not null" json:"name"`
	Code        string       `gorm:"size:50;index" json:"code"`
	Leader      string       `gorm:"size:50" json:"leader"`
	Phone       string       `gorm:"size:20" json:"phone"`
	Email       string       `gorm:"size:100" json:"email"`
	Sort        int          `gorm:"default:0" json:"sort"`
	Status      int8         `gorm:"default:1" json:"status"`
	Children    []Department `gorm:"-" json:"children,omitempty"`
}

func (Department) TableName() string {
	return "departments"
}

// UserDepartment 用户部门关联表
type UserDepartment struct {
	UserID       uint `gorm:"primaryKey" json:"user_id"`
	DepartmentID uint `gorm:"primaryKey" json:"department_id"`
}

func (UserDepartment) TableName() string {
	return "user_departments"
}

// DataScope 数据权限范围常量
const (
	DataScopeAll      = 1 // 全部数据
	DataScopeDeptTree = 2 // 本部门及下级
	DataScopeDept     = 3 // 本部门
	DataScopeSelf     = 4 // 仅自己
)
