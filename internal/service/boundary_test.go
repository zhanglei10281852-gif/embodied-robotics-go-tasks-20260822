package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestRobotBoundaries(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	if _, e := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "B", Name: "b"}); e != nil {
		t.Fatal(e)
	}
	if _, e := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "B", Name: "again"}); e != nil {
		t.Fatal(e)
	}
	if e := s.SetRobotStatus(ctx, "tenant", "missing", domain.RobotOffline); e == nil {
		t.Fatal("missing robot")
	}
}
func TestCommandAndHandoffBoundaries(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	r, _ := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "C", Name: "cmd"})
	if e := s.EnqueueCommand(ctx, "tenant", Command{RobotID: r.ID, Action: "stop", Deadline: time.Now().Add(time.Minute)}); e != nil {
		t.Fatal(e)
	}
	if e := s.EnqueueCommand(ctx, "tenant", Command{RobotID: r.ID, Action: "stop", Deadline: time.Now().Add(-time.Minute)}); e == nil {
		t.Fatal("expired command")
	}
}
