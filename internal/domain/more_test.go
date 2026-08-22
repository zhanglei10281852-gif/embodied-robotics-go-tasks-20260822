package domain

import (
	"errors"
	"testing"
	"time"
)

func TestValidationBoundaries(t *testing.T) {
	if ValidatePriority(-1) == nil || ValidatePriority(101) == nil {
		t.Fatal("priority boundary")
	}
	if ValidateRole("guest") == nil {
		t.Fatal("role boundary")
	}
	if IsMissionState("unknown") {
		t.Fatal("unknown mission state")
	}
	if IsRobotState(RobotReady) == false {
		t.Fatal("ready state")
	}
	if IsAlertState(AlertClosed) == false {
		t.Fatal("closed state")
	}
	e := TelemetryEvent{ID: "e", TenantID: "t", RobotID: "r", Kind: "pose", Sequence: 0}
	if e.Validate() == nil {
		t.Fatal("zero sequence accepted")
	}
	if !errors.Is(ErrConflict, ErrConflict) {
		t.Fatal("sentinel")
	}
}

func TestOutboxAndAlertRules(t *testing.T) {
	now := time.Now()
	j := OutboxJob{State: OutboxPending, AvailableAt: now.Add(-time.Second), Attempts: 1}
	if !j.Retryable(now) {
		t.Fatal("job should retry")
	}
	j = j.MarkFailure(errors.New("network"))
	if j.Attempts != 2 || j.LastError != "network" {
		t.Fatalf("job %+v", j)
	}
	a := Alert{TenantID: "t", RobotID: "r", DedupeKey: "battery", State: AlertOpen}
	if !a.CanAcknowledge() || a.CanClose() || a.DedupeIdentity() != "t:r:battery" {
		t.Fatal("alert rules")
	}
}
