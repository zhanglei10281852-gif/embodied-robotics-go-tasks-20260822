package worker

import (
	"context"
	"errors"
	"testing"
)

func TestMaintenanceRunPropagatesCancellation(t *testing.T) {
	m := NewMaintenance()
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	err := m.Run(ctx, "tenant", func(got context.Context) error {
		select {
		case <-got.Done():
			return got.Err()
		default:
			return nil
		}
	})
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
}
