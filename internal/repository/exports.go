package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

type ExportRow struct {
	ID, Kind, Payload string
	At                time.Time
}

func (r *Repositories) ExportTelemetry(ctx context.Context, tenant, robot string, limit int) ([]ExportRow, error) {
	if limit < 1 {
		limit = 100
	}
	rows, e := r.DB.SQL.QueryContext(ctx, `SELECT id,kind,payload_json,recorded_at FROM telemetry_events WHERE tenant_id=? AND robot_id=? ORDER BY recorded_at,id LIMIT ?`, tenant, robot, limit)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []ExportRow{}
	for rows.Next() {
		var x ExportRow
		var ts string
		if e := rows.Scan(&x.ID, &x.Kind, &x.Payload, &ts); e != nil {
			return nil, e
		}
		x.At, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repositories) ExportTelemetryJSON(ctx context.Context, tenant, robot string, limit int) ([]byte, error) {
	rows, e := r.ExportTelemetry(ctx, tenant, robot, limit)
	if e != nil {
		return nil, e
	}
	b, e := json.Marshal(rows)
	if e != nil {
		return nil, fmt.Errorf("marshal telemetry export: %w", e)
	}
	return b, nil
}
func (r *Repositories) InsertAuditAndOutbox(ctx context.Context, a domain.AuditEvent, j domain.OutboxJob) error {
	return r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		if _, e := tx.ExecContext(ctx, `INSERT INTO audit_events(id,tenant_id,actor_id,action,object_type,object_id,detail_json,occurred_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.TenantID, a.ActorID, a.Action, a.ObjectType, a.ObjectID, a.DetailJSON, a.OccurredAt.UTC().Format(time.RFC3339Nano)); e != nil {
			return e
		}
		_, e := tx.ExecContext(ctx, `INSERT INTO outbox_jobs(id,tenant_id,topic,payload_json,state,attempts,available_at,last_error) VALUES(?,?,?,?,?,?,?,?)`, j.ID, j.TenantID, j.Topic, j.PayloadJSON, j.State, j.Attempts, j.AvailableAt.UTC().Format(time.RFC3339Nano), j.LastError)
		return e
	})
}
