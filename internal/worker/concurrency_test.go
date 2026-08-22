package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sync"
	"testing"
	"time"
)

func TestAggregatorConcurrentAdds(t *testing.T) {
	a := NewAggregator()
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			_ = a.Add(context.Background(), domain.TelemetryEvent{RobotID: "r", Sequence: int64(n + 1)})
		}(i)
	}
	wg.Wait()
	if got := len(a.Snapshot("r")); got != 20 {
		t.Fatal(got)
	}
}
func TestBusBackpressureDoesNotBlock(t *testing.T) {
	b := NewBus()
	_ = b.Subscribe()
	done := make(chan struct{})
	go func() {
		for i := 0; i < 100; i++ {
			b.Publish(Event{Topic: "x"})
		}
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("publish blocked")
	}
	b.Close()
}
