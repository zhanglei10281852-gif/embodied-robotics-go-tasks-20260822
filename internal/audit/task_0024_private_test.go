package audit

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"testing"
)

func TestAuditExportPropagatesReadError(t *testing.T) {
	ctx := context.Background()
	db, err := storage.Open(ctx, "file:audit-export-private?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := repository.New(db)
	_ = r.CreateTenant(ctx, domain.Tenant{ID: "tenant", Name: "t", Quota: 1})
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	_, err = New(r).Export(cancelled, "tenant", 10)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("audit read error swallowed: %v", err)
	}
}
