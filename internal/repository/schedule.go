package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"Test/internal/domain"
	"Test/internal/model"

	"github.com/Masterminds/squirrel"
	"github.com/jackc/pgx/v5/pgconn"
)

// uniqueViolationCode adalah kode error PostgreSQL untuk pelanggaran unique constraint.
const uniqueViolationCode = "23505"

// scheduleColumns adalah kolom schedule yang dipakai pada semua query.
// show_date dan show_time diformat menjadi string agar response konsisten
// dalam bentuk "YYYY-MM-DD" dan "HH:MM:SS".
var scheduleColumns = []string{
	"id",
	"movie_id",
	"branch_id",
	"studio_id",
	"to_char(show_date, 'YYYY-MM-DD') AS show_date",
	"show_time::text AS show_time",
	"ticket_price",
	"created_at",
	"updated_at",
}

// ScheduleRepository menangani akses data jadwal tayang.
type ScheduleRepository struct {
	db *sql.DB
}

func NewScheduleRepository(db *sql.DB) *ScheduleRepository {
	return &ScheduleRepository{db: db}
}

// List mengambil daftar schedule sesuai filter.
func (r *ScheduleRepository) List(ctx context.Context, filter model.ScheduleFilter) ([]model.Schedule, error) {
	query := builder.Select(scheduleColumns...).From("schedules")

	if filter.MovieID != nil {
		query = query.Where("movie_id = ?", *filter.MovieID)
	}
	if filter.BranchID != nil {
		query = query.Where("branch_id = ?", *filter.BranchID)
	}
	if filter.StudioID != nil {
		query = query.Where("studio_id = ?", *filter.StudioID)
	}
	if filter.ShowDate != "" {
		query = query.Where("show_date = ?::date", filter.ShowDate)
	}

	query = query.OrderBy("show_date ASC", "show_time ASC", "id ASC")

	sqlQuery, args, err := query.ToSql()
	if err != nil {
		return nil, fmt.Errorf("gagal membangun query schedule: %w", err)
	}

	rows, err := r.db.QueryContext(ctx, sqlQuery, args...)
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil daftar schedule: %w", err)
	}
	defer rows.Close()

	schedules := make([]model.Schedule, 0)
	for rows.Next() {
		schedule, err := scanSchedule(rows)
		if err != nil {
			return nil, fmt.Errorf("gagal membaca data schedule: %w", err)
		}
		schedules = append(schedules, schedule)
	}
	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("gagal iterasi data schedule: %w", err)
	}

	return schedules, nil
}

// GetByID mengambil satu schedule berdasarkan id.
func (r *ScheduleRepository) GetByID(ctx context.Context, id int64) (*model.Schedule, error) {
	query, args, err := builder.Select(scheduleColumns...).From("schedules").
		Where("id = ?", id).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("gagal membangun query schedule: %w", err)
	}

	schedule, err := scanSchedule(r.db.QueryRowContext(ctx, query, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, fmt.Errorf("gagal mengambil schedule: %w", err)
	}

	return &schedule, nil
}

// Create menyimpan schedule baru dan mengembalikan data yang tersimpan.
func (r *ScheduleRepository) Create(ctx context.Context, input model.ScheduleInput) (*model.Schedule, error) {
	query, args, err := builder.Insert("schedules").
		Columns("movie_id", "branch_id", "studio_id", "show_date", "show_time", "ticket_price").
		Values(
			input.MovieID,
			input.BranchID,
			input.StudioID,
			squirrel.Expr("?::date", input.ShowDate),
			squirrel.Expr("?::time", input.ShowTime),
			input.TicketPrice,
		).
		Suffix("RETURNING " + strings.Join(scheduleColumns, ", ")).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("gagal membangun query insert schedule: %w", err)
	}

	schedule, err := scanSchedule(r.db.QueryRowContext(ctx, query, args...))
	if err != nil {
		return nil, mapScheduleWriteError(err)
	}

	return &schedule, nil
}

// Update mengubah schedule berdasarkan id dan mengembalikan data terbaru.
func (r *ScheduleRepository) Update(ctx context.Context, id int64, input model.ScheduleInput) (*model.Schedule, error) {
	query, args, err := builder.Update("schedules").
		Set("movie_id", input.MovieID).
		Set("branch_id", input.BranchID).
		Set("studio_id", input.StudioID).
		Set("show_date", squirrel.Expr("?::date", input.ShowDate)).
		Set("show_time", squirrel.Expr("?::time", input.ShowTime)).
		Set("ticket_price", input.TicketPrice).
		Where("id = ?", id).
		Suffix("RETURNING " + strings.Join(scheduleColumns, ", ")).
		ToSql()
	if err != nil {
		return nil, fmt.Errorf("gagal membangun query update schedule: %w", err)
	}

	schedule, err := scanSchedule(r.db.QueryRowContext(ctx, query, args...))
	if errors.Is(err, sql.ErrNoRows) {
		return nil, domain.ErrNotFound
	}
	if err != nil {
		return nil, mapScheduleWriteError(err)
	}

	return &schedule, nil
}

// Delete menghapus schedule berdasarkan id.
func (r *ScheduleRepository) Delete(ctx context.Context, id int64) error {
	query, args, err := builder.Delete("schedules").Where("id = ?", id).ToSql()
	if err != nil {
		return fmt.Errorf("gagal membangun query delete schedule: %w", err)
	}

	result, err := r.db.ExecContext(ctx, query, args...)
	if err != nil {
		return fmt.Errorf("gagal menghapus schedule: %w", err)
	}

	affected, err := result.RowsAffected()
	if err != nil {
		return fmt.Errorf("gagal menghapus schedule: %w", err)
	}
	if affected == 0 {
		return domain.ErrNotFound
	}

	return nil
}

// ExistsSlot memeriksa apakah sudah ada schedule lain pada studio, tanggal,
// dan jam yang sama. excludeID dipakai saat update agar data itu sendiri tidak dihitung.
func (r *ScheduleRepository) ExistsSlot(ctx context.Context, studioID int64, showDate, showTime string, excludeID int64) (bool, error) {
	query := builder.Select("1").From("schedules").
		Where("studio_id = ?", studioID).
		Where("show_date = ?::date", showDate).
		Where("show_time = ?::time", showTime)

	if excludeID > 0 {
		query = query.Where("id <> ?", excludeID)
	}

	var found int
	err := query.RunWith(r.db).QueryRowContext(ctx).Scan(&found)
	if errors.Is(err, sql.ErrNoRows) {
		return false, nil
	}
	if err != nil {
		return false, fmt.Errorf("gagal memeriksa jadwal bentrok: %w", err)
	}

	return true, nil
}

type rowScanner interface {
	Scan(dest ...any) error
}

func scanSchedule(row rowScanner) (model.Schedule, error) {
	var schedule model.Schedule
	err := row.Scan(
		&schedule.ID,
		&schedule.MovieID,
		&schedule.BranchID,
		&schedule.StudioID,
		&schedule.ShowDate,
		&schedule.ShowTime,
		&schedule.TicketPrice,
		&schedule.CreatedAt,
		&schedule.UpdatedAt,
	)
	return schedule, err
}

// mapScheduleWriteError menerjemahkan error constraint database menjadi error domain.
func mapScheduleWriteError(err error) error {
	var pgErr *pgconn.PgError
	if errors.As(err, &pgErr) && pgErr.Code == uniqueViolationCode {
		return domain.ErrConflict
	}
	return fmt.Errorf("gagal menyimpan schedule: %w", err)
}
