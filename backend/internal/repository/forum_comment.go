package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

// --- Comments ---

// CommentWithMeta is a comment with associated creator info.
type CommentWithMeta struct {
	model.Comment
	Creator  model.Member      `json:"creator"`
	Children []CommentWithMeta `json:"children"`
}

// CreateComment inserts a new comment.
func (r *ForumRepo) CreateComment(comment *model.Comment) error {
	return pkg.DB.Create(comment).Error
}

// DeleteComment soft-deletes a comment.
func (r *ForumRepo) DeleteComment(id uint) error {
	// Soft delete children (level 2) first
	pkg.DB.Where("parent_id = ?", id).Delete(&model.Comment{})
	return pkg.DB.Delete(&model.Comment{}, id).Error
}

// FindCommentByID finds a comment by ID.
func (r *ForumRepo) FindCommentByID(id uint) (*model.Comment, error) {
	var c model.Comment
	err := pkg.DB.First(&c, id).Error
	return &c, err
}

// FindCommentWithCreator finds a comment with creator metadata.
func (r *ForumRepo) FindCommentWithCreator(id uint) (*CommentWithMeta, error) {
	var c model.Comment
	if err := pkg.DB.Preload("Creator").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &CommentWithMeta{Comment: c, Creator: c.Creator}, nil
}

// ListComments returns comments for a target type and ID.
func (r *ForumRepo) ListComments(targetType string, targetID uint) ([]CommentWithMeta, error) {
	var comments []model.Comment
	err := pkg.DB.Preload("Creator").
		Where("target_type = ? AND target_id = ? AND parent_id IS NULL", targetType, targetID).
		Order("created_at ASC").
		Find(&comments).Error
	if err != nil {
		return nil, err
	}

	result := make([]CommentWithMeta, 0, len(comments))
	for _, c := range comments {
		cm := CommentWithMeta{Comment: c, Creator: c.Creator}
		// Fetch children
		var children []model.Comment
		pkg.DB.Preload("Creator").
			Where("target_type = ? AND target_id = ? AND parent_id = ?", targetType, targetID, c.ID).
			Order("created_at ASC").
			Find(&children)
		for _, ch := range children {
			cm.Children = append(cm.Children, CommentWithMeta{Comment: ch, Creator: ch.Creator})
		}
		result = append(result, cm)
	}
	return result, nil
}
