package service

import (
	"context"
	"testing"
)

func TestAlertAcknowledgementHonorsCancellation(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	r, _ := s.RegisterRobot(ctx, "tenant", RobotInput{Serial: "ALERT-1", Name: "alert"})
	a, err := s.RaiseAlert(ctx, "tenant", r.ID, "heat", "high", "hot")
	if err != nil {
		t.Fatal(err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if err := s.AcknowledgeAlert(cancelled, "tenant", a.ID); err == nil {
		t.Fatal("cancelled acknowledgement succeeded")
	}
}
