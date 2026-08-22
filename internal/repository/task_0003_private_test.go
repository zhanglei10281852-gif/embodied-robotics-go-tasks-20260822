package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestTransitionRejectsStaleSourceState(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	if err := r.Robots().Create(ctx, domain.Robot{ID: "stale-r", TenantID: "t", Serial: "STALE", Name: "stale", Status: domain.RobotReady, Version: 1}); err != nil {
		t.Fatal(err)
	}
	m := domain.Mission{ID: "stale-m", TenantID: "t", RobotID: "stale-r", Status: domain.MissionDraft, Priority: 1, PolicyVersion: 1, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := r.Missions().Create(ctx, m, nil); err != nil {
		t.Fatal(err)
	}
	if err := r.Missions().Transition(ctx, "t", m.ID, domain.MissionRunning, domain.MissionCancelled, 1); err == nil {
		t.Fatal("stale transition was accepted")
	}
}
