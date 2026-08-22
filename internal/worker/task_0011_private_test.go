package worker

import (
	"testing"
	"time"
)

func TestUnsubscribeDoesNotRaceWithPublish(t *testing.T) {
	b := NewBus()
	ch := b.Subscribe()
	stop := make(chan struct{})
	go func() {
		for {
			select {
			case <-stop:
				return
			default:
				b.Publish(Event{Topic: "pose"})
			}
		}
	}()
	time.Sleep(time.Millisecond)
	b.Unsubscribe(ch)
	close(stop)
}
