package domain

import (
	"errors"
	"strings"
	"time"
)

type CommandRequest struct {
	TenantID, RobotID, Action, IdempotencyKey string
	Payload                                   []byte
	Deadline                                  time.Time
}

func (c CommandRequest) Validate() error {
	if strings.TrimSpace(c.TenantID) == "" || strings.TrimSpace(c.RobotID) == "" {
		return errors.New("command scope missing")
	}
	if strings.TrimSpace(c.Action) == "" {
		return errors.New("command action missing")
	}
	if len(c.Payload) > 1<<20 {
		return errors.New("command payload too large")
	}
	if c.Deadline.IsZero() || !c.Deadline.After(time.Now()) {
		return ErrExpired
	}
	return nil
}

type HandoffRequest struct {
	TenantID, MissionID, OperatorID string
	TTL                             time.Duration
}

func (h HandoffRequest) Validate() error {
	if h.TenantID == "" || h.MissionID == "" || h.OperatorID == "" {
		return errors.New("handoff scope missing")
	}
	if h.TTL <= 0 {
		return ErrExpired
	}
	return nil
}

type BatchResult struct {
	Accepted, Rejected int
	Errors             []error
}

func (b *BatchResult) Add(err error) {
	if err == nil {
		b.Accepted++
	} else {
		b.Rejected++
		b.Errors = append(b.Errors, err)
	}
}
