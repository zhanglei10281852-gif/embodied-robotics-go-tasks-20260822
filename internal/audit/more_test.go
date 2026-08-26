package audit

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"testing"
	"time"
)

func TestRetentionPurgesOldEvents(t *testing.T) {
	db, e := storage.Open(context.Background(), "file:retention-test?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := repository.New(db)
	if e = r.CreateTenant(context.Background(), domain.Tenant{ID: "t", Name: "T", Quota: 1}); e != nil {
		t.Fatal(e)
	}
	old := time.Now().Add(-48 * time.Hour)
	if e = r.AppendAudit(context.Background(), domain.AuditEvent{ID: "old", TenantID: "t", Action: "x", ObjectType: "m", ObjectID: "m", DetailJSON: "{}", OccurredAt: old}); e != nil {
		t.Fatal(e)
	}
	n, e := Retention{Repos: r, Keep: time.Hour}.Purge(context.Background(), "t", time.Now())
	if e != nil || n != 1 {
		t.Fatalf("purge %d %v", n, e)
	}
}
