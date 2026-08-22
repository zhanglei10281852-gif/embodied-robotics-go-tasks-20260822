package service

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestPolicyEvaluationStopsOnCancellation(t *testing.T) {
	s := newService(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := domain.Policy{ID: "p", TenantID: "tenant", Name: "safe", Expression: "allow", Revision: 1, State: "approved"}
	m := domain.Mission{TenantID: "tenant", PolicyVersion: 1}
	if _, err := s.EvaluatePolicy(ctx, p, m); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation ignored: %v", err)
	}
}
