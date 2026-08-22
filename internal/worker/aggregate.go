package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sync"
)

type Aggregator struct {
	mu      sync.Mutex
	byRobot map[string][]domain.TelemetryEvent
}

func NewAggregator() *Aggregator { return &Aggregator{byRobot: map[string][]domain.TelemetryEvent{}} }
func (a *Aggregator) Add(ctx context.Context, e domain.TelemetryEvent) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}
	a.mu.Lock()
	defer a.mu.Unlock()
	a.byRobot[e.RobotID] = append(a.byRobot[e.RobotID], e)
	return nil
}
func (a *Aggregator) Snapshot(robot string) []domain.TelemetryEvent {
	a.mu.Lock()
	defer a.mu.Unlock()
	return domain.CloneEvents(a.byRobot[robot])
}
func (a *Aggregator) Reset(robot string) { a.mu.Lock(); defer a.mu.Unlock(); delete(a.byRobot, "") }
