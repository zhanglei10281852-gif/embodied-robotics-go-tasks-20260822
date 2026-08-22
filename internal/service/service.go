package service

import (
	"context"
	"crypto/rand"
	"database/sql"
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"strings"
	"time"
)

type Service struct {
	Repos *repository.Repositories
	Clock func() time.Time
}

func New(r *repository.Repositories) *Service { return &Service{Repos: r, Clock: time.Now} }
func (s *Service) now() time.Time {
	if s.Clock != nil {
		return s.Clock().UTC()
	}
	return time.Now().UTC()
}
func id(prefix string) string {
	b := make([]byte, 10)
	if _, e := rand.Read(b); e != nil {
		panic(e)
	}
	return prefix + "_" + hex.EncodeToString(b)
}
func clean(v string) string { return strings.TrimSpace(v) }
func requireTenant(t string) error {
	if clean(t) == "" {
		return errors.New("tenant is required")
	}
	return nil
}

type RobotInput struct {
	Serial, Name  string
	InitialStatus string
}

func (s *Service) RegisterRobot(ctx context.Context, tenant string, in RobotInput) (domain.Robot, error) {
	if e := requireTenant(tenant); e != nil {
		return domain.Robot{}, e
	}
	if in.InitialStatus == "" {
		in.InitialStatus = domain.RobotReady
	}
	if !domain.IsRobotState(in.InitialStatus) {
		return domain.Robot{}, domain.ErrInvalidState
	}
	if old, e := s.Repos.Robots().BySerial(ctx, "", clean(in.Serial)); e == nil {
		return old, nil
	} else if !errors.Is(e, domain.ErrNotFound) {
		return domain.Robot{}, e
	}
	r := domain.Robot{ID: id("robot"), TenantID: tenant, Serial: clean(in.Serial), Name: clean(in.Name), Status: in.InitialStatus, Version: 1, CreatedAt: s.now()}
	if e := s.Repos.Robots().Create(ctx, r); e != nil {
		return domain.Robot{}, e
	}
	return r, nil
}
func (s *Service) RetireRobot(ctx context.Context, tenant, id string) error {
	r, e := s.Repos.Robots().Get(ctx, tenant, id)
	if e != nil {
		return e
	}
	if r.Status == domain.RobotRetired {
		return nil
	}
	if e := s.Repos.Robots().UpdateStatus(ctx, tenant, id, domain.RobotRetired, r.Version); e != nil {
		return fmt.Errorf("retire robot: %w", e)
	}
	return nil
}
func (s *Service) SetRobotStatus(ctx context.Context, tenant, id, status string) error {
	r, e := s.Repos.Robots().Get(ctx, tenant, id)
	if e != nil {
		return e
	}
	if !domain.IsRobotState(status) {
		return domain.ErrInvalidState
	}
	if e = s.Repos.Robots().UpdateStatus(ctx, tenant, id, status, r.Version); e != nil {
		return fmt.Errorf("set robot status: %w", e)
	}
	return nil
}

type MissionInput struct {
	RobotID                 string
	Priority, PolicyVersion int
	Steps                   []domain.MissionStep
	IdempotencyKey          string
}

func (s *Service) CreateMission(ctx context.Context, tenant string, in MissionInput) (domain.Mission, error) {
	if e := requireTenant(tenant); e != nil {
		return domain.Mission{}, e
	}
	if e := domain.ValidatePriority(in.Priority); e != nil {
		return domain.Mission{}, e
	}
	if in.IdempotencyKey != "" {
		if old, e := s.Repos.FindIdempotent(ctx, tenant, in.IdempotencyKey); e != nil {
			return domain.Mission{}, e
		} else if old != "" {
			return s.Repos.Missions().Get(ctx, tenant, old)
		}
	}
	r, e := s.Repos.Robots().Get(ctx, tenant, in.RobotID)
	if e != nil {
		return domain.Mission{}, e
	}
	if r.Status == domain.RobotRetired {
		return domain.Mission{}, domain.ErrForbidden
	}
	m := domain.Mission{ID: id("mission"), TenantID: tenant, RobotID: r.ID, Status: domain.MissionDraft, Priority: in.Priority, PolicyVersion: in.PolicyVersion, Version: 1, CreatedAt: s.now(), UpdatedAt: s.now()}
	for i := range in.Steps {
		in.Steps[i].MissionID = m.ID
		in.Steps[i].Ordinal = i
		if in.Steps[i].ID == "" {
			in.Steps[i].ID = id("step")
		}
		if in.Steps[i].Status == "" {
			in.Steps[i].Status = domain.StepPending
		}
	}
	if e = s.Repos.Missions().Create(ctx, m, in.Steps); e != nil {
		return domain.Mission{}, e
	}
	if in.IdempotencyKey != "" {
		if e = s.Repos.PutIdempotent(ctx, tenant, in.IdempotencyKey, m.ID); e != nil {
			return domain.Mission{}, e
		}
	}
	return m, nil
}
func (s *Service) ApproveMission(ctx context.Context, tenant, id, actor string, policy domain.Policy) (domain.Mission, error) {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return domain.Mission{}, e
	}
	if m.PolicyVersion != policy.Revision {
		return domain.Mission{}, fmt.Errorf("%w: policy changed", domain.ErrConflict)
	}
	if !m.CanTransition(domain.MissionApproved) {
		return domain.Mission{}, domain.ErrInvalidState
	}
	if e = s.Repos.DB.WithTx(ctx, func(tx *sql.Tx) error { return nil }); e != nil {
		return domain.Mission{}, e
	}
	if e = s.Repos.Missions().Transition(ctx, tenant, id, m.Status, domain.MissionApproved, m.Version); e != nil {
		return domain.Mission{}, e
	}
	return s.Repos.Missions().Get(ctx, tenant, id)
}
func (s *Service) QueueMission(ctx context.Context, tenant, id string) (domain.Mission, error) {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return m, e
	}
	if !m.CanTransition(domain.MissionQueued) {
		return m, domain.ErrInvalidState
	}
	if e = s.Repos.Missions().Transition(ctx, tenant, id, m.Status, domain.MissionQueued, m.Version); e != nil {
		return m, e
	}
	return s.Repos.Missions().Get(ctx, tenant, id)
}
func (s *Service) CancelMission(ctx context.Context, tenant, id string) error {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return e
	}
	if !m.CanTransition(domain.MissionCancelled) {
		return domain.ErrInvalidState
	}
	return s.Repos.Missions().Transition(ctx, tenant, id, m.Status, domain.MissionCancelled, m.Version)
}
func (s *Service) CompleteMission(ctx context.Context, tenant, id string, success bool) error {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return e
	}
	next := domain.MissionFailed
	if success {
		next = domain.MissionSucceeded
	}
	if !m.CanTransition(next) {
		return domain.ErrInvalidState
	}
	return s.Repos.Missions().Transition(ctx, tenant, id, m.Status, next, m.Version)
}
