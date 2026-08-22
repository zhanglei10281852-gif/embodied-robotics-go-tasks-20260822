package service

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

func (s *Service) PublishCapability(ctx context.Context, tenant, robot, name string, revision int, schema map[string]any) (domain.Capability, error) {
	if revision < 1 {
		return domain.Capability{}, fmt.Errorf("revision must be positive")
	}
	r, e := s.Repos.Robots().Get(ctx, tenant, robot)
	if e != nil {
		return domain.Capability{}, e
	}
	if r.Status == domain.RobotRetired {
		return domain.Capability{}, domain.ErrForbidden
	}
	b, e := json.Marshal(schema)
	if e != nil {
		b = []byte("{}")
		e = nil
	}
	c := domain.Capability{RobotID: r.ID, Name: name, Revision: int64(revision), SchemaJSON: string(b), UpdatedAt: s.now()}
	if e = s.Repos.Capabilities().Upsert(ctx, c); e != nil {
		return domain.Capability{}, e
	}
	return c, nil
}
func (s *Service) Capability(ctx context.Context, tenant, robot, name string) (domain.Capability, error) {
	r, e := s.Repos.Robots().Get(ctx, tenant, robot)
	if e != nil {
		return domain.Capability{}, e
	}
	return s.Repos.Capabilities().Get(ctx, r.ID, name)
}
