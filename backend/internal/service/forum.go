package service

import (
	"errors"

	"warmisle/internal/model"
	"warmisle/internal/repository"
)

var (
	ErrForumPostNotFound       = errors.New("post not found")
	ErrForumTopicNotFound      = errors.New("topic not found")
	ErrForumCommentNotFound    = errors.New("comment not found")
	ErrForumVoteNotFound       = errors.New("vote not found")
	ErrForumTagNotFound        = errors.New("tag not found")
	ErrForumPermissionDenied   = errors.New("permission denied")
	ErrForumContentRequired    = errors.New("content is required")
	ErrForumTitleRequired      = errors.New("title is required")
	ErrForumNestingTooDeep     = errors.New("cannot reply to a reply")
	ErrForumVoteDeadlinePassed = errors.New("voting deadline has passed")
	ErrForumAlreadyVoted       = errors.New("already voted")
	ErrForumTagInUse           = errors.New("tag is in use")
	ErrForumTagNameTaken       = errors.New("tag name already exists")
	ErrForumInvalidTargetType  = errors.New("invalid target type")
)

var validTargetTypes = map[string]bool{"post": true, "topic": true, "comment": true}

type ForumService struct {
	repo *repository.ForumRepo
}

func NewForumService() *ForumService {
	return &ForumService{repo: &repository.ForumRepo{}}
}

// --- Feed ---

func (s *ForumService) GetFeed(page, pageSize int) (*repository.FeedResponse, error) {
	return s.repo.GetFeed(page, pageSize)
}

// --- Tags ---

func (s *ForumService) ListTags() ([]model.Tag, error) {
	return s.repo.ListTags()
}

func (s *ForumService) CreateTag(name string) (*model.Tag, error) {
	if name == "" {
		return nil, ErrForumContentRequired
	}
	tag := &model.Tag{Name: name}
	if err := s.repo.CreateTag(tag); err != nil {
		return nil, ErrForumTagNameTaken
	}
	return tag, nil
}

func (s *ForumService) UpdateTag(id uint, name string) (*model.Tag, error) {
	if name == "" {
		return nil, ErrForumContentRequired
	}
	tag, err := s.repo.FindTagByID(id)
	if err != nil {
		return nil, ErrForumTagNotFound
	}
	tag.Name = name
	if err := s.repo.UpdateTag(tag); err != nil {
		return nil, ErrForumTagNameTaken
	}
	return tag, nil
}

func (s *ForumService) DeleteTag(id uint) error {
	count, err := s.repo.DeleteTag(id)
	if err != nil {
		return err
	}
	if count > 0 {
		return ErrForumTagInUse
	}
	return nil
}
