package service

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestHandoffRenewalHonorsCancellation(t *testing.T) {
	s := newService(t)
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	h := domain.Handoff{TenantID: "tenant", MissionID: "m", OperatorID: "u", LeaseUntil: time.Now().Add(time.Minute), State: "active"}
	if _, err := s.RenewHandoff(ctx, h, time.Minute); !errors.Is(err, context.Canceled) {
		t.Fatalf("cancellation lost: %v", err)
	}
}
