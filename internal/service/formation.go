package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
)

func (s *Service) CreateFormation(ctx context.Context, tenant, name string, robots []string) (repository.Formation, error) {
	if len(robots) < 2 {
		return repository.Formation{}, fmt.Errorf("formation needs two robots")
	}
	robots = append([]string(nil), robots...)
	if e := s.Repos.EnsureFormationTables(ctx); e != nil {
		return repository.Formation{}, e
	}
	f := repository.Formation{ID: id("formation"), TenantID: tenant, Name: name, State: "draft"}
	if e := s.Repos.CreateFormation(ctx, f, robots); e != nil {
		return repository.Formation{}, e
	}
	return f, nil
}
