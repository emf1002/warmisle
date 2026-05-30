package handler

import (
	"errors"

	"github.com/gin-gonic/gin"
	"warmisle/internal/pkg"
)

// serviceError maps a service error to an HTTP response.
type serviceError struct {
	err      error
	httpCode int
	bizCode  int
	msg      string
}

// handleServiceError checks err against mappings and returns the first match.
// If no match, returns 500.
func handleServiceError(c *gin.Context, err error, mappings ...serviceError) {
	for _, m := range mappings {
		if errors.Is(err, m.err) {
			pkg.Error(c, m.httpCode, m.bizCode, m.msg)
			return
		}
	}
	pkg.Error(c, 500, 50001, "服务器内部错误")
}

// getMemberID extracts member_id from Gin context (set by JWT middleware).
func getMemberID(c *gin.Context) uint {
	val, _ := c.Get("member_id")
	return val.(uint)
}

// getRole extracts role from Gin context (set by JWT middleware).
func getRole(c *gin.Context) string {
	val, _ := c.Get("role")
	return val.(string)
}
