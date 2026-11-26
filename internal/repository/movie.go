package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
)

// MovieRepository menangani pengecekan data film.
type MovieRepository struct {
	db *sql.DB
}

func NewMovieRepository(db *sql.DB) *MovieRepository {
	return &MovieRepository{db: db}
}

// Exists memeriksa apakah movie dengan id tertentu tersedia.
func (r *MovieRepository) Exists(ctx context.Context, id int64) (bool, error) {
	var found int
	err := builder.Select("1").From("movies").Where("id = ?", id).
		RunWith(r.db).QueryRowContext(ctx).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("gagal memeriksa movie: %w", err)
	}
	return true, nil
}
