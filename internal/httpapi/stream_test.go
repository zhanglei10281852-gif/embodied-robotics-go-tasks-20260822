package httpapi

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/worker"
	"net/http/httptest"
	"testing"
	"time"
)

func TestStreamCancellation(t *testing.T) {
	b := worker.NewBus()
	s := Stream{Bus: b}
	ctx, cancel := context.WithCancel(context.Background())
	req := httptest.NewRequest("GET", "/events", nil).WithContext(ctx)
	rw := httptest.NewRecorder()
	done := make(chan struct{})
	go func() { s.ServeHTTP(rw, req); close(done) }()
	time.Sleep(10 * time.Millisecond)
	cancel()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("stream not cancelled")
	}
}
