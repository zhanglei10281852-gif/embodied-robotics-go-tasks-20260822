package telemetry

import (
	"context"
	"testing"
	"time"
)

func TestWindowsAndClientCopies(t *testing.T) {
	now := time.Now()
	w := CurrentWindow(now, time.Minute)
	if !w.Contains(now.Add(-time.Second)) {
		t.Fatal("window")
	}
	parts := Split(w, 4)
	if len(parts) != 4 || parts[0].Start != w.Start {
		t.Fatal("split")
	}
	c := &MemoryClient{}
	payload := []byte("robot")
	if e := c.Send(context.Background(), "r", payload); e != nil {
		t.Fatal(e)
	}
	payload[0] = 'X'
	if c.Count() != 1 {
		t.Fatal("count")
	}
}
func TestDecodeInvalid(t *testing.T) {
	if _, e := Decode([]byte("not gzip")); e == nil {
		t.Fatal("invalid accepted")
	}
}
