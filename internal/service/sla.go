package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sort"
	"time"
)

type SLAWindow struct {
	Start, End time.Time
	Target     time.Duration
}

func (w SLAWindow) Validate() error {
	if w.Start.IsZero() || w.End.IsZero() || !w.End.After(w.Start) {
		return fmt.Errorf("invalid SLA window")
	}
	if w.Target <= 0 {
		return fmt.Errorf("SLA target positive")
	}
	return nil
}
func (w SLAWindow) Contains(t time.Time) bool { return !t.Before(w.Start) && t.Before(w.End) }

type MissionSLA struct {
	MissionID                  string
	Created, Started, Finished time.Time
	Target                     time.Duration
	Breached                   bool
}

func (s MissionSLA) Duration() time.Duration {
	end := s.Finished
	if end.IsZero() {
		end = time.Now().UTC()
	}
	return end.Sub(s.Created)
}
func (s MissionSLA) Evaluate() MissionSLA { s.Breached = s.Duration() > s.Target; return s }
func (s *Service) EvaluateMissionSLA(ctx context.Context, tenant, id string, target time.Duration) (MissionSLA, error) {
	m, e := s.Repos.Missions().Get(ctx, tenant, id)
	if e != nil {
		return MissionSLA{}, e
	}
	sla := MissionSLA{MissionID: id, Created: m.CreatedAt, Target: target}
	if m.Status == domain.MissionRunning {
		sla.Started = m.UpdatedAt
	}
	if m.IsTerminal() {
		sla.Finished = m.UpdatedAt
	}
	return sla.Evaluate(), nil
}
func SortSLA(in []MissionSLA) []MissionSLA {
	out := append([]MissionSLA(nil), in...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Duration() > out[j].Duration() })
	return out
}
func PercentileSLA(in []MissionSLA, p float64) time.Duration {
	if len(in) == 0 {
		return 0
	}
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	sorted := SortSLA(in)
	idx := int(float64(len(sorted)-1) * p)
	return sorted[idx].Duration()
}
