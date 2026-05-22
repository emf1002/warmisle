package handler

import (
	"fmt"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Todo_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"买菜","description":"去超市","priority":"important"}`
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "买菜", data["Title"])
	assert.Equal(t, "important", data["Priority"])
	assert.Equal(t, "pending", data["Status"])
}

func TestHandler_Todo_Create_EmptyTitle(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"","description":"空标题"}`
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Todo_Create_InvalidPriority(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"测试","priority":"super-urgent"}`
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Todo_Create_WithAssignee(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	body := fmt.Sprintf(`{"title":"指派任务","assignee_id":%d}`, member.ID)
	w := testutil.MakeRequest(r, "POST", "/api/todos", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["Assignee"])
}

func TestHandler_Todo_List_WithFilter(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建几条待办
	testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"待办1","priority":"urgent"}`, memberToken)
	testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"待办2","priority":"normal"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/todos", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["list"])
}

func TestHandler_Todo_Toggle_CompleteAndUncomplete(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建待办
	w := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"完成测试"}`, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	todoID := data["ID"].(float64)

	// 完成
	w = testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d/toggle", int(todoID)), nil, memberToken)
	toggleResp := testutil.AssertSuccessResponse(t, w)
	toggleData := testutil.ParseDataMap(toggleResp)
	assert.Equal(t, "completed", toggleData["Status"])

	// 恢复
	w = testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d/toggle", int(todoID)), nil, memberToken)
	toggleResp2 := testutil.AssertSuccessResponse(t, w)
	toggleData2 := testutil.ParseDataMap(toggleResp2)
	assert.Equal(t, "pending", toggleData2["Status"])
}

func TestHandler_Todo_Claim_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	// 创建未指派的待办
	w := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"认领测试"}`, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	todoID := data["ID"].(float64)

	// 创建第二个成员来认领
	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "claimer", Password: hash, Name: "认领者", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	w = testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d/claim", int(todoID)), nil, m2Token)
	claimResp := testutil.AssertSuccessResponse(t, w)
	claimData := testutil.ParseDataMap(claimResp)
	assert.NotNil(t, claimData["Assignee"])
}

func TestHandler_Todo_Update_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"原始标题"}`, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	todoID := data["ID"].(float64)

	body := `{"title":"修改后标题","priority":"urgent"}`
	w = testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/todos/%d", int(todoID)), body, memberToken)
	updateResp := testutil.AssertSuccessResponse(t, w)
	updateData := testutil.ParseDataMap(updateResp)
	assert.Equal(t, "修改后标题", updateData["Title"])
}

func TestHandler_Todo_Delete_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "POST", "/api/todos", `{"title":"待删除"}`, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	todoID := data["ID"].(float64)

	w = testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/todos/%d", int(todoID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}
