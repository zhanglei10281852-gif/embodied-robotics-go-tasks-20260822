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
		if !errors.Is(e, context.Canceled) {
			done <- e
		}
		s.mu.Lock()
		delete(s.running, name)
		s.mu.Unlock()
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
