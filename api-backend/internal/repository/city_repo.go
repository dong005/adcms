package repository

import (
	"adcms/internal/model"
	"adcms/pkg/database"

	"gorm.io/gorm"
)

type CityRepository struct {
	db *gorm.DB
}

func NewCityRepository() *CityRepository {
	return &CityRepository{db: database.DB}
}

func (r *CityRepository) ListByPID(pid uint) ([]model.City, error) {
	var cities []model.City
	err := r.db.Where("pid = ?", pid).Order("sort ASC, id ASC").Find(&cities).Error
	return cities, err
}

func (r *CityRepository) FindByID(id uint) (*model.City, error) {
	var c model.City
	err := r.db.First(&c, id).Error
	return &c, err
}

func (r *CityRepository) Create(c *model.City) error {
	return r.db.Create(c).Error
}

func (r *CityRepository) Update(c *model.City) error {
	return r.db.Save(c).Error
}

func (r *CityRepository) Delete(id uint) error {
	// 递归删除子节点
	var children []model.City
	r.db.Where("pid = ?", id).Find(&children)
	for _, child := range children {
		r.Delete(child.ID)
	}
	return r.db.Delete(&model.City{}, id).Error
}

type CityTree struct {
	model.City
	Children []CityTree `json:"children,omitempty"`
}

func (r *CityRepository) Tree(pid uint, maxLevel int8) ([]CityTree, error) {
	var cities []model.City
	err := r.db.Where("pid = ?", pid).Order("sort ASC, id ASC").Find(&cities).Error
	if err != nil {
		return nil, err
	}
	var tree []CityTree
	for _, c := range cities {
		node := CityTree{City: c}
		if c.Level < maxLevel {
			children, _ := r.Tree(c.ID, maxLevel)
			node.Children = children
		}
		tree = append(tree, node)
	}
	return tree, nil
}
