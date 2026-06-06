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
	"crypto/md5"
	"time"

	"github.com/gin-gonic/gin"
	_ "modernc.org/sqlite"
	"github.com/pressly/goose/v3"

	"warmisle/internal/pkg"
	"warmisle/internal/routes"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

//go:embed migrations/*.sql
var migrationFS embed.FS

func main() {
	// 读取配置
	dbPath := getEnv("HC_DB_PATH", "./data/warmisle.db")
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

	// 迁移由 goose 管理，无需 GORM AutoMigrate
	// 如需添加字段，创建新的 goose 迁移文件

	r := gin.Default()

	// API 路由
	routes.Register(r)

	// 前端静态文件（从 embed 提供）
	dist, _ := fs.Sub(frontendFS, "frontend/dist")
	r.Use(func(c *gin.Context) {
		// API 请求跳过静态文件处理
		if strings.HasPrefix(c.Request.URL.Path, "/api") {
			c.Next()
			return
		}
		// 尝试读取静态文件，SPA fallback 到 index.html
		path := c.Request.URL.Path
		if path == "/" {
			path = "/index.html"
		}
		data, err := fs.ReadFile(dist, strings.TrimPrefix(path, "/"))
		if err != nil {
			// SPA fallback: 前端路由交由 index.html 处理
			data, _ = fs.ReadFile(dist, "index.html")
			path = "/index.html"
		}

		// ETag 支持
		hash := md5.Sum(data)
		etag := `"` + hex.EncodeToString(hash[:]) + `"`
		if match := c.GetHeader("If-None-Match"); match == etag {
			c.Status(http.StatusNotModified)
			c.Abort()
			return
		}
		c.Header("ETag", etag)

		// Cache-Control: Vite 构建的文件带内容哈希，可长期缓存；index.html 不缓存
		if strings.HasSuffix(path, ".html") || path == "/index.html" {
			c.Header("Cache-Control", "no-cache")
		} else {
			c.Header("Cache-Control", "public, max-age=31536000, immutable")
		}

		c.Data(http.StatusOK, getContentType(path), data)
		c.Abort()
	})

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

func getContentType(path string) string {
	ext := filepath.Ext(path)
	// 常见的前端资源 MIME 类型
	switch ext {
	case ".html":
		return "text/html; charset=utf-8"
	case ".js":
		return "application/javascript; charset=utf-8"
	case ".css":
		return "text/css; charset=utf-8"
	case ".json":
		return "application/json; charset=utf-8"
	case ".png":
		return "image/png"
	case ".jpg", ".jpeg":
		return "image/jpeg"
	case ".svg":
		return "image/svg+xml"
	case ".ico":
		return "image/x-icon"
	case ".woff":
		return "font/woff"
	case ".woff2":
		return "font/woff2"
	default:
		return "text/plain; charset=utf-8"
	}
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

	sqlDB, err := sql.Open("sqlite", dbPath)
	if err != nil {
		return err
	}
	defer sqlDB.Close()

	goose.SetBaseFS(migrationFS)
	goose.SetDialect("sqlite3")
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
