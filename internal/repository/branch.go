package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// BranchRepository menangani pengecekan data cabang bioskop.
type BranchRepository struct {
	db *sql.DB
}

func NewBranchRepository(db *sql.DB) *BranchRepository {
	return &BranchRepository{db: db}
}

// Exists memeriksa apakah branch dengan id tertentu tersedia.
func (r *BranchRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var found int
	err := builder.Select("1").From("branches").Where("id = ?", id).
		RunWith(r.db).QueryRowContext(ctx).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("gagal memeriksa branch: %w", err)
	}
	return true, nil
}
