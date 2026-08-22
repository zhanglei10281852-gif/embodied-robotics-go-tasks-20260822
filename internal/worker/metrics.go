package worker

import (
	"context"
	"sync"
	"time"
)

type Counter struct {
	mu     sync.Mutex
	values map[string]int64
}

func NewCounter() *Counter                  { return &Counter{values: map[string]int64{}} }
func (c *Counter) Add(name string, n int64) { c.mu.Lock(); defer c.mu.Unlock(); c.values[name] += n }
func (c *Counter) Get(name string) int64    { c.mu.Lock(); defer c.mu.Unlock(); return c.values[name] }
func (c *Counter) Snapshot() map[string]int64 {
	c.mu.Lock()
	defer c.mu.Unlock()
	out := map[string]int64{}
	for k, v := range c.values {
		out[k] = v
	}
	return out
}

type Meter struct {
	Counter  *Counter
	Interval time.Duration
}

func (m Meter) Run(ctx context.Context, fn func() map[string]int64) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		tick := time.NewTicker(m.Interval)
		defer tick.Stop()
		for {
			select {
			case <-ctx.Done():
				return
			case <-tick.C:
				for k, v := range fn() {
					m.Counter.Add(k, v)
				}
			}
		}
	}()
	return done
}
