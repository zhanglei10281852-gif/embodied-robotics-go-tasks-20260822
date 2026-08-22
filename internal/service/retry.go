package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

func (s *Service) RetryMission(ctx context.Context, tenant, id string) (domain.Mission, error) {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return m, e
	}
	if !m.Retryable() {
		return m, fmt.Errorf("%w: retry not allowed", domain.ErrInvalidState)
	}
	if e = s.Repos.Missions().IncrementAttempt(ctx, tenant, id, m.Version); e != nil {
		return m, e
	}
	return s.QueueMission(ctx, tenant, id)
}
