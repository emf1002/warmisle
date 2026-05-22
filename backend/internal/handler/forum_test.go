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

// === Posts ===

func TestHandler_Forum_CreatePost_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"content":"今天天气真好！"}`
	w := testutil.MakeRequest(r, "POST", "/api/posts", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "今天天气真好！", data["Content"])
}

func TestHandler_Forum_CreatePost_EmptyContent(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"content":""}`
	w := testutil.MakeRequest(r, "POST", "/api/posts", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_DeletePost_ByCreator(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"待删除"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	postID := data["ID"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/posts/%d", int(postID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Forum_DeletePost_ByAdmin(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"管理员删除"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	postID := data["ID"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/posts/%d", int(postID)), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Forum_DeletePost_Forbidden(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	resp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"不能删"}`, memberToken)
	data := testutil.ParseDataMap(testutil.ParseResponse(t, resp))
	postID := data["ID"].(float64)

	// 创建另一个成员
	hash, _ := pkg.HashPassword("testpass123")
	m2 := model.Member{Username: "other", Password: hash, Name: "其他", Avatar: "👶", Role: "member", Status: "active"}
	pkg.DB.Create(&m2)
	m2Token := testutil.GenerateTestToken(m2)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/posts/%d", int(postID)), nil, m2Token)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

// === Topics ===

func TestHandler_Forum_CreateTopic_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"讨论话题","content":"话题内容"}`
	w := testutil.MakeRequest(r, "POST", "/api/topics", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "讨论话题", data["Title"])
}

func TestHandler_Forum_CreateTopic_EmptyTitle(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"","content":"内容"}`
	w := testutil.MakeRequest(r, "POST", "/api/topics", body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_CreateTopic_WithTag(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	// 创建标签
	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["ID"].(float64)

	body := fmt.Sprintf(`{"title":"育儿话题","content":"讨论育儿","tag_id":%d}`, int(tagID))
	w := testutil.MakeRequest(r, "POST", "/api/topics", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["tag"])
}

func TestHandler_Forum_TogglePin_AdminOnly(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	// 创建话题
	topicResp := testutil.MakeRequest(r, "POST", "/api/topics", `{"title":"可置顶","content":"内容"}`, memberToken)
	topicData := testutil.ParseDataMap(testutil.ParseResponse(t, topicResp))
	topicID := topicData["ID"].(float64)

	// 普通成员不能置顶
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/topics/%d/pin", int(topicID)), nil, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)

	// 管理员可以置顶
	w = testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/topics/%d/pin", int(topicID)), nil, adminToken)
	pinResp := testutil.AssertSuccessResponse(t, w)
	pinData := testutil.ParseDataMap(pinResp)
	assert.Equal(t, true, pinData["IsPinned"])
}

func TestHandler_Forum_GetTopic(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	topicResp := testutil.MakeRequest(r, "POST", "/api/topics", `{"title":"查看话题","content":"内容"}`, memberToken)
	topicData := testutil.ParseDataMap(testutil.ParseResponse(t, topicResp))
	topicID := topicData["ID"].(float64)

	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/topics/%d", int(topicID)), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "查看话题", data["Title"])
}

// === Comments ===

func TestHandler_Forum_CreateComment_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"评论目标"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["ID"].(float64)

	body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"好文"}`, int(postID))
	w := testutil.MakeRequest(r, "POST", "/api/comments", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "好文", data["Content"])
}

func TestHandler_Forum_CreateComment_ReplyToLevel1(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"回复测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["ID"].(float64)

	// 一级评论
	c1Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"一级"}`, int(postID))
	c1Resp := testutil.MakeRequest(r, "POST", "/api/comments", c1Body, memberToken)
	c1Data := testutil.ParseDataMap(testutil.ParseResponse(t, c1Resp))
	c1ID := c1Data["ID"].(float64)

	// 二级回复
	c2Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"二级回复"}`, int(postID), int(c1ID))
	w := testutil.MakeRequest(r, "POST", "/api/comments", c2Body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["ParentID"])
}

func TestHandler_Forum_CreateComment_NestingTooDeep(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"嵌套测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["ID"].(float64)

	// 一级
	c1Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"一级"}`, int(postID))
	c1Resp := testutil.MakeRequest(r, "POST", "/api/comments", c1Body, memberToken)
	c1Data := testutil.ParseDataMap(testutil.ParseResponse(t, c1Resp))
	c1ID := c1Data["ID"].(float64)

	// 二级
	c2Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"二级"}`, int(postID), int(c1ID))
	c2Resp := testutil.MakeRequest(r, "POST", "/api/comments", c2Body, memberToken)
	c2Data := testutil.ParseDataMap(testutil.ParseResponse(t, c2Resp))
	c2ID := c2Data["ID"].(float64)

	// 三级应失败
	c3Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"三级不允许"}`, int(postID), int(c2ID))
	w := testutil.MakeRequest(r, "POST", "/api/comments", c3Body, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_DeleteComment_SyncDeleteReplies(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"删除测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["ID"].(float64)

	c1Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"content":"一级"}`, int(postID))
	c1Resp := testutil.MakeRequest(r, "POST", "/api/comments", c1Body, memberToken)
	c1Data := testutil.ParseDataMap(testutil.ParseResponse(t, c1Resp))
	c1ID := c1Data["ID"].(float64)

	c2Body := fmt.Sprintf(`{"target_type":"post","target_id":%d,"parent_id":%d,"content":"二级"}`, int(postID), int(c1ID))
	c2Resp := testutil.MakeRequest(r, "POST", "/api/comments", c2Body, memberToken)
	c2Data := testutil.ParseDataMap(testutil.ParseResponse(t, c2Resp))
	c2ID := c2Data["ID"].(float64)

	// 删除一级
	testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/comments/%d", int(c1ID)), nil, memberToken)

	// 二级应被同步软删除
	var count int64
	pkg.DB.Unscoped().Model(&model.Comment{}).Where("id = ? AND deleted_at IS NOT NULL", int(c2ID)).Count(&count)
	assert.Equal(t, int64(1), count)
}

// === Likes ===

func TestHandler_Forum_ToggleLike(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	postResp := testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"点赞测试"}`, memberToken)
	postData := testutil.ParseDataMap(testutil.ParseResponse(t, postResp))
	postID := postData["ID"].(float64)

	// 点赞
	body := fmt.Sprintf(`{"target_type":"post","target_id":%d}`, int(postID))
	w := testutil.MakeRequest(r, "POST", "/api/likes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, true, data["liked"])

	// 取消点赞
	w = testutil.MakeRequest(r, "POST", "/api/likes", body, memberToken)
	resp = testutil.AssertSuccessResponse(t, w)
	data = testutil.ParseDataMap(resp)
	assert.Equal(t, false, data["liked"])
}

// === Votes ===

func TestHandler_Forum_CreateVote_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"title":"周末去哪","options":["公园","海边","爬山"],"is_multi":false}`
	w := testutil.MakeRequest(r, "POST", "/api/votes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "周末去哪", data["Title"])
}

func TestHandler_Forum_CreateVote_WithDeadline(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	deadline := time.Now().Add(24 * time.Hour).Format(time.RFC3339)
	body := fmt.Sprintf(`{"title":"限时投票","options":["A","B"],"deadline":"%s"}`, deadline)
	w := testutil.MakeRequest(r, "POST", "/api/votes", body, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["Deadline"])
}

func TestHandler_Forum_Vote_CastAndDuplicate(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	voteResp := testutil.MakeRequest(r, "POST", "/api/votes", `{"title":"投票","options":["A","B"]}`, memberToken)
	voteData := testutil.ParseDataMap(testutil.ParseResponse(t, voteResp))
	voteID := voteData["ID"].(float64)
	options := voteData["options"].([]interface{})
	optionID := options[0].(map[string]interface{})["ID"].(float64)

	// 投票
	castBody := fmt.Sprintf(`{"option_id":%d}`, int(optionID))
	w := testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/votes/%d/vote", int(voteID)), castBody, memberToken)
	testutil.AssertSuccessResponse(t, w)

	// 重复投票
	w = testutil.MakeRequest(r, "POST", fmt.Sprintf("/api/votes/%d/vote", int(voteID)), castBody, memberToken)
	testutil.AssertErrorResponse(t, w, 400, 40001)
}

func TestHandler_Forum_DeleteVote_BeforeDeadline(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	voteResp := testutil.MakeRequest(r, "POST", "/api/votes", `{"title":"可删除","options":["A","B"]}`, memberToken)
	voteData := testutil.ParseDataMap(testutil.ParseResponse(t, voteResp))
	voteID := voteData["ID"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/votes/%d", int(voteID)), nil, memberToken)
	testutil.AssertSuccessResponse(t, w)
}

func TestHandler_Forum_GetVote(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	voteResp := testutil.MakeRequest(r, "POST", "/api/votes", `{"title":"查看投票","options":["A","B"]}`, memberToken)
	voteData := testutil.ParseDataMap(testutil.ParseResponse(t, voteResp))
	voteID := voteData["ID"].(float64)

	w := testutil.MakeRequest(r, "GET", fmt.Sprintf("/api/votes/%d", int(voteID)), nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "查看投票", data["Title"])
}

// === Feed ===

func TestHandler_Forum_Feed(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/posts", `{"content":"动态1"}`, memberToken)
	testutil.MakeRequest(r, "POST", "/api/topics", `{"title":"话题1","content":"内容"}`, memberToken)

	w := testutil.MakeRequest(r, "GET", "/api/feed", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.NotNil(t, data["items"])
}

// === Tags ===

func TestHandler_Tag_Create_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	body := `{"name":"育儿"}`
	w := testutil.MakeRequest(r, "POST", "/api/tags", body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "育儿", data["Name"])
}

func TestHandler_Tag_Create_ForbiddenForMember(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, _, memberToken := testutil.SeedAdminAndMember(t)

	body := `{"name":"育儿"}`
	w := testutil.MakeRequest(r, "POST", "/api/tags", body, memberToken)
	testutil.AssertErrorResponse(t, w, 403, 40301)
}

func TestHandler_Tag_Create_Duplicate(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	w := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	testutil.AssertErrorResponse(t, w, 409, 40002)
}

func TestHandler_Tag_List(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"家务"}`, adminToken)

	w := testutil.MakeRequest(r, "GET", "/api/tags", nil, memberToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataArray(resp)
	assert.Len(t, data, 2)
}

func TestHandler_Tag_Update_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"旧名"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["ID"].(float64)

	body := `{"name":"新名"}`
	w := testutil.MakeRequest(r, "PUT", fmt.Sprintf("/api/tags/%d", int(tagID)), body, adminToken)
	resp := testutil.AssertSuccessResponse(t, w)
	data := testutil.ParseDataMap(resp)
	assert.Equal(t, "新名", data["Name"])
}

func TestHandler_Tag_Delete_InUse(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, memberToken := testutil.SeedAdminAndMember(t)

	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"育儿"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["ID"].(float64)

	// 创建引用该标签的话题
	testutil.MakeRequest(r, "POST", "/api/topics",
		fmt.Sprintf(`{"title":"育儿话题","content":"内容","tag_id":%d}`, int(tagID)), memberToken)

	// 尝试删除标签
	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/tags/%d", int(tagID)), nil, adminToken)
	testutil.AssertErrorResponse(t, w, 400, 40004)
}

func TestHandler_Tag_Delete_Success(t *testing.T) {
	setupMemberTest()
	defer testutil.TeardownTestDB()

	r := setupTestRouter()
	_, _, adminToken, _ := testutil.SeedAdminAndMember(t)

	tagResp := testutil.MakeRequest(r, "POST", "/api/tags", `{"name":"可删除"}`, adminToken)
	tagData := testutil.ParseDataMap(testutil.ParseResponse(t, tagResp))
	tagID := tagData["ID"].(float64)

	w := testutil.MakeRequest(r, "DELETE", fmt.Sprintf("/api/tags/%d", int(tagID)), nil, adminToken)
	testutil.AssertSuccessResponse(t, w)
}
