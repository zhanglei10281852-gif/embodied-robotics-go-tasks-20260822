package worker

import (
	"context"
	"testing"
	"time"
)

func TestSupervisorStop(t *testing.T) {
	s := NewSupervisor()
	done := s.Start("x", func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() })
	s.Stop("x")
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("not stopped")
	}
}
