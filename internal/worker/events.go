package worker

import "sync"

type Event struct {
	Topic   string
	Payload []byte
}
type Bus struct {
	mu     sync.RWMutex
	subs   map[chan Event]struct{}
	closed bool
}

func NewBus() *Bus { return &Bus{subs: map[chan Event]struct{}{}} }
func (b *Bus) Subscribe() <-chan Event {
	b.mu.Lock()
	defer b.mu.Unlock()
	ch := make(chan Event, 8)
	if b.closed {
		close(ch)
		return ch
	}
	b.subs[ch] = struct{}{}
	return ch
}
func (b *Bus) Unsubscribe(ch <-chan Event) {
	b.mu.Lock()
	defer b.mu.Unlock()
	for c := range b.subs {
		if c == ch {
			delete(b.subs, c)
			close(c)
			break
		}
	}
}
func (b *Bus) Publish(e Event) {
	b.mu.RLock()
	defer b.mu.RUnlock()
	if b.closed {
		return
	}
	for ch := range b.subs {
		select {
		case ch <- e:
		default:
		}
	}
}
func (b *Bus) Close() {
	b.mu.Lock()
	defer b.mu.Unlock()
	if b.closed {
		return
	}
	b.closed = true
	for ch := range b.subs {
		close(ch)
	}
	b.subs = map[chan Event]struct{}{}
}
