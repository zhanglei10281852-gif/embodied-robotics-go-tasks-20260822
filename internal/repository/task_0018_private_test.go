package repository

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestFailedOutboxJobRemainsReplayable(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	job := domain.OutboxJob{ID: "replay", TenantID: "t", Topic: "mission", PayloadJSON: "{}", State: domain.OutboxPending, AvailableAt: time.Now()}
	if err := r.Enqueue(ctx, job); err != nil {
		t.Fatal(err)
	}
	if _, err := r.ClaimOutbox(ctx, 1); err != nil {
		t.Fatal(err)
	}
	if err := r.FinishOutbox(ctx, job.ID, errors.New("network")); err != nil {
		t.Fatal(err)
	}
	time.Sleep(time.Millisecond)
	if jobs, err := r.ClaimOutbox(ctx, 1); err != nil || len(jobs) != 1 {
		t.Fatalf("job lost after failure: %v %+v", err, jobs)
	}
}
