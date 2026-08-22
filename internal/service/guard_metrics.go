package service

import (
	"sync"
	"time"
)

type GuardMetrics struct {
	mu                       sync.Mutex
	Allowed, Denied, Expired int64
	Last                     time.Time
}

func (g *GuardMetrics) Allow() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Allowed++
	g.Last = time.Now().UTC()
}
func (g *GuardMetrics) Deny() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Denied++
	g.Last = time.Now().UTC()
}
func (g *GuardMetrics) Expire() {
	g.mu.Lock()
	defer g.mu.Unlock()
	g.Expired++
	g.Last = time.Now().UTC()
}
func (g *GuardMetrics) Snapshot() (int64, int64, int64, time.Time) {
	g.mu.Lock()
	defer g.mu.Unlock()
	return g.Allowed, g.Denied, g.Expired, g.Last
}
