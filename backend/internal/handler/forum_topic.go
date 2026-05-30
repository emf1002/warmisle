package handler

import (
	"strconv"

	"warmisle/internal/pkg"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

// POST /api/topics
func (h *ForumHandler) CreateTopic(c *gin.Context) {
	var req struct {
		Title   string `json:"title" binding:"required"`
		Content string `json:"content"`
		TagID   *uint  `json:"tag_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.CreateTopic(req.Title, req.Content, req.TagID, getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumTitleRequired, 400, 40001, "标题不能为空"},
		)
		return
	}
	pkg.Success(c, result)
}

// PUT /api/topics/:id
func (h *ForumHandler) UpdateTopic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req struct {
		Title   *string `json:"title"`
		Content *string `json:"content"`
		TagID   *uint   `json:"tag_id"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.UpdateTopic(uint(id), req.Title, req.Content, req.TagID, getMemberID(c), getRole(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumTopicNotFound, 404, 40401, "话题不存在"},
			serviceError{service.ErrForumPermissionDenied, 403, 40301, "只能编辑自己的话题"},
			serviceError{service.ErrForumTitleRequired, 400, 40001, "标题不能为空"},
		)
		return
	}
	pkg.Success(c, result)
}

// DELETE /api/topics/:id
func (h *ForumHandler) DeleteTopic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.DeleteTopic(uint(id), getMemberID(c), getRole(c)); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumTopicNotFound, 404, 40401, "话题不存在"},
			serviceError{service.ErrForumPermissionDenied, 403, 40301, "只能删除自己的话题"},
		)
		return
	}
	pkg.Success(c, nil)
}

// PUT /api/topics/:id/pin
func (h *ForumHandler) TogglePin(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.TogglePin(uint(id), getRole(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumTopicNotFound, 404, 40401, "话题不存在"},
			serviceError{service.ErrForumPermissionDenied, 403, 40301, "权限不足"},
		)
		return
	}
	pkg.Success(c, result)
}

// GET /api/topics/:id
func (h *ForumHandler) GetTopic(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.GetTopic(uint(id), getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumTopicNotFound, 404, 40401, "话题不存在"},
		)
		return
	}
	pkg.Success(c, result)
}
