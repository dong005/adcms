package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"
	"fmt"

	"gorm.io/gorm"
)

type PermissionRepository struct {
	db *gorm.DB
}

func NewPermissionRepository() *PermissionRepository {
	return &PermissionRepository{db: database.DB}
}

func (r *PermissionRepository) FindAll() ([]model.Permission, error) {
	var permissions []model.Permission
	err := r.db.Order("id ASC").Find(&permissions).Error
	return permissions, err
}

func (r *PermissionRepository) FindByID(id uint) (*model.Permission, error) {
	var permission model.Permission
	err := r.db.First(&permission, id).Error
	return &permission, err
}

func (r *PermissionRepository) Create(perm *model.Permission) error {
	return r.db.Create(perm).Error
}

func (r *PermissionRepository) Update(perm *model.Permission) error {
	return r.db.Save(perm).Error
}

func (r *PermissionRepository) Delete(id uint) error {
	// 检查是否有子权限
	var count int64
	r.db.Model(&model.Permission{}).Where("parent_id = ?", id).Count(&count)
	if count > 0 {
		return fmt.Errorf("该权限下有子权限，请先删除子权限")
	}
	return r.db.Delete(&model.Permission{}, id).Error
}

func (r *PermissionRepository) HasChildren(id uint) bool {
	var count int64
	r.db.Model(&model.Permission{}).Where("parent_id = ?", id).Count(&count)
	return count > 0
}

func BuildPermissionTree(permissions []model.Permission, parentID uint) []map[string]interface{} {
	var tree []map[string]interface{}
	for _, p := range permissions {
		if p.ParentID == parentID {
			node := map[string]interface{}{
				"id":          p.ID,
				"name":        p.Name,
				"code":        p.Code,
				"type":        p.Type,
				"parent_id":   p.ParentID,
				"path":        p.Path,
				"method":      p.Method,
				"description": p.Description,
				"children":    BuildPermissionTree(permissions, p.ID),
			}
			tree = append(tree, node)
		}
	}
	return tree
}
