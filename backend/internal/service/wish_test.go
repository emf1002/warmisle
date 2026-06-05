package service

import (
	"testing"

	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupWishTest() (*WishService, func()) {
	testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewWishService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestWishService_Create_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher", "Wisher", "member")

	wish, err := svc.Create("想买iPad", "画画用", "item", "important", "personal", nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "想买iPad", wish.Title)
	assert.Equal(t, "item", wish.Category)
	assert.Equal(t, "personal", wish.Type, "default type should be personal")
	assert.Equal(t, "pending", wish.Status)
	assert.Equal(t, creator.ID, wish.CreatorID)
}

func TestWishService_Create_TitleRequired(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher2", "Wisher2", "member")

	_, err := svc.Create("", "", "other", "normal", "", nil, creator.ID)
	assert.ErrorIs(t, err, ErrWishTitleRequired)
}

func TestWishService_Create_InvalidCategory(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher3", "Wisher3", "member")

	_, err := svc.Create("test", "", "invalid_category", "normal", "", nil, creator.ID)
	assert.ErrorIs(t, err, ErrWishInvalidCategory)
}

func TestWishService_Create_DefaultCategory(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher4", "Wisher4", "member")

	wish, err := svc.Create("test", "", "", "normal", "", nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "other", wish.Category)
}

func TestWishService_Create_ValidCategories(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wisher5", "Wisher5", "member")

	for _, cat := range []string{"item", "travel", "experience", "other"} {
		wish, err := svc.Create("test "+cat, "", cat, "normal", "", nil, creator.ID)
		require.NoError(t, err)
		assert.Equal(t, cat, wish.Category)
	}
}

func TestWishService_Promote_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "promoter", "Promoter", "member")

	wish, _ := svc.Create("提升测试", "", "other", "normal", "", nil, creator.ID)
	assert.Equal(t, "personal", wish.Type)

	promoted, err := svc.Promote(wish.ID, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "family", promoted.Type)
}

func TestWishService_Promote_NotCreator(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "promoter2", "Promoter2", "member")
	other := testutil.CreateTestMember(pkg.DB, "other", "Other", "member")

	wish, _ := svc.Create("不能提升别人的愿望", "", "other", "normal", "", nil, creator.ID)

	_, err := svc.Promote(wish.ID, other.ID)
	assert.ErrorIs(t, err, ErrWishPermissionDenied)
}

func TestWishService_Vote_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "voter_creator", "VCreator", "member")
	voter := testutil.CreateTestMember(pkg.DB, "voter", "Voter", "member")

	wish, _ := svc.Create("投票测试", "", "other", "normal", "", nil, creator.ID)

	result, err := svc.Vote(wish.ID, voter.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.VoteCount)
}

func TestWishService_Vote_Duplicate(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vc2", "VC2", "member")
	voter := testutil.CreateTestMember(pkg.DB, "dvoter", "DVoter", "member")

	wish, _ := svc.Create("重复投票", "", "other", "normal", "", nil, creator.ID)

	svc.Vote(wish.ID, voter.ID)
	_, err := svc.Vote(wish.ID, voter.ID)
	assert.ErrorIs(t, err, ErrWishAlreadyVoted)
}

func TestWishService_Unvote_Success(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "vc3", "VC3", "member")
	voter := testutil.CreateTestMember(pkg.DB, "unvoter", "Unvoter", "member")

	wish, _ := svc.Create("取消投票", "", "other", "normal", "", nil, creator.ID)
	svc.Vote(wish.ID, voter.ID)

	result, err := svc.Unvote(wish.ID, voter.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(0), result.VoteCount)
}

func TestWishService_UpdateStatus_AdminAnyStatus(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc3", "WC3", "member")
	admin := testutil.CreateTestMember(pkg.DB, "admin_wish", "AdminW", "admin")

	wish, _ := svc.Create("状态测试", "", "other", "normal", "", nil, creator.ID)

	for _, status := range []string{"pending", "agreed", "achieved", "abandoned"} {
		result, err := svc.UpdateStatus(wish.ID, status, admin.ID, "admin")
		require.NoError(t, err, "admin should be able to set status=%s", status)
		assert.Equal(t, status, result.Status)
	}
}

func TestWishService_UpdateStatus_CreatorOnlyAbandon(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc4", "WC4", "member")

	wish, _ := svc.Create("放弃测试", "", "other", "normal", "", nil, creator.ID)

	result, err := svc.UpdateStatus(wish.ID, "abandoned", creator.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, "abandoned", result.Status)
}

func TestWishService_UpdateStatus_CreatorCannotSetAgreed(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc5", "WC5", "member")

	wish, _ := svc.Create("不能改已同意", "", "other", "normal", "", nil, creator.ID)

	_, err := svc.UpdateStatus(wish.ID, "agreed", creator.ID, "member")
	assert.ErrorIs(t, err, ErrWishPermissionDenied)
}

func TestWishService_UpdateStatus_InvalidStatus(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc6", "WC6", "member")
	admin := testutil.CreateTestMember(pkg.DB, "admin6", "Admin6", "admin")

	wish, _ := svc.Create("无效状态", "", "other", "normal", "", nil, creator.ID)

	_, err := svc.UpdateStatus(wish.ID, "deleted", admin.ID, "admin")
	assert.ErrorIs(t, err, ErrWishInvalidStatus)
}

func TestWishService_Delete_ByCreator(t *testing.T) {
	svc, teardown := setupWishTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "wc7", "WC7", "member")

	wish, _ := svc.Create("删除愿望", "", "other", "normal", "", nil, creator.ID)

	err := svc.Delete(wish.ID, creator.ID, "member")
	require.NoError(t, err)

	_, err = svc.FindByID(wish.ID)
	assert.ErrorIs(t, err, ErrWishNotFound)
}
