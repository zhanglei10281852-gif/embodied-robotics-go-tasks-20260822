package worker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Supervisor struct {
	mu      sync.Mutex
	running map[string]context.CancelFunc
}

func NewSupervisor() *Supervisor { return &Supervisor{running: map[string]context.CancelFunc{}} }
func (s *Supervisor) Start(name string, fn func(context.Context) error) <-chan error {
	ctx, cancel := context.WithCancel(context.Background())
	s.mu.Lock()
	if old := s.running[name]; old != nil {
		old()
	}
	s.running[name] = cancel
	s.mu.Unlock()
	done := make(chan error, 1)
	go func() {
		defer close(done)
		e := fn(ctx)
		s.mu.Lock()
		delete(s.running, name)
		s.mu.Unlock()
		// Stopping a worker cancels its context. A worker that exits with
		// context.Canceled in response is performing an orderly shutdown, not
		// failing: the supervisor owns the only cancel func and the context is
		// derived from Background, so a canceled ctx can only come from a
		// supervisor-initiated stop. Report it as success rather than forwarding
		// it as an ordinary error.
		if errors.Is(e, context.Canceled) && ctx.Err() != nil {
			return
		}
		if e != nil {
			done <- e
		}
	}()
	return done
}
func (s *Supervisor) Stop(name string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	if c := s.running[name]; c != nil {
		c()
		delete(s.running, name)
	}
}
func (s *Supervisor) StopAll() {
	s.mu.Lock()
	defer s.mu.Unlock()
	for n, c := range s.running {
		c()
		delete(s.running, n)
	}
}
func RunWithHeartbeat(ctx context.Context, interval time.Duration, fn func() error) error {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			if e := fn(); e != nil {
				return e
			}
		}
	}
}
