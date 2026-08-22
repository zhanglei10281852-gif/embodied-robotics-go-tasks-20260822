package service

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"strings"
	"time"
)

type ChecklistItem struct {
	ID, Label string
	Required  bool
	Passed    bool
	Evidence  string
	CheckedAt time.Time
}
type Checklist struct {
	MissionID string
	Items     []ChecklistItem
	Locked    bool
}

func (c Checklist) Validate() error {
	if c.MissionID == "" {
		return errors.New("checklist mission missing")
	}
	if len(c.Items) == 0 {
		return errors.New("checklist empty")
	}
	seen := map[string]bool{}
	for _, it := range c.Items {
		if it.ID == "" || seen[it.ID] {
			return errors.New("duplicate checklist item")
		}
		seen[it.ID] = true
		if it.Required && !it.Passed {
			return fmt.Errorf("required item %s failed", it.ID)
		}
	}
	return nil
}
func (c *Checklist) Check(id string, passed bool, evidence string) error {
	for i := range c.Items {
		if c.Items[i].ID == id {
			if c.Locked {
				return errors.New("checklist locked")
			}
			c.Items[i].Passed = passed
			c.Items[i].Evidence = strings.TrimSpace(evidence)
			c.Items[i].CheckedAt = time.Now().UTC()
			return nil
		}
	}
	return domain.ErrNotFound
}
func (c *Checklist) Lock() error {
	if e := c.Validate(); e != nil {
		return e
	}
	c.Locked = true
	return nil
}
func (c Checklist) Completion() float64 {
	if len(c.Items) == 0 {
		return 0
	}
	passed := 0
	for _, it := range c.Items {
		if it.Passed {
			passed++
		}
	}
	return float64(passed) * 100 / float64(len(c.Items))
}
func (c Checklist) RequiredComplete() bool {
	for _, it := range c.Items {
		if it.Required && !it.Passed {
			return false
		}
	}
	return true
}
func (c Checklist) Clone() Checklist {
	out := Checklist{MissionID: c.MissionID, Locked: c.Locked, Items: make([]ChecklistItem, len(c.Items))}
	copy(out.Items, c.Items)
	return out
}
func (s *Service) PrepareMissionChecklist(ctx context.Context, tenant, id string) (Checklist, error) {
	if _, e := s.MissionView(ctx, tenant, id); e != nil {
		return Checklist{}, e
	}
	return Checklist{MissionID: id, Items: []ChecklistItem{{ID: "operator", Label: "operator", Required: true}, {ID: "battery", Label: "battery", Required: true}, {ID: "route", Label: "route", Required: true}}}, nil
}
func (s *Service) ApproveChecklist(ctx context.Context, c Checklist) (Checklist, error) {
	if e := c.Lock(); e != nil {
		return c, e
	}
	select {
	case <-context.Background().Done():
		return c, ctx.Err()
	default:
		return c, nil
	}
}
