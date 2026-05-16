package repository

import (
	"home-center/internal/model"
	"home-center/internal/pkg"
	"time"

	"gorm.io/gorm"
)

type ForumRepo struct{}

func tagPtr(t model.Topic) *model.Tag {
	if t.TagID == nil {
		return nil
	}
	return &t.Tag
}

// --- Feed ---

type FeedItem struct {
	Type      string       `json:"type"` // post / topic
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Creator   model.Member `json:"creator"`
	Tag       *model.Tag   `json:"tag,omitempty"`
	IsPinned  bool         `json:"is_pinned"`
	CreatedAt time.Time    `json:"created_at"`
}

type FeedResponse struct {
	Pinned []FeedItem `json:"pinned"`
	Items  []FeedItem `json:"items"`
	Total  int64      `json:"total"`
}

func (r *ForumRepo) GetFeed(page, pageSize int) (*FeedResponse, error) {
	// 1. Pinned topics
	var pinnedTopics []model.Topic
	pkg.DB.Preload("Creator").Preload("Tag").
		Where("is_pinned = ?", true).
		Order("created_at DESC").
		Find(&pinnedTopics)

	pinned := make([]FeedItem, 0, len(pinnedTopics))
	for _, t := range pinnedTopics {
		pinned = append(pinned, FeedItem{
			Type:      "topic",
			ID:        t.ID,
			Title:     t.Title,
			Content:   t.Content,
			Creator:   t.Creator,
			Tag:       tagPtr(t),
			IsPinned:  true,
			CreatedAt: t.CreatedAt,
		})
	}

	// 2. Count non-pinned items
	var postCount, topicCount int64
	pkg.DB.Model(&model.Post{}).Count(&postCount)
	pkg.DB.Model(&model.Topic{}).Where("is_pinned = ?", false).Count(&topicCount)
	total := postCount + topicCount

	// 3. Fetch posts and non-pinned topics as separate queries, then merge in Go
	var posts []model.Post
	pkg.DB.Preload("Creator").Order("created_at DESC").Find(&posts)

	var topics []model.Topic
	pkg.DB.Preload("Creator").Preload("Tag").
		Where("is_pinned = ?", false).
		Order("created_at DESC").
		Find(&topics)

	// Merge and sort
	allItems := make([]FeedItem, 0, len(posts)+len(topics))
	for _, p := range posts {
		allItems = append(allItems, FeedItem{
			Type:      "post",
			ID:        p.ID,
			Content:   p.Content,
			Creator:   p.Creator,
			CreatedAt: p.CreatedAt,
		})
	}
	for _, t := range topics {
		allItems = append(allItems, FeedItem{
			Type:      "topic",
			ID:        t.ID,
			Title:     t.Title,
			Content:   t.Content,
			Creator:   t.Creator,
			Tag:       tagPtr(t),
			CreatedAt: t.CreatedAt,
		})
	}

	// Sort by created_at DESC
	for i := 0; i < len(allItems); i++ {
		for j := i + 1; j < len(allItems); j++ {
			if allItems[j].CreatedAt.After(allItems[i].CreatedAt) {
				allItems[i], allItems[j] = allItems[j], allItems[i]
			}
		}
	}

	// Paginate
	offset := (page - 1) * pageSize
	items := make([]FeedItem, 0)
	if offset < len(allItems) {
		end := offset + pageSize
		if end > len(allItems) {
			end = len(allItems)
		}
		items = allItems[offset:end]
	}

	return &FeedResponse{
		Pinned: pinned,
		Items:  items,
		Total:  total,
	}, nil
}

// --- Posts ---

type PostWithMeta struct {
	model.Post
	Creator     model.Member `json:"creator"`
	LikeCount   int64        `json:"like_count"`
	CommentCount int64       `json:"comment_count"`
	Liked       bool         `json:"liked"`
}

func (r *ForumRepo) CreatePost(post *model.Post) error {
	return pkg.DB.Create(post).Error
}

func (r *ForumRepo) FindPostByID(id uint, currentMemberID uint) (*PostWithMeta, error) {
	var post model.Post
	if err := pkg.DB.Preload("Creator").First(&post, id).Error; err != nil {
		return nil, err
	}
	return r.enrichPost(&post, currentMemberID), nil
}

func (r *ForumRepo) UpdatePost(post *model.Post) error {
	return pkg.DB.Save(post).Error
}

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

// --- Comments ---

type CommentWithMeta struct {
	model.Comment
	Creator  model.Member      `json:"creator"`
	Children []CommentWithMeta `json:"children"`
}

func (r *ForumRepo) CreateComment(comment *model.Comment) error {
	return pkg.DB.Create(comment).Error
}

func (r *ForumRepo) DeleteComment(id uint) error {
	// Soft delete children (level 2) first
	pkg.DB.Where("parent_id = ?", id).Delete(&model.Comment{})
	return pkg.DB.Delete(&model.Comment{}, id).Error
}

func (r *ForumRepo) FindCommentByID(id uint) (*model.Comment, error) {
	var c model.Comment
	err := pkg.DB.First(&c, id).Error
	return &c, err
}

func (r *ForumRepo) FindCommentWithCreator(id uint) (*CommentWithMeta, error) {
	var c model.Comment
	if err := pkg.DB.Preload("Creator").First(&c, id).Error; err != nil {
		return nil, err
	}
	return &CommentWithMeta{Comment: c, Creator: c.Creator}, nil
}

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

// --- Tags ---

func (r *ForumRepo) ListTags() ([]model.Tag, error) {
	var tags []model.Tag
	err := pkg.DB.Order("id ASC").Find(&tags).Error
	return tags, err
}

func (r *ForumRepo) CreateTag(tag *model.Tag) error {
	return pkg.DB.Create(tag).Error
}

func (r *ForumRepo) UpdateTag(tag *model.Tag) error {
	return pkg.DB.Save(tag).Error
}

func (r *ForumRepo) FindTagByID(id uint) (*model.Tag, error) {
	var tag model.Tag
	err := pkg.DB.First(&tag, id).Error
	return &tag, err
}

func (r *ForumRepo) DeleteTag(id uint) (int64, error) {
	var count int64
	pkg.DB.Model(&model.Topic{}).Where("tag_id = ?", id).Count(&count)
	if count > 0 {
		return count, nil
	}
	return 0, pkg.DB.Delete(&model.Tag{}, id).Error
}
