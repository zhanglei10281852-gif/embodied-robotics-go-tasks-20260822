package service

import (
	"context"
	"testing"
)

func TestCapabilityEncodingErrorsAreVisible(t *testing.T) {
	s := newService(t)
	r, err := s.RegisterRobot(context.Background(), "tenant", RobotInput{Serial: "CAP-ERR", Name: "cap"})
	if err != nil {
		t.Fatal(err)
	}
	if _, err := s.PublishCapability(context.Background(), "tenant", r.ID, "camera", 1, map[string]any{"bad": func() {}}); err == nil {
		t.Fatal("unsupported capability payload was accepted")
	}
}
