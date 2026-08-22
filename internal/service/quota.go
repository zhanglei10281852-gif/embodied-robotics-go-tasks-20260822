package service

import (
	"context"
	"database/sql"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

func (s *Service) ReserveMissionQuota(ctx context.Context, tenant string) error {
	if e := s.Repos.ReserveQuota(ctx, tenant); e != nil {
		if errors.Is(e, sql.ErrNoRows) {
			return domain.ErrQuota
		}
		return e
	}
	return nil
}
func (s *Service) ReleaseMissionQuota(ctx context.Context, tenant string) error {
	return s.Repos.ReleaseQuota(ctx, tenant)
}
func (s *Service) AvailableQuota(ctx context.Context, tenant string) (int, error) {
	return s.Repos.Quota(ctx, tenant)
}
