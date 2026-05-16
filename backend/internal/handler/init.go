package handler

import (
	"home-center/internal/pkg"
	"home-center/internal/service"

	"github.com/gin-gonic/gin"
)

type InitHandler struct {
	svc *service.InitService
}

func NewInitHandler() *InitHandler {
	return &InitHandler{svc: &service.InitService{}}
}

type setupRequest struct {
	AdminName string `json:"admin_name" binding:"required"`
	Username  string `json:"username" binding:"required"`
	Password  string `json:"password" binding:"required"`
}

func (h *InitHandler) Setup(c *gin.Context) {
	var req setupRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数不能为空")
		return
	}

	// 检查是否需要初始化
	authSvc := service.NewAuthService()
	needInit, err := authSvc.InitCheck()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	if !needInit {
		pkg.Error(c, 400, 40001, "系统已初始化")
		return
	}

	// 在事务中创建管理员
	tx := pkg.DB.Begin()
	if tx.Error != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	admin, err := h.svc.Setup(tx, req.AdminName, req.Username, req.Password)
	if err != nil {
		tx.Rollback()
		pkg.Error(c, 500, 50001, "创建管理员失败")
		return
	}

	if err := tx.Commit().Error; err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	// 生成 token
	token, err := pkg.GenerateToken(admin.ID, admin.Username, admin.Role)
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, gin.H{
		"token": token,
		"member": gin.H{
			"id":       admin.ID,
			"username": admin.Username,
			"name":     admin.Name,
			"role":     admin.Role,
		},
	})
}
