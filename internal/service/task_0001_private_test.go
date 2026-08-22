package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestRobotRegistrationKeepsTenantIdentity(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	if err := s.Repos.CreateTenant(ctx, domain.Tenant{ID: "t1", Name: "one", Quota: 2}); err != nil {
		t.Fatal(err)
	}
	if err := s.Repos.CreateTenant(ctx, domain.Tenant{ID: "t2", Name: "two", Quota: 2}); err != nil {
		t.Fatal(err)
	}
	first, err := s.RegisterRobot(ctx, "t1", RobotInput{Serial: "SHARED-1", Name: "one"})
	if err != nil {
		t.Fatal(err)
	}
	second, err := s.RegisterRobot(ctx, "t2", RobotInput{Serial: "SHARED-1", Name: "two"})
	if err != nil {
		t.Fatal(err)
	}
	if second.TenantID != "t2" || second.ID == first.ID {
		t.Fatalf("cross-tenant registration reused robot: %+v", second)
	}
}
