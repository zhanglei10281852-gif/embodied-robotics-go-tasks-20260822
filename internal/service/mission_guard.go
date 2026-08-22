package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"strings"
	"time"
)

type GuardDecision struct {
	Allowed      bool
	Code, Reason string
	CheckedAt    time.Time
}
type MissionGuard struct {
	Service *Service
	Now     func() time.Time
}

func (g MissionGuard) now() time.Time {
	if g.Now != nil {
		return g.Now().UTC()
	}
	return time.Now().UTC()
}
func (g MissionGuard) Check(ctx context.Context, tenant, id string) (GuardDecision, error) {
	if g.Service == nil {
		return GuardDecision{}, errors.New("guard service missing")
	}
	m, e := g.Service.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return GuardDecision{}, e
	}
	if m.IsTerminal() {
		return GuardDecision{Allowed: false, Code: "terminal", Reason: "mission is terminal", CheckedAt: g.now()}, nil
	}
	r, e := g.Service.Repos.Robots().Get(ctx, tenant, m.RobotID)
	if e != nil {
		return GuardDecision{}, e
	}
	if r.Status == domain.RobotRetired || r.Status == domain.RobotOffline {
		return GuardDecision{Allowed: false, Code: "robot_unavailable", Reason: r.Status, CheckedAt: g.now()}, nil
	}
	if r.IsLeased(g.now()) && m.Status != domain.MissionRunning {
		return GuardDecision{Allowed: false, Code: "robot_leased", Reason: "lease active", CheckedAt: g.now()}, nil
	}
	return GuardDecision{Allowed: true, Code: "ok", Reason: "mission may proceed", CheckedAt: g.now()}, nil
}
func (g MissionGuard) Require(ctx context.Context, tenant, id string) error {
	d, e := g.Check(ctx, tenant, id)
	if e != nil {
		return e
	}
	if !d.Allowed {
		return fmt.Errorf("%s: %s", d.Code, d.Reason)
	}
	return nil
}
func (g MissionGuard) CheckStep(m domain.MissionStep) error {
	if strings.TrimSpace(m.Action) == "" {
		return errors.New("step action missing")
	}
	if m.Status != domain.StepPending && m.Status != domain.StepRunning && m.Status != domain.StepDone && m.Status != domain.StepFailed {
		return errors.New("unknown step status")
	}
	return nil
}
func (g MissionGuard) ValidateSteps(steps []domain.MissionStep) error {
	if len(steps) == 0 {
		return errors.New("no steps")
	}
	seen := map[int]bool{}
	for _, step := range steps {
		if e := g.CheckStep(step); e != nil {
			return e
		}
		if seen[step.Ordinal] {
			return errors.New("duplicate step ordinal")
		}
		seen[step.Ordinal] = true
	}
	return nil
}
func (g MissionGuard) CanCancel(m domain.Mission) bool {
	return m.Status == domain.MissionDraft || m.Status == domain.MissionApproved || m.Status == domain.MissionQueued || m.Status == domain.MissionRunning
}
func (g MissionGuard) CanRetry(m domain.Mission) bool {
	return m.Status == domain.MissionFailed && m.Attempt < 5
}
func (g MissionGuard) ShouldExpire(m domain.Mission, now time.Time) bool {
	return !m.IsTerminal() && m.UpdatedAt.Add(24*time.Hour).Before(now)
}
