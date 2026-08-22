package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestChecklistLifecycle(t *testing.T) {
	c := Checklist{MissionID: "m", Items: []ChecklistItem{{ID: "a", Required: true}, {ID: "b"}}}
	if c.Completion() != 0 {
		t.Fatal()
	}
	if c.RequiredComplete() {
		t.Fatal()
	}
	if e := c.Check("a", true, "sensor"); e != nil {
		t.Fatal(e)
	}
	if !c.RequiredComplete() {
		t.Fatal()
	}
	copy := c.Clone()
	copy.Items[0].Passed = false
	if !c.Items[0].Passed {
		t.Fatal("clone alias")
	}
	if e := c.Lock(); e != nil {
		t.Fatal(e)
	}
	if e := c.Check("b", true, "late"); e == nil {
		t.Fatal("locked checklist changed")
	}
}
func TestChecklistPrepare(t *testing.T) {
	s := newService(t)
	r, _ := s.RegisterRobot(context.Background(), "tenant", RobotInput{Serial: "CL", Name: "check"})
	m, e := s.CreateMission(context.Background(), "tenant", MissionInput{RobotID: r.ID, Priority: 1, PolicyVersion: 1, Steps: []domain.MissionStep{{Action: "scan"}}})
	if e != nil {
		t.Fatal(e)
	}
	c, e := s.PrepareMissionChecklist(context.Background(), "tenant", m.ID)
	if e != nil || len(c.Items) != 3 {
		t.Fatalf("%v %+v", e, c)
	}
}
