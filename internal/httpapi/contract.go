package httpapi

import (
	"encoding/json"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"net/http"
	"strconv"
	"strings"
	"time"
)

type ErrorEnvelope struct {
	Code      string            `json:"code"`
	Message   string            `json:"message"`
	RequestID string            `json:"request_id,omitempty"`
	Details   map[string]string `json:"details,omitempty"`
}

func decodeLimit(r *http.Request, def, max int) int {
	n, err := strconv.Atoi(r.URL.Query().Get("limit"))
	if err != nil || n < 1 {
		return def
	}
	if n > max {
		return max
	}
	return n
}
func parseTimeQuery(r *http.Request, key string, def time.Time) time.Time {
	v := strings.TrimSpace(r.URL.Query().Get(key))
	if v == "" {
		return def
	}
	t, e := time.Parse(time.RFC3339Nano, v)
	if e != nil {
		return def
	}
	return t
}
func writeErrorEnvelope(w http.ResponseWriter, status int, code, msg string, details map[string]string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(ErrorEnvelope{Code: code, Message: msg, RequestID: w.Header().Get("X-Request-ID"), Details: details})
}
func missionSummary(m domain.Mission) map[string]any {
	return map[string]any{"id": m.ID, "status": m.Status, "priority": m.Priority, "version": m.Version, "updated_at": m.UpdatedAt.UTC().Format(time.RFC3339Nano)}
}
func robotSummary(r domain.Robot) map[string]any {
	return map[string]any{"id": r.ID, "serial": r.Serial, "name": r.Name, "status": r.Status, "version": r.Version}
}
