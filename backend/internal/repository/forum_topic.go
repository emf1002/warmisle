package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

// --- Topics ---

type TopicWithMeta struct {
	model.Topic
	Creator      model.Member `json:"creator"`
	Tag          *model.Tag   `json:"tag"`
	LikeCount    int64        `json:"like_count"`
	CommentCount int64        `json:"comment_count"`
	Liked        bool         `json:"liked"`
}

func (r *ForumRepo) CreateTopic(topic *model.Topic) error {
	return pkg.DB.Create(topic).Error
}

func (r *ForumRepo) FindTopicByID(id uint, currentMemberID uint) (*TopicWithMeta, error) {
	var topic model.Topic
	if err := pkg.DB.Preload("Creator").Preload("Tag").First(&topic, id).Error; err != nil {
		return nil, err
	}
	var likeCount, commentCount int64
	pkg.DB.Model(&model.Like{}).Where("target_type = 'topic' AND target_id = ?", topic.ID).Count(&likeCount)
	pkg.DB.Model(&model.Comment{}).Where("target_type = 'topic' AND target_id = ?", topic.ID).Count(&commentCount)
	var liked int64
	pkg.DB.Model(&model.Like{}).Where("target_type = 'topic' AND target_id = ? AND member_id = ?", topic.ID, currentMemberID).Count(&liked)
	return &TopicWithMeta{
		Topic:        topic,
		Creator:      topic.Creator,
		Tag:          tagPtr(topic),
		LikeCount:    likeCount,
		CommentCount: commentCount,
		Liked:        liked > 0,
	}, nil
}

func (r *ForumRepo) UpdateTopic(topic *model.Topic) error {
	return pkg.DB.Save(topic).Error
}

func (r *ForumRepo) DeleteTopic(id uint) error {
	return pkg.DB.Delete(&model.Topic{}, id).Error
}
