package handler

import (
	"strconv"

	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

// GET /api/feed
func (h *ForumHandler) Feed(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	result, err := h.svc.GetFeed(page, pageSize)
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, result)
}
