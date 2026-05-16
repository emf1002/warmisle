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
}
