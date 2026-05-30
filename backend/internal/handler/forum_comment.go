package handler

import (
	"strconv"

	"warmisle/internal/pkg"
	"warmisle/internal/repository"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

// GET /api/comments?target_type=&target_id=
func (h *ForumHandler) ListComments(c *gin.Context) {
	targetType := c.Query("target_type")
	targetIDStr := c.Query("target_id")
	if targetType == "" || targetIDStr == "" {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}
	targetID, err := strconv.ParseUint(targetIDStr, 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.ListComments(targetType, uint(targetID))
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	if result == nil {
		result = []repository.CommentWithMeta{}
	}
	pkg.Success(c, result)
}

// POST /api/comments
func (h *ForumHandler) CreateComment(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
		ParentID   *uint  `json:"parent_id"`
		Content    string `json:"content" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.CreateComment(req.TargetType, req.TargetID, req.ParentID, req.Content, getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumInvalidTargetType, 400, 40001, "无效的评论目标类型"},
			serviceError{service.ErrForumContentRequired, 400, 40001, "内容不能为空"},
			serviceError{service.ErrForumCommentNotFound, 404, 40401, "父评论不存在"},
			serviceError{service.ErrForumNestingTooDeep, 400, 40001, "不能回复二级评论"},
		)
		return
	}
	pkg.Success(c, result)
}

// DELETE /api/comments/:id
func (h *ForumHandler) DeleteComment(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.DeleteComment(uint(id), getMemberID(c), getRole(c)); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumCommentNotFound, 404, 40401, "评论不存在"},
			serviceError{service.ErrForumPermissionDenied, 403, 40301, "只能删除自己的评论"},
		)
		return
	}
	pkg.Success(c, nil)
}
