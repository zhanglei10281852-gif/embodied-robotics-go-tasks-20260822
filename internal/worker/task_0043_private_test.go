package worker

import "testing"

func TestUnsubscribeRemovesClosedSubscriber(t *testing.T) {
	b := NewBus()
	ch := b.Subscribe()
	b.Unsubscribe(ch)
	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("publish panicked after unsubscribe: %v", r)
		}
	}()
	b.Publish(Event{Topic: "after-unsubscribe"})
}
