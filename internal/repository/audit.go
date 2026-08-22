package repository

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (r *Repositories) AppendAudit(ctx context.Context, e domain.AuditEvent) error {
	_, err := r.DB.SQL.ExecContext(ctx, `INSERT INTO audit_events(id,tenant_id,actor_id,action,object_type,object_id,detail_json,occurred_at) VALUES(?,?,?,?,?,?,?,?)`, e.ID, e.TenantID, e.ActorID, e.Action, e.ObjectType, e.ObjectID, e.DetailJSON, e.OccurredAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00"))
	return wrap("append audit", err)
}
func (r *Repositories) ListAudit(ctx context.Context, tenant string, limit int) ([]domain.AuditEvent, error) {
	rows, e := r.DB.SQL.QueryContext(ctx, `SELECT id,tenant_id,actor_id,action,object_type,object_id,detail_json,occurred_at FROM audit_events WHERE tenant_id=? ORDER BY occurred_at DESC,id DESC LIMIT ?`, tenant, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.AuditEvent{}
	for rows.Next() {
		var a domain.AuditEvent
		var ts string
		if e := rows.Scan(&a.ID, &a.TenantID, &a.ActorID, &a.Action, &a.ObjectType, &a.ObjectID, &a.DetailJSON, &ts); e != nil {
			return nil, e
		}
		a.OccurredAt, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, a)
	}
	return out, rows.Err()
}
