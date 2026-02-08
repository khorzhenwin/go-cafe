package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAuthConfig_MissingSecret(t *testing.T) {
	os.Clearenv()
	_, err := LoadAuthConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "JWT_SECRET")
}

func TestLoadAuthConfig_Ok(t *testing.T) {
	os.Clearenv()
	os.Setenv("JWT_SECRET", "my-secret-key")
	defer os.Clearenv()

	cfg, err := LoadAuthConfig()
	require.NoError(t, err)
	assert.Equal(t, []byte("my-secret-key"), cfg.JWTSecret)
	assert.NotZero(t, cfg.JWTExpiry)
}
