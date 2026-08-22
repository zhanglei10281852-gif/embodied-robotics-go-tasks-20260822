package audit

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sort"
	"time"
)

type Chain struct {
	Previous string
	Events   []domain.AuditEvent
}

func (c *Chain) Append(e domain.AuditEvent) error {
	if e.ID == "" || e.TenantID == "" || e.OccurredAt.IsZero() {
		return errors.New("invalid audit event")
	}
	if len(c.Events) > 0 && e.OccurredAt.Before(c.Events[len(c.Events)-1].OccurredAt) {
		return errors.New("audit order regression")
	}
	c.Events = append(c.Events, e)
	return nil
}
func (c Chain) Digest() string {
	h := sha256.New()
	for _, e := range c.Events {
		b, _ := json.Marshal(e)
		h.Write([]byte(c.Previous))
		h.Write(b)
		c.Previous = hex.EncodeToString(h.Sum(nil))
	}
	return c.Previous
}
func (c Chain) Sorted() Chain {
	out := Chain{Previous: c.Previous, Events: append([]domain.AuditEvent(nil), c.Events...)}
	sort.SliceStable(out.Events, func(i, j int) bool { return out.Events[i].OccurredAt.Before(out.Events[j].OccurredAt) })
	return out
}
func (c Chain) Since(at time.Time) []domain.AuditEvent {
	out := []domain.AuditEvent{}
	for _, e := range c.Events {
		if !e.OccurredAt.Before(at) {
			out = append(out, e)
		}
	}
	return out
}
