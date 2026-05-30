package service

import (
	"errors"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/repository"
)

var (
	ErrLedgerNotFound         = errors.New("ledger not found")
	ErrInvalidAmount          = errors.New("amount must be positive")
	ErrLedgerCategoryNotFound = errors.New("category not found")
	ErrLedgerPermissionDenied = errors.New("permission denied")
)

type LedgerService struct {
	repo         *repository.LedgerRepo
	categoryRepo *repository.CategoryRepo
}

func NewLedgerService() *LedgerService {
	return &LedgerService{
		repo:         &repository.LedgerRepo{},
		categoryRepo: &repository.CategoryRepo{},
	}
}

func (s *LedgerService) List(filter repository.LedgerFilter) (*repository.ListResult, error) {
	return s.repo.List(filter)
}

func (s *LedgerService) FindByID(id uint) (*repository.LedgerWithAssoc, error) {
	result, err := s.repo.FindByID(id)
	if err != nil {
		return nil, wrapNotFound(err, ErrLedgerNotFound)
	}
	return result, nil
}

func (s *LedgerService) Create(amount int64, note string, categoryID uint, occurredAt time.Time, creatorID uint) (*repository.LedgerWithAssoc, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Validate category exists
	_, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		return nil, wrapNotFound(err, ErrLedgerCategoryNotFound)
	}

	ledger := &model.Ledger{
		Amount:     amount,
		Note:       note,
		CategoryID: categoryID,
		CreatorID:  creatorID,
		OccurredAt: model.FromTime(occurredAt),
	}

	if err := s.repo.Create(ledger); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ledger.ID)
}

func (s *LedgerService) Update(id uint, amount *int64, note *string, categoryID *uint, occurredAt *time.Time, currentMemberID uint, currentRole string) (*repository.LedgerWithAssoc, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return nil, wrapNotFound(err, ErrLedgerNotFound)
	}

	// Permission check: only creator or admin can update
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return nil, ErrLedgerPermissionDenied
	}

	if amount != nil {
		if *amount <= 0 {
			return nil, ErrInvalidAmount
		}
		existing.Amount = *amount
	}

	if note != nil {
		existing.Note = *note
	}

	if categoryID != nil {
		_, err := s.categoryRepo.FindByID(*categoryID)
		if err != nil {
			return nil, wrapNotFound(err, ErrLedgerCategoryNotFound)
		}
		existing.CategoryID = *categoryID
	}

	if occurredAt != nil {
		existing.OccurredAt = model.FromTime(*occurredAt)
	}

	if err := s.repo.Update(&existing.Ledger); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *LedgerService) Delete(id uint, currentMemberID uint, currentRole string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		return wrapNotFound(err, ErrLedgerNotFound)
	}

	// Permission check: only creator or admin can delete
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrLedgerPermissionDenied
	}

	return s.repo.Delete(id)
}
