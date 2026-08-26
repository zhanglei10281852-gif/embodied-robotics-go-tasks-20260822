package worker

import (
	"context"
	"errors"
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

// TestSupervisorStopReportsClean verifies that stopping a worker blocked on
// cancellation is reported as a normal stop (no error on done) rather than
// forwarded as an ordinary context.Canceled failure.
func TestSupervisorStopReportsClean(t *testing.T) {
	s := NewSupervisor()
	done := s.Start("w", func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	s.Stop("w")
	select {
	case e, ok := <-done:
		if ok && e != nil {
			t.Fatalf("expected clean stop, got %v", e)
		}
	case <-time.After(time.Second):
		t.Fatal("worker not stopped")
	}
}

// TestSupervisorForwardRealErrors verifies that genuine failures are still
// surfaced on the done channel.
func TestSupervisorForwardRealErrors(t *testing.T) {
	s := NewSupervisor()
	done := s.Start("f", func(ctx context.Context) error {
		return errors.New("boom")
	})
	select {
	case e := <-done:
		if e == nil || e.Error() != "boom" {
			t.Fatalf("expected boom, got %v", e)
		}
	case <-time.After(time.Second):
		t.Fatal("no result")
	}
}

// TestSupervisorStopAllReportsClean verifies StopAll also yields no error for
// cancellation-triggered exits.
func TestSupervisorStopAllReportsClean(t *testing.T) {
	s := NewSupervisor()
	done := s.Start("a", func(ctx context.Context) error {
		<-ctx.Done()
		return ctx.Err()
	})
	s.StopAll()
	select {
	case e, ok := <-done:
		if ok && e != nil {
			t.Fatalf("expected clean stop, got %v", e)
		}
	case <-time.After(time.Second):
		t.Fatal("worker not stopped")
	}
}
