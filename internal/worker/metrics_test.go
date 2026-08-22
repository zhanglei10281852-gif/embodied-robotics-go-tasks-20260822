package worker

import (
	"context"
	"testing"
	"time"
)

func TestCounterSnapshot(t *testing.T) {
	c := NewCounter()
	c.Add("missions", 2)
	if c.Get("missions") != 2 {
		t.Fatal()
	}
	snap := c.Snapshot()
	snap["missions"] = 0
	if c.Get("missions") != 2 {
		t.Fatal("alias")
	}
	ctx, cancel := context.WithCancel(context.Background())
	done := Meter{Counter: c, Interval: time.Millisecond}.Run(ctx, func() map[string]int64 { return map[string]int64{"ticks": 1} })
	time.Sleep(3 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("meter")
	}
}
