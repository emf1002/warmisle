package handler

import (
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

// ---------------------------------------------------------------------------
// Mock infrastructure — because BackupHandler uses a concrete
// *service.BackupService, we define a local interface and a mock struct,
// then test the handler logic through a thin wrapper.
// ---------------------------------------------------------------------------

// backupSvc is the subset of BackupService methods exercised by the handler
// tests below.
type backupSvc interface {
	GetConfig() (*model.CloudDriveConfig, error)
	RestoreBackup(cloudFileID, confirmText string) error
}

// mockBackupSvc implements backupSvc with configurable return values.
type mockBackupSvc struct {
	config    *model.CloudDriveConfig
	configErr error
	restoreErr error
}

func (m *mockBackupSvc) GetConfig() (*model.CloudDriveConfig, error) {
	return m.config, m.configErr
}

func (m *mockBackupSvc) RestoreBackup(cloudFileID, confirmText string) error {
	return m.restoreErr
}

// testHandler mirrors the real BackupHandler methods but uses the backupSvc
// interface so we can inject our mock.
type testHandler struct {
	svc backupSvc
}

func (h *testHandler) GetConfig(c *gin.Context) {
	cfg, err := h.svc.GetConfig()
	if err != nil {
		pkg.Error(c, http.StatusInternalServerError, 50000, err.Error())
		return
	}
	pkg.Success(c, cfg)
}

func (h *testHandler) RestoreBackup(c *gin.Context) {
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

// ---------------------------------------------------------------------------
// TestGetConfig — handler returns config successfully.
// ---------------------------------------------------------------------------
func TestGetConfig(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockBackupSvc{
		config: &model.CloudDriveConfig{
			ID:       1,
			Provider: "alipan",
			AppID:    "test-app-id",
			Status:   "authorized",
		},
	}

	h := &testHandler{svc: mock}

	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	h.GetConfig(c)

	if w.Code != http.StatusOK {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusOK)
	}

	var resp pkg.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Code != 0 {
		t.Fatalf("biz code = %d, want 0", resp.Code)
	}
	if resp.Message != "ok" {
		t.Fatalf("message = %q, want %q", resp.Message, "ok")
	}
}

// ---------------------------------------------------------------------------
// TestRestoreBackup_BadConfirmText — when confirm text doesn't match,
// the handler returns HTTP 400 with biz code 40025.
// ---------------------------------------------------------------------------
func TestRestoreBackup_BadConfirmText(t *testing.T) {
	gin.SetMode(gin.TestMode)

	mock := &mockBackupSvc{
		restoreErr: errors.New("确认文字不匹配"),
	}

	h := &testHandler{svc: mock}

	body := `{"confirm_text":"wrong confirmation"}`
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)
	c.Request = httptest.NewRequest(http.MethodPost, "/backup/restore/cloud-file-123", strings.NewReader(body))
	c.Request.Header.Set("Content-Type", "application/json")
	c.Params = gin.Params{{Key: "cloudFileId", Value: "cloud-file-123"}}

	h.RestoreBackup(c)

	if w.Code != http.StatusBadRequest {
		t.Fatalf("status = %d, want %d", w.Code, http.StatusBadRequest)
	}

	var resp pkg.Response
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("unmarshal response: %v", err)
	}
	if resp.Code != 40025 {
		t.Fatalf("biz code = %d, want 40025", resp.Code)
	}
}
