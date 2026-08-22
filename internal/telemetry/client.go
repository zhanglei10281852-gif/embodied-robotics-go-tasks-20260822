package telemetry

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrDisconnected = errors.New("robot transport disconnected")

type Client interface {
	Send(context.Context, string, []byte) error
	Close() error
}
type MemoryClient struct {
	mu     sync.Mutex
	closed bool
	Sent   [][]byte
	Delay  time.Duration
}

func (c *MemoryClient) Send(ctx context.Context, _ string, p []byte) error {
	c.mu.Lock()
	closed := c.closed
	d := c.Delay
	c.mu.Unlock()
	if closed {
		return ErrDisconnected
	}
	timer := time.NewTimer(d)
	defer timer.Stop()
	select {
	case <-ctx.Done():
		return ctx.Err()
	case <-timer.C:
	}
	c.mu.Lock()
	defer c.mu.Unlock()
	if c.closed {
		return ErrDisconnected
	}
	v := append([]byte(nil), p...)
	c.Sent = append(c.Sent, v)
	return nil
}
func (c *MemoryClient) Close() error { c.mu.Lock(); defer c.mu.Unlock(); c.closed = true; return nil }
func (c *MemoryClient) Count() int   { c.mu.Lock(); defer c.mu.Unlock(); return len(c.Sent) }

type Window struct{ Start, End time.Time }

func (w Window) Contains(t time.Time) bool { return !t.Before(w.Start) && t.Before(w.End) }
