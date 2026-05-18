package service

import (
	"errors"

	"warmisle/internal/model"
	"warmisle/internal/repository"

	"gorm.io/gorm"
)

var (
	ErrWishNotFound         = errors.New("wish not found")
	ErrWishPermissionDenied = errors.New("permission denied")
	ErrWishTitleRequired    = errors.New("title is required")
	ErrWishInvalidPriority  = errors.New("invalid priority value")
	ErrWishInvalidCategory  = errors.New("invalid category value")
	ErrWishAlreadyVoted     = errors.New("already voted")
	ErrWishNotVoted         = errors.New("not voted yet")
	ErrWishInvalidStatus    = errors.New("invalid status value")
)

var validWishCategories = map[string]bool{"item": true, "travel": true, "experience": true, "other": true}
var validWishStatuses = map[string]bool{"pending": true, "agreed": true, "achieved": true, "abandoned": true}

type WishService struct {
	repo       *repository.WishRepo
	memberRepo *repository.MemberRepo
}

func NewWishService() *WishService {
	return &WishService{
		repo:       &repository.WishRepo{},
		memberRepo: &repository.MemberRepo{},
	}
}

func validWishPriority(p string) bool {
	return p == "normal" || p == "important" || p == "urgent"
}

func (s *WishService) FindByID(id uint) (*repository.WishWithAssoc, error) {
	result, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWishNotFound
		}
		return nil, err
	}
	return result, nil
}

func (s *WishService) List(filter repository.WishFilter) (*repository.WishListResult, error) {
	return s.repo.List(filter)
}

func (s *WishService) Create(title, description, category, priority string, amount *int64, creatorID uint) (*repository.WishWithAssoc, error) {
	if title == "" {
		return nil, ErrWishTitleRequired
	}

	if priority == "" {
		priority = "normal"
	}
	if !validWishPriority(priority) {
		return nil, ErrWishInvalidPriority
	}

	if category == "" {
		category = "other"
	}
	if !validWishCategories[category] {
		return nil, ErrWishInvalidCategory
	}

	wish := &model.Wish{
		Title:       title,
		Description: description,
		Category:    category,
		Priority:    priority,
		Amount:      amount,
		Type:        "personal",
		Status:      "pending",
		CreatorID:   creatorID,
	}

	if err := s.repo.Create(wish); err != nil {
		return nil, err
	}

	return s.repo.FindByID(wish.ID)
}

func (s *WishService) Update(id uint, title, description, category, priority *string, amount *int64, currentMemberID uint, currentRole string) (*repository.WishWithAssoc, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWishNotFound
		}
		return nil, err
	}

	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return nil, ErrWishPermissionDenied
	}

	if title != nil {
		if *title == "" {
			return nil, ErrWishTitleRequired
		}
		existing.Title = *title
	}
	if description != nil {
		existing.Description = *description
	}
	if category != nil {
		if !validWishCategories[*category] {
			return nil, ErrWishInvalidCategory
		}
		existing.Category = *category
	}
	if priority != nil {
		if !validWishPriority(*priority) {
			return nil, ErrWishInvalidPriority
		}
		existing.Priority = *priority
	}
	if amount != nil {
		existing.Amount = amount
	}

	if err := s.repo.Update(&existing.Wish); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *WishService) Delete(id uint, currentMemberID uint, currentRole string) error {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrWishNotFound
		}
		return err
	}

	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrWishPermissionDenied
	}

	return s.repo.Delete(id)
}

func (s *WishService) Promote(id uint, currentMemberID uint) (*repository.WishWithAssoc, error) {
	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWishNotFound
		}
		return nil, err
	}

	if existing.CreatorID != currentMemberID {
		return nil, ErrWishPermissionDenied
	}

	existing.Type = "family"

	if err := s.repo.Update(&existing.Wish); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *WishService) UpdateStatus(id uint, status string, currentMemberID uint, currentRole string) (*repository.WishWithAssoc, error) {
	if !validWishStatuses[status] {
		return nil, ErrWishInvalidStatus
	}

	existing, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWishNotFound
		}
		return nil, err
	}

	if currentRole != "admin" {
		// 创建者仅可标记放弃
		if existing.CreatorID != currentMemberID || status != "abandoned" {
			return nil, ErrWishPermissionDenied
		}
	}

	existing.Status = status

	if err := s.repo.Update(&existing.Wish); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *WishService) Vote(id, currentMemberID uint) (*repository.WishWithAssoc, error) {
	_, err := s.repo.FindByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWishNotFound
		}
		return nil, err
	}

	hasVoted, err := s.repo.HasVoted(id, currentMemberID)
	if err != nil {
		return nil, err
	}
	if hasVoted {
		return nil, ErrWishAlreadyVoted
	}

	vote := &model.WishVote{
		WishID:   id,
		MemberID: currentMemberID,
	}
	if err := s.repo.CreateVote(vote); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}

func (s *WishService) Unvote(id, currentMemberID uint) (*repository.WishWithAssoc, error) {
	if _, err := s.repo.FindByID(id); err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrWishNotFound
		}
		return nil, err
	}

	hasVoted, err := s.repo.HasVoted(id, currentMemberID)
	if err != nil {
		return nil, err
	}
	if !hasVoted {
		return nil, ErrWishNotVoted
	}

	if err := s.repo.DeleteVote(id, currentMemberID); err != nil {
		return nil, err
	}

	return s.repo.FindByID(id)
}
