package repository

import (
	"encoding/base64"
	"encoding/json"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"gorm.io/gorm"
)

// CursorData encodes pagination position as (occurred_at, id) tuple.
type CursorData struct {
	OccurredAt string `json:"occurred_at"` // "2006-01-02 15:04:05"
	ID         uint   `json:"id"`
}

// EncodeCursor base64-encodes a CursorData.
func EncodeCursor(c CursorData) string {
	b, _ := json.Marshal(c)
	return base64.StdEncoding.EncodeToString(b)
}

// DecodeCursor decodes a base64 cursor string.
func DecodeCursor(s string) (*CursorData, error) {
	b, err := base64.StdEncoding.DecodeString(s)
	if err != nil {
		return nil, err
	}
	var c CursorData
	if err := json.Unmarshal(b, &c); err != nil {
		return nil, err
	}
	return &c, nil
}

type LedgerRepo struct{}

type LedgerFilter struct {
	StartDate  string // "2026-05-01"
	EndDate    string // "2026-06-01" (exclusive upper bound)
	CategoryID *uint
	CreatorID  *uint
	Limit      int         // page size, default 20
	Cursor     *CursorData // nil = first page
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
}

type ListResult struct {
	Summary    LedgerSummary `json:"summary"`
	Groups     []LedgerGroup `json:"groups"`
	NextCursor *string       `json:"next_cursor"`
	HasMore    bool          `json:"has_more"`
}

type LedgerSummary struct {
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
	Balance int64 `json:"balance"`
}

// applyOptionalFilters appends category_id / creator_id conditions when set.
func (r *LedgerRepo) applyOptionalFilters(query *gorm.DB, filter LedgerFilter) *gorm.DB {
	if filter.CategoryID != nil {
		query = query.Where("ledgers.category_id = ?", *filter.CategoryID)
	}
	if filter.CreatorID != nil {
		query = query.Where("ledgers.creator_id = ?", *filter.CreatorID)
	}
	return query
}

// computeSummary returns income/expense totals with the same filters as the list query.
func (r *LedgerRepo) computeSummary(filter LedgerFilter) (LedgerSummary, error) {
	var summary LedgerSummary
	type summaryRow struct {
		Type   string
		Amount int64
	}
	var rows []summaryRow

	query := pkg.DB.Table("ledgers").
		Select("categories.type, COALESCE(SUM(ledgers.amount), 0) as amount").
		Joins("JOIN categories ON ledgers.category_id = categories.id AND categories.deleted_at IS NULL").
		Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", filter.StartDate, filter.EndDate)
	query = r.applyOptionalFilters(query, filter)

	if err := query.Group("categories.type").Scan(&rows).Error; err != nil {
		return summary, err
	}
	for _, row := range rows {
		if row.Type == "income" {
			summary.Income = row.Amount
		} else if row.Type == "expense" {
			summary.Expense = row.Amount
		}
	}
	summary.Balance = summary.Income - summary.Expense
	return summary, nil
}

// calcDailyTotal sums income (+) and expense (-) for a set of items.
func (r *LedgerRepo) calcDailyTotal(items []LedgerWithAssoc) int64 {
	var total int64
	for _, item := range items {
		if item.Category.Type == "income" {
			total += item.Amount
		} else {
			total -= item.Amount
		}
	}
	return total
}

func (r *LedgerRepo) List(filter LedgerFilter) (*ListResult, error) {
	// 1. Compute summary (now applies category/creator filters)
	summary, err := r.computeSummary(filter)
	if err != nil {
		return nil, err
	}

	// 2. Determine limit
	limit := filter.Limit
	if limit <= 0 {
		limit = 20
	}

	// 3. Build base query
	query := pkg.DB.Model(&model.Ledger{}).
		Preload("Category").
		Preload("Creator").
		Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", filter.StartDate, filter.EndDate)
	query = r.applyOptionalFilters(query, filter)

	// 4. Apply cursor condition
	if filter.Cursor != nil {
		query = query.Where(
			"(ledgers.occurred_at < ? OR (ledgers.occurred_at = ? AND ledgers.id < ?))",
			filter.Cursor.OccurredAt, filter.Cursor.OccurredAt, filter.Cursor.ID,
		)
	}

	// 5. Fetch limit + 1 records
	var records []model.Ledger
	if err := query.
		Order("ledgers.occurred_at DESC, ledgers.id DESC").
		Limit(limit + 1).
		Find(&records).Error; err != nil {
		return nil, err
	}

	// 6. Determine hasMore and keep the (limit+1)th record for date check
	hasMore := len(records) > limit
	var extraRecord *model.Ledger
	if hasMore {
		extraRecord = &records[limit]
		records = records[:limit]
	}

	// 7. Group by date
	groups := r.groupByDate(records)

	// 8. If last date group was cut off, fetch remaining records (补全)
	if hasMore && len(groups) > 0 && extraRecord != nil {
		lastGroup := &groups[len(groups)-1]
		extraDate := time.Time(extraRecord.OccurredAt).Format("2006-01-02")
		if extraDate == lastGroup.Date {
			minId := lastGroup.Items[len(lastGroup.Items)-1].ID
			nextDay := time.Time(extraRecord.OccurredAt).AddDate(0, 0, 1).Format("2006-01-02")

		补全Q := pkg.DB.Model(&model.Ledger{}).
				Preload("Category").
				Preload("Creator").
				Where("ledgers.occurred_at >= ? AND ledgers.occurred_at < ?", lastGroup.Date, nextDay).
				Where("ledgers.id < ?", minId)
			补全Q = r.applyOptionalFilters(补全Q, filter)

			var extraRecords []model.Ledger
			if err := 补全Q.
				Order("ledgers.occurred_at DESC, ledgers.id DESC").
				Find(&extraRecords).Error; err != nil {
				return nil, err
			}

			for _, er := range extraRecords {
				lastGroup.Items = append(lastGroup.Items, LedgerWithAssoc{
					Ledger:   er,
					Category: er.Category,
					Creator:  er.Creator,
				})
			}
			lastGroup.DailyTotal = r.calcDailyTotal(lastGroup.Items)
		}
	}

	// 9. Compute nextCursor
	var nextCursor *string
	if hasMore {
		var lastItem LedgerWithAssoc
		if len(groups) > 0 {
			lastGroup := groups[len(groups)-1]
			lastItem = lastGroup.Items[len(lastGroup.Items)-1]
		}
		cursor := EncodeCursor(CursorData{
			OccurredAt: time.Time(lastItem.OccurredAt).Format("2006-01-02 15:04:05"),
			ID:         lastItem.ID,
		})
		nextCursor = &cursor
	}

	return &ListResult{
		Summary:    summary,
		Groups:     groups,
		NextCursor: nextCursor,
		HasMore:    hasMore,
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

func (r *LedgerRepo) Create(ledger *model.Ledger) error {
	return pkg.DB.Create(ledger).Error
}

func (r *LedgerRepo) Update(ledger *model.Ledger) error {
	return pkg.DB.Save(ledger).Error
}

func (r *LedgerRepo) FindByID(id uint) (*LedgerWithAssoc, error) {
	var ledger model.Ledger
	err := pkg.DB.
		Preload("Category").
		Preload("Creator").
		First(&ledger, id).Error
	if err != nil {
		return nil, err
	}
	return &LedgerWithAssoc{
		Ledger:   ledger,
		Category: ledger.Category,
		Creator:  ledger.Creator,
	}, nil
}

func (r *LedgerRepo) Delete(id uint) error {
	return pkg.DB.Delete(&model.Ledger{}, id).Error
}

// Ensure LedgerWithAssoc embeds include OccurredAt for date grouping
var _ time.Time = time.Time(model.Ledger{}.OccurredAt)
