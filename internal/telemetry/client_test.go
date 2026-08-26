package telemetry

import (
	"context"
	"testing"
	"time"
)

func TestClientCancellationAndClose(t *testing.T) {
	c := &MemoryClient{Delay: time.Second}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := c.Send(ctx, "r", []byte("x")); e == nil {
		t.Fatal("expected cancellation")
	}
	if e := c.Close(); e != nil {
		t.Fatal(e)
	}
	if e := c.Send(context.Background(), "r", nil); e != ErrDisconnected {
		t.Fatal(e)
	}
}
