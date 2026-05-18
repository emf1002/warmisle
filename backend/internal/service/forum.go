package service

import (
	"errors"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/repository"

	"gorm.io/gorm"
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
	ErrForumVoteDeadlinePassed  = errors.New("voting deadline has passed")
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

// --- Posts ---

func (s *ForumService) CreatePost(content string, creatorID uint) (*repository.PostWithMeta, error) {
	if content == "" {
		return nil, ErrForumContentRequired
	}
	post := &model.Post{Content: content, CreatorID: creatorID}
	if err := s.repo.CreatePost(post); err != nil {
		return nil, err
	}
	return s.repo.FindPostByID(post.ID, creatorID)
}

func (s *ForumService) UpdatePost(id uint, content string, currentMemberID uint, currentRole string) (*repository.PostWithMeta, error) {
	existing, err := s.repo.FindPostByID(id, currentMemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrForumPostNotFound
		}
		return nil, err
	}
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return nil, ErrForumPermissionDenied
	}
	if content == "" {
		return nil, ErrForumContentRequired
	}
	existing.Content = content
	if err := s.repo.UpdatePost(&existing.Post); err != nil {
		return nil, err
	}
	return s.repo.FindPostByID(id, currentMemberID)
}

func (s *ForumService) DeletePost(id uint, currentMemberID uint, currentRole string) error {
	existing, err := s.repo.FindPostByID(id, currentMemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrForumPostNotFound
		}
		return err
	}
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrForumPermissionDenied
	}
	return s.repo.DeletePost(id)
}

// --- Topics ---

func (s *ForumService) CreateTopic(title, content string, tagID *uint, creatorID uint) (*repository.TopicWithMeta, error) {
	if title == "" {
		return nil, ErrForumTitleRequired
	}
	topic := &model.Topic{Title: title, Content: content, TagID: tagID, CreatorID: creatorID}
	if err := s.repo.CreateTopic(topic); err != nil {
		return nil, err
	}
	return s.repo.FindTopicByID(topic.ID, creatorID)
}

func (s *ForumService) UpdateTopic(id uint, title, content *string, tagID *uint, currentMemberID uint, currentRole string) (*repository.TopicWithMeta, error) {
	existing, err := s.repo.FindTopicByID(id, currentMemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrForumTopicNotFound
		}
		return nil, err
	}
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return nil, ErrForumPermissionDenied
	}
	if title != nil {
		if *title == "" {
			return nil, ErrForumTitleRequired
		}
		existing.Title = *title
	}
	if content != nil {
		existing.Content = *content
	}
	if tagID != nil {
		existing.TagID = tagID
	}
	if err := s.repo.UpdateTopic(&existing.Topic); err != nil {
		return nil, err
	}
	return s.repo.FindTopicByID(id, currentMemberID)
}

func (s *ForumService) DeleteTopic(id uint, currentMemberID uint, currentRole string) error {
	existing, err := s.repo.FindTopicByID(id, currentMemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrForumTopicNotFound
		}
		return err
	}
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrForumPermissionDenied
	}
	return s.repo.DeleteTopic(id)
}

func (s *ForumService) TogglePin(id uint, currentRole string) (*repository.TopicWithMeta, error) {
	if currentRole != "admin" {
		return nil, ErrForumPermissionDenied
	}
	existing, err := s.repo.FindTopicByID(id, 0)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrForumTopicNotFound
		}
		return nil, err
	}
	existing.IsPinned = !existing.IsPinned
	if err := s.repo.UpdateTopic(&existing.Topic); err != nil {
		return nil, err
	}
	return s.repo.FindTopicByID(id, 0)
}

func (s *ForumService) GetTopic(id uint, currentMemberID uint) (*repository.TopicWithMeta, error) {
	result, err := s.repo.FindTopicByID(id, currentMemberID)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, ErrForumTopicNotFound
		}
		return nil, err
	}
	return result, nil
}

// --- Comments ---

func (s *ForumService) CreateComment(targetType string, targetID uint, parentID *uint, content string, creatorID uint) (*repository.CommentWithMeta, error) {
	if !validTargetTypes[targetType] && targetType != "wish" {
		return nil, ErrForumInvalidTargetType
	}
	if content == "" {
		return nil, ErrForumContentRequired
	}
	// Validate 2-level nesting
	if parentID != nil {
		parent, err := s.repo.FindCommentByID(*parentID)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, ErrForumCommentNotFound
			}
			return nil, err
		}
		if parent.ParentID != nil {
			return nil, ErrForumNestingTooDeep
		}
	}
	comment := &model.Comment{
		TargetType: targetType,
		TargetID:   targetID,
		ParentID:   parentID,
		Content:    content,
		CreatorID:  creatorID,
	}
	if err := s.repo.CreateComment(comment); err != nil {
		return nil, err
	}
	// Return the created comment with creator info
	return s.repo.FindCommentWithCreator(comment.ID)
}

func (s *ForumService) DeleteComment(id uint, currentMemberID uint, currentRole string) error {
	comment, err := s.repo.FindCommentByID(id)
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return ErrForumCommentNotFound
		}
		return err
	}
	if comment.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrForumPermissionDenied
	}
	return s.repo.DeleteComment(id)
}

func (s *ForumService) ListComments(targetType string, targetID uint) ([]repository.CommentWithMeta, error) {
	return s.repo.ListComments(targetType, targetID)
}

// --- Likes ---

func (s *ForumService) ToggleLike(targetType string, targetID, memberID uint) (bool, error) {
	if !validTargetTypes[targetType] {
		return false, ErrForumInvalidTargetType
	}
	return s.repo.ToggleLike(targetType, targetID, memberID)
}

// --- Votes ---

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
		Deadline:  deadline,
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
	if vote.Deadline != nil && time.Now().After(*vote.Deadline) {
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
	if vote.Deadline != nil && time.Now().After(*vote.Deadline) {
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
