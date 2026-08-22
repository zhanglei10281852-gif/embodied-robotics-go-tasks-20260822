package middleware

import (
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestRevokedSessionCannotBeValidated(t *testing.T) {
	now := time.Now().UTC()
	s := domain.Session{Token: "revoked", UserID: "u", TenantID: "t", ExpiresAt: now.Add(time.Minute), RevokedAt: &now}
	if err := s.Valid(now); err != domain.ErrRevoked {
		t.Fatalf("revoked session accepted: %v", err)
	}
}
