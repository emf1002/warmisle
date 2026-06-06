package service

import (
	"fmt"
	"sort"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/repository"
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

// monthBounds returns the start and end time for a "YYYY-MM" string.
func monthBounds(month string) (time.Time, time.Time, error) {
	t, err := time.Parse("2006-01", month)
	if err != nil {
		return time.Time{}, time.Time{}, fmt.Errorf("invalid month format: %w", err)
	}
	start := t
	end := t.AddDate(0, 1, 0)
	return start, end, nil
}

func (s *DashboardService) GetSummary(month string) (map[string]int64, error) {
	start, end, err := monthBounds(month)
	if err != nil {
		return nil, err
	}

	type row struct {
		Type   string
		Amount int64
	}
	var rows []row
	err = pkg.DB.Table("ledgers").
		Select("categories.type, COALESCE(SUM(ledgers.amount), 0) as amount").
		Joins("JOIN categories ON ledgers.category_id = categories.id AND categories.deleted_at IS NULL").
		Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", start, end).
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
	start, end, err := monthBounds(month)
	if err != nil {
		return nil, err
	}

	var rows []CategorySum
	err = pkg.DB.Table("ledgers").
		Select("categories.id as category_id, categories.name as category_name, categories.icon, COALESCE(SUM(ledgers.amount), 0) as amount").
		Joins("JOIN categories ON ledgers.category_id = categories.id AND categories.deleted_at IS NULL").
		Where("categories.type = 'expense'").
		Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", start, end).
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

	// Batch load vote counts to avoid N+1
	wishIDs := make([]uint, len(wishes))
	for i, w := range wishes {
		wishIDs[i] = w.ID
	}
	type voteRow struct {
		WishID    uint
		VoteCount int64
	}
	var voteRows []voteRow
	pkg.DB.Model(&model.WishVote{}).
		Select("wish_id, COUNT(*) as vote_count").
		Where("wish_id IN ?", wishIDs).
		Group("wish_id").
		Scan(&voteRows)

	voteMap := make(map[uint]int64, len(voteRows))
	for _, vr := range voteRows {
		voteMap[vr.WishID] = vr.VoteCount
	}

	trends := make([]WishTrend, 0, len(wishes))
	for _, w := range wishes {
		trends = append(trends, WishTrend{
			Wish:      w,
			Creator:   w.Creator,
			VoteCount: voteMap[w.ID],
		})
	}
	return trends, nil
}

func (s *DashboardService) GetForumHot() ([]repository.FeedItem, error) {
	var items []repository.FeedItem

	var posts []model.Post
	pkg.DB.Preload("Creator").Order("created_at DESC").Limit(5).Find(&posts)
	for _, p := range posts {
		content := p.Content
		if len([]rune(content)) > 100 {
			content = string([]rune(content)[:100]) + "..."
		}
		items = append(items, repository.FeedItem{
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
		items = append(items, repository.FeedItem{
			Type:      "topic",
			ID:        t.ID,
			Title:     t.Title,
			Content:   content,
			Creator:   t.Creator,
			CreatedAt: t.CreatedAt,
		})
	}

	// Sort merged list by created_at DESC
	sort.Slice(items, func(i, j int) bool {
		return items[i].CreatedAt.After(items[j].CreatedAt)
	})
	if len(items) > 5 {
		items = items[:5]
	}
	if items == nil {
		items = []repository.FeedItem{}
	}

	return items, nil
}
