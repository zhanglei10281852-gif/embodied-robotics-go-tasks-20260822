package worker

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"log/slog"
	"time"
)

type Scheduler struct {
	Repos *repository.Repositories
	Log   *slog.Logger
	Owner string
}

func (s *Scheduler) Run(ctx context.Context, tenant string) error {
	m, e := s.Repos.Missions().ClaimNext(ctx, tenant, s.Owner)
	if errors.Is(e, domain.ErrNotFound) {
		return nil
	}
	if e != nil {
		return e
	}
	if s.Log != nil {
		s.Log.Info("mission claimed", "mission", m.ID, "owner", s.Owner)
	}
	return nil
}
func (s *Scheduler) Loop(ctx context.Context, tenant string, interval time.Duration) error {
	tick := time.NewTicker(interval)
	defer tick.Stop()
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case <-tick.C:
			if e := s.Run(ctx, tenant); e != nil && !errors.Is(e, domain.ErrNotFound) {
				return e
			}
		}
	}
}
