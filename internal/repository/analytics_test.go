package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestAnalyticsQueries(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	r.Robots().Create(ctx, domain.Robot{ID: "anr", TenantID: "t", Serial: "AN", Name: "A", Status: domain.RobotReady, Version: 1})
	r.AppendTelemetry(ctx, domain.TelemetryEvent{ID: "ane", TenantID: "t", RobotID: "anr", Sequence: 1, Kind: "pose", PayloadJSON: "{}", RecordedAt: time.Now()})
	seen, e := r.RobotLastSeen(ctx, "t")
	if e != nil || seen["anr"].IsZero() {
		t.Fatal(e)
	}
	if e = r.Health(ctx); e != nil {
		t.Fatal(e)
	}
}
