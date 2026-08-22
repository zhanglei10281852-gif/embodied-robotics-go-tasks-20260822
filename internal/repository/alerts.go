package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (r *Repositories) UpsertAlert(ctx context.Context, a domain.Alert) error {
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO alerts(id,tenant_id,robot_id,dedupe_key,severity,message,state,created_at) VALUES(?,?,?,?,?,?,?,?) ON CONFLICT(tenant_id,dedupe_key) DO UPDATE SET severity=excluded.severity,message=excluded.message,state=excluded.state`, a.ID, a.TenantID, a.RobotID, a.DedupeKey, a.Severity, a.Message, a.State, a.CreatedAt.UTC().Format(time.RFC3339Nano))
	return wrap("upsert alert", e)
}
func (r *Repositories) GetAlert(ctx context.Context, tenant, id string) (domain.Alert, error) {
	var a domain.Alert
	var created, acked sql.NullString
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT id,tenant_id,robot_id,dedupe_key,severity,message,state,created_at,acknowledged_at FROM alerts WHERE tenant_id=? AND id=?`, tenant, id).Scan(&a.ID, &a.TenantID, &a.RobotID, &a.DedupeKey, &a.Severity, &a.Message, &a.State, &created, &acked)
	if e == sql.ErrNoRows {
		return a, domain.ErrNotFound
	}
	a.CreatedAt, _ = time.Parse(time.RFC3339Nano, created.String)
	if acked.Valid {
		v, _ := time.Parse(time.RFC3339Nano, acked.String)
		a.AcknowledgedAt = &v
	}
	return a, e
}
func (r *Repositories) AckAlert(ctx context.Context, tenant, id string, at time.Time) error {
	res, e := r.DB.SQL.ExecContext(ctx, `UPDATE alerts SET state=?,acknowledged_at=? WHERE tenant_id=? AND id=? AND state=?`, domain.AlertAcknowledged, at.UTC().Format(time.RFC3339Nano), tenant, id, domain.AlertOpen)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r *Repositories) CloseAlert(ctx context.Context, tenant, id string) error {
	res, e := r.DB.SQL.ExecContext(ctx, `UPDATE alerts SET state=? WHERE tenant_id=? AND id=? AND state=?`, domain.AlertClosed, tenant, id, domain.AlertAcknowledged)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
