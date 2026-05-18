package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

type CategoryRepo struct{}

func (r *CategoryRepo) List() ([]model.Category, error) {
	var list []model.Category
	err := pkg.DB.Order("type, sort_order").Find(&list).Error
	return list, err
}

func (r *CategoryRepo) FindByID(id uint) (*model.Category, error) {
	var c model.Category
	err := pkg.DB.First(&c, id).Error
	return &c, err
}

func (r *CategoryRepo) Create(c *model.Category) error {
	return pkg.DB.Create(c).Error
}

func (r *CategoryRepo) Update(c *model.Category) error {
	return pkg.DB.Save(c).Error
}

// SoftDelete returns count of associated ledger records; if >0, refuse deletion
func (r *CategoryRepo) SoftDelete(id uint) (int64, error) {
	var count int64
	pkg.DB.Model(&model.Ledger{}).Where("category_id = ?", id).Count(&count)
	if count > 0 {
		return count, nil
	}
	return 0, pkg.DB.Delete(&model.Category{}, id).Error
}

func (r *CategoryRepo) FindByNameAndType(name, typ string) (*model.Category, error) {
	var c model.Category
	err := pkg.DB.Where("name = ? AND type = ?", name, typ).First(&c).Error
	if err != nil {
		return nil, err
	}
	return &c, nil
}
