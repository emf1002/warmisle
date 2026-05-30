package service

import (
	"warmisle/internal/model"
	"warmisle/internal/repository"
)

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
		return nil, wrapNotFound(err, ErrForumPostNotFound)
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
		return wrapNotFound(err, ErrForumPostNotFound)
	}
	if existing.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrForumPermissionDenied
	}
	return s.repo.DeletePost(id)
}
