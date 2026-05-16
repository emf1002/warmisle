package routes

import (
	"home-center/internal/handler"

	"github.com/gin-gonic/gin"
)

func Register(r *gin.Engine) {
	api := r.Group("/api")

	auth := handler.NewAuthHandler()
	init := handler.NewInitHandler()

	api.GET("/init/check", auth.InitCheck)
	api.POST("/init/setup", init.Setup)
	api.POST("/auth/login", auth.Login)
}
