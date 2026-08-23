package worker

import (
	"context"
	"errors"
	"sync"
	"time"
)

type Maintenance struct {
	mu      sync.Mutex
	last    map[string]time.Time
	running bool
}

func NewMaintenance() *Maintenance { return &Maintenance{last: map[string]time.Time{}} }
func (m *Maintenance) Start(tenant string, now time.Time) bool {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.running {
		return false
	}
	if last := m.last[tenant]; !last.IsZero() && now.Sub(last) < time.Minute {
		return false
	}
	m.running = true
	m.last[tenant] = now
	return true
}
func (m *Maintenance) Finish() { m.mu.Lock(); m.running = false; m.mu.Unlock() }
func (m *Maintenance) Run(ctx context.Context, tenant string, fn func(context.Context) error) error {
	if !m.Start(tenant, time.Now().UTC()) {
		return errors.New("maintenance already running")
	}
	defer m.Finish()
	if e := ctx.Err(); e != nil {
		return e
	}
	return fn(ctx)
}
