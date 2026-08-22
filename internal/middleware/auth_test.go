package middleware

import (
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestBearerAndSession(t *testing.T) {
	if tok, e := Bearer("Bearer abc"); e != nil || tok != "abc" {
		t.Fatalf("%v %s", e, tok)
	}
	if _, e := Bearer("Basic abc"); e == nil {
		t.Fatal("scheme accepted")
	}
	s := domain.Session{Token: "x", UserID: "u", TenantID: "t", ExpiresAt: time.Now().Add(time.Minute)}
	if e := s.Valid(time.Now()); e != nil {
		t.Fatal(e)
	}
}
