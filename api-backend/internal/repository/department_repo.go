package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type DepartmentRepository struct {
	db *gorm.DB
}

func NewDepartmentRepository() *DepartmentRepository {
	return &DepartmentRepository{db: database.DB}
}

func (r *DepartmentRepository) Create(dept *model.Department) error {
	return r.db.Create(dept).Error
}

func (r *DepartmentRepository) Update(dept *model.Department) error {
	return r.db.Save(dept).Error
}

func (r *DepartmentRepository) Delete(id uint) error {
	return r.db.Delete(&model.Department{}, id).Error
}

func (r *DepartmentRepository) FindByID(id uint) (*model.Department, error) {
	var dept model.Department
	err := r.db.First(&dept, id).Error
	return &dept, err
}

func (r *DepartmentRepository) FindAll(tenantID uint) ([]model.Department, error) {
	var depts []model.Department
	query := r.db.Order("sort ASC, id ASC")
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	err := query.Find(&depts).Error
	return depts, err
}

func (r *DepartmentRepository) HasChildren(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Department{}).Where("parent_id = ?", id).Count(&count).Error
	return count > 0, err
}

func (r *DepartmentRepository) HasUsers(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.User{}).Where("department_id = ?", id).Count(&count).Error
	return count > 0, err
}

// GetChildDeptIDs 获取部门及所有下级部门ID（递归）
func (r *DepartmentRepository) GetChildDeptIDs(tenantID, deptID uint) []uint {
	allDepts, _ := r.FindAll(tenantID)
	ids := []uint{deptID}
	r.collectChildIDs(allDepts, deptID, &ids)
	return ids
}

func (r *DepartmentRepository) collectChildIDs(depts []model.Department, parentID uint, ids *[]uint) {
	for _, dept := range depts {
		if dept.ParentID == parentID {
			*ids = append(*ids, dept.ID)
			r.collectChildIDs(depts, dept.ID, ids)
		}
	}
}

// BuildTree 构建部门树
func BuildDepartmentTree(depts []model.Department, parentID uint) []model.Department {
	var tree []model.Department
	for _, dept := range depts {
		if dept.ParentID == parentID {
			dept.Children = BuildDepartmentTree(depts, dept.ID)
			tree = append(tree, dept)
		}
	}
	return tree
}
