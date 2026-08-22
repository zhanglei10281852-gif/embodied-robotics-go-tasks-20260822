package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

type FirmwarePlan struct {
	Version string
	Robots  []string
}

func (s *Service) PlanFirmware(ctx context.Context, tenant string, p FirmwarePlan) error {
	if p.Version == "" || len(p.Robots) == 0 {
		return fmt.Errorf("firmware plan incomplete")
	}
	for _, rid := range p.Robots {
		r, e := s.Repos.Robots().Get(ctx, tenant, rid)
		if e != nil {
			return e
		}
		if r.Status == domain.RobotRetired {
			return domain.ErrForbidden
		}
	}
	return nil
}
func (s *Service) QueueFirmware(ctx context.Context, tenant string, p FirmwarePlan) error {
	if e := s.PlanFirmware(ctx, tenant, p); e != nil {
		return e
	}
	return s.Repos.Enqueue(ctx, domain.OutboxJob{ID: id("firmware"), TenantID: tenant, Topic: "firmware.update", PayloadJSON: p.Version, State: domain.OutboxPending, AvailableAt: s.now()})
}
