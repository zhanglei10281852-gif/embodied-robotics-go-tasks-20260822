package config

import (
	"os"
	"testing"
	"time"
)

func TestLoadDefaultsAndOverrides(t *testing.T) {
	for _, k := range []string{"PORT", "DATABASE_URL", "SESSION_TTL", "WORKER_INTERVAL"} {
		t.Setenv(k, "")
	}
	c := Load()
	if c.Port != "8080" || c.DatabaseURL == "" || c.SessionTTL != 30*time.Minute {
		t.Fatalf("defaults: %+v", c)
	}
	t.Setenv("PORT", "9090")
	t.Setenv("SESSION_TTL", "2m")
	t.Setenv("WORKER_INTERVAL", "1s")
	o := Load()
	if o.PortNumber() != 9090 || o.SessionTTL != 2*time.Minute || o.WorkerInterval != time.Second {
		t.Fatalf("overrides: %+v", o)
	}
	_ = os.Getenv("PORT")
}

func TestConfigValidation(t *testing.T) {
	c := Config{DatabaseURL: "file:test.db", SessionTTL: time.Minute, WorkerInterval: time.Second}
	if err := c.Validate(); err != nil {
		t.Fatal(err)
	}
	c.SessionTTL = 0
	if err := c.Validate(); err == nil {
		t.Fatal("expected invalid duration")
	}
}
