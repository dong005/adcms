package model

type Menu struct {
	BaseModel
	TenantID         uint   `gorm:"index;default:0" json:"tenant_id"`
	ParentID         uint   `gorm:"default:0" json:"parent_id"`
	Name             string `gorm:"size:50;not null" json:"name"`
	Path             string `gorm:"size:255" json:"path"`
	Component        string `gorm:"size:255" json:"component"`
	Redirect         string `gorm:"size:255" json:"redirect"`
	Icon             string `gorm:"size:50" json:"icon"`
	Title            string `gorm:"size:50" json:"title"`
	HideInMenu       int8   `gorm:"default:0" json:"hide_in_menu"`
	HideInTab        int8   `gorm:"default:0" json:"hide_in_tab"`
	HideInBreadcrumb int8   `gorm:"default:0" json:"hide_in_breadcrumb"`
	KeepAlive        int8   `gorm:"default:1" json:"keep_alive"`
	FrameSrc         string `gorm:"size:255" json:"frame_src"`
	Sort             int    `gorm:"default:0" json:"sort"`
	Status           int8   `gorm:"default:1" json:"status"`
	PermissionCode   string `gorm:"size:100" json:"permission_code"`
	Children         []Menu `gorm:"-" json:"children,omitempty"`
}

func (Menu) TableName() string {
	return "menus"
}

type MenuMeta struct {
	Title            string `json:"title"`
	Icon             string `json:"icon,omitempty"`
	HideInMenu       bool   `json:"hideInMenu,omitempty"`
	HideInTab        bool   `json:"hideInTab,omitempty"`
	HideInBreadcrumb bool   `json:"hideInBreadcrumb,omitempty"`
	KeepAlive        bool   `json:"keepAlive,omitempty"`
	FrameSrc         string `json:"frameSrc,omitempty"`
}

type MenuTree struct {
	ID        uint       `json:"id"`
	ParentID  uint       `json:"parentId"`
	Name      string     `json:"name"`
	Path      string     `json:"path"`
	Component string     `json:"component,omitempty"`
	Redirect  string     `json:"redirect,omitempty"`
	Meta      MenuMeta   `json:"meta"`
	Children  []MenuTree `json:"children,omitempty"`
}
