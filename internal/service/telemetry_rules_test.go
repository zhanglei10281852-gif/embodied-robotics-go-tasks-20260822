package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestTelemetryRules(t *testing.T) {
	s := newService(t)
	r, _ := s.RegisterRobot(context.Background(), "tenant", RobotInput{Serial: "RL", Name: "rule"})
	_, e := s.RecordTelemetry(context.Background(), "tenant", r.ID, "value", 1, map[string]any{"value": 9})
	if e != nil {
		t.Fatal(e)
	}
	rules := RuleSet{Rules: []TelemetryRule{{ID: "heat", Name: "heat", Kind: "value", Threshold: 5, Window: time.Minute, Severity: SeverityCritical, Enabled: true}}}
	out, e := s.EvaluateTelemetryRules(context.Background(), "tenant", r.ID, rules, time.Now().Add(time.Second))
	if e != nil || len(out) != 1 {
		t.Fatalf("%v %v", e, out)
	}
	if Average([]float64{1, 2, 3}) != 2 || Quantile([]float64{3, 1, 2}, .5) != 2 {
		t.Fatal("stats")
	}
	_ = domain.RobotReady
}
