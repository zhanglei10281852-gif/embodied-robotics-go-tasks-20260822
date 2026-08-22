package service

import (
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestMergeTelemetryDoesNotAliasCallerSlices(t *testing.T) {
	first := make([]domain.TelemetryEvent, 1, 2)
	first[0] = domain.TelemetryEvent{ID: "first", RecordedAt: time.Now()}
	merged := New(nil).MergeTelemetry(first, []domain.TelemetryEvent{{ID: "second", RecordedAt: time.Now().Add(time.Second)}})
	merged[0].ID = "changed"
	if first[0].ID != "first" {
		t.Fatalf("caller slice was mutated: %+v", first)
	}
}
