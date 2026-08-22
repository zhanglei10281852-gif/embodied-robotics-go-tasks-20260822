package httpapi

import (
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/middleware"
	"net/http"
	"strings"
)

func TenantHeader(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		tenant := strings.TrimSpace(r.Header.Get("X-Tenant-ID"))
		if tenant == "" {
			writeErrorEnvelope(w, 400, "tenant_missing", "tenant header required", nil)
			return
		}
		next.ServeHTTP(w, r)
	})
}
func NoCache(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Cache-Control", "no-store")
		next.ServeHTTP(w, r)
	})
}
func Compose(h http.Handler) http.Handler {
	return middleware.Recover(middleware.RequestID(NoCache(TenantHeader(h))))
}
