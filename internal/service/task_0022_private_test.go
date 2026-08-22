package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestFormationKeepsEveryMember(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	for _, r := range []struct{ id, serial string }{{"f1", "F1"}, {"f2", "F2"}} {
		if err := s.Repos.Robots().Create(ctx, domain.Robot{ID: r.id, TenantID: "tenant", Serial: r.serial, Name: r.id, Status: domain.RobotReady, Version: 1}); err != nil {
			t.Fatal(err)
		}
	}
	f, err := s.CreateFormation(ctx, "tenant", "pair", []string{"f1", "f2"})
	if err != nil {
		t.Fatal(err)
	}
	var count int
	if err := s.Repos.DB.SQL.QueryRowContext(ctx, "SELECT count(*) FROM formation_members WHERE formation_id=?", f.ID).Scan(&count); err != nil {
		t.Fatal(err)
	}
	if count != 2 {
		t.Fatalf("formation lost member: %d", count)
	}
}
