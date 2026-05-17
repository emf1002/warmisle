package service

import (
	"testing"

	"home-center/internal/model"
	"home-center/internal/pkg"
	"home-center/internal/testutil"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func setupAuthTest() (*AuthService, func()) {
	_ = testutil.SetupTestDB()
	pkg.InitJWT("test-jwt-secret-for-auth-tests")
	svc := NewAuthService()
	return svc, func() { testutil.TeardownTestDB() }
}

func TestAuthService_InitCheck_NeedInit(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	needInit, err := svc.InitCheck()
	require.NoError(t, err)
	assert.True(t, needInit, "empty database should need init")
}

func TestAuthService_InitCheck_NoNeedInit(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("pass123")
	pkg.DB.Create(&model.Member{
		Username: "existing", Password: hash, Name: "Exist", Role: "admin", Status: "active",
	})

	needInit, err := svc.InitCheck()
	require.NoError(t, err)
	assert.False(t, needInit, "database with members should not need init")
}

func TestAuthService_Login_Success(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("correctpw")
	pkg.DB.Create(&model.Member{
		Username: "alice", Password: hash, Name: "Alice", Role: "member", Status: "active",
	})

	token, err := svc.Login("alice", "correctpw")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
}

func TestAuthService_Login_WrongPassword(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("correctpw")
	pkg.DB.Create(&model.Member{
		Username: "bob", Password: hash, Name: "Bob", Role: "member", Status: "active",
	})

	_, err := svc.Login("bob", "wrongpw")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_UserNotFound(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	_, err := svc.Login("nonexistent", "any")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
}

func TestAuthService_Login_DisabledMember(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("pass123")
	pkg.DB.Create(&model.Member{
		Username: "disabled_user", Password: hash, Name: "Disabled", Role: "member", Status: "disabled",
	})

	_, err := svc.Login("disabled_user", "pass123")
	assert.ErrorIs(t, err, ErrInvalidCredentials, "disabled member should not login")
}

func TestAuthService_Login_LockAfter5Failures(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("secret")
	pkg.DB.Create(&model.Member{
		Username: "locktest", Password: hash, Name: "LockTest", Role: "member", Status: "active",
	})

	// 5 consecutive wrong passwords
	for i := 0; i < 5; i++ {
		_, err := svc.Login("locktest", "wrong")
		assert.ErrorIs(t, err, ErrInvalidCredentials)
	}

	// 6th attempt should be locked
	_, err := svc.Login("locktest", "secret")
	assert.ErrorIs(t, err, ErrAccountLocked, "should be locked after 5 failures")
}

func TestAuthService_Login_ResetAfterSuccess(t *testing.T) {
	svc, teardown := setupAuthTest()
	defer teardown()

	hash, _ := pkg.HashPassword("goodpw")
	pkg.DB.Create(&model.Member{
		Username: "reset_test", Password: hash, Name: "ResetTest", Role: "member", Status: "active",
	})

	// 4 failures
	for i := 0; i < 4; i++ {
		svc.Login("reset_test", "wrong")
	}

	// 5th attempt succeeds — should reset counter
	token, err := svc.Login("reset_test", "goodpw")
	require.NoError(t, err)
	assert.NotEmpty(t, token)

	// After success, another failure should NOT immediately lock
	_, err = svc.Login("reset_test", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials)
	// Confirm not locked (only 1 failure after reset)
	_, err = svc.Login("reset_test", "wrong")
	assert.ErrorIs(t, err, ErrInvalidCredentials, "should not be locked after only 2 failures")
}
