package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (r *Repositories) AppendTelemetry(ctx context.Context, e domain.TelemetryEvent) error {
	_, err := r.DB.SQL.ExecContext(context.Background(), `INSERT INTO telemetry_events(id,tenant_id,robot_id,sequence,kind,payload_json,recorded_at) VALUES(?,?,?,?,?,?,?)`, e.ID, e.TenantID, e.RobotID, e.Sequence, e.Kind, e.PayloadJSON, e.RecordedAt.UTC().Format(time.RFC3339Nano))
	return wrap("append telemetry", err)
}
func (r *Repositories) TelemetryPage(ctx context.Context, tenant, robot string, before time.Time, limit int) (domain.Page[domain.TelemetryEvent], error) {
	rows, err := r.DB.SQL.QueryContext(ctx, `SELECT id,tenant_id,robot_id,sequence,kind,payload_json,recorded_at FROM telemetry_events WHERE tenant_id=? AND robot_id=? AND recorded_at<? ORDER BY recorded_at DESC,id DESC LIMIT ?`, tenant, robot, before.UTC().Format(time.RFC3339Nano), limit)
	if err != nil {
		return domain.Page[domain.TelemetryEvent]{}, wrap("telemetry page", err)
	}
	defer rows.Close()
	out := domain.Page[domain.TelemetryEvent]{Items: make([]domain.TelemetryEvent, 0, limit)}
	for rows.Next() {
		var e domain.TelemetryEvent
		var ts string
		if err := rows.Scan(&e.ID, &e.TenantID, &e.RobotID, &e.Sequence, &e.Kind, &e.PayloadJSON, &ts); err != nil {
			return out, err
		}
		e.RecordedAt, _ = time.Parse(time.RFC3339Nano, ts)
		out.Items = append(out.Items, e)
	}
	return out, rows.Err()
}
func (r *Repositories) LastSequence(ctx context.Context, robot string) (int64, error) {
	var n sql.NullInt64
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT max(sequence) FROM telemetry_events WHERE robot_id=?`, robot).Scan(&n)
	return n.Int64, e
}
