package httpapi

import (
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"net/http/httptest"
	"testing"
	"time"
)

func TestSessionStore(t *testing.T) {
	s := &SessionStore{}
	now := time.Now().Add(time.Minute)
	s.Add(domain.Session{Token: "a", UserID: "u", TenantID: "t", ExpiresAt: now})
	got, e := s.Lookup(nil, "a")
	if e != nil || got.UserID != "u" {
		t.Fatalf("%v %+v", e, got)
	}
	s.Revoke("a")
	got, e = s.Lookup(nil, "a")
	if e != nil || got.RevokedAt == nil {
		t.Fatal("revoke")
	}
}
func TestErrorEnvelope(t *testing.T) {
	rw := httptest.NewRecorder()
	writeErrorEnvelope(rw, 409, "conflict", "version", map[string]string{"field": "version"})
	if rw.Code != 409 || rw.Body.Len() == 0 {
		t.Fatal("envelope")
	}
}
