package config

import (
	"fmt"
	"net/url"
	"os"
)

type DBConfig struct {
	Host     string
	Port     string
	User     string
	Password string
	Name     string
	SSLMode  string
	DSN      string // postgres URL with sslrootcert (for reference; connection uses GetFormattedDSN)
}

// GetFormattedDSN returns libpq-style connection string used by the driver (no sslrootcert; matches gold-digger).
func (c *DBConfig) GetFormattedDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode,
	)
}

func LoadAWSConfig() (*DBConfig, error) {
	cfg := &DBConfig{
		Host:     os.Getenv("DB_HOST"),
		Port:     os.Getenv("DB_PORT"),
		User:     os.Getenv("DB_USER"),
		Password: os.Getenv("DB_PASSWORD"),
		Name:     os.Getenv("DB_NAME"),
		SSLMode:  os.Getenv("DB_SSL"),
	}
	if cfg.Port == "" {
		cfg.Port = "5432"
	}
	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}
	// URL form with cert (same as gold-digger); actual connection uses GetFormattedDSN() which has no cert
	sslRootCert := os.Getenv("DB_SSL_ROOT_CERT")
	if sslRootCert == "" {
		sslRootCert = "global-bundle.pem"
	}
	cfg.DSN = fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s&sslrootcert=%s",
		cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Name, cfg.SSLMode, sslRootCert)

	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.Name == "" {
		return nil, fmt.Errorf("incomplete DB config: need DB_HOST, DB_USER, DB_PASSWORD, DB_NAME")
	}

	return cfg, nil
}

func LoadLocalDBConfig() (*DBConfig, error) {
	cfg := &DBConfig{
		Host:     os.Getenv("LOCAL_DB_HOST"),
		Port:     os.Getenv("LOCAL_DB_PORT"),
		User:     os.Getenv("LOCAL_DB_USER"),
		Password: os.Getenv("LOCAL_DB_PASSWORD"),
		Name:     os.Getenv("LOCAL_DB_NAME"),
		SSLMode:  os.Getenv("LOCAL_DB_SSL"),
	}

	if cfg.SSLMode == "" {
		cfg.SSLMode = "disable"
	}

	// Local DSN does not use sslrootcert
	if cfg.Host == "" || cfg.User == "" || cfg.Password == "" || cfg.Name == "" {
		return nil, fmt.Errorf("incomplete LOCAL DB config")
	}

	return cfg, nil
}

// GetLocalDSN returns libpq connection string for local DB (no cert).
func (c *DBConfig) GetLocalDSN() string {
	return fmt.Sprintf(
		"host=%s user=%s password=%s dbname=%s port=%s sslmode=%s",
		c.Host, c.User, c.Password, c.Name, c.Port, c.SSLMode,
	)
}

// GetMigrationDSN returns a postgres URL for golang-migrate (no sslrootcert).
func (c *DBConfig) GetMigrationDSN() string {
	password := url.QueryEscape(c.Password)
	return fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.User, password, c.Host, c.Port, c.Name, c.SSLMode)
}
