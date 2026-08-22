package service

import (
	"context"
	"testing"
)

func TestQuotaErrorsReachMissionCaller(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := s.ReserveMissionQuota(ctx, "tenant"); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.ReserveMissionQuota(ctx, "tenant"); err == nil {
		t.Fatal("quota exhaustion was reported as success")
	}
}
