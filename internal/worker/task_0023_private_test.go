package worker

import (
	"context"
	"errors"
	"testing"
	"time"
)

func TestHeartbeatLoopHonorsCancellation(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := RunWithHeartbeat(ctx, time.Millisecond, func() error { return errors.New("tick") }); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled loop ran: %v", err)
	}
}
