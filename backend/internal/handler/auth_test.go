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
	json.Unmarshal(w.Body.Bytes(), &resp)
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

	r := setupAuthRouter()
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

	r := setupAuthRouter()
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

	r := setupAuthRouter()
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
	json.Unmarshal(w.Body.Bytes(), &resp)
	assert.Equal(t, float64(40001), resp["code"])
}
