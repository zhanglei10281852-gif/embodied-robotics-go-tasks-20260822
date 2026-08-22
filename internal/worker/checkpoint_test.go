package worker

import (
	"context"
	"testing"
)

func TestCheckpointMonotonic(t *testing.T) {
	c := NewCheckpoint()
	if e := c.Advance("r", 2); e != nil {
		t.Fatal(e)
	}
	if e := c.Advance("r", 1); e == nil {
		t.Fatal("regression")
	}
	if v, ok := c.Read("r"); !ok || v != 2 {
		t.Fatal(v, ok)
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if e := c.Wait(ctx, "r", 3); e == nil {
		t.Fatal("cancel")
	}
}
