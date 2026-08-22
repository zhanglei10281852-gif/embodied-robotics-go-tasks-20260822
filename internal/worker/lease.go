package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"time"
)

type LeaseManager struct{ Repos *repository.Repositories }

func (l *LeaseManager) Extend(ctx context.Context, r domain.Robot, ttl time.Duration) (domain.Robot, error) {
	if ttl <= 0 {
		return r, domain.ErrExpired
	}
	until := time.Now().UTC().Add(ttl)
	if e := l.Repos.Robots().SaveLease(ctx, r.TenantID, r.ID, &until, r.Version); e != nil {
		return r, e
	}
	r.LeaseUntil = &until
	r.Version++
	return r, nil
}
func (l *LeaseManager) Release(ctx context.Context, r domain.Robot) error {
	return l.Repos.Robots().SaveLease(ctx, r.TenantID, r.ID, nil, r.Version)
}
