package domain

import "testing"

func TestRetryStopsAtConfiguredMaximum(t *testing.T) {
	m := Mission{Status: MissionFailed, Attempt: 5}
	if m.Retryable() {
		t.Fatal("retry limit was exceeded")
	}
}
