package service

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"gorm.io/gorm"
)

type InitService struct{}

func (s *InitService) Setup(tx *gorm.DB, adminName, username, password string) (*model.Member, error) {
	hash, err := pkg.HashPassword(password)
	if err != nil {
		return nil, err
	}
	admin := model.Member{
		Username: username,
		Password: hash,
		Name:     adminName,
		Avatar:   "👨",
		Role:     "admin",
		Status:   "active",
	}
	if err := tx.Create(&admin).Error; err != nil {
		return nil, err
	}
	return &admin, nil
}
