package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"sort"
	"time"
)

type Snapshot struct {
	TenantID    string
	GeneratedAt time.Time
	Robots      []domain.Robot
	Missions    []domain.Mission
	Alerts      int
	Digest      string
}

func (s *Service) BuildSnapshot(ctx context.Context, tenant string) (Snapshot, error) {
	robots, e := s.Repos.Robots().List(ctx, tenant, 1000)
	if e != nil {
		return Snapshot{}, e
	}
	page, e := s.Search(ctx, tenant, repository.MissionFilter{Limit: 1000})
	if e != nil {
		return Snapshot{}, e
	}
	alerts, e := s.Repos.CountOpenAlerts(ctx, tenant)
	if e != nil {
		return Snapshot{}, e
	}
	sort.Slice(robots, func(i, j int) bool { return robots[i].ID < robots[j].ID })
	raw, _ := json.Marshal(struct {
		R []domain.Robot
		M []domain.Mission
		A int
	}{robots, page.Items, alerts})
	return Snapshot{TenantID: tenant, GeneratedAt: s.now(), Robots: robots, Missions: page.Items, Alerts: alerts, Digest: domain.StableDigest(string(raw))}, nil
}
func (s Snapshot) Validate() error {
	if s.TenantID == "" || s.GeneratedAt.IsZero() {
		return fmt.Errorf("snapshot metadata missing")
	}
	if s.Digest == "" {
		return fmt.Errorf("snapshot digest missing")
	}
	return nil
}
func (s Snapshot) RobotStatusCounts() map[string]int {
	out := map[string]int{}
	for _, r := range s.Robots {
		out[r.Status]++
	}
	return out
}
