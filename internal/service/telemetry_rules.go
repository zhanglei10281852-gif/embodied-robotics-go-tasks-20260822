package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"math"
	"strings"
	"time"
)

type RuleSeverity string

const (
	SeverityInfo     RuleSeverity = "info"
	SeverityWarn     RuleSeverity = "warn"
	SeverityCritical RuleSeverity = "critical"
)

type TelemetryRule struct {
	ID, Name, Kind string
	Threshold      float64
	Window         time.Duration
	Severity       RuleSeverity
	Enabled        bool
}

func (r TelemetryRule) Validate() error {
	if r.ID == "" || r.Name == "" || r.Kind == "" {
		return errors.New("rule identity missing")
	}
	if r.Window <= 0 {
		return errors.New("rule window positive")
	}
	if r.Severity != SeverityInfo && r.Severity != SeverityWarn && r.Severity != SeverityCritical {
		return errors.New("rule severity invalid")
	}
	return nil
}
func (r TelemetryRule) Match(ev domain.TelemetryEvent) bool {
	if !r.Enabled || ev.Kind != r.Kind {
		return false
	}
	v := 0.0
	fmt.Sscanf(ev.PayloadJSON, `{"value":%f}`, &v)
	return math.Abs(v) >= r.Threshold
}

type RuleResult struct {
	RuleID, RobotID string
	Matched         bool
	Value           float64
	At              time.Time
	Reason          string
}

func (r TelemetryRule) Evaluate(ev domain.TelemetryEvent) RuleResult {
	out := RuleResult{RuleID: r.ID, RobotID: ev.RobotID, At: ev.RecordedAt}
	if !r.Match(ev) {
		out.Reason = "below threshold"
		return out
	}
	out.Matched = true
	out.Reason = "threshold exceeded"
	return out
}

type RuleSet struct{ Rules []TelemetryRule }

func (s RuleSet) Validate() error {
	if len(s.Rules) == 0 {
		return errors.New("empty rule set")
	}
	seen := map[string]bool{}
	for _, r := range s.Rules {
		if e := r.Validate(); e != nil {
			return e
		}
		if seen[r.ID] {
			return errors.New("duplicate rule")
		}
		seen[r.ID] = true
	}
	return nil
}
func (s RuleSet) Evaluate(events []domain.TelemetryEvent) []RuleResult {
	out := []RuleResult{}
	for _, ev := range events {
		for _, r := range s.Rules {
			if result := r.Evaluate(ev); result.Matched {
				out = append(out, result)
			}
		}
	}
	return out
}
func (s *Service) EvaluateTelemetryRules(ctx context.Context, tenant, robot string, rules RuleSet, before time.Time) ([]RuleResult, error) {
	if e := rules.Validate(); e != nil {
		return nil, e
	}
	page, e := s.ReadTelemetry(ctx, tenant, robot, before, 500)
	if e != nil {
		return nil, e
	}
	return rules.Evaluate(page.Items), nil
}

type SafetyWindow struct {
	Start, End time.Time
	Reason     string
}

func (w SafetyWindow) Allows(t time.Time) bool { return !t.Before(w.Start) && t.Before(w.End) }
func (w SafetyWindow) Validate() error {
	if w.Start.IsZero() || w.End.IsZero() || !w.End.After(w.Start) {
		return errors.New("invalid safety window")
	}
	if strings.TrimSpace(w.Reason) == "" {
		return errors.New("safety reason required")
	}
	return nil
}

type LimitSet struct{ Speed, Acceleration, Payload float64 }

func (l LimitSet) Validate() error {
	if l.Speed <= 0 || l.Acceleration <= 0 || l.Payload <= 0 {
		return errors.New("limits positive")
	}
	return nil
}
func (l LimitSet) ClampSpeed(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > l.Speed {
		return l.Speed
	}
	return v
}
func (l LimitSet) ClampAcceleration(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > l.Acceleration {
		return l.Acceleration
	}
	return v
}
func (l LimitSet) AcceptPayload(v float64) bool { return v >= 0 && v <= l.Payload }

type Calibration struct {
	RobotID                   string
	OffsetX, OffsetY, OffsetZ float64
	AppliedAt                 time.Time
}

func (c Calibration) Apply(p [3]float64) [3]float64 {
	return [3]float64{p[0] + c.OffsetX, p[1] + c.OffsetY, p[2] + c.OffsetZ}
}
func (c Calibration) Validate() error {
	if c.RobotID == "" || c.AppliedAt.IsZero() {
		return errors.New("calibration incomplete")
	}
	return nil
}
func Average(values []float64) float64 {
	if len(values) == 0 {
		return 0
	}
	var total float64
	for _, v := range values {
		total += v
	}
	return total / float64(len(values))
}
func Quantile(values []float64, p float64) float64 {
	if len(values) == 0 {
		return 0
	}
	if p < 0 {
		p = 0
	}
	if p > 1 {
		p = 1
	}
	copyValues := append([]float64(nil), values...)
	for i := 1; i < len(copyValues); i++ {
		for j := i; j > 0 && copyValues[j] < copyValues[j-1]; j-- {
			copyValues[j], copyValues[j-1] = copyValues[j-1], copyValues[j]
		}
	}
	idx := int(float64(len(copyValues)-1) * p)
	return copyValues[idx]
}
