package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"Test/internal/domain"
	"Test/internal/model"
)

// UserRepository menangani akses data user.
type UserRepository struct {
	db *sql.DB
}

func NewUserRepository(db *sql.DB) *UserRepository {
	return &UserRepository{db: db}
}

// GetByEmail mencari user berdasarkan email. Email dinormalisasi menjadi huruf kecil.
func (r *UserRepository) GetByEmail(ctx context.Context, email string) (*model.User, error) {
	query, args, err := builder.
		Select("id", "name", "email", "password_hash", "created_at", "updated_at").
		From("users").
		Where("LOWER(email) = ?", strings.ToLower(strings.TrimSpace(email))).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("gagal membangun query user: %w", err)
	}

	var user model.User
	err = r.db.QueryRowContext(ctx, query, args...).Scan(
		&user.ID,
		&user.Name,
		&user.Email,
		&user.Password,
		&user.CreatedAt,
		&user.UpdatedAt,
	)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil user: %w", err)
	}

	return &user, nil
}
