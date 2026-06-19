// Package middleware provides HTTP middleware for warmisle.
package middleware

import (
	"os"
	"sync"
	"time"

	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

// rateLimitBucket 滑动窗口计数器
type rateLimitBucket struct {
	count    int
	resetAt  time.Time
}

// RateLimiter 基于 IP 的滑动窗口速率限制器
type RateLimiter struct {
	mu       sync.Mutex
	buckets  map[string]*rateLimitBucket
	limit    int           // 窗口内最大请求数
	window   time.Duration // 时间窗口
	cleanupInterval time.Duration
}

// NewRateLimiter 创建一个新的速率限制器
func NewRateLimiter(limit int, window time.Duration) *RateLimiter {
	rl := &RateLimiter{
		buckets:  make(map[string]*rateLimitBucket),
		limit:    limit,
		window:   window,
		cleanupInterval: 10 * time.Minute,
	}
	// 后台清理过期桶
	go rl.cleanup()
	return rl
}

// allow 检查指定 key 是否允许通过
func (rl *RateLimiter) allow(key string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	bucket, exists := rl.buckets[key]

	// 窗口已过期，重置
	if !exists || now.After(bucket.resetAt) {
		rl.buckets[key] = &rateLimitBucket{
			count:   1,
			resetAt: now.Add(rl.window),
		}
		return true
	}

	// 窗口内检查
	if bucket.count >= rl.limit {
		return false
	}

	bucket.count++
	return true
}

// cleanup 定期清理过期的桶
func (rl *RateLimiter) cleanup() {
	ticker := time.NewTicker(rl.cleanupInterval)
	defer ticker.Stop()
	for range ticker.C {
		rl.mu.Lock()
		now := time.Now()
		for key, bucket := range rl.buckets {
			if now.After(bucket.resetAt) {
				delete(rl.buckets, key)
			}
		}
		rl.mu.Unlock()
	}
}

// getClientIP 获取客户端真实 IP
func getClientIP(c *gin.Context) string {
	// 优先从 X-Forwarded-For 获取（反向代理场景）
	if ip := c.GetHeader("X-Forwarded-For"); ip != "" {
		return ip
	}
	if ip := c.GetHeader("X-Real-IP"); ip != "" {
		return ip
	}
	return c.ClientIP()
}

// 全局速率限制器实例
var (
	// 通用 API 限制：100 请求/分钟
	globalLimiter = NewRateLimiter(100, 1*time.Minute)
	// 登录接口严格限制：5 请求/分钟
	loginLimiter = NewRateLimiter(5, 1*time.Minute)
)

// RateLimit 通用速率限制中间件
func RateLimit() gin.HandlerFunc {
	// 测试模式下跳过限流
	if os.Getenv("HC_TEST_MODE") == "true" {
		return func(c *gin.Context) { c.Next() }
	}
	return func(c *gin.Context) {
		ip := getClientIP(c)
		if !globalLimiter.allow(ip) {
			pkg.Error(c, 429, 42901, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}

// StrictRateLimit 严格速率限制中间件（用于登录等敏感端点）
func StrictRateLimit() gin.HandlerFunc {
	// 测试模式下跳过限流
	if os.Getenv("HC_TEST_MODE") == "true" {
		return func(c *gin.Context) { c.Next() }
	}
	return func(c *gin.Context) {
		ip := getClientIP(c)
		if !loginLimiter.allow(ip) {
			pkg.Error(c, 429, 42901, "请求过于频繁，请稍后再试")
			c.Abort()
			return
		}
		c.Next()
	}
}
