package pkg

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_Success(t *testing.T) {
	hash, err := HashPassword("test123")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.Contains(t, hash, "$2a$")
}

func TestHashPassword_EmptyPassword(t *testing.T) {
	hash, err := HashPassword("")
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
}

func TestCheckPassword_Match(t *testing.T) {
	hash, _ := HashPassword("secure456")
	assert.True(t, CheckPassword(hash, "secure456"))
}

func TestCheckPassword_NoMatch(t *testing.T) {
	hash, _ := HashPassword("secure456")
	assert.False(t, CheckPassword(hash, "wrongpass"))
}

func TestCheckPassword_InvalidHash(t *testing.T) {
	assert.False(t, CheckPassword("not-a-valid-hash", "anything"))
}

func TestDefaultPassword_Value(t *testing.T) {
	assert.Equal(t, "home123", DefaultPassword)
}
