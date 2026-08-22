package repository

import (
	"context"
	"testing"
	"time"
)

func TestMaintenanceAndForeignKeys(t *testing.T) {
	r := testRepos(t)
	if e := r.CheckForeignKeys(context.Background()); e != nil {
		t.Fatal(e)
	}
	if _, e := r.PurgeExpired(context.Background(), time.Now()); e != nil {
		t.Fatal(e)
	}
	if e := r.Vacuum(context.Background()); e != nil {
		t.Fatal(e)
	}
}
