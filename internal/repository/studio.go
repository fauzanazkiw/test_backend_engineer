package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"

	"Test/internal/domain"
)

// StudioRepository menangani pengecekan data studio.
type StudioRepository struct {
	db *sql.DB
}

func NewStudioRepository(db *sql.DB) *StudioRepository {
	return &StudioRepository{db: db}
}

// BranchID mengembalikan branch_id dari studio, atau domain.ErrNotFound bila studio tidak ada.
func (r *StudioRepository) BranchID(ctx context.Context, id int64) (int64, error) {
	var branchID int64
	err := builder.Select("branch_id").From("studios").Where("id = ?", id).
		RunWith(r.db).QueryRowContext(ctx).Scan(&branchID)
	if errors.Is(err, sql.ErrNoRows) {
		return 0, domain.ErrNotFound
	}
	if err != nil {
		return 0, fmt.Errorf("gagal mengambil studio: %w", err)
	}
	return branchID, nil
}
