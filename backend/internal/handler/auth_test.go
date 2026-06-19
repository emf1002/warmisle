package handler

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"strings"
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

func setupAuthRouter() *gin.Engine {
	r := gin.New()
	authH := NewAuthHandler()
	r.GET("/api/init/check", authH.InitCheck)
	r.POST("/api/auth/login", authH.Login)
	return r
}

func TestHandler_InitCheck_NeedInit(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/init/check", nil)
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp) //nolint:errcheck
	assert.Equal(t, float64(0), resp["code"])
	data := resp["data"].(map[string]interface{})
	assert.True(t, data["need_init"].(bool))
}

func TestHandler_InitCheck_NoNeedInit(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	hash, _ := pkg.HashPassword("pass")
	pkg.DB.Create(&model.Member{Username: "existing", Password: hash, Name: "Exist", Role: "admin", Status: "active"})

	r := setupAuthRouter()
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/api/init/check", nil)
	r.ServeHTTP(w, req)

	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp) //nolint:errcheck
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

	r := setupAuthRouter()
	body := `{"username":"loginuser","password":"correctpw"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 200, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp) //nolint:errcheck
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

	r := setupAuthRouter()
	body := `{"username":"wrongpw","password":"wrongpassword"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 401, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp) //nolint:errcheck
	assert.Equal(t, float64(40101), resp["code"])
}

func TestHandler_Login_EmptyCredentials(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	r := setupAuthRouter()
	body := `{"username":"","password":""}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp) //nolint:errcheck
	assert.Equal(t, float64(40001), resp["code"])
}

func TestHandler_Login_MissingPassword(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	r := setupAuthRouter()
	body := `{"username":"someone"}`
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("POST", "/api/auth/login", strings.NewReader(body))
	req.Header.Set("Content-Type", "application/json")
	r.ServeHTTP(w, req)

	assert.Equal(t, 400, w.Code)
	var resp map[string]interface{}
	json.Unmarshal(w.Body.Bytes(), &resp) //nolint:errcheck
	assert.Equal(t, float64(40001), resp["code"])
}

func TestHandler_Login_DisabledUser(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	r := setupAuthRouter()

	// 先初始化系统
	initBody := `{"admin_name":"管理员","username":"admin","password":"admin123"}`
	testutil.MakeRequest(r, "POST", "/api/init/setup", initBody, "")

	// 创建禁用成员
	testutil.SeedDisabledMember(t)

	loginBody := `{"username":"disabled","password":"testpass123"}`
	w := testutil.MakeRequest(r, "POST", "/api/auth/login", loginBody, "")
	testutil.AssertErrorResponse(t, w, 401, 40101)
}

func TestHandler_ProtectedEndpoint_NoToken(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	// 需要用完整路由才能测试受保护端点
	fullRouter := setupTestRouter()
	w := testutil.MakeRequest(fullRouter, "GET", "/api/members", nil, "")
	testutil.AssertErrorResponse(t, w, 401, 40101)
}

func TestHandler_ProtectedEndpoint_InvalidToken(t *testing.T) {
	testutil.SetupTestDB()
	defer testutil.TeardownTestDB()
	pkg.InitJWT("test-secret")

	fullRouter := setupTestRouter()
	w := testutil.MakeRequest(fullRouter, "GET", "/api/members", nil, "invalid-token")
	testutil.AssertErrorResponse(t, w, 401, 40101)
}
