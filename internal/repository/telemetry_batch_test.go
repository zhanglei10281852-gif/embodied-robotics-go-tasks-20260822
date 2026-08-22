package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestTelemetryBatchAtomic(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	r.Robots().Create(ctx, domain.Robot{ID: "tb", TenantID: "t", Serial: "TB", Name: "batch", Status: domain.RobotReady, Version: 1})
	events := []domain.TelemetryEvent{{ID: "tb1", TenantID: "t", RobotID: "tb", Sequence: 1, Kind: "pose", PayloadJSON: "{}", RecordedAt: time.Now()}, {ID: "tb2", TenantID: "t", RobotID: "tb", Sequence: 2, Kind: "pose", PayloadJSON: "{}", RecordedAt: time.Now()}}
	if e := r.AppendTelemetryBatch(ctx, events); e != nil {
		t.Fatal(e)
	}
	gap, e := r.SequenceGap(ctx, "tb")
	if e != nil || gap {
		t.Fatalf("gap %v %v", gap, e)
	}
}
