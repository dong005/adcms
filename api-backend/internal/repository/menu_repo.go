package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type MenuRepository struct {
	db *gorm.DB
}

func NewMenuRepository() *MenuRepository {
	return &MenuRepository{db: database.DB}
}

func (r *MenuRepository) Create(menu *model.Menu) error {
	return r.db.Create(menu).Error
}

func (r *MenuRepository) Update(menu *model.Menu) error {
	return r.db.Save(menu).Error
}

func (r *MenuRepository) Delete(id uint) error {
	return r.db.Delete(&model.Menu{}, id).Error
}

func (r *MenuRepository) FindByID(id uint) (*model.Menu, error) {
	var menu model.Menu
	err := r.db.First(&menu, id).Error
	return &menu, err
}

func (r *MenuRepository) FindAll(tenantID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("tenant_id = 0 OR tenant_id = ?", tenantID).Order("sort ASC, id ASC").Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) FindByRoleIDs(roleIDs []uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Distinct().
		Joins("JOIN role_menus ON role_menus.menu_id = menus.id").
		Where("role_menus.role_id IN ?", roleIDs).
		Where("menus.status = 1").
		Order("menus.sort ASC, menus.id ASC").
		Find(&menus).Error
	if err != nil {
		return menus, err
	}

	// 防御性补全：自动补充缺失的父菜单，确保菜单树完整
	menus = r.completeParentMenus(menus)
	return menus, nil
}

// completeParentMenus 补全菜单列表中缺失的父菜单
func (r *MenuRepository) completeParentMenus(menus []model.Menu) []model.Menu {
	existIDs := make(map[uint]bool)
	var missingParentIDs []uint
	for _, m := range menus {
		existIDs[m.ID] = true
	}
	for _, m := range menus {
		if m.ParentID != 0 && !existIDs[m.ParentID] {
			missingParentIDs = append(missingParentIDs, m.ParentID)
			existIDs[m.ParentID] = true
		}
	}
	if len(missingParentIDs) == 0 {
		return menus
	}
	// 递归查找所有缺失的祖先菜单
	var parentMenus []model.Menu
	r.db.Where("id IN ? AND status = 1", missingParentIDs).Find(&parentMenus)
	menus = append(menus, parentMenus...)
	// 递归补全（处理多层嵌套的情况）
	return r.completeParentMenus(menus)
}

func (r *MenuRepository) FindChildren(parentID uint) ([]model.Menu, error) {
	var menus []model.Menu
	err := r.db.Where("parent_id = ?", parentID).Find(&menus).Error
	return menus, err
}

func (r *MenuRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Menu{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func BuildMenuTree(menus []model.Menu, parentID uint) []model.Menu {
	var tree []model.Menu
	for _, menu := range menus {
		if menu.ParentID == parentID {
			menu.Children = BuildMenuTree(menus, menu.ID)
			tree = append(tree, menu)
		}
	}
	return tree
}

func ConvertToMenuTree(menus []model.Menu) []model.MenuTree {
	menuMap := make(map[uint]*model.MenuTree)
	var roots []model.MenuTree

	for _, menu := range menus {
		// 将 Title 转换为 i18n key 格式（menu.{Name}），让前端能正确翻译
		i18nKey := "menu." + menu.Name

		menuTree := model.MenuTree{
			ID:        menu.ID,
			ParentID:  menu.ParentID,
			Name:      menu.Name,
			Path:      menu.Path,
			Component: menu.Component,
			Redirect:  menu.Redirect,
			Meta: model.MenuMeta{
				Title:            i18nKey,
				Icon:             menu.Icon,
				HideInMenu:       menu.HideInMenu == 1,
				HideInTab:        menu.HideInTab == 1,
				HideInBreadcrumb: menu.HideInBreadcrumb == 1,
				KeepAlive:        menu.KeepAlive == 1,
				FrameSrc:         menu.FrameSrc,
			},
		}
		menuMap[menu.ID] = &menuTree
	}

	for _, menu := range menus {
		if menu.ParentID == 0 {
			roots = append(roots, *menuMap[menu.ID])
		} else {
			if parent, ok := menuMap[menu.ParentID]; ok {
				parent.Children = append(parent.Children, *menuMap[menu.ID])
				menuMap[menu.ParentID] = parent
			}
		}
	}

	var result []model.MenuTree
	for _, root := range roots {
		result = append(result, buildTreeWithChildren(menuMap, root.ID))
	}

	return result
}

func buildTreeWithChildren(menuMap map[uint]*model.MenuTree, id uint) model.MenuTree {
	menu := *menuMap[id]
	var children []model.MenuTree
	for _, child := range menu.Children {
		children = append(children, buildTreeWithChildren(menuMap, child.ID))
	}
	menu.Children = children
	return menu
}
