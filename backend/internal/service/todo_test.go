package service

import (
	"testing"

	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupTodoTest() (*TodoService, func()) {
	testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewTodoService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestTodoService_Create_Success(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator", "Creator", "member")

	todo, err := svc.Create("买菜", "去超市买菜", "important", nil, nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "买菜", todo.Title)
	assert.Equal(t, "important", todo.Priority)
	assert.Equal(t, "pending", todo.Status)
	assert.Equal(t, creator.ID, todo.CreatorID)
}

func TestTodoService_Create_TitleRequired(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator2", "Creator2", "member")

	_, err := svc.Create("", "描述", "normal", nil, nil, creator.ID)
	assert.ErrorIs(t, err, ErrTodoTitleRequired)
}

func TestTodoService_Create_InvalidPriority(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator3", "Creator3", "member")

	_, err := svc.Create("测试", "", "super-urgent", nil, nil, creator.ID)
	assert.ErrorIs(t, err, ErrTodoInvalidPriority)
}

func TestTodoService_Create_DefaultPriority(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator4", "Creator4", "member")

	todo, err := svc.Create("默认优先级", "", "", nil, nil, creator.ID)
	require.NoError(t, err)
	assert.Equal(t, "normal", todo.Priority)
}

func TestTodoService_Create_WithAssignee(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator5", "Creator5", "member")
	assignee := testutil.CreateTestMember(pkg.DB, "assignee", "Assignee", "member")
	assigneeID := assignee.ID

	todo, err := svc.Create("指派任务", "分配给某人", "urgent", &assigneeID, nil, creator.ID)
	require.NoError(t, err)
	require.NotNil(t, todo.Assignee)
	assert.Equal(t, assignee.ID, todo.Assignee.ID)
}

func TestTodoService_Toggle_CompleteAndUncomplete(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator6", "Creator6", "member")

	todo, _ := svc.Create("完成测试", "", "normal", nil, nil, creator.ID)

	// Complete
	completed, err := svc.Toggle(todo.ID, creator.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, "completed", completed.Status)
	assert.NotNil(t, completed.CompletedAt)

	// Uncomplete
	uncompleted, err := svc.Toggle(todo.ID, creator.ID, "member")
	require.NoError(t, err)
	assert.Equal(t, "pending", uncompleted.Status)
	assert.Nil(t, uncompleted.CompletedAt)
}

func TestTodoService_Toggle_PermissionDenied(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator7", "Creator7", "member")
	other := testutil.CreateTestMember(pkg.DB, "other7", "Other7", "member")

	todo, _ := svc.Create("权限测试", "", "normal", nil, nil, creator.ID)

	_, err := svc.Toggle(todo.ID, other.ID, "member")
	assert.ErrorIs(t, err, ErrTodoPermissionDenied)
}

func TestTodoService_Claim_Success(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator8", "Creator8", "member")
	claimer := testutil.CreateTestMember(pkg.DB, "claimer", "Claimer", "member")

	todo, _ := svc.Create("认领测试", "", "normal", nil, nil, creator.ID)

	claimed, err := svc.Claim(todo.ID, claimer.ID)
	require.NoError(t, err)
	require.NotNil(t, claimed.Assignee)
	assert.Equal(t, claimer.ID, claimed.Assignee.ID)
}

func TestTodoService_Claim_AlreadyAssigned(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator9", "Creator9", "member")
	assignee := testutil.CreateTestMember(pkg.DB, "assignee9", "Assignee9", "member")
	claimer := testutil.CreateTestMember(pkg.DB, "claimer9", "Claimer9", "member")
	assigneeID := assignee.ID

	todo, _ := svc.Create("已指派", "", "normal", &assigneeID, nil, creator.ID)

	_, err := svc.Claim(todo.ID, claimer.ID)
	assert.ErrorIs(t, err, ErrTodoAlreadyAssigned)
}

func TestTodoService_Delete_ByCreator(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator10", "Creator10", "member")

	todo, _ := svc.Create("删除测试", "", "normal", nil, nil, creator.ID)

	err := svc.Delete(todo.ID, creator.ID, "member")
	require.NoError(t, err)

	_, err = svc.FindByID(todo.ID)
	assert.ErrorIs(t, err, ErrTodoNotFound)
}

func TestTodoService_Delete_ByNonCreator(t *testing.T) {
	svc, teardown := setupTodoTest()
	defer teardown()

	creator := testutil.CreateTestMember(pkg.DB, "creator11", "Creator11", "member")
	other := testutil.CreateTestMember(pkg.DB, "other11", "Other11", "member")

	todo, _ := svc.Create("删除权限测试", "", "normal", nil, nil, creator.ID)

	err := svc.Delete(todo.ID, other.ID, "member")
	assert.ErrorIs(t, err, ErrTodoPermissionDenied)
}
