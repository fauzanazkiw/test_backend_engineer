package server

import (
	"Test/internal/handler"
	"context"
	"errors"
	"log"
	"net/http"
	"time"
)

type Handlers struct {
	Health   *handler.HealthHandler
	Auth     *handler.AuthHandler
	Schedule *handler.ScheduleHandler
}

type Server struct {
	HttpServer *http.Server
	Port       string
}

func New(port string, handlers Handlers, jwtSecret []byte) *Server {
	return &Server{
		Port: port,
		HttpServer: &http.Server{
			Addr:              ":" + port,
			Handler:           routes(handlers, jwtSecret),
			ReadHeaderTimeout: 5 * time.Second,
			ReadTimeout:       10 * time.Second,
			WriteTimeout:      15 * time.Second,
			IdleTimeout:       60 * time.Second,
		},
	}
}

// Run menjalankan HTTP server dan melakukan graceful shutdown saat context dibatalkan.
func (s *Server) Run(ctx context.Context) error {
	errCh := make(chan error, 1)
	go func() {
		log.Printf("server berjalan pada port %s", s.Port)
		if err := s.HttpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			errCh <- err
		}
	}()

	select {
	case err := <-errCh:
		return err

	case <-ctx.Done():
		log.Println("sinyal berhenti diterima, mematikan server...")
	}

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	return s.HttpServer.Shutdown(shutdownCtx)
}
