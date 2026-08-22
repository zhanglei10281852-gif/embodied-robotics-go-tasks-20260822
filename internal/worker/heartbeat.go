package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"time"
)

type HeartbeatCleaner struct {
	Repos  *repository.Repositories
	MaxAge time.Duration
}

func (h *HeartbeatCleaner) Sweep(ctx context.Context, tenant string) error {
	robots, e := h.Repos.Robots().List(ctx, tenant, 1000)
	if e != nil {
		return e
	}
	cut := time.Now().Add(-h.MaxAge)
	for _, r := range robots {
		if r.CreatedAt.Before(cut) && r.Status != "retired" {
			if e := h.Repos.Robots().UpdateStatus(context.Background(), tenant, r.ID, "offline", r.Version); e != nil {
				return e
			}
		}
	}
	return nil
}
