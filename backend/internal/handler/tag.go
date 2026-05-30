package handler

import (
	"strconv"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

type TagHandler struct {
	svc *service.ForumService
}

func NewTagHandler() *TagHandler {
	return &TagHandler{svc: service.NewForumService()}
}

// GET /api/tags
func (h *TagHandler) List(c *gin.Context) {
	tags, err := h.svc.ListTags()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}
	if tags == nil {
		tags = []model.Tag{}
	}
	pkg.Success(c, tags)
}

// POST /api/tags
func (h *TagHandler) Create(c *gin.Context) {
	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	tag, err := h.svc.CreateTag(req.Name)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumContentRequired, 400, 40001, "名称不能为空"},
			serviceError{service.ErrForumTagNameTaken, 409, 40002, "标签名已存在"},
		)
		return
	}
	pkg.Success(c, tag)
}

// PUT /api/tags/:id
func (h *TagHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	tag, err := h.svc.UpdateTag(uint(id), req.Name)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumTagNotFound, 404, 40401, "标签不存在"},
			serviceError{service.ErrForumContentRequired, 400, 40001, "名称不能为空"},
			serviceError{service.ErrForumTagNameTaken, 409, 40002, "标签名已存在"},
		)
		return
	}
	pkg.Success(c, tag)
}

// DELETE /api/tags/:id
func (h *TagHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.DeleteTag(uint(id)); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrForumTagInUse, 400, 40004, "该标签被话题引用，不能删除"},
		)
		return
	}
	pkg.Success(c, nil)
}
