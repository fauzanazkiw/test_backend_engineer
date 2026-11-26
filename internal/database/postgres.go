package database

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Driver PostgreSQL (pgx) untuk database/sql.
	_ "github.com/jackc/pgx/v5/stdlib"
)

// Connect membuka koneksi ke PostgreSQL dan memastikan koneksi benar-benar tersedia.
func Connect(ctx context.Context, databaseURL string, maxConns int) (*sql.DB, error) {
	db, err := sql.Open("pgx", databaseURL)
	if err != nil {
		return nil, fmt.Errorf("gagal membuka koneksi database: %w", err)
	}

	db.SetMaxOpenConns(maxConns)
	db.SetMaxIdleConns(maxConns / 2)
	db.SetConnMaxLifetime(30 * time.Minute)
	db.SetConnMaxIdleTime(5 * time.Minute)

	pingCtx, cancel := context.WithTimeout(ctx, 5*time.Second)
	defer cancel()

	if err := db.PingContext(pingCtx); err != nil {
		db.Close()
		return nil, fmt.Errorf("gagal terhubung ke database: %w", err)
	}

	return db, nil
}
