package handler

import (
	"net/http"
	"os"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

type seedLedgersRequest struct {
	Count     int    `json:"count" binding:"required"`
	StartDate string `json:"start_date"`
	EndDate   string `json:"end_date"`
}

type seedSummary struct {
	Income  int64 `json:"income"`
	Expense int64 `json:"expense"`
	Balance int64 `json:"balance"`
}

type seedResult struct {
	Count              int         `json:"count"`
	Summary            seedSummary `json:"summary"`
	ExpenseCategoryCnt int         `json:"expense_category_count"`
	IncomeCategoryCnt  int         `json:"income_category_count"`
}

// TestSeedLedgers 批量播种账单数据（仅测试模式可用）
func TestSeedLedgers(c *gin.Context) {
	if os.Getenv("HC_TEST_MODE") != "true" {
		pkg.Error(c, http.StatusNotFound, 404, "not found")
		return
	}

	var req seedLedgersRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Count <= 0 {
		pkg.Error(c, 400, 40001, "count 必须大于 0")
		return
	}

	db := pkg.DB

	// 解析日期范围，默认近 7 天
	endDate := time.Now()
	if req.EndDate != "" {
		var err error
		endDate, err = time.Parse("2006-01-02", req.EndDate)
		if err != nil {
			pkg.Error(c, 400, 40001, "end_date 格式错误，应为 YYYY-MM-DD")
			return
		}
	}

	startDate := endDate.AddDate(0, 0, -6) // 默认 7 天前
	if req.StartDate != "" {
		var err error
		startDate, err = time.Parse("2006-01-02", req.StartDate)
		if err != nil {
			pkg.Error(c, 400, 40001, "start_date 格式错误，应为 YYYY-MM-DD")
			return
		}
	}

	days := int(endDate.Sub(startDate).Hours()/24) + 1
	if days <= 0 {
		days = 1
	}

	// 获取预置分类
	var expenseCats, incomeCats []model.Category
	db.Where("type = ? AND deleted_at IS NULL", "expense").Order("sort_order").Find(&expenseCats)
	db.Where("type = ? AND deleted_at IS NULL", "income").Order("sort_order").Find(&incomeCats)

	if len(expenseCats) == 0 || len(incomeCats) == 0 {
		pkg.Error(c, 400, 40002, "请先调用 /api/test/reset 创建预置分类")
		return
	}

	// 获取第一个成员作为 creator
	var member model.Member
	if err := db.Where("deleted_at IS NULL").Order("id").First(&member).Error; err != nil {
		pkg.Error(c, 400, 40003, "请先创建至少一个成员")
		return
	}

	// 批量创建记录
	ledgers := make([]model.Ledger, 0, req.Count)
	expIdx := 0
	incIdx := 0

	for i := 0; i < req.Count; i++ {
		dayOffset := i * days / req.Count
		occurred := time.Date(
			startDate.Year(), startDate.Month(), startDate.Day()+dayOffset,
			12, 0, 0, 0, time.UTC,
		)

		var cat model.Category
		var amount int64

		if i%5 == 4 {
			// 每 5 条中 1 条收入
			cat = incomeCats[incIdx%len(incomeCats)]
			amount = 50000 // 500 元
			incIdx++
		} else {
			// 4 条支出
			cat = expenseCats[expIdx%len(expenseCats)]
			amount = int64((i%10 + 1) * 100) // 100~1000 分
			expIdx++
		}

		ledgers = append(ledgers, model.Ledger{
			Amount:     amount,
			Note:       "测试账单-" + time.Now().Format("150405") + "-" + string(rune('A'+i%26)),
			CategoryID: cat.ID,
			CreatorID:  member.ID,
			OccurredAt: model.FromTime(occurred),
		})
	}

	// 批量插入
	if err := db.CreateInBatches(&ledgers, 100).Error; err != nil {
		pkg.Error(c, 500, 50001, "播种数据失败: "+err.Error())
		return
	}

	// 计算 summary
	var income, expense int64
	db.Table("ledgers").
		Select("COALESCE(SUM(ledgers.amount), 0)").
		Joins("JOIN categories ON categories.id = ledgers.category_id").
		Where("categories.type = ? AND ledgers.deleted_at IS NULL", "income").
		Scan(&income)

	db.Table("ledgers").
		Select("COALESCE(SUM(ledgers.amount), 0)").
		Joins("JOIN categories ON categories.id = ledgers.category_id").
		Where("categories.type = ? AND ledgers.deleted_at IS NULL", "expense").
		Scan(&expense)

	pkg.Success(c, seedResult{
		Count: req.Count,
		Summary: seedSummary{
			Income:  income,
			Expense: expense,
			Balance: income - expense,
		},
		ExpenseCategoryCnt: len(expenseCats),
		IncomeCategoryCnt:  len(incomeCats),
	})
}
