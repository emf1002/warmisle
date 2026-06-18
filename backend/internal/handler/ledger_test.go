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

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":3550,"note":"午餐","category_id":%d,"occurred_at":"%s"}`,
		cat.ID, time.Now().Format(time.RFC3339))
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, float64(3550), data["amount"])
	assert.Equal(t, "午餐", data["note"])
}

func TestHandler_Ledger_Create_ZeroAmount(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":0,"note":"零元","category_id":%d}`, cat.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_NegativeAmount(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":-10,"note":"负数","category_id":%d}`, cat.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_CategoryNotFound(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := fmt.Sprintf(`{"amount":10,"note":"分类不存在","category_id":99999}`)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_List_ByMonth(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	// 创建当月记录 (amount in cents)
	now := time.Now()
	ledger := model.Ledger{Amount: 3000, Note: "当月", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: model.FromTime(now)}
	pkg.DB.Create(&ledger)

	startDate := now.Format("2006-01") + "-01"
	firstOfNext := time.Date(now.Year(), now.Month()+1, 1, 0, 0, 0, 0, time.UTC)
	endDate := firstOfNext.Format("2006-01-02")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?start_date=%s&end_date=%s&limit=20", startDate, endDate), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["summary"])
	assert.NotNil(t, data["groups"])
	assert.NotNil(t, data["has_more"])
}

func TestHandler_Ledger_List_CursorPagination(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	// Create 12 records on different dates
	for i := 1; i <= 12; i++ {
		ledger := model.Ledger{
			Amount:     int64(i * 100),
			Note:       fmt.Sprintf("记录%d", i),
			CategoryID: cat.ID,
			CreatorID:  member.ID,
			OccurredAt: model.FromTime(time.Date(2026, 5, i, 12, 0, 0, 0, time.UTC)),
		}
		pkg.DB.Create(&ledger)
	}

	startDate := "2026-05-01"
	endDate := "2026-06-01"

	// Page 1: limit=5
	w1 := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?start_date=%s&end_date=%s&limit=5", startDate, endDate), nil, memberToken)
	resp1 := testutil.AssertSuccessResponse(t, w1)
	data1 := testutil.ParseDataMap(resp1)

	groups1 := data1["groups"].([]interface{})
	totalItems1 := 0
	for _, g := range groups1 {
		items := g.(map[string]interface{})["items"].([]interface{})
		totalItems1 += len(items)
	}
	assert.Equal(t, 5, totalItems1)
	assert.Equal(t, true, data1["has_more"])
	assert.NotNil(t, data1["next_cursor"])

	// Page 2: use cursor from page 1
	cursor := data1["next_cursor"].(string)
	w2 := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?start_date=%s&end_date=%s&limit=5&cursor=%s", startDate, endDate, cursor), nil, memberToken)
	resp2 := testutil.AssertSuccessResponse(t, w2)
	data2 := testutil.ParseDataMap(resp2)

	groups2 := data2["groups"].([]interface{})
	totalItems2 := 0
	for _, g := range groups2 {
		items := g.(map[string]interface{})["items"].([]interface{})
		totalItems2 += len(items)
	}
	assert.Equal(t, 5, totalItems2)
	assert.Equal(t, true, data2["has_more"])
}

func TestHandler_Ledger_List_InvalidCursor(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	startDate := "2026-05-01"
	endDate := "2026-06-01"
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?start_date=%s&end_date=%s&cursor=not-a-valid-cursor", startDate, endDate), nil, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Update_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	// 创建记录 (amount in cents)
	ledger := model.Ledger{Amount: 2000, Note: "原始", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: model.FromTime(time.Now())}
	pkg.DB.Create(&ledger)

	// Update: amount 5000 cents = 50.00 yuan
	body := `{"amount":5000,"note":"修改后"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/ledgers/%d", ledger.ID), body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "修改后", data["note"])
	assert.Equal(t, float64(5000), data["amount"])
}

func TestHandler_Ledger_Update_ByNonCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, _ := testutil.SeedAdminAndMember(t)

	// 创建第二个成员
	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "member2", Password: hash, Name: "成员2", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	ledger := model.Ledger{Amount: 2000, Note: "原始", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: model.FromTime(time.Now())}
	pkg.DB.Create(&ledger)

	body := `{"amount":50.00}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/ledgers/%d", ledger.ID), body, m2Token)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Ledger_Delete_ByAdmin(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	ledger := model.Ledger{Amount: 1000, Note: "待删除", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: model.FromTime(time.Now())}
	pkg.DB.Create(&ledger)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/ledgers/%d", ledger.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)

	// 验证软删除
	var count int64
	pkg.DB.Unscoped().Model(&model.Ledger{}).Where("id = ? AND deleted_at IS NOT NULL", ledger.ID).Count(&count)
	require.Equal(t, int64(1), count)
}
