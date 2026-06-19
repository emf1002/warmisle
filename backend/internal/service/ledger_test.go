package service

import (
	"fmt"
	"testing"
	"time"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/repository"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupLedgerTest() (*LedgerService, func()) {
	testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewLedgerService()
	return svc, func() { testutil.TeardownTestDB() }
}

func seedLedgerFixtures() (member1, member2 model.Member, cat model.Category) {
	member1 = testutil.CreateTestMember(pkg.DB, "user1", "User1", "member")
	member2 = testutil.CreateTestMember(pkg.DB, "user2", "User2", "member")
	cat = testutil.CreateTestCategory(pkg.DB, "expense", "餐饮", "🍱", 1)
	return
}

func TestLedgerService_Create_Success(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	result, err := svc.Create(3550, "午餐", cat.ID, time.Now(), m1.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3550), result.Amount)
	assert.Equal(t, "午餐", result.Note)
	assert.Equal(t, m1.ID, result.CreatorID)
}

func TestLedgerService_Create_ZeroAmount(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	_, err := svc.Create(0, "免费", cat.ID, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestLedgerService_Create_NegativeAmount(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	_, err := svc.Create(-100, "负数", cat.ID, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestLedgerService_Create_CategoryNotFound(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, _ := seedLedgerFixtures()

	_, err := svc.Create(1000, "不存在的分类", 99999, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrLedgerCategoryNotFound)
}

func TestLedgerService_Update_ByCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	created, _ := svc.Create(2000, "原始", cat.ID, time.Now(), m1.ID)

	newAmount := int64(3000)
	newNote := "修改后"
	updated, err := svc.Update(created.ID, &newAmount, &newNote, nil, nil, m1.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, int64(3000), updated.Amount)
	assert.Equal(t, "修改后", updated.Note)
}

func TestLedgerService_Update_ByNonCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, m2, cat := seedLedgerFixtures()
	created, _ := svc.Create(2000, "原始", cat.ID, time.Now(), m1.ID)

	newAmount := int64(3000)
	_, err := svc.Update(created.ID, &newAmount, nil, nil, nil, m2.ID, "member")
	assert.ErrorIs(t, err, ErrLedgerPermissionDenied)
}

func TestLedgerService_Update_ByAdmin(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	admin := testutil.CreateTestMember(pkg.DB, "admin", "Admin", "admin")
	created, _ := svc.Create(2000, "原始", cat.ID, time.Now(), m1.ID)

	newAmount := int64(5000)
	updated, err := svc.Update(created.ID, &newAmount, nil, nil, nil, admin.ID, "admin")
	require.NoError(t, err)
	assert.Equal(t, int64(5000), updated.Amount)
}

func TestLedgerService_Delete_ByCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	created, _ := svc.Create(1000, "待删除", cat.ID, time.Now(), m1.ID)

	err := svc.Delete(created.ID, m1.ID, "member")
	require.NoError(t, err)

	// Verify soft delete — FindByID should return not found
	_, err = svc.FindByID(created.ID)
	assert.ErrorIs(t, err, ErrLedgerNotFound)
}

func TestLedgerService_Delete_ByNonCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, m2, cat := seedLedgerFixtures()
	created, _ := svc.Create(1000, "待删除", cat.ID, time.Now(), m1.ID)

	err := svc.Delete(created.ID, m2.ID, "member")
	assert.ErrorIs(t, err, ErrLedgerPermissionDenied)
}

func TestLedgerService_Delete_ByAdmin(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	admin := testutil.CreateTestMember(pkg.DB, "admin2", "Admin2", "admin")
	created, _ := svc.Create(1000, "待删除", cat.ID, time.Now(), m1.ID)

	err := svc.Delete(created.ID, admin.ID, "admin")
	require.NoError(t, err)
}

func TestLedgerService_List_ByMonth(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	// Create records in different months
	svc.Create(1000, "5月", cat.ID, time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC), m1.ID) //nolint:errcheck
	svc.Create(2000, "6月", cat.ID, time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), m1.ID) //nolint:errcheck

	result, err := svc.List(repository.LedgerFilter{StartDate: "2026-05-01", EndDate: "2026-06-01", Limit: 20})
	require.NoError(t, err)
	assert.Len(t, result.Groups, 1)
	assert.False(t, result.HasMore)
	assert.Nil(t, result.NextCursor)
}

func TestLedgerService_List_CursorPagination(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	// Create 15 records on different dates within May to trigger cursor pagination
	for i := 1; i <= 15; i++ {
		svc.Create(int64(i*100), fmt.Sprintf("记录%d", i), cat.ID, //nolint:errcheck
			time.Date(2026, 5, i, 12, 0, 0, 0, time.UTC), m1.ID)
	}

	// Page 1: limit=10
	page1, err := svc.List(repository.LedgerFilter{
		StartDate: "2026-05-01", EndDate: "2026-06-01", Limit: 10,
	})
	require.NoError(t, err)
	assert.True(t, page1.HasMore)
	assert.NotNil(t, page1.NextCursor)

	// Count items across all groups on page 1
	totalItemsPage1 := 0
	for _, g := range page1.Groups {
		totalItemsPage1 += len(g.Items)
	}
	assert.Equal(t, 10, totalItemsPage1)

	// Page 2: use cursor from page 1
	cursor, err := repository.DecodeCursor(*page1.NextCursor)
	require.NoError(t, err)

	page2, err := svc.List(repository.LedgerFilter{
		StartDate: "2026-05-01", EndDate: "2026-06-01", Limit: 10, Cursor: cursor,
	})
	require.NoError(t, err)

	totalItemsPage2 := 0
	for _, g := range page2.Groups {
		totalItemsPage2 += len(g.Items)
	}
	assert.Equal(t, 5, totalItemsPage2)
	assert.False(t, page2.HasMore)
	assert.Nil(t, page2.NextCursor)

	// Verify no overlap: collect all IDs from both pages
	seenIDs := make(map[uint]bool)
	for _, g := range page1.Groups {
		for _, item := range g.Items {
			assert.False(t, seenIDs[item.ID], "duplicate ID %d on page 1", item.ID)
			seenIDs[item.ID] = true
		}
	}
	for _, g := range page2.Groups {
		for _, item := range g.Items {
			assert.False(t, seenIDs[item.ID], "duplicate ID %d across pages", item.ID)
			seenIDs[item.ID] = true
		}
	}
	assert.Equal(t, 15, len(seenIDs))
}
