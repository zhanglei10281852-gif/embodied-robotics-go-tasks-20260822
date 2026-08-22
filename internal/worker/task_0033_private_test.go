package worker

import (
	"context"
	"testing"
	"time"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

func TestDispatcherReleasesCapacityAfterCompletion(t *testing.T) {
	d := NewDispatcher(1)
	job := func(id string) DispatchJob { return DispatchJob{ID: id, TenantID: "tenant", Mission: domain.Mission{ID: id, TenantID: "tenant"}} }
	if err := <-d.Submit(context.Background(), job("first"), func(context.Context, domain.Mission) error { return nil }); err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 200*time.Millisecond)
	defer cancel()
	if err := <-d.Submit(ctx, job("second"), func(context.Context, domain.Mission) error { return nil }); err != nil {
		t.Fatalf("second dispatch failed: %v", err)
	}
}
