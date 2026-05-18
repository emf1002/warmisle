package service

import (
	"errors"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrLedgerNotFound         = errors.New("ledger not found")
	ErrInvalidAmount          = errors.New("amount must be positive")
	ErrLedgerCategoryNotFound = errors.New("category not found")
	ErrNoMembers              = errors.New("at least one member required")
	ErrLedgerPermissionDenied = errors.New("permission denied")
)

type LedgerService struct {
	repo         *repository.LedgerRepo
	categoryRepo *repository.CategoryRepo
	memberRepo   *repository.MemberRepo
}

func NewLedgerService() *LedgerService {
	return &LedgerService{
		repo:         &repository.LedgerRepo{},
		categoryRepo: &repository.CategoryRepo{},
		memberRepo:   &repository.MemberRepo{},
	}
}

func (s *LedgerService) List(filter repository.LedgerFilter) (*repository.ListResult, error) {
	return s.repo.List(filter)
}

func (s *LedgerService) FindByID(id uint) (*repository.LedgerWithAssoc, error) {
	result, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLedgerNotFound
		}
		return nil, err
	}
	return result, nil
}

func (s *LedgerService) Create(amount int64, note string, categoryID uint, memberIDs []uint, occurredAt time.Time, creatorID uint) (*repository.LedgerWithAssoc, error) {
	if amount <= 0 {
		return nil, ErrInvalidAmount
	}

	// Validate category exists
	_, err := s.categoryRepo.FindByID(categoryID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLedgerCategoryNotFound
		}
		return nil, err
	}

	if len(memberIDs) == 0 {
		return nil, ErrNoMembers
	}

	ledger := &model.Ledger{
		Amount:     amount,
		Note:       note,
		CategoryID: categoryID,
		CreatorID:  creatorID,
		OccurredAt: occurredAt,
	}

	if err := s.repo.Create(ledger, memberIDs); err != nil {
		return nil, err
	}

	return s.repo.FindByID(ledger.ID)
}

func (s *LedgerService) Update(id uint, amount *int64, note *string, categoryID *uint, memberIDs []uint, occurredAt *time.Time, currentMemberID uint, currentRole string) (*repository.LedgerWithAssoc, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrLedgerNotFound
		}
		return nil, err
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
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrLedgerCategoryNotFound
			}
			return nil, err
		}
		existing.CategoryID = *categoryID
	}

	if occurredAt != nil {
		existing.OccurredAt = *occurredAt
	}

	// Members: if nil, keep existing; if empty, error; otherwise replace
	if memberIDs != nil {
		if len(memberIDs) == 0 {
			return nil, ErrNoMembers
		}
	} else {
		memberIDs = make([]uint, len(existing.Members))
		for i, m := range existing.Members {
			memberIDs[i] = m.ID
		}
	}

	if err := s.repo.Update(&existing.Ledger, memberIDs); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *LedgerService) Delete(id uint, currentMemberID uint, currentRole string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrLedgerNotFound
		}
		return err
	}

	// Permission check: only creator or admin can delete
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrLedgerPermissionDenied
	}

	return s.repo.Delete(id)
}
