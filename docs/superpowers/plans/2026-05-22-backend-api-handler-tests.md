# 后端 API Handler 集成测试实施计划

> **For agentic workers:** REQUIRED SUB-SKILL: Use superpowers:subagent-driven-development (recommended) or superpowers:executing-plans to implement this plan task-by-task. Steps use checkbox (`- [ ]`) syntax for tracking.

**Goal:** 对所有 57 个后端 API 端点编写 HTTP 级别集成测试，确保接口满足 PRD 和技术文档的目标要求。

**Architecture:** 使用 Go 标准 `testing` + `httptest` + `testify/assert`。通过 `testutil.SetupTestDB()` 搭建 in-memory SQLite，复用 `routes.Register()` 注册完整路由（含中间件），每个测试用例独立初始化。辅助函数提供 token 生成和 JSON 响应解析。

**Tech Stack:** Go testing, httptest, testify/assert, testify/require, Gin test mode, SQLite :memory:

**现状分析：**
- 现有 89 个 service 层测试 + 6 个 handler 测试（仅 auth）+ 18 个 pkg 测试
- 缺失：57 个 API 端点中 52 个无 HTTP 级别测试
- 缺失：middleware 独立测试
- 缺失：dashboard service 测试、init service 测试

---

## Task 1: 测试基础设施 — 路由与辅助函数

**Files:**
- Create: `backend/internal/testutil/router.go`
- Create: `backend/internal/testutil/helpers.go`

- [ ] **Step 1: 创建全路由测试路由器**

`backend/internal/testutil/router.go`:
```go
package testutil

import (
	"warmisle/internal/middleware"
	"warmisle/internal/handler"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
)

func init() {
	gin.SetMode(gin.TestMode)
}

// SetupTestRouter 创建完整的 Gin 路由，注册所有端点和中间件
func SetupTestRouter() *gin.Engine {
	r := gin.New()
	r.Use(gin.Recovery())

	api := r.Group("/api")

	authH := handler.NewAuthHandler()
	initH := handler.NewInitHandler()
	memberH := handler.NewMemberHandler()
	categoryH := handler.NewCategoryHandler()
	ledgerH := handler.NewLedgerHandler()
	dashboardH := handler.NewDashboardHandler()
	todoH := handler.NewTodoHandler()
	wishH := handler.NewWishHandler()
	forumH := handler.NewForumHandler()
	tagH := handler.NewTagHandler()

	// Public
	api.GET("/init/check", authH.InitCheck)
	api.POST("/init/setup", initH.Setup)
	api.POST("/auth/login", authH.Login)

	authRequired := middleware.AuthRequired()
	adminRequired := middleware.AdminRequired()

	// Members
	api.GET("/members", authRequired, memberH.List)
	api.POST("/members", authRequired, adminRequired, memberH.Create)
	api.PUT("/members/:id", authRequired, adminRequired, memberH.Update)
	api.DELETE("/members/:id", authRequired, adminRequired, memberH.Delete)
	api.PUT("/members/:id/disable", authRequired, adminRequired, memberH.Disable)
	api.PUT("/members/:id/enable", authRequired, adminRequired, memberH.Enable)
	api.PUT("/members/:id/reset-pwd", authRequired, adminRequired, memberH.ResetPassword)

	// Profile
	api.GET("/profile", authRequired, memberH.GetProfile)
	api.PUT("/profile", authRequired, memberH.UpdateProfile)
	api.PUT("/profile/password", authRequired, memberH.ChangePassword)

	// Categories
	api.GET("/categories", authRequired, categoryH.List)
	api.POST("/categories", authRequired, adminRequired, categoryH.Create)
	api.PUT("/categories/:id", authRequired, adminRequired, categoryH.Update)
	api.DELETE("/categories/:id", authRequired, adminRequired, categoryH.Delete)

	// Ledger
	api.GET("/ledgers", authRequired, ledgerH.List)
	api.POST("/ledgers", authRequired, ledgerH.Create)
	api.GET("/ledgers/:id", authRequired, ledgerH.GetByID)
	api.PUT("/ledgers/:id", authRequired, ledgerH.Update)
	api.DELETE("/ledgers/:id", authRequired, ledgerH.Delete)

	// Dashboard
	api.GET("/dashboard/summary", authRequired, dashboardH.Summary)
	api.GET("/dashboard/expense-chart", authRequired, dashboardH.ExpenseChart)
	api.GET("/dashboard/upcoming-todos", authRequired, dashboardH.UpcomingTodos)
	api.GET("/dashboard/wish-trends", authRequired, dashboardH.WishTrends)
	api.GET("/dashboard/forum-hot", authRequired, dashboardH.ForumHot)

	// Todos
	api.GET("/todos", authRequired, todoH.List)
	api.POST("/todos", authRequired, todoH.Create)
	api.PUT("/todos/:id", authRequired, todoH.Update)
	api.DELETE("/todos/:id", authRequired, todoH.Delete)
	api.PUT("/todos/:id/toggle", authRequired, todoH.Toggle)
	api.PUT("/todos/:id/claim", authRequired, todoH.Claim)

	// Wishes
	api.GET("/wishes", authRequired, wishH.List)
	api.POST("/wishes", authRequired, wishH.Create)
	api.PUT("/wishes/:id", authRequired, wishH.Update)
	api.DELETE("/wishes/:id", authRequired, wishH.Delete)
	api.POST("/wishes/:id/promote", authRequired, wishH.Promote)
	api.PUT("/wishes/:id/status", authRequired, wishH.UpdateStatus)
	api.POST("/wishes/:id/vote", authRequired, wishH.Vote)
	api.DELETE("/wishes/:id/vote", authRequired, wishH.Unvote)

	// Forum
	api.GET("/feed", authRequired, forumH.Feed)
	api.POST("/posts", authRequired, forumH.CreatePost)
	api.PUT("/posts/:id", authRequired, forumH.UpdatePost)
	api.DELETE("/posts/:id", authRequired, forumH.DeletePost)
	api.POST("/topics", authRequired, forumH.CreateTopic)
	api.PUT("/topics/:id", authRequired, forumH.UpdateTopic)
	api.DELETE("/topics/:id", authRequired, forumH.DeleteTopic)
	api.PUT("/topics/:id/pin", authRequired, adminRequired, forumH.TogglePin)
	api.GET("/topics/:id", authRequired, forumH.GetTopic)
	api.GET("/comments", authRequired, forumH.ListComments)
	api.POST("/comments", authRequired, forumH.CreateComment)
	api.DELETE("/comments/:id", authRequired, forumH.DeleteComment)
	api.POST("/likes", authRequired, forumH.ToggleLike)
	api.POST("/votes", authRequired, forumH.CreateVote)
	api.DELETE("/votes/:id", authRequired, forumH.DeleteVote)
	api.POST("/votes/:id/vote", authRequired, forumH.Vote)
	api.GET("/votes/:id", authRequired, forumH.GetVote)

	// Tags
	api.GET("/tags", authRequired, tagH.List)
	api.POST("/tags", authRequired, adminRequired, tagH.Create)
	api.PUT("/tags/:id", authRequired, adminRequired, tagH.Update)
	api.DELETE("/tags/:id", authRequired, adminRequired, tagH.Delete)

	return r
}
```

- [ ] **Step 2: 创建 HTTP 测试辅助函数**

`backend/internal/testutil/helpers.go`:
```go
package testutil

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

// APIResponse 统一响应结构
type APIResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

// MakeRequest 发起 HTTP 请求并返回响应
func MakeRequest(r *gin.Engine, method, path string, body interface{}, token string) *httptest.ResponseRecorder {
	var reqBody *bytes.Buffer
	if body != nil {
		switch v := body.(type) {
		case string:
			reqBody = bytes.NewBufferString(v)
		default:
			b, _ := json.Marshal(body)
			reqBody = bytes.NewBuffer(b)
		}
	} else {
		reqBody = bytes.NewBuffer(nil)
	}

	req, _ := http.NewRequest(method, path, reqBody)
	req.Header.Set("Content-Type", "application/json")
	if token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ParseResponse 解析 JSON 响应
func ParseResponse(t *testing.T, w *httptest.ResponseRecorder) APIResponse {
	var resp APIResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err, "failed to parse response: %s", w.Body.String())
	return resp
}

// ParseDataMap 将 resp.Data 转为 map
func ParseDataMap(resp APIResponse) map[string]interface{} {
	if resp.Data == nil {
		return nil
	}
	return resp.Data.(map[string]interface{})
}

// ParseDataArray 将 resp.Data 转为 []interface{}
func ParseDataArray(resp APIResponse) []interface{} {
	if resp.Data == nil {
		return nil
	}
	return resp.Data.([]interface{})
}

// GenerateTestToken 为测试成员生成 JWT token
func GenerateTestToken(member model.Member) string {
	token, _ := pkg.GenerateToken(member.ID, member.Username, member.Role)
	return token
}

// SeedAdminAndMember 创建一个管理员和一个普通成员，返回 (admin, member, adminToken, memberToken)
func SeedAdminAndMember(t *testing.T) (model.Member, model.Member, string, string) {
	db := pkg.DB
	hash, err := pkg.HashPassword("testpass123")
	require.NoError(t, err)

	admin := model.Member{
		Username: "admin", Password: hash, Name: "管理员",
		Avatar: "👨", Role: "admin", Status: "active",
	}
	db.Create(&admin)

	member := model.Member{
		Username: "member", Password: hash, Name: "普通成员",
		Avatar: "👩", Role: "member", Status: "active",
	}
	db.Create(&member)

	return admin, member, GenerateTestToken(admin), GenerateTestToken(member)
}

// SeedDisabledMember 创建一个禁用成员
func SeedDisabledMember(t *testing.T) (model.Member, string) {
	db := pkg.DB
	hash, _ := pkg.HashPassword("testpass123")
	disabled := model.Member{
		Username: "disabled", Password: hash, Name: "禁用成员",
		Avatar: "👶", Role: "member", Status: "disabled",
	}
	db.Create(&disabled)
	return disabled, GenerateTestToken(disabled)
}

// SeedTestCategory 创建测试分类
func SeedTestCategory(ctype, name, icon string, sortOrder int) model.Category {
	cat := model.Category{Type: ctype, Name: name, Icon: icon, SortOrder: sortOrder}
	pkg.DB.Create(&cat)
	return cat
}

// AssertErrorResponse 断言错误响应
func AssertErrorResponse(t *testing.T, w *httptest.ResponseRecorder, expectedHTTPCode, expectedBizCode int) {
	require.Equal(t, expectedHTTPCode, w.Code, fmt.Sprintf("expected HTTP %d, got %d: %s", expectedHTTPCode, w.Code, w.Body.String()))
	resp := ParseResponse(t, w)
	require.Equal(t, expectedBizCode, resp.Code, fmt.Sprintf("expected biz code %d, got %d", expectedBizCode, resp.Code))
}

// AssertSuccessResponse 断言成功响应
func AssertSuccessResponse(t *testing.T, w *httptest.ResponseRecorder) APIResponse {
	require.Equal(t, 200, w.Code, fmt.Sprintf("expected HTTP 200, got %d: %s", w.Code, w.Body.String()))
	resp := ParseResponse(t, w)
	require.Equal(t, 0, resp.Code, fmt.Sprintf("expected code 0, got %d", resp.Code))
	return resp
}
```

- [ ] **Step 3: 验证编译通过**

```bash
cd backend && go build ./internal/testutil/
```

Expected: 编译成功

- [ ] **Step 4: Commit**

```bash
git add backend/internal/testutil/router.go backend/internal/testutil/helpers.go
git commit -m "test: add full router setup and HTTP test helpers for handler integration tests"
```

---

## Task 2: Init & Auth Handler 测试（扩展现有）

**Files:**
- Modify: `backend/internal/handler/auth_test.go`
- Create: `backend/internal/handler/init_test.go`

- [ ] **Step 1: 创建 Init handler 测试**

`backend/internal/handler/init_test.go`:
```go
package handler

import (
	"testing"

	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Setup_Success(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	initJWT()

	r := testutil.SetupTestRouter()

	body := `{"admin_name":"管理员","username":"admin","password":"admin123"}`
	w := testutil.MakeRequest(r, "POST", "/api/init/setup", body, "")
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.NotEmpty(t, data["token"])
	member := data["member"].(map[string]interface{})
	assert.Equal(t, "管理员", member["name"])
	assert.Equal(t, "admin", member["role"])
}

func TestHandler_Setup_AlreadyInitialized(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	initJWT()

	r := testutil.SetupTestRouter()

	// 第一次初始化
	body := `{"admin_name":"管理员","username":"admin","password":"admin123"}`
	testutil.MakeRequest(r, "POST", "/api/init/setup", body, "")

	// 第二次应失败
	w := testutil.MakeRequest(r, "POST", "/api/init/setup", body, "")
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Setup_MissingFields(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	initJWT()

	r := testutil.SetupTestRouter()

	body := `{"username":"admin"}`
	w := testutil.MakeRequest(r, "POST", "/api/init/setup", body, "")
	testutil.AssertErrorResponse(t, w, 400, 40001)
}
```

- [ ] **Step 2: 扩展 Auth handler 测试 — 添加禁用用户登录和 token 过期测试**

在 `auth_test.go` 中添加（保留现有测试，追加）：
```go
func TestHandler_Login_DisabledUser(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	initJWT()

	r := testutil.SetupTestRouter()

	// 先初始化系统
	initBody := `{"admin_name":"管理员","username":"admin","password":"admin123"}`
	testutil.MakeRequest(r, "POST", "/api/init/setup", initBody, "")

	// 创建禁用成员
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)
	disabled, _ := testutil.SeedDisabledMember(t)
	_ = disabled

	loginBody := `{"username":"disabled","password":"testpass123"}`
	w := testutil.MakeRequest(r, "POST", "/api/auth/login", loginBody, "")
	testutil.AssertErrorResponse(t, w, 401, 40101)
}

func TestHandler_ProtectedEndpoint_NoToken(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	initJWT()

	r := testutil.SetupTestRouter()

	w := testutil.MakeRequest(r, "GET", "/api/members", nil, "")
	testutil.AssertErrorResponse(t, w, 401, 40101)
}

func TestHandler_ProtectedEndpoint_InvalidToken(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	initJWT()

	r := testutil.SetupTestRouter()

	w := testutil.MakeRequest(r, "GET", "/api/members", nil, "invalid-token")
	testutil.AssertErrorResponse(t, w, 401, 40101)
}
```

- [ ] **Step 3: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Setup|TestHandler_Login_Disabled|TestHandler_Protected" -v
```

Expected: 所有新测试 PASS

- [ ] **Step 4: Commit**

```bash
git add backend/internal/handler/init_test.go backend/internal/handler/auth_test.go
git commit -m "test: add init handler tests and expand auth handler tests for disabled users and token validation"
```

---

## Task 3: Member & Profile Handler 测试

**Files:**
- Create: `backend/internal/handler/member_test.go`

- [ ] **Step 1: 编写 Member handler 测试**

`backend/internal/handler/member_test.go`:
```go
package handler

import (
	"fmt"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func initJWT() {
	pkg.InitJWT("test-secret-for-handler-tests")
}

func setupMemberTest() {
	testutil.SetupTestDB()
	initJWT()
}

func TestHandler_Member_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"username":"newuser","password":"pass123","name":"新成员","avatar":"👩"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "newuser", data["username"])
	assert.Equal(t, "新成员", data["name"])
}

func TestHandler_Member_Create_DuplicateUsername(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"username":"admin","password":"pass123","name":"重复"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, adminToken)
	testutil.AssertErrorResponse(t, w, 409, 40002)
}

func TestHandler_Member_Create_InvalidUsername(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"username":"ab","password":"pass123","name":"短用户名"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Member_Create_ForbiddenForMember(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"username":"newuser","password":"pass123","name":"新成员"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Member_List_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "GET", "/api/members", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataArray(resp)
	assert.Len(t, data, 2)
}

func TestHandler_Member_Update_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"name":"修改后名称","avatar":"🐶"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d", member.ID), body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "修改后名称", data["name"])
}

func TestHandler_Member_Update_CannotRemoveLastAdmin(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	admin, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"role":"member"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d", admin.ID), body, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40003)
}

func TestHandler_Member_Delete_ForbiddenForMember(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/members/%d", member.ID), nil, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Member_Delete_WithActivityRecords(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	// 给成员创建一条记账记录
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)
	pkg.DB.Create(&model.Ledger{Amount: 1000, CategoryID: cat.ID, CreatorID: member.ID})

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/members/%d", member.ID), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40004)
}

func TestHandler_Member_Disable_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/disable", member.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Member_Disable_Self(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	admin, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/disable", admin.ID), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40005)
}

func TestHandler_Member_Enable_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	// 先禁用
	testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/disable", member.ID), nil, adminToken)

	// 再启用
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/enable", member.ID), nil, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "active", data["status"])
}

func TestHandler_Member_ResetPassword(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/reset-pwd", member.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Profile_Get(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "GET", "/api/profile", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "member", data["username"])
}

func TestHandler_Profile_Update(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"name":"新名字","avatar":"🔥"}`
	w := testutil.MakeRequest(r, "PUT", "/api/profile", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "新名字", data["name"])
}

func TestHandler_Profile_ChangePassword_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"old_password":"testpass123","new_password":"newpass456"}`
	w := testutil.MakeRequest(r, "PUT", "/api/profile/password", body, memberToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Profile_ChangePassword_WrongOld(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"old_password":"wrongold","new_password":"newpass456"}`
	w := testutil.MakeRequest(r, "PUT", "/api/profile/password", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40010)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Member|TestHandler_Profile" -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/member_test.go
git commit -m "test: add member and profile handler HTTP integration tests"
```

---

## Task 4: Category Handler 测试

**Files:**
- Create: `backend/internal/handler/category_test.go`

- [ ] **Step 1: 编写 Category handler 测试**

`backend/internal/handler/category_test.go`:
```go
package handler

import (
	"fmt"
	"testing"

	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Category_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"type":"expense","name":"交通","icon":"🚗","sort_order":2}`
	w := testutil.MakeRequest(r, "POST", "/api/categories", body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "expense", data["type"])
	assert.Equal(t, "交通", data["name"])
}

func TestHandler_Category_Create_ForbiddenForMember(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"type":"expense","name":"交通","icon":"🚗"}`
	w := testutil.MakeRequest(r, "POST", "/api/categories", body, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Category_Create_DuplicateName(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"type":"expense","name":"餐饮","icon":"🍱"}`
	testutil.MakeRequest(r, "POST", "/api/categories", body, adminToken)

	w := testutil.MakeRequest(r, "POST", "/api/categories", body, adminToken)
	testutil.AssertErrorResponse(t, w, 409, 40002)
}

func TestHandler_Category_Create_InvalidType(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"type":"invalid","name":"测试","icon":"🔧"}`
	w := testutil.MakeRequest(r, "POST", "/api/categories", body, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Category_List_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)
	testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)
	testutil.SeedTestCategory("income", "工资", "💰", 1)

	w := testutil.MakeRequest(r, "GET", "/api/categories", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.Len(t, data, 2)
}

func TestHandler_Category_Update_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := `{"name":"美食","icon":"🍕"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/categories/%d", cat.ID), body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "美食", data["name"])
}

func TestHandler_Category_Delete_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "临时", "🔧", 99)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/categories/%d", cat.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Category_Delete_InUse(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	// 创建一条使用该分类的记账记录
	pkg.DB.Create(&model.Ledger{Amount: 1000, CategoryID: cat.ID, CreatorID: member.ID})

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/categories/%d", cat.ID), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}
```

需要在文件顶部添加 import `"warmisle/internal/model"` 和 `"warmisle/internal/pkg"`。

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Category" -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/category_test.go
git commit -m "test: add category handler HTTP integration tests"
```

---

## Task 5: Ledger Handler 测试

**Files:**
- Create: `backend/internal/handler/ledger_test.go`

- [ ] **Step 1: 编写 Ledger handler 测试**

`backend/internal/handler/ledger_test.go`:
```go
package handler

import (
	"fmt"
	"testing"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHandler_Ledger_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":35.50,"note":"午餐","category_id":%d,"member_ids":[%d],"occurred_at":"%s"}`,
		cat.ID, member.ID, time.Now().Format(time.RFC3339))
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, float64(3550), data["amount"])
	assert.Equal(t, "午餐", data["note"])
}

func TestHandler_Ledger_Create_ZeroAmount(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":0,"note":"零元","category_id":%d,"member_ids":[%d]}`, cat.ID, member.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_NegativeAmount(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":-10,"note":"负数","category_id":%d,"member_ids":[%d]}`, cat.ID, member.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_NoMembers(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":10,"note":"无成员","category_id":%d,"member_ids":[]}`, cat.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_CategoryNotFound(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	body := fmt.Sprintf(`{"amount":10,"note":"分类不存在","category_id":99999,"member_ids":[%d]}`, member.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_List_ByMonth(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	// 创建当月记录
	now := time.Now()
	pkg.DB.Create(&model.Ledger{Amount: 3000, Note: "当月", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: now})
	pkg.DB.Create(&model.LedgerMember{LedgerID: 1, MemberID: member.ID})

	month := now.Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["summary"])
}

func TestHandler_Ledger_Update_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	// 创建记录
	ledger := model.Ledger{Amount: 2000, Note: "原始", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: time.Now()}
	pkg.DB.Create(&ledger)
	pkg.DB.Create(&model.LedgerMember{LedgerID: ledger.ID, MemberID: member.ID})

	body := `{"amount":50.00,"note":"修改后"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/ledgers/%d", ledger.ID), body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "修改后", data["note"])
}

func TestHandler_Ledger_Update_ByNonCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	admin, member, _, _ := testutil.SeedAdminAndMember(t)
	_, _, member2Token := func() (model.Member, model.Member, string) {
		hash, _ := pkg.HashPassword("testpass123")
		m2 := model.Member{Username: "member2", Password: hash, Name: "成员2", Avatar: "👶", Role: "member", Status: "active"}
		pkg.DB.Create(&m2)
		return model.Member{}, m2, testutil.GenerateTestToken(m2)
	}()
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	ledger := model.Ledger{Amount: 2000, Note: "原始", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: time.Now()}
	pkg.DB.Create(&ledger)

	body := `{"amount":50.00}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/ledgers/%d", ledger.ID), body, member2Token)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Ledger_Delete_ByAdmin(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	ledger := model.Ledger{Amount: 1000, Note: "待删除", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: time.Now()}
	pkg.DB.Create(&ledger)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/ledgers/%d", ledger.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)

	// 验证软删除
	var count int64
	pkg.DB.Unscoped().Model(&model.Ledger{}).Where("id = ? AND deleted_at IS NOT NULL", ledger.ID).Count(&count)
	require.Equal(t, int64(1), count)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Ledger" -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/ledger_test.go
git commit -m "test: add ledger handler HTTP integration tests"
```

---

## Task 6: Todo Handler 测试

**Files:**
- Create: `backend/internal/handler/todo_test.go`

- [ ] **Step 1: 编写 Todo handler 测试**

`backend/internal/handler/todo_test.go`:
```go
package handler

import (
	"fmt"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Todo_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"买菜","description":"去超市","priority":"important"}`
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "买菜", data["title"])
	assert.Equal(t, "important", data["priority"])
	assert.Equal(t, "pending", data["status"])
}

func TestHandler_Todo_Create_EmptyTitle(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"","description":"空标题"}`
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Todo_Create_InvalidPriority(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"测试","priority":"super-urgent"}`
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Todo_Create_WithAssignee(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	body := fmt.Sprintf(`{"title":"指派任务","assignee_id":%d}`, member.ID)
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["assignee"])
}

func TestHandler_Todo_List_WithFilter(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建几条待办
	testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"待办1","priority":"urgent"}`, memberToken)
	testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"待办2","priority":"normal"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/todos", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["list"])
}

func TestHandler_Todo_Toggle_CompleteAndUncomplete(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建待办
	resp := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"完成测试"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	todoID := data["id"].(float64)

	// 完成
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d/toggle", int(todoID)), nil, memberToken)
	toggleResp := testutil.AssertSuccessResponse(t, w)
	toggleData := testutil.ParseDataMap(toggleResp)
	assert.Equal(t, "completed", toggleData["status"])

	// 恢复
	w = testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d/toggle", int(todoID)), nil, memberToken)
	toggleResp2 := testutil.AssertSuccessResponse(t, w)
	toggleData2 := testutil.ParseDataMap(toggleResp2)
	assert.Equal(t, "pending", toggleData2["status"])
}

func TestHandler_Todo_Claim_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建未指派的待办
	resp := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"认领测试"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	todoID := data["id"].(float64)

	// 创建第二个成员来认领
	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "claimer", Password: hash, Name: "认领者", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d/claim", int(todoID)), nil, m2Token)
	claimResp := testutil.AssertSuccessResponse(t, w)
	claimData := testutil.ParseDataMap(claimResp)
	assert.NotNil(t, claimData["assignee"])
}

func TestHandler_Todo_Update_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"原始标题"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	todoID := data["id"].(float64)

	body := `{"title":"修改后标题","priority":"urgent"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d", int(todoID)), body, memberToken)
	updateResp := testutil.AssertSuccessResponse(t, w)
	updateData := testutil.ParseDataMap(updateResp)
	assert.Equal(t, "修改后标题", updateData["title"])
}

func TestHandler_Todo_Delete_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"待删除"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	todoID := data["id"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/todos/%d", int(todoID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Todo" -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/todo_test.go
git commit -m "test: add todo handler HTTP integration tests"
```

---

## Task 7: Wish Handler 测试

**Files:**
- Create: `backend/internal/handler/wish_test.go`

- [ ] **Step 1: 编写 Wish handler 测试**

`backend/internal/handler/wish_test.go`:
```go
package handler

import (
	"fmt"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Wish_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"想买iPad","description":"画画用","category":"item","priority":"important"}`
	w := testutil.MakeRequest(r, "POST", "/api/wishes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "想买iPad", data["title"])
	assert.Equal(t, "personal", data["type"])
	assert.Equal(t, "pending", data["status"])
}

func TestHandler_Wish_Create_EmptyTitle(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"","description":"空标题"}`
	w := testutil.MakeRequest(r, "POST", "/api/wishes", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Wish_Create_InvalidCategory(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"测试","category":"invalid_cat"}`
	w := testutil.MakeRequest(r, "POST", "/api/wishes", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Wish_Promote_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建个人愿望
	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"提升测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)

	// 提升为家庭愿望
	w := testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/promote", int(wishID)), nil, memberToken)
	promoteResp := testutil.AssertSuccessResponse(t, w)
	promoteData := testutil.ParseDataMap(promoteResp)
	assert.Equal(t, "family", promoteData["type"])
}

func TestHandler_Wish_Promote_NotCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建愿望
	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"不能提升","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)

	// 另一个成员尝试提升
	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "other", Password: hash, Name: "其他", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	w := testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/promote", int(wishID)), nil, m2Token)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Wish_Vote_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建家庭愿望
	resp := testutil.MakeRequest(r, "POST", `/api/wishes`, `{"title":"投票测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)
	// 先提升
	testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/promote", int(wishID)), nil, memberToken)

	// 第二个成员投票
	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "voter", Password: hash, Name: "投票者", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	w := testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/vote", int(wishID)), nil, m2Token)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Wish_Vote_Duplicate(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"重复投票","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)
	testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/promote", int(wishID)), nil, memberToken)

	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "dvoter", Password: hash, Name: "重复投票者", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	// 第一次投票
	testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/vote", int(wishID)), nil, m2Token)
	// 第二次应失败
	w := testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/vote", int(wishID)), nil, m2Token)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Wish_Unvote_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"取消投票","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)
	testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/promote", int(wishID)), nil, memberToken)

	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "unvoter", Password: hash, Name: "取消投票者", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/vote", int(wishID)), nil, m2Token)
	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/wishes/%d/vote", int(wishID)), nil, m2Token)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Wish_UpdateStatus_AdminAnyStatus(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"状态测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)

	for _, status := range []string{"agreed", "achieved"} {
		body := fmt.Sprintf(`{"status":"%s"}`, status)
		w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/wishes/%d/status", int(wishID)), body, adminToken)
		statusResp := testutil.AssertSuccessResponse(t, w)
		statusData := testutil.ParseDataMap(statusResp)
		assert.Equal(t, status, statusData["status"])
	}
}

func TestHandler_Wish_UpdateStatus_CreatorOnlyAbandon(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"放弃测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)

	// 创建者可以标记放弃
	body := `{"status":"abandoned"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/wishes/%d/status", int(wishID)), body, memberToken)
	statusResp := testutil.AssertSuccessResponse(t, w)
	statusData := testutil.ParseDataMap(statusResp)
	assert.Equal(t, "abandoned", statusData["status"])
}

func TestHandler_Wish_Delete_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"删除测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["id"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/wishes/%d", int(wishID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Wish" -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/wish_test.go
git commit -m "test: add wish handler HTTP integration tests"
```

---

## Task 8: Dashboard Handler 测试

**Files:**
- Create: `backend/internal/handler/dashboard_test.go`

- [ ] **Step 1: 编写 Dashboard handler 测试**

`backend/internal/handler/dashboard_test.go`:
```go
package handler

import (
	"fmt"
	"testing"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Dashboard_Summary_Empty(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	month := time.Now().Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/dashboard/summary?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, float64(0), data["income"])
	assert.Equal(t, float64(0), data["expense"])
	assert.Equal(t, float64(0), data["balance"])
}

func TestHandler_Dashboard_Summary_WithData(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	incomeCat := testutil.SeedTestCategory("income", "工资", "💰", 1)
	expenseCat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	now := time.Now()
	pkg.DB.Create(&model.Ledger{Amount: 10000, Note: "工资", CategoryID: incomeCat.ID, CreatorID: member.ID, OccurredAt: now})
	pkg.DB.Create(&model.Ledger{Amount: 3000, Note: "午餐", CategoryID: expenseCat.ID, CreatorID: member.ID, OccurredAt: now})

	month := now.Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/dashboard/summary?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, float64(10000), data["income"])
	assert.Equal(t, float64(3000), data["expense"])
	assert.Equal(t, float64(7000), data["balance"])
}

func TestHandler_Dashboard_ExpenseChart(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)
	pkg.DB.Create(&model.Ledger{Amount: 5000, Note: "餐饮", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: time.Now()})

	month := time.Now().Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/dashboard/expense-chart?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}

func TestHandler_Dashboard_UpcomingTodos(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建待办
	testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"近期待办","priority":"urgent"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/dashboard/upcoming-todos", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}

func TestHandler_Dashboard_WishTrends(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"愿望动态","category":"other"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/dashboard/wish-trends", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}

func TestHandler_Dashboard_ForumHot(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"论坛热点测试"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/dashboard/forum-hot", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Dashboard" -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/dashboard_test.go
git commit -m "test: add dashboard handler HTTP integration tests"
```

---

## Task 9: Forum Handler 测试（Post + Topic + Comment + Like + Vote + Tag）

**Files:**
- Create: `backend/internal/handler/forum_test.go`

- [ ] **Step 1: 编写 Forum handler 测试**

`backend/internal/handler/forum_test.go`:
```go
package handler

import (
	"fmt"
	"testing"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

// === Posts ===

func TestHandler_Forum_CreatePost_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"content":"今天天气真好！"}`
	w := testutil.MakeRequest(r, "POST", "/api/posts", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "今天天气真好！", data["content"])
}

func TestHandler_Forum_CreatePost_EmptyContent(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"content":""}`
	w := testutil.MakeRequest(r, "POST", "/api/posts", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_DeletePost_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"待删除"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	postID := data["id"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/posts/%d", int(postID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Forum_DeletePost_ByAdmin(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"管理员删除"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	postID := data["id"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/posts/%d", int(postID)), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Forum_DeletePost_Forbidden(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"不能删"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	postID := data["id"].(float64)

	// 创建另一个成员
	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "other", Password: hash, Name: "其他", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/posts/%d", int(postID)), nil, m2Token)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

// === Topics ===

func TestHandler_Forum_CreateTopic_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"讨论话题","content":"话题内容"}`
	w := testutil.MakeRequest(r, "POST", "/api/topics", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "讨论话题", data["title"])
}

func TestHandler_Forum_CreateTopic_EmptyTitle(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"","content":"内容"}`
	w := testutil.MakeRequest(r, "POST", "/api/topics", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_CreateTopic_WithTag(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	// 创建标签
	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["id"].(float64)

	body := fmt.Sprintf(`{"title":"育儿话题","content":"讨论育儿","tag_id":%d}`, int(tagID))
	w := testutil.MakeRequest(r, "POST", "/api/topics", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["tag"])
}

func TestHandler_Forum_TogglePin_AdminOnly(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	// 创建话题
	topicResp := testutil.MakeRequest(r, "POST", "/api/topics", `{"title":"可置顶","content":"内容"}`, memberToken)
	topicData := testutil.ParseDataMap(testutil.ParseResponse(t, topicResp))
	topicID := topicData["id"].(float64)

	// 普通成员不能置顶
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/topics/%d/pin", int(topicID)), nil, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)

	// 管理员可以置顶
	w = testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/topics/%d/pin", int(topicID)), nil, adminToken)
	pinResp := testutil.AssertSuccessResponse(t, w)
	pinData := testutil.ParseDataMap(pinResp)
	assert.Equal(t, true, pinData["is_pinned"])
}

func TestHandler_Forum_GetTopic(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	topicResp := testutil.MakeRequest(r, "POST", "/api/topics", `{"title":"查看话题","content":"内容"}`, memberToken)
	topicData := testutil.ParseDataMap(testutil.ParseResponse(t, topicResp))
	topicID := topicData["id"].(float64)

	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/topics/%d", int(topicID)), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "查看话题", data["title"])
}

// === Comments ===

func TestHandler_Forum_CreateComment_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"评论目标"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["id"].(float64)

	body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"好文"}`, int(postID))
	w := testutil.MakeRequest(r, "POST", "/api/comments", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "好文", data["content"])
}

func TestHandler_Forum_CreateComment_ReplyToLevel1(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"回复测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["id"].(float64)

	// 一级评论
	c1Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"一级"}`, int(postID))
	c1Resp := testutil.MakeRequest(r, "POST", "/api/comments", c1Body, memberToken)
	c1Data := testutil.ParseDataMap(testutil.ParseResponse(t, c1Resp))
	c1ID := c1Data["id"].(float64)

	// 二级回复
	c2Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"二级回复"}`, int(postID), int(c1ID))
	w := testutil.MakeRequest(r, "POST", "/api/comments", c2Body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["parent_id"])
}

func TestHandler_Forum_CreateComment_NestingTooDeep(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"嵌套测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["id"].(float64)

	// 一级
	c1Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"一级"}`, int(postID))
	c1Resp := testutil.MakeRequest(r, "POST", "/api/comments", c1Body, memberToken)
	c1Data := testutil.ParseDataMap(testutil.ParseResponse(t, c1Resp))
	c1ID := c1Data["id"].(float64)

	// 二级
	c2Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"二级"}`, int(postID), int(c1ID))
	c2Resp := testutil.MakeRequest(r, "POST", "/api/comments", c2Body, memberToken)
	c2Data := testutil.ParseDataMap(testutil.ParseResponse(t, c2Resp))
	c2ID := c2Data["id"].(float64)

	// 三级应失败
	c3Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"三级不允许"}`, int(postID), int(c2ID))
	w := testutil.MakeRequest(r, "POST", "/api/comments", c3Body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_DeleteComment_SyncDeleteReplies(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"删除测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["id"].(float64)

	c1Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"一级"}`, int(postID))
	c1Resp := testutil.MakeRequest(r, "POST", "/api/comments", c1Body, memberToken)
	c1Data := testutil.ParseDataMap(testutil.ParseResponse(t, c1Resp))
	c1ID := c1Data["id"].(float64)

	c2Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"二级"}`, int(postID), int(c1ID))
	c2Resp := testutil.MakeRequest(r, "POST", "/api/comments", c2Body, memberToken)
	c2Data := testutil.ParseDataMap(testutil.ParseResponse(t, c2Resp))
	c2ID := c2Data["id"].(float64)

	// 删除一级
	testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/comments/%d", int(c1ID)), nil, memberToken)

	// 二级应被同步软删除
	var count int64
	pkg.DB.Unscoped().Model(&model.Comment{}).Where("id = ? AND deleted_at IS NOT NULL", int(c2ID)).Count(&count)
	assert.Equal(t, int64(1), count)
}

// === Likes ===

func TestHandler_Forum_ToggleLike(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"点赞测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["id"].(float64)

	// 点赞
	body := fmt.Sprintf(`{"target_type":"post","target_id":%d}`, int(postID))
	w := testutil.MakeRequest(r, "POST", "/api/likes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, true, data["liked"])

	// 取消点赞
	w = testutil.MakeRequest(r, "POST", "/api/likes", body, memberToken)
	resp = testutil.AssertSuccessResponse(t, w)
	data = testutil.ParseDataMap(resp)
	assert.Equal(t, false, data["liked"])
}

// === Votes ===

func TestHandler_Forum_CreateVote_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"周末去哪","options":["公园","海边","爬山"],"is_multi":false}`
	w := testutil.MakeRequest(r, "POST", "/api/votes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "周末去哪", data["title"])
}

func TestHandler_Forum_CreateVote_WithDeadline(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	deadline := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	body := fmt.Sprintf(`{"title":"限时投票","options":["A","B"],"deadline":"%s"}`, deadline)
	w := testutil.MakeRequest(r, "POST", "/api/votes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["deadline"])
}

func TestHandler_Forum_Vote_CastAndDuplicate(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	voteResp := testutil.MakeRequest(r, "POST", "/api/votes", `{"title":"投票","options":["A","B"]}`, memberToken)
	voteData := testutil.ParseDataMap(testutil.ParseResponse(t, voteResp))
	voteID := voteData["id"].(float64)
	options := voteData["options"].([]interface{})
	optionID := options[0].(map[string]interface{})["id"].(float64)

	// 投票
	castBody := fmt.Sprintf(`{"option_id":%d}`, int(optionID))
	w := testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/votes/%d/vote", int(voteID)), castBody, memberToken)
	testutil.AssertSuccessResponse(t, w)

	// 重复投票
	w = testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/votes/%d/vote", int(voteID)), castBody, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_DeleteVote_BeforeDeadline(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	voteResp := testutil.MakeRequest(r, "POST", "/api/votes", `{"title":"可删除","options":["A","B"]}`, memberToken)
	voteData := testutil.ParseDataMap(testutil.ParseResponse(t, voteResp))
	voteID := voteData["id"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/votes/%d", int(voteID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Forum_GetVote(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	voteResp := testutil.MakeRequest(r, "POST", "/api/votes", `{"title":"查看投票","options":["A","B"]}`, memberToken)
	voteData := testutil.ParseDataMap(testutil.ParseResponse(t, voteResp))
	voteID := voteData["id"].(float64)

	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/votes/%d", int(voteID)), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "查看投票", data["title"])
}

// === Feed ===

func TestHandler_Forum_Feed(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"动态1"}`, memberToken)
	testutil.MakeRequest(r, "POST", "/api/topics", `{"title":"话题1","content":"内容"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/feed", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["list"])
}

// === Tags ===

func TestHandler_Tag_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"name":"育儿"}`
	w := testutil.MakeRequest(r, "POST", "/api/tags", body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "育儿", data["name"])
}

func TestHandler_Tag_Create_ForbiddenForMember(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"name":"育儿"}`
	w := testutil.MakeRequest(r, "POST", "/api/tags", body, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Tag_Create_Duplicate(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	w := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	testutil.AssertErrorResponse(t, w, 409, 40002)
}

func TestHandler_Tag_List(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"家务"}`, adminToken)

	w := testutil.MakeRequest(r, "GET", "/api/tags", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.Len(t, data, 2)
}

func TestHandler_Tag_Update_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"旧名"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["id"].(float64)

	body := `{"name":"新名"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/tags/%d", int(tagID)), body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "新名", data["name"])
}

func TestHandler_Tag_Delete_InUse(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["id"].(float64)

	// 创建引用该标签的话题
	testutil.MakeRequest(r, "POST", "/api/topics",
		fmt.Sprintf(`{"title":"育儿话题","content":"内容","tag_id":%d}`, int(tagID)), memberToken)

	// 尝试删除标签
	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/tags/%d", int(tagID)), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40004)
}

func TestHandler_Tag_Delete_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := testutil.SetupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"可删除"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["id"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/tags/%d", int(tagID)), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/handler/ -run "TestHandler_Forum|TestHandler_Tag" -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/handler/forum_test.go
git commit -m "test: add forum handler HTTP integration tests (posts, topics, comments, likes, votes, tags, feed)"
```

---

## Task 10: Middleware 测试

**Files:**
- Create: `backend/internal/middleware/auth_test.go`

- [ ] **Step 1: 编写 Middleware 测试**

`backend/internal/middleware/auth_test.go`:
```go
package middleware

import (
	"net/http"
	"net/http/httptest"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
)

func init() {
	gin.SetMode(gin.TestMode)
}

func setupMiddlewareTest() {
	testutil.SetupTestDB()
	pkg.InitJWT("test-middleware-secret")
}

func TestAuthRequired_NoToken(t *testing.T) {
	setupMiddlewareTest()
	defer testutil.TeardownTestDB()

	r := gin.New()
	r.Use(AuthRequired())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_InvalidToken(t *testing.T) {
	setupMiddlewareTest()
	defer testutil.TeardownTestDB()

	r := gin.New()
	r.Use(AuthRequired())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer invalid-token")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
}

func TestAuthRequired_ValidToken(t *testing.T) {
	setupMiddlewareTest()
	defer testutil.TeardownTestDB()

	hash, _ := pkg.HashPassword("testpass")
	member := model.Member{Username: "testuser", Password: hash, Name: "Test", Role: "member", Status: "active"}
	pkg.DB.Create(&member)
	token, _ := pkg.GenerateToken(member.ID, member.Username, member.Role)

	r := gin.New()
	r.Use(AuthRequired())
	r.GET("/test", func(c *gin.Context) {
		memberID, _ := c.Get("member_id")
		role, _ := c.Get("role")
		c.JSON(200, gin.H{"member_id": memberID, "role": role})
	})

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAuthRequired_DisabledMember(t *testing.T) {
	setupMiddlewareTest()
	defer testutil.TeardownTestDB()

	hash, _ := pkg.HashPassword("testpass")
	member := model.Member{Username: "disabled", Password: hash, Name: "Disabled", Role: "member", Status: "disabled"}
	pkg.DB.Create(&member)
	token, _ := pkg.GenerateToken(member.ID, member.Username, member.Role)

	r := gin.New()
	r.Use(AuthRequired())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}

func TestAdminRequired_AsAdmin(t *testing.T) {
	setupMiddlewareTest()
	defer testutil.TeardownTestDB()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "admin")
		c.Next()
	})
	r.Use(AdminRequired())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
}

func TestAdminRequired_AsMember(t *testing.T) {
	setupMiddlewareTest()
	defer testutil.TeardownTestDB()

	r := gin.New()
	r.Use(func(c *gin.Context) {
		c.Set("role", "member")
		c.Next()
	})
	r.Use(AdminRequired())
	r.GET("/test", func(c *gin.Context) { c.JSON(200, gin.H{"ok": true}) })

	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/test", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 403, w.Code)
}
```

- [ ] **Step 2: 运行测试**

```bash
cd backend && go test ./internal/middleware/ -v
```

Expected: 所有测试 PASS

- [ ] **Step 3: Commit**

```bash
git add backend/internal/middleware/auth_test.go
git commit -m "test: add middleware AuthRequired and AdminRequired tests"
```

---

## Task 11: 全量测试验证与覆盖率报告

- [ ] **Step 1: 运行全部后端测试**

```bash
cd backend && go test ./... -v -count=1 2>&1 | tail -30
```

Expected: 所有测试 PASS（约 180+ 条测试）

- [ ] **Step 2: 检查测试覆盖率**

```bash
cd backend && go test ./... -coverprofile=coverage.out
go tool cover -func=coverage.out | tail -30
```

- [ ] **Step 3: 验证所有 57 个端点均有测试覆盖**

```bash
cd backend && go test ./internal/handler/ -v -count=1 2>&1 | grep "=== RUN" | wc -l
```

Expected: handler 测试数量 ≥ 50

- [ ] **Step 4: Commit 最终测试状态**

```bash
git add backend/coverage.out
git commit -m "test: add comprehensive backend API handler integration tests

- 57 API endpoints covered with HTTP-level tests
- Middleware tests (AuthRequired, AdminRequired)
- Tests verify PRD acceptance criteria and business rules
- All tests passing"
```

---

## 测试覆盖矩阵

| 模块 | 端点数 | Handler 测试数 | 覆盖的 PRD 验收标准 |
|------|--------|---------------|-------------------|
| Init | 1 | 3 | A-01 初始化 |
| Auth | 2 | 8 | A-02 登录, A-03 禁用用户, Token 校验 |
| Member | 7 | 12 | M-01~M-05, BR-P02, BR-P08 |
| Profile | 3 | 5 | P-01~P-03 |
| Category | 4 | 7 | C-01~C-04 |
| Ledger | 5 | 9 | L-01~L-07, BR-M01~M05 |
| Todo | 6 | 9 | T-01~T-06 |
| Wish | 8 | 12 | W-01~W-07 |
| Dashboard | 5 | 6 | D-01~D-05 |
| Forum | 17 | 22 | F-01~F-10 |
| Middleware | 2 | 6 | BR-A03, BR-A04, BR-P01 |
| **总计** | **60** | **~99** | |

## Task 依赖关系

```
Task 1 (基础设施)
  ↓
Task 2-9 (各模块 Handler 测试，均依赖 Task 1，但互相独立)
  ↓
Task 10 (Middleware 测试，独立)
  ↓
Task 11 (全量验证)
```
