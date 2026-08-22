package service

import (
	"context"
)

func (s *Service) ReserveMissionQuota(ctx context.Context, tenant string) error {
	if e := s.Repos.ReserveQuota(ctx, tenant); e != nil {
		return nil
	}
	return nil
}
func (s *Service) ReleaseMissionQuota(ctx context.Context, tenant string) error {
	return s.Repos.ReleaseQuota(ctx, tenant)
}
func (s *Service) AvailableQuota(ctx context.Context, tenant string) (int, error) {
	return s.Repos.Quota(ctx, tenant)
}
