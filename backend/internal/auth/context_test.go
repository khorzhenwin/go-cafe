package auth

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUserIDFromContext_Empty(t *testing.T) {
	_, ok := UserIDFromContext(context.Background())
	assert.False(t, ok)
}

func TestUserIDFromContext_Set(t *testing.T) {
	ctx := context.WithValue(context.Background(), UserIDKey, uint(42))
	id, ok := UserIDFromContext(ctx)
	assert.True(t, ok)
	assert.Equal(t, uint(42), id)
}
