package repository

import (
	"fmt"

	"warmisle/internal/model"

	"gorm.io/gorm"
)

// BackupRepo handles backup-related database operations.
type BackupRepo struct {
	db *gorm.DB
}

// NewBackupRepo creates a new BackupRepo with the given gorm.DB instance.
func NewBackupRepo(db *gorm.DB) *BackupRepo {
	return &BackupRepo{db: db}
}

// GetConfig retrieves the single cloud drive configuration row for the alipan provider.
func (r *BackupRepo) GetConfig() (*model.CloudDriveConfig, error) {
	var cfg model.CloudDriveConfig
	if err := r.db.Where("provider = ?", "alipan").First(&cfg).Error; err != nil {
		return nil, err
	}
	return &cfg, nil
}

// SaveConfig creates or updates a cloud drive configuration.
func (r *BackupRepo) SaveConfig(cfg *model.CloudDriveConfig) error {
	return r.db.Save(cfg).Error
}

// CreateRecord inserts a new backup record.
func (r *BackupRepo) CreateRecord(record *model.BackupRecord) error {
	return r.db.Create(record).Error
}

// ListRecords returns paginated backup records ordered by created_at descending.
func (r *BackupRepo) ListRecords(page, pageSize int) ([]model.BackupRecord, int64, error) {
	var records []model.BackupRecord
	var total int64

	// Count total records
	if err := r.db.Model(&model.BackupRecord{}).Count(&total).Error; err != nil {
		return nil, 0, err
	}

	// Calculate offset
	offset := (page - 1) * pageSize
	if offset < 0 {
		offset = 0
	}

	// Fetch paginated records
	if err := r.db.Order("created_at DESC").Offset(offset).Limit(pageSize).Find(&records).Error; err != nil {
		return nil, 0, err
	}

	return records, total, nil
}

// DeleteRecord deletes a backup record by ID.
func (r *BackupRepo) DeleteRecord(id uint) error {
	return r.db.Delete(&model.BackupRecord{}, id).Error
}

// GetRecordByID finds a backup record by its ID.
func (r *BackupRepo) GetRecordByID(id uint) (*model.BackupRecord, error) {
	var record model.BackupRecord
	if err := r.db.First(&record, id).Error; err != nil {
		return nil, err
	}
	return &record, nil
}

// ListExpiredRecords returns non-pre_restore records older than the given number of days.
func (r *BackupRepo) ListExpiredRecords(beforeDays int) ([]model.BackupRecord, error) {
	var records []model.BackupRecord
	if err := r.db.
		Where("is_pre_restore = ?", 0).
		Where("created_at < datetime('now', ? || ' days')", fmt.Sprintf("-%d", beforeDays)).
		Find(&records).Error; err != nil {
		return nil, err
	}
	return records, nil
}
