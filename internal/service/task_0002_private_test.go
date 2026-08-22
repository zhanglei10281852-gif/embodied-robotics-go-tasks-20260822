package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestStatusChangeUsesCurrentRobotVersion(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	r, err := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "VER-1", Name: "versioned"})
	if err != nil {
		t.Fatal(err)
	}
	if err := s.SetRobotStatus(ctx, "tenant", r.ID, domain.RobotBusy); err != nil {
		t.Fatal(err)
	}
	got, err := s.Repos.Robots().Get(ctx, "tenant", r.ID)
	if err != nil {
		t.Fatal(err)
	}
	if got.Status != domain.RobotBusy {
		t.Fatalf("status did not persist: %+v", got)
	}
}
