package handler

import (
	"errors"
	"strconv"
	"time"

	"warmisle/internal/pkg"
	"warmisle/internal/repository"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

type TodoHandler struct {
	svc *service.TodoService
}

func NewTodoHandler() *TodoHandler {
	return &TodoHandler{svc: service.NewTodoService()}
}

// GET /api/todos
func (h *TodoHandler) List(c *gin.Context) {
	var req struct {
		Status     string `form:"status"`
		AssigneeID *uint  `form:"assignee_id"`
		Page       int    `form:"page"`
		PageSize   int    `form:"page_size"`
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

	result, err := h.svc.List(repository.TodoFilter{
		Status:     req.Status,
		AssigneeID: req.AssigneeID,
		Page:       req.Page,
		PageSize:   req.PageSize,
	})
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, result)
}

type createTodoRequest struct {
	Title       string  `json:"title" binding:"required"`
	Description string  `json:"description"`
	Priority    string  `json:"priority"`
	AssigneeID  *uint   `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

// POST /api/todos
func (h *TodoHandler) Create(c *gin.Context) {
	var req createTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Title == "" {
		pkg.Error(c, 400, 40001, "标题不能为空")
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			t, err = time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				pkg.Error(c, 400, 40001, "日期格式错误")
				return
			}
		}
		dueDate = &t
	}

	creatorID := getMemberID(c)

	result, err := h.svc.Create(req.Title, req.Description, req.Priority, req.AssigneeID, dueDate, creatorID)
	if err != nil {
		if errors.Is(err, service.ErrTodoTitleRequired) {
			pkg.Error(c, 400, 40001, "标题不能为空")
			return
		}
		if errors.Is(err, service.ErrTodoInvalidPriority) {
			pkg.Error(c, 400, 40001, "优先级无效")
			return
		}
		if errors.Is(err, service.ErrTodoAssigneeNotFound) {
			pkg.Error(c, 400, 40001, "指派的成员不存在")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, result)
}

type updateTodoRequest struct {
	Title       *string `json:"title"`
	Description *string `json:"description"`
	Priority    *string `json:"priority"`
	AssigneeID  *uint   `json:"assignee_id"`
	DueDate     *string `json:"due_date"`
}

// PUT /api/todos/:id
func (h *TodoHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req updateTodoRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Title != nil && *req.Title == "" {
		pkg.Error(c, 400, 40001, "标题不能为空")
		return
	}

	var dueDate *time.Time
	if req.DueDate != nil && *req.DueDate != "" {
		t, err := time.Parse("2006-01-02", *req.DueDate)
		if err != nil {
			t, err = time.Parse(time.RFC3339, *req.DueDate)
			if err != nil {
				pkg.Error(c, 400, 40001, "日期格式错误")
				return
			}
		}
		dueDate = &t
	}

	currentMemberID := getMemberID(c)
	currentRole := getRole(c)

	result, err := h.svc.Update(uint(id), req.Title, req.Description, req.Priority, req.AssigneeID, dueDate, currentMemberID, currentRole)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			pkg.Error(c, 404, 40001, "待办事项不存在")
			return
		}
		if errors.Is(err, service.ErrTodoPermissionDenied) {
			pkg.Error(c, 403, 40301, "只能编辑自己创建或负责的待办")
			return
		}
		if errors.Is(err, service.ErrTodoTitleRequired) {
			pkg.Error(c, 400, 40001, "标题不能为空")
			return
		}
		if errors.Is(err, service.ErrTodoInvalidPriority) {
			pkg.Error(c, 400, 40001, "优先级无效")
			return
		}
		if errors.Is(err, service.ErrTodoAssigneeNotFound) {
			pkg.Error(c, 400, 40001, "指派的成员不存在")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, result)
}

// DELETE /api/todos/:id
func (h *TodoHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	currentMemberID := getMemberID(c)
	currentRole := getRole(c)

	if err := h.svc.Delete(uint(id), currentMemberID, currentRole); err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			pkg.Error(c, 404, 40001, "待办事项不存在")
			return
		}
		if errors.Is(err, service.ErrTodoPermissionDenied) {
			pkg.Error(c, 403, 40301, "只能删除自己创建的待办")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, nil)
}

// PUT /api/todos/:id/toggle
func (h *TodoHandler) Toggle(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	currentMemberID := getMemberID(c)
	currentRole := getRole(c)

	result, err := h.svc.Toggle(uint(id), currentMemberID, currentRole)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			pkg.Error(c, 404, 40001, "待办事项不存在")
			return
		}
		if errors.Is(err, service.ErrTodoPermissionDenied) {
			pkg.Error(c, 403, 40301, "只能操作自己创建或负责的待办")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, result)
}

// PUT /api/todos/:id/claim
func (h *TodoHandler) Claim(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	currentMemberID := getMemberID(c)

	result, err := h.svc.Claim(uint(id), currentMemberID)
	if err != nil {
		if errors.Is(err, service.ErrTodoNotFound) {
			pkg.Error(c, 404, 40001, "待办事项不存在")
			return
		}
		if errors.Is(err, service.ErrTodoAlreadyAssigned) {
			pkg.Error(c, 400, 40001, "该待办已被认领")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, result)
}
