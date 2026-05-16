package routes

import (
	"home-center/internal/handler"
	"home-center/internal/middleware"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	auth := handler.NewAuthHandler()
	init := handler.NewInitHandler()
	member := handler.NewMemberHandler()
	category := handler.NewCategoryHandler()
	ledger := handler.NewLedgerHandler()
	dashboard := handler.NewDashboardHandler()
	todo := handler.NewTodoHandler()
	wish := handler.NewWishHandler()

	// Auth & Init
	api.GET("/init/check", auth.InitCheck)
	api.POST("/init/setup", init.Setup)
	api.POST("/auth/login", auth.Login)

	authRequired := middleware.AuthRequired()
	adminRequired := middleware.AdminRequired()

	// Member management
	api.GET("/members", authRequired, member.List)
	api.POST("/members", authRequired, adminRequired, member.Create)
	api.PUT("/members/:id", authRequired, adminRequired, member.Update)
	api.DELETE("/members/:id", authRequired, adminRequired, member.Delete)
	api.PUT("/members/:id/disable", authRequired, adminRequired, member.Disable)
	api.PUT("/members/:id/enable", authRequired, adminRequired, member.Enable)
	api.PUT("/members/:id/reset-pwd", authRequired, adminRequired, member.ResetPassword)

	// Profile
	api.GET("/profile", authRequired, member.GetProfile)
	api.PUT("/profile", authRequired, member.UpdateProfile)

	// Category management
	api.GET("/categories", authRequired, category.List)
	api.POST("/categories", authRequired, adminRequired, category.Create)
	api.PUT("/categories/:id", authRequired, adminRequired, category.Update)
	api.DELETE("/categories/:id", authRequired, adminRequired, category.Delete)

	// Ledger management
	api.GET("/ledgers", authRequired, ledger.List)
	api.POST("/ledgers", authRequired, ledger.Create)
	api.GET("/ledgers/:id", authRequired, ledger.GetByID)
	api.PUT("/ledgers/:id", authRequired, ledger.Update)
	api.DELETE("/ledgers/:id", authRequired, ledger.Delete)

	// Dashboard
	api.GET("/dashboard/summary", authRequired, dashboard.Summary)
	api.GET("/dashboard/expense-chart", authRequired, dashboard.ExpenseChart)
	api.GET("/dashboard/upcoming-todos", authRequired, dashboard.UpcomingTodos)
	api.GET("/dashboard/wish-trends", authRequired, dashboard.WishTrends)
	api.GET("/dashboard/forum-hot", authRequired, dashboard.ForumHot)

	// Todo management
	api.GET("/todos", authRequired, todo.List)
	api.POST("/todos", authRequired, todo.Create)
	api.PUT("/todos/:id", authRequired, todo.Update)
	api.DELETE("/todos/:id", authRequired, todo.Delete)
	api.PUT("/todos/:id/toggle", authRequired, todo.Toggle)
	api.PUT("/todos/:id/claim", authRequired, todo.Claim)

	// Wish management
	api.GET("/wishes", authRequired, wish.List)
	api.POST("/wishes", authRequired, wish.Create)
	api.PUT("/wishes/:id", authRequired, wish.Update)
	api.DELETE("/wishes/:id", authRequired, wish.Delete)
	api.POST("/wishes/:id/promote", authRequired, wish.Promote)
	api.PUT("/wishes/:id/status", authRequired, wish.UpdateStatus)
	api.POST("/wishes/:id/vote", authRequired, wish.Vote)
	api.DELETE("/wishes/:id/vote", authRequired, wish.Unvote)
}
