package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestAggregatorResetRemovesRobotState(t *testing.T) {
	a := NewAggregator()
	if err := a.Add(context.Background(), domain.TelemetryEvent{RobotID: "robot", Sequence: 1}); err != nil {
		t.Fatal(err)
	}
	a.Reset("robot")
	if got := a.Snapshot("robot"); len(got) != 0 {
		t.Fatalf("stale aggregate remained: %+v", got)
	}
}
