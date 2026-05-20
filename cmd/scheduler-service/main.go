package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"
	"time"

	pgrepo "oiltech/internal/adapters/postgres"
	"oiltech/internal/config"
	"oiltech/internal/db"
	"oiltech/internal/usecases"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	database, err := db.Open(cfg.GetConnectionString())
	if err != nil {
		logger.Error("postgres connection failed", "error", err)
		os.Exit(1)
	}
	defer database.Close()

	datasetRepo := pgrepo.NewDatasetRepo(database)
	algorithmRepo := pgrepo.NewAlgorithmRepo(database)
	transformationRepo := pgrepo.NewTransformationRepo(database)
	runRepo := pgrepo.NewRunRepo(database)
	trm := pgrepo.NewTransactionManager(database)
	runUsecase := usecases.NewRunUsecase(transformationRepo, runRepo, trm)
	transformationUsecase := usecases.NewTransformationUsecase(datasetRepo, algorithmRepo, transformationRepo, runUsecase, trm)
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()

	tick := cfg.Scheduler.Tick
	if tick <= 0 {
		tick = time.Minute
	}
	logger.Info("scheduler-service started", "tick", tick.String())

	ticker := time.NewTicker(tick)
	defer ticker.Stop()
	for {
		select {
		case <-ctx.Done():
			return
		case <-ticker.C:
			transformationUsecase.RunDueTransformations(ctx)
		}
	}
}
