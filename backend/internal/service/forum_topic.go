package service

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/repository"
)

func (s *ForumService) CreateTopic(title, content string, tagID *uint, creatorID uint) (*repository.TopicWithMeta, error) {
	if title == "" {
		return nil, ErrForumTitleRequired
	}
	// XSS 防护：过滤 HTML 内容
	sanitizedContent := pkg.SanitizeHTML(content)
	topic := &model.Topic{Title: title, Content: sanitizedContent, TagID: tagID, CreatorID: creatorID}
	if err := s.repo.CreateTopic(topic); err != nil {
		return nil, err
	}
	return s.repo.FindTopicByID(topic.ID, creatorID)
}

func (s *ForumService) UpdateTopic(id uint, title, content *string, tagID *uint, currentMemberID uint, currentRole string) (*repository.TopicWithMeta, error) {
	existing, err := s.repo.FindTopicByID(id, currentMemberID)
	if err != nil {
		return nil, wrapNotFound(err, ErrForumTopicNotFound)
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
		sanitized := pkg.SanitizeHTML(*content)
		existing.Content = sanitized
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
		return wrapNotFound(err, ErrForumTopicNotFound)
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
		return nil, wrapNotFound(err, ErrForumTopicNotFound)
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
		return nil, wrapNotFound(err, ErrForumTopicNotFound)
	}
	return result, nil
}
