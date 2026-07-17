package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestCheckPasswordHash(t *testing.T) {
	password := "password1234!"
	hash, err := HashPassword(password)
	require.NoError(t, err)

	// Test: Good Password
	ok, err := CheckPasswordHash(password, hash)
	require.NoError(t, err)
	assert.True(t, ok)

	// Test: Bad Password
	ok, err = CheckPasswordHash("hello1234!", hash)
	require.NoError(t, err)
	assert.False(t, ok)

	// Test: Bad Hash
	_, err = CheckPasswordHash(password, "BadHash")
	require.Error(t, err)
}
