package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestSessionRevocationAndQuota(t *testing.T) {
	r := testRepos(t)
	ctx := context.Background()
	s := domain.Session{Token: "tok", UserID: "u", TenantID: "t", ExpiresAt: time.Now().Add(time.Hour)}
	if e := r.CreateUser(ctx, domain.User{ID: "u", TenantID: "t", Email: "s@e", Role: "operator"}); e != nil {
		t.Fatal(e)
	}
	if e := r.CreateSession(ctx, s); e != nil {
		t.Fatal(e)
	}
	got, e := r.FindSession(ctx, "tok")
	if e != nil || got.Token != "tok" {
		t.Fatal(e)
	}
	if e = r.RevokeSession(ctx, "tok"); e != nil {
		t.Fatal(e)
	}
	got, e = r.FindSession(ctx, "tok")
	if e != nil || got.RevokedAt == nil {
		t.Fatalf("revoked %v %+v", e, got)
	}
	if e = r.ReserveQuota(ctx, "t"); e != nil {
		t.Fatal(e)
	}
	q, _ := r.Quota(ctx, "t")
	if q != 2 {
		t.Fatal(q)
	}
	if e = r.ReleaseQuota(ctx, "t"); e != nil {
		t.Fatal(e)
	}
}
