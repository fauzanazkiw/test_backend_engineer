package handler

import (
	"Test/internal/response"
	"context"
	"net/http"
	"time"
)

// Pinger adalah dependency minimal yang dibutuhkan health handler.
type Pinger interface {
	PingContext(ctx context.Context) error
}

type HealthHandler struct {
	db Pinger
}

func NewHealthHandler(db Pinger) *HealthHandler {
	return &HealthHandler{db: db}
}

// Check menangani GET /health dengan memastikan koneksi database tersedia.
func (h *HealthHandler) Check(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), 2*time.Second)
	defer cancel()

	if err := h.db.PingContext(ctx); err != nil {
		response.Error(w, http.StatusServiceUnavailable, "database is not available")
		return
	}

	response.JSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
