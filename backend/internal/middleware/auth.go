package middleware

import (
	"strings"

	"home-center/internal/model"
	"home-center/internal/pkg"

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
		// 检查成员是否被禁用
		var member model.Member
		if err := pkg.DB.First(&member, claims.MemberID).Error; err != nil || member.Status == "disabled" {
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
