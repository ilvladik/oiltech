package config

import (
	"fmt"
	"os"
	"time"
)

type Config struct {
	Server struct {
		CoreAddr             string
		CollectorAddr        string
		EquipWorkPeriodsAddr string
	}
	Database struct {
		Host     string
		Port     string
		Name     string
		User     string
		Password string
		SSLMode  string
	}
	MQTT struct {
		Broker      string
		TopicPrefix string
	}
	Scheduler struct {
		Tick time.Duration
	}
}

func Load() Config {
	var cfg Config
	cfg.Server.CoreAddr = getEnv("CORE_HTTP_ADDR", ":8080")
	cfg.Server.CollectorAddr = getEnv("COLLECTOR_HTTP_ADDR", ":8081")
	cfg.Server.EquipWorkPeriodsAddr = getEnv("EQUIP_WORK_PERIODS_HTTP_ADDR", ":9000")
	cfg.Database.Host = getEnv("POSTGRES_HOST", "postgres")
	cfg.Database.Port = getEnv("POSTGRES_PORT", "5432")
	cfg.Database.Name = getEnv("POSTGRES_DB", "oiltech")
	cfg.Database.User = getEnv("POSTGRES_USER", "oiltech")
	cfg.Database.Password = getEnv("POSTGRES_PASSWORD", "oiltech")
	cfg.Database.SSLMode = getEnv("POSTGRES_SSLMODE", "disable")
	cfg.MQTT.Broker = getEnv("MQTT_BROKER", "tcp://mqtt:1883")
	cfg.MQTT.TopicPrefix = getEnv("MQTT_TOPIC_PREFIX", "oiltech/datasets")

	tick, err := time.ParseDuration(getEnv("SCHEDULER_TICK", "1m"))
	if err != nil {
		tick = time.Minute
	}
	cfg.Scheduler.Tick = tick
	return cfg
}

func (c *Config) GetConnectionString() string {
	if dsn := os.Getenv("POSTGRES_DSN"); dsn != "" {
		return dsn
	}
	return fmt.Sprintf(
		"postgres://%s:%s@%s:%s/%s?sslmode=%s",
		c.Database.User,
		c.Database.Password,
		c.Database.Host,
		c.Database.Port,
		c.Database.Name,
		c.Database.SSLMode,
	)
}

func getEnv(key, fallback string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return fallback
}
