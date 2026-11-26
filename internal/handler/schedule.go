package handler

import (
	"errors"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"Test/internal/domain"
	"Test/internal/model"
	"Test/internal/response"
	"Test/internal/service"

	"github.com/go-chi/chi/v5"
)

// ScheduleHandler menangani endpoint CRUD jadwal tayang.
type ScheduleHandler struct {
	schedules *service.ScheduleService
}

func NewScheduleHandler(scheduleService *service.ScheduleService) *ScheduleHandler {
	return &ScheduleHandler{schedules: scheduleService}
}

type scheduleRequest struct {
	MovieID     *int64 `json:"movie_id"`
	BranchID    *int64 `json:"branch_id"`
	StudioID    *int64 `json:"studio_id"`
	ShowDate    string `json:"show_date"`
	ShowTime    string `json:"show_time"`
	TicketPrice *int64 `json:"ticket_price"`
}

func (req scheduleRequest) toPayload() service.SchedulePayload {
	payload := service.SchedulePayload{
		ShowDate:    req.ShowDate,
		ShowTime:    req.ShowTime,
		TicketPrice: req.TicketPrice,
	}
	if req.MovieID != nil {
		payload.MovieID = *req.MovieID
	}
	if req.BranchID != nil {
		payload.BranchID = *req.BranchID
	}
	if req.StudioID != nil {
		payload.StudioID = *req.StudioID
	}
	return payload
}

func (h *ScheduleHandler) List(w http.ResponseWriter, r *http.Request) {
	filter, err := parseScheduleFilter(r)
	if err != nil {
		response.Error(w, http.StatusBadRequest, err.Error())
		return
	}

	schedules, err := h.schedules.List(r.Context(), filter)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.Data(w, http.StatusOK, schedules)
}

func (h *ScheduleHandler) Get(w http.ResponseWriter, r *http.Request) {
	id, ok := parseScheduleID(w, r)
	if !ok {
		return
	}

	schedule, err := h.schedules.Get(r.Context(), id)
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.Data(w, http.StatusOK, schedule)
}

func (h *ScheduleHandler) Create(w http.ResponseWriter, r *http.Request) {
	var req scheduleRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	schedule, err := h.schedules.Create(r.Context(), req.toPayload())
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.MessageData(w, http.StatusCreated, "schedule created successfully", schedule)
}

func (h *ScheduleHandler) Update(w http.ResponseWriter, r *http.Request) {
	id, ok := parseScheduleID(w, r)
	if !ok {
		return
	}

	var req scheduleRequest
	if err := decodeJSON(r, &req); err != nil {
		response.Error(w, http.StatusBadRequest, "invalid request body")
		return
	}

	schedule, err := h.schedules.Update(r.Context(), id, req.toPayload())
	if err != nil {
		h.writeError(w, err)
		return
	}

	response.MessageData(w, http.StatusOK, "schedule updated successfully", schedule)
}

func (h *ScheduleHandler) Delete(w http.ResponseWriter, r *http.Request) {
	id, ok := parseScheduleID(w, r)
	if !ok {
		return
	}

	if err := h.schedules.Delete(r.Context(), id); err != nil {
		h.writeError(w, err)
		return
	}

	response.Message(w, http.StatusOK, "schedule deleted successfully")
}

func (h *ScheduleHandler) writeError(w http.ResponseWriter, err error) {
	var validationErrs *domain.ValidationErrors
	switch {
	case errors.As(err, &validationErrs):
		response.ValidationError(w, *validationErrs)
	case errors.Is(err, domain.ErrNotFound):
		response.Error(w, http.StatusNotFound, "schedule not found")
	case errors.Is(err, domain.ErrConflict):
		response.Error(w, http.StatusConflict, "schedule already exists for this studio, date, and time")
	default:
		log.Printf("schedule handler: %v", err)
		response.Error(w, http.StatusInternalServerError, "internal server error")
	}
}

func parseScheduleID(w http.ResponseWriter, r *http.Request) (int64, bool) {
	id, err := strconv.ParseInt(chi.URLParam(r, "id"), 10, 64)
	if err != nil || id <= 0 {
		response.Error(w, http.StatusBadRequest, "invalid schedule id")
		return 0, false
	}
	return id, true
}

func parseScheduleFilter(r *http.Request) (model.ScheduleFilter, error) {
	query := r.URL.Query()
	filter := model.ScheduleFilter{}

	for key, target := range map[string]**int64{
		"movie_id":  &filter.MovieID,
		"branch_id": &filter.BranchID,
		"studio_id": &filter.StudioID,
	} {
		raw := strings.TrimSpace(query.Get(key))
		if raw == "" {
			continue
		}
		value, err := strconv.ParseInt(raw, 10, 64)
		if err != nil || value <= 0 {
			return filter, errors.New("invalid " + key + " filter")
		}
		*target = &value
	}

	if raw := strings.TrimSpace(query.Get("show_date")); raw != "" {
		if _, err := time.Parse("2006-01-02", raw); err != nil {
			return filter, errors.New("invalid show_date filter, expected YYYY-MM-DD")
		}
		filter.ShowDate = raw
	}

	return filter, nil
}
