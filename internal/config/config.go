package config

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"time"

	"github.com/joho/godotenv"
)

const (
	defaultAppPort    = "8080"
	defaultJWTExpiry  = 24 * time.Hour
	defaultDBMaxConns = 10
)

// Config menyimpan seluruh konfigurasi aplikasi yang berasal dari environment variable.
type Config struct {
	AppPort     string
	DatabaseURL string
	JWTSecret   string
	JWTExpiry   time.Duration
	DBMaxConns  int
}

// Load membaca konfigurasi dari environment variable.
// File .env bersifat opsional dan hanya dipakai untuk memudahkan development.
func Load() (*Config, error) {
	_ = godotenv.Load()

	cfg := &Config{
		AppPort:     getEnv("APP_PORT", defaultAppPort),
		DatabaseURL: os.Getenv("DATABASE_URL"),
		JWTSecret:   os.Getenv("JWT_SECRET"),
		JWTExpiry:   defaultJWTExpiry,
		DBMaxConns:  defaultDBMaxConns,
	}

	if cfg.DatabaseURL == "" {
		return nil, errors.New("DATABASE_URL wajib diisi")
	}
	if cfg.JWTSecret == "" {
		return nil, errors.New("JWT_SECRET wajib diisi")
	}

	if raw := os.Getenv("JWT_EXPIRES_IN"); raw != "" {
		expiry, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("JWT_EXPIRES_IN tidak valid: %w", err)
		}
		if expiry <= 0 {
			return nil, errors.New("JWT_EXPIRES_IN harus lebih besar dari 0")
		}
		cfg.JWTExpiry = expiry
	}

	if raw := os.Getenv("DB_MAX_CONNS"); raw != "" {
		maxConns, err := strconv.Atoi(raw)
		if err != nil || maxConns <= 0 {
			return nil, errors.New("DB_MAX_CONNS harus berupa angka positif")
		}
		cfg.DBMaxConns = maxConns
	}

	return cfg, nil
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
