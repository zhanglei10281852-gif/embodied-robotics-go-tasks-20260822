package service

import (
	"context"
	"testing"
)

func TestQuotaNeverGoesNegative(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	for i := 0; i < 3; i++ {
		if err := s.ReserveMissionQuota(ctx, "tenant"); err != nil {
			t.Fatal(err)
		}
	}
	if err := s.ReserveMissionQuota(ctx, "tenant"); err == nil {
		t.Fatal("quota reservation exceeded limit")
	}
	if quota, err := s.AvailableQuota(ctx, "tenant"); err != nil || quota < 0 {
		t.Fatalf("quota=%d err=%v", quota, err)
	}
}
