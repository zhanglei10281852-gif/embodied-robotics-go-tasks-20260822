package httpapi

import (
	"context"
	"encoding/json"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/middleware"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/service"
	"net/http"
	"time"
)

type Server struct {
	Service  *service.Service
	Sessions *SessionStore
}

func New(s *service.Service) *Server { return &Server{Service: s, Sessions: &SessionStore{}} }
func (s *Server) Routes() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", s.health)
	mux.HandleFunc("POST /v1/robots", s.createRobot)
	mux.HandleFunc("POST /v1/missions", s.createMission)
	mux.HandleFunc("POST /v1/missions/{id}/cancel", s.cancelMission)
	mux.HandleFunc("POST /v1/telemetry", s.telemetry)
	return middleware.Recover(middleware.RequestID(mux))
}
func (s *Server) health(w http.ResponseWriter, r *http.Request) {
	write(w, 200, map[string]any{"status": "ok", "time": time.Now().UTC()})
}
func (s *Server) createRobot(w http.ResponseWriter, r *http.Request) {
	var in struct{ TenantID, Serial, Name string }
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	robot, e := s.Service.RegisterRobot(r.Context(), in.TenantID, service.RobotInput{Serial: in.Serial, Name: in.Name})
	if e != nil {
		writeServiceErr(w, e)
		return
	}
	write(w, 201, robot)
}
func (s *Server) createMission(w http.ResponseWriter, r *http.Request) {
	var in struct {
		TenantID, RobotID       string
		Priority, PolicyVersion int
		Steps                   []domain.MissionStep
		IdempotencyKey          string
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	m, e := s.Service.CreateMission(r.Context(), in.TenantID, service.MissionInput{RobotID: in.RobotID, Priority: in.Priority, PolicyVersion: in.PolicyVersion, Steps: in.Steps, IdempotencyKey: in.IdempotencyKey})
	if e != nil {
		writeServiceErr(w, e)
		return
	}
	write(w, 201, m)
}
func (s *Server) cancelMission(w http.ResponseWriter, r *http.Request) {
	tenant := r.URL.Query().Get("tenant")
	if e := s.Service.CancelMission(r.Context(), tenant, r.PathValue("id")); e != nil {
		writeServiceErr(w, e)
		return
	}
	write(w, 200, map[string]string{"status": "cancelled"})
}
func (s *Server) telemetry(w http.ResponseWriter, r *http.Request) {
	var in struct {
		TenantID, RobotID, Kind string
		Sequence                int64
		Payload                 map[string]any
	}
	if json.NewDecoder(r.Body).Decode(&in) != nil {
		writeErr(w, 400, "invalid json")
		return
	}
	e, err := s.Service.RecordTelemetry(r.Context(), in.TenantID, in.RobotID, in.Kind, in.Sequence, in.Payload)
	if err != nil {
		writeServiceErr(w, err)
		return
	}
	write(w, 201, e)
}
func write(w http.ResponseWriter, status int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	_ = json.NewEncoder(w).Encode(v)
}
func writeErr(w http.ResponseWriter, status int, msg string) {
	write(w, status, map[string]string{"error": msg})
}
func writeServiceErr(w http.ResponseWriter, e error) {
	status := 500
	if errors.Is(e, domain.ErrNotFound) {
		status = 404
	}
	if errors.Is(e, domain.ErrConflict) {
		status = 409
	}
	if errors.Is(e, domain.ErrForbidden) {
		status = 403
	}
	writeErr(w, status, e.Error())
}

type SessionStore struct{ sessions map[string]domain.Session }

func (s *SessionStore) Lookup(_ context.Context, token string) (domain.Session, error) {
	if s.sessions == nil {
		return domain.Session{}, domain.ErrRevoked
	}
	v, ok := s.sessions[token]
	if !ok {
		return domain.Session{}, domain.ErrRevoked
	}
	return v, nil
}
func (s *SessionStore) Add(v domain.Session) {
	if s.sessions == nil {
		s.sessions = map[string]domain.Session{}
	}
	s.sessions[v.Token] = v
}
func (s *SessionStore) Revoke(token string) {
	v, ok := s.sessions[token]
	if ok {
		now := time.Now().UTC()
		v.RevokedAt = &now
		s.sessions[token] = v
	}
}

var _ middleware.SessionLookup = (*SessionStore)(nil)
