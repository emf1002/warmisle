package middleware

import (
	"os"
	"strings"

	"github.com/gin-gonic/gin"
)

// CORS 跨域资源共享中间件
// 生产环境（嵌入前端）：仅允许同源，不设置 CORS 头（浏览器默认同源策略）
// 开发环境：读取 HC_CORS_ORIGINS 环境变量，支持白名单配置
func CORS() gin.HandlerFunc {
	// 读取配置
	originsEnv := os.Getenv("HC_CORS_ORIGINS")

	// 生产模式：嵌入前端时仅同源访问，不设置 CORS 头
	// 注意：若前端使用 Vite dev server (port 3000) 代理 /api，需设置 HC_CORS_ORIGINS=http://localhost:3000
	if originsEnv == "" {
		// 同源模式：不添加任何 CORS 响应头
		return func(c *gin.Context) {
			c.Next()
		}
	}

	// 开发模式：解析白名单
	allowedOrigins := make(map[string]bool)
	if originsEnv == "*" {
		allowedOrigins["*"] = true
	} else {
		for _, origin := range strings.Split(originsEnv, ",") {
			allowedOrigins[strings.TrimSpace(origin)] = true
		}
	}

	return func(c *gin.Context) {
		origin := c.GetHeader("Origin")
		if origin == "" {
			c.Next()
			return
		}

		// 检查白名单
		allowOrigin := ""
		if allowedOrigins["*"] {
			allowOrigin = "*"
		} else if allowedOrigins[origin] {
			allowOrigin = origin
		}

		if allowOrigin == "" {
			// 不在白名单中，不添加 CORS 头，浏览器将阻止
			c.Next()
			return
		}

		c.Header("Access-Control-Allow-Origin", allowOrigin)
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Content-Type, Authorization")
		c.Header("Access-Control-Max-Age", "86400") // 24小时

		// 预检请求直接返回
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}
