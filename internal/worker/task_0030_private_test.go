package worker

import (
	"context"
	"testing"
	"time"
)

func TestSupervisorCancellationIsNotReportedAsFailure(t *testing.T) {
	s := NewSupervisor()
	done := s.Start("cancel", func(ctx context.Context) error { <-ctx.Done(); return ctx.Err() })
	s.Stop("cancel")
	select {
	case err, ok := <-done:
		if ok && err != nil {
			t.Fatalf("cancellation surfaced as failure: %v", err)
		}
	case <-time.After(time.Second):
		t.Fatal("supervisor did not stop")
	}
}
