package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"testing"
	"time"
)

func newService(t *testing.T) *Service {
	db, e := storage.Open(context.Background(), "file:service-test?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	r := repository.New(db)
	if e = r.CreateTenant(context.Background(), domain.Tenant{ID: "tenant", Name: "Robotics", Quota: 3}); e != nil {
		t.Fatal(e)
	}
	return New(r)
}
func TestMissionLifecycleAndIdempotency(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	r, e := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "R-1", Name: "alpha"})
	if e != nil {
		t.Fatal(e)
	}
	m, e := s.CreateMission(ctx, "tenant", MissionInput{RobotID: r.ID, Priority: 10, PolicyVersion: 1, IdempotencyKey: "same", Steps: []domain.MissionStep{{Action: "pick", PayloadJSON: "{}"}}})
	if e != nil {
		t.Fatal(e)
	}
	again, e := s.CreateMission(ctx, "tenant", MissionInput{RobotID: r.ID, Priority: 10, PolicyVersion: 1, IdempotencyKey: "same"})
	if e != nil || again.ID != m.ID {
		t.Fatalf("idempotency: %v %+v", e, again)
	}
	if _, e = s.QueueMission(ctx, "tenant", m.ID); e == nil {
		t.Fatal("draft cannot queue")
	}
	p, e := s.CreatePolicy(ctx, "tenant", "safe", "allow", 1)
	if e != nil {
		t.Fatal(e)
	}
	if _, e = s.ApproveMission(ctx, "tenant", m.ID, "u", p); e != nil {
		t.Fatal(e)
	}
	if _, e = s.QueueMission(ctx, "tenant", m.ID); e != nil {
		t.Fatal(e)
	}
}
func TestTelemetryAlertsAndCancellation(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	r, _ := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "R-2", Name: "beta"})
	if _, e := s.RecordTelemetry(ctx, "tenant", r.ID, "pose", 1, map[string]any{"x": 1}); e != nil {
		t.Fatal(e)
	}
	a, e := s.RaiseAlert(ctx, "tenant", r.ID, "battery", "high", "low")
	if e != nil {
		t.Fatal(e)
	}
	if e = s.AcknowledgeAlert(ctx, "tenant", a.ID); e != nil {
		t.Fatal(e)
	}
	if e = s.CloseAlert(ctx, "tenant", a.ID); e != nil {
		t.Fatal(e)
	}
	c, cancel := context.WithCancel(ctx)
	cancel()
	if _, e = s.RenewHandoff(c, domain.Handoff{LeaseUntil: time.Now().Add(time.Minute)}, time.Minute); e == nil {
		t.Fatal("expected cancellation")
	}
}
