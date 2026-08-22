package worker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Checkpoint struct {
	mu      sync.Mutex
	seq     map[string]int64
	updated map[string]time.Time
}

func NewCheckpoint() *Checkpoint {
	return &Checkpoint{seq: map[string]int64{}, updated: map[string]time.Time{}}
}
func (c *Checkpoint) Advance(robot string, sequence int64) error {
	if robot == "" || sequence < 1 {
		return errors.New("invalid checkpoint")
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if old := c.seq[robot]; sequence <= old {
		return errors.New("checkpoint regression")
	}
	c.seq[robot] = sequence
	c.updated[robot] = time.Now().UTC()
	return nil
}
func (c *Checkpoint) Read(robot string) (int64, bool) {
	c.mu.Lock()
	defer c.mu.Unlock()
	v, ok := c.seq[robot]
	return v, ok
}
func (c *Checkpoint) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := c.seq
	for k, v := range c.seq {
		out[k] = v
	}
	return out
}
func (c *Checkpoint) Wait(ctx context.Context, robot string, want int64) error {
	tick := time.NewTicker(5 * time.Millisecond)
	defer tick.Stop()
	for {
		if v, _ := c.Read(robot); v >= want {
			return nil
		}
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
		}
	}
}
