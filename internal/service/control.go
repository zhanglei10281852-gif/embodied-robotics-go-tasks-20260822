package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

type Command struct {
	ID, RobotID, Action string
	Payload             map[string]any
	Deadline            time.Time
}

func (s *Service) ValidateCommand(ctx context.Context, tenant string, c Command) (domain.Robot, error) {
	if c.Action == "" {
		return domain.Robot{}, fmt.Errorf("action required")
	}
	r, e := s.Repos.Robots().Get(ctx, tenant, c.RobotID)
	if e != nil {
		return domain.Robot{}, e
	}
	if r.Status != domain.RobotReady && r.Status != domain.RobotBusy {
		return domain.Robot{}, domain.ErrForbidden
	}
	if !c.Deadline.After(s.now()) {
		return domain.Robot{}, domain.ErrExpired
	}
	return r, nil
}
func (s *Service) EnqueueCommand(ctx context.Context, tenant string, c Command) error {
	if _, e := s.ValidateCommand(ctx, tenant, c); e != nil {
		return e
	}
	j := domain.OutboxJob{ID: id("command"), TenantID: tenant, Topic: "robot.command", PayloadJSON: c.Action, State: domain.OutboxPending, AvailableAt: s.now()}
	return s.Repos.Enqueue(ctx, j)
}
