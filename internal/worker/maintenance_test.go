package worker

import (
	"context"
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
