package service

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/repository"
)

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
			return nil, wrapNotFound(err, ErrForumCommentNotFound)
		}
		if parent.ParentID != nil {
			return nil, ErrForumNestingTooDeep
		}
	}
	// XSS 防护：过滤 HTML 内容
	sanitizedContent := pkg.SanitizeHTML(content)
	comment := &model.Comment{
		TargetType: targetType,
		TargetID:   targetID,
		ParentID:   parentID,
		Content:    sanitizedContent,
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
		return wrapNotFound(err, ErrForumCommentNotFound)
	}
	if comment.CreatorID != currentMemberID && currentRole != "admin" {
		return ErrForumPermissionDenied
	}
	return s.repo.DeleteComment(id)
}

func (s *ForumService) ListComments(targetType string, targetID uint) ([]repository.CommentWithMeta, error) {
	return s.repo.ListComments(targetType, targetID)
}
