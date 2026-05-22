package handler

import (
	"fmt"
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
)

func TestHandler_Member_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"username":"newuser","password":"pass123","name":"新成员","avatar":"👩"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "newuser", data["username"])
	assert.Equal(t, "新成员", data["name"])
}

func TestHandler_Member_Create_DuplicateUsername(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"username":"admin","password":"pass123","name":"重复"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, adminToken)
	testutil.AssertErrorResponse(t, w, 409, 40002)
}

func TestHandler_Member_Create_InvalidUsername(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"username":"ab","password":"pass123","name":"短用户名"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Member_Create_ForbiddenForMember(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"username":"newuser","password":"pass123","name":"新成员"}`
	w := testutil.MakeRequest(r, "POST", "/api/members", body, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Member_List_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "GET", "/api/members", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataArray(resp)
	assert.Len(t, data, 2)
}

func TestHandler_Member_Update_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"name":"修改后名称","avatar":"🐶"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d", member.ID), body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)

	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "修改后名称", data["name"])
}

func TestHandler_Member_Update_CannotRemoveLastAdmin(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	admin, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"role":"member"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d", admin.ID), body, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40003)
}

func TestHandler_Member_Delete_ForbiddenForMember(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/members/%d", member.ID), nil, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Member_Delete_WithActivityRecords(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	// 给成员创建一条记账记录
	cat := testutil.SeedTestCategory("expense", "餐饮", "🍱", 1)
	pkg.DB.Create(&model.Ledger{Amount: 1000, CategoryID: cat.ID, CreatorID: member.ID})

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/members/%d", member.ID), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40004)
}

func TestHandler_Member_Disable_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/disable", member.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Member_Disable_Self(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	admin, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/disable", admin.ID), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40005)
}

func TestHandler_Member_Enable_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	// 先禁用
	testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/disable", member.ID), nil, adminToken)

	// 再启用
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/enable", member.ID), nil, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "active", data["status"])
}

func TestHandler_Member_ResetPassword(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, member, adminToken, _ := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/members/%d/reset-pwd", member.ID), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Profile_Get(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	w := testutil.MakeRequest(r, "GET", "/api/profile", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "member", data["username"])
}

func TestHandler_Profile_Update(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"name":"新名字","avatar":"🔥"}`
	w := testutil.MakeRequest(r, "PUT", "/api/profile", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "新名字", data["name"])
}

func TestHandler_Profile_ChangePassword_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"old_password":"testpass123","new_password":"newpass456"}`
	w := testutil.MakeRequest(r, "PUT", "/api/profile/password", body, memberToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Profile_ChangePassword_WrongOld(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"old_password":"wrongold","new_password":"newpass456"}`
	w := testutil.MakeRequest(r, "PUT", "/api/profile/password", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40010)
}
