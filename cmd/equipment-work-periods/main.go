package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"time"

	pgrepo "oiltech/internal/adapters/postgres"
	"oiltech/internal/algokit"
	"oiltech/internal/config"
	"oiltech/internal/db"
	"oiltech/internal/equipworkperiods"
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
	transformationRepo := pgrepo.NewTransformationRepo(database)
	runRepo := pgrepo.NewRunRepo(database)
	kit := algokit.NewKit(datasetRepo, transformationRepo, runRepo)
	runner := equipworkperiods.NewRunner(kit, database)

	addr := cfg.Server.EquipWorkPeriodsAddr

	mux := http.NewServeMux()
	mux.HandleFunc("POST /run", runner.Handle)
	mux.HandleFunc("GET /healthz", func(w http.ResponseWriter, _ *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		_, _ = w.Write([]byte(`{"status":"ok"}`))
	})

	server := &http.Server{
		Addr:              addr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}

	go func() {
		logger.Info("equipment-work-periods started", "addr", addr)
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
