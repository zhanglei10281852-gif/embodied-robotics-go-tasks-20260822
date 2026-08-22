package config

import (
	"os"
	"strconv"
	"time"
)

type Config struct {
	Port           string
	DatabaseURL    string
	SessionTTL     time.Duration
	WorkerInterval time.Duration
}

func Load() Config {
	return Config{
		Port:           value("PORT", "8080"),
		DatabaseURL:    value("DATABASE_URL", "file:robotics.db?_pragma=foreign_keys(1)"),
		SessionTTL:     duration("SESSION_TTL", 30*time.Minute),
		WorkerInterval: duration("WORKER_INTERVAL", 250*time.Millisecond),
	}
}

func value(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}

func duration(key string, fallback time.Duration) time.Duration {
	v := os.Getenv(key)
	if v == "" {
		return fallback
	}
	if d, err := time.ParseDuration(v); err == nil && d > 0 {
		return d
	}
	return fallback
}

func (c Config) PortNumber() int {
	n, err := strconv.Atoi(c.Port)
	if err != nil || n < 1 {
		return 8080
	}
	return n
}

func (c Config) Validate() error {
	if c.DatabaseURL == "" {
		return ErrDatabaseURL
	}
	if c.SessionTTL <= 0 || c.WorkerInterval <= 0 {
		return ErrDuration
	}
	return nil
}
