package service

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestSignerAndAuthorize(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	u := domain.User{ID: "sec", TenantID: "tenant", Email: "sec@e", Role: "operator"}
	if e := s.Repos.CreateUser(ctx, u); e != nil {
		t.Fatal(e)
	}
	session, e := s.CreateSession(ctx, u, time.Minute)
	if e != nil {
		t.Fatal(e)
	}
	signer := TokenSigner{Secret: []byte("secret")}
	token := signer.Sign(session)
	if signer.Verify(token) != session.Token {
		t.Fatal("signature")
	}
	if signer.Verify(token+"x") != "" {
		t.Fatal("tamper")
	}
	if e = s.Authorize(ctx, session, "tenant", "operator"); e != nil {
		t.Fatal(e)
	}
	if e = s.Authorize(ctx, session, "other", "operator"); e == nil {
		t.Fatal("cross tenant")
	}
}
