package handler

import (
	"strconv"

	"warmisle/internal/pkg"
	"warmisle/internal/repository"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

type WishHandler struct {
	svc *service.WishService
}

func NewWishHandler() *WishHandler {
	return &WishHandler{svc: service.NewWishService()}
}

// GET /api/wishes
func (h *WishHandler) List(c *gin.Context) {
	var req struct {
		Type      string `form:"type"`
		Status    string `form:"status"`
		CreatorID *uint  `form:"creator_id"`
		Page      int    `form:"page"`
		PageSize  int    `form:"page_size"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	result, err := h.svc.List(repository.WishFilter{
		Type:      req.Type,
		Status:    req.Status,
		CreatorID: req.CreatorID,
		Page:      req.Page,
		PageSize:  req.PageSize,
	})
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, result)
}

type createWishRequest struct {
	Title       string `json:"title" binding:"required"`
	Description string `json:"description"`
	Category    string `json:"category"`
	Priority    string `json:"priority"`
	Type        string `json:"type"`
	Amount      *int64 `json:"amount"`
}

// POST /api/wishes
func (h *WishHandler) Create(c *gin.Context) {
	var req createWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.Create(req.Title, req.Description, req.Category, req.Priority, req.Type, req.Amount, getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrWishTitleRequired, 400, 40001, "标题不能为空"},
			serviceError{service.ErrWishInvalidPriority, 400, 40001, "优先级无效"},
			serviceError{service.ErrWishInvalidCategory, 400, 40001, "分类无效"},
		)
		return
	}

	pkg.Success(c, result)
}

type updateWishRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Category    *string `json:"category"`
	Priority    *string `json:"priority"`
	Amount      *int64  `json:"amount"`
}

// PUT /api/wishes/:id
func (h *WishHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req updateWishRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Title != nil && *req.Title == "" {
		pkg.Error(c, 400, 40001, "标题不能为空")
		return
	}

	result, err := h.svc.Update(uint(id), req.Title, req.Description, req.Category, req.Priority, req.Amount, getMemberID(c), getRole(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrWishNotFound, 404, 40401, "愿望不存在"},
			serviceError{service.ErrWishPermissionDenied, 403, 40301, "只能编辑自己创建的愿望"},
			serviceError{service.ErrWishTitleRequired, 400, 40001, "标题不能为空"},
			serviceError{service.ErrWishInvalidPriority, 400, 40001, "优先级无效"},
			serviceError{service.ErrWishInvalidCategory, 400, 40001, "分类无效"},
		)
		return
	}

	pkg.Success(c, result)
}

// DELETE /api/wishes/:id
func (h *WishHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.Delete(uint(id), getMemberID(c), getRole(c)); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrWishNotFound, 404, 40401, "愿望不存在"},
			serviceError{service.ErrWishPermissionDenied, 403, 40301, "只能删除自己创建的愿望"},
		)
		return
	}

	pkg.Success(c, nil)
}

// POST /api/wishes/:id/promote
func (h *WishHandler) Promote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.Promote(uint(id), getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrWishNotFound, 404, 40401, "愿望不存在"},
			serviceError{service.ErrWishPermissionDenied, 403, 40301, "只能提升自己创建的愿望"},
		)
		return
	}

	pkg.Success(c, result)
}

type updateWishStatusRequest struct {
	Status string `json:"status" binding:"required"`
}

// PUT /api/wishes/:id/status
func (h *WishHandler) UpdateStatus(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req updateWishStatusRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	result, err := h.svc.UpdateStatus(uint(id), req.Status, getMemberID(c), getRole(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrWishNotFound, 404, 40401, "愿望不存在"},
			serviceError{service.ErrWishPermissionDenied, 403, 40301, "权限不足"},
			serviceError{service.ErrWishInvalidStatus, 400, 40001, "状态值无效"},
		)
		return
	}

	pkg.Success(c, result)
}

// POST /api/wishes/:id/vote
func (h *WishHandler) Vote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.Vote(uint(id), getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrWishNotFound, 404, 40401, "愿望不存在"},
			serviceError{service.ErrWishAlreadyVoted, 400, 40001, "已投票"},
		)
		return
	}

	pkg.Success(c, result)
}

// DELETE /api/wishes/:id/vote
func (h *WishHandler) Unvote(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.Unvote(uint(id), getMemberID(c))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrWishNotFound, 404, 40401, "愿望不存在"},
			serviceError{service.ErrWishNotVoted, 400, 40001, "未投票"},
		)
		return
	}

	pkg.Success(c, result)
}
