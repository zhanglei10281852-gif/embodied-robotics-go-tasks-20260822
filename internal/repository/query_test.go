package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestMissionSearchAndSteps(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	r.Robots().Create(ctx, domain.Robot{ID: "qr", TenantID: "t", Serial: "Q", Name: "Q", Status: domain.RobotReady, Version: 1})
	m := domain.Mission{ID: "qm", TenantID: "t", RobotID: "qr", Status: domain.MissionDraft, Priority: 9, PolicyVersion: 1, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if e := r.Missions().Create(ctx, m, []domain.MissionStep{{ID: "qs", MissionID: "qm", Ordinal: 0, Action: "scan", PayloadJSON: "{}", Status: domain.StepPending}}); e != nil {
		t.Fatal(e)
	}
	page, e := r.SearchMissions(ctx, "t", MissionFilter{Status: domain.MissionDraft, Limit: 10})
	if e != nil || len(page.Items) != 1 {
		t.Fatalf("page %v %+v", e, page)
	}
	steps, e := r.ListSteps(ctx, "t", "qm")
	if e != nil || len(steps) != 1 {
		t.Fatal(e)
	}
}
