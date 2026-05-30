package handler

import (
	"fmt"
	"testing"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Dashboard_Summary_Empty(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	month := time.Now().Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/dashboard/summary?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, float64(0), data["income"])
	assert.Equal(t, float64(0), data["expense"])
	assert.Equal(t, float64(0), data["balance"])
}

func TestHandler_Dashboard_Summary_WithData(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	incomeCat := testutil.SeedTestCategory("income", "工资", "💰", 1)
	expenseCat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)

	now := model.FromTime(time.Now())
	pkg.DB.Create(&model.Ledger{Amount: 10000, Note: "工资", CategoryID: incomeCat.ID, CreatorID: member.ID, OccurredAt: now})
	pkg.DB.Create(&model.Ledger{Amount: 3000, Note: "午餐", CategoryID: expenseCat.ID, CreatorID: member.ID, OccurredAt: now})

	month := now.ToTime().Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/dashboard/summary?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, float64(10000), data["income"])
	assert.Equal(t, float64(3000), data["expense"])
	assert.Equal(t, float64(7000), data["balance"])
}

func TestHandler_Dashboard_ExpenseChart(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)
	pkg.DB.Create(&model.Ledger{Amount: 5000, Note: "餐饮", CategoryID: cat.ID, CreatorID: member.ID, OccurredAt: model.FromTime(time.Now())})

	month := time.Now().Format("2006-01")
	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/dashboard/expense-chart?month=%s", month), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}

func TestHandler_Dashboard_UpcomingTodos(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建待办
	testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"近期待办","priority":"urgent"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/dashboard/upcoming-todos", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}

func TestHandler_Dashboard_WishTrends(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/wishes", `{"title":"愿望动态","category":"other"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/dashboard/wish-trends", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}

func TestHandler_Dashboard_ForumHot(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"论坛热点测试"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/dashboard/forum-hot", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.NotEmpty(t, data)
}
