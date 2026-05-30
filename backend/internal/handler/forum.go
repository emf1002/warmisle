package handler

import (
	"strconv"
	"time"

	"warmisle/internal/pkg"
	"warmisle/internal/repository"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

type ForumHandler struct {
	svc *service.ForumService
}

func NewForumHandler() *ForumHandler {
	return &ForumHandler{svc: service.NewForumService()}
}

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

// POST /api/votes
func (h *ForumHandler) CreateVote(c *gin.Context) {
	var req struct {
		Title   string    `json:"title" binding:"required"`
		Options []string  `json:"options" binding:"required"`
		IsMulti bool      `json:"is_multi"`
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
