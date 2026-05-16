package service

import (
	"home-center/internal/model"
	"home-center/internal/pkg"
	"time"
)

type DashboardService struct{}

func NewDashboardService() *DashboardService {
	return &DashboardService{}
}

type CategorySum struct {
	CategoryID   uint   `json:"category_id"`
	CategoryName string `json:"category_name"`
	Icon         string `json:"icon"`
	Amount       int64  `json:"amount"`
}

type WishTrend struct {
	model.Wish
	Creator   model.Member `json:"creator"`
	VoteCount int64        `json:"vote_count"`
}

type FeedItem struct {
	Type      string       `json:"type"`
	ID        uint         `json:"id"`
	Title     string       `json:"title"`
	Content   string       `json:"content"`
	Creator   model.Member `json:"creator"`
	CreatedAt time.Time    `json:"created_at"`
}

func (s *DashboardService) GetSummary(month string) (map[string]int64, error) {
	type row struct {
		Type   string
		Amount int64
	}
	var rows []row
	err := pkg.DB.Table("ledgers").
		Select("categories.type, COALESCE(SUM(ledgers.amount), 0) as amount").
		Joins("JOIN categories ON ledgers.category_id = categories.id AND categories.deleted_at IS NULL").
		Where("strftime('%Y-%m', ledgers.occurred_at) = ?", month).
		Group("categories.type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}

	result := map[string]int64{"income": 0, "expense": 0, "balance": 0}
	for _, r := range rows {
		if r.Type == "income" {
			result["income"] = r.Amount
		} else if r.Type == "expense" {
			result["expense"] = r.Amount
		}
	}
	result["balance"] = result["income"] - result["expense"]
	return result, nil
}

func (s *DashboardService) GetExpenseChart(month string) ([]CategorySum, error) {
	var rows []CategorySum
	err := pkg.DB.Table("ledgers").
		Select("categories.id as category_id, categories.name as category_name, categories.icon, COALESCE(SUM(ledgers.amount), 0) as amount").
		Joins("JOIN categories ON ledgers.category_id = categories.id AND categories.deleted_at IS NULL").
		Where("categories.type = 'expense'").
		Where("strftime('%Y-%m', ledgers.occurred_at) = ?", month).
		Group("categories.id").
		Order("amount DESC").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	if rows == nil {
		rows = []CategorySum{}
	}
	return rows, nil
}

func (s *DashboardService) GetUpcomingTodos() ([]model.Todo, error) {
	var todos []model.Todo
	err := pkg.DB.
		Preload("Assignee").
		Where("status = 'pending'").
		Order("CASE priority WHEN 'urgent' THEN 0 WHEN 'important' THEN 1 ELSE 2 END, due_date ASC, created_at DESC").
		Limit(5).
		Find(&todos).Error
	if err != nil {
		return nil, err
	}
	if todos == nil {
		todos = []model.Todo{}
	}
	return todos, nil
}

func (s *DashboardService) GetWishTrends() ([]WishTrend, error) {
	var wishes []model.Wish
	err := pkg.DB.
		Preload("Creator").
		Order("created_at DESC").
		Limit(10).
		Find(&wishes).Error
	if err != nil {
		return nil, err
	}

	trends := make([]WishTrend, 0, len(wishes))
	for _, w := range wishes {
		var voteCount int64
		pkg.DB.Model(&model.WishVote{}).Where("wish_id = ?", w.ID).Count(&voteCount)
		trends = append(trends, WishTrend{
			Wish:      w,
			Creator:   w.Creator,
			VoteCount: voteCount,
		})
	}
	return trends, nil
}

func (s *DashboardService) GetForumHot() ([]FeedItem, error) {
	var items []FeedItem

	var posts []model.Post
	pkg.DB.Preload("Creator").Order("created_at DESC").Limit(5).Find(&posts)
	for _, p := range posts {
		content := p.Content
		if len([]rune(content)) > 100 {
			content = string([]rune(content)[:100]) + "..."
		}
		items = append(items, FeedItem{
			Type:      "post",
			ID:        p.ID,
			Title:     "",
			Content:   content,
			Creator:   p.Creator,
			CreatedAt: p.CreatedAt,
		})
	}

	var topics []model.Topic
	pkg.DB.Preload("Creator").Order("created_at DESC").Limit(5).Find(&topics)
	for _, t := range topics {
		content := t.Content
		if len([]rune(content)) > 100 {
			content = string([]rune(content)[:100]) + "..."
		}
		items = append(items, FeedItem{
			Type:      "topic",
			ID:        t.ID,
			Title:     t.Title,
			Content:   content,
			Creator:   t.Creator,
			CreatedAt: t.CreatedAt,
		})
	}

	// Sort merged list by created_at DESC
	for i := 0; i < len(items); i++ {
		for j := i + 1; j < len(items); j++ {
			if items[j].CreatedAt.After(items[i].CreatedAt) {
				items[i], items[j] = items[j], items[i]
			}
		}
	}
	if len(items) > 5 {
		items = items[:5]
	}
	if items == nil {
		items = []FeedItem{}
	}

	return items, nil
}
