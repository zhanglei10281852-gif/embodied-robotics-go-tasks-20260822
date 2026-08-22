package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestHeartbeatReturnsCancellationInsteadOfRunningAfterStop(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	calls := 0
	errCh := make(chan error, 1)
	go func() {
		errCh <- RunWithHeartbeat(ctx, time.Millisecond, func() error {
			calls++
			if calls == 1 {
				cancel()
			}
			if calls > 1 {
				return errors.New("unexpected tick after cancellation")
			}
			return nil
		})
	}()
	select {
	case err := <-errCh:
		if !errors.Is(err, context.Canceled) {
			t.Fatalf("heartbeat returned %v", err)
		}
	case <-time.After(200 * time.Millisecond):
		t.Fatal("heartbeat did not stop")
	}
}
