package handler

import (
	"strconv"
	"time"

	"warmisle/internal/pkg"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

// POST /api/likes
func (h *ForumHandler) ToggleLike(c *gin.Context) {
	var req struct {
		TargetType string `json:"target_type" binding:"required"`
		TargetID   uint   `json:"target_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	liked, err := h.svc.ToggleLike(req.TargetType, req.TargetID, getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumInvalidTargetType, 400, 40001, "无效的目标类型"},
		)
		return
	}
	pkg.Success(c, gin.H{"liked": liked})
}

// GET /api/votes
func (h *ForumHandler) ListVotes(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 {
		pageSize = 20
	}

	items, total, err := h.svc.ListVotes(page, pageSize, getMemberID(c))
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, gin.H{
		"items": items,
		"total": total,
		"page":  page,
		"page_size": pageSize,
	})
}

// POST /api/votes
func (h *ForumHandler) CreateVote(c *gin.Context) {
	var req struct {
		Title    string   `json:"title" binding:"required"`
		Options  []string `json:"options" binding:"required"`
		IsMulti  bool     `json:"is_multi"`
		Deadline *string  `json:"deadline"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	var deadline *time.Time
	if req.Deadline != nil {
		t, err := time.Parse(time.RFC3339, *req.Deadline)
		if err != nil {
			pkg.Error(c, 400, 40001, "日期格式错误")
			return
		}
		deadline = &t
	}

	result, err := h.svc.CreateVote(req.Title, req.Options, req.IsMulti, deadline, getMemberID(c))
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	pkg.Success(c, result)
}

// DELETE /api/votes/:id
func (h *ForumHandler) DeleteVote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.DeleteVote(uint(id), getMemberID(c), getRole(c)); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumVoteNotFound, 404, 40401, "投票不存在"},
			serviceError{service.ErrForumVoteDeadlinePassed, 400, 40001, "投票已截止，不能删除"},
			serviceError{service.ErrForumPermissionDenied, 403, 40301, "只能删除自己创建的投票"},
		)
		return
	}
	pkg.Success(c, nil)
}

// POST /api/votes/:id/vote
func (h *ForumHandler) Vote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req struct {
		OptionID uint `json:"option_id" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.Vote(uint(id), req.OptionID, getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumVoteNotFound, 404, 40401, "投票不存在"},
			serviceError{service.ErrForumVoteDeadlinePassed, 400, 40001, "投票已截止"},
			serviceError{service.ErrForumAlreadyVoted, 400, 40001, "已投票"},
		)
		return
	}
	pkg.Success(c, result)
}

// GET /api/votes/:id
func (h *ForumHandler) GetVote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.GetVote(uint(id), getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumVoteNotFound, 404, 40401, "投票不存在"},
		)
		return
	}
	pkg.Success(c, result)
}
