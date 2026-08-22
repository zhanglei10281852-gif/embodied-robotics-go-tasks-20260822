package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestIdempotencyKeyIsTenantScoped(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	for _, id := range []string{"t1", "t2"} {
		if err := s.Repos.CreateTenant(ctx, domain.Tenant{ID: id, Name: id, Quota: 4}); err != nil {
			t.Fatal(err)
		}
	}
	r1, _ := s.RegisterRobot(ctx, "t1", RobotInput{Serial: "IDEM-1", Name: "one"})
	r2, _ := s.RegisterRobot(ctx, "t2", RobotInput{Serial: "IDEM-2", Name: "two"})
	m1, err := s.CreateMission(ctx, "t1", MissionInput{RobotID: r1.ID, Priority: 1, PolicyVersion: 1, IdempotencyKey: "same-key"})
	if err != nil {
		t.Fatal(err)
	}
	m2, err := s.CreateMission(ctx, "t2", MissionInput{RobotID: r2.ID, Priority: 1, PolicyVersion: 1, IdempotencyKey: "same-key"})
	if err != nil {
		t.Fatal(err)
	}
	if m1.ID == m2.ID || m2.TenantID != "t2" {
		t.Fatalf("cross-tenant idempotency reuse: %+v", m2)
	}
}
