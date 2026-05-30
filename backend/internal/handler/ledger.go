package handler

import (
	"strconv"
	"time"

	"warmisle/internal/pkg"
	"warmisle/internal/repository"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

type LedgerHandler struct {
	svc *service.LedgerService
}

func NewLedgerHandler() *LedgerHandler {
	return &LedgerHandler{svc: service.NewLedgerService()}
}

func getCurrentMonth() string {
	return time.Now().Format("2006-01")
}

// GET /api/ledgers
func (h *LedgerHandler) List(c *gin.Context) {
	var req struct {
		Month      string `form:"month"`
		MemberID   *uint  `form:"member_id"`
		CategoryID *uint  `form:"category_id"`
		CreatorID  *uint  `form:"creator_id"`
		Page       int    `form:"page"`
		PageSize   int    `form:"page_size"`
	}

	if err := c.ShouldBindQuery(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Month == "" {
		req.Month = getCurrentMonth()
	}
	if req.Page <= 0 {
		req.Page = 1
	}
	if req.PageSize <= 0 {
		req.PageSize = 20
	}

	filter := repository.LedgerFilter{
		Month:      req.Month,
		MemberID:   req.MemberID,
		CategoryID: req.CategoryID,
		CreatorID:  req.CreatorID,
		Page:       req.Page,
		PageSize:   req.PageSize,
	}

	result, err := h.svc.List(filter)
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	if result.Groups == nil {
		result.Groups = []repository.LedgerGroup{}
	}

	pkg.Success(c, result)
}

// GET /api/ledgers/:id
func (h *LedgerHandler) GetByID(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	result, err := h.svc.FindByID(uint(id))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrLedgerNotFound, 404, 40401, "记账记录不存在"},
		)
		return
	}

	pkg.Success(c, result)
}

type createLedgerRequest struct {
	Amount     float64 `json:"amount" binding:"required"`
	Note       string  `json:"note"`
	CategoryID uint    `json:"category_id" binding:"required"`
	MemberIDs  []uint  `json:"member_ids" binding:"required"`
	OccurredAt string  `json:"occurred_at"`
}

// POST /api/ledgers
func (h *LedgerHandler) Create(c *gin.Context) {
	var req createLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Amount <= 0 {
		pkg.Error(c, 400, 40001, "金额必须大于 0")
		return
	}
	if len(req.MemberIDs) == 0 {
		pkg.Error(c, 400, 40001, "请至少选择一位关联成员")
		return
	}

	amountCents := int64(req.Amount)

	occurredAt := time.Now()
	if req.OccurredAt != "" {
		var err error
		occurredAt, err = time.Parse(time.RFC3339, req.OccurredAt)
		if err != nil {
			// Try local datetime format
			occurredAt, err = time.Parse("2006-01-02T15:04:05Z", req.OccurredAt)
			if err != nil {
				occurredAt, err = time.Parse("2006-01-02 15:04:05", req.OccurredAt)
				if err != nil {
					pkg.Error(c, 400, 40001, "日期格式错误")
					return
				}
			}
		}
	}

	creatorID := getMemberID(c)

	result, err := h.svc.Create(amountCents, req.Note, req.CategoryID, req.MemberIDs, occurredAt, creatorID)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrInvalidAmount, 400, 40001, "金额必须大于 0"},
			serviceError{service.ErrLedgerCategoryNotFound, 400, 40001, "分类不存在"},
			serviceError{service.ErrNoMembers, 400, 40001, "请至少选择一位关联成员"},
		)
		return
	}

	pkg.Success(c, result)
}

type updateLedgerRequest struct {
	Amount     *float64 `json:"amount"`
	Note       *string  `json:"note"`
	CategoryID *uint    `json:"category_id"`
	MemberIDs  []uint   `json:"member_ids"`
	OccurredAt *string  `json:"occurred_at"`
}

// PUT /api/ledgers/:id
func (h *LedgerHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req updateLedgerRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	var amountCents *int64
	if req.Amount != nil {
		if *req.Amount <= 0 {
			pkg.Error(c, 400, 40001, "金额必须大于 0")
			return
		}
		cents := int64(*req.Amount)
		amountCents = &cents
	}

	var occurredAt *time.Time
	if req.OccurredAt != nil && *req.OccurredAt != "" {
		t, err := time.Parse(time.RFC3339, *req.OccurredAt)
		if err != nil {
			t, err = time.Parse("2006-01-02T15:04:05Z", *req.OccurredAt)
			if err != nil {
				t, err = time.Parse("2006-01-02 15:04:05", *req.OccurredAt)
				if err != nil {
					pkg.Error(c, 400, 40001, "日期格式错误")
					return
				}
			}
		}
		occurredAt = &t
	}

	currentMemberID := getMemberID(c)
	currentRole := getRole(c)

	result, err := h.svc.Update(uint(id), amountCents, req.Note, req.CategoryID, req.MemberIDs, occurredAt, currentMemberID, currentRole)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrLedgerNotFound, 404, 40401, "记账记录不存在"},
			serviceError{service.ErrLedgerPermissionDenied, 403, 40301, "只能修改自己创建的记录"},
			serviceError{service.ErrInvalidAmount, 400, 40001, "金额必须大于 0"},
			serviceError{service.ErrLedgerCategoryNotFound, 400, 40001, "分类不存在"},
			serviceError{service.ErrNoMembers, 400, 40001, "请至少选择一位关联成员"},
		)
		return
	}

	pkg.Success(c, result)
}

// DELETE /api/ledgers/:id
func (h *LedgerHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	currentMemberID := getMemberID(c)
	currentRole := getRole(c)

	if err := h.svc.Delete(uint(id), currentMemberID, currentRole); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrLedgerNotFound, 404, 40401, "记账记录不存在"},
			serviceError{service.ErrLedgerPermissionDenied, 403, 40301, "只能删除自己创建的记录"},
		)
		return
	}

	pkg.Success(c, nil)
}
