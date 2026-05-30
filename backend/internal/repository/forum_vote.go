package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"gorm.io/gorm"
)

// --- Votes ---

type VoteWithDetail struct {
	model.Vote
	Creator    model.Member        `json:"creator"`
	Options    []VoteOptionSummary `json:"options"`
	TotalVotes int64               `json:"total_votes"`
}

type VoteOptionSummary struct {
	model.VoteOption
	VoteCount int64 `json:"vote_count"`
}

func (r *ForumRepo) CreateVote(vote *model.Vote, options []model.VoteOption) error {
	return pkg.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(vote).Error; err != nil {
			return err
		}
		for i := range options {
			options[i].VoteID = vote.ID
		}
		return tx.Create(&options).Error
	})
}

func (r *ForumRepo) FindVoteByID(id uint, currentMemberID uint) (*VoteWithDetail, error) {
	var vote model.Vote
	if err := pkg.DB.Preload("Creator").Preload("Options").First(&vote, id).Error; err != nil {
		return nil, err
	}
	var totalVotes int64
	pkg.DB.Model(&model.VoteRecord{}).Where("vote_id = ?", id).Count(&totalVotes)

	options := make([]VoteOptionSummary, 0, len(vote.Options))
	for _, opt := range vote.Options {
		var count int64
		pkg.DB.Model(&model.VoteRecord{}).Where("option_id = ?", opt.ID).Count(&count)
		options = append(options, VoteOptionSummary{VoteOption: opt, VoteCount: count})
	}

	return &VoteWithDetail{
		Vote:       vote,
		Creator:    vote.Creator,
		Options:    options,
		TotalVotes: totalVotes,
	}, nil
}

func (r *ForumRepo) DeleteVote(id uint) error {
	return pkg.DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Where("vote_id = ?", id).Delete(&model.VoteRecord{}).Error; err != nil {
			return err
		}
		if err := tx.Where("vote_id = ?", id).Delete(&model.VoteOption{}).Error; err != nil {
			return err
		}
		return tx.Delete(&model.Vote{}, id).Error
	})
}

func (r *ForumRepo) RecordVote(voteID, optionID, memberID uint) error {
	record := model.VoteRecord{VoteID: voteID, OptionID: optionID, MemberID: memberID}
	return pkg.DB.Create(&record).Error
}

func (r *ForumRepo) HasVotedForVote(voteID, memberID uint) (bool, error) {
	var count int64
	err := pkg.DB.Model(&model.VoteRecord{}).
		Where("vote_id = ? AND member_id = ?", voteID, memberID).
		Count(&count).Error
	return count > 0, err
}

// --- Likes ---

func (r *ForumRepo) ToggleLike(targetType string, targetID, memberID uint) (bool, error) {
	var existing model.Like
	err := pkg.DB.Where("target_type = ? AND target_id = ? AND member_id = ?", targetType, targetID, memberID).
		First(&existing).Error
	if err == nil {
		// Unlike
		return false, pkg.DB.Delete(&existing).Error
	}
	// Like
	like := model.Like{TargetType: targetType, TargetID: targetID, MemberID: memberID}
	return true, pkg.DB.Create(&like).Error
}

func (r *ForumRepo) GetLikeCount(targetType string, targetID uint) (int64, error) {
	var count int64
	err := pkg.DB.Model(&model.Like{}).
		Where("target_type = ? AND target_id = ?", targetType, targetID).
		Count(&count).Error
	return count, err
}
