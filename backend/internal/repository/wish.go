package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

// WishRepo handles wish data access.
type WishRepo struct{}

// WishFilter specifies query filters for wishes.
type WishFilter struct {
	Type      string
	Status    string
	CreatorID *uint
	Page      int
	PageSize  int
}

// WishWithAssoc is a wish with associated creator and vote count.
type WishWithAssoc struct {
	model.Wish
	Creator   model.Member `json:"creator"`
	VoteCount int64        `json:"vote_count"`
}

// WishListResult is the paginated wish list response.
type WishListResult struct {
	List     []WishWithAssoc `json:"list"`
	Total    int64           `json:"total"`
	Page     int             `json:"page"`
	PageSize int             `json:"page_size"`
}

// List returns wishes matching the given filter.
func (r *WishRepo) List(filter WishFilter) (*WishListResult, error) {
	query := pkg.DB.Model(&model.Wish{}).Preload("Creator")

	if filter.Type != "" {
		query = query.Where("wishes.type = ?", filter.Type)
	}
	if filter.Status != "" {
		query = query.Where("wishes.status = ?", filter.Status)
	}
	if filter.CreatorID != nil {
		query = query.Where("wishes.creator_id = ?", *filter.CreatorID)
	}

	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	var wishes []model.Wish
	err := query.
		Order("wishes.created_at DESC").
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Find(&wishes).Error
	if err != nil {
		return nil, err
	}

	list := make([]WishWithAssoc, 0, len(wishes))
	for _, w := range wishes {
		var voteCount int64
		pkg.DB.Model(&model.WishVote{}).Where("wish_id = ?", w.ID).Count(&voteCount)
		list = append(list, WishWithAssoc{
			Wish:      w,
			Creator:   w.Creator,
			VoteCount: voteCount,
		})
	}

	return &WishListResult{
		List:     list,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

// FindByID finds a wish by ID with metadata.
func (r *WishRepo) FindByID(id uint) (*WishWithAssoc, error) {
	var wish model.Wish
	err := pkg.DB.Preload("Creator").First(&wish, id).Error
	if err != nil {
		return nil, err
	}
	var voteCount int64
	pkg.DB.Model(&model.WishVote{}).Where("wish_id = ?", id).Count(&voteCount)
	return &WishWithAssoc{
		Wish:      wish,
		Creator:   wish.Creator,
		VoteCount: voteCount,
	}, nil
}

// Create inserts a new wish.
func (r *WishRepo) Create(wish *model.Wish) error {
	return pkg.DB.Create(wish).Error
}

// Update modifies an existing wish.
func (r *WishRepo) Update(wish *model.Wish) error {
	return pkg.DB.Save(wish).Error
}

// Delete soft-deletes a wish.
func (r *WishRepo) Delete(id uint) error {
	return pkg.DB.Delete(&model.Wish{}, id).Error
}

// HasVoted checks if a member has voted on a wish.
func (r *WishRepo) HasVoted(wishID, memberID uint) (bool, error) {
	var count int64
	err := pkg.DB.Model(&model.WishVote{}).
		Where("wish_id = ? AND member_id = ?", wishID, memberID).
		Count(&count).Error
	return count > 0, err
}

// CreateVote casts a vote on a wish.
func (r *WishRepo) CreateVote(vote *model.WishVote) error {
	return pkg.DB.Create(vote).Error
}

// DeleteVote removes a vote from a wish.
func (r *WishRepo) DeleteVote(wishID, memberID uint) error {
	return pkg.DB.Where("wish_id = ? AND member_id = ?", wishID, memberID).Delete(&model.WishVote{}).Error
}

// CountByStatus counts wishes by ID and status.
func (r *WishRepo) CountByStatus(wishID uint, status string) (int64, error) {
	var count int64
	err := pkg.DB.Model(&model.Wish{}).Where("id = ? AND status = ?", wishID, status).Count(&count).Error
	return count, err
}
