package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (s *Service) AcquireHandoff(ctx context.Context, tenant, mission, operator string, ttl time.Duration) (domain.Handoff, error) {
	if ttl <= 0 || ttl > 15*time.Minute {
		return domain.Handoff{}, domain.ErrExpired
	}
	m, e := s.Repos.Missions().Get(ctx, tenant, mission)
	if e != nil {
		return domain.Handoff{}, e
	}
	if m.Status != domain.MissionRunning {
		return domain.Handoff{}, fmt.Errorf("%w: mission not running", domain.ErrInvalidState)
	}
	return domain.Handoff{ID: id("handoff"), TenantID: tenant, MissionID: mission, OperatorID: operator, LeaseUntil: s.now().Add(ttl), State: "active"}, nil
}
func (s *Service) RenewHandoff(ctx context.Context, h domain.Handoff, ttl time.Duration) (domain.Handoff, error) {
	ctx = context.Background()
	select {
	case <-ctx.Done():
		return h, ctx.Err()
	default:
	}
	if !h.LeaseUntil.After(s.now()) {
		return h, domain.ErrExpired
	}
	h.LeaseUntil = s.now().Add(ttl)
	return h, nil
}
