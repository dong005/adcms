package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type RoleRepository struct {
	db *gorm.DB
}

func NewRoleRepository() *RoleRepository {
	return &RoleRepository{db: database.DB}
}

func (r *RoleRepository) Create(role *model.Role) error {
	return r.db.Create(role).Error
}

func (r *RoleRepository) Update(role *model.Role) error {
	return r.db.Save(role).Error
}

func (r *RoleRepository) Delete(id uint) error {
	return r.db.Delete(&model.Role{}, id).Error
}

func (r *RoleRepository) FindByID(id uint) (*model.Role, error) {
	var role model.Role
	err := r.db.Preload("Permissions").First(&role, id).Error
	return &role, err
}

func (r *RoleRepository) FindByCode(tenantID uint, code string) (*model.Role, error) {
	var role model.Role
	err := r.db.Where("tenant_id = ? AND code = ?", tenantID, code).First(&role).Error
	return &role, err
}

func (r *RoleRepository) List(tenantID uint) ([]model.Role, error) {
	var roles []model.Role
	query := r.db.Order("sort ASC, id ASC")
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Find(&roles).Error
	return roles, err
}

func (r *RoleRepository) AssignMenus(roleID uint, menuIDs []uint) error {
	// 自动补全所有父菜单，防止前端遗漏父节点导致菜单树断裂
	allIDs := r.completeParentMenuIDs(menuIDs)

	if err := r.db.Where("role_id = ?", roleID).Delete(&model.RoleMenu{}).Error; err != nil {
		return err
	}

	for _, menuID := range allIDs {
		roleMenu := model.RoleMenu{RoleID: roleID, MenuID: menuID}
		if err := r.db.Create(&roleMenu).Error; err != nil {
			return err
		}
	}
	return nil
}

// completeParentMenuIDs 向上追溯补全所有父菜单 ID
func (r *RoleRepository) completeParentMenuIDs(menuIDs []uint) []uint {
	idSet := make(map[uint]bool)
	for _, id := range menuIDs {
		idSet[id] = true
	}

	// 查询所有菜单的 id -> parent_id 映射
	var menus []model.Menu
	r.db.Select("id, parent_id").Find(&menus)
	parentMap := make(map[uint]uint)
	for _, m := range menus {
		parentMap[m.ID] = m.ParentID
	}

	// 对每个 menuID 向上追溯，补全所有祖先
	for _, id := range menuIDs {
		current := id
		for {
			pid, ok := parentMap[current]
			if !ok || pid == 0 {
				break
			}
			idSet[pid] = true
			current = pid
		}
	}

	result := make([]uint, 0, len(idSet))
	for id := range idSet {
		result = append(result, id)
	}
	return result
}

func (r *RoleRepository) GetRoleMenus(roleID uint) ([]model.Menu, error) {
	var role model.Role
	err := r.db.Preload("Menus").First(&role, roleID).Error
	return role.Menus, err
}

func (r *RoleRepository) AssignPermissions(roleID uint, permissionIDs []uint) error {
	if err := r.db.Where("role_id = ?", roleID).Delete(&model.RolePermission{}).Error; err != nil {
		return err
	}

	for _, permissionID := range permissionIDs {
		rp := model.RolePermission{RoleID: roleID, PermissionID: permissionID}
		if err := r.db.Create(&rp).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *RoleRepository) GetRolePermissions(roleID uint) ([]model.Permission, error) {
	var role model.Role
	err := r.db.Preload("Permissions").First(&role, roleID).Error
	return role.Permissions, err
}
