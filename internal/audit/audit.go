package audit

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"time"
)

type Logger struct{ Repos *repository.Repositories }

func New(r *repository.Repositories) *Logger { return &Logger{Repos: r} }
func (l *Logger) Record(ctx context.Context, tenant, actor, action, typ, obj string, detail any) error {
	b, e := json.Marshal(detail)
	if e != nil {
		return fmt.Errorf("audit detail: %w", e)
	}
	return l.Repos.AppendAudit(ctx, domain.AuditEvent{ID: fmt.Sprintf("audit-%d", time.Now().UnixNano()), TenantID: tenant, ActorID: actor, Action: action, ObjectType: typ, ObjectID: obj, DetailJSON: string(b), OccurredAt: time.Now().UTC()})
}
func (l *Logger) RecordTransition(ctx context.Context, tenant, actor, object string, from, to string) error {
	return l.Record(ctx, tenant, actor, "state.transition", "mission", object, map[string]string{"from": from, "to": to})
}
