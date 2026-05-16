package handler

import (
	"time"

	"home-center/internal/pkg"
	"home-center/internal/service"

	"github.com/gin-gonic/gin"
)

type DashboardHandler struct {
	svc *service.DashboardService
}

func NewDashboardHandler() *DashboardHandler {
	return &DashboardHandler{svc: service.NewDashboardService()}
}

func getDefaultMonth(c *gin.Context) string {
	month := c.Query("month")
	if month == "" {
		month = time.Now().Format("2006-01")
	}
	return month
}

// GET /api/dashboard/summary?month=2026-05
func (h *DashboardHandler) Summary(c *gin.Context) {
	month := getDefaultMonth(c)
	result, err := h.svc.GetSummary(month)
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, result)
}

// GET /api/dashboard/expense-chart?month=2026-05
func (h *DashboardHandler) ExpenseChart(c *gin.Context) {
	month := getDefaultMonth(c)
	result, err := h.svc.GetExpenseChart(month)
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, result)
}

// GET /api/dashboard/upcoming-todos
func (h *DashboardHandler) UpcomingTodos(c *gin.Context) {
	result, err := h.svc.GetUpcomingTodos()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, result)
}

// GET /api/dashboard/wish-trends
func (h *DashboardHandler) WishTrends(c *gin.Context) {
	result, err := h.svc.GetWishTrends()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, result)
}

// GET /api/dashboard/forum-hot
func (h *DashboardHandler) ForumHot(c *gin.Context) {
	result, err := h.svc.GetForumHot()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, result)
}
