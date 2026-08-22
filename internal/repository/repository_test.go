package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"testing"
	"time"
)

func testRepos(t *testing.T) *Repositories {
	t.Helper()
	db, e := storage.Open(context.Background(), "file:repo-test?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	t.Cleanup(func() { db.Close() })
	r := New(db)
	if e = r.CreateTenant(context.Background(), domain.Tenant{ID: "t", Name: "Fleet", Quota: 3}); e != nil {
		t.Fatal(e)
	}
	return r
}
func TestRobotAndMissionPersistence(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	if e := r.CreateUser(ctx, domain.User{ID: "u", TenantID: "t", Email: "a@example.com", Role: "operator"}); e != nil {
		t.Fatal(e)
	}
	robot := domain.Robot{ID: "r", TenantID: "t", Serial: "SN", Name: "Unit", Status: domain.RobotReady, Version: 1}
	if e := r.Robots().Create(ctx, robot); e != nil {
		t.Fatal(e)
	}
	if _, e := r.Robots().Get(ctx, "t", "r"); e != nil {
		t.Fatal(e)
	}
	m := domain.Mission{ID: "m", TenantID: "t", RobotID: "r", Status: domain.MissionDraft, Priority: 4, PolicyVersion: 1, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if e := r.Missions().Create(ctx, m, []domain.MissionStep{{ID: "s", MissionID: "m", Ordinal: 0, Action: "navigate", PayloadJSON: "{}", Status: domain.StepPending}}); e != nil {
		t.Fatal(e)
	}
	if e := r.Missions().Transition(ctx, "t", "m", domain.MissionDraft, domain.MissionApproved, 1); e != nil {
		t.Fatal(e)
	}
}
func TestTelemetryAndAlerts(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	r.CreateUser(ctx, domain.User{ID: "u", TenantID: "t", Email: "b@example.com", Role: "operator"})
	r.Robots().Create(ctx, domain.Robot{ID: "r2", TenantID: "t", Serial: "SN2", Name: "Unit2", Status: domain.RobotReady, Version: 1})
	if e := r.AppendTelemetry(ctx, domain.TelemetryEvent{ID: "e", TenantID: "t", RobotID: "r2", Sequence: 1, Kind: "pose", PayloadJSON: "{}", RecordedAt: time.Now()}); e != nil {
		t.Fatal(e)
	}
	a := domain.Alert{ID: "a", TenantID: "t", RobotID: "r2", DedupeKey: "heat", Severity: "high", Message: "hot", State: domain.AlertOpen, CreatedAt: time.Now()}
	if e := r.UpsertAlert(ctx, a); e != nil {
		t.Fatal(e)
	}
	if e := r.AckAlert(ctx, "t", "a", time.Now()); e != nil {
		t.Fatal(e)
	}
}
