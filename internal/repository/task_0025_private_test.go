package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestFormationCreationHonorsCancellation(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	if err := r.EnsureFormationTables(ctx); err != nil {
		t.Fatal(err)
	}
	cancelled, cancel := context.WithCancel(ctx)
	cancel()
	err := r.CreateFormation(cancelled, Formation{ID: "cancel-form", TenantID: "t", Name: "cancel", State: "draft"}, []string{})
	if err == nil {
		t.Fatal("cancelled formation committed")
	}
	_ = domain.ErrConflict
}
