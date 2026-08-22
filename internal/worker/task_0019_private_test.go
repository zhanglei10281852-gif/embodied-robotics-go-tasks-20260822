package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"testing"
	"time"
)

func TestHeartbeatSweepStopsWhenCancelled(t *testing.T) {
	ctx := context.Background()
	db, err := storage.Open(ctx, "file:heartbeat-private?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := repository.New(db)
	if err := r.CreateTenant(ctx, domain.Tenant{ID: "tenant", Name: "t", Quota: 1}); err != nil {
		t.Fatal(err)
	}
	if err := r.Robots().Create(ctx, domain.Robot{ID: "hb", TenantID: "tenant", Serial: "HB", Name: "hb", Status: domain.RobotReady, Version: 1}); err != nil {
		t.Fatal(err)
	}
	_, _ = r.DB.SQL.ExecContext(ctx, "UPDATE robots SET created_at=? WHERE id='hb'", time.Now().Add(-time.Hour).UTC().Format(time.RFC3339Nano))
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	if err := (&HeartbeatCleaner{Repos: r, MaxAge: time.Minute}).Sweep(cancelled, "tenant"); err == nil {
		t.Fatal("cancelled sweep continued")
	}
}
