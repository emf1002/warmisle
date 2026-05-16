package handler

import (
	"errors"
	"strconv"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/service"

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
		if errors.Is(err, service.ErrInvalidUsername) {
			pkg.Error(c, 400, 40001, "用户名格式不正确（3-20位字母、数字、下划线）")
			return
		}
		if errors.Is(err, service.ErrInvalidPassword) {
			pkg.Error(c, 400, 40001, "密码长度需在6-32位之间")
			return
		}
		if errors.Is(err, service.ErrUsernameTaken) {
			pkg.Error(c, 409, 40002, "用户名已存在")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
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
		if errors.Is(err, service.ErrMemberNotFound) {
			pkg.Error(c, 404, 40001, "成员不存在")
			return
		}
		if errors.Is(err, service.ErrCannotDeleteLastAdmin) {
			pkg.Error(c, 400, 40003, "不能修改最后一个管理员的角色")
			return
		}
		if errors.Is(err, service.ErrInvalidName) {
			pkg.Error(c, 400, 40001, "名称长度需在1-20位之间")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
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
		if errors.Is(err, service.ErrMemberNotFound) {
			pkg.Error(c, 404, 40001, "成员不存在")
			return
		}
		if errors.Is(err, service.ErrCannotDeleteLastAdmin) {
			pkg.Error(c, 400, 40003, "不能删除最后一个管理员")
			return
		}
		if errors.Is(err, service.ErrHasActivityRecords) {
			pkg.Error(c, 400, 40004, "该成员有活动记录，建议禁用而非删除")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
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
		if errors.Is(err, service.ErrMemberNotFound) {
			pkg.Error(c, 404, 40001, "成员不存在")
			return
		}
		if errors.Is(err, service.ErrCannotDisableSelf) {
			pkg.Error(c, 400, 40005, "不能禁用自己")
			return
		}
		if errors.Is(err, service.ErrCannotDeleteLastAdmin) {
			pkg.Error(c, 400, 40003, "不能禁用最后一个管理员")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
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
		if errors.Is(err, service.ErrMemberNotFound) {
			pkg.Error(c, 404, 40001, "成员不存在")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
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
		if errors.Is(err, service.ErrMemberNotFound) {
			pkg.Error(c, 404, 40001, "成员不存在")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
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
		if errors.Is(err, service.ErrInvalidName) {
			pkg.Error(c, 400, 40001, "名称长度需在1-20位之间")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, member)
}

func (h *MemberHandler) GetProfile(c *gin.Context) {
	memberIDVal, _ := c.Get("member_id")
	memberID := memberIDVal.(uint)

	member, err := h.svc.GetProfile(memberID)
	if err != nil {
		if errors.Is(err, service.ErrMemberNotFound) {
			pkg.Error(c, 404, 40001, "成员不存在")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
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
		if errors.Is(err, service.ErrMemberNotFound) {
			pkg.Error(c, 404, 40001, "成员不存在")
			return
		}
		if errors.Is(err, service.ErrIncorrectPassword) {
			pkg.Error(c, 400, 40010, "原密码错误")
			return
		}
		if errors.Is(err, service.ErrInvalidPassword) {
			pkg.Error(c, 400, 40001, "新密码长度需在6-32位之间")
			return
		}
		pkg.Error(c, 500, 50001, "服务器内部错误")
		return
	}

	pkg.Success(c, nil)
}
