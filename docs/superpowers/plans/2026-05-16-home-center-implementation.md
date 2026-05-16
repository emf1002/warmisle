# 家庭数字中心 V1 Implementation Plan

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 实现一个包含记账、待办、愿望清单、论坛、仪表盘的家庭协作平台，Go + Vue 单二进制部署。

**Architecture:** 经典分层架构 (handler → service → repository) + 前端 Vue 3。Go 后端嵌入前端构建产物，编译为单二进制。SQLite WAL 模式存储，GORM 操作数据库，goose 管理迁移。

**Tech Stack:** Go 1.22+, Gin, GORM, SQLite, goose, Vue 3, Ant Design Vue, Vite, Pinia, Axios

**修订记录**：

| 版本 | 日期 | 修改内容 |
|------|------|----------|
| V1.1.0 | 2026-05-16 | 同步技术设计 V1.0.1 和 UI 风格指南 V1.1.1：1) 预置数据移至 migration 创建；2) 新增数据库索引任务；3) CLI 补充解锁同步；4) 仪表盘 API 补充 month 参数；5) 信息流 API 补充分页/响应结构；6) 评论/点赞 API 补充请求体定义；7) JWT 中间件补充 disabled 校验说明

---

## Phase 0: 项目脚手架

### Task 0.1: 初始化后端项目

**Files:**
- Create: `backend/go.mod`
- Create: `backend/cmd/server/main.go`
- Create: `backend/internal/pkg/response.go`
- Create: `backend/internal/routes/router.go`

- [ ] **Step 1: 初始化 Go module**

```bash
cd backend
go mod init home-center
go get github.com/gin-gonic/gin
go get gorm.io/gorm
go get gorm.io/driver/sqlite
go get github.com/golang-jwt/jwt/v5
go get golang.org/x/crypto/bcrypt
go get github.com/pressly/goose/v3
```

- [ ] **Step 2: 创建入口文件**

`backend/cmd/server/main.go`:
```go
package main

import (
    "embed"
    "home-center/internal/routes"
    "io/fs"
    "net/http"

    "github.com/gin-gonic/gin"
)

//go:embed frontend/dist/*
var frontendFS embed.FS

func main() {
    r := gin.Default()

    // API 路由
    routes.Register(r)

    // 前端静态文件
    dist, _ := fs.Sub(frontendFS, "frontend/dist")
    r.StaticFS("/", http.FS(dist))

    r.Run(":8080")
}
```

- [ ] **Step 3: 统一响应工具**

`backend/internal/pkg/response.go`:
```go
package pkg

import "github.com/gin-gonic/gin"

type Response struct {
    Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data"`
}

type PageData struct {
    List     interface{} `json:"list"`
    Total    int64       `json:"total"`
    Page     int         `json:"page"`
    PageSize int         `json:"page_size"`
}

func Success(c *gin.Context, data interface{}) {
    c.JSON(200, Response{Code: 0, Message: "ok", Data: data})
}

func Error(c *gin.Context, httpCode, bizCode int, msg string) {
    c.JSON(httpCode, Response{Code: bizCode, Message: msg})
}
```

- [ ] **Step 4: 路由空壳**

`backend/internal/routes/router.go`:
```go
package routes

import "github.com/gin-gonic/gin"

func Register(r *gin.Engine) {
    api := r.Group("/api")
    // 各模块路由后续注册
    _ = api
}
```

- [ ] **Step 5: 验证编译通过**

```bash
cd backend && go build ./cmd/server/
```

- [ ] **Step 6: Commit**

```bash
git init
git add -A
git commit -m "chore: initialize Go backend scaffold"
```

### Task 0.2: 初始化前端项目

**Files:**
- Create: `frontend/package.json`
- Create: `frontend/vite.config.ts`
- Create: `frontend/index.html`
- Create: `frontend/src/main.ts`
- Create: `frontend/src/App.vue`
- Create: `frontend/tsconfig.json`

- [ ] **Step 1: 创建 Vue 3 项目**

```bash
cd frontend
npm create vite@latest . -- --template vue-ts
npm install vue-router@4 pinia axios ant-design-vue @ant-design/icons-vue
```

- [ ] **Step 2: 配置 vite.config.ts**

```typescript
import { defineConfig } from 'vite'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  base: '/',
  server: {
    port: 3000,
    proxy: {
      '/api': 'http://localhost:8080'
    }
  },
  build: {
    outDir: 'dist',
    assetsDir: 'assets'
  }
})
```

- [ ] **Step 3: 验证构建**

```bash
cd frontend && npm run build
```

- [ ] **Step 4: 建立 embed 路径**

确保前端的 build 输出目录能被后端 embed 找到。在 `backend/cmd/server/main.go` 中，`//go:embed frontend/dist/*` 路径对应关系为：

```
home-center-v1/
├── backend/
│   └── cmd/server/main.go   ← 这里 embed "frontend/dist/*"
└── frontend/
    └── dist/                ← 构建产物
```

需要调整 go.work 或使用相对路径。推荐方式：在 `backend/` 目录下 `go run ./cmd/server/`，Vite 构建输出到 `backend/frontend/dist/`。

修改 `frontend/vite.config.ts` 的 `outDir`：`../backend/frontend/dist`。

- [ ] **Step 5: Commit**

```bash
git add -A
git commit -m "chore: initialize Vue 3 frontend scaffold"
```

### Task 0.3: Makefile 与开发工具链

**Files:**
- Create: `Makefile`

- [ ] **Step 1: 创建 Makefile**

```makefile
.PHONY: dev build clean

dev:
	@echo "Starting dev server..."
	@cd frontend && npm run dev & cd backend && air -- cmd/server/main.go

build:
	@cd frontend && npm run build
	@cd backend && go build -o home-center ./cmd/server/

clean:
	@rm -rf frontend/dist backend/home-center
```

- [ ] **Step 2: 创建 .gitignore**

```
frontend/node_modules/
frontend/dist/
backend/home-center
backend/tmp/
*.db
*.db-shm
*.db-wal
.env
```

- [ ] **Step 3: Commit**

```bash
git add -A
git commit -m "chore: add Makefile and .gitignore"
```

---

## Phase 1: 数据库与模型层

### Task 1.1: 数据库连接与初始化

**Files:**
- Create: `backend/internal/pkg/database.go`
- Create: `backend/internal/model/member.go`
- Create: `backend/migrations/001_init.up.sql`
- Create: `backend/migrations/001_init.down.sql`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: 数据库连接**

`backend/internal/pkg/database.go`:
```go
package pkg

import (
    "gorm.io/driver/sqlite"
    "gorm.io/gorm"
    "gorm.io/gorm/logger"
)

var DB *gorm.DB

func InitDatabase(dbPath string) error {
    var err error
    DB, err = gorm.Open(sqlite.Open(dbPath+"?_journal_mode=WAL"), &gorm.Config{
        Logger: logger.Default.LogMode(logger.Info),
    })
    return err
}
```

- [ ] **Step 2: 迁移工具函数**

在 `backend/cmd/server/main.go` 中，启动时先备份再调用 goose 迁移：
```go
import (
    "database/sql"
    "embed"
    "fmt"
    "io"
    "os"
    "path/filepath"
    "sort"
    "time"
    "github.com/pressly/goose/v3"
)

//go:embed migrations/*.sql
var migrationFS embed.FS

func backupDB(dbPath string) error {
    backupDir := "backups"
    if err := os.MkdirAll(backupDir, 0755); err != nil {
        return err
    }
    backupFile := filepath.Join(backupDir,
        fmt.Sprintf("backup-%s.db", time.Now().Format("20060102_150405")))
    src, err := os.Open(dbPath)
    if err != nil {
        if os.IsNotExist(err) { return nil } // 首次启动无现有 DB
        return err
    }
    defer src.Close()
    dst, err := os.Create(backupFile)
    if err != nil { return err }
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
    if err != nil { return err }
    defer sqlDB.Close()

    goose.SetBaseFS(migrationFS)
    if err := goose.Up(sqlDB, "migrations"); err != nil {
        return fmt.Errorf("migration failed: %w", err)
    }
    return nil
}

func cleanupBackups() {
    entries, _ := filepath.Glob("backups/backup-*.db")
    if len(entries) <= 7 { return }
    // 按文件名排序（含时间戳），保留最近 7 份
    sort.Strings(entries)
    for _, f := range entries[:len(entries)-7] {
        os.Remove(f)
    }
}
```

> 迁移失败时启动中止，日志输出错误信息，提示用户手动用 `backups/` 目录下的备份文件恢复。

- [ ] **Step 3: 创建迁移文件**

`backend/migrations/001_init.up.sql` — 创建所有基础表（见下方 DDL 汇总）。

- [ ] **Step 4: 注册模型 AutoMigrate**

在 `InitDatabase` 后执行 `DB.AutoMigrate(&model.Member{}, ...)` 作为二次校验。

- [ ] **Step 5: 数据库索引**

在迁移文件 `001_init.up.sql` 中创建以下索引（技术设计 3.4 节）：

| 表 | 索引 | 理由 |
|------|------|------|
| ledgers | `(deleted_at, occurred_at)` | 按月查询+排序 |
| ledgers | `(creator_id, deleted_at)` | 按记录者筛选 |
| ledgers | `(category_id, deleted_at)` | 按分类筛选+删除校验 |
| ledger_members | `(member_id)` | 按成员筛选 |
| todos | `(deleted_at, status, priority, due_date)` | 列表排序 |
| todos | `(assignee_id, deleted_at)` | 按指派人筛选 |
| posts | `(deleted_at, created_at)` | 信息流排序 |
| topics | `(deleted_at, created_at)` | 信息流排序 |
| topics | `(is_pinned, deleted_at, created_at)` | 公告置顶查询 |
| comments | `(target_type, target_id, deleted_at)` | 评论列表 |
| likes | `(target_type, target_id)` | 点赞数统计 |
| wish_votes | `(wish_id)` | 愿望投票数 |
| vote_records | `(vote_id, member_id)` | 投票去重+结果 |
| members | `(username, deleted_at)` | 登录查询+用户名唯一校验 |

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: add database connection and migration setup"
```

---

## Phase 2: 认证与初始化

### Task 2.1: 认证模块 — 后端

**Files:**
- Create: `backend/internal/model/member.go`
- Create: `backend/internal/pkg/jwt.go`
- Create: `backend/internal/pkg/password.go`
- Create: `backend/internal/middleware/auth.go`
- Create: `backend/internal/repository/auth.go`
- Create: `backend/internal/service/auth.go`
- Create: `backend/internal/handler/auth.go`
- Modify: `backend/internal/routes/router.go`
- Modify: `backend/cmd/server/main.go`

- [ ] **Step 1: Member 数据模型**

`backend/internal/model/member.go`:
```go
package model

import (
    "time"
    "gorm.io/gorm"
)

type Member struct {
    ID        uint           `gorm:"primaryKey"`
    Username  string         `gorm:"uniqueIndex;size:20"`
    Password  string         `gorm:"size:255"`
    Name      string         `gorm:"size:20"`
    Avatar    string         `gorm:"size:10;default:👨"`
    Role      string         `gorm:"size:10;default:member"`
    Status    string         `gorm:"size:10;default:active"`
    LastLogin *time.Time
    CreatedAt time.Time
    UpdatedAt time.Time
    DeletedAt gorm.DeletedAt `gorm:"index"`
}

func (Member) TableName() string { return "members" }
```

- [ ] **Step 2: JWT 工具**

`backend/internal/pkg/jwt.go`:
```go
package pkg

import (
    "time"
    "github.com/golang-jwt/jwt/v5"
)

var jwtSecret []byte

func InitJWT(secret string) {
    jwtSecret = []byte(secret)
}

type Claims struct {
    MemberID uint   `json:"member_id"`
    Username string `json:"username"`
    Role     string `json:"role"`
    jwt.RegisteredClaims
}

func GenerateToken(memberID uint, username, role string) (string, error) {
    claims := Claims{
        MemberID: memberID,
        Username: username,
        Role:     role,
        RegisteredClaims: jwt.RegisteredClaims{
            ExpiresAt: jwt.NewNumericDate(time.Now().Add(7 * 24 * time.Hour)),
        },
    }
    token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
    return token.SignedString(jwtSecret)
}

func ParseToken(tokenStr string) (*Claims, error) {
    token, err := jwt.ParseWithClaims(tokenStr, &Claims{}, func(t *jwt.Token) (interface{}, error) {
        return jwtSecret, nil
    })
    if err != nil { return nil, err }
    claims, ok := token.Claims.(*Claims)
    if !ok || !token.Valid { return nil, jwt.ErrSignatureInvalid }
    return claims, nil
}
```

- [ ] **Step 3: 密码工具**

`backend/internal/pkg/password.go`:
```go
package pkg

import "golang.org/x/crypto/bcrypt"

const DefaultPassword = "home123"

func HashPassword(pwd string) (string, error) {
    bytes, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
    return string(bytes), err
}

func CheckPassword(hash, pwd string) bool {
    return bcrypt.CompareHashAndPassword([]byte(hash), []byte(pwd)) == nil
}
```

- [ ] **Step 4: JWT 中间件**

`backend/internal/middleware/auth.go`:
```go
package middleware

import (
    "strings"
    "github.com/gin-gonic/gin"
    "home-center/internal/pkg"
)

func AuthRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        auth := c.GetHeader("Authorization")
        if !strings.HasPrefix(auth, "Bearer ") {
            pkg.Error(c, 401, 40101, "未登录")
            c.Abort()
            return
        }
        claims, err := pkg.ParseToken(auth[7:])
        if err != nil {
            pkg.Error(c, 401, 40101, "Token 已过期")
            c.Abort()
            return
        }
        // 检查成员是否被禁用 — 通过注入的 repository 查询
        var member model.Member
        if err := pkg.DB.First(&member, claims.MemberID).Error; err != nil || member.Status == "disabled" {
            pkg.Error(c, 403, 40301, "账号已被禁用")
            c.Abort()
            return
        }
        c.Set("member_id", claims.MemberID)
        c.Set("role", claims.Role)
        c.Next()
    }
}

func AdminRequired() gin.HandlerFunc {
    return func(c *gin.Context) {
        role, _ := c.Get("role")
        if role != "admin" {
            pkg.Error(c, 403, 40301, "权限不足")
            c.Abort()
            return
        }
        c.Next()
    }
}
```

- [ ] **Step 5: Repository**

`backend/internal/repository/auth.go`:
```go
package repository

import (
    "home-center/internal/model"
    "home-center/internal/pkg"
)

type AuthRepo struct{}

func (r *AuthRepo) FindByUsername(username string) (*model.Member, error) {
    var m model.Member
    err := pkg.DB.Where("username = ?", username).First(&m).Error
    if err != nil { return nil, err }
    return &m, nil
}

func (r *AuthRepo) FindByID(id uint) (*model.Member, error) {
    var m model.Member
    err := pkg.DB.First(&m, id).Error
    if err != nil { return nil, err }
    return &m, nil
}

func (r *AuthRepo) Count() (int64, error) {
    var count int64
    err := pkg.DB.Model(&model.Member{}).Count(&count).Error
    return count, err
}

func (r *AuthRepo) Create(member *model.Member) error {
    return pkg.DB.Create(member).Error
}

func (r *AuthRepo) UpdatePassword(id uint, hash string) error {
    return pkg.DB.Model(&model.Member{}).Where("id = ?", id).Update("password", hash).Error
}
```

- [ ] **Step 6: Service**

`backend/internal/service/auth.go`:
```go
package service

import "home-center/internal/repository"

type AuthService struct {
    repo *repository.AuthRepo
}

func NewAuthService() *AuthService {
    return &AuthService{repo: &repository.AuthRepo{}}
}

// Login 返回 token
func (s *AuthService) Login(username, password string) (string, error) {
    member, err := s.repo.FindByUsername(username)
    if err != nil { return "", ErrInvalidCredentials }
    if !pkg.CheckPassword(member.Password, password) { return "", ErrInvalidCredentials }
    return pkg.GenerateToken(member.ID, member.Username, member.Role)
}

// InitCheck 确认是否需要初始化
func (s *AuthService) InitCheck() (bool, error) {
    count, err := s.repo.Count()
    return count == 0, err
}
```

- [ ] **Step 6a: 登录锁定策略**

在 `backend/internal/service/auth.go` 中添加内存限流，登录失败计数按用户名存储：

```go
package service

import (
    "sync"
    "time"
)

type loginAttempt struct {
    count      int
    lockUntil  time.Time
}

var (
    attemptMu     sync.Mutex
    loginAttempts = make(map[string]*loginAttempt)
)

const maxAttempts = 5
const lockDuration = 15 * time.Minute

func (s *AuthService) isLocked(username string) bool {
    attemptMu.Lock()
    defer attemptMu.Unlock()
    a, ok := loginAttempts[username]
    if !ok { return false }
    if time.Now().Before(a.lockUntil) { return true }
    delete(loginAttempts, username)
    return false
}

func (s *AuthService) recordFailed(username string) {
    attemptMu.Lock()
    defer attemptMu.Unlock()
    a, ok := loginAttempts[username]
    if !ok {
        loginAttempts[username] = &loginAttempt{count: 1}
        return
    }
    a.count++
    if a.count >= maxAttempts {
        a.lockUntil = time.Now().Add(lockDuration)
    }
}

func (s *AuthService) clearAttempts(username string) {
    attemptMu.Lock()
    defer attemptMu.Unlock()
    delete(loginAttempts, username)
}
```

修改 `Login` 方法：先检查 `isLocked` → 校验密码 → 失败调用 `recordFailed` / 成功调用 `clearAttempts`。

> 服务重启后计数重置（内存存储）。管理员被锁定可通过 CLI 工具 `reset-password` 解除。

- [ ] **Step 7: Handler**

`backend/internal/handler/auth.go`:
```go
package handler

import (
    "net/http"
    "github.com/gin-gonic/gin"
    "home-center/internal/pkg"
    "home-center/internal/service"
)

type AuthHandler struct {
    svc *service.AuthService
}

func NewAuthHandler() *AuthHandler {
    return &AuthHandler{svc: service.NewAuthService()}
}

func (h *AuthHandler) Login(c *gin.Context) {
    var req struct {
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        pkg.Error(c, 400, 40001, "参数错误")
        return
    }
    token, err := h.svc.Login(req.Username, req.Password)
    if err != nil {
        pkg.Error(c, 401, 40101, "用户名或密码错误")
        return
    }
    pkg.Success(c, gin.H{"token": token})
}

func (h *AuthHandler) InitCheck(c *gin.Context) {
    needInit, _ := h.svc.InitCheck()
    pkg.Success(c, gin.H{"need_init": needInit})
}
```

- [ ] **Step 8: 注册路由**

在 `backend/internal/routes/router.go` 中添加：
```go
import "home-center/internal/handler"

func Register(r *gin.Engine) {
    api := r.Group("/api")
    auth := handler.NewAuthHandler()

    api.GET("/init/check", auth.InitCheck)
    api.POST("/auth/login", auth.Login)
}
```

- [ ] **Step 9: 启动时初始化 DB + JWT**

`backend/cmd/server/main.go` 中 main 函数：
```go
func main() {
    pkg.InitDatabase("data.db")
    pkg.InitJWT(os.Getenv("HC_JWT_SECRET")) // 空则自动生成随机 secret

    // 数据库迁移：先 goose 迁移 SQL，再 AutoMigrate 补充校验
    if err := runMigrations("data.db"); err != nil {
        log.Fatalf("migration failed: %v", err)
    }
    pkg.DB.AutoMigrate(&model.Member{}, &model.Category{}, &model.Tag{})

    r := gin.Default()
    routes.Register(r)
    r.Run(":8080")
}
```

- [ ] **Step 10: 验证编译 + 测试登录接口**

```bash
cd backend && go build ./cmd/server/
# 启动后测试
curl -X POST http://localhost:8080/api/auth/login -d '{"username":"admin","password":"home123"}'
```

- [ ] **Step 11: Commit**

```bash
git add -A
git commit -m "feat: add authentication module with JWT login"
```

### Task 2.2: 系统初始化流程

**Files:**
- Create: `backend/internal/service/init.go`
- Create: `backend/internal/handler/init.go`
- Modify: `backend/internal/routes/router.go`
- Modify: `backend/migrations/001_init.up.sql`（补充预置分类和标签的 INSERT OR IGNORE）

- [ ] **Step 1: 初始化 Service**

> **关键变更（技术设计 V1.0.1）**：预置数据（20 个分类 + 10 个标签）由 goose migration 在数据库迁移阶段创建，`POST /api/init/setup` 仅创建管理员。

`backend/internal/service/init.go` — 仅创建管理员：
```go
package service

import (
    "home-center/internal/model"
    "home-center/internal/pkg"
    "gorm.io/gorm"
)

type InitService struct{}

func (s *InitService) Setup(tx *gorm.DB, adminName, username, password string) (*model.Member, error) {
    hash, _ := pkg.HashPassword(password)
    admin := model.Member{
        Username: username,
        Password: hash,
        Name:     adminName,
        Avatar:   "👨",
        Role:     "admin",
        Status:   "active",
    }
    if err := tx.Create(&admin).Error; err != nil {
        return nil, err
    }
    return &admin, nil
}
```

> **事务说明**：`Setup` 接收 `tx *gorm.DB` 参数，调用方需用 `pkg.DB.Transaction(...)` 包裹。事务内仅创建管理员，预置数据已在 migration 中创建。

Handler 调用示例：
```go
func (h *InitHandler) Setup(c *gin.Context) {
    var req struct {
        Name     string `json:"name" binding:"required"`
        Username string `json:"username" binding:"required"`
        Password string `json:"password" binding:"required"`
    }
    if err := c.ShouldBindJSON(&req); err != nil {
        pkg.Error(c, 400, 40001, "参数错误")
        return
    }
    var token string
    err := pkg.DB.Transaction(func(tx *gorm.DB) error {
        admin, err := h.svc.Setup(tx, req.Name, req.Username, req.Password)
        if err != nil {
            return err
        }
        token, _ = pkg.GenerateToken(admin.ID, req.Username, "admin")
        return nil
    })
    if err != nil {
        pkg.Error(c, 500, 40005, "初始化失败: "+err.Error())
        return
    }
    // 首次初始化后需创建预置数据
    pkg.Success(c, gin.H{"token": token})
}
```

预置数据在迁移文件 `backend/migrations/001_init.up.sql` 中通过 `INSERT OR IGNORE` 创建，确保幂等（技术设计 3.3 节）：

```sql
-- 预置分类（支出 15 个 + 收入 5 个）
INSERT OR IGNORE INTO categories (type, name, icon, sort_order, preset, created_at, updated_at) VALUES
('expense', '餐饮', '🍱', 1, 1, datetime('now'), datetime('now')),
-- ... 其余 19 个
('income', '其他收入', '📦', 20, 1, datetime('now'), datetime('now'));

-- 预置标签（10 个）
INSERT OR IGNORE INTO tags (name, preset, created_at) VALUES
('家务', 1, datetime('now')),
-- ... 其余 9 个
('宠物', 1, datetime('now'));
```

预置分类和标签的完整列表详见 PRD 第 6.3 节和 6.7 节。

- [ ] **Step 2: 初始化 Handler + 路由**

`backend/internal/handler/init.go` + 注册 `POST /api/init/setup`

- [ ] **Step 3: 验证完整初始化流程**

- [ ] **Step 4: Commit**

### Task 2.3: 认证模块 — 前端

**Files:**
- Create: `frontend/src/api/auth.ts`
- Create: `frontend/src/stores/auth.ts`
- Create: `frontend/src/router/index.ts`
- Create: `frontend/src/views/auth/Login.vue`
- Create: `frontend/src/views/auth/Init.vue`
- Create: `frontend/src/App.vue` (布局框架)
- Create: `frontend/src/utils/request.ts` (Axios 拦截器)

- [ ] **Step 1: Axios 实例与拦截器**

`frontend/src/utils/request.ts` — 自动携带 Token、401 跳转登录

- [ ] **Step 2: Auth API 与 Store**

`frontend/src/api/auth.ts` — 封装 login/initCheck/setup 请求

`frontend/src/stores/auth.ts` (Pinia) — 管理 token、memberInfo、登录状态

- [ ] **Step 3: 路由配置**

```typescript
const routes = [
  { path: '/login', component: Login },
  { path: '/init', component: Init },
  { path: '/', component: Dashboard, meta: { requiresAuth: true } },
  // ...
]
```

路由守卫：检查 token，无 token 跳转 /login；首次访问检测 need_init 跳转 /init。

- [ ] **Step 4: 登录页 + 初始化页**

`Login.vue` — 用户名+密码表单，登录后跳转仪表盘

`Init.vue` — 创建管理员表单（用户名+姓名+密码）

- [ ] **Step 5: 验证前端登录流程**

```bash
cd frontend && npm run dev
# 浏览器访问 http://localhost:3000，验证登录→仪表盘
```

- [ ] **Step 6: Commit**

```bash
git add -A
git commit -m "feat: add frontend auth pages with login flow"
```

---

## Phase 3: 成员管理

### Task 3.1: 成员管理 — 后端

**Files:**
- Create: `backend/internal/repository/member.go`
- Create: `backend/internal/service/member.go`
- Create: `backend/internal/handler/member.go`
- Modify: `backend/internal/routes/router.go`

- [ ] **Step 1: Member Repository**

`CRUD` + `FindByID` + `List` + `Count` + 查询活动记录方法

- [ ] **Step 2: Member Service**

业务逻辑：添加校验(用户名唯一、密码复杂度)、编辑校验(最后管理员保护)、删除校验(有活动记录则拒绝)、禁用/启用逻辑

- [ ] **Step 3: Member Handler + 路由**

```go
api.GET("/members", memberHandler.List)       // 所有登录用户
api.POST("/members", auth.MemberHandler)       // admin
api.PUT("/members/:id", memberHandler.Update)  // admin
api.DELETE("/members/:id", memberHandler.Delete) // admin
api.PUT("/members/:id/disable", memberHandler.Disable)
api.PUT("/members/:id/enable", memberHandler.Enable)
api.PUT("/members/:id/reset-pwd", memberHandler.ResetPassword)
api.PUT("/profile", memberHandler.UpdateProfile) // 本人
```

- [ ] **Step 4: 验证**

```bash
curl http://localhost:8080/api/members -H "Authorization: Bearer $TOKEN"
```

- [ ] **Step 5: Commit**

### Task 3.2: 成员管理 — 前端

**Files:**
- Create: `frontend/src/views/member/Index.vue`
- Create: `frontend/src/api/member.ts`
- Create: `frontend/src/api/member.ts`

- [ ] **Step 1: Member API**

`frontend/src/api/member.ts` — 封装 CRUD + 禁用/启用/重置密码

- [ ] **Step 2: 成员管理页面**

表格展示成员列表，显示头像、姓名、角色标签、状态标签(已禁用)。管理员显示"添加成员"按钮和操作列。

- [ ] **Step 3: 添加/编辑对话框**

表单：姓名、头像(emoji 选择器)、角色、用户名、初始密码

- [ ] **Step 4: Commit**

---

## Phase 4: 分类管理

### Task 4.1: 分类管理 — 后端

**Files:**
- Create: `backend/internal/repository/category.go`
- Create: `backend/internal/service/category.go`
- Create: `backend/internal/handler/category.go`
- Modify: `backend/internal/routes/router.go`

- [ ] **Step 1: Category Repository**

`backend/internal/repository/category.go`:
```go
package repository

import (
    "home-center/internal/model"
    "home-center/internal/pkg"
)

type CategoryRepo struct{}

func (r *CategoryRepo) List() ([]model.Category, error) {
    var list []model.Category
    err := pkg.DB.Order("type, sort_order").Find(&list).Error
    return list, err
}

func (r *CategoryRepo) FindByID(id uint) (*model.Category, error) {
    var c model.Category
    err := pkg.DB.First(&c, id).Error
    return &c, err
}

func (r *CategoryRepo) Create(c *model.Category) error {
    return pkg.DB.Create(c).Error
}

func (r *CategoryRepo) Update(c *model.Category) error {
    return pkg.DB.Save(c).Error
}

// SoftDelete 返回关联记录数，>0 则拒绝删除
func (r *CategoryRepo) SoftDelete(id uint) (int64, error) {
    var count int64
    pkg.DB.Model(&model.Ledger{}).Where("category_id = ?", id).Count(&count)
    if count > 0 { return count, nil }
    return 0, pkg.DB.Delete(&model.Category{}, id).Error
}
```

- [ ] **Step 2: Category Service**

校验：同 type 下名称唯一、删除检查引用

- [ ] **Step 3: Category Handler + 路由**

```go
api.GET("/categories", categoryHandler.List)
api.POST("/categories", auth.Middleware(categoryHandler.Create))  // admin
api.PUT("/categories/:id", auth.Middleware(categoryHandler.Update)) // admin
api.DELETE("/categories/:id", auth.Middleware(categoryHandler.Delete)) // admin
```

- [ ] **Step 4: 验证 CRUD**

```bash
curl -X POST http://localhost:8080/api/categories \
  -H "Authorization: Bearer $TOKEN" \
  -d '{"type":"expense","name":"测试","icon":"📦","sort_order":99}'
```

- [ ] **Step 5: Commit**

### Task 4.2: 分类管理 — 前端

**Files:**
- Create: `frontend/src/views/category/Index.vue`

(表格 + 添加/编辑对话框 + 按收支分组展示)

---

## Phase 5: 记账本

### Task 5.1: 记账本 — 后端

**Files:**
- Create: `backend/internal/model/ledger.go`
- Create: `backend/internal/repository/ledger.go`
- Create: `backend/internal/service/ledger.go`
- Create: `backend/internal/handler/ledger.go`
- Modify: `backend/internal/routes/router.go`

- [ ] **Step 1: Ledger 模型 + LedgerMember 模型**

```go
type Ledger struct {
    ID         uint           `gorm:"primaryKey"`
    Amount     int64          `gorm:"not null"`
    Note       string         `gorm:"size:200"`
    CategoryID uint           `gorm:"index"`
    CreatorID  uint           `gorm:"index"`
    OccurredAt time.Time      `gorm:"index"`
    CreatedAt  time.Time
    UpdatedAt  time.Time
    DeletedAt  gorm.DeletedAt `gorm:"index"`
    // 关联
    Category   Category       `gorm:"foreignKey:CategoryID"`
    Creator    Member         `gorm:"foreignKey:CreatorID"`
    Members    []Member       `gorm:"many2many:ledger_members;"`
}
```

- [ ] **Step 2: Ledger Repository**

按月查询：`DB.Where("strftime('%Y-%m', occurred_at) = ?", month)`
筛选: member_id / category_id / creator_id 可组合
日期分组: 查询后按 occurred_at 分组（Go 代码中分组，非 SQL）

- [ ] **Step 3: Ledger Service**

校验：金额>0、分类存在、关联成员至少一个
权限：仅创建者和管理员可编辑/删除

- [ ] **Step 4: Ledger Handler + 路由**

`GET /api/ledgers?month=&member_id=&category_id=&creator_id=&page=&page_size=`

响应格式：
```json
{
  "summary": { "income": 10000, "expense": 6000, "balance": 4000 },
  "groups": [
    { "date": "2026-05-16", "daily_total": 500, "items": [...] }
  ],
  "total": 100,
  "page": 1,
  "page_size": 20
}
```

- [ ] **Step 5: 验证记账 CRUD + 筛选**
- [ ] **Step 6: Commit**

### Task 5.2: 记账本 — 前端

**Files:**
- Create: `frontend/src/views/ledger/Index.vue`
- Create: `frontend/src/views/ledger/LedgerForm.vue`
- Create: `frontend/src/api/ledger.ts`

- [ ] **Step 1: Ledger API**
- [ ] **Step 2: 记账列表页**

月度小计 sticky 顶部 → 日期分组列表 → 每条显示: 分类名+金额+备注+记录者头像
日期格式: 今天/昨天/M月D日 周X

日期分组标头视觉规范（UI 风格指南 3.8a）：标头背景 `colorBgLayout`（#F5F5F5），日期文字 14px/500，日小计金额 14px/500 跟随收支颜色

- [ ] **Step 3: 记账表单对话框**

选择分类(默认支出)→输入金额→选择关联成员(多选)→选择日期→备注

- [ ] **Step 4: 筛选器**

月份选择器、成员筛选、分类筛选、记录者筛选

- [ ] **Step 5: Commit**

---

## Phase 6: 待办管理

### Task 6.1: 待办管理 — 后端

**Files:**
- Create: `backend/internal/model/todo.go`
- Create: `backend/internal/repository/todo.go`
- Create: `backend/internal/service/todo.go`
- Create: `backend/internal/handler/todo.go`
- Modify: `backend/internal/routes/router.go`

- [ ] **Step 1: Todo + TodoLog 模型**
- [ ] **Step 2: Todo Repository**

列表排序：优先级(紧急>重要>普通) → 有截止日期升序 → 创建时间倒序

- [ ] **Step 3: Todo Service**

认领逻辑(仅未指派)、完成切换(仅创建者/被指派人/admin)、编辑权限校验
变更日志记录（创建待办时设置 assignee、编辑待办时记录 TodoLog）

- [ ] **Step 4: Todo Handler + 路由**

`GET /api/todos?status=pending|completed&assignee_id=&page=1&page_size=20` — 列表查询
`POST /api/todos` — 创建
`PUT /api/todos/:id` — 编辑（创建者/被指派人/admin）
`DELETE /api/todos/:id` — 删除（创建者/admin）
`PUT /api/todos/:id/toggle` — 完成/恢复切换（创建者/被指派人/admin）
`PUT /api/todos/:id/claim` — 认领（登录用户，仅未指派）

- [ ] **Step 5: Commit**

### Task 6.2: 待办管理 — 前端

**Files:**
- Create: `frontend/src/views/todo/Index.vue`
- Create: `frontend/src/api/todo.ts`

- [ ] **Step 1: Todo API**
- [ ] **Step 2: 待办列表页**

列表显示：标题、优先级标签、指派人、截止日期(过期红)、完成勾选框
筛选：待办/已完成、指派人

- [ ] **Step 3: 编辑对话框 + 认领按钮**

创建者和被指派人可见编辑按钮；未指派待办显示"认领"按钮

- [ ] **Step 4: Commit**

---

## Phase 7: 仪表盘

### Task 7.1: 仪表盘 — 后端

**Files:**
- Create: `backend/internal/handler/dashboard.go`
- Create: `backend/internal/service/dashboard.go`
- Modify: `backend/internal/routes/router.go`

- [ ] **Step 1: Dashboard Service**

5个聚合查询方法，直接使用 GORM 聚合（无需额外 repository 层，因为是只读统计）：
```go
func (s *DashboardService) GetSummary(month string) (map[string]int64, error)
func (s *DashboardService) GetExpenseChart(month string) ([]CategorySum, error)
func (s *DashboardService) GetUpcomingTodos() ([]model.Todo, error)  // LIMIT 5
func (s *DashboardService) GetWishTrends() ([]model.WishVote, error) // 按投票时间
func (s *DashboardService) GetForumHot() (map[string]interface{}, error)
```

- [ ] **Step 2: Dashboard Handler + 路由**

路由注册：
```go
api.GET("/dashboard/summary", dashboardHandler.Summary)           // ?month=2026-05
api.GET("/dashboard/expense-chart", dashboardHandler.ExpenseChart) // ?month=2026-05
api.GET("/dashboard/upcoming-todos", dashboardHandler.UpcomingTodos)
api.GET("/dashboard/wish-trends", dashboardHandler.WishTrends)
api.GET("/dashboard/forum-hot", dashboardHandler.ForumHot)
```

`summary` 和 `expense-chart` 接受可选 `month` 查询参数（格式 `YYYY-MM`），缺省使用当前月份。

- [ ] **Step 3: Commit**

### Task 7.2: 仪表盘 — 前端

**Files:**
- Create: `frontend/src/views/dashboard/Index.vue`
- Create: `frontend/src/api/dashboard.ts`

- [ ] **Step 1: Dashboard API**
- [ ] **Step 2: 仪表盘页面**

三张统计卡片(收入/支出/结余) → 支出饼图(使用 Ant Design Vue Chart 或简单 CSS 实现) → 近期待办列表(5条) → 愿望动态(5条，可点击跳转愿望列表) → 论坛热点

> **移动端愿望清单入口**（UI 风格指南 4.0）：愿望清单在移动端不设独立 Tab，通过仪表盘页面内嵌"愿望动态"区域进入。点击任意愿望卡片跳转完整愿望列表页 `/wish`（顶部有返回按钮）。

- [ ] **Step 3: 月份切换器（仅仪表盘独立，与记账本月份各自独立不联动）**
- [ ] **Step 4: Commit**

---

## Phase 8: 愿望清单

### Task 8.1: 愿望清单 — 后端

**Files:**
- Create: `backend/internal/model/wish.go`
- Create: `backend/internal/repository/wish.go`
- Create: `backend/internal/service/wish.go`
- Create: `backend/internal/handler/wish.go`
- Modify: `backend/internal/routes/router.go`

- [ ] **Step 1: Wish + WishVote 模型**
- [ ] **Step 2: Wish Repository + Service**

状态流转校验（管理员可任意转、创建者仅可标记放弃）
提升为家庭愿望（单向不可逆）
投票（每人一票，可取消）

- [ ] **Step 3: Wish Handler + 路由**

`GET /api/wishes?type=personal|family&status=pending|agreed|achieved|abandoned&creator_id=&page=1&page_size=20` — 列表查询
`POST /api/wishes` — 创建
`PUT /api/wishes/:id` — 编辑（创建者/admin）
`DELETE /api/wishes/:id` — 删除（创建者/admin）
`POST /api/wishes/:id/promote` — 提升为家庭愿望（创建者，单向不可逆）
`PUT /api/wishes/:id/status` — 状态变更（admin 任意/创建者仅放弃）
`POST /api/wishes/:id/vote` — 投票
`DELETE /api/wishes/:id/vote` — 取消投票

- [ ] **Step 4: Commit**

### Task 8.2: 愿望清单 — 前端

**Files:**
- Create: `frontend/src/views/wish/Index.vue`
- Create: `frontend/src/api/wish.ts`

- [ ] **Step 1: Wish API**
- [ ] **Step 2: 愿望列表页**

个人愿望/家庭愿望两个 tab。个人愿望展示创建者标识(只读)，家庭愿望展示投票按钮和人数。
卡片形式展示：标题、分类、金额、优先级标签（复用待办优先级配色：紧急=red、重要=orange、普通=default）、状态标签、投票人数

- [ ] **Step 3: 创建/编辑对话框 + 提升按钮**
- [ ] **Step 4: Commit**

---

## Phase 9: 家庭论坛

### Task 9.1: 论坛 — 后端

**Files:**
- Create: `backend/internal/model/forum.go` (Post, Topic, Vote, VoteOption, VoteRecord, Comment, Like, Tag)
- Create: `backend/internal/repository/forum.go`
- Create: `backend/internal/service/forum.go`
- Create: `backend/internal/handler/forum.go`
- Modify: `backend/internal/routes/router.go`

这是最大的模块，涉及 8 个模型、信息流混合查询。

- [ ] **Step 1: 论坛数据模型**

(已在设计文档中完整定义，按表创建)

- [ ] **Step 2: Forum Repository**

关键查询——信息流（分页 + 置顶公告 + 动态/话题混合）：

```go
func (r *ForumRepo) GetFeed(page, pageSize int) (*FeedResponse, error) {
    // 1. 置顶公告: topics WHERE is_pinned=true AND deleted_at IS NULL ORDER BY created_at DESC
    //    → pinned 返回所有有效公告（不分页）
    // 2. 动态+话题混合: posts + topics (WHERE is_pinned=false) UNION ORDER BY created_at DESC
    //    → items 按分页返回
    // 3. 总条数: total (动态+非置顶话题的未删除总数)
}

type FeedResponse struct {
    Pinned []FeedItem `json:"pinned"`   // 所有有效公告，按创建时间倒序
    Items  []FeedItem `json:"items"`    // 动态+话题混合，按创建时间倒序，分页
    Total  int64      `json:"total"`    // items 总数（不含 pinned）
}
```

`GET /api/feed?page=1&page_size=20`，前端若 `pinned` 为空数组则隐藏公告区。

- [ ] **Step 3: Forum Service**

公告规则：置顶时不可评论/点赞，取消置顶后变为普通话题
投票规则：截止后不可投票；截止前创建者和管理员可删除
评论：2层嵌套

- [ ] **Step 4: Forum Handler + 路由**

- [ ] **Step 4: Forum Handler + 路由**

信息流 `GET /api/feed?page=1&page_size=20` → 公告置顶在前 + 动态/话题混合按时间倒序，响应含 `pinned` + `items` + `total`。

评论请求体（`POST /api/comments`）：
```json
{
  "target_type": "post",
  "target_id": 1,
  "parent_id": null,
  "content": "说得好"
}
```
- `target_type` 枚举：`post` / `topic` / `wish`
- `parent_id` 为 null → 一级评论；指向一级评论 ID → 二级评论
- 后端校验：`parent_id` 指向的评论的 `parent_id` 必须为 null（禁止三级嵌套）

点赞请求体（`POST /api/likes`）：
```json
{
  "target_type": "post",
  "target_id": 1
}
```
`DELETE /api/likes` 请求体同上，通过 `(target_type, target_id, member_id)` 唯一约束定位。

- [ ] **Step 5: Commit**

### Task 9.2: 论坛 — 前端

**Files:**
- Create: `frontend/src/views/forum/Index.vue`
- Create: `frontend/src/views/forum/TopicDetail.vue`
- Create: `frontend/src/api/forum.ts`

- [ ] **Step 1: Forum API**
- [ ] **Step 2: 信息流页**

卡片列表：公告(特殊样式+置顶) → 动态卡片(内容+作者+时间+评论数+点赞) → 话题卡片(标题+摘要+标签+作者+时间)
发帖按钮：发动态/发话题

- [ ] **Step 3: 话题详情页**

话题内容 + 评论列表(2层嵌套) + 评论输入框

- [ ] **Step 4: 投票组件**

创建投票表单(标题+选项+截止日期+单选/多选)
投票卡片(选项列表+投票按钮+实时结果)
截止后展示结果

- [ ] **Step 5: Commit**

### Task 9.3: 标签管理 — 后端+前端

**Files:**
- 后端: handler/tag.go, service/tag.go, repository/tag.go
- 前端: 在论坛模块内嵌标签管理对话框

- [ ] **Step 1: Tag CRUD (admin only)**
- [ ] **Step 2: 前端标签管理**
- [ ] **Step 3: Commit**

---

## Phase 10: 个人中心

### Task 10.1: 个人中心 — 后端

**Files:**
- Create: `backend/internal/handler/profile.go`
- Modify: `backend/internal/routes/router.go`

简化实现——复用 Member 的 repository/service，仅添加：

```go
// GET /api/profile
func (h *ProfileHandler) GetProfile(c *gin.Context) {
    memberID := c.GetUint("member_id")
    // 返回 member 信息
}

// PUT /api/profile
func (h *ProfileHandler) UpdateProfile(c *gin.Context) {
    // 修改姓名和头像
}
```

- [ ] **Step 1: Profile Handler**
- [ ] **Step 2: Commit**

### Task 10.2: 个人中心 — 前端

**Files:**
- Create: `frontend/src/views/profile/Index.vue`

- [ ] **Step 1: 个人中心页面**

头像+姓名+角色标签+用户名(只读) → 列表项：账单分类(跳转)、家庭成员(跳转)、修改个人信息、修改密码、退出登录

- [ ] **Step 2: 修改个人信息对话框**
- [ ] **Step 3: 修改密码对话框**
- [ ] **Step 4: Commit**

---

## Phase 11: 布局与响应式

### Task 11.1: 全局布局

**Files:**
- Create: `frontend/src/layouts/MainLayout.vue`
- Create: `frontend/src/layouts/AuthLayout.vue`
- Modify: `frontend/src/App.vue`

- [ ] **Step 1: AuthLayout**

简洁居中布局，用于登录页和初始化页

- [ ] **Step 2: 桌面端 MainLayout**

左侧侧边栏(Logo + 导航菜单) + 顶部栏(用户头像下拉菜单) + 主内容区(router-view)

- [ ] **Step 3: 移动端 MainLayout**

底部 TabBar(仪表盘/记账/待办/论坛/我的) + 顶部栏 + 主内容区

- [ ] **Step 4: 响应式切换**

CSS media query 断点：>= 768px 侧边栏，< 768px 底栏

- [ ] **Step 5: Commit**

### Task 11.2: 空状态与全局优化

- [ ] **Step 1: 空状态通用组件**

各模块配置对应的空状态文案和引导按钮

- [ ] **Step 2: 全局样式主题**

Ant Design Vue 主题色配置、字体、间距统一

- [ ] **Step 3: Commit**

---

## 附录：Task 依赖关系

```
Task 0.1, 0.2  ← Phase 0 脚手架
    ↓
Task 1.1  ← 数据库
    ↓
Task 2.1, 2.2  ← 认证后端(必须先于任何模块)
    ↓
Task 2.3  ← 认证前端
    ↓
Task 3.1, 4.1  ← 成员+分类后端(可并行)
    ↓       ↓
Task 3.2    Task 4.2  ← 成员+分类前端(可并行)
    ↓
Task 5.1, 6.1, 7.1  ← 记账/待办/仪表盘(可并行)
    ↓       ↓       ↓
Task 5.2    Task 6.2    Task 7.2
    ↓
Task 8.1, 9.1  ← 愿望+论坛(可并行)
    ↓       ↓
Task 8.2    Task 9.2, 9.3
    ↓
Task 10.1, 10.2  ← 个人中心
    ↓
Task 11.1, 11.2  ← 布局与响应式
```

注意：前端各模块页面依赖于布局组件(Task 11.1)，但可以先用简单布局开发，最后统一替换。
