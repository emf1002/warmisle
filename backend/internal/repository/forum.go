package repository

import (
	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"time"
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
