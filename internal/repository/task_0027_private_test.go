package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestMissionSearchWithoutCursorIncludesFutureRows(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	if err := r.Robots().Create(ctx, domain.Robot{ID: "future-r", TenantID: "t", Serial: "FUTURE", Name: "future", Status: domain.RobotReady, Version: 1}); err != nil {
		t.Fatal(err)
	}
	m := domain.Mission{ID: "future-m", TenantID: "t", RobotID: "future-r", Status: domain.MissionDraft, Priority: 1, PolicyVersion: 1, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if err := r.Missions().Create(ctx, m, nil); err != nil {
		t.Fatal(err)
	}
	_, _ = r.DB.SQL.ExecContext(ctx, "UPDATE missions SET updated_at=? WHERE id=?", time.Now().Add(time.Hour).UTC().Format(time.RFC3339Nano), m.ID)
	page, err := r.SearchMissions(ctx, "t", MissionFilter{Limit: 10})
	if err != nil {
		t.Fatal(err)
	}
	if len(page.Items) != 1 {
		t.Fatalf("future mission was filtered: %+v", page.Items)
	}
}
