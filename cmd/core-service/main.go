package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	pgrepo "oiltech/internal/adapters/postgres"
	"oiltech/internal/config"
	"oiltech/internal/db"
	"oiltech/internal/handlers"
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

	datasetUsecase := usecases.NewDatasetUsecase(datasetRepo, trm)
	algorithmUsecase := usecases.NewAlgorithmUsecase(algorithmRepo)
	runUsecase := usecases.NewRunUsecase(transformationRepo, runRepo, trm)
	transformationUsecase := usecases.NewTransformationUsecase(datasetRepo, algorithmRepo, transformationRepo, runUsecase, trm)

	server := &http.Server{
		Addr:              cfg.Server.CoreAddr,
		Handler:           handlers.NewHandler(datasetUsecase, algorithmUsecase, transformationUsecase, runUsecase).Routes(),
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("core-service started", "addr", cfg.Server.CoreAddr)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("http server failed", "error", err)
			os.Exit(1)
		}
	}()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	<-ctx.Done()

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	_ = server.Shutdown(shutdownCtx)
}
