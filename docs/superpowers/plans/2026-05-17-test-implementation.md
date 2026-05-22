# 家庭数字中心 V1 — 测试实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 为家庭数字中心 V1 建立完整的测试体系，覆盖后端所有 service 层业务逻辑、API handler 集成测试、核心工具函数，以及前端关键组件。

**Architecture:** 采用 Go 标准 `testing` 包 + `testify/assert` 断言库。后端通过 in-memory SQLite 搭建测试数据库，每个测试用例独立初始化。Handler 使用 `httptest.NewRecorder` + Gin test mode。前端使用 vitest + @vue/test-utils 进行组件测试。

**Tech Stack:** Go testing, testify/assert, testify/require, SQLite :memory:, httptest, vitest, @vue/test-utils

**测试策略总则：**
1. Service 层测试：核心业务逻辑全覆盖，使用 in-memory SQLite
2. Handler 层测试：API 请求/响应格式、错误码、权限校验
3. Repository 层：不单独测试（由 service 测试间接覆盖）
4. 工具函数测试：JWT、密码、金额转换等纯函数
5. 前端组件测试：关键交互组件（登录表单、待办列表、记账表单）

---

## Phase 0: 测试基础设施

### Task 0.1: 添加 Go 测试依赖

**Files:**
- Modify: `backend/go.mod`

- [ ] **Step 1: 安装 testify 断言库**

```bash
cd backend && go get github.com/stretchr/testify
```

- [ ] **Step 2: 验证依赖安装成功**

```bash
cd backend && go mod tidy
```

Expected: `go.mod` 中出现 `github.com/stretchr/testify`

- [ ] **Step 3: Commit**

```bash
git add backend/go.mod backend/go.sum
git commit -m "chore: add testify test dependency"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 0.2: 创建测试辅助工具

**Files:**
- Create: `backend/internal/testutil/db.go`

- [ ] **Step 1: 创建测试数据库初始化工具**

`backend/internal/testutil/db.go`:
```go
package testutil

import (
	"home-center/internal/model"
	"home-center/internal/pkg"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

// SetupTestDB 创建 in-memory SQLite 并自动迁移所有模型，设置 pkg.DB 全局变量
func SetupTestDB() *gorm.DB {
	db, err := gorm.Open(sqlite.Open(":memory:"), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		panic("failed to connect test database: " + err.Error())
	}

	// 自动迁移所有模型
	err = db.AutoMigrate(
		&model.Member{},
		&model.Category{},
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
		&model.Tag{},
	)
	if err != nil {
		panic("failed to migrate test database: " + err.Error())
	}

	pkg.DB = db
	return db
}

// TeardownTestDB 清理全局 DB 引用
func TeardownTestDB() {
	pkg.DB = nil
}

// SeedMembers 创建测试成员，返回 ID 列表
func SeedMembers(db *gorm.DB, members []model.Member) []model.Member {
	for i := range members {
		db.Create(&members[i])
	}
	return members
}

// SeedCategories 创建测试分类
func SeedCategories(db *gorm.DB, categories []model.Category) []model.Category {
	for i := range categories {
		db.Create(&categories[i])
	}
	return categories
}

// CreateTestMember 快速创建单个测试成员
func CreateTestMember(db *gorm.DB, username, name, role string) model.Member {
	m := model.Member{
		Username: username,
		Password: "$2a$10$dummyhashedpassword1234567890abcdef",
		Name:     name,
		Avatar:   "👨",
		Role:     role,
		Status:   "active",
	}
	db.Create(&m)
	return m
}

// CreateTestCategory 快速创建单个测试分类
func CreateTestCategory(db *gorm.DB, ctype, name, icon string, sortOrder int) model.Category {
	c := model.Category{
		Type:      ctype,
		Name:      name,
		Icon:      icon,
		SortOrder: sortOrder,
		Preset:    false,
	}
	db.Create(&c)
	return c
}
```

- [ ] **Step 2: 验证编译通过**

```bash
cd backend && go build ./internal/testutil/
```

Expected: 编译成功

- [ ] **Step 3: Commit**

```bash
git add backend/internal/testutil/db.go
git commit -m "feat: add test database utility for in-memory SQLite setup"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

---

## Phase 1: 工具函数测试

### Task 1.1: 密码工具测试

**Files:**
- Create: `backend/internal/pkg/password_test.go`

- [ ] **Step 1: 编写密码哈希与验证测试**

`backend/internal/pkg/password_test.go`:
```go
package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_Success(t *testing.T) {
	hash, err := HashPassword("test123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	// bcrypt hash 以 $2a$ 开头
	assert.Contains(t, hash, "$2a$")
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	require.NoError(t, err)
	// 空密码也能哈希
	assert.NotEmpty(t, hash)
}

func TestCheckPassword_Match(t *testing.T) {
	hash, _ := HashPassword("secure456")
	assert.True(t, CheckPassword(hash, "secure456"))
}

func TestCheckPassword_NoMatch(t *testing.T) {
	hash, _ := HashPassword("secure456")
	assert.False(t, CheckPassword(hash, "wrongpass"))
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	assert.False(t, CheckPassword("not-a-valid-hash", "anything"))
}

func TestDefaultPassword_Value(t *testing.T) {
	assert.Equal(t, "home123", DefaultPassword)
}
```

- [ ] **Step 2: 运行测试验证通过**

```bash
cd backend && go test ./internal/pkg/ -run "TestHash|TestCheck|TestDefault" -v
```

Expected: 7 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/pkg/password_test.go
git commit -m "test: add password hash and verify unit tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 1.2: JWT 工具测试

**Files:**
- Create: `backend/internal/pkg/jwt_test.go`

- [ ] **Step 1: 编写 JWT 生成与解析测试**

`backend/internal/pkg/jwt_test.go`:
```go
package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	InitJWT("test-secret-key-for-testing-purposes")
}

func TestGenerateToken_Success(t *testing.T) {
	token, err := GenerateToken(1, "testuser", "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	// JWT 应包含两个点分隔符
	assert.Contains(t, token, ".")
}

func TestParseToken_Success(t *testing.T) {
	token, _ := GenerateToken(42, "alice", "member")
	claims, err := ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.MemberID)
	assert.Equal(t, "alice", claims.Username)
	assert.Equal(t, "member", claims.Role)
}

func TestParseToken_Expired(t *testing.T) {
	// 使用已过期的 token
	oldSecret := jwtSecret
	jwtSecret = []byte("expired-test-secret")
	token, _ := GenerateToken(1, "test", "member")
	// 清理解析器使用的密钥使其不匹配
	jwtSecret = []byte("different-secret")
	_, err := ParseToken(token)
	assert.Error(t, err)

	// 恢复
	jwtSecret = oldSecret
}

func TestParseToken_InvalidFormat(t *testing.T) {
	_, err := ParseToken("not.a.valid.token")
	assert.Error(t, err)
}

func TestParseToken_EmptyString(t *testing.T) {
	_, err := ParseToken("")
	assert.Error(t, err)
}

func TestGenerateToken_RolePreserved(t *testing.T) {
	token, err := GenerateToken(99, "bob", "member")
	require.NoError(t, err)
	claims, _ := ParseToken(token)
	assert.Equal(t, "member", claims.Role)
	assert.Equal(t, uint(99), claims.MemberID)
}

func TestTokenExpiry_SetTo7Days(t *testing.T) {
	token, _ := GenerateToken(1, "user", "member")
	claims, _ := ParseToken(token)
	// 验证过期时间设置在未来的约 7 天
	exp := claims.ExpiresAt.Time
	assert.True(t, exp.After(time.Now()))
	assert.True(t, exp.Before(time.Now().Add(8*24*time.Hour)))
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/pkg/ -run "TestGenerate|TestParse|TestToken" -v
```

Expected: 所有 JWT 测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/pkg/jwt_test.go
git commit -m "test: add JWT token generation and parsing unit tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 1.3: 响应工具测试

**Files:**
- Create: `backend/internal/pkg/response_test.go`

- [ ] **Step 1: 编写统一响应格式测试**

`backend/internal/pkg/response_test.go`:
```go
package pkg

import (
	"encoding/json"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func TestSuccess_ReturnsCorrectFormat(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, gin.H{"key": "value"})

	assert.Equal(t, 200, w.Code)

	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Equal(t, "ok", resp.Message)
	data, ok := resp.Data.(map[string]interface{})
	assert.True(t, ok)
	assert.Equal(t, "value", data["key"])
}

func TestSuccess_WithNilData(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Success(c, nil)

	assert.Equal(t, 200, w.Code)
	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 0, resp.Code)
	assert.Nil(t, resp.Data)
}

func TestError_ReturnsCorrectFormat(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, 400, 40001, "参数错误")

	assert.Equal(t, 400, w.Code)
	var resp Response
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, 40001, resp.Code)
	assert.Equal(t, "参数错误", resp.Message)
	assert.Nil(t, resp.Data)
}

func TestError_401Unauthorized(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, 401, 40101, "未登录")

	assert.Equal(t, 401, w.Code)
}

func TestError_403Forbidden(t *testing.T) {
	w := httptest.NewRecorder()
	c, _ := gin.CreateTestContext(w)

	Error(c, 403, 40301, "权限不足")

	assert.Equal(t, 403, w.Code)
}

func TestPageData_Structure(t *testing.T) {
	pd := PageData{
		List:     []string{"a", "b"},
		Total:    2,
		Page:     1,
		PageSize: 20,
	}
	assert.Equal(t, int64(2), pd.Total)
	assert.Equal(t, 1, pd.Page)
	assert.Equal(t, 20, pd.PageSize)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/pkg/ -run "TestSuccess|TestError|TestPage" -v
```

Expected: 所有响应工具测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/pkg/response_test.go
git commit -m "test: add response utility unit tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

---

## Phase 2: Service 层测试

### Task 2.1: AuthService 测试

**Files:**
- Create: `backend/internal/service/auth_test.go`

- [ ] **Step 1: 编写认证服务测试**

`backend/internal/service/auth_test.go`:
```go
package service

import (
	"testing"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTest() (*AuthService, func()) {
	db := testutil.SetupTestDB()
	pkg.InitJWT("test-jwt-secret-for-auth-tests")
	svc := NewAuthService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestAuthService_InitCheck_NeedInit(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	needInit, err := svc.InitCheck()
	require.NoError(t, err)
	assert.True(t, needInit, "empty database should need init")
}

func TestAuthService_InitCheck_NoNeedInit(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	// 种子数据：创建一个成员
	hash, _ := pkg.HashPassword("pass123")
	pkg.DB.Create(&model.Member{
		Username: "existing", Password: hash, Name: "Exist", Role: "admin", Status: "active",
	})

	needInit, err := svc.InitCheck()
	require.NoError(t, err)
	assert.False(t, needInit, "database with members should not need init")
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("correctpw")
	pkg.DB.Create(&model.Member{
		Username: "alice", Password: hash, Name: "Alice", Role: "member", Status: "active",
	})

	token, err := svc.Login("alice", "correctpw")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("correctpw")
	pkg.DB.Create(&model.Member{
		Username: "bob", Password: hash, Name: "Bob", Role: "member", Status: "active",
	})

	_, err := svc.Login("bob", "wrongpw")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	_, err := svc.Login("nonexistent", "any");
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_DisabledMember(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("pass123")
	pkg.DB.Create(&model.Member{
		Username: "disabled_user", Password: hash, Name: "Disabled", Role: "member", Status: "disabled",
	})

	_, err := svc.Login("disabled_user", "pass123")
	assert.ErrorIs(t, err, ErrInvalidCredentials, "disabled member should not login")
}

func TestAuthService_Login_LockAfter5Failures(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("secret")
	pkg.DB.Create(&model.Member{
		Username: "locktest", Password: hash, Name: "LockTest", Role: "member", Status: "active",
	})

	// 连续 5 次错误密码
	for i := 0; i < 5; i++ {
		_, err := svc.Login("locktest", "wrong")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	}

	// 第 6 次应该被锁定
	_, err := svc.Login("locktest", "secret")
	assert.ErrorIs(t, err, ErrAccountLocked, "should be locked after 5 failures")
}

func TestAuthService_Login_ResetAfterSuccess(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("goodpw")
	pkg.DB.Create(&model.Member{
		Username: "reset_test", Password: hash, Name: "ResetTest", Role: "member", Status: "active",
	})

	// 4 次失败
	for i := 0; i < 4; i++ {
		svc.Login("reset_test", "wrong")
	}

	// 第 5 次成功 — 应当重置计数
	token, err := svc.Login("reset_test", "goodpw")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// 再次失败 — 应该从 0 重新计数，不会立即锁定
	_, err = svc.Login("reset_test", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	// 确认不会被锁定（只有 1 次失败）
	_, err = svc.Login("reset_test", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials, "should not be locked after only 2 failures")
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/service/ -run "TestAuth" -v
```

Expected: 8 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/auth_test.go
git commit -m "test: add auth service login and lockout tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 2.2: MemberService 测试

**Files:**
- Create: `backend/internal/service/member_test.go`

- [ ] **Step 1: 编写成员管理服务测试**

`backend/internal/service/member_test.go`:
```go
package service

import (
	"testing"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMemberTest() (*MemberService, func()) {
	db := testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewMemberService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestMemberService_Create_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, err := svc.Create("newuser", "password1", "New User", "👩", "member")
	require.NoError(t, err)
	assert.Equal(t, "newuser", member.Username)
	assert.Equal(t, "New User", member.Name)
	assert.Equal(t, "member", member.Role)
	assert.Equal(t, "active", member.Status)
	assert.NotEmpty(t, member.Password)
}

func TestMemberService_Create_DuplicateUsername(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	svc.Create("dup", "pass1234", "First", "", "")
	_, err := svc.Create("dup", "pass5678", "Second", "", "")
	assert.ErrorIs(t, err, ErrUsernameTaken)
}

func TestMemberService_Create_InvalidUsername(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	_, err := svc.Create("ab", "pass1234", "Short", "", "")
	assert.ErrorIs(t, err, ErrInvalidUsername, "username too short")

	_, err = svc.Create("user@name", "pass1234", "Bad", "", "")
	assert.ErrorIs(t, err, ErrInvalidUsername, "username has invalid chars")
}

func TestMemberService_Create_InvalidPassword(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	_, err := svc.Create("validuser", "12345", "Valid", "", "")
	assert.ErrorIs(t, err, ErrInvalidPassword, "password too short")
}

func TestMemberService_Create_DefaultValues(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	// 测试默认头像、角色当传空时的行为
	member, err := svc.Create("minuser", "password1", "", "", "")
	require.NoError(t, err)
	assert.Equal(t, "👨", member.Avatar, "default avatar when empty")
	assert.Equal(t, "member", member.Role, "default role when empty")
}

func TestMemberService_Update_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	created, _ := svc.Create("edit_user", "password1", "Old Name", "👨", "member")

	updated, err := svc.Update(created.ID, "New Name", "👩", "member")
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "👩", updated.Avatar)
}

func TestMemberService_Update_CannotRemoveLastAdmin(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	admin, _ := svc.Create("onlyadmin", "password1", "Admin", "👨", "admin")

	_, err := svc.Update(admin.ID, "", "", "member")
	assert.ErrorIs(t, err, ErrCannotDeleteLastAdmin)
}

func TestMemberService_Delete_WithActivityRecords(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	// 创建一个带活动记录的成员（例如有记账记录）
	member, _ := svc.Create("active_member", "password1", "Active", "👨", "member")
	// 为其创建一条记账记录
	pkg.DB.Create(&model.Ledger{
		Amount: 1000, CategoryID: 1, CreatorID: member.ID,
	})

	err := svc.Delete(member.ID, 1)
	assert.ErrorIs(t, err, ErrHasActivityRecords)
}

func TestMemberService_Disable_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("todisable", "password1", "ToDisable", "👨", "member")

	err := svc.Disable(member.ID, 1) // currentMemberID=1 (admin, not self)
	require.NoError(t, err)

	// 验证状态已变更
	updated, _ := svc.GetProfile(member.ID)
	assert.Equal(t, "disabled", updated.Status)
}

func TestMemberService_Disable_CannotDisableSelf(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("self_disable", "password1", "Self", "👨", "admin")

	err := svc.Disable(member.ID, member.ID)
	assert.ErrorIs(t, err, ErrCannotDisableSelf)
}

func TestMemberService_Disable_CannotDisableLastAdmin(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	admin, _ := svc.Create("lastadmin", "password1", "LastAdmin", "👨", "admin")

	err := svc.Disable(admin.ID, 2) // currentMemberID=2 (other admin)
	assert.ErrorIs(t, err, ErrCannotDeleteLastAdmin)
}

func TestMemberService_Enable_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("toenable", "password1", "ToEnable", "👨", "member")
	svc.Disable(member.ID, 1)

	enabled, err := svc.Enable(member.ID)
	require.NoError(t, err)
	assert.Equal(t, "active", enabled.Status)
}

func TestMemberService_ResetPassword(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("resetpwd", "oldpassword", "Reset", "👨", "member")

	err := svc.ResetPassword(member.ID)
	require.NoError(t, err)

	// 验证密码已变为默认密码
	updated, _ := svc.GetProfile(member.ID)
	assert.True(t, pkg.CheckPassword(updated.Password, pkg.DefaultPassword))
}

func TestMemberService_ChangePassword_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("changepwd", "oldpassword", "Change", "👨", "member")

	err := svc.ChangePassword(member.ID, "oldpassword", "newpassword")
	require.NoError(t, err)

	// 验证新密码生效
	updated, _ := svc.GetProfile(member.ID)
	assert.True(t, pkg.CheckPassword(updated.Password, "newpassword"))
}

func TestMemberService_ChangePassword_WrongOldPassword(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("wrongold", "correctold", "WrongOld", "👨", "member")

	err := svc.ChangePassword(member.ID, "wrongold", "newpassword")
	assert.ErrorIs(t, err, ErrIncorrectPassword)
}

func TestMemberService_List(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	svc.Create("user1", "password1", "User1", "", "member")
	svc.Create("user2", "password1", "User2", "", "admin")

	list, err := svc.List()
	require.NoError(t, err)
	assert.Len(t, list, 2)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/service/ -run "TestMember" -v
```

Expected: 16 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/member_test.go
git commit -m "test: add member service CRUD and validation tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 2.3: LedgerService 测试

**Files:**
- Create: `backend/internal/service/ledger_test.go`

- [ ] **Step 1: 编写记账服务测试**

`backend/internal/service/ledger_test.go`:
```go
package service

import (
	"testing"
	"time"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLedgerTest() (*LedgerService, func()) {
	db := testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewLedgerService()
	return svc, func() { testutil.TeardownTestDB() }
}

func seedLedgerFixtures(db interface{}) (member1, member2 model.Member, cat model.Category) {
	member1 = testutil.CreateTestMember(pkg.DB, "user1", "User1", "member")
	member2 = testutil.CreateTestMember(pkg.DB, "user2", "User2", "member")
	cat = testutil.CreateTestCategory(pkg.DB, "expense", "餐饮", "🍱", 1)
	return
}

func TestLedgerService_Create_Success(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, m2, cat := seedLedgerFixtures(nil)

	result, err := svc.Create(3550, "午餐", cat.ID, []uint{m1.ID, m2.ID}, time.Now(), m1.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3550), result.Amount)
	assert.Equal(t, "午餐", result.Note)
	assert.Equal(t, m1.ID, result.CreatorID)
	assert.Len(t, result.Members, 2)
}

func TestLedgerService_Create_ZeroAmount(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures(nil)

	_, err := svc.Create(0, "免费", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestLedgerService_Create_NegativeAmount(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures(nil)

	_, err := svc.Create(-100, "负数", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestLedgerService_Create_NoMembers(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	_, _, cat := seedLedgerFixtures(nil)

	_, err := svc.Create(1000, "无成员", cat.ID, []uint{}, time.Now(), 1)
	assert.ErrorIs(t, err, ErrNoMembers)
}

func TestLedgerService_Create_CategoryNotFound(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, _ := seedLedgerFixtures(nil)

	_, err := svc.Create(1000, "不存在的分类", 99999, []uint{m1.ID}, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrLedgerCategoryNotFound)
}

func TestLedgerService_Update_ByCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures(nil)
	created, _ := svc.Create(2000, "原始", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	newAmount := int64(3000)
	newNote := "修改后"
	updated, err := svc.Update(created.ID, &newAmount, &newNote, nil, nil, nil, m1.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, int64(3000), updated.Amount)
	assert.Equal(t, "修改后", updated.Note)
}

func TestLedgerService_Update_ByNonCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, m2, cat := seedLedgerFixtures(nil)
	created, _ := svc.Create(2000, "原始", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	newAmount := int64(3000)
	_, err := svc.Update(created.ID, &newAmount, nil, nil, nil, nil, m2.ID, "member")
	assert.ErrorIs(t, err, ErrLedgerPermissionDenied)
}

func TestLedgerService_Update_ByAdmin(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures(nil)
	admin := testutil.CreateTestMember(pkg.DB, "admin", "Admin", "admin")
	created, _ := svc.Create(2000, "原始", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	newAmount := int64(5000)
	updated, err := svc.Update(created.ID, &newAmount, nil, nil, nil, nil, admin.ID, "admin")
	require.NoError(t, err)
	assert.Equal(t, int64(5000), updated.Amount)
}

func TestLedgerService_Delete_ByCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures(nil)
	created, _ := svc.Create(1000, "待删除", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	err := svc.Delete(created.ID, m1.ID, "member")
	require.NoError(t, err)

	// 软删除后 FindByID 应返回 not found
	_, err = svc.FindByID(created.ID)
	assert.ErrorIs(t, err, ErrLedgerNotFound)
}

func TestLedgerService_Delete_ByNonCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, m2, cat := seedLedgerFixtures(nil)
	created, _ := svc.Create(1000, "待删除", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	err := svc.Delete(created.ID, m2.ID, "member")
	assert.ErrorIs(t, err, ErrLedgerPermissionDenied)
}

func TestLedgerService_Delete_ByAdmin(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures(nil)
	admin := testutil.CreateTestMember(pkg.DB, "admin2", "Admin2", "admin")
	created, _ := svc.Create(1000, "待删除", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	err := svc.Delete(created.ID, admin.ID, "admin")
	require.NoError(t, err)
}

func TestLedgerService_List_ByMonth(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures(nil)

	// 创建两条不同月份的记录
	svc.Create(1000, "5月", cat.ID, []uint{m1.ID}, time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC), m1.ID)
	svc.Create(2000, "6月", cat.ID, []uint{m1.ID}, time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), m1.ID)

	result, err := svc.List(repository.LedgerFilter{Month: "2026-05", Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/service/ -run "TestLedger" -v
```

Expected: 11 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/ledger_test.go
git commit -m "test: add ledger service CRUD and permission tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 2.4: TodoService 测试

**Files:**
- Create: `backend/internal/service/todo_test.go`

- [ ] **Step 1: 编写待办服务测试**

`backend/internal/service/todo_test.go`:
```go
package service

import (
	"testing"
	"time"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTodoTest() (*TodoService, func()) {
	db := testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewTodoService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestTodoService_Create_Success(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator", "Creator", "member")

	todo, err := svc.Create("买菜", "去超市买菜", "important", nil, nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "买菜", todo.Title)
	assert.Equal(t, "important", todo.Priority)
	assert.Equal(t, "pending", todo.Status)
	assert.Equal(t, creator.ID, todo.CreatorID)
}

func TestTodoService_Create_TitleRequired(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator2", "Creator2", "member")

	_, err := svc.Create("", "描述", "normal", nil, nil, creator.ID)
	assert.ErrorIs(t, err, ErrTodoTitleRequired)
}

func TestTodoService_Create_InvalidPriority(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator3", "Creator3", "member")

	_, err := svc.Create("测试", "", "super-urgent", nil, nil, creator.ID)
	assert.ErrorIs(t, err, ErrTodoInvalidPriority)
}

func TestTodoService_Create_DefaultPriority(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator4", "Creator4", "member")

	todo, err := svc.Create("默认优先级", "", "", nil, nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "normal", todo.Priority)
}

func TestTodoService_Create_WithAssignee(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator5", "Creator5", "member")
	assignee := testutil.CreateTestMember(pkg.DB, "assignee", "Assignee", "member")
	assigneeID := assignee.ID

	todo, err := svc.Create("指派任务", "分配给某人", "urgent", &assigneeID, nil, creator.ID)
	require.NoError(t, err)
	require.NotNil(t, todo.Assignee)
	assert.Equal(t, assignee.ID, todo.Assignee.ID)
}

func TestTodoService_Toggle_CompleteAndUncomplete(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator6", "Creator6", "member")

	todo, _ := svc.Create("完成测试", "", "normal", nil, nil, creator.ID)

	// 完成
	completed, err := svc.Toggle(todo.ID, creator.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, "completed", completed.Status)
	assert.NotNil(t, completed.CompletedAt)

	// 恢复
	uncompleted, err := svc.Toggle(todo.ID, creator.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, "pending", uncompleted.Status)
	assert.Nil(t, uncompleted.CompletedAt)
}

func TestTodoService_Toggle_PermissionDenied(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator7", "Creator7", "member")
	other := testutil.CreateTestMember(pkg.DB, "other7", "Other7", "member")

	todo, _ := svc.Create("权限测试", "", "normal", nil, nil, creator.ID)

	_, err := svc.Toggle(todo.ID, other.ID, "member")
	assert.ErrorIs(t, err, ErrTodoPermissionDenied)
}

func TestTodoService_Claim_Success(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator8", "Creator8", "member")
	claimer := testutil.CreateTestMember(pkg.DB, "claimer", "Claimer", "member")

	todo, _ := svc.Create("认领测试", "", "normal", nil, nil, creator.ID)

	claimed, err := svc.Claim(todo.ID, claimer.ID)
	require.NoError(t, err)
	require.NotNil(t, claimed.Assignee)
	assert.Equal(t, claimer.ID, claimed.Assignee.ID)
}

func TestTodoService_Claim_AlreadyAssigned(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator9", "Creator9", "member")
	assignee := testutil.CreateTestMember(pkg.DB, "assignee9", "Assignee9", "member")
	claimer := testutil.CreateTestMember(pkg.DB, "claimer9", "Claimer9", "member")
	assigneeID := assignee.ID

	todo, _ := svc.Create("已指派", "", "normal", &assigneeID, nil, creator.ID)

	_, err := svc.Claim(todo.ID, claimer.ID)
	assert.ErrorIs(t, err, ErrTodoAlreadyAssigned)
}

func TestTodoService_Delete_ByCreator(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator10", "Creator10", "member")

	todo, _ := svc.Create("删除测试", "", "normal", nil, nil, creator.ID)

	err := svc.Delete(todo.ID, creator.ID, "member")
	require.NoError(t, err)

	_, err = svc.FindByID(todo.ID)
	assert.ErrorIs(t, err, ErrTodoNotFound)
}

func TestTodoService_Delete_ByNonCreator(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator11", "Creator11", "member")
	other := testutil.CreateTestMember(pkg.DB, "other11", "Other11", "member")

	todo, _ := svc.Create("删除权限测试", "", "normal", nil, nil, creator.ID)

	err := svc.Delete(todo.ID, other.ID, "member")
	assert.ErrorIs(t, err, ErrTodoPermissionDenied)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/service/ -run "TestTodo" -v
```

Expected: 11 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/todo_test.go
git commit -m "test: add todo service CRUD, toggle, and claim tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 2.5: WishService 测试

**Files:**
- Create: `backend/internal/service/wish_test.go`

- [ ] **Step 1: 编写愿望服务测试**

`backend/internal/service/wish_test.go`:
```go
package service

import (
	"testing"

	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWishTest() (*WishService, func()) {
	db := testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewWishService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestWishService_Create_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher", "Wisher", "member")

	wish, err := svc.Create("想买iPad", "画画用", "item", "important", nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "想买iPad", wish.Title)
	assert.Equal(t, "item", wish.Category)
	assert.Equal(t, "personal", wish.Type, "default type should be personal")
	assert.Equal(t, "pending", wish.Status)
	assert.Equal(t, creator.ID, wish.CreatorID)
}

func TestWishService_Create_TitleRequired(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher2", "Wisher2", "member")

	_, err := svc.Create("", "", "other", "normal", nil, creator.ID)
	assert.ErrorIs(t, err, ErrWishTitleRequired)
}

func TestWishService_Create_InvalidCategory(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher3", "Wisher3", "member")

	_, err := svc.Create("test", "", "invalid_category", "normal", nil, creator.ID)
	assert.ErrorIs(t, err, ErrWishInvalidCategory)
}

func TestWishService_Create_DefaultCategory(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher4", "Wisher4", "member")

	wish, err := svc.Create("test", "", "", "normal", nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "other", wish.Category)
}

func TestWishService_Create_ValidCategories(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher5", "Wisher5", "member")

	for _, cat := range []string{"item", "travel", "experience", "other"} {
		wish, err := svc.Create("test "+cat, "", cat, "normal", nil, creator.ID)
		require.NoError(t, err)
		assert.Equal(t, cat, wish.Category)
	}
}

func TestWishService_Promote_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "promoter", "Promoter", "member")

	wish, _ := svc.Create("提升测试", "", "other", "normal", nil, creator.ID)
	assert.Equal(t, "personal", wish.Type)

	promoted, err := svc.Promote(wish.ID, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "family", promoted.Type)
}

func TestWishService_Promote_NotCreator(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "promoter2", "Promoter2", "member")
	other := testutil.CreateTestMember(pkg.DB, "other", "Other", "member")

	wish, _ := svc.Create("不能提升别人的愿望", "", "other", "normal", nil, creator.ID)

	_, err := svc.Promote(wish.ID, other.ID)
	assert.ErrorIs(t, err, ErrWishPermissionDenied)
}

func TestWishService_Vote_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "voter_creator", "VCreator", "member")
	voter := testutil.CreateTestMember(pkg.DB, "voter", "Voter", "member")

	wish, _ := svc.Create("投票测试", "", "other", "normal", nil, creator.ID)

	result, err := svc.Vote(wish.ID, voter.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, result.VoteCount)
}

func TestWishService_Vote_Duplicate(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vc2", "VC2", "member")
	voter := testutil.CreateTestMember(pkg.DB, "dvoter", "DVoter", "member")

	wish, _ := svc.Create("重复投票", "", "other", "normal", nil, creator.ID)

	svc.Vote(wish.ID, voter.ID)
	_, err := svc.Vote(wish.ID, voter.ID)
	assert.ErrorIs(t, err, ErrWishAlreadyVoted)
}

func TestWishService_Unvote_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vc3", "VC3", "member")
	voter := testutil.CreateTestMember(pkg.DB, "unvoter", "Unvoter", "member")

	wish, _ := svc.Create("取消投票", "", "other", "normal", nil, creator.ID)
	svc.Vote(wish.ID, voter.ID)

	result, err := svc.Unvote(wish.ID, voter.ID)
	require.NoError(t, err)
	assert.Equal(t, 0, result.VoteCount)
}

func TestWishService_UpdateStatus_AdminAnyStatus(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc3", "WC3", "member")
	admin := testutil.CreateTestMember(pkg.DB, "admin_wish", "AdminW", "admin")

	wish, _ := svc.Create("状态测试", "", "other", "normal", nil, creator.ID)

	for _, status := range []string{"pending", "agreed", "achieved", "abandoned"} {
		result, err := svc.UpdateStatus(wish.ID, status, admin.ID, "admin")
		require.NoError(t, err, "admin should be able to set status=%s", status)
		assert.Equal(t, status, result.Status)
	}
}

func TestWishService_UpdateStatus_CreatorOnlyAbandon(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc4", "WC4", "member")

	wish, _ := svc.Create("放弃测试", "", "other", "normal", nil, creator.ID)

	// 创建者可以标记放弃
	result, err := svc.UpdateStatus(wish.ID, "abandoned", creator.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, "abandoned", result.Status)
}

func TestWishService_UpdateStatus_CreatorCannotSetAgreed(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc5", "WC5", "member")

	wish, _ := svc.Create("不能改已同意", "", "other", "normal", nil, creator.ID)

	_, err := svc.UpdateStatus(wish.ID, "agreed", creator.ID, "member")
	assert.ErrorIs(t, err, ErrWishPermissionDenied)
}

func TestWishService_UpdateStatus_InvalidStatus(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc6", "WC6", "member")
	admin := testutil.CreateTestMember(pkg.DB, "admin6", "Admin6", "admin")

	wish, _ := svc.Create("无效状态", "", "other", "normal", nil, creator.ID)

	_, err := svc.UpdateStatus(wish.ID, "deleted", admin.ID, "admin")
	assert.ErrorIs(t, err, ErrWishInvalidStatus)
}

func TestWishService_Delete_ByCreator(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc7", "WC7", "member")

	wish, _ := svc.Create("删除愿望", "", "other", "normal", nil, creator.ID)

	err := svc.Delete(wish.ID, creator.ID, "member")
	require.NoError(t, err)

	_, err = svc.FindByID(wish.ID)
	assert.ErrorIs(t, err, ErrWishNotFound)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/service/ -run "TestWish" -v
```

Expected: 15 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/wish_test.go
git commit -m "test: add wish service CRUD, vote, promote, and status tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 2.6: ForumService 测试

**Files:**
- Create: `backend/internal/service/forum_test.go`

- [ ] **Step 1: 编写论坛服务测试**

`backend/internal/service/forum_test.go`:
```go
package service

import (
	"testing"
	"time"

	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupForumTest() (*ForumService, func()) {
	db := testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewForumService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestForumService_CreatePost_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "poster", "Poster", "member")

	post, err := svc.CreatePost("Hello world!", creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "Hello world!", post.Content)
	assert.Equal(t, creator.ID, post.CreatorID)
}

func TestForumService_CreatePost_EmptyContent(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "poster2", "Poster2", "member")

	_, err := svc.CreatePost("", creator.ID)
	assert.ErrorIs(t, err, ErrForumContentRequired)
}

func TestForumService_DeletePost_ByCreator(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "poster3", "Poster3", "member")

	post, _ := svc.CreatePost("Delete me", creator.ID)

	err := svc.DeletePost(post.ID, creator.ID, "member")
	require.NoError(t, err)

	// 验证软删除 — 查询应返回 RecordNotFound
	var deletedPost model.Post
	err = pkg.DB.First(&deletedPost, post.ID).Error
	assert.Error(t, err)
}

func TestForumService_DeletePost_ByNonCreator(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "poster4", "Poster4", "member")
	other := testutil.CreateTestMember(pkg.DB, "other4", "Other4", "member")

	post, _ := svc.CreatePost("Can't delete", creator.ID)

	err := svc.DeletePost(post.ID, other.ID, "member")
	assert.ErrorIs(t, err, ErrForumPermissionDenied)
}

func TestForumService_DeletePost_ByAdmin(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "poster5", "Poster5", "member")
	admin := testutil.CreateTestMember(pkg.DB, "admin5", "Admin5", "admin")

	post, _ := svc.CreatePost("Admin can delete", creator.ID)

	err := svc.DeletePost(post.ID, admin.ID, "admin")
	require.NoError(t, err)
}

func TestForumService_CreateTopic_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "topicer", "Topicer", "member")

	topic, err := svc.CreateTopic("讨论标题", "讨论内容", nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "讨论标题", topic.Title)
	assert.Equal(t, "讨论内容", topic.Content)
	assert.Equal(t, creator.ID, topic.CreatorID)
}

func TestForumService_CreateTopic_TitleRequired(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "topicer2", "Topicer2", "member")

	_, err := svc.CreateTopic("", "内容", nil, creator.ID)
	assert.ErrorIs(t, err, ErrForumTitleRequired)
}

func TestForumService_TogglePin(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "pinner_c", "PinnerC", "member")
	admin := testutil.CreateTestMember(pkg.DB, "pinner_a", "PinnerA", "admin")

	topic, _ := svc.CreateTopic("可置顶", "内容", nil, creator.ID)
	assert.False(t, topic.IsPinned)

	pinned, err := svc.TogglePin(topic.ID, "admin")
	require.NoError(t, err)
	assert.True(t, pinned.IsPinned)

	unpinned, err := svc.TogglePin(topic.ID, "admin")
	require.NoError(t, err)
	assert.False(t, unpinned.IsPinned)
}

func TestForumService_TogglePin_NotAdmin(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "member_pin", "MemberPin", "member")

	topic, _ := svc.CreateTopic("不能置顶", "内容", nil, creator.ID)

	_, err := svc.TogglePin(topic.ID, "member")
	assert.ErrorIs(t, err, ErrForumPermissionDenied)
}

func TestForumService_CreateComment_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "commenter", "Commenter", "member")
	post, _ := svc.CreatePost("文章", creator.ID)

	comment, err := svc.CreateComment("post", post.ID, nil, "好文章", creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "好文章", comment.Content)
	assert.Equal(t, "post", comment.TargetType)
}

func TestForumService_CreateComment_ReplyToLevel1(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "c1", "C1", "member")
	replyer := testutil.CreateTestMember(pkg.DB, "c2", "C2", "member")
	post, _ := svc.CreatePost("文章", creator.ID)
	level1, _ := svc.CreateComment("post", post.ID, nil, "一级评论", creator.ID)

	level2, err := svc.CreateComment("post", post.ID, &level1.ID, "回复你", replyer.ID)
	require.NoError(t, err)
	assert.NotNil(t, level2.ParentID)
	assert.Equal(t, level1.ID, *level2.ParentID)
}

func TestForumService_CreateComment_NestingTooDeep(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "nc1", "NC1", "member")
	replyer := testutil.CreateTestMember(pkg.DB, "nc2", "NC2", "member")
	third := testutil.CreateTestMember(pkg.DB, "nc3", "NC3", "member")
	post, _ := svc.CreatePost("嵌套测试", creator.ID)
	level1, _ := svc.CreateComment("post", post.ID, nil, "一级", creator.ID)
	level2, _ := svc.CreateComment("post", post.ID, &level1.ID, "二级", replyer.ID)

	// 尝试创建三级（回复二级）应失败
	_, err := svc.CreateComment("post", post.ID, &level2.ID, "三级不允许", third.ID)
	assert.ErrorIs(t, err, ErrForumNestingTooDeep)
}

func TestForumService_DeleteComment_SyncDeleteReplies(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "dc1", "DC1", "member")
	replyer := testutil.CreateTestMember(pkg.DB, "dc2", "DC2", "member")
	post, _ := svc.CreatePost("文章", creator.ID)
	level1, _ := svc.CreateComment("post", post.ID, nil, "一级", creator.ID)
	level2, _ := svc.CreateComment("post", post.ID, &level1.ID, "二级", replyer.ID)

	// 删除一级评论
	err := svc.DeleteComment(level1.ID, creator.ID, "member")
	require.NoError(t, err)

	// 二级评论应被同步软删除
	_, err = svc.repo.FindCommentWithCreator(level2.ID)
	assert.Error(t, err)
}

func TestForumService_ToggleLike(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "liker_c", "LikerC", "member")
	liker := testutil.CreateTestMember(pkg.DB, "liker", "Liker", "member")
	post, _ := svc.CreatePost("点赞测试", creator.ID)

	// 点赞
	liked, err := svc.ToggleLike("post", post.ID, liker.ID)
	require.NoError(t, err)
	assert.True(t, liked)

	// 取消点赞
	unliked, err := svc.ToggleLike("post", post.ID, liker.ID)
	require.NoError(t, err)
	assert.False(t, unliked)
}

func TestForumService_CreateVote_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vote_c", "VoteC", "member")

	vote, err := svc.CreateVote("周末去哪", []string{"公园", "海边", "爬山"}, false, nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "周末去哪", vote.Title)
	assert.Len(t, vote.Options, 3)
}

func TestForumService_CreateVote_TitleRequired(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vote_c2", "VoteC2", "member")

	_, err := svc.CreateVote("", []string{"A", "B"}, false, nil, creator.ID)
	assert.ErrorIs(t, err, ErrForumTitleRequired)
}

func TestForumService_CreateVote_MinOptions(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vote_c3", "VoteC3", "member")

	_, err := svc.CreateVote("只有一项", []string{"A"}, false, nil, creator.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 2 options")
}

func TestForumService_Vote_CastAndDuplicate(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vc_c", "VCC", "member")
	voter := testutil.CreateTestMember(pkg.DB, "vc_v", "VCV", "member")

	vote, _ := svc.CreateVote("选哪个", []string{"A", "B"}, false, nil, creator.ID)

	result, err := svc.Vote(vote.ID, vote.Options[0].ID, voter.ID)
	require.NoError(t, err)
	assert.Equal(t, 1, result.TotalVoters)

	// 重复投票
	_, err = svc.Vote(vote.ID, vote.Options[1].ID, voter.ID)
	assert.ErrorIs(t, err, ErrForumAlreadyVoted)
}

func TestForumService_DeleteVote_BeforeDeadline(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "dv_c", "DVC", "member")

	vote, _ := svc.CreateVote("可删除", []string{"A", "B"}, false, nil, creator.ID)

	err := svc.DeleteVote(vote.ID, creator.ID, "member")
	require.NoError(t, err)
}

func TestForumService_DeleteVote_AfterDeadline(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "dd_c", "DDC", "member")

	pastDeadline := time.Now().Add(-1 * time.Hour)
	vote, _ := svc.CreateVote("已过期", []string{"A", "B"}, false, &pastDeadline, creator.ID)

	err := svc.DeleteVote(vote.ID, creator.ID, "member")
	assert.ErrorIs(t, err, ErrForumVoteDeadlinePassed)
}

func TestForumService_DeleteTag_InUse(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "tag_c", "TagC", "member")

	tag, _ := svc.CreateTag("育儿")
	svc.CreateTopic("育儿话题", "内容", &tag.ID, creator.ID)

	err := svc.DeleteTag(tag.ID)
	assert.ErrorIs(t, err, ErrForumTagInUse)
}

func TestForumService_DeleteTag_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()

	tag, _ := svc.CreateTag("无引用标签")

	err := svc.DeleteTag(tag.ID)
	require.NoError(t, err)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/service/ -run "TestForum" -v
```

Expected: 24 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/service/forum_test.go
git commit -m "test: add forum service posts, topics, comments, likes, votes, and tags tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

---

## Phase 3: Handler 层集成测试

### Task 3.1: Auth Handler 测试

**Files:**
- Create: `backend/internal/handler/auth_test.go`

- [ ] **Step 1: 编写认证 API 集成测试**

`backend/internal/handler/auth_test.go`:
```go
package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() { gin.SetMode(gin.TestMode) }

func setupRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	authH := NewAuthHandler()
	r.GET("/api/init/check", authH.InitCheck)
	r.POST("/api/auth/login", authH.Login)
	return r
}

func TestHandler_InitCheck_NeedInit(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/init/check", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
	data := resp["data"].(map[string]interface{})
	assert.True(t, data["need_init"].(bool))
}

func TestHandler_InitCheck_NoNeedInit(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	// 创建已有成员
	hash, _ := pkg.HashPassword("pass")
	pkg.DB.Create(&model.Member{Username: "existing", Password: hash, Name: "Exist", Role: "admin", Status: "active"})

	r := setupRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/init/check", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	data := resp["data"].(map[string]interface{})
	assert.False(t, data["need_init"].(bool))
}

func TestHandler_Login_Success(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	hash, _ := pkg.HashPassword("correctpw")
	pkg.DB.Create(&model.Member{
		Username: "loginuser", Password: hash, Name: "LoginUser", Role: "member", Status: "active",
	})

	r := setupRouter()
	body := `{"username":"loginuser","password":"correctpw"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(0), resp["code"])
	data := resp["data"].(map[string]interface{})
	assert.NotEmpty(t, data["token"])
}

func TestHandler_Login_WrongPassword(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	hash, _ := pkg.HashPassword("correctpw")
	pkg.DB.Create(&model.Member{
		Username: "wrongpw", Password: hash, Name: "WrongPW", Role: "member", Status: "active",
	})

	r := setupRouter()
	body := `{"username":"wrongpw","password":"wrongpassword"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40101), resp["code"])
}

func TestHandler_Login_EmptyCredentials(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	r := setupRouter()
	body := `{"username":"","password":""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40001), resp["code"])
}

func TestHandler_Login_MissingFields(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	r := setupRouter()
	body := `{"username":"someone"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40001), resp["code"])
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler" -v
```

Expected: 6 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/auth_test.go
git commit -m "test: add auth handler API integration tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

---

## Phase 4: 前端组件测试

### Task 4.1: 前端测试基础设施

**Files:**
- Modify: `frontend/package.json`
- Create: `frontend/vitest.config.ts`
- Create: `frontend/src/test/setup.ts`

- [ ] **Step 1: 安装前端测试依赖**

```bash
cd frontend && npm install -D vitest @vue/test-utils happy-dom jsdom
```

- [ ] **Step 2: 创建 vitest 配置**

`frontend/vitest.config.ts`:
```typescript
import { defineConfig } from 'vitest/config'
import vue from '@vitejs/plugin-vue'

export default defineConfig({
  plugins: [vue()],
  test: {
    environment: 'jsdom',
    globals: true,
    setupFiles: ['./src/test/setup.ts'],
  },
})
```

- [ ] **Step 3: 创建测试 setup 文件**

`frontend/src/test/setup.ts`:
```typescript
import { config } from '@vue/test-utils'

// 全局 mock Ant Design Vue 的某些组件
config.global.stubs = {
  'a-button': true,
  'a-input': true,
  'a-form': true,
  'a-form-item': true,
  'a-modal': true,
  'a-tag': true,
  'a-avatar': true,
  'a-card': true,
  'a-empty': true,
  'a-spin': true,
  'a-select': true,
  'a-date-picker': true,
  'a-input-number': true,
  'a-checkbox': true,
  'a-checkbox-group': true,
  'a-radio-group': true,
  'router-link': true,
  'router-view': true,
}
```

- [ ] **Step 4: 在 package.json 添加测试脚本**

修改 `frontend/package.json` 的 `scripts`:
```json
{
  "scripts": {
    "dev": "vite",
    "build": "vue-tsc && vite build",
    "preview": "vite preview",
    "test": "vitest run",
    "test:watch": "vitest"
  }
}
```

- [ ] **Step 5: 验证测试框架可用**

```bash
cd frontend && npx vitest run --reporter=verbose
```

Expected: "No test files found" (正常，尚未创建测试文件)

- [ ] **Step 6: Commit**

```bash
git add frontend/package.json frontend/vitest.config.ts frontend/src/test/setup.ts
git commit -m "chore: add vitest and @vue/test-utils for frontend testing"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 4.2: 前端工具函数测试

**Files:**
- Create: `frontend/src/utils/__tests__/amount.test.ts`
- Create: `frontend/src/utils/amount.ts` (if not exists)

- [ ] **Step 1: 编写金额格式化测试**

如果 `frontend/src/utils/amount.ts` 不存在，先创建：

`frontend/src/utils/amount.ts`:
```typescript
/** 分转元格式化 */
export function formatAmount(cents: number): string {
  const yuan = cents / 100
  return `¥${yuan.toFixed(2)}`
}

/** 收入/支出带符号格式化 */
export function formatLedgerAmount(cents: number, type: 'income' | 'expense'): string {
  const yuan = cents / 100
  return type === 'income' ? `+¥${yuan.toFixed(2)}` : `-¥${yuan.toFixed(2)}`
}

/** 元转分 */
export function yuanToCents(yuan: number): number {
  return Math.round(yuan * 100)
}
```

`frontend/src/utils/__tests__/amount.test.ts`:
```typescript
import { describe, it, expect } from 'vitest'
import { formatAmount, formatLedgerAmount, yuanToCents } from '../amount'

describe('formatAmount', () => {
  it('应该将分转换为元格式', () => {
    expect(formatAmount(3550)).toBe('¥35.50')
  })

  it('应该处理零', () => {
    expect(formatAmount(0)).toBe('¥0.00')
  })

  it('应该处理整数元', () => {
    expect(formatAmount(10000)).toBe('¥100.00')
  })

  it('应该处理大额金额', () => {
    expect(formatAmount(99999999)).toBe('¥999999.99')
  })
})

describe('formatLedgerAmount', () => {
  it('收入应该带 + 前缀', () => {
    expect(formatLedgerAmount(10000, 'income')).toBe('+¥100.00')
  })

  it('支出应该带 - 前缀', () => {
    expect(formatLedgerAmount(5000, 'expense')).toBe('-¥50.00')
  })
})

describe('yuanToCents', () => {
  it('应该将元转为分', () => {
    expect(yuanToCents(35.5)).toBe(3550)
  })

  it('应该处理整数', () => {
    expect(yuanToCents(100)).toBe(10000)
  })

  it('应该正确处理浮点精度', () => {
    expect(yuanToCents(0.01)).toBe(1)
    expect(yuanToCents(0.99)).toBe(99)
  })
})
```

- [ ] **Step 2: 运行测试**

```bash
cd frontend && npx vitest run src/utils/__tests__/amount.test.ts
```

Expected: 9 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add frontend/src/utils/amount.ts frontend/src/utils/__tests__/amount.test.ts
git commit -m "test: add amount formatting utility and tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

### Task 4.3: EmptyState 组件测试

**Files:**
- Create: `frontend/src/components/__tests__/EmptyState.test.ts`

- [ ] **Step 1: 先确认 EmptyState 组件存在**

```bash
ls frontend/src/components/EmptyState.vue
```

If it doesn't exist, create a minimal version:

`frontend/src/components/EmptyState.vue`:
```vue
<template>
  <div class="empty-state" :class="type">
    <div class="empty-icon">📭</div>
    <p class="empty-message">{{ message }}</p>
    <a-button v-if="actionText && type === 'no-data'" type="primary" @click="$emit('action')">
      {{ actionText }}
    </a-button>
    <a v-if="type === 'no-result'" @click="$emit('clear')">清除筛选</a>
  </div>
</template>

<script setup lang="ts">
defineProps<{
  message: string
  actionText?: string
  type?: 'no-data' | 'no-result'
}>()

defineEmits<{
  action: []
  clear: []
}>()
</script>
```

`frontend/src/components/__tests__/EmptyState.test.ts`:
```typescript
import { describe, it, expect } from 'vitest'
import { mount } from '@vue/test-utils'
import EmptyState from '../EmptyState.vue'

describe('EmptyState', () => {
  it('应该渲染消息文本', () => {
    const wrapper = mount(EmptyState, {
      props: { message: '暂无数据' },
    })
    expect(wrapper.text()).toContain('暂无数据')
  })

  it('no-data 类型应该显示按钮', () => {
    const wrapper = mount(EmptyState, {
      props: { message: '暂无记录', actionText: '记一笔', type: 'no-data' },
    })
    expect(wrapper.text()).toContain('记一笔')
  })

  it('no-data 类型无 actionText 时不显示按钮', () => {
    const wrapper = mount(EmptyState, {
      props: { message: '暂无数据', type: 'no-data' },
    })
    expect(wrapper.find('a-button-stub').exists()).toBe(false)
  })

  it('no-result 类型应该显示清除筛选链接', () => {
    const wrapper = mount(EmptyState, {
      props: { message: '没有找到', type: 'no-result' },
    })
    expect(wrapper.text()).toContain('清除筛选')
  })

  it('no-result 类型不应该显示创建按钮', () => {
    const wrapper = mount(EmptyState, {
      props: { message: '没有找到', type: 'no-result', actionText: '创建' },
    })
    // no-result 下不渲染 action 按钮
    expect(wrapper.find('a-button-stub').exists()).toBe(false)
  })

  it('点击清除筛选应该触发 clear 事件', async () => {
    const wrapper = mount(EmptyState, {
      props: { message: '没有找到', type: 'no-result' },
    })
    await wrapper.find('a').trigger('click')
    expect(wrapper.emitted('clear')).toBeTruthy()
  })
})
```

- [ ] **Step 2: 运行测试**

```bash
cd frontend && npx vitest run src/components/__tests__/EmptyState.test.ts
```

Expected: 6 条测试全部 PASS

- [ ] **Step 3: Commit**

```bash
git add frontend/src/components/EmptyState.vue frontend/src/components/__tests__/EmptyState.test.ts
git commit -m "test: add EmptyState component unit tests"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

---

## Phase 5: 集成测试运行与验证

### Task 5.1: 运行全部测试并生成报告

- [ ] **Step 1: 运行全部后端测试**

```bash
cd backend && go test ./... -v -count=1 2>&1 | tee test-output.log
```

Expected: 所有测试 PASS (约 97+ 条测试)

- [ ] **Step 2: 运行全部前端测试**

```bash
cd frontend && npx vitest run --reporter=verbose 2>&1 | tee test-output.log
```

Expected: 所有测试 PASS (约 15 条测试)

- [ ] **Step 3: 检查测试覆盖率（可选）**

```bash
cd backend && go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -20
```

- [ ] **Step 4: Commit 最终测试状态**

```bash
git add backend/coverage.out test-output.log
git commit -m "test: add test reports and coverage data

Total: 112+ tests (97 backend + 15 frontend), all passing"

Co-Authored-By: Claude Opus 4.7 <noreply@anthropic.com>
```

---

## 附录 A: 测试覆盖率矩阵

| 模块 | Service 测试 | Handler 测试 | 工具测试 | 前端测试 | 总计 |
|------|-------------|-------------|---------|---------|------|
| 密码工具 | — | — | 7 | — | 7 |
| JWT 工具 | — | — | 7 | — | 7 |
| 响应工具 | — | — | 6 | — | 6 |
| 认证 | 8 | 6 | — | — | 14 |
| 成员管理 | 16 | — | — | — | 16 |
| 记账本 | 11 | — | — | — | 11 |
| 待办管理 | 11 | — | — | — | 11 |
| 愿望清单 | 15 | — | — | — | 15 |
| 论坛 | 24 | — | — | — | 24 |
| 金额工具 | — | — | — | 9 | 9 |
| EmptyState | — | — | — | 6 | 6 |
| **总计** | **85** | **6** | **20** | **15** | **126** |

## 附录 B: 关键测试场景对照表

| PRD 验收标准 | 覆盖任务 |
|------------|---------|
| A-01 首次初始化 | Task 2.1 (InitCheck) |
| A-02 用户名+密码登录 | Task 2.1 (Login), Task 3.1 (Login handler) |
| A-04 成员密码复杂度 | Task 2.2 (Create_InvalidPassword) |
| 登录5次失败锁定15分钟 | Task 2.1 (LockAfter5Failures) |
| 禁用成员无法登录 | Task 2.1 (DisabledMember) |
| M-01 添加成员校验 | Task 2.2 (Create系列) |
| M-02 禁止删除最后管理员 | Task 2.2 (CannotRemoveLastAdmin) |
| M-03a 禁止禁用自己 | Task 2.2 (CannotDisableSelf) |
| M-03 有活动记录不能删除 | Task 2.2 (Delete_WithActivityRecords) |
| L-01 金额正整数校验 | Task 2.3 (ZeroAmount, NegativeAmount) |
| L-01 至少选一个成员 | Task 2.3 (NoMembers) |
| L-04 只能编辑自己记录 | Task 2.3 (Update_ByNonCreator) |
| L-05 管理员可编辑任意 | Task 2.3 (Update_ByAdmin) |
| T-01 待办优先级校验 | Task 2.4 (InvalidPriority) |
| T-03 认领已指派待办 | Task 2.4 (Claim_AlreadyAssigned) |
| W-02 提升为家庭愿望(单向) | Task 2.5 (Promote系列) |
| W-03 每人一票 | Task 2.5 (Vote_Duplicate) |
| W-04a 创建者仅可标记放弃 | Task 2.5 (CreatorCannotSetAgreed) |
| F-06 评论2层嵌套限制 | Task 2.6 (NestingTooDeep) |
| F-07 点赞切换 | Task 2.6 (ToggleLike) |
| F-04a 截止后不可删除投票 | Task 2.6 (DeleteVote_AfterDeadline) |
| BR-CM04 删一级评论同步删二级 | Task 2.6 (SyncDeleteReplies) |
| 金额格式(分转元) | Task 4.2 |
| 空状态组件 | Task 4.3 |

## 附录 C: Task 依赖关系

```
Phase 0: Task 0.1 → Task 0.2 (测试基础设施)
              ↓
Phase 1: Task 1.1, 1.2, 1.3 (工具测试，可并行)
              ↓
Phase 2: Task 2.1 → Task 2.2 → Task 2.3 → Task 2.4 → Task 2.5 → Task 2.6 (Service 测试，按模块依赖顺序)
              ↓
Phase 3: Task 3.1 (Handler 测试)
              ↓
Phase 4: Task 4.1 → Task 4.2, 4.3 (前端测试)
              ↓
Phase 5: Task 5.1 (最终验证)
```
