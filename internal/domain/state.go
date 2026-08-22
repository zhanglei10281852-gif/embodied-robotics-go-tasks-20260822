package domain

import "time"

type StateHistory struct {
	From, To, Actor string
	At              time.Time
	Reason          string
}

func AllowedMissionStates() []string {
	return []string{MissionDraft, MissionApproved, MissionQueued, MissionRunning, MissionSucceeded, MissionFailed, MissionCancelled}
}
func AllowedRobotStates() []string {
	return []string{RobotOffline, RobotReady, RobotBusy, RobotRetired}
}
func AllowedAlertStates() []string { return []string{AlertOpen, AlertAcknowledged, AlertClosed} }
func IsMissionState(s string) bool {
	for _, x := range AllowedMissionStates() {
		if x == s {
			return true
		}
	}
	return false
}
func IsRobotState(s string) bool {
	for _, x := range AllowedRobotStates() {
		if x == s {
			return true
		}
	}
	return false
}
func IsAlertState(s string) bool {
	for _, x := range AllowedAlertStates() {
		if x == s {
			return true
		}
	}
	return false
}
