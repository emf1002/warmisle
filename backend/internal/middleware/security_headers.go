// Package middleware provides HTTP middleware for warmisle.
package middleware

import (
	"github.com/gin-gonic/gin"
)

// SecurityHeaders 设置安全相关的 HTTP 响应头
func SecurityHeaders() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 防止 MIME 类型嗅探
		c.Header("X-Content-Type-Options", "nosniff")

		// 点击劫持防护
		c.Header("X-Frame-Options", "DENY")

		// 旧版浏览器 XSS 过滤器
		c.Header("X-XSS-Protection", "1; mode=block")

		// HTTP 严格传输安全（1年 + 子域名 + 预加载）
		// 注意：无 HTTPS 部署时，此头可能造成问题。如无 HTTPS，建议移除或设为 max-age=0
		c.Header("Strict-Transport-Security", "max-age=31536000; includeSubDomains; preload")

		// 内容安全策略：仅允许同源资源
		// 如果前端需要加载外部资源（CDN、图片等），需相应调整
		csp := "default-src 'self'; " +
			"script-src 'self'; " +
			"style-src 'self' 'unsafe-inline'; " +
			"img-src 'self' data: https:; " +
			"font-src 'self' data:; " +
			"connect-src 'self'; " +
			"frame-ancestors 'none'; " +
			"base-uri 'self'; " +
			"form-action 'self'"
		c.Header("Content-Security-Policy", csp)

		// Referrer 策略：同源发送完整 URL，跨域仅发送源
		c.Header("Referrer-Policy", "strict-origin-when-cross-origin")

		// 权限策略：禁用不必要 API
		c.Header("Permissions-Policy", "camera=(), microphone=(), geolocation=(), payment=()")

		// 移除服务器版本信息（Gin 自动处理，这里显式设置）
		c.Header("Server", "")

		c.Next()
	}
}
