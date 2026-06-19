package repository

import (
	"testing"
	"time"

	"warmisle/internal/model"

	"github.com/glebarez/sqlite"
	"gorm.io/gorm"
)

// setupTestDB creates an in-memory SQLite database with auto-migrated tables
// and returns a gorm.DB instance and a BackupRepo.
func setupTestDB(t *testing.T) (*gorm.DB, *BackupRepo) {
	t.Helper()

	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{})
	if err != nil {
		t.Fatalf("failed to open in-memory SQLite: %v", err)
	}

	// Auto-migrate the required models.
	if err := db.AutoMigrate(&model.CloudDriveConfig{}, &model.BackupRecord{}); err != nil {
		t.Fatalf("AutoMigrate() error = %v", err)
	}

	repo := NewBackupRepo(db)
	return db, repo
}

// seedConfig creates and returns a default CloudDriveConfig row.
func seedConfig(t *testing.T, db *gorm.DB) *model.CloudDriveConfig {
	t.Helper()

	cfg := &model.CloudDriveConfig{
		Provider:   "alipan",
		AppID:      "test-app-id",
		Status:     "unconfigured",
		BackupDir:  "/test-backups/",
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}
	if err := db.Create(cfg).Error; err != nil {
		t.Fatalf("seedConfig() error = %v", err)
	}
	return cfg
}

// ---------------------------------------------------------------------------
// TestGetConfig verifies that GetConfig returns the initial configuration.
// ---------------------------------------------------------------------------
func TestGetConfig(t *testing.T) {
	db, repo := setupTestDB(t)
	_ = seedConfig(t, db)

	cfg, err := repo.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() error = %v", err)
	}
	if cfg == nil {
		t.Fatal("GetConfig() returned nil config")
	}
	if cfg.Provider != "alipan" {
		t.Fatalf("Provider = %q, want %q", cfg.Provider, "alipan")
	}
	if cfg.AppID != "test-app-id" {
		t.Fatalf("AppID = %q, want %q", cfg.AppID, "test-app-id")
	}
}

// ---------------------------------------------------------------------------
// TestSaveConfig verifies that updating a config via SaveConfig is persisted
// and re-reading returns the updated values.
// ---------------------------------------------------------------------------
func TestSaveConfig(t *testing.T) {
	db, repo := setupTestDB(t)
	cfg := seedConfig(t, db)

	// Update fields.
	cfg.AppID = "updated-app-id"
	cfg.Status = "authorized"
	cfg.EncryptedSecret = "enc-secret-123"
	if err := repo.SaveConfig(cfg); err != nil {
		t.Fatalf("SaveConfig() error = %v", err)
	}

	// Re-read and verify.
	got, err := repo.GetConfig()
	if err != nil {
		t.Fatalf("GetConfig() after save error = %v", err)
	}
	if got.AppID != "updated-app-id" {
		t.Fatalf("AppID = %q, want %q", got.AppID, "updated-app-id")
	}
	if got.Status != "authorized" {
		t.Fatalf("Status = %q, want %q", got.Status, "authorized")
	}
	if got.EncryptedSecret != "enc-secret-123" {
		t.Fatalf("EncryptedSecret = %q, want %q", got.EncryptedSecret, "enc-secret-123")
	}
}

// ---------------------------------------------------------------------------
// TestCreateRecord verifies that creating a backup record persists it and it
// can be retrieved afterward.
// ---------------------------------------------------------------------------
func TestCreateRecord(t *testing.T) {
	_, repo := setupTestDB(t)

	now := time.Now()
	record := &model.BackupRecord{
		FileName:     "backup_20250101.db",
		CloudFileID:  "cloud-file-123",
		FileSize:     4096,
		BackupType:   "manual",
		UploadStatus: "completed",
		IntegrityOk:  1,
		IsPreRestore: 0,
		CreatedAt:    now,
	}

	if err := repo.CreateRecord(record); err != nil {
		t.Fatalf("CreateRecord() error = %v", err)
	}

	// Verify it was assigned an ID.
	if record.ID == 0 {
		t.Fatal("CreateRecord() did not assign ID")
	}

	// Retrieve by ID.
	got, err := repo.GetRecordByID(record.ID)
	if err != nil {
		t.Fatalf("GetRecordByID() error = %v", err)
	}
	if got.FileName != "backup_20250101.db" {
		t.Fatalf("FileName = %q, want %q", got.FileName, "backup_20250101.db")
	}
	if got.CloudFileID != "cloud-file-123" {
		t.Fatalf("CloudFileID = %q, want %q", got.CloudFileID, "cloud-file-123")
	}
	if got.BackupType != "manual" {
		t.Fatalf("BackupType = %q, want %q", got.BackupType, "manual")
	}
}

// ---------------------------------------------------------------------------
// TestListRecords verifies paginated listing of backup records.
// ---------------------------------------------------------------------------
func TestListRecords(t *testing.T) {
	_, repo := setupTestDB(t)

	// Create 5 records with staggered timestamps.
	base := time.Now()
	for i := 0; i < 5; i++ {
		r := &model.BackupRecord{
			FileName:     "file_" + string(rune('A'+i)) + ".db",
			BackupType:   "manual",
			UploadStatus: "completed",
			CreatedAt:    base.Add(time.Duration(i) * time.Second),
		}
		if err := repo.CreateRecord(r); err != nil {
			t.Fatalf("CreateRecord(%d) error = %v", i, err)
		}
	}

	// Page 1, size 2.
	records, total, err := repo.ListRecords(1, 2)
	if err != nil {
		t.Fatalf("ListRecords(1,2) error = %v", err)
	}
	if total != 5 {
		t.Fatalf("total = %d, want 5", total)
	}
	if len(records) != 2 {
		t.Fatalf("len(records) = %d, want 2", len(records))
	}

	// Page 3, size 2 (should have 1 remaining).
	records2, total2, err := repo.ListRecords(3, 2)
	if err != nil {
		t.Fatalf("ListRecords(3,2) error = %v", err)
	}
	if total2 != 5 {
		t.Fatalf("total = %d, want 5", total2)
	}
	if len(records2) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records2))
	}

	// Verify descending order: page 1 should have the newest records.
	if len(records) >= 2 {
		if records[0].CreatedAt.Before(records[1].CreatedAt) {
			t.Error("records not in descending order by created_at")
		}
	}
}

// ---------------------------------------------------------------------------
// TestListExpiredRecords verifies that only non-pre_restore records older
// than the given number of days are returned.
// ---------------------------------------------------------------------------
func TestListExpiredRecords(t *testing.T) {
	_, repo := setupTestDB(t)

	now := time.Now()

	// Record 1: old, non-pre_restore → should be returned.
	r1 := &model.BackupRecord{
		FileName:     "old_normal.db",
		BackupType:   "manual",
		UploadStatus: "completed",
		IsPreRestore: 0,
		CreatedAt:    now.Add(-40 * 24 * time.Hour),
	}
	if err := repo.CreateRecord(r1); err != nil {
		t.Fatalf("CreateRecord(r1) error = %v", err)
	}

	// Record 2: old, pre_restore → should NOT be returned.
	r2 := &model.BackupRecord{
		FileName:     "old_pre_restore.db",
		BackupType:   "pre_restore",
		UploadStatus: "completed",
		IsPreRestore: 1,
		CreatedAt:    now.Add(-40 * 24 * time.Hour),
	}
	if err := repo.CreateRecord(r2); err != nil {
		t.Fatalf("CreateRecord(r2) error = %v", err)
	}

	// Record 3: recent, non-pre_restore → should NOT be returned.
	r3 := &model.BackupRecord{
		FileName:     "recent_normal.db",
		BackupType:   "manual",
		UploadStatus: "completed",
		IsPreRestore: 0,
		CreatedAt:    now.Add(-5 * 24 * time.Hour),
	}
	if err := repo.CreateRecord(r3); err != nil {
		t.Fatalf("CreateRecord(r3) error = %v", err)
	}

	// Query with beforeDays=30.
	records, err := repo.ListExpiredRecords(30)
	if err != nil {
		t.Fatalf("ListExpiredRecords(30) error = %v", err)
	}

	// Only r1 should be returned.
	if len(records) != 1 {
		t.Fatalf("len(records) = %d, want 1", len(records))
	}
	if records[0].FileName != "old_normal.db" {
		t.Fatalf("FileName = %q, want %q", records[0].FileName, "old_normal.db")
	}
}
