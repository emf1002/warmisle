package routes

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
	api := r.Group("/api")
	// 各模块路由后续注册
	_ = api
}
