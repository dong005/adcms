package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type DictRepository struct {
	db *gorm.DB
}

func NewDictRepository() *DictRepository {
	return &DictRepository{db: database.DB}
}

// ========== DictType ==========

func (r *DictRepository) ListTypes(tenantID uint, keyword string) ([]model.DictType, error) {
	var types []model.DictType
	query := r.db.Model(&model.DictType{})
	if tenantID > 0 {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if keyword != "" {
		query = query.Where("name LIKE ? OR code LIKE ?", "%"+keyword+"%", "%"+keyword+"%")
	}
	err := query.Order("sort ASC, id ASC").Find(&types).Error
	return types, err
}

func (r *DictRepository) FindTypeByID(id uint) (*model.DictType, error) {
	var t model.DictType
	err := r.db.First(&t, id).Error
	return &t, err
}

func (r *DictRepository) FindTypeByCode(code string) (*model.DictType, error) {
	var t model.DictType
	err := r.db.Where("code = ?", code).First(&t).Error
	return &t, err
}

func (r *DictRepository) CreateType(t *model.DictType) error {
	return r.db.Create(t).Error
}

func (r *DictRepository) UpdateType(t *model.DictType) error {
	return r.db.Save(t).Error
}

func (r *DictRepository) DeleteType(id uint) error {
	// 先删除该类型下的所有字典数据
	if err := r.db.Where("dict_type_id = ?", id).Delete(&model.Dict{}).Error; err != nil {
		return err
	}
	return r.db.Delete(&model.DictType{}, id).Error
}

// ========== Dict ==========

func (r *DictRepository) ListDicts(typeID uint) ([]model.Dict, error) {
	var dicts []model.Dict
	err := r.db.Where("dict_type_id = ?", typeID).Order("sort ASC, id ASC").Find(&dicts).Error
	return dicts, err
}

func (r *DictRepository) FindDictByID(id uint) (*model.Dict, error) {
	var d model.Dict
	err := r.db.First(&d, id).Error
	return &d, err
}

func (r *DictRepository) CreateDict(d *model.Dict) error {
	return r.db.Create(d).Error
}

func (r *DictRepository) UpdateDict(d *model.Dict) error {
	return r.db.Save(d).Error
}

func (r *DictRepository) DeleteDict(id uint) error {
	return r.db.Delete(&model.Dict{}, id).Error
}

// GetDictsByCode 根据字典类型编码获取字典数据（前端下拉用）
func (r *DictRepository) GetDictsByCode(code string) ([]model.Dict, error) {
	var dictType model.DictType
	if err := r.db.Where("code = ? AND status = 1", code).First(&dictType).Error; err != nil {
		return nil, err
	}
	var dicts []model.Dict
	err := r.db.Where("dict_type_id = ? AND status = 1", dictType.ID).Order("sort ASC, id ASC").Find(&dicts).Error
	return dicts, err
}
