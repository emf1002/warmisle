package middleware

import (
	"strings"

	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

func AuthRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		auth := c.GetHeader("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			pkg.Error(c, 401, 40101, "未登录")
			c.Abort()
			return
		}
		claims, err := pkg.ParseToken(auth[7:])
		if err != nil {
			pkg.Error(c, 401, 40101, "Token 已过期")
			c.Abort()
			return
		}
		// 只查询 status 字段检查是否被禁用（避免加载整个 member 记录）
		var status string
		if err := pkg.DB.Model(&model.Member{}).Select("status").Where("id = ?", claims.MemberID).Scan(&status).Error; err != nil || status == "disabled" {
			pkg.Error(c, 403, 40301, "账号已被禁用")
			c.Abort()
			return
		}
		c.Set("member_id", claims.MemberID)
		c.Set("role", claims.Role)
		c.Next()
	}
}

func AdminRequired() gin.HandlerFunc {
	return func(c *gin.Context) {
		role, _ := c.Get("role")
		if role != "admin" {
			pkg.Error(c, 403, 40301, "权限不足")
			c.Abort()
			return
		}
		c.Next()
	}
}
