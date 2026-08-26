package worker

import (
	"context"
	"errors"
	"sync/atomic"
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

func TestRunWithHeartbeatCancelAfterFirstTick(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())

	var calls int32
	started := make(chan struct{})
	err := make(chan error, 1)
	go func() {
		close(started)
		err <- RunWithHeartbeat(ctx, 10*time.Millisecond, func() error {
			atomic.AddInt32(&calls, 1)
			return nil
		})
	}()

	<-started
	// wait for the first tick, then cancel the caller context.
	time.Sleep(50 * time.Millisecond)
	before := atomic.LoadInt32(&calls)
	if before < 1 {
		t.Fatalf("expected at least one callback, got %d", before)
	}
	cancel()

	select {
	case e := <-err:
		if !errors.Is(e, context.Canceled) {
			t.Fatalf("expected context.Canceled, got %v", e)
		}
	case <-time.After(time.Second):
		t.Fatal("RunWithHeartbeat did not return after cancel")
	}

	// give the loop a chance to wrongly fire again after cancellation.
	time.Sleep(50 * time.Millisecond)
	after := atomic.LoadInt32(&calls)
	if after != before {
		t.Fatalf("callback fired after cancel: before=%d after=%d", before, after)
	}
}
