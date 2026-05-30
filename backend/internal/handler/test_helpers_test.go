package handler

import (
	"warmisle/internal/middleware"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/gin-gonic/gin"
)

func initJWT() {
	pkg.InitJWT("test-secret-for-handler-tests")
}

func setupMemberTest() {
	testutil.SetupTestDB()
	initJWT()
}

// setupTestRouter creates a full Gin router with all endpoints and middleware.
// Note: This mirrors routes.Register() but cannot import routes directly
// due to the circular dependency (routes -> handler -> routes).
func setupTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")

	authH := NewAuthHandler()
	initH := NewInitHandler()
	memberH := NewMemberHandler()
	categoryH := NewCategoryHandler()
	ledgerH := NewLedgerHandler()
	dashboardH := NewDashboardHandler()
	todoH := NewTodoHandler()
	wishH := NewWishHandler()
	forumH := NewForumHandler()
	tagH := NewTagHandler()

	// Public
	api.GET("/init/check", authH.InitCheck)
	api.POST("/init/setup", initH.Setup)
	api.POST("/auth/login", authH.Login)

	authRequired := middleware.AuthRequired()
	adminRequired := middleware.AdminRequired()

	// Members
	api.GET("/members", authRequired, memberH.List)
	api.POST("/members", authRequired, adminRequired, memberH.Create)
	api.PUT("/members/:id", authRequired, adminRequired, memberH.Update)
	api.DELETE("/members/:id", authRequired, adminRequired, memberH.Delete)
	api.PUT("/members/:id/disable", authRequired, adminRequired, memberH.Disable)
	api.PUT("/members/:id/enable", authRequired, adminRequired, memberH.Enable)
	api.PUT("/members/:id/reset-pwd", authRequired, adminRequired, memberH.ResetPassword)

	// Profile
	api.GET("/profile", authRequired, memberH.GetProfile)
	api.PUT("/profile", authRequired, memberH.UpdateProfile)
	api.PUT("/profile/password", authRequired, memberH.ChangePassword)

	// Categories
	api.GET("/categories", authRequired, categoryH.List)
	api.POST("/categories", authRequired, adminRequired, categoryH.Create)
	api.PUT("/categories/:id", authRequired, adminRequired, categoryH.Update)
	api.DELETE("/categories/:id", authRequired, adminRequired, categoryH.Delete)

	// Ledger
	api.GET("/ledgers", authRequired, ledgerH.List)
	api.POST("/ledgers", authRequired, ledgerH.Create)
	api.GET("/ledgers/:id", authRequired, ledgerH.GetByID)
	api.PUT("/ledgers/:id", authRequired, ledgerH.Update)
	api.DELETE("/ledgers/:id", authRequired, ledgerH.Delete)

	// Dashboard
	api.GET("/dashboard/summary", authRequired, dashboardH.Summary)
	api.GET("/dashboard/expense-chart", authRequired, dashboardH.ExpenseChart)
	api.GET("/dashboard/upcoming-todos", authRequired, dashboardH.UpcomingTodos)
	api.GET("/dashboard/wish-trends", authRequired, dashboardH.WishTrends)
	api.GET("/dashboard/forum-hot", authRequired, dashboardH.ForumHot)

	// Todos
	api.GET("/todos", authRequired, todoH.List)
	api.POST("/todos", authRequired, todoH.Create)
	api.PUT("/todos/:id", authRequired, todoH.Update)
	api.DELETE("/todos/:id", authRequired, todoH.Delete)
	api.PUT("/todos/:id/toggle", authRequired, todoH.Toggle)
	api.PUT("/todos/:id/claim", authRequired, todoH.Claim)

	// Wishes
	api.GET("/wishes", authRequired, wishH.List)
	api.POST("/wishes", authRequired, wishH.Create)
	api.PUT("/wishes/:id", authRequired, wishH.Update)
	api.DELETE("/wishes/:id", authRequired, wishH.Delete)
	api.POST("/wishes/:id/promote", authRequired, wishH.Promote)
	api.PUT("/wishes/:id/status", authRequired, wishH.UpdateStatus)
	api.POST("/wishes/:id/vote", authRequired, wishH.Vote)
	api.DELETE("/wishes/:id/vote", authRequired, wishH.Unvote)

	// Forum
	api.GET("/feed", authRequired, forumH.Feed)
	api.POST("/posts", authRequired, forumH.CreatePost)
	api.PUT("/posts/:id", authRequired, forumH.UpdatePost)
	api.DELETE("/posts/:id", authRequired, forumH.DeletePost)
	api.POST("/topics", authRequired, forumH.CreateTopic)
	api.PUT("/topics/:id", authRequired, forumH.UpdateTopic)
	api.DELETE("/topics/:id", authRequired, forumH.DeleteTopic)
	api.PUT("/topics/:id/pin", authRequired, adminRequired, forumH.TogglePin)
	api.GET("/topics/:id", authRequired, forumH.GetTopic)
	api.GET("/comments", authRequired, forumH.ListComments)
	api.POST("/comments", authRequired, forumH.CreateComment)
	api.DELETE("/comments/:id", authRequired, forumH.DeleteComment)
	api.POST("/likes", authRequired, forumH.ToggleLike)
	api.POST("/votes", authRequired, forumH.CreateVote)
	api.DELETE("/votes/:id", authRequired, forumH.DeleteVote)
	api.POST("/votes/:id/vote", authRequired, forumH.Vote)
	api.GET("/votes/:id", authRequired, forumH.GetVote)

	// Tags
	api.GET("/tags", authRequired, tagH.List)
	api.POST("/tags", authRequired, adminRequired, tagH.Create)
	api.PUT("/tags/:id", authRequired, adminRequired, tagH.Update)
	api.DELETE("/tags/:id", authRequired, adminRequired, tagH.Delete)

	return r
}
