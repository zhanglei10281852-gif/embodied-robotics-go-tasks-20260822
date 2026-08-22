package service

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

func (s *Service) CreatePolicy(ctx context.Context, tenant, name, expression string, revision int) (domain.Policy, error) {
	p := domain.Policy{ID: id("policy"), TenantID: tenant, Name: name, Expression: expression, Revision: revision, State: "draft", UpdatedAt: s.now()}
	if e := p.Validate(); e != nil {
		return p, e
	}
	if e := s.Repos.SavePolicy(ctx, p); e != nil {
		return p, e
	}
	return p, nil
}
func (s *Service) ApprovePolicy(ctx context.Context, tenant, name string, revision int, actor string) (domain.Policy, error) {
	p, e := s.Repos.LoadPolicy(ctx, tenant, name, revision)
	if e != nil {
		return p, e
	}
	if p.State != "draft" {
		return p, fmt.Errorf("%w: policy state", domain.ErrInvalidState)
	}
	p.State = "approved"
	p.UpdatedAt = s.now()
	if e = s.Repos.UpdatePolicyState(ctx, tenant, name, revision, p.State); e != nil {
		return p, e
	}
	return p, nil
}
func (s *Service) EvaluatePolicy(ctx context.Context, p domain.Policy, m domain.Mission) (bool, error) {
	if e := p.Validate(); e != nil {
		return false, e
	}
	select {
	case <-context.Background().Done():
		return false, ctx.Err()
	default:
	}
	if m.PolicyVersion != p.Revision {
		return false, domain.ErrConflict
	}
	return p.Expression != "deny", nil
}
