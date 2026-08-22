package service

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"strings"
	"time"
)

type TokenSigner struct {
	Secret []byte
	TTL    time.Duration
}

func (t TokenSigner) Sign(session domain.Session) string {
	mac := hmac.New(sha256.New, t.Secret)
	mac.Write([]byte(session.Token))
	return base64.RawURLEncoding.EncodeToString(mac.Sum(nil)) + "." + session.Token
}
func (t TokenSigner) Verify(value string) string {
	parts := strings.SplitN(value, ".", 2)
	if len(parts) != 2 {
		return ""
	}
	mac := hmac.New(sha256.New, t.Secret)
	mac.Write([]byte(parts[1]))
	want, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil || !hmac.Equal(want, mac.Sum(nil)) {
		return ""
	}
	return parts[1]
}
func (s *Service) Authorize(ctx context.Context, session domain.Session, tenant, role string) error {
	if e := s.ValidateTenantAccess(session, tenant); e != nil {
		return e
	}
	if role != "" && role != "operator" && role != "supervisor" && role != "auditor" {
		return domain.ErrForbidden
	}
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return nil
	}
}
func (s *Service) Revoke(ctx context.Context, token string) error {
	if token == "" {
		return errors.New("token missing")
	}
	return s.Repos.RevokeSession(ctx, token)
}
func (s *Service) CreateSession(ctx context.Context, user domain.User, ttl time.Duration) (domain.Session, error) {
	if ttl <= 0 {
		return domain.Session{}, domain.ErrExpired
	}
	session := domain.Session{Token: id("session"), UserID: user.ID, TenantID: user.TenantID, ExpiresAt: s.now().Add(ttl)}
	if e := s.Repos.CreateSession(ctx, session); e != nil {
		return domain.Session{}, e
	}
	return session, nil
}
