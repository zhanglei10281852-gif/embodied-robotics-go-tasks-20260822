package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestPolicyRevisionConflictIsPreserved(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	p := domain.Policy{ID: "p-old", TenantID: "tenant", Name: "safe", Expression: "allow", Revision: 1, State: "approved"}
	m := domain.Mission{ID: "m-new", TenantID: "tenant", RobotID: "r", PolicyVersion: 2}
	if allowed, err := s.EvaluatePolicy(ctx, p, m); err == nil || allowed {
		t.Fatalf("stale policy accepted: allowed=%v err=%v", allowed, err)
	}
}
