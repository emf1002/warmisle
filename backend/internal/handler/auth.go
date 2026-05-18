package handler

import (
	"errors"

	"warmisle/internal/pkg"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

type AuthHandler struct {
	svc *service.AuthService
}

func NewAuthHandler() *AuthHandler {
	return &AuthHandler{svc: service.NewAuthService()}
}

type loginRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

func (h *AuthHandler) Login(c *gin.Context) {
	var req loginRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "用户名和密码不能为空")
		return
	}

	token, err := h.svc.Login(req.Username, req.Password)
	if err != nil {
		if errors.Is(err, service.ErrAccountLocked) {
			pkg.Error(c, 429, 40102, "账号已被锁定，请15分钟后重试")
			return
		}
		if errors.Is(err, service.ErrInvalidCredentials) {
			pkg.Error(c, 401, 40101, "用户名或密码错误")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) InitCheck(c *gin.Context) {
	needInit, err := h.svc.InitCheck()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, gin.H{"need_init": needInit})
}
