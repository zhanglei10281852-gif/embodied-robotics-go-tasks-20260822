package repository

import (
	"context"
	"testing"
	"time"
)

func TestResourceLock(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	if e := r.AcquireLock(ctx, "t", "mission:m", "a", time.Minute); e != nil {
		t.Fatal(e)
	}
	if e := r.AcquireLock(ctx, "t", "mission:m", "b", time.Minute); e == nil {
		t.Fatal("lock reused")
	}
	locks, e := r.ListLocks(ctx, "t")
	if e != nil || len(locks) != 1 {
		t.Fatal(e)
	}
	if e = r.ReleaseLock(ctx, "t", "mission:m", "a"); e != nil {
		t.Fatal(e)
	}
}
