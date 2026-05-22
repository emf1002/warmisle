package handler

import (
	"fmt"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Category_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
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

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"type":"expense","name":"交通","icon":"🚗"}`
	w := testutil.MakeRequest(r, "POST", "/api/categories", body, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Category_Create_DuplicateName(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"type":"expense","name":"餐饮","icon":"🍱"}`
	testutil.MakeRequest(r, "POST", "/api/categories", body, adminToken)

	w := testutil.MakeRequest(r, "POST", "/api/categories", body, adminToken)
	testutil.AssertErrorResponse(t, w, 409, 40002)
}

func TestHandler_Category_Create_InvalidType(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"type":"invalid","name":"测试","icon":"🔧"}`
	w := testutil.MakeRequest(r, "POST", "/api/categories", body, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Category_List_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
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

	r := setupTestRouter()
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

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "临时", "🔧", 99)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/categories/%d", cat.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Category_Delete_InUse(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	// 创建一条使用该分类的记账记录
	pkg.DB.Create(&model.Ledger{Amount: 1000, CategoryID: cat.ID, CreatorID: member.ID})

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/categories/%d", cat.ID), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}
