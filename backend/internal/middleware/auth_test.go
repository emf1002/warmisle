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
