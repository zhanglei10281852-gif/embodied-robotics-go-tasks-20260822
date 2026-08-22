package domain

import (
	"testing"
	"time"
)

func TestMetricSeries(t *testing.T) {
	now := time.Now()
	s := MetricSeries{Name: "latency", Samples: []MetricSample{{Value: 3, At: now.Add(2 * time.Second)}, {Value: 1, At: now}, {Value: 2, At: now.Add(time.Second)}}}
	if s.Average() != 2 || s.P95() != 3 {
		t.Fatalf("avg %.1f p95 %.1f", s.Average(), s.P95())
	}
	w := s.Window(now, now.Add(2*time.Second))
	if len(w.Samples) != 2 {
		t.Fatal(len(w.Samples))
	}
}
