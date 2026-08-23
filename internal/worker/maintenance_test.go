package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestMaintenanceSingleFlight(t *testing.T) {
	m := NewMaintenance()
	if !m.Start("t", time.Now()) {
		t.Fatal("first start")
	}
	if m.Start("t", time.Now()) {
		t.Fatal("second start")
	}
	m.Finish()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := m.Run(ctx, "t", func(context.Context) error { return nil }); e == nil {
		t.Fatal("cooldown should apply")
	}
}

func TestMaintenanceCancelledContextDoesNotRunCallback(t *testing.T) {
	m := NewMaintenance()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	ran := false
	e := m.Run(ctx, "tenant", func(context.Context) error {
		ran = true
		return nil
	})
	if !errors.Is(e, context.Canceled) {
		t.Fatalf("expected context.Canceled, got %v", e)
	}
	if ran {
		t.Fatal("callback should not run when context is cancelled before start")
	}
}
