package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestCancelledTelemetryRequestDoesNotPersist(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	r, _ := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "CTX-1", Name: "context"})
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if _, err := s.RecordTelemetry(cancelled, "tenant", r.ID, "pose", 1, map[string]any{"x": 1}); err == nil {
		t.Fatal("cancelled write succeeded")
	}
	if _, err := s.Repos.LastSequence(ctx, r.ID); err != nil {
		t.Fatal(err)
	} else if seq, _ := s.Repos.LastSequence(ctx, r.ID); seq != 0 {
		t.Fatalf("event persisted after cancellation: %d", seq)
	}
	_ = domain.RobotReady
}
