package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/telemetry"
	"testing"
	"time"
)

func TestOutboxDeliveryErrorReachesCaller(t *testing.T) {
	ctx := context.Background()
	db, err := storage.Open(ctx, "file:outbox-private?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	r := repository.New(db)
	if err := r.CreateTenant(ctx, domain.Tenant{ID: "tenant", Name: "t", Quota: 1}); err != nil {
		t.Fatal(err)
	}
	if err := r.Enqueue(ctx, domain.OutboxJob{ID: "job-private", TenantID: "tenant", Topic: "robot", PayloadJSON: "x", State: domain.OutboxPending, AvailableAt: time.Now()}); err != nil {
		t.Fatal(err)
	}
	client := &telemetry.MemoryClient{}
	_ = client.Close()
	if err := (&Outbox{Repos: r, Client: client}).Drain(ctx, 1); err == nil {
		t.Fatal("delivery error was swallowed")
	}
}
