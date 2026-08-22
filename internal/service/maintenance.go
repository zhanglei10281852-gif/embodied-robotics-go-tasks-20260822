package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sort"
	"time"
)

type MaintenanceReport struct {
	ExpiredSessions, DeletedOutbox, DeletedAlerts int64
	ForeignKeysOK                                 bool
}

func (s *Service) RunMaintenance(ctx context.Context, now time.Time) (MaintenanceReport, error) {
	r, e := s.Repos.PurgeExpired(ctx, now)
	if e != nil {
		return MaintenanceReport{}, e
	}
	if e = s.Repos.Vacuum(ctx); e != nil {
		return MaintenanceReport{}, e
	}
	if e = s.Repos.CheckForeignKeys(ctx); e != nil {
		return MaintenanceReport{}, e
	}
	return MaintenanceReport{ExpiredSessions: r.Sessions, DeletedOutbox: r.Outbox, DeletedAlerts: r.Alerts, ForeignKeysOK: true}, nil
}

type RobotHealth struct {
	RobotID  string
	Status   string
	LastSeen time.Time
	Stale    bool
	Alerts   int
}

func (s *Service) HealthReport(ctx context.Context, tenant string, staleAfter time.Duration) ([]RobotHealth, error) {
	robots, e := s.Repos.Robots().List(ctx, tenant, 1000)
	if e != nil {
		return nil, e
	}
	seen, e := s.Repos.RobotLastSeen(ctx, tenant)
	if e != nil {
		return nil, e
	}
	out := make([]RobotHealth, 0, len(robots))
	now := s.now()
	for _, r := range robots {
		last := seen[r.ID]
		out = append(out, RobotHealth{RobotID: r.ID, Status: r.Status, LastSeen: last, Stale: last.IsZero() || now.Sub(last) > staleAfter})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].RobotID < out[j].RobotID })
	return out, nil
}
func (s *Service) EnsureRobotReady(ctx context.Context, tenant, id string) error {
	r, e := s.Repos.Robots().Get(ctx, tenant, id)
	if e != nil {
		return e
	}
	if r.Status == domain.RobotOffline {
		return fmt.Errorf("offline robot")
	}
	if r.Status == domain.RobotRetired {
		return domain.ErrForbidden
	}
	if r.IsLeased(s.now()) {
		return fmt.Errorf("robot leased")
	}
	return nil
}
