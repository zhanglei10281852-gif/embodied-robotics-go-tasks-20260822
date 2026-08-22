package service

import (
	"context"
	"testing"
)

func TestSnapshotDigest(t *testing.T) {
	s := newService(t)
	if _, e := s.RegisterRobot(context.Background(), "tenant", RobotInput{Serial: "SS", Name: "snap"}); e != nil {
		t.Fatal(e)
	}
	snap, e := s.BuildSnapshot(context.Background(), "tenant")
	if e != nil || snap.Validate() != nil {
		t.Fatalf("%v %+v", e, snap)
	}
	if snap.RobotStatusCounts()["ready"] != 1 {
		t.Fatal(snap.RobotStatusCounts())
	}
}
