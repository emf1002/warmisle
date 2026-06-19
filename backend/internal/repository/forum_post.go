package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

// --- Posts ---

// PostWithMeta is a post with associated creator info.
type PostWithMeta struct {
	model.Post
	Creator      model.Member `json:"creator"`
	LikeCount    int64        `json:"like_count"`
	CommentCount int64        `json:"comment_count"`
	Liked        bool         `json:"liked"`
}

// CreatePost inserts a new post.
func (r *ForumRepo) CreatePost(post *model.Post) error {
	return pkg.DB.Create(post).Error
}

// FindPostByID finds a post by ID with metadata.
func (r *ForumRepo) FindPostByID(id uint, currentMemberID uint) (*PostWithMeta, error) {
	var post model.Post
	if err := pkg.DB.Preload("Creator").First(&post, id).Error; err != nil {
		return nil, err
	}
	return r.enrichPost(&post, currentMemberID), nil
}

// UpdatePost modifies an existing post.
func (r *ForumRepo) UpdatePost(post *model.Post) error {
	return pkg.DB.Save(post).Error
}

// DeletePost soft-deletes a post.
func (r *ForumRepo) DeletePost(id uint) error {
	return pkg.DB.Delete(&model.Post{}, id).Error
}

func (r *ForumRepo) enrichPost(post *model.Post, currentMemberID uint) *PostWithMeta {
	var likeCount, commentCount int64
	pkg.DB.Model(&model.Like{}).Where("target_type = 'post' AND target_id = ?", post.ID).Count(&likeCount)
	pkg.DB.Model(&model.Comment{}).Where("target_type = 'post' AND target_id = ?", post.ID).Count(&commentCount)
	var liked int64
	pkg.DB.Model(&model.Like{}).Where("target_type = 'post' AND target_id = ? AND member_id = ?", post.ID, currentMemberID).Count(&liked)
	return &PostWithMeta{
		Post:         *post,
		Creator:      post.Creator,
		LikeCount:    likeCount,
		CommentCount: commentCount,
		Liked:        liked > 0,
	}
}
