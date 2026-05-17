package service

import (
	"testing"
	"time"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/repository"
	"home-center/internal/testutil"

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

	m1, m2, cat := seedLedgerFixtures()

	result, err := svc.Create(3550, "午餐", cat.ID, []uint{m1.ID, m2.ID}, time.Now(), m1.ID)
	require.NoError(t, err)
	assert.Equal(t, int64(3550), result.Amount)
	assert.Equal(t, "午餐", result.Note)
	assert.Equal(t, m1.ID, result.CreatorID)
	assert.Len(t, result.Members, 2)
}

func TestLedgerService_Create_ZeroAmount(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	_, err := svc.Create(0, "免费", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestLedgerService_Create_NegativeAmount(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	_, err := svc.Create(-100, "负数", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrInvalidAmount)
}

func TestLedgerService_Create_NoMembers(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	_, _, cat := seedLedgerFixtures()

	_, err := svc.Create(1000, "无成员", cat.ID, []uint{}, time.Now(), 1)
	assert.ErrorIs(t, err, ErrNoMembers)
}

func TestLedgerService_Create_CategoryNotFound(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, _ := seedLedgerFixtures()

	_, err := svc.Create(1000, "不存在的分类", 99999, []uint{m1.ID}, time.Now(), m1.ID)
	assert.ErrorIs(t, err, ErrLedgerCategoryNotFound)
}

func TestLedgerService_Update_ByCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	created, _ := svc.Create(2000, "原始", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	newAmount := int64(3000)
	newNote := "修改后"
	updated, err := svc.Update(created.ID, &newAmount, &newNote, nil, nil, nil, m1.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, int64(3000), updated.Amount)
	assert.Equal(t, "修改后", updated.Note)
}

func TestLedgerService_Update_ByNonCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, m2, cat := seedLedgerFixtures()
	created, _ := svc.Create(2000, "原始", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	newAmount := int64(3000)
	_, err := svc.Update(created.ID, &newAmount, nil, nil, nil, nil, m2.ID, "member")
	assert.ErrorIs(t, err, ErrLedgerPermissionDenied)
}

func TestLedgerService_Update_ByAdmin(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	admin := testutil.CreateTestMember(pkg.DB, "admin", "Admin", "admin")
	created, _ := svc.Create(2000, "原始", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	newAmount := int64(5000)
	updated, err := svc.Update(created.ID, &newAmount, nil, nil, nil, nil, admin.ID, "admin")
	require.NoError(t, err)
	assert.Equal(t, int64(5000), updated.Amount)
}

func TestLedgerService_Delete_ByCreator(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	created, _ := svc.Create(1000, "待删除", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

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
	created, _ := svc.Create(1000, "待删除", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	err := svc.Delete(created.ID, m2.ID, "member")
	assert.ErrorIs(t, err, ErrLedgerPermissionDenied)
}

func TestLedgerService_Delete_ByAdmin(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()
	admin := testutil.CreateTestMember(pkg.DB, "admin2", "Admin2", "admin")
	created, _ := svc.Create(1000, "待删除", cat.ID, []uint{m1.ID}, time.Now(), m1.ID)

	err := svc.Delete(created.ID, admin.ID, "admin")
	require.NoError(t, err)
}

func TestLedgerService_List_ByMonth(t *testing.T) {
	svc, teardown := setupLedgerTest()
	defer teardown()

	m1, _, cat := seedLedgerFixtures()

	// Create records in different months
	svc.Create(1000, "5月", cat.ID, []uint{m1.ID}, time.Date(2026, 5, 15, 0, 0, 0, 0, time.UTC), m1.ID)
	svc.Create(2000, "6月", cat.ID, []uint{m1.ID}, time.Date(2026, 6, 10, 0, 0, 0, 0, time.UTC), m1.ID)

	result, err := svc.List(repository.LedgerFilter{Month: "2026-05", Page: 1, PageSize: 20})
	require.NoError(t, err)
	assert.Equal(t, int64(1), result.Total)
}
