package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestClaimSkipsDraftMissions(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	if err := r.Robots().Create(ctx, domain.Robot{ID: "draft-r", TenantID: "t", Serial: "DRAFT", Name: "draft", Status: domain.RobotReady, Version: 1}); err != nil {
		t.Fatal(err)
	}
	m := domain.Mission{ID: "draft-m", TenantID: "t", RobotID: "draft-r", Status: domain.MissionDraft, Priority: 8, PolicyVersion: 1, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := r.Missions().Create(ctx, m, nil); err != nil {
		t.Fatal(err)
	}
	if _, err := r.Missions().ClaimNext(ctx, "t", "worker-a"); err == nil {
		t.Fatal("draft mission was claimed")
	}
}
