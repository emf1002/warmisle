package service

import (
	"errors"

	"home-center/internal/model"
	"home-center/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrCategoryNotFound = errors.New("category not found")
	ErrCategoryExists   = errors.New("category with same name and type already exists")
	ErrCategoryInUse    = errors.New("category is in use by ledger records")
)

type CategoryService struct {
	repo *repository.CategoryRepo
}

func NewCategoryService() *CategoryService {
	return &CategoryService{repo: &repository.CategoryRepo{}}
}

func (s *CategoryService) List() ([]model.Category, error) {
	return s.repo.List()
}

func (s *CategoryService) Create(typ, name, icon string, sortOrder int) (*model.Category, error) {
	// Check uniqueness: same type + same name
	existing, err := s.repo.FindByNameAndType(name, typ)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	}
	if existing != nil {
		return nil, ErrCategoryExists
	}

	c := &model.Category{
		Type:      typ,
		Name:      name,
		Icon:      icon,
		SortOrder: sortOrder,
		Preset:    false,
	}
	if err := s.repo.Create(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Update(id uint, typ, name, icon string, sortOrder *int) (*model.Category, error) {
	c, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrCategoryNotFound
		}
		return nil, err
	}

	// Check uniqueness if name or type changed
	if (typ != "" && typ != c.Type) || (name != "" && name != c.Name) {
		checkType := c.Type
		if typ != "" {
			checkType = typ
		}
		checkName := c.Name
		if name != "" {
			checkName = name
		}
		existing, err := s.repo.FindByNameAndType(checkName, checkType)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, err
		}
		if existing != nil && existing.ID != id {
			return nil, ErrCategoryExists
		}
	}

	if typ != "" {
		c.Type = typ
	}
	if name != "" {
		c.Name = name
	}
	if icon != "" {
		c.Icon = icon
	}
	if sortOrder != nil {
		c.SortOrder = *sortOrder
	}

	if err := s.repo.Update(c); err != nil {
		return nil, err
	}
	return c, nil
}

func (s *CategoryService) Delete(id uint) error {
	count, err := s.repo.SoftDelete(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrCategoryInUse
	}
	return nil
}
