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

	r := setupTestRouter()

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

	r := setupTestRouter()

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

	r := setupTestRouter()

	body := `{"username":"admin"}`
	w := testutil.MakeRequest(r, "POST", "/api/init/setup", body, "")
	testutil.AssertErrorResponse(t, w, 400, 40001)
}
