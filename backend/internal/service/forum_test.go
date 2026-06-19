package service

import (
	"testing"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupForumTest() (*ForumService, func()) {
	testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewForumService()
	return svc, func() { testutil.TeardownTestDB() }
}

// --- Posts ---

func TestForumService_CreatePost_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "poster", "Poster", "member")
	post, err := svc.CreatePost("Hello world!", creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "Hello world!", post.Content)
	assert.Equal(t, creator.ID, post.CreatorID)
}

func TestForumService_CreatePost_EmptyContent(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "poster2", "Poster2", "member")
	_, err := svc.CreatePost("", creator.ID)
	assert.ErrorIs(t, err, ErrForumContentRequired)
}

func TestForumService_DeletePost_ByCreator(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "poster3", "Poster3", "member")
	post, _ := svc.CreatePost("Delete me", creator.ID)
	err := svc.DeletePost(post.ID, creator.ID, "member")
	require.NoError(t, err)
	// Verify soft delete via DB
	var p model.Post
	err = pkg.DB.First(&p, post.ID).Error
	assert.Error(t, err)
}

func TestForumService_DeletePost_ByNonCreator(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "poster4", "Poster4", "member")
	other := testutil.CreateTestMember(pkg.DB, "other4", "Other4", "member")
	post, _ := svc.CreatePost("Can't delete", creator.ID)
	err := svc.DeletePost(post.ID, other.ID, "member")
	assert.ErrorIs(t, err, ErrForumPermissionDenied)
}

func TestForumService_DeletePost_ByAdmin(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "poster5", "Poster5", "member")
	admin := testutil.CreateTestMember(pkg.DB, "admin5", "Admin5", "admin")
	post, _ := svc.CreatePost("Admin can delete", creator.ID)
	err := svc.DeletePost(post.ID, admin.ID, "admin")
	require.NoError(t, err)
}

// --- Topics ---

func TestForumService_CreateTopic_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "topicer", "Topicer", "member")
	topic, err := svc.CreateTopic("讨论标题", "讨论内容", nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "讨论标题", topic.Title)
	assert.Equal(t, "讨论内容", topic.Content)
	assert.Equal(t, creator.ID, topic.CreatorID)
}

func TestForumService_CreateTopic_TitleRequired(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "topicer2", "Topicer2", "member")
	_, err := svc.CreateTopic("", "内容", nil, creator.ID)
	assert.ErrorIs(t, err, ErrForumTitleRequired)
}

func TestForumService_TogglePin(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "pinner_c", "PinnerC", "member")
	testutil.CreateTestMember(pkg.DB, "pinner_a", "PinnerA", "admin")
	topic, _ := svc.CreateTopic("可置顶", "内容", nil, creator.ID)
	assert.False(t, topic.IsPinned)
	pinned, err := svc.TogglePin(topic.ID, "admin")
	require.NoError(t, err)
	assert.True(t, pinned.IsPinned)
	unpinned, err := svc.TogglePin(topic.ID, "admin")
	require.NoError(t, err)
	assert.False(t, unpinned.IsPinned)
}

func TestForumService_TogglePin_NotAdmin(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "member_pin", "MemberPin", "member")
	topic, _ := svc.CreateTopic("不能置顶", "内容", nil, creator.ID)
	_, err := svc.TogglePin(topic.ID, "member")
	assert.ErrorIs(t, err, ErrForumPermissionDenied)
}

// --- Comments ---

func TestForumService_CreateComment_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "commenter", "Commenter", "member")
	post, _ := svc.CreatePost("文章", creator.ID)
	comment, err := svc.CreateComment("post", post.ID, nil, "好文章", creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "好文章", comment.Content)
	assert.Equal(t, "post", comment.TargetType)
}

func TestForumService_CreateComment_ReplyToLevel1(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "c1", "C1", "member")
	replyer := testutil.CreateTestMember(pkg.DB, "c2", "C2", "member")
	post, _ := svc.CreatePost("文章", creator.ID)
	level1, _ := svc.CreateComment("post", post.ID, nil, "一级评论", creator.ID)
	level2, err := svc.CreateComment("post", post.ID, &level1.ID, "回复你", replyer.ID)
	require.NoError(t, err)
	require.NotNil(t, level2.ParentID)
	assert.Equal(t, level1.ID, *level2.ParentID)
}

func TestForumService_CreateComment_NestingTooDeep(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "nc1", "NC1", "member")
	replyer := testutil.CreateTestMember(pkg.DB, "nc2", "NC2", "member")
	third := testutil.CreateTestMember(pkg.DB, "nc3", "NC3", "member")
	post, _ := svc.CreatePost("嵌套测试", creator.ID)
	level1, _ := svc.CreateComment("post", post.ID, nil, "一级", creator.ID)
	level2, _ := svc.CreateComment("post", post.ID, &level1.ID, "二级", replyer.ID)
	_, err := svc.CreateComment("post", post.ID, &level2.ID, "三级不允许", third.ID)
	assert.ErrorIs(t, err, ErrForumNestingTooDeep)
}

func TestForumService_DeleteComment_SyncDeleteReplies(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "dc1", "DC1", "member")
	replyer := testutil.CreateTestMember(pkg.DB, "dc2", "DC2", "member")
	post, _ := svc.CreatePost("文章", creator.ID)
	level1, _ := svc.CreateComment("post", post.ID, nil, "一级", creator.ID)
	level2, _ := svc.CreateComment("post", post.ID, &level1.ID, "二级", replyer.ID)
	// Delete level 1 comment
	err := svc.DeleteComment(level1.ID, creator.ID, "member")
	require.NoError(t, err)
	// Level 2 should be synced soft-deleted — query via repo
	var c model.Comment
	err = pkg.DB.First(&c, level2.ID).Error
	assert.Error(t, err)
}

// --- Likes ---

func TestForumService_ToggleLike(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "liker_c", "LikerC", "member")
	liker := testutil.CreateTestMember(pkg.DB, "liker", "Liker", "member")
	post, _ := svc.CreatePost("点赞测试", creator.ID)
	liked, err := svc.ToggleLike("post", post.ID, liker.ID)
	require.NoError(t, err)
	assert.True(t, liked)
	unliked, err := svc.ToggleLike("post", post.ID, liker.ID)
	require.NoError(t, err)
	assert.False(t, unliked)
}

// --- Votes ---

func TestForumService_CreateVote_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "vote_c", "VoteC", "member")
	vote, err := svc.CreateVote("周末去哪", []string{"公园", "海边", "爬山"}, false, nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "周末去哪", vote.Title)
	assert.Len(t, vote.Options, 3)
}

func TestForumService_CreateVote_TitleRequired(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "vote_c2", "VoteC2", "member")
	_, err := svc.CreateVote("", []string{"A", "B"}, false, nil, creator.ID)
	assert.ErrorIs(t, err, ErrForumTitleRequired)
}

func TestForumService_CreateVote_MinOptions(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "vote_c3", "VoteC3", "member")
	_, err := svc.CreateVote("只有一项", []string{"A"}, false, nil, creator.ID)
	assert.Error(t, err)
	assert.Contains(t, err.Error(), "at least 2 options")
}

func TestForumService_Vote_CastAndDuplicate(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "vc_c", "VCC", "member")
	voter := testutil.CreateTestMember(pkg.DB, "vc_v", "VCV", "member")
	vote, _ := svc.CreateVote("选哪个", []string{"A", "B"}, false, nil, creator.ID)
	result, err := svc.Vote(vote.ID, []uint{vote.Options[0].ID}, voter.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.TotalVotes)
	_, err = svc.Vote(vote.ID, []uint{vote.Options[1].ID}, voter.ID)
	assert.ErrorIs(t, err, ErrForumAlreadyVoted)
}

func TestForumService_DeleteVote_BeforeDeadline(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "dv_c", "DVC", "member")
	vote, _ := svc.CreateVote("可删除", []string{"A", "B"}, false, nil, creator.ID)
	err := svc.DeleteVote(vote.ID, creator.ID, "member")
	require.NoError(t, err)
}

func TestForumService_DeleteVote_AfterDeadline(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "dd_c", "DDC", "member")
	pastDeadline := time.Now().Add(-1 * time.Hour)
	vote, _ := svc.CreateVote("已过期", []string{"A", "B"}, false, &pastDeadline, creator.ID)
	err := svc.DeleteVote(vote.ID, creator.ID, "member")
	assert.ErrorIs(t, err, ErrForumVoteDeadlinePassed)
}

// --- Tags ---

func TestForumService_CreateTag_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	tag, err := svc.CreateTag("育儿")
	require.NoError(t, err)
	assert.Equal(t, "育儿", tag.Name)
}

func TestForumService_DeleteTag_InUse(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "tag_c", "TagC", "member")
	tag, _ := svc.CreateTag("育儿")
	svc.CreateTopic("育儿话题", "内容", &tag.ID, creator.ID) //nolint:errcheck
	err := svc.DeleteTag(tag.ID)
	assert.ErrorIs(t, err, ErrForumTagInUse)
}

func TestForumService_DeleteTag_Success(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	tag, _ := svc.CreateTag("无引用标签")
	err := svc.DeleteTag(tag.ID)
	require.NoError(t, err)
}

// --- Feed ---

func TestForumService_GetFeed_Empty(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	feed, err := svc.GetFeed(1, 10)
	require.NoError(t, err)
	assert.Empty(t, feed.Pinned)
	assert.Empty(t, feed.Items)
}

func TestForumService_GetFeed_WithContent(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "feed_c", "FeedC", "member")
	svc.CreatePost("Hello feed", creator.ID) //nolint:errcheck
	feed, err := svc.GetFeed(1, 10)
	require.NoError(t, err)
	assert.NotEmpty(t, feed.Items)
	assert.Equal(t, "post", feed.Items[0].Type)
}

func TestForumService_GetFeed_PinnedTopics(t *testing.T) {
	svc, teardown := setupForumTest()
	defer teardown()
	creator := testutil.CreateTestMember(pkg.DB, "feed_pin", "FeedPin", "admin")
	topic, _ := svc.CreateTopic("置顶话题", "内容", nil, creator.ID)
	svc.TogglePin(topic.ID, "admin") //nolint:errcheck
	feed, err := svc.GetFeed(1, 10)
	require.NoError(t, err)
	assert.Len(t, feed.Pinned, 1)
	assert.Equal(t, "置顶话题", feed.Pinned[0].Title)
}

