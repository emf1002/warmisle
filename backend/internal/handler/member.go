package handler

import (
	"strconv"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

type MemberHandler struct {
	svc *service.MemberService
}

func NewMemberHandler() *MemberHandler {
	return &MemberHandler{svc: service.NewMemberService()}
}

type createMemberRequest struct {
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Name     string `json:"name"`
	Avatar   string `json:"avatar"`
	Role     string `json:"role"`
}

func (h *MemberHandler) Create(c *gin.Context) {
	var req createMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	member, err := h.svc.Create(req.Username, req.Password, req.Name, req.Avatar, req.Role)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrInvalidUsername, 400, 40001, "用户名格式不正确（3-20位字母、数字、下划线）"},
			serviceError{service.ErrInvalidPassword, 400, 40001, "密码长度需在6-32位之间"},
			serviceError{service.ErrUsernameTaken, 409, 40002, "用户名已存在"},
		)
		return
	}

	pkg.Success(c, member)
}

type updateMemberRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
	Role   string `json:"role"`
}

func (h *MemberHandler) Update(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	var req updateMemberRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	member, err := h.svc.Update(uint(id), req.Name, req.Avatar, req.Role)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrMemberNotFound, 404, 40001, "成员不存在"},
			serviceError{service.ErrCannotDeleteLastAdmin, 400, 40003, "不能修改最后一个管理员的角色"},
			serviceError{service.ErrInvalidName, 400, 40001, "名称长度需在1-20位之间"},
		)
		return
	}

	pkg.Success(c, member)
}

func (h *MemberHandler) Delete(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	memberIDVal, _ := c.Get("member_id")
	currentMemberID := memberIDVal.(uint)

	if err := h.svc.Delete(uint(id), currentMemberID); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrMemberNotFound, 404, 40001, "成员不存在"},
			serviceError{service.ErrCannotDeleteLastAdmin, 400, 40003, "不能删除最后一个管理员"},
			serviceError{service.ErrHasActivityRecords, 400, 40004, "该成员有活动记录，建议禁用而非删除"},
		)
		return
	}

	pkg.Success(c, nil)
}

func (h *MemberHandler) Disable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	memberIDVal, _ := c.Get("member_id")
	currentMemberID := memberIDVal.(uint)

	if err := h.svc.Disable(uint(id), currentMemberID); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrMemberNotFound, 404, 40001, "成员不存在"},
			serviceError{service.ErrCannotDisableSelf, 400, 40005, "不能禁用自己"},
			serviceError{service.ErrCannotDeleteLastAdmin, 400, 40003, "不能禁用最后一个管理员"},
		)
		return
	}

	pkg.Success(c, nil)
}

func (h *MemberHandler) Enable(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	member, err := h.svc.Enable(uint(id))
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrMemberNotFound, 404, 40001, "成员不存在"},
		)
		return
	}

	pkg.Success(c, member)
}

func (h *MemberHandler) ResetPassword(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 64)
	if err != nil {
		pkg.Error(c, 400, 40001, "参数错误")
		return
	}

	if err := h.svc.ResetPassword(uint(id)); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrMemberNotFound, 404, 40001, "成员不存在"},
		)
		return
	}

	pkg.Success(c, nil)
}

func (h *MemberHandler) List(c *gin.Context) {
	list, err := h.svc.List()
	if err != nil {
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	if list == nil {
		list = []model.Member{}
	}

	pkg.Success(c, list)
}

type updateProfileRequest struct {
	Name   string `json:"name"`
	Avatar string `json:"avatar"`
}

func (h *MemberHandler) UpdateProfile(c *gin.Context) {
	memberIDVal, _ := c.Get("member_id")
	memberID := memberIDVal.(uint)

	var req updateProfileRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	member, err := h.svc.UpdateProfile(memberID, req.Name, req.Avatar)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrInvalidName, 400, 40001, "名称长度需在1-20位之间"},
		)
		return
	}

	pkg.Success(c, member)
}

func (h *MemberHandler) GetProfile(c *gin.Context) {
	memberIDVal, _ := c.Get("member_id")
	memberID := memberIDVal.(uint)

	member, err := h.svc.GetProfile(memberID)
	if err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrMemberNotFound, 404, 40001, "成员不存在"},
		)
		return
	}

	pkg.Success(c, member)
}

type changePasswordRequest struct {
	OldPassword string `json:"old_password" binding:"required"`
	NewPassword string `json:"new_password" binding:"required"`
}

func (h *MemberHandler) ChangePassword(c *gin.Context) {
	memberIDVal, _ := c.Get("member_id")
	memberID := memberIDVal.(uint)

	var req changePasswordRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, 400, 40001, "参数校验失败")
		return
	}

	if err := h.svc.ChangePassword(memberID, req.OldPassword, req.NewPassword); err != nil {
		handleServiceError(c, err,
			serviceError{service.ErrMemberNotFound, 404, 40001, "成员不存在"},
			serviceError{service.ErrIncorrectPassword, 400, 40010, "原密码错误"},
			serviceError{service.ErrInvalidPassword, 400, 40001, "新密码长度需在6-32位之间"},
		)
		return
	}

	pkg.Success(c, nil)
}
