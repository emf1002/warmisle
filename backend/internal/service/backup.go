// Package service implements the business logic for cloud backup operations.
package service

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"sync"
	"time"

	"gorm.io/gorm"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/plugin"
	"warmisle/internal/plugin/alipan"
	"warmisle/internal/repository"
)

// BackupService manages cloud backup and restore operations.
// It coordinates between the local SQLite database, the cloud drive plugin,
// and the backup record repository.
type BackupService struct {
	repo       *repository.BackupRepo
	db         *gorm.DB
	mu         sync.Mutex
	restoring  bool
	restoreMu  sync.Mutex
	dbPath     string
	backupDir  string
	stateStore map[string]*oauthState // state -> OAuth metadata
	stateMu    sync.Mutex
}

// oauthState holds metadata associated with an OAuth2 authorization state.
type oauthState struct {
	timestamp    time.Time
	codeVerifier string // PKCE code_verifier, empty for confidential client mode
}

// NewBackupService creates a new BackupService instance.
func NewBackupService(repo *repository.BackupRepo, db *gorm.DB, dbPath string) *BackupService {
	return &BackupService{
		repo:       repo,
		db:         db,
		dbPath:     dbPath,
		backupDir:  "backups",
		stateStore: make(map[string]*oauthState),
	}
}

// GetConfig retrieves the current cloud drive configuration.
// EncryptedSecret and EncryptedToken are cleared before returning
// for security — they should never be exposed to the API layer.
func (s *BackupService) GetConfig() (*model.CloudDriveConfig, error) {
	cfg, err := s.repo.GetConfig()
	if err != nil {
		return nil, err
	}
	// Never expose encrypted secrets to callers.
	cfg.EncryptedSecret = ""
	cfg.EncryptedToken = ""
	return cfg, nil
}

// SaveConfig creates or updates the cloud drive configuration.
// If appSecret is non-empty, it is encrypted and stored.
// If redirectURI or backupDir are non-empty, they are set on the config.
// AppID is always updated.
func (s *BackupService) SaveConfig(appID, appSecret, redirectURI, backupDir string) (*model.CloudDriveConfig, error) {
	cfg, err := s.repo.GetConfig()
	if err != nil {
		// No config yet — create a new one.
		cfg = &model.CloudDriveConfig{
			Provider: "alipan",
		}
	}

	cfg.AppID = appID

	if appSecret != "" {
		encrypted, encErr := pkg.Encrypt(appSecret, pkg.GetBackupKey())
		if encErr != nil {
			return nil, fmt.Errorf("加密 app_secret 失败: %w", encErr)
		}
		cfg.EncryptedSecret = encrypted
	}

	// If AppID is provided (and secret isn't set, meaning PKCE mode),
	// or if a new secret is provided, reset to pending_auth.
	if appID != "" && (appSecret != "" || cfg.Status == "unconfigured") {
		cfg.Status = "pending_auth"
	}

	if redirectURI != "" {
		cfg.RedirectURI = redirectURI
	}

	if backupDir != "" {
		cfg.BackupDir = backupDir
	}

	cfg.UpdatedAt = time.Now()

	if saveErr := s.repo.SaveConfig(cfg); saveErr != nil {
		return nil, saveErr
	}

	return cfg, nil
}

// GetAuthURL generates an OAuth authorization URL and an anti-CSRF state token.
// The state is stored in-memory and must be presented back in HandleCallback.
func (s *BackupService) GetAuthURL() (authURL, state string, err error) {
	cfg, err := s.repo.GetConfig()
	if err != nil {
		return "", "", fmt.Errorf("获取配置失败: %w", err)
	}

	var appSecret string
	if cfg.EncryptedSecret != "" {
		var decErr error
		appSecret, decErr = pkg.Decrypt(cfg.EncryptedSecret, pkg.GetBackupKey())
		if decErr != nil {
			return "", "", fmt.Errorf("解密 app_secret 失败: %w", decErr)
		}
	}

	oauth := alipan.NewAlipanOAuth(cfg.AppID, appSecret, cfg.RedirectURI)

	// Generate a random 8-byte state token.
	b := make([]byte, 8)
	if _, randErr := rand.Read(b); randErr != nil {
		return "", "", fmt.Errorf("生成随机 state 失败: %w", randErr)
	}
	state = hex.EncodeToString(b)

	s.stateMu.Lock()
	s.stateStore[state] = &oauthState{
		timestamp:    time.Now(),
		codeVerifier: oauth.GetCodeVerifier(),
	}
	s.stateMu.Unlock()

	authURL = oauth.GetAuthURL(state)
	return authURL, state, nil
}

// HandleCallback processes the OAuth callback after the user authorizes the application.
// It validates the state token, exchanges the authorization code for tokens,
// and persists the encrypted token to the configuration.
func (s *BackupService) HandleCallback(code, state string) error {
	s.stateMu.Lock()
	stored, ok := s.stateStore[state]
	if !ok {
		s.stateMu.Unlock()
		return errors.New("无效的授权状态")
	}
	codeVerifier := stored.codeVerifier
	delete(s.stateStore, state)
	s.stateMu.Unlock()

	cfg, err := s.repo.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	var appSecret string
	if cfg.EncryptedSecret != "" {
		var decErr error
		appSecret, decErr = pkg.Decrypt(cfg.EncryptedSecret, pkg.GetBackupKey())
		if decErr != nil {
			return fmt.Errorf("解密 app_secret 失败: %w", decErr)
		}
	}

	oauth := alipan.NewAlipanOAuth(cfg.AppID, appSecret, cfg.RedirectURI)
	if appSecret == "" && codeVerifier != "" {
		oauth.SetCodeVerifier(codeVerifier)
	}

	if exchErr := oauth.ExchangeCode(code); exchErr != nil {
		return fmt.Errorf("交换授权码失败: %w", exchErr)
	}

	tokenInfo := oauth.GetTokenInfo()
	tokenJSON, marshalErr := json.Marshal(tokenInfo)
	if marshalErr != nil {
		return fmt.Errorf("序列化令牌失败: %w", marshalErr)
	}

	encToken, encErr := pkg.Encrypt(string(tokenJSON), pkg.GetBackupKey())
	if encErr != nil {
		return fmt.Errorf("加密令牌失败: %w", encErr)
	}

	cfg.EncryptedToken = encToken
	cfg.Status = "authorized"

	// TokenExpiry is stored as *time.Time in the model.
	expiry := time.Now().Add(time.Duration(tokenInfo.ExpiresIn) * time.Second)
	cfg.TokenExpiry = &expiry

	cfg.UpdatedAt = time.Now()

	if saveErr := s.repo.SaveConfig(cfg); saveErr != nil {
		return fmt.Errorf("保存配置失败: %w", saveErr)
	}

	return nil
}

// TriggerBackup performs a manual backup: it runs a SQLite integrity check,
// creates a VACUUM INTO backup file, uploads it to the cloud drive, and
// persists a backup record. It also cleans up expired backups afterward.
func (s *BackupService) TriggerBackup() (*model.BackupRecord, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg, err := s.repo.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	if cfg.Status != "authorized" {
		return nil, errors.New("网盘尚未完成授权，请先完成 OAuth 授权")
	}

	record := model.BackupRecord{
		FileName:     "",
		CloudFileID:  "",
		FileSize:     0,
		BackupType:   "manual",
		UploadStatus: "pending",
		IntegrityOk:  0,
		IsPreRestore: 0,
		CreatedAt:    time.Now(),
	}

	if createErr := s.repo.CreateRecord(&record); createErr != nil {
		return nil, fmt.Errorf("创建备份记录失败: %w", createErr)
	}

	// Step 1: SQLite integrity check.
	type integrityRow struct {
		Integrity string `gorm:"column:integrity_check"`
	}
	var row integrityRow
	if execErr := s.db.Raw("PRAGMA integrity_check").Scan(&row).Error; execErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("完整性检查执行失败: %v", execErr)
		s.db.Save(&record)
		return &record, fmt.Errorf("完整性检查执行失败: %w", execErr)
	}

	if row.Integrity != "ok" {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("完整性检查失败: %s", row.Integrity)
		s.db.Save(&record)
		return &record, errors.New(record.ErrorMessage)
	}

	// Step 2: VACUUM INTO to create a clean backup file.
	if mkdirErr := os.MkdirAll(s.backupDir, 0755); mkdirErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("创建备份目录失败: %v", mkdirErr)
		s.db.Save(&record)
		return &record, mkdirErr
	}

	fileName := fmt.Sprintf("warmisle_backup_%s.db", time.Now().Format("20060102_150405"))
	backupPath := filepath.Join(s.backupDir, fileName)
	// Escape single quotes in the path for SQL safety.
	sql := fmt.Sprintf("VACUUM INTO '%s'", stringsReplaceSingleQuote(backupPath))
	if execErr := s.db.Exec(sql).Error; execErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("VACUUM 备份失败: %v", execErr)
		s.db.Save(&record)
		return &record, execErr
	}

	// Step 3: Verify the backup file exists and is non-empty.
	fileInfo, statErr := os.Stat(backupPath)
	if statErr != nil || fileInfo.Size() == 0 {
		errMsg := "备份文件创建失败"
		if statErr != nil {
			errMsg = fmt.Sprintf("备份文件创建失败: %v", statErr)
		}
		record.UploadStatus = "failed"
		record.ErrorMessage = errMsg
		s.db.Save(&record)
		if statErr != nil {
			return &record, statErr
		}
		return &record, errors.New(errMsg)
	}
	fileSize := fileInfo.Size()

	// Step 4: Initialize cloud drive client.
	client, oauthProv, err := s.initCloudClient(cfg)
	if err != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("初始化网盘客户端失败: %v", err)
		s.db.Save(&record)
		return &record, err
	}

	// Step 5: Upload the backup file.
	cloudFile, uploadErr := client.Upload(backupPath, cfg.BackupDir)
	if uploadErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = uploadErr.Error()
		s.db.Save(&record)
		return &record, uploadErr
	}

	// Step 6: After a successful upload, the OAuth tokens may have been refreshed.
	// Persist the refreshed tokens back to the config.
	if err := s.persistRefreshedTokens(cfg, oauthProv); err != nil {
		log.Printf("[backup] 保存刷新后的令牌失败: %v", err)
		// Non-fatal — upload already succeeded.
	}

	// Step 7: Update and save the backup record.
	record.CloudFileID = cloudFile.FileID
	record.FileSize = fileSize
	record.FileName = fileName
	record.UploadStatus = "completed"
	record.IntegrityOk = 1
	s.db.Save(&record)

	// Clean up expired backups in the background (best-effort).
	defer func() { s.cleanExpiredBackups() }()

	return &record, nil
}

// RestoreBackup downloads a backup file from the cloud, creates a pre-restore
// emergency backup of the current database, replaces the database file, and
// triggers an application restart. The confirmText must match exactly the
// Chinese confirmation string to prevent accidental restores.
func (s *BackupService) RestoreBackup(cloudFileID, confirmText string) error {
	if confirmText != "我已了解风险，确认恢复" {
		return errors.New("确认文字不匹配")
	}

	s.restoreMu.Lock()
	if s.restoring {
		s.restoreMu.Unlock()
		return errors.New("正在恢复中，请稍后重试")
	}
	s.restoring = true
	s.restoreMu.Unlock()

	defer func() {
		s.restoreMu.Lock()
		s.restoring = false
		s.restoreMu.Unlock()
	}()

	// Initialize cloud drive client.
	cfg, err := s.repo.GetConfig()
	if err != nil {
		return fmt.Errorf("获取配置失败: %w", err)
	}

	client, _, err := s.initCloudClient(cfg)
	if err != nil {
		return fmt.Errorf("初始化网盘客户端失败: %w", err)
	}

	// Step 1: Download the backup file to a temporary location.
	tmpDir := filepath.Join(s.backupDir, ".restore-tmp")
	if mkdirErr := os.MkdirAll(tmpDir, 0755); mkdirErr != nil {
		return fmt.Errorf("创建临时目录失败: %w", mkdirErr)
	}

	tmpFile := filepath.Join(tmpDir, fmt.Sprintf("warmisle_restore_%s.db", time.Now().Format("20060102_150405")))

	if dlErr := client.Download(cloudFileID, tmpFile); dlErr != nil {
		return fmt.Errorf("下载备份文件失败: %w", dlErr)
	}

	// Step 2: Verify the downloaded file.
	info, statErr := os.Stat(tmpFile)
	if statErr != nil || info.Size() == 0 {
		return errors.New("下载文件校验失败")
	}

	// Step 3: Create an emergency pre-restore backup of the current database.
	preName := fmt.Sprintf("warmisle_pre_restore_%s.db", time.Now().Format("20060102_150405"))
	prePath := filepath.Join(s.backupDir, preName)
	preSQL := fmt.Sprintf("VACUUM INTO '%s'", stringsReplaceSingleQuote(prePath))
	if execErr := s.db.Exec(preSQL).Error; execErr != nil {
		return fmt.Errorf("创建恢复前备份失败: %w", execErr)
	}

	preRecord := model.BackupRecord{
		FileName:     preName,
		BackupType:   "pre_restore",
		UploadStatus: "completed",
		FileSize:     0,
		IsPreRestore: 1,
		IntegrityOk:  0,
		CreatedAt:    time.Now(),
	}
	if createErr := s.repo.CreateRecord(&preRecord); createErr != nil {
		log.Printf("[backup] 创建恢复前备份记录失败: %v", createErr)
	}

	// Step 4: Close the current database connection.
	sqlDB, dbErr := s.db.DB()
	if dbErr != nil {
		return fmt.Errorf("获取底层数据库连接失败: %w", dbErr)
	}
	if closeErr := sqlDB.Close(); closeErr != nil {
		return fmt.Errorf("关闭数据库连接失败: %w", closeErr)
	}

	// Step 5: Replace the database file with the downloaded backup.
	src, openErr := os.Open(tmpFile)
	if openErr != nil {
		return fmt.Errorf("打开下载文件失败: %w", openErr)
	}
	_ = src.Close()

	dst, createErr := os.Create(s.dbPath)
	if createErr != nil {
		_ = src.Close()
		return fmt.Errorf("创建目标数据库文件失败: %w", createErr)
	}

	if _, copyErr := io.Copy(dst, src); copyErr != nil {
		_ = dst.Close()
		_ = src.Close()
		return fmt.Errorf("复制数据库文件失败: %w", copyErr)
	}
	_ = src.Close()
	_ = dst.Close()

	// Step 6: Write a marker file to signal the application to perform
	// post-restore initialization on next startup.
	if mkdirErr := os.MkdirAll("data", 0755); mkdirErr != nil {
		return fmt.Errorf("创建 data 目录失败: %w", mkdirErr)
	}
	if writeErr := os.WriteFile("data/.restore_complete", []byte(time.Now().String()), 0644); writeErr != nil {
		return fmt.Errorf("写入恢复标记文件失败: %w", writeErr)
	}

	// Give the file system a moment to flush, then exit to trigger restart.
	time.Sleep(500 * time.Millisecond)
	os.Exit(0)

	return nil
}

// ScheduleBackup performs an automatic backup initiated by the scheduler.
// It is identical to TriggerBackup except the backup type is "scheduled"
// and failures are logged rather than returned as errors.
func (s *BackupService) ScheduleBackup() {
	s.mu.Lock()
	defer s.mu.Unlock()

	cfg, err := s.repo.GetConfig()
	if err != nil {
		log.Printf("[backup] 定时备份: 获取配置失败: %v", err)
		return
	}

	if cfg.Status != "authorized" {
		log.Printf("[backup] 定时备份: 网盘尚未授权，跳过")
		return
	}

	record := model.BackupRecord{
		FileName:     "",
		CloudFileID:  "",
		FileSize:     0,
		BackupType:   "scheduled",
		UploadStatus: "pending",
		IntegrityOk:  0,
		IsPreRestore: 0,
		CreatedAt:    time.Now(),
	}

	if createErr := s.repo.CreateRecord(&record); createErr != nil {
		log.Printf("[backup] 定时备份: 创建记录失败: %v", createErr)
		return
	}

	// Integrity check.
	type integrityRow struct {
		Integrity string `gorm:"column:integrity_check"`
	}
	var row integrityRow
	if execErr := s.db.Raw("PRAGMA integrity_check").Scan(&row).Error; execErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("完整性检查执行失败: %v", execErr)
		s.db.Save(&record)
		log.Printf("[backup] 定时备份: %v", execErr)
		return
	}

	if row.Integrity != "ok" {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("完整性检查失败: %s", row.Integrity)
		s.db.Save(&record)
		log.Printf("[backup] 定时备份: 完整性检查失败: %s", row.Integrity)
		return
	}

	// VACUUM INTO.
	if mkdirErr := os.MkdirAll(s.backupDir, 0755); mkdirErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("创建备份目录失败: %v", mkdirErr)
		s.db.Save(&record)
		log.Printf("[backup] 定时备份: %v", mkdirErr)
		return
	}

	fileName := fmt.Sprintf("warmisle_backup_%s.db", time.Now().Format("20060102_150405"))
	backupPath := filepath.Join(s.backupDir, fileName)
	sql := fmt.Sprintf("VACUUM INTO '%s'", stringsReplaceSingleQuote(backupPath))
	if execErr := s.db.Exec(sql).Error; execErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("VACUUM 备份失败: %v", execErr)
		s.db.Save(&record)
		log.Printf("[backup] 定时备份: %v", execErr)
		return
	}

	fileInfo, statErr := os.Stat(backupPath)
	if statErr != nil || fileInfo.Size() == 0 {
		errMsg := "备份文件创建失败"
		if statErr != nil {
			errMsg = fmt.Sprintf("备份文件创建失败: %v", statErr)
		}
		record.UploadStatus = "failed"
		record.ErrorMessage = errMsg
		s.db.Save(&record)
		log.Printf("[backup] 定时备份: %s", errMsg)
		return
	}
	fileSize := fileInfo.Size()

	// Initialize cloud client.
	client, oauthProv, err := s.initCloudClient(cfg)
	if err != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = fmt.Sprintf("初始化网盘客户端失败: %v", err)
		s.db.Save(&record)
		log.Printf("[backup] 定时备份: %v", err)
		return
	}

	// Upload.
	cloudFile, uploadErr := client.Upload(backupPath, cfg.BackupDir)
	if uploadErr != nil {
		record.UploadStatus = "failed"
		record.ErrorMessage = uploadErr.Error()
		s.db.Save(&record)
		log.Printf("[backup] 定时备份: 上传失败: %v", uploadErr)
		return
	}

	// Persist refreshed tokens.
	if persistErr := s.persistRefreshedTokens(cfg, oauthProv); persistErr != nil {
		log.Printf("[backup] 定时备份: 保存刷新后的令牌失败: %v", persistErr)
	}

	// Update record.
	record.CloudFileID = cloudFile.FileID
	record.FileSize = fileSize
	record.FileName = fileName
	record.UploadStatus = "completed"
	record.IntegrityOk = 1
	s.db.Save(&record)

	log.Printf("[backup] 定时备份: 成功上传 %s (%s)", fileName, cloudFile.FileID)

	// Clean up expired backups.
	s.cleanExpiredBackups()
}

// ListCloudFiles returns the list of files in the configured cloud backup directory.
func (s *BackupService) ListCloudFiles() ([]plugin.CloudFileInfo, error) {
	cfg, err := s.repo.GetConfig()
	if err != nil {
		return nil, fmt.Errorf("获取配置失败: %w", err)
	}

	client, _, err := s.initCloudClient(cfg)
	if err != nil {
		return nil, fmt.Errorf("初始化网盘客户端失败: %w", err)
	}

	return client.List(cfg.BackupDir)
}

// ListHistory returns paginated backup history records.
func (s *BackupService) ListHistory(page, pageSize int) ([]model.BackupRecord, int64, error) {
	return s.repo.ListRecords(page, pageSize)
}

// DeleteHistory deletes a backup record by ID. If the record has an associated
// cloud file, it is also deleted from the cloud drive.
func (s *BackupService) DeleteHistory(id uint) error {
	record, err := s.repo.GetRecordByID(id)
	if err != nil {
		return fmt.Errorf("获取备份记录失败: %w", err)
	}

	// If the record has an associated cloud file, delete it.
	if record.CloudFileID != "" {
		cfg, cfgErr := s.repo.GetConfig()
		if cfgErr != nil {
			log.Printf("[backup] 删除云端文件时获取配置失败: %v", cfgErr)
		} else {
			client, _, clientErr := s.initCloudClient(cfg)
			if clientErr != nil {
				log.Printf("[backup] 删除云端文件时初始化客户端失败: %v", clientErr)
			} else {
				if delErr := client.Delete(record.CloudFileID); delErr != nil {
					log.Printf("[backup] 删除云端文件失败 (file_id=%s): %v", record.CloudFileID, delErr)
				}
			}
		}
	}

	return s.repo.DeleteRecord(id)
}

// cleanExpiredBackups removes backup records and their cloud files that are
// older than the configured retention period. This is best-effort and will
// not propagate errors to callers.
func (s *BackupService) cleanExpiredBackups() {
	cfg, err := s.repo.GetConfig()
	if err != nil {
		log.Printf("[backup] cleanExpiredBackups: 获取配置失败: %v", err)
		return
	}

	records, listErr := s.repo.ListExpiredRecords(cfg.RetentionDays)
	if listErr != nil {
		log.Printf("[backup] cleanExpiredBackups: 查询过期记录失败: %v", listErr)
		return
	}

	if len(records) == 0 {
		return
	}

	// Initialize cloud client only once for batch deletion.
	client, _, clientErr := s.initCloudClient(cfg)
	if clientErr != nil {
		log.Printf("[backup] cleanExpiredBackups: 初始化客户端失败，仅清理本地记录: %v", clientErr)
		// Still delete local records even if cloud client init fails.
		for _, r := range records {
			if delErr := s.repo.DeleteRecord(r.ID); delErr != nil {
				log.Printf("[backup] cleanExpiredBackups: 删除本地记录失败 (id=%d): %v", r.ID, delErr)
			}
		}
		return
	}

	for _, r := range records {
		if r.CloudFileID != "" {
			if delErr := client.Delete(r.CloudFileID); delErr != nil {
				log.Printf("[backup] cleanExpiredBackups: 删除云端文件失败 (file_id=%s): %v", r.CloudFileID, delErr)
			}
		}
		if delErr := s.repo.DeleteRecord(r.ID); delErr != nil {
			log.Printf("[backup] cleanExpiredBackups: 删除本地记录失败 (id=%d): %v", r.ID, delErr)
		}
	}

	log.Printf("[backup] cleanExpiredBackups: 清理了 %d 条过期记录", len(records))
}

// ---------------------------------------------------------------------------
// Internal helpers
// ---------------------------------------------------------------------------

// initCloudClient decrypts credentials from the config, initializes the
// AlipanOAuth provider and AlipanClient, and returns the client along with
// a reference to the oauth provider for token persistence.
func (s *BackupService) initCloudClient(cfg *model.CloudDriveConfig) (*alipan.AlipanClient, *alipan.AlipanOAuth, error) {
	appSecret, err := pkg.Decrypt(cfg.EncryptedSecret, pkg.GetBackupKey())
	if err != nil {
		return nil, nil, fmt.Errorf("解密 app_secret 失败: %w", err)
	}

	tokenJSON, err := pkg.Decrypt(cfg.EncryptedToken, pkg.GetBackupKey())
	if err != nil {
		return nil, nil, fmt.Errorf("解密令牌失败: %w", err)
	}

	var token alipan.AlipanToken
	if unmarshalErr := json.Unmarshal([]byte(tokenJSON), &token); unmarshalErr != nil {
		return nil, nil, fmt.Errorf("解析令牌失败: %w", unmarshalErr)
	}

	oauth := alipan.NewAlipanOAuth(cfg.AppID, appSecret, cfg.RedirectURI)

	// cfg.TokenExpiry is *time.Time; LoadTokens takes time.Time.
	tokenExpiry := time.Time{}
	if cfg.TokenExpiry != nil {
		tokenExpiry = *cfg.TokenExpiry
	}
	oauth.LoadTokens(token.AccessToken, token.RefreshToken, tokenExpiry)

	client := alipan.NewAlipanClient(oauth)
	return client, oauth, nil
}

// persistRefreshedTokens serializes the current token state from the OAuth
// provider and saves it to the cloud drive configuration.
func (s *BackupService) persistRefreshedTokens(cfg *model.CloudDriveConfig, oauth *alipan.AlipanOAuth) error {
	newToken := oauth.GetTokenInfo()
	newTokenJSON, marshalErr := json.Marshal(newToken)
	if marshalErr != nil {
		return fmt.Errorf("序列化令牌失败: %w", marshalErr)
	}

	encToken, encErr := pkg.Encrypt(string(newTokenJSON), pkg.GetBackupKey())
	if encErr != nil {
		return fmt.Errorf("加密令牌失败: %w", encErr)
	}

	cfg.EncryptedToken = encToken
	// oauth.TokenExpiry is time.Time; cfg.TokenExpiry is *time.Time.
	expiry := oauth.TokenExpiry
	cfg.TokenExpiry = &expiry

	return s.repo.SaveConfig(cfg)
}

// stringsReplaceSingleQuote escapes single quotes in a string by doubling them.
// This is required for safely embedding file paths in SQLite SQL statements.
func stringsReplaceSingleQuote(s string) string {
	result := make([]byte, 0, len(s)+8)
	for i := 0; i < len(s); i++ {
		if s[i] == '\'' {
			result = append(result, '\'', '\'')
		} else {
			result = append(result, s[i])
		}
	}
	return string(result)
}

// UpdateSchedule updates the scheduled backup configuration without touching
// OAuth credentials or other config fields.
func (s *BackupService) UpdateSchedule(scheduleEnabled bool, scheduleTime string, retentionDays int) error {
	cfg, err := s.repo.GetConfig()
	if err != nil {
		return err
	}

	if scheduleEnabled {
		cfg.ScheduleEnabled = 1
	} else {
		cfg.ScheduleEnabled = 0
	}

	if scheduleTime != "" {
		cfg.ScheduleTime = scheduleTime
	}

	if retentionDays > 0 {
		cfg.RetentionDays = retentionDays
	}

	cfg.UpdatedAt = time.Now()
	return s.repo.SaveConfig(cfg)
}
