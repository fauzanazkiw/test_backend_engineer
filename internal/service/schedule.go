package service

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"Test/internal/domain"
	"Test/internal/model"
	"Test/internal/repository"
)

const (
	dateLayout = "2006-01-02"
	timeLayout = "15:04:05"
)

// Payload jadwal tayang.
type SchedulePayload struct {
	MovieID     int64
	BranchID    int64
	StudioID    int64
	ShowDate    string
	ShowTime    string
	TicketPrice *int64
}

// ScheduleService menangani business logic jadwal tayang.
type ScheduleService struct {
	schedules *repository.ScheduleRepository
	movies    *repository.MovieRepository
	branches  *repository.BranchRepository
	studios   *repository.StudioRepository
}

func NewScheduleService(
	schedules *repository.ScheduleRepository,
	movies *repository.MovieRepository,
	branches *repository.BranchRepository,
	studios *repository.StudioRepository,
) *ScheduleService {
	return &ScheduleService{
		schedules: schedules,
		movies:    movies,
		branches:  branches,
		studios:   studios,
	}
}

func (s *ScheduleService) List(ctx context.Context, filter model.ScheduleFilter) ([]model.Schedule, error) {
	return s.schedules.List(ctx, filter)
}

func (s *ScheduleService) Get(ctx context.Context, id int64) (*model.Schedule, error) {
	return s.schedules.GetByID(ctx, id)
}

func (s *ScheduleService) Create(ctx context.Context, payload SchedulePayload) (*model.Schedule, error) {
	input, validationErrs, err := s.validate(ctx, payload)
	if err != nil {
		return nil, err
	}
	if !validationErrs.Empty() {
		return nil, validationErrs
	}

	conflict, err := s.schedules.ExistsSlot(ctx, input.StudioID, input.ShowDate, input.ShowTime, 0)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, domain.ErrConflict
	}

	return s.schedules.Create(ctx, input)
}

func (s *ScheduleService) Update(ctx context.Context, id int64, payload SchedulePayload) (*model.Schedule, error) {
	if _, err := s.schedules.GetByID(ctx, id); err != nil {
		return nil, err
	}

	input, validationErrs, err := s.validate(ctx, payload)
	if err != nil {
		return nil, err
	}
	if !validationErrs.Empty() {
		return nil, validationErrs
	}

	conflict, err := s.schedules.ExistsSlot(ctx, input.StudioID, input.ShowDate, input.ShowTime, id)
	if err != nil {
		return nil, err
	}
	if conflict {
		return nil, domain.ErrConflict
	}

	return s.schedules.Update(ctx, id, input)
}

func (s *ScheduleService) Delete(ctx context.Context, id int64) error {
	return s.schedules.Delete(ctx, id)
}

// validate memeriksa seluruh input dan memastikan referensi movie, branch,
// serta studio benar-benar tersedia sebelum data disimpan.
func (s *ScheduleService) validate(ctx context.Context, payload SchedulePayload) (model.ScheduleInput, *domain.ValidationErrors, error) {
	errs := domain.NewValidationErrors() // sekarang *ValidationErrors
	input := model.ScheduleInput{
		MovieID:  payload.MovieID,
		BranchID: payload.BranchID,
		StudioID: payload.StudioID,
	}

	if payload.MovieID <= 0 {
		errs.Add("movie_id", "movie_id is required")
	}
	if payload.BranchID <= 0 {
		errs.Add("branch_id", "branch_id is required")
	}
	if payload.StudioID <= 0 {
		errs.Add("studio_id", "studio_id is required")
	}

	showDate := strings.TrimSpace(payload.ShowDate)
	switch {
	case showDate == "":
		errs.Add("show_date", "show_date is required")
	default:
		parsed, err := time.Parse(dateLayout, showDate)
		if err != nil {
			errs.Add("show_date", "invalid date format, expected YYYY-MM-DD")
		} else {
			input.ShowDate = parsed.Format(dateLayout)
		}
	}

	showTime := strings.TrimSpace(payload.ShowTime)
	switch {
	case showTime == "":
		errs.Add("show_time", "show_time is required")
	default:
		normalized, err := normalizeShowTime(showTime)
		if err != nil {
			errs.Add("show_time", "invalid time format, expected HH:MM or HH:MM:SS")
		} else {
			input.ShowTime = normalized
		}
	}

	switch {
	case payload.TicketPrice == nil:
		errs.Add("ticket_price", "ticket_price is required")
	case *payload.TicketPrice < 0:
		errs.Add("ticket_price", "ticket_price must be greater than or equal to 0")
	default:
		input.TicketPrice = *payload.TicketPrice
	}

	if !errs.Empty() {
		return input, errs, nil
	}

	movieExists, err := s.movies.Exists(ctx, payload.MovieID)
	if err != nil {
		return input, nil, err
	}
	if !movieExists {
		errs.Add("movie_id", fmt.Sprintf("movie with id %d not found", payload.MovieID))
	}

	branchExists, err := s.branches.Exists(ctx, payload.BranchID)
	if err != nil {
		return input, nil, err
	}
	if !branchExists {
		errs.Add("branch_id", fmt.Sprintf("branch with id %d not found", payload.BranchID))
	}

	studioBranchID, err := s.studios.BranchID(ctx, payload.StudioID)
	if errors.Is(err, domain.ErrNotFound) {
		errs.Add("studio_id", fmt.Sprintf("studio with id %d not found", payload.StudioID))
	} else if err != nil {
		return input, nil, err
	} else if studioBranchID != payload.BranchID {
		errs.Add("studio_id", "studio does not belong to the selected branch")
	}

	return input, errs, nil
}

// normalizeShowTime menerima format "HH:MM" maupun "HH:MM:SS" dan
// mengembalikan format standar "HH:MM:SS".
func normalizeShowTime(raw string) (string, error) {
	for _, layout := range []string{timeLayout, "15:04"} {
		parsed, err := time.Parse(layout, raw)
		if err == nil {
			return parsed.Format(timeLayout), nil
		}
	}
	return "", errors.New("invalid time format")
}
