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

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"想买iPad","description":"画画用","category":"item","priority":"important"}`
	w := testutil.MakeRequest(r, "POST", "/api/wishes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "想买iPad", data["Title"])
	assert.Equal(t, "personal", data["Type"])
	assert.Equal(t, "pending", data["Status"])
}

func TestHandler_Wish_Create_EmptyTitle(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"","description":"空标题"}`
	w := testutil.MakeRequest(r, "POST", "/api/wishes", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Wish_Create_InvalidCategory(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"测试","category":"invalid_cat"}`
	w := testutil.MakeRequest(r, "POST", "/api/wishes", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Wish_Promote_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建个人愿望
	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"提升测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)

	// 提升为家庭愿望
	w := testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/wishes/%d/promote", int(wishID)), nil, memberToken)
	promoteResp := testutil.AssertSuccessResponse(t, w)
	promoteData := testutil.ParseDataMap(promoteResp)
	assert.Equal(t, "family", promoteData["Type"])
}

func TestHandler_Wish_Promote_NotCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建愿望
	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"不能提升","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)

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

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建家庭愿望
	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"投票测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)
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

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"重复投票","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)
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

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"取消投票","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)
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

	r := setupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"状态测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)

	for _, status := range []string{"agreed", "achieved"} {
		body := fmt.Sprintf(`{"status":"%s"}`, status)
		w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/wishes/%d/status", int(wishID)), body, adminToken)
		statusResp := testutil.AssertSuccessResponse(t, w)
		statusData := testutil.ParseDataMap(statusResp)
		assert.Equal(t, status, statusData["Status"])
	}
}

func TestHandler_Wish_UpdateStatus_CreatorOnlyAbandon(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"放弃测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)

	// 创建者可以标记放弃
	body := `{"status":"abandoned"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/wishes/%d/status", int(wishID)), body, memberToken)
	statusResp := testutil.AssertSuccessResponse(t, w)
	statusData := testutil.ParseDataMap(statusResp)
	assert.Equal(t, "abandoned", statusData["Status"])
}

func TestHandler_Wish_Delete_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"删除测试","category":"other"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	wishID := data["ID"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/wishes/%d", int(wishID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}
