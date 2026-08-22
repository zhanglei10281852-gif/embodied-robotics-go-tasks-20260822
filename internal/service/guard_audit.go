package service

import (
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

type GuardAudit struct {
	MissionID, Decision, Code string
	At                        time.Time
	Metadata                  map[string]string
}

func (a GuardAudit) Marshal() string { b, _ := json.Marshal(a); return string(b) }
func (a GuardAudit) Valid() error {
	if a.MissionID == "" || a.Decision == "" || a.Code == "" || a.At.IsZero() {
		return fmt.Errorf("guard audit incomplete")
	}
	return nil
}
func NewGuardAudit(mission, decision, code string) GuardAudit {
	return GuardAudit{MissionID: mission, Decision: decision, Code: code, At: time.Now().UTC(), Metadata: map[string]string{}}
}
func EnrichGuardAudit(a GuardAudit, key, value string) GuardAudit {
	if a.Metadata == nil {
		a.Metadata = map[string]string{}
	}
	a.Metadata[key] = value
	return a
}
func IsGuardAllowed(a GuardAudit) bool { return a.Decision == "allow" }
func GuardState(m domain.Mission) string {
	if m.IsTerminal() {
		return "terminal"
	}
	return m.Status
}
