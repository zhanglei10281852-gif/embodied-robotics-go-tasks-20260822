package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestAggregatorCopiesAndCancellation(t *testing.T) {
	a := NewAggregator()
	e := domain.TelemetryEvent{RobotID: "r", Sequence: 1}
	if err := a.Add(context.Background(), e); err != nil {
		t.Fatal(err)
	}
	x := a.Snapshot("r")
	x[0].Sequence = 9
	if a.Snapshot("r")[0].Sequence != 1 {
		t.Fatal("snapshot alias")
	}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	if err := a.Add(ctx, e); err == nil {
		t.Fatal("expected cancellation")
	}
}
func TestBusCloseAndBackoff(t *testing.T) {
	b := NewBus()
	ch := b.Subscribe()
	b.Publish(Event{Topic: "x"})
	select {
	case <-ch:
	case <-time.After(time.Second):
		t.Fatal("event not delivered")
	}
	b.Close()
	if _, ok := <-ch; ok {
		t.Fatal("channel open")
	}
	if d := (Backoff{Base: time.Second, Max: 4 * time.Second}).Delay(3); d != 4*time.Second {
		t.Fatal(d)
	}
}
