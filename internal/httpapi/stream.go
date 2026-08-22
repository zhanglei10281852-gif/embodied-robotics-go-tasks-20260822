package httpapi

import (
	"bufio"
	"encoding/json"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/worker"
	"net/http"
	"time"
)

type Stream struct{ Bus *worker.Bus }

func (s Stream) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if s.Bus == nil {
		writeErrorEnvelope(w, 500, "bus_missing", "event bus unavailable", nil)
		return
	}
	w.Header().Set("Content-Type", "text/event-stream")
	w.Header().Set("Cache-Control", "no-cache")
	flusher, ok := w.(http.Flusher)
	if !ok {
		writeErrorEnvelope(w, 500, "stream_unsupported", "stream unsupported", nil)
		return
	}
	ch := s.Bus.Subscribe()
	defer s.Bus.Unsubscribe(ch)
	enc := json.NewEncoder(bufio.NewWriter(w))
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-r.Context().Done():
			return
		case ev, ok := <-ch:
			if !ok {
				return
			}
			_ = enc.Encode(map[string]any{"topic": ev.Topic, "payload": string(ev.Payload)})
			flusher.Flush()
		case <-ticker.C:
			_, _ = w.Write([]byte(": heartbeat\n\n"))
			flusher.Flush()
		}
	}
}
func StreamHeaders(w http.ResponseWriter) {
	w.Header().Set("Connection", "keep-alive")
	w.Header().Set("X-Accel-Buffering", "no")
}
