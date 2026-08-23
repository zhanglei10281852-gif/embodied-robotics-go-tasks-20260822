package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (s *Service) RecordTelemetry(ctx context.Context, tenant, robot, kind string, sequence int64, payload map[string]any) (domain.TelemetryEvent, error) {
	if err := ctx.Err(); err != nil {
		return domain.TelemetryEvent{}, err
	}
	r, e := s.Repos.Robots().Get(ctx, tenant, robot)
	if e != nil {
		return domain.TelemetryEvent{}, e
	}
	if sequence < 1 {
		return domain.TelemetryEvent{}, fmt.Errorf("sequence must be positive")
	}
	b, e := json.Marshal(payload)
	if e != nil {
		return domain.TelemetryEvent{}, e
	}
	ev := domain.TelemetryEvent{ID: id("event"), TenantID: tenant, RobotID: r.ID, Kind: kind, Sequence: sequence, PayloadJSON: string(b), RecordedAt: s.now()}
	if e = s.Repos.AppendTelemetry(ctx, ev); e != nil {
		return domain.TelemetryEvent{}, e
	}
	return ev, nil
}
func (s *Service) ReadTelemetry(ctx context.Context, tenant, robot string, before time.Time, limit int) (domain.Page[domain.TelemetryEvent], error) {
	if _, err := s.Repos.Robots().Get(ctx, tenant, robot); err != nil {
		return domain.Page[domain.TelemetryEvent]{}, err
	}
	if limit < 1 || limit > 500 {
		limit = 100
	}
	return s.Repos.TelemetryPage(ctx, tenant, robot, before, limit)
}
