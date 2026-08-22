package service

import (
	"context"
	"testing"
	"time"
)

func TestAnalyticsEmptyAndPercent(t *testing.T) {
	s := newService(t)
	if got, e := s.Availability(context.Background(), "tenant"); e != nil || got.Ready != 0 {
		t.Fatalf("%v %+v", e, got)
	}
	if Percent(1, 4) != 25 {
		t.Fatal(Percent(1, 4))
	}
	if Percent(1, 0) != 0 {
		t.Fatal("zero")
	}
	_, e := s.MotionStatistics(context.Background(), "tenant", "missing", time.Now())
	if e == nil {
		t.Fatal("missing robot should error")
	}
}
