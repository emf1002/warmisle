package service

import (
	"testing"

	"warmisle/internal/model"
	"warmisle/internal/pkg"
	"warmisle/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupMemberTest() (*MemberService, func()) {
	testutil.SetupTestDB()
	pkg.InitJWT("test-secret")
	svc := NewMemberService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestMemberService_Create_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, err := svc.Create("newuser", "password1", "New User", "👩", "member")
	require.NoError(t, err)
	assert.Equal(t, "newuser", member.Username)
	assert.Equal(t, "New User", member.Name)
	assert.Equal(t, "member", member.Role)
	assert.Equal(t, "active", member.Status)
	assert.NotEmpty(t, member.Password)
}

func TestMemberService_Create_DuplicateUsername(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	svc.Create("dup", "pass1234", "First", "", "") //nolint:errcheck
	_, err := svc.Create("dup", "pass5678", "Second", "", "")
	assert.ErrorIs(t, err, ErrUsernameTaken)
}

func TestMemberService_Create_InvalidUsername(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	_, err := svc.Create("ab", "pass1234", "Short", "", "")
	assert.ErrorIs(t, err, ErrInvalidUsername, "username too short")

	_, err = svc.Create("user@name", "pass1234", "Bad", "", "")
	assert.ErrorIs(t, err, ErrInvalidUsername, "username has invalid chars")
}

func TestMemberService_Create_InvalidPassword(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	_, err := svc.Create("validuser", "12345", "Valid", "", "")
	assert.ErrorIs(t, err, ErrInvalidPassword, "password too short")
}

func TestMemberService_Create_DefaultValues(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, err := svc.Create("minuser", "password1", "", "", "")
	require.NoError(t, err)
	assert.Equal(t, "👨", member.Avatar, "default avatar when empty")
	assert.Equal(t, "member", member.Role, "default role when empty")
}

func TestMemberService_Update_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	created, _ := svc.Create("edit_user", "password1", "Old Name", "👨", "member")

	updated, err := svc.Update(created.ID, "New Name", "👩", "member")
	require.NoError(t, err)
	assert.Equal(t, "New Name", updated.Name)
	assert.Equal(t, "👩", updated.Avatar)
}

func TestMemberService_Update_CannotRemoveLastAdmin(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	admin, _ := svc.Create("onlyadmin", "password1", "Admin", "👨", "admin")

	_, err := svc.Update(admin.ID, "", "", "member")
	assert.ErrorIs(t, err, ErrCannotDeleteLastAdmin)
}

func TestMemberService_Delete_WithActivityRecords(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	// Create a member with activity records (ledger entry)
	member, _ := svc.Create("active_member", "password1", "Active", "👨", "member")
	pkg.DB.Create(&model.Ledger{
		Amount: 1000, CategoryID: 1, CreatorID: member.ID,
	})

	err := svc.Delete(member.ID, 1)
	assert.ErrorIs(t, err, ErrHasActivityRecords)
}

func TestMemberService_Disable_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	// Create a second member as the operator, so operator ID != target member ID
	svc.Create("operator", "password1", "Operator", "👨", "member") //nolint:errcheck
	member, _ := svc.Create("todisable", "password1", "ToDisable", "👨", "member")

	err := svc.Disable(member.ID, 1)
	require.NoError(t, err)

	updated, _ := svc.GetProfile(member.ID)
	assert.Equal(t, "disabled", updated.Status)
}

func TestMemberService_Disable_CannotDisableSelf(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("self_disable", "password1", "Self", "👨", "admin")

	err := svc.Disable(member.ID, member.ID)
	assert.ErrorIs(t, err, ErrCannotDisableSelf)
}

func TestMemberService_Disable_CannotDisableLastAdmin(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	admin, _ := svc.Create("lastadmin", "password1", "LastAdmin", "👨", "admin")

	err := svc.Disable(admin.ID, 2)
	assert.ErrorIs(t, err, ErrCannotDeleteLastAdmin)
}

func TestMemberService_Enable_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("toenable", "password1", "ToEnable", "👨", "member")
	_ = svc.Disable(member.ID, 1)

	enabled, err := svc.Enable(member.ID)
	require.NoError(t, err)
	assert.Equal(t, "active", enabled.Status)
}

func TestMemberService_ResetPassword(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("resetpwd", "oldpassword", "Reset", "👨", "member")

	err := svc.ResetPassword(member.ID)
	require.NoError(t, err)

	updated, _ := svc.GetProfile(member.ID)
	assert.True(t, pkg.CheckPassword(updated.Password, pkg.DefaultPassword))
}

func TestMemberService_ChangePassword_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("changepwd", "oldpassword", "Change", "👨", "member")

	err := svc.ChangePassword(member.ID, "oldpassword", "newpassword")
	require.NoError(t, err)

	updated, _ := svc.GetProfile(member.ID)
	assert.True(t, pkg.CheckPassword(updated.Password, "newpassword"))
}

func TestMemberService_ChangePassword_WrongOldPassword(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("wrongold", "correctold", "WrongOld", "👨", "member")

	err := svc.ChangePassword(member.ID, "wrongold", "newpassword")
	assert.ErrorIs(t, err, ErrIncorrectPassword)
}

func TestMemberService_List(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	svc.Create("user1", "password1", "User1", "", "member") //nolint:errcheck
	svc.Create("user2", "password1", "User2", "", "admin") //nolint:errcheck

	list, err := svc.List()
	require.NoError(t, err)
	assert.Len(t, list, 2)
}

func TestMemberService_UpdateProfile_Success(t *testing.T) {
	svc, teardown := setupMemberTest()
	defer teardown()

	member, _ := svc.Create("prof_user", "password1", "Old", "👨", "member")

	updated, err := svc.UpdateProfile(member.ID, "NewName", "👩")
	require.NoError(t, err)
	assert.Equal(t, "NewName", updated.Name)
	assert.Equal(t, "👩", updated.Avatar)
}
