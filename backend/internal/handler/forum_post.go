package handler

import (
	"strconv"

	"warmisle/internal/pkg"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

// POST /api/posts
func (h *ForumHandler) CreatePost(c *gin.Context) {
	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.CreatePost(req.Content, getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumContentRequired, 400, 40001, "内容不能为空"},
		)
		return
	}
	pkg.Success(c, result)
}

// PUT /api/posts/:id
func (h *ForumHandler) UpdatePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req struct {
		Content string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.UpdatePost(uint(id), req.Content, getMemberID(c), getRole(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumPostNotFound, 404, 40401, "动态不存在"},
			serviceError{service.ErrForumPermissionDenied, 403, 40301, "只能编辑自己的动态"},
			serviceError{service.ErrForumContentRequired, 400, 40001, "内容不能为空"},
		)
		return
	}
	pkg.Success(c, result)
}

// DELETE /api/posts/:id
func (h *ForumHandler) DeletePost(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.DeletePost(uint(id), getMemberID(c), getRole(c)); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumPostNotFound, 404, 40401, "动态不存在"},
			serviceError{service.ErrForumPermissionDenied, 403, 40301, "只能删除自己的动态"},
		)
		return
	}
	pkg.Success(c, nil)
}
