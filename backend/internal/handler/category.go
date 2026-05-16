package handler

import (
	"errors"
	"strconv"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/service"

	"github.com/gin-gonic/gin"
)

type CategoryHandler struct {
	svc *service.CategoryService
}

func NewCategoryHandler() *CategoryHandler {
	return &CategoryHandler{svc: service.NewCategoryService()}
}

type createCategoryRequest struct {
	Type      string `json:"type" binding:"required"`
	Name      string `json:"name" binding:"required"`
	Icon      string `json:"icon"`
	SortOrder int    `json:"sort_order"`
}

func (h *CategoryHandler) Create(c *gin.Context) {
	var req createCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Type != "income" && req.Type != "expense" {
		pkg.Error(c, 400, 40001, "类型必须为 income 或 expense")
		return
	}

	cat, err := h.svc.Create(req.Type, req.Name, req.Icon, req.SortOrder)
	if err != nil {
		if errors.Is(err, service.ErrCategoryExists) {
			pkg.Error(c, 409, 40002, "同类型下已存在同名分类")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, cat)
}

type updateCategoryRequest struct {
	Type      string `json:"type"`
	Name      string `json:"name"`
	Icon      string `json:"icon"`
	SortOrder *int   `json:"sort_order"`
}

func (h *CategoryHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req updateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if req.Type != "" && req.Type != "income" && req.Type != "expense" {
		pkg.Error(c, 400, 40001, "类型必须为 income 或 expense")
		return
	}

	cat, err := h.svc.Update(uint(id), req.Type, req.Name, req.Icon, req.SortOrder)
	if err != nil {
		if errors.Is(err, service.ErrCategoryNotFound) {
			pkg.Error(c, 404, 40001, "分类不存在")
			return
		}
		if errors.Is(err, service.ErrCategoryExists) {
			pkg.Error(c, 409, 40002, "同类型下已存在同名分类")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, cat)
}

func (h *CategoryHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.Delete(uint(id)); err != nil {
		if errors.Is(err, service.ErrCategoryInUse) {
			pkg.Error(c, 400, 40001, "该分类下有记账记录，无法删除")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, nil)
}

func (h *CategoryHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	if list == nil {
		list = []model.Category{}
	}

	pkg.Success(c, list)
}
