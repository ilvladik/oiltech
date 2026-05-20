package main

import (
	"context"
	"log/slog"
	"os"
	"os/signal"

	"oiltech/internal/config"
	"oiltech/internal/simulator"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{Level: slog.LevelInfo}))
	cfg := config.Load()

	sim := simulator.New(cfg.MQTT.Broker, cfg.MQTT.TopicPrefix, logger)
	if err := sim.Start(); err != nil {
		logger.Error("simulator start failed", "error", err)
		os.Exit(1)
	}
	defer sim.Stop()

	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt)
	defer stop()
	sim.PublishLoop(ctx.Done())
}
