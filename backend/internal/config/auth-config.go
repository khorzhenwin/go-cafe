package config

import (
	"fmt"
	"os"
	"time"
)

type AuthConfig struct {
	JWTSecret []byte
	JWTExpiry time.Duration
}

func LoadAuthConfig() (*AuthConfig, error) {
	secret := os.Getenv("JWT_SECRET")
	if secret == "" {
		return nil, fmt.Errorf("JWT_SECRET is required")
	}
	expiryStr := os.Getenv("JWT_EXPIRY")
	if expiryStr == "" {
		expiryStr = "24h"
	}
	expiry, err := time.ParseDuration(expiryStr)
	if err != nil {
		expiry = 24 * time.Hour
	}
	return &AuthConfig{
		JWTSecret: []byte(secret),
		JWTExpiry: expiry,
	}, nil
}
