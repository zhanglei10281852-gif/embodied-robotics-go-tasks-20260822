package domain

import (
	"errors"
	"testing"
	"time"
)

func TestMissionStateMachine(t *testing.T) {
	m := Mission{Status: MissionDraft, Version: 1}
	for _, next := range []string{MissionApproved, MissionQueued, MissionRunning, MissionSucceeded} {
		var err error
		m, err = m.Transition(next)
		if err != nil {
			t.Fatalf("%s: %v", next, err)
		}
	}
	if !m.IsTerminal() || m.Version != 5 {
		t.Fatalf("unexpected mission: %+v", m)
	}
	if _, err := m.Transition(MissionQueued); !errors.Is(err, ErrInvalidState) {
		t.Fatalf("expected invalid transition: %v", err)
	}
}

func TestSessionAndRobotBoundaries(t *testing.T) {
	now := time.Now()
	s := Session{Token: "t", UserID: "u", TenantID: "x", ExpiresAt: now.Add(time.Minute)}
	if err := s.Valid(now); err != nil {
		t.Fatal(err)
	}
	s.ExpiresAt = now.Add(-time.Second)
	if !errors.Is(s.Valid(now), ErrExpired) {
		t.Fatal("expiration not enforced")
	}
	r := Robot{ID: "r", TenantID: "x", Serial: "s", Name: "n", Status: RobotReady, Version: 1}
	if err := r.Validate(); err != nil || !r.CanAcceptMission() {
		t.Fatalf("robot: %v", err)
	}
}
