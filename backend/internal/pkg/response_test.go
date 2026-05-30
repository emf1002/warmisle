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
