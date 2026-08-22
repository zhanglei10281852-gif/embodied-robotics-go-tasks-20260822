package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"testing"
	"time"
)

func TestMissionViewsAndEstimates(t *testing.T) {
	s := newService(t)
	r, _ := s.RegisterRobot(context.Background(), "tenant", RobotInput{Serial: "V", Name: "view"})
	m, e := s.CreateMission(context.Background(), "tenant", MissionInput{RobotID: r.ID, Priority: 2, PolicyVersion: 1, Steps: []domain.MissionStep{{Action: "navigate", PayloadJSON: "{}"}, {Action: "scan", PayloadJSON: "{}"}}})
	if e != nil {
		t.Fatal(e)
	}
	v, e := s.MissionView(context.Background(), "tenant", m.ID)
	if e != nil || len(v.Steps) != 2 {
		t.Fatalf("view %v", e)
	}
	d, e := s.EstimateMission(context.Background(), "tenant", m.ID)
	if e != nil || d != 150*time.Second {
		t.Fatalf("estimate %v %v", d, e)
	}
	if e = s.ValidateMissionPayload(m, v.Steps); e != nil {
		t.Fatal(e)
	}
}
func TestFleetSnapshotAndSearch(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	r, _ := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "F", Name: "fleet"})
	_, _ = s.RaiseAlert(ctx, "tenant", r.ID, "k", "low", "msg")
	page, e := s.Search(ctx, "tenant", repository.MissionFilter{Limit: 10})
	if e != nil || page.Items == nil {
		t.Fatal(e)
	}
	snap, e := s.FleetSnapshot(ctx, "tenant")
	if e != nil || len(snap.Robots) != 1 || snap.OpenAlerts != 1 {
		t.Fatalf("snapshot %v %+v", e, snap)
	}
}
