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
	"oiltech/internal/mqttcollector"
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
	trm := pgrepo.NewTransactionManager(database)
	datasetUsecase := usecases.NewDatasetUsecase(datasetRepo, trm)
	collector := mqttcollector.New(datasetUsecase, cfg.MQTT.Broker, cfg.MQTT.TopicPrefix, logger)
	if err := collector.Start(); err != nil {
		logger.Error("collector start failed", "error", err)
		os.Exit(1)
	}
	defer collector.Stop()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", mqttcollector.NewHealthHandler(collector))
	server := &http.Server{
		Addr:              cfg.Server.CollectorAddr,
		Handler:           mux,
		ReadHeaderTimeout: 5 * time.Second,
	}
	go func() {
		logger.Info("collector-service started", "addr", cfg.Server.CollectorAddr)
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
