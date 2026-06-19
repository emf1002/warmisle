package handler

import (
	"net/http"
	"strconv"

	"warmisle/internal/pkg"
	"warmisle/internal/service"

	"github.com/gin-gonic/gin"
)

// BackupHandler handles HTTP requests for cloud drive backup operations.
type BackupHandler struct {
	svc *service.BackupService
}

// NewBackupHandler creates a new BackupHandler.
func NewBackupHandler(svc *service.BackupService) *BackupHandler {
	return &BackupHandler{svc: svc}
}

// GetConfig returns the cloud drive configuration (secrets masked).
func (h *BackupHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, cfg)
}

// SaveConfig updates the cloud drive configuration.
func (h *BackupHandler) SaveConfig(c *gin.Context) {
	var req struct {
		AppID       string `json:"app_id"`
		AppSecret   string `json:"app_secret"`
		RedirectURI string `json:"redirect_uri"`
		BackupDir   string `json:"backup_dir"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 40000, "参数错误: "+err.Error())
		return
	}
	cfg, err := h.svc.SaveConfig(req.AppID, req.AppSecret, req.RedirectURI, req.BackupDir)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, cfg)
}

// GetAuthURL generates the Aliyun Drive OAuth2 authorization URL.
func (h *BackupHandler) GetAuthURL(c *gin.Context) {
	authURL, state, err := h.svc.GetAuthURL()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, gin.H{"auth_url": authURL, "state": state})
}

// Callback handles the OAuth2 redirect from Aliyun Drive.
// Aliyun Drive redirects via GET with query parameters: ?code=xxx&state=yyy
func (h *BackupHandler) Callback(c *gin.Context) {
	code := c.Query("code")
	state := c.Query("state")
	if code == "" || state == "" {
		pkg.Error(c, http.StatusBadRequest, 40000, "缺少 code 或 state 参数")
		return
	}
	if err := h.svc.HandleCallback(code, state); err != nil {
		pkg.Error(c, http.StatusBadRequest, 40020, err.Error())
		return
	}
	// Redirect to the backup page after successful authorization
	c.Redirect(http.StatusFound, "/backup")
}

// TriggerBackup initiates a manual database backup to the cloud drive.
func (h *BackupHandler) TriggerBackup(c *gin.Context) {
	record, err := h.svc.TriggerBackup()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50010, err.Error())
		return
	}
	pkg.Success(c, record)
}

// ListCloudFiles lists backup files stored on the cloud drive.
func (h *BackupHandler) ListCloudFiles(c *gin.Context) {
	files, err := h.svc.ListCloudFiles()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, files)
}

// RestoreBackup downloads a backup from the cloud and restores the database.
func (h *BackupHandler) RestoreBackup(c *gin.Context) {
	cloudFileID := c.Param("cloudFileId")
	var req struct {
		ConfirmText string `json:"confirm_text"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 40000, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.RestoreBackup(cloudFileID, req.ConfirmText); err != nil {
		pkg.Error(c, http.StatusBadRequest, 40025, err.Error())
		return
	}
	pkg.Success(c, nil)
}

// ListHistory returns the local backup history with pagination.
func (h *BackupHandler) ListHistory(c *gin.Context) {
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	pageSize, _ := strconv.Atoi(c.DefaultQuery("page_size", "20"))
	records, total, err := h.svc.ListHistory(page, pageSize)
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, gin.H{"list": records, "total": total, "page": page, "page_size": pageSize})
}

// DeleteHistory deletes a backup record and its cloud file.
func (h *BackupHandler) DeleteHistory(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		pkg.Error(c, http.StatusBadRequest, 40000, "无效的ID")
		return
	}
	if err := h.svc.DeleteHistory(uint(id)); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, nil)
}

// GetSchedule returns the scheduled backup configuration.
func (h *BackupHandler) GetSchedule(c *gin.Context) {
	cfg, err := h.svc.GetConfig()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, gin.H{
		"schedule_enabled": cfg.ScheduleEnabled,
		"schedule_time":    cfg.ScheduleTime,
		"retention_days":   cfg.RetentionDays,
	})
}

// SaveSchedule updates the scheduled backup configuration.
func (h *BackupHandler) SaveSchedule(c *gin.Context) {
	var req struct {
		ScheduleEnabled bool   `json:"schedule_enabled"`
		ScheduleTime    string `json:"schedule_time"`
		RetentionDays   int    `json:"retention_days"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		pkg.Error(c, http.StatusBadRequest, 40000, "参数错误: "+err.Error())
		return
	}
	if err := h.svc.UpdateSchedule(req.ScheduleEnabled, req.ScheduleTime, req.RetentionDays); err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, gin.H{
		"schedule_enabled": req.ScheduleEnabled,
		"schedule_time":    req.ScheduleTime,
		"retention_days":   req.RetentionDays,
	})
}
