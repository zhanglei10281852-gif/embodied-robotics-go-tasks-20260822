package worker

import "testing"

func TestPublishAfterCloseIsSafe(t *testing.T) {
	b := NewBus()
	b.Close()
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("publish panicked after close: %v", r)
		}
	}()
	b.Publish(Event{Topic: "shutdown"})
}
