package server

import (
	"net/http"
	"time"

	"Test/internal/middleware"
	"Test/internal/response"

	"github.com/go-chi/chi/v5"
	chimw "github.com/go-chi/chi/v5/middleware"
)

func routes(h Handlers, jwtSecret []byte) http.Handler {
	r := chi.NewRouter()
	r.Use(chimw.RequestID)
	r.Use(chimw.Logger)
	r.Use(chimw.RealIP)
	r.Use(chimw.Recoverer)
	r.Use(chimw.Timeout(15 * time.Second))

	r.Get("/", func(w http.ResponseWriter, r *http.Request) {
		response.JSON(w, http.StatusOK, map[string]string{"message": "schedule api is running"})
	})
	r.Get("/health", h.Health.Check)

	r.Route("/api/v1", func(r chi.Router) {
		r.Post("/auth/login", h.Auth.Login)

		r.Group(func(r chi.Router) {
			r.Use(middleware.Auth(jwtSecret))
			r.Get("/schedules", h.Schedule.List)
			r.Get("/schedules/{id}", h.Schedule.Get)
			r.Post("/schedules", h.Schedule.Create)
			r.Put("/schedules/{id}", h.Schedule.Update)
			r.Delete("/schedules/{id}", h.Schedule.Delete)
		})
	})

	r.NotFound(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusNotFound, "endpoint not found")
	})
	r.MethodNotAllowed(func(w http.ResponseWriter, r *http.Request) {
		response.Error(w, http.StatusMethodNotAllowed, "method not allowed")
	})

	return r
}
