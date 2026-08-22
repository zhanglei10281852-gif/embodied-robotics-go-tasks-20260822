package audit

import (
	"context"
	"encoding/json"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

func (l *Logger) Export(ctx context.Context, tenant string, limit int) ([]byte, error) {
	events, e := l.Repos.ListAudit(ctx, tenant, limit)
	if e != nil {
		return nil, e
	}
	return json.Marshal(events)
}
func Summarize(events []domain.AuditEvent) map[string]int {
	out := map[string]int{}
	for _, e := range events {
		out[e.Action]++
	}
	return out
}
