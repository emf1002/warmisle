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
