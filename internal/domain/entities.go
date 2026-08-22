package domain

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

type Tenant struct {
	ID, Name  string
	Quota     int
	CreatedAt time.Time
}
type User struct {
	ID, TenantID, Email, Role string
	RevokedAt                 *time.Time
}
type Session struct {
	Token, UserID, TenantID string
	ExpiresAt               time.Time
	RevokedAt               *time.Time
}
type Robot struct {
	ID, TenantID, Serial, Name, Status string
	Version                            int64
	LeaseUntil                         *time.Time
	CreatedAt                          time.Time
}
type Capability struct {
	RobotID, Name, SchemaJSON string
	Revision                  int64
	UpdatedAt                 time.Time
}
type Mission struct {
	ID, TenantID, RobotID, Status, Owner string
	Priority, PolicyVersion, Attempt     int
	Version                              int64
	CreatedAt, UpdatedAt                 time.Time
}
type MissionStep struct {
	ID, MissionID, Action, PayloadJSON, Status string
	Ordinal                                    int
}
type TelemetryEvent struct {
	ID, TenantID, RobotID, Kind, PayloadJSON string
	Sequence                                 int64
	RecordedAt                               time.Time
}
type Policy struct {
	ID, TenantID, Name, Expression, State string
	Revision                              int
	UpdatedAt                             time.Time
}
type Approval struct {
	ID, TenantID, MissionID, PolicyID, Decision, DecidedBy string
	PolicyRevision                                         int
	DecidedAt                                              *time.Time
}
type Handoff struct {
	ID, TenantID, MissionID, OperatorID, State string
	LeaseUntil                                 time.Time
}
type Alert struct {
	ID, TenantID, RobotID, DedupeKey, Severity, Message, State string
	CreatedAt                                                  time.Time
	AcknowledgedAt                                             *time.Time
}
type AuditEvent struct {
	ID, TenantID, ActorID, Action, ObjectType, ObjectID, DetailJSON string
	OccurredAt                                                      time.Time
}
type OutboxJob struct {
	ID, TenantID, Topic, PayloadJSON, State, LastError string
	Attempts                                           int
	AvailableAt                                        time.Time
}
type Page[T any] struct {
	Items      []T
	NextCursor string
	Total      int
}

var (
	ErrInvalidState = errors.New("invalid state transition")
	ErrConflict     = errors.New("resource version conflict")
	ErrNotFound     = errors.New("resource not found")
	ErrForbidden    = errors.New("operation forbidden")
	ErrQuota        = errors.New("tenant quota exceeded")
	ErrExpired      = errors.New("lease expired")
	ErrRevoked      = errors.New("session revoked")
)

const (
	RobotOffline      = "offline"
	RobotReady        = "ready"
	RobotBusy         = "busy"
	RobotRetired      = "retired"
	MissionDraft      = "draft"
	MissionApproved   = "approved"
	MissionQueued     = "queued"
	MissionRunning    = "running"
	MissionSucceeded  = "succeeded"
	MissionFailed     = "failed"
	MissionCancelled  = "cancelled"
	StepPending       = "pending"
	StepRunning       = "running"
	StepDone          = "done"
	StepFailed        = "failed"
	AlertOpen         = "open"
	AlertAcknowledged = "acknowledged"
	AlertClosed       = "closed"
	OutboxPending     = "pending"
	OutboxSending     = "sending"
	OutboxDone        = "done"
	OutboxDead        = "dead"
)

func (r Robot) Validate() error {
	if strings.TrimSpace(r.ID) == "" || strings.TrimSpace(r.TenantID) == "" {
		return errors.New("robot identity is required")
	}
	if strings.TrimSpace(r.Serial) == "" || strings.TrimSpace(r.Name) == "" {
		return errors.New("robot serial and name are required")
	}
	if r.Status == "" {
		return errors.New("robot status is required")
	}
	if r.Version < 1 {
		return errors.New("robot version must be positive")
	}
	return nil
}

func (r Robot) CanAcceptMission() bool      { return r.Status == RobotReady && r.LeaseUntil == nil }
func (r Robot) IsLeased(now time.Time) bool { return r.LeaseUntil != nil && r.LeaseUntil.After(now) }
func (r Robot) String() string              { return fmt.Sprintf("%s/%s", r.TenantID, r.Serial) }

func (m Mission) CanTransition(next string) bool {
	switch m.Status {
	case MissionDraft:
		return next == MissionApproved || next == MissionCancelled
	case MissionApproved:
		return next == MissionQueued || next == MissionCancelled
	case MissionQueued:
		return next == MissionRunning || next == MissionCancelled || next == MissionFailed
	case MissionRunning:
		return next == MissionSucceeded || next == MissionFailed || next == MissionCancelled
	case MissionFailed:
		return next == MissionQueued || next == MissionCancelled
	default:
		return false
	}
}

func (m Mission) Transition(next string) (Mission, error) {
	if !m.CanTransition(next) {
		return Mission{}, fmt.Errorf("%w: %s -> %s", ErrInvalidState, m.Status, next)
	}
	m.Status = next
	m.Version++
	m.UpdatedAt = time.Now().UTC()
	return m, nil
}

func (m Mission) IsTerminal() bool {
	return m.Status == MissionSucceeded || m.Status == MissionCancelled
}
func (m Mission) Retryable() bool { return m.Status == MissionFailed && m.Attempt <= 5 }
func (m Mission) Cursor() string {
	return fmt.Sprintf("%s:%d", m.UpdatedAt.UTC().Format(time.RFC3339Nano), m.Version)
}

func (p Policy) Validate() error {
	if p.TenantID == "" || p.Name == "" || p.Expression == "" {
		return errors.New("policy fields are required")
	}
	if p.Revision < 1 {
		return errors.New("policy revision must be positive")
	}
	return nil
}

func (a Alert) CanAcknowledge() bool   { return a.State == AlertOpen }
func (a Alert) CanClose() bool         { return a.State == AlertAcknowledged }
func (a Alert) DedupeIdentity() string { return a.TenantID + ":" + a.RobotID + ":" + a.DedupeKey }

func (s Session) Valid(now time.Time) error {
	if s.Token == "" || s.UserID == "" || s.TenantID == "" {
		return ErrRevoked
	}
	if s.RevokedAt != nil {
		return ErrRevoked
	}
	if !s.ExpiresAt.After(now) {
		return ErrExpired
	}
	return nil
}

func (e TelemetryEvent) Validate() error {
	if e.ID == "" || e.TenantID == "" || e.RobotID == "" || e.Kind == "" {
		return errors.New("telemetry identity is required")
	}
	if e.Sequence < 1 {
		return errors.New("telemetry sequence must be positive")
	}
	return nil
}

func (j OutboxJob) Retryable(now time.Time) bool {
	return j.State == OutboxPending && !j.AvailableAt.After(now) && j.Attempts < 8
}
func (j OutboxJob) MarkFailure(err error) OutboxJob {
	j.Attempts++
	j.LastError = err.Error()
	if j.Attempts >= 8 {
		j.State = OutboxDead
	}
	return j
}

func ValidateRole(role string) error {
	if role != "operator" && role != "supervisor" && role != "auditor" {
		return ErrForbidden
	}
	return nil
}

func ValidatePriority(p int) error {
	if p < 0 || p > 100 {
		return errors.New("priority outside range")
	}
	return nil
}

func CloneStrings(in []string) []string { out := make([]string, len(in)); copy(out, in); return out }
func CloneEvents(in []TelemetryEvent) []TelemetryEvent {
	out := make([]TelemetryEvent, len(in))
	copy(out, in)
	return out
}
