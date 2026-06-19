package model

import "time"

// CloudDriveConfig represents a single-row cloud drive configuration.
// There should only be one row for the 'alipan' provider.
type CloudDriveConfig struct {
	ID              uint       `gorm:"primaryKey" json:"id"`
	Provider        string     `gorm:"not null;default:alipan" json:"provider"`
	AppID           string     `gorm:"not null;default:''" json:"app_id"`
	EncryptedSecret string     `gorm:"not null;default:''" json:"encrypted_secret"`
	RedirectURI     string     `gorm:"not null;default:''" json:"redirect_uri"`
	EncryptedToken  string     `gorm:"not null;default:''" json:"encrypted_token"`
	TokenExpiry     *time.Time `json:"token_expiry"`
	Status          string     `gorm:"not null;default:unconfigured" json:"status"`
	BackupDir       string     `gorm:"not null;default:/warmisle-backups/" json:"backup_dir"`
	ScheduleEnabled int        `gorm:"not null;default:0" json:"schedule_enabled"`
	ScheduleTime    string     `gorm:"not null;default:03:00" json:"schedule_time"`
	RetentionDays   int        `gorm:"not null;default:30" json:"retention_days"`
	CreatedAt       time.Time  `json:"created_at"`
	UpdatedAt       time.Time  `json:"updated_at"`
}

// TableName returns the table name for CloudDriveConfig.
func (CloudDriveConfig) TableName() string { return "cloud_drive_configs" }

// BackupRecord represents a single backup record entry.
type BackupRecord struct {
	ID            uint      `gorm:"primaryKey" json:"id"`
	FileName      string    `gorm:"not null" json:"file_name"`
	CloudFileID   string    `gorm:"not null;default:''" json:"cloud_file_id"`
	FileSize      int64     `gorm:"not null;default:0" json:"file_size"`
	BackupType    string    `gorm:"not null;default:manual" json:"backup_type"`
	UploadStatus  string    `gorm:"not null;default:pending" json:"upload_status"`
	IntegrityOk   int       `gorm:"not null;default:0" json:"integrity_ok"`
	ErrorMessage  string    `gorm:"not null;default:''" json:"error_message"`
	IsPreRestore  int       `gorm:"not null;default:0" json:"is_pre_restore"`
	CreatedAt     time.Time `json:"created_at"`
}

// TableName returns the table name for BackupRecord.
func (BackupRecord) TableName() string { return "backup_records" }
