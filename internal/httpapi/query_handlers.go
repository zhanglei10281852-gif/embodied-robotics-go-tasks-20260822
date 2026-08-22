package httpapi

import (
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"net/http"
	"time"
)

func (s *Server) queryMissions(w http.ResponseWriter, r *http.Request) {
	tenant := r.URL.Query().Get("tenant")
	f := repository.MissionFilter{Status: r.URL.Query().Get("status"), RobotID: r.URL.Query().Get("robot"), Limit: decodeLimit(r, 50, 500), Before: parseTimeQuery(r, "before", time.Now().UTC().Add(time.Second))}
	page, e := s.Service.Search(r.Context(), tenant, f)
	if e != nil {
		writeServiceErr(w, e)
		return
	}
	write(w, 200, page)
}
func (s *Server) missionDetails(w http.ResponseWriter, r *http.Request) {
	v, e := s.Service.MissionView(r.Context(), r.URL.Query().Get("tenant"), r.PathValue("id"))
	if e != nil {
		writeServiceErr(w, e)
		return
	}
	write(w, 200, map[string]any{"mission": missionSummary(v.Mission), "robot": robotSummary(v.Robot), "steps": v.Steps})
}
