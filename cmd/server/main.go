package main

import (
	"context"
	"log"
	"os"
	"os/signal"
	"syscall"

	"Test/internal/config"
	"Test/internal/database"
	"Test/internal/handler"
	"Test/internal/repository"
	"Test/internal/server"
	"Test/internal/service"
)

func main() {
	if err := run(); err != nil {
		log.Fatalf("server gagal berjalan: %v", err)
	}
}

func run() error {
	cfg, err := config.Load()
	if err != nil {
		return err
	}

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
	defer stop()

	db, err := database.Connect(ctx, cfg.DatabaseURL, cfg.DBMaxConns)
	if err != nil {
		return err
	}
	defer db.Close()

	userRepo := repository.NewUserRepository(db)
	scheduleRepo := repository.NewScheduleRepository(db)
	movieRepo := repository.NewMovieRepository(db)
	branchRepo := repository.NewBranchRepository(db)
	studioRepo := repository.NewStudioRepository(db)

	authService := service.NewAuthService(userRepo, cfg.JWTSecret, cfg.JWTExpiry)
	scheduleService := service.NewScheduleService(scheduleRepo, movieRepo, branchRepo, studioRepo)

	handlers := server.Handlers{
		Health:   handler.NewHealthHandler(db),
		Auth:     handler.NewAuthHandler(authService),
		Schedule: handler.NewScheduleHandler(scheduleService),
	}

	return server.New(cfg.AppPort, handlers, []byte(cfg.JWTSecret)).Run(ctx)
}
