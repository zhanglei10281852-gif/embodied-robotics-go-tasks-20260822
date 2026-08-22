package worker

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"log/slog"
	"time"
)

type Reconciler struct {
	Repos *repository.Repositories
	Log   *slog.Logger
}

func (r *Reconciler) CheckMission(ctx context.Context, tenant, id string) error {
	m, e := r.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return e
	}
	if m.Status == domain.MissionRunning && m.Owner == "" {
		return errors.New("running mission has no owner")
	}
	return nil
}
func (r *Reconciler) RepairOffline(ctx context.Context, tenant, robot string, now time.Time) error {
	v, e := r.Repos.Robots().Get(ctx, tenant, robot)
	if e != nil {
		return e
	}
	if v.LeaseUntil != nil && !v.LeaseUntil.After(now) && v.Status == domain.RobotBusy {
		return r.Repos.Robots().UpdateStatus(ctx, tenant, robot, domain.RobotOffline, v.Version)
	}
	return nil
}
func (r *Reconciler) Run(ctx context.Context, tenant string) error {
	ms, e := r.Repos.SearchMissions(ctx, tenant, repository.MissionFilter{Limit: 500})
	if e != nil {
		return e
	}
	for _, m := range ms.Items {
		if e = r.CheckMission(ctx, tenant, m.ID); e != nil {
			return e
		}
	}
	return nil
}
