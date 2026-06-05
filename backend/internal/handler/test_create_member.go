package handler

import (
	"net/http"
	"os"
	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

type createTestMemberRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name" binding:"required"`
}

// TestCreateMember 创建测试用 member 用户并返回 JWT token
// 仅在 HC_TEST_MODE=true 时可用
func TestCreateMember(c *gin.Context) {
	if os.Getenv("HC_TEST_MODE") != "true" {
		pkg.Error(c, http.StatusNotFound, 404, "not found")
		return
	}

	var req createTestMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 400, "参数错误")
		return
	}

	db := pkg.DB

	// 检查用户名是否已存在
	var count int64
	db.Model(&model.Member{}).Where("username = ?", req.Username).Count(&count)
	if count > 0 {
		pkg.Error(c, http.StatusConflict, 40901, "用户名已被使用")
		return
	}

	// 加密密码
	hashedPwd, err := pkg.HashPassword(req.Password)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "密码加密失败")
		return
	}

	// 创建 member
	member := model.Member{
		Username: req.Username,
		Password: hashedPwd,
		Name:     req.Name,
		Role:     "member",
	}
	if err := db.Create(&member).Error; err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "创建失败")
		return
	}

	// 签发 token
	token, err := pkg.GenerateToken(member.ID, member.Username, member.Role)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 500, "token 签发失败")
		return
	}

	pkg.Success(c, gin.H{
		"token":     token,
		"member_id": member.ID,
	})
}
