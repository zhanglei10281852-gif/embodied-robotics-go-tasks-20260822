package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestMissionTransactionCommit(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	r.Robots().Create(ctx, domain.Robot{ID: "tr", TenantID: "t", Serial: "TR", Name: "tx", Status: domain.RobotReady, Version: 1})
	m := domain.Mission{ID: "tm", TenantID: "t", RobotID: "tr", Status: domain.MissionDraft, Priority: 1, PolicyVersion: 1, Version: 1, CreatedAt: time.Now(), UpdatedAt: time.Now()}
	if e := r.Missions().Create(ctx, m, nil); e != nil {
		t.Fatal(e)
	}
	tx, e := r.BeginMission(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if e = tx.SetStatus(ctx, "t", "tm", domain.MissionApproved, 1); e != nil {
		t.Fatal(e)
	}
	if e = tx.Rollback(); e != nil {
		t.Fatal(e)
	}
}
func TestMissionTransactionRollback(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	tx, e := r.BeginMission(ctx)
	if e != nil {
		t.Fatal(e)
	}
	if e = tx.GuardTenant(ctx, "t"); e != nil {
		t.Fatal(e)
	}
	if e = tx.Rollback(); e != nil {
		t.Fatal(e)
	}
}
