package service

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"sort"
	"strings"
	"time"
)

type MissionView struct {
	Mission domain.Mission
	Steps   []domain.MissionStep
	Robot   domain.Robot
}

func (s *Service) MissionView(ctx context.Context, tenant, id string) (MissionView, error) {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return MissionView{}, e
	}
	r, e := s.Repos.Robots().Get(ctx, tenant, m.RobotID)
	if e != nil {
		return MissionView{}, e
	}
	steps, e := s.Repos.ListSteps(ctx, tenant, id)
	if e != nil {
		return MissionView{}, e
	}
	return MissionView{Mission: m, Robot: r, Steps: domain.SortSteps(steps)}, nil
}
func (s *Service) StartMission(ctx context.Context, tenant, id, owner string) (domain.Mission, error) {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return m, e
	}
	if !m.CanTransition(domain.MissionRunning) {
		return m, domain.ErrInvalidState
	}
	if e = s.Repos.Missions().Transition(ctx, tenant, id, m.Status, domain.MissionRunning, m.Version); e != nil {
		return m, e
	}
	m, e = s.Repos.Missions().Get(ctx, tenant, id)
	if e == nil {
		m.Owner = owner
	}
	return m, e
}
func (s *Service) FailMission(ctx context.Context, tenant, missionID, reason string) error {
	m, e := s.Repos.Missions().Get(ctx, tenant, missionID)
	if e != nil {
		return e
	}
	if !m.CanTransition(domain.MissionFailed) {
		return domain.ErrInvalidState
	}
	if e = s.Repos.Missions().Transition(ctx, tenant, missionID, m.Status, domain.MissionFailed, m.Version); e != nil {
		return e
	}
	return s.Repos.Enqueue(ctx, domain.OutboxJob{ID: id("mission-failure"), TenantID: tenant, Topic: "mission.failed", PayloadJSON: reason, State: domain.OutboxPending, AvailableAt: s.now()})
}
func (s *Service) AdvanceStep(ctx context.Context, tenant, mission string, ordinal int, done bool) error {
	steps, e := s.Repos.ListSteps(ctx, tenant, mission)
	if e != nil {
		return e
	}
	for _, st := range steps {
		if st.Ordinal == ordinal {
			next := domain.StepRunning
			if done {
				next = domain.StepDone
			}
			return s.Repos.UpdateStep(ctx, tenant, st.ID, next)
		}
	}
	return domain.ErrNotFound
}
func (s *Service) ValidateMissionPayload(m domain.Mission, steps []domain.MissionStep) error {
	if m.TenantID == "" || m.RobotID == "" {
		return errors.New("mission identity missing")
	}
	if len(steps) == 0 {
		return errors.New("mission needs steps")
	}
	seen := map[int]bool{}
	for _, st := range steps {
		if seen[st.Ordinal] || strings.TrimSpace(st.Action) == "" {
			return errors.New("invalid mission steps")
		}
		seen[st.Ordinal] = true
	}
	return nil
}
func (s *Service) EstimateMission(ctx context.Context, tenant, id string) (time.Duration, error) {
	v, e := s.MissionView(ctx, tenant, id)
	if e != nil {
		return 0, e
	}
	var d time.Duration
	for _, st := range v.Steps {
		switch st.Action {
		case "navigate":
			d += 2 * time.Minute
		case "pick":
			d += 90 * time.Second
		case "scan":
			d += 30 * time.Second
		default:
			d += time.Minute
		}
	}
	return d, nil
}

type TelemetrySummary struct {
	RobotID     string
	Samples     int
	First, Last time.Time
	Kinds       map[string]int
}

func (s *Service) SummarizeTelemetry(ctx context.Context, tenant, robot string, window time.Duration) (TelemetrySummary, error) {
	page, e := s.ReadTelemetry(ctx, tenant, robot, s.now().Add(time.Second), 500)
	if e != nil {
		return TelemetrySummary{}, e
	}
	sum := TelemetrySummary{RobotID: robot, Kinds: map[string]int{}}
	cut := s.now().Add(-window)
	for _, ev := range page.Items {
		if ev.RecordedAt.Before(cut) {
			continue
		}
		sum.Samples++
		sum.Kinds[ev.Kind]++
		if sum.First.IsZero() || ev.RecordedAt.Before(sum.First) {
			sum.First = ev.RecordedAt
		}
		if ev.RecordedAt.After(sum.Last) {
			sum.Last = ev.RecordedAt
		}
	}
	return sum, nil
}
func (s *Service) DecodeTelemetryPayload(ev domain.TelemetryEvent) (map[string]any, error) {
	var out map[string]any
	if e := json.Unmarshal([]byte(ev.PayloadJSON), &out); e != nil {
		return nil, fmt.Errorf("decode telemetry: %w", e)
	}
	return out, nil
}
func (s *Service) MergeTelemetry(a, b []domain.TelemetryEvent) []domain.TelemetryEvent {
	out := append(a, b...)
	sort.SliceStable(out, func(i, j int) bool {
		if out[i].RecordedAt.Equal(out[j].RecordedAt) {
			return out[i].Sequence < out[j].Sequence
		}
		return out[i].RecordedAt.Before(out[j].RecordedAt)
	})
	return out
}

type FleetSnapshot struct {
	TenantID   string
	Robots     []domain.Robot
	OpenAlerts int
	Quota      int
}

func (s *Service) FleetSnapshot(ctx context.Context, tenant string) (FleetSnapshot, error) {
	robots, e := s.Repos.Robots().List(ctx, tenant, 500)
	if e != nil {
		return FleetSnapshot{}, e
	}
	alerts, e := s.Repos.CountOpenAlerts(ctx, tenant)
	if e != nil {
		return FleetSnapshot{}, e
	}
	quota, e := s.AvailableQuota(ctx, tenant)
	if e != nil {
		return FleetSnapshot{}, e
	}
	return FleetSnapshot{TenantID: tenant, Robots: robots, OpenAlerts: alerts, Quota: quota}, nil
}
func (s *Service) RetireFleet(ctx context.Context, tenant string) error {
	robots, e := s.Repos.Robots().List(ctx, tenant, 500)
	if e != nil {
		return e
	}
	for _, r := range robots {
		if e = s.RetireRobot(ctx, tenant, r.ID); e != nil {
			return e
		}
	}
	return nil
}

type PolicyDecision struct {
	Allowed  bool
	Reason   string
	Revision int
}

func (s *Service) Decide(ctx context.Context, p domain.Policy, m domain.Mission) (PolicyDecision, error) {
	ok, e := s.EvaluatePolicy(ctx, p, m)
	if e != nil {
		return PolicyDecision{}, e
	}
	d := PolicyDecision{Allowed: ok, Revision: p.Revision}
	if !ok {
		d.Reason = "policy denied"
	}
	return d, nil
}
func (s *Service) RequireApproval(ctx context.Context, tenant, id string, policy domain.Policy) (domain.Mission, error) {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return m, e
	}
	d, e := s.Decide(ctx, policy, m)
	if e != nil {
		return m, e
	}
	if !d.Allowed {
		return m, domain.ErrForbidden
	}
	return s.ApproveMission(ctx, tenant, id, "policy", policy)
}

func (s *Service) Search(ctx context.Context, tenant string, f repository.MissionFilter) (domain.Page[domain.Mission], error) {
	if e := requireTenant(tenant); e != nil {
		return domain.Page[domain.Mission]{}, e
	}
	return s.Repos.SearchMissions(ctx, tenant, f)
}
func (s *Service) ValidateTenantAccess(session domain.Session, tenant string) error {
	if session.TenantID != tenant {
		return domain.ErrForbidden
	}
	return session.Valid(s.now())
}
func (s *Service) CancelIfExpired(ctx context.Context, tenant, id string, now time.Time) error {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return e
	}
	if m.UpdatedAt.Add(24 * time.Hour).After(now) {
		return nil
	}
	return s.CancelMission(ctx, tenant, id)
}
