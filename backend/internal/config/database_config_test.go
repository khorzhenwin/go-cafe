package config

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestLoadAWSConfig_Incomplete(t *testing.T) {
	os.Clearenv()
	_, err := LoadAWSConfig()
	require.Error(t, err)
	assert.Contains(t, err.Error(), "incomplete DB config")
}

func TestLoadAWSConfig_Complete(t *testing.T) {
	os.Clearenv()
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5433")
	os.Setenv("DB_USER", "u")
	os.Setenv("DB_PASSWORD", "p")
	os.Setenv("DB_NAME", "db")
	os.Setenv("DB_SSL", "disable")
	defer os.Clearenv()

	cfg, err := LoadAWSConfig()
	require.NoError(t, err)
	assert.Equal(t, "localhost", cfg.Host)
	assert.Equal(t, "5433", cfg.Port)
	assert.Equal(t, "u", cfg.User)
	assert.Equal(t, "p", cfg.Password)
	assert.Equal(t, "db", cfg.Name)
	assert.Equal(t, "disable", cfg.SSLMode)
	assert.Contains(t, cfg.DSN, "postgres://")
	assert.Contains(t, cfg.DSN, "sslrootcert")
}

func TestGetFormattedDSN_NoCert(t *testing.T) {
	cfg := &DBConfig{Host: "h", Port: "5432", User: "u", Password: "p", Name: "n", SSLMode: "disable"}
	dsn := cfg.GetFormattedDSN()
	assert.Contains(t, dsn, "host=h")
	assert.Contains(t, dsn, "user=u")
	assert.Contains(t, dsn, "dbname=n")
	assert.NotContains(t, dsn, "sslrootcert")
}

func TestGetMigrationDSN_EscapesPassword(t *testing.T) {
	cfg := &DBConfig{Host: "h", Port: "5432", User: "u", Password: "p@ss", Name: "n", SSLMode: "disable"}
	dsn := cfg.GetMigrationDSN()
	assert.Contains(t, dsn, "postgres://")
	assert.Contains(t, dsn, "p%40ss")
}
