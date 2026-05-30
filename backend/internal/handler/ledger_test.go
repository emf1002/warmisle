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
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":3550,"note":"午餐","category_id":%d,"member_ids":[%d],"occurred_at":"%s"}`,
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

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":0,"note":"零元","category_id":%d,"member_ids":[%d]}`, cat.ID, member.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_NegativeAmount(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":-10,"note":"负数","category_id":%d,"member_ids":[%d]}`, cat.ID, member.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_NoMembers(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	body := fmt.Sprintf(`{"amount":10,"note":"无成员","category_id":%d,"member_ids":[]}`, cat.ID)
	w := testutil.MakeRequest(r, "POST", "/api/ledgers", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Ledger_Create_CategoryNotFound(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	body := fmt.Sprintf(`{"amount":10,"note":"分类不存在","category_id":99999,"member_ids":[%d]}`, member.ID)
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
	pkg.DB.Create(&model.LedgerMember{LedgerID: ledger.ID, MemberID: member.ID})

	month := now.Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/ledgers?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["summary"])
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
	pkg.DB.Create(&model.LedgerMember{LedgerID: ledger.ID, MemberID: member.ID})

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
	pkg.DB.Create(&model.LedgerMember{LedgerID: ledger.ID, MemberID: member.ID})

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
	pkg.DB.Create(&model.LedgerMember{LedgerID: ledger.ID, MemberID: member.ID})

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/ledgers/%d", ledger.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)

	// 验证软删除
	var count int64
	pkg.DB.Unscoped().Model(&model.Ledger{}).Where("id = ? AND deleted_at IS NOT NULL", ledger.ID).Count(&count)
	require.Equal(t, int64(1), count)
}
