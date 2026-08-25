package domain

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strings"
	"time"
)

var serialPattern = regexp.MustCompile(`^[A-Z0-9][A-Z0-9-]{2,63}$`)

type RobotSerial string

func NewRobotSerial(v string) (RobotSerial, error) {
	v = strings.ToUpper(strings.TrimSpace(v))
	if !serialPattern.MatchString(v) {
		return "", errors.New("invalid robot serial")
	}
	return RobotSerial(v), nil
}
func (s RobotSerial) String() string { return string(s) }

type TenantID string

func (t TenantID) Valid() bool { return strings.TrimSpace(string(t)) != "" }

type MissionID string

func (m MissionID) String() string { return string(m) }

type Cursor struct {
	At time.Time
	ID string
}

func (c Cursor) Encode() string {
	return fmt.Sprintf("%s|%s", c.At.UTC().Format(time.RFC3339Nano), c.ID)
}
func ParseCursor(v string) (Cursor, error) {
	parts := strings.SplitN(v, "|", 2)
	if len(parts) != 2 {
		return Cursor{}, errors.New("invalid cursor")
	}
	at, e := time.Parse(time.RFC3339Nano, parts[0])
	if e != nil || parts[1] == "" {
		return Cursor{}, errors.New("invalid cursor")
	}
	return Cursor{At: at, ID: parts[1]}, nil
}

type CapabilitySet map[string]int64

func (c CapabilitySet) Clone() CapabilitySet {
	out := CapabilitySet{}
	for k, v := range c {
		out[k] = v
	}
	return out
}
func (c CapabilitySet) Supports(name string, revision int64) bool { return c[name] >= revision }
func StableDigest(parts ...string) string {
	h := sha256.New()
	for _, p := range parts {
		h.Write([]byte(p))
		h.Write([]byte{0})
	}
	return hex.EncodeToString(h.Sum(nil))
}

type StepPlan struct {
	Ordinal  int
	Action   string
	Timeout  time.Duration
	Required []string
}

func (p StepPlan) Validate() error {
	if p.Ordinal < 0 || strings.TrimSpace(p.Action) == "" {
		return errors.New("invalid step plan")
	}
	if p.Timeout <= 0 {
		return errors.New("step timeout must be positive")
	}
	if len(p.Required) == 0 {
		return errors.New("step capabilities required")
	}
	return nil
}
func SortSteps(steps []MissionStep) []MissionStep {
	out := append([]MissionStep(nil), steps...)
	sort.SliceStable(out, func(i, j int) bool { return out[i].Ordinal < out[j].Ordinal })
	return out
}
func NormalizeLabels(in map[string]string) map[string]string {
	out := map[string]string{}
	for k, v := range in {
		key := strings.ToLower(strings.TrimSpace(k))
		val := strings.TrimSpace(v)
		if key != "" && val != "" {
			out[key] = val
		}
	}
	return out
}
