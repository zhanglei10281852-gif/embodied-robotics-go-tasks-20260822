package config

import (
	"testing"
	"time"
)

func TestDurationFallbacks(t *testing.T) {
	t.Setenv("SESSION_TTL", "bad")
	t.Setenv("WORKER_INTERVAL", "-1s")
	c := Load()
	if c.SessionTTL <= 0 || c.WorkerInterval <= 0 {
		t.Fatal("fallback not applied")
	}
}
func TestAddressAndOrigins(t *testing.T) {
	c := Config{Port: "8123", DatabaseURL: "file:x", SessionTTL: time.Minute, WorkerInterval: time.Second}
	o := c.Options()
	if o.Address() != "0.0.0.0:8123" {
		t.Fatal(o.Address())
	}
	if o.HasOrigin("http://LOCALHOST") != true {
		t.Fatal("origin case")
	}
	if o.HasOrigin("https://example") {
		t.Fatal("unexpected origin")
	}
}
