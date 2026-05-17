package pkg

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func init() {
	InitJWT("test-secret-key-for-testing-purposes")
}

func TestGenerateToken_Success(t *testing.T) {
	token, err := GenerateToken(1, "testuser", "admin")
	require.NoError(t, err)
	assert.NotEmpty(t, token)
	assert.Contains(t, token, ".")
}

func TestParseToken_Success(t *testing.T) {
	token, _ := GenerateToken(42, "alice", "member")
	claims, err := ParseToken(token)
	require.NoError(t, err)
	assert.Equal(t, uint(42), claims.MemberID)
	assert.Equal(t, "alice", claims.Username)
	assert.Equal(t, "member", claims.Role)
}

func TestParseToken_InvalidFormat(t *testing.T) {
	_, err := ParseToken("not.a.valid.token")
	assert.Error(t, err)
}

func TestParseToken_EmptyString(t *testing.T) {
	_, err := ParseToken("")
	assert.Error(t, err)
}

func TestGenerateToken_RolePreserved(t *testing.T) {
	token, err := GenerateToken(99, "bob", "member")
	require.NoError(t, err)
	claims, _ := ParseToken(token)
	assert.Equal(t, "member", claims.Role)
	assert.Equal(t, uint(99), claims.MemberID)
}

func TestTokenExpiry_SetTo7Days(t *testing.T) {
	token, _ := GenerateToken(1, "user", "member")
	claims, _ := ParseToken(token)
	exp := claims.ExpiresAt.Time
	assert.True(t, exp.After(time.Now()))
	assert.True(t, exp.Before(time.Now().Add(8*24*time.Hour)))
}
