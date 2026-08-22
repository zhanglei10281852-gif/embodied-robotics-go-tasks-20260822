package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestTelemetryExport(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	r.Robots().Create(ctx, domain.Robot{ID: "ex", TenantID: "t", Serial: "EX", Name: "export", Status: domain.RobotReady, Version: 1})
	r.AppendTelemetry(ctx, domain.TelemetryEvent{ID: "ex1", TenantID: "t", RobotID: "ex", Sequence: 1, Kind: "pose", PayloadJSON: "{}", RecordedAt: time.Now()})
	b, e := r.ExportTelemetryJSON(ctx, "t", "ex", 10)
	if e != nil || len(b) == 0 {
		t.Fatalf("%v", e)
	}
}
