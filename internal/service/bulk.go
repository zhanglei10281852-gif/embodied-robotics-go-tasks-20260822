package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sync"
)

type BulkRegistration struct{ Serial, Name string }

func (s *Service) RegisterFleet(ctx context.Context, tenant string, inputs []BulkRegistration) ([]domain.Robot, error) {
	if len(inputs) == 0 {
		return nil, fmt.Errorf("fleet empty")
	}
	out := make([]domain.Robot, 0, len(inputs))
	for _, in := range inputs {
		r, e := s.RegisterRobot(ctx, tenant, RobotInput{Serial: in.Serial, Name: in.Name})
		if e != nil {
			return nil, e
		}
		out = append(out, r)
	}
	return out, nil
}
func (s *Service) RegisterFleetConcurrent(ctx context.Context, tenant string, inputs []BulkRegistration) ([]domain.Robot, error) {
	out := make([]domain.Robot, len(inputs))
	errs := make(chan error, len(inputs))
	var wg sync.WaitGroup
	for i, in := range inputs {
		wg.Add(1)
		go func(i int, in BulkRegistration) {
			defer wg.Done()
			r, e := s.RegisterRobot(ctx, tenant, RobotInput{Serial: in.Serial, Name: in.Name})
			if e != nil {
				errs <- e
				return
			}
			out[i] = r
		}(i, in)
	}
	wg.Wait()
	close(errs)
	for e := range errs {
		if e != nil {
			return nil, e
		}
	}
	return out, nil
}
func (s *Service) CancelFleet(ctx context.Context, tenant string, ids []string) domain.BatchResult {
	out := domain.BatchResult{}
	for _, id := range ids {
		out.Add(s.CancelMission(ctx, tenant, id))
	}
	return out
}
func (s *Service) ValidateFleet(ctx context.Context, tenant string, robots []BulkRegistration) error {
	seen := map[string]bool{}
	for _, r := range robots {
		if r.Serial == "" || seen[r.Serial] {
			return fmt.Errorf("duplicate fleet serial")
		}
		seen[r.Serial] = true
	}
	_, e := s.Repos.Robots().List(ctx, tenant, 1)
	return e
}
