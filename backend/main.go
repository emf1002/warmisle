package main

import (
	"crypto/rand"
	"database/sql"
	"embed"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	_ "github.com/mattn/go-sqlite3"
	"github.com/pressly/goose/v3"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/routes"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

//go:embed migrations/*.sql
var migrationFS embed.FS

func main() {
	// 读取配置
	dbPath := getEnv("HC_DB_PATH", "./data/home-center.db")
	port := getEnv("HC_PORT", "8080")
	jwtSecret := getEnv("HC_JWT_SECRET", "")

	// 确保 data 目录存在
	if err := os.MkdirAll(filepath.Dir(dbPath), 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}

	// 初始化 JWT secret：环境变量为空则自动生成或从文件读取
	if jwtSecret == "" {
		jwtSecret = loadOrGenerateSecret(filepath.Dir(dbPath))
	}
	pkg.InitJWT(jwtSecret)

	// 初始化数据库
	if err := pkg.InitDatabase(dbPath); err != nil {
		log.Fatalf("failed to init database: %v", err)
	}

	// 执行数据库迁移
	if err := runMigrations(dbPath); err != nil {
		log.Fatalf("migration failed: %v", err)
	}

	// GORM AutoMigrate 作为二次校验
	if err := pkg.DB.AutoMigrate(
		&model.Member{},
		&model.Category{},
		&model.Tag{},
		&model.Ledger{},
		&model.LedgerMember{},
		&model.Todo{},
		&model.TodoLog{},
		&model.Wish{},
		&model.WishVote{},
		&model.Post{},
		&model.Topic{},
		&model.Vote{},
		&model.VoteOption{},
		&model.VoteRecord{},
		&model.Comment{},
		&model.Like{},
	); err != nil {
		log.Fatalf("failed to auto-migrate: %v", err)
	}

	r := gin.Default()

	// API 路由
	routes.Register(r)

	// 前端静态文件（从 embed 提供）
	dist, _ := fs.Sub(frontendFS, "frontend/dist")
	r.StaticFS("/", http.FS(dist))

	log.Printf("Server starting on :%s", port)
	if err := r.Run(":" + port); err != nil {
		log.Fatalf("failed to start server: %v", err)
	}
}

func getEnv(key, defaultVal string) string {
	if val := os.Getenv(key); val != "" {
		return val
	}
	return defaultVal
}

func loadOrGenerateSecret(dataDir string) string {
	secretFile := filepath.Join(dataDir, "secret.key")

	// 尝试从文件读取
	if data, err := os.ReadFile(secretFile); err == nil && len(data) > 0 {
		log.Println("JWT secret loaded from file")
		return strings.TrimSpace(string(data))
	}

	// 生成新的随机 secret
	randomBytes := make([]byte, 32)
	if _, err := rand.Read(randomBytes); err != nil {
		log.Fatalf("failed to generate random secret: %v", err)
	}
	secret := hex.EncodeToString(randomBytes)

	// 持久化到文件
	if err := os.MkdirAll(dataDir, 0755); err != nil {
		log.Fatalf("failed to create data directory: %v", err)
	}
	if err := os.WriteFile(secretFile, []byte(secret), 0600); err != nil {
		log.Fatalf("failed to write secret file: %v", err)
	}
	log.Println("JWT secret generated and persisted to file")
	return secret
}

func backupDB(dbPath string) error {
	backupDir := "backups"
	if err := os.MkdirAll(backupDir, 0755); err != nil {
		return err
	}
	backupFile := filepath.Join(backupDir,
		fmt.Sprintf("backup-%s.db", time.Now().Format("20060102_150405")))
	src, err := os.Open(dbPath)
	if err != nil {
		if os.IsNotExist(err) {
			return nil // 首次启动无现有 DB
		}
		return err
	}
	defer src.Close()
	dst, err := os.Create(backupFile)
	if err != nil {
		return err
	}
	defer dst.Close()
	_, err = io.Copy(dst, src)
	return err
}

func runMigrations(dbPath string) error {
	// 迁移前自动备份
	if err := backupDB(dbPath); err != nil {
		log.Printf("warning: backup failed (continuing): %v", err)
	}
	// 清理旧备份（保留最近 7 份）
	cleanupBackups()

	sqlDB, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	goose.SetBaseFS(migrationFS)
	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return fmt.Errorf("migration failed: %w", err)
	}
	return nil
}

func cleanupBackups() {
	entries, _ := filepath.Glob("backups/backup-*.db")
	if len(entries) <= 7 {
		return
	}
	// 按文件名排序（含时间戳），保留最近 7 份
	sort.Strings(entries)
	for _, f := range entries[:len(entries)-7] {
		os.Remove(f)
	}
}
