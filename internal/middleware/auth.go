package middleware

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"net/http"
	"strings"
	"time"
)

type SessionLookup interface {
	Lookup(context.Context, string) (domain.Session, error)
}
type ContextKey string

const SessionKey ContextKey = "robotics.session"

func WithSession(ctx context.Context, s domain.Session) context.Context {
	return context.WithValue(ctx, SessionKey, s)
}
func FromContext(ctx context.Context) (domain.Session, bool) {
	s, ok := ctx.Value(SessionKey).(domain.Session)
	return s, ok
}
func Bearer(header string) (string, error) {
	parts := strings.Fields(header)
	if len(parts) != 2 || !strings.EqualFold(parts[0], "Bearer") || parts[1] == "" {
		return "", errors.New("invalid bearer token")
	}
	return parts[1], nil
}
func RequireSession(lookup SessionLookup, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token, e := Bearer(r.Header.Get("Authorization"))
		if e != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		s, e := lookup.Lookup(r.Context(), token)
		if e != nil || s.Valid(time.Now().UTC()) != nil {
			http.Error(w, "unauthorized", http.StatusUnauthorized)
			return
		}
		next.ServeHTTP(w, r.WithContext(WithSession(r.Context(), s)))
	})
}
