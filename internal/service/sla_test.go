package service

import (
	"context"
	"testing"
	"time"
)

func TestSLAWindows(t *testing.T) {
	w := SLAWindow{Start: time.Now(), End: time.Now().Add(time.Hour), Target: time.Minute}
	if e := w.Validate(); e != nil {
		t.Fatal(e)
	}
	if !w.Contains(w.Start.Add(time.Second)) {
		t.Fatal("contains")
	}
	s := newService(t)
	_, e := s.EvaluateMissionSLA(context.Background(), "tenant", "missing", time.Minute)
	if e == nil {
		t.Fatal("missing mission")
	}
}
