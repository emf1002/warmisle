package repository

import (
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
)

type LedgerRepo struct{}

type LedgerFilter struct {
	Month      string // "2026-05"
	MemberID   *uint
	CategoryID *uint
	CreatorID  *uint
	Page       int
	PageSize   int
}

type LedgerGroup struct {
	Date       string           `json:"date"`
	DailyTotal int64            `json:"daily_total"`
	Items      []LedgerWithAssoc `json:"items"`
}

type LedgerWithAssoc struct {
	model.Ledger
	Category model.Category `json:"category"`
	Creator  model.Member   `json:"creator"`
	Members  []model.Member `json:"members"`
}

type ListResult struct {
	Summary  LedgerSummary `json:"summary"`
	Groups   []LedgerGroup `json:"groups"`
	Total    int64         `json:"total"`
	Page     int           `json:"page"`
	PageSize int           `json:"page_size"`
}

type LedgerSummary struct {
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
	Balance int64 `json:"balance"`
}

func (r *LedgerRepo) List(filter LedgerFilter) (*ListResult, error) {
	// Calculate summary for the month
	var summary LedgerSummary
	type summaryRow struct {
		Type   string
		Amount int64
	}
	var rows []summaryRow
	err := pkg.DB.Table("ledgers").
		Select("categories.type, COALESCE(SUM(ledgers.amount), 0) as amount").
		Joins("JOIN categories ON ledgers.category_id = categories.id AND categories.deleted_at IS NULL").
		Where("strftime('%Y-%m', ledgers.occurred_at) = ?", filter.Month).
		Group("categories.type").
		Scan(&rows).Error
	if err != nil {
		return nil, err
	}
	for _, row := range rows {
		if row.Type == "income" {
			summary.Income = row.Amount
		} else if row.Type == "expense" {
			summary.Expense = row.Amount
		}
	}
	summary.Balance = summary.Income - summary.Expense

	// Build query
	query := pkg.DB.Model(&model.Ledger{}).
		Preload("Category").
		Preload("Creator").
		Preload("Members").
		Where("strftime('%Y-%m', ledgers.occurred_at) = ?", filter.Month)

	if filter.CategoryID != nil {
		query = query.Where("ledgers.category_id = ?", *filter.CategoryID)
	}
	if filter.CreatorID != nil {
		query = query.Where("ledgers.creator_id = ?", *filter.CreatorID)
	}
	if filter.MemberID != nil {
		subQuery := pkg.DB.Table("ledger_members").Select("ledger_id").Where("member_id = ?", *filter.MemberID)
		query = query.Where("ledgers.id IN (?)", subQuery)
	}

	// Count total
	var total int64
	if err := query.Count(&total).Error; err != nil {
		return nil, err
	}

	// Fetch with pagination
	var ledgers []model.Ledger
	err = query.
		Offset((filter.Page - 1) * filter.PageSize).
		Limit(filter.PageSize).
		Order("ledgers.occurred_at DESC, ledgers.id DESC").
		Find(&ledgers).Error
	if err != nil {
		return nil, err
	}

	// Group by date
	groups := r.groupByDate(ledgers)

	return &ListResult{
		Summary:  summary,
		Groups:   groups,
		Total:    total,
		Page:     filter.Page,
		PageSize: filter.PageSize,
	}, nil
}

func (r *LedgerRepo) groupByDate(ledgers []model.Ledger) []LedgerGroup {
	groupMap := make(map[string][]LedgerWithAssoc)
	var order []string

	for _, l := range ledgers {
		dateKey := time.Time(l.OccurredAt).Format("2006-01-02")
		item := LedgerWithAssoc{
			Ledger:   l,
			Category: l.Category,
			Creator:  l.Creator,
			Members:  l.Members,
		}
		if _, exists := groupMap[dateKey]; !exists {
			order = append(order, dateKey)
		}
		groupMap[dateKey] = append(groupMap[dateKey], item)
	}

	groups := make([]LedgerGroup, 0, len(order))
	for _, dateKey := range order {
		items := groupMap[dateKey]
		// Calculate daily total: income - expense for the day
		var dailyTotal int64
		for _, item := range items {
			if item.Category.Type == "income" {
				dailyTotal += item.Amount
			} else {
				dailyTotal -= item.Amount
			}
		}
		groups = append(groups, LedgerGroup{
			Date:       dateKey,
			DailyTotal: dailyTotal,
			Items:      items,
		})
	}
	return groups
}

func (r *LedgerRepo) Create(ledger *model.Ledger, memberIDs []uint) error {
	if err := pkg.DB.Create(ledger).Error; err != nil {
		return err
	}

	for _, memberID := range memberIDs {
		lm := model.LedgerMember{
			LedgerID: ledger.ID,
			MemberID: memberID,
		}
		if err := pkg.DB.Create(&lm).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *LedgerRepo) Update(ledger *model.Ledger, memberIDs []uint) error {
	if err := pkg.DB.Save(ledger).Error; err != nil {
		return err
	}

	// Clean up and re-insert ledger_members
	if err := pkg.DB.Where("ledger_id = ?", ledger.ID).Delete(&model.LedgerMember{}).Error; err != nil {
		return err
	}

	for _, memberID := range memberIDs {
		lm := model.LedgerMember{
			LedgerID: ledger.ID,
			MemberID: memberID,
		}
		if err := pkg.DB.Create(&lm).Error; err != nil {
			return err
		}
	}
	return nil
}

func (r *LedgerRepo) FindByID(id uint) (*LedgerWithAssoc, error) {
	var ledger model.Ledger
	err := pkg.DB.
		Preload("Category").
		Preload("Creator").
		Preload("Members").
		First(&ledger, id).Error
	if err != nil {
		return nil, err
	}
	return &LedgerWithAssoc{
		Ledger:   ledger,
		Category: ledger.Category,
		Creator:  ledger.Creator,
		Members:  ledger.Members,
	}, nil
}

func (r *LedgerRepo) Delete(id uint) error {
	return pkg.DB.Delete(&model.Ledger{}, id).Error
}

// Ensure LedgerWithAssoc embeds include OccurredAt for date grouping
var _ time.Time = time.Time(model.Ledger{}.OccurredAt)
