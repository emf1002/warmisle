package service

import (
	"errors"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/repository"
)

func (s *ForumService) ToggleLike(targetType string, targetID, memberID uint) (bool, error) {
	if !validTargetTypes[targetType] {
		return false, ErrForumInvalidTargetType
	}
	return s.repo.ToggleLike(targetType, targetID, memberID)
}

func (s *ForumService) CreateVote(title string, options []string, isMulti bool, deadline *time.Time, creatorID uint) (*repository.VoteWithDetail, error) {
	if title == "" {
		return nil, ErrForumTitleRequired
	}
	if len(options) < 2 {
		return nil, errors.New("at least 2 options required")
	}
	vote := &model.Vote{
		Title:     title,
		CreatorID: creatorID,
		IsMulti:   isMulti,
		Deadline:  model.FromTimePtr(deadline),
	}
	voteOptions := make([]model.VoteOption, len(options))
	for i, opt := range options {
		voteOptions[i] = model.VoteOption{Content: opt, SortOrder: i}
	}
	if err := s.repo.CreateVote(vote, voteOptions); err != nil {
		return nil, err
	}
	return s.repo.FindVoteByID(vote.ID, creatorID)
}

func (s *ForumService) DeleteVote(id uint, currentMemberID uint, currentRole string) error {
	vote, err := s.repo.FindVoteByID(id, currentMemberID)
	if err != nil {
		return ErrForumVoteNotFound
	}
	if vote.Deadline != nil && time.Now().After(vote.Deadline.ToTime()) {
		return ErrForumVoteDeadlinePassed
	}
	if vote.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrForumPermissionDenied
	}
	return s.repo.DeleteVote(id)
}

func (s *ForumService) Vote(id uint, optionID, memberID uint) (*repository.VoteWithDetail, error) {
	vote, err := s.repo.FindVoteByID(id, memberID)
	if err != nil {
		return nil, ErrForumVoteNotFound
	}
	if vote.Deadline != nil && time.Now().After(vote.Deadline.ToTime()) {
		return nil, ErrForumVoteDeadlinePassed
	}
	hasVoted, err := s.repo.HasVotedForVote(id, memberID)
	if err != nil {
		return nil, err
	}
	if hasVoted {
		return nil, ErrForumAlreadyVoted
	}
	if err := s.repo.RecordVote(id, optionID, memberID); err != nil {
		return nil, err
	}
	return s.repo.FindVoteByID(id, memberID)
}

func (s *ForumService) GetVote(id uint, memberID uint) (*repository.VoteWithDetail, error) {
	result, err := s.repo.FindVoteByID(id, memberID)
	if err != nil {
		return nil, ErrForumVoteNotFound
	}
	return result, nil
}
