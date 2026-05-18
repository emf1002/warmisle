package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

type AuthRepo struct{}

func (r *AuthRepo) FindByUsername(username string) (*model.Member, error) {
	var m model.Member
	err := pkg.DB.Where("username = ?", username).First(&m).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *AuthRepo) FindByID(id uint) (*model.Member, error) {
	var m model.Member
	err := pkg.DB.First(&m, id).Error
	if err != nil {
		return nil, err
	}
	return &m, nil
}

func (r *AuthRepo) Count() (int64, error) {
	var count int64
	err := pkg.DB.Model(&model.Member{}).Count(&count).Error
	return count, err
}

func (r *AuthRepo) Create(member *model.Member) error {
	return pkg.DB.Create(member).Error
}

func (r *AuthRepo) UpdatePassword(id uint, hash string) error {
	return pkg.DB.Model(&model.Member{}).Where("id = ?", id).Update("password", hash).Error
}
