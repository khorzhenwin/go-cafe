package auth

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestHashPassword_and_CheckPassword(t *testing.T) {
	password := "secret123"
	hash, err := HashPassword(password)
	require.NoError(t, err)
	assert.NotEmpty(t, hash)
	assert.NotEqual(t, password, hash)

	assert.True(t, CheckPassword(password, hash))
	assert.False(t, CheckPassword("wrong", hash))
	assert.False(t, CheckPassword(password, "not-a-hash"))
}

func TestHashPassword_DifferentEachTime(t *testing.T) {
	h1, _ := HashPassword("same")
	h2, _ := HashPassword("same")
	assert.NotEqual(t, h1, h2)
	// but both verify
	assert.True(t, CheckPassword("same", h1))
	assert.True(t, CheckPassword("same", h2))
}
