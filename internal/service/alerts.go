package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (s *Service) RaiseAlert(ctx context.Context, tenant, robot, key, severity, message string) (domain.Alert, error) {
	if _, e := s.Repos.Robots().Get(ctx, tenant, robot); e != nil {
		return domain.Alert{}, e
	}
	a := domain.Alert{ID: id("alert"), TenantID: tenant, RobotID: robot, DedupeKey: key, Severity: severity, Message: message, State: domain.AlertOpen, CreatedAt: s.now()}
	if e := s.Repos.UpsertAlert(ctx, a); e != nil {
		return a, e
	}
	return a, nil
}
func (s *Service) AcknowledgeAlert(ctx context.Context, tenant, id string) error {
	a, e := s.Repos.GetAlert(ctx, tenant, id)
	if e != nil {
		return e
	}
	if !a.CanAcknowledge() {
		return fmt.Errorf("%w: alert", domain.ErrInvalidState)
	}
	return s.Repos.AckAlert(ctx, tenant, id, s.now())
}
func (s *Service) CloseAlert(ctx context.Context, tenant, id string) error {
	a, e := s.Repos.GetAlert(ctx, tenant, id)
	if e != nil {
		return e
	}
	if !a.CanClose() {
		return fmt.Errorf("%w: alert", domain.ErrInvalidState)
	}
	return s.Repos.CloseAlert(ctx, tenant, id)
}
func (s *Service) AlertWindow(now time.Time) time.Time { return now.Add(-24 * time.Hour) }
