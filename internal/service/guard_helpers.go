package service

import (
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func GuardError(code string) error {
	if code == "" {
		return errors.New("guard code missing")
	}
	return fmt.Errorf("guard denied: %s", code)
}
func MissionAge(m domain.Mission, now time.Time) time.Duration {
	if now.Before(m.CreatedAt) {
		return 0
	}
	return now.Sub(m.CreatedAt)
}
func MissionFresh(m domain.Mission, now time.Time, max time.Duration) bool {
	return MissionAge(m, now) <= max
}
func RobotReadyFor(m domain.Mission, r domain.Robot, now time.Time) bool {
	if m.TenantID != r.TenantID || m.RobotID != r.ID {
		return false
	}
	if r.Status != domain.RobotReady {
		return false
	}
	return !r.IsLeased(now)
}
func SameTenant(tenant string, values ...string) bool {
	for _, v := range values {
		if v != tenant {
			return false
		}
	}
	return tenant != ""
}
func NonTerminal(status string) bool {
	return status != domain.MissionSucceeded && status != domain.MissionCancelled
}
func Terminal(status string) bool             { return !NonTerminal(status) }
func AllowedTransition(from, to string) bool  { return (domain.Mission{Status: from}).CanTransition(to) }
func RetryCountAllowed(attempt, max int) bool { return attempt >= 0 && attempt < max }
func DeadlineRemaining(deadline, now time.Time) time.Duration {
	if deadline.Before(now) {
		return 0
	}
	return deadline.Sub(now)
}
