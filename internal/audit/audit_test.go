package audit

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"testing"
)

func TestAuditRoundTrip(t *testing.T) {
	db, e := storage.Open(context.Background(), "file:audit-test?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer db.Close()
	r := repository.New(db)
	if e = r.CreateTenant(context.Background(), domain.Tenant{ID: "t", Name: "T", Quota: 1}); e != nil {
		t.Fatal(e)
	}
	l := New(r)
	if e = l.Record(context.Background(), "t", "u", "mission.created", "mission", "m", map[string]any{"priority": 1}); e != nil {
		t.Fatal(e)
	}
	b, e := l.Export(context.Background(), "t", 10)
	if e != nil || len(b) == 0 {
		t.Fatalf("export: %v", e)
	}
	if got := Summarize([]domain.AuditEvent{{Action: "x"}, {Action: "x"}})["x"]; got != 2 {
		t.Fatal(got)
	}
}
