package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type StatusCount struct {
	Status string
	Count  int
}

func (r *Repositories) MissionStatusCounts(ctx context.Context, tenant string) ([]StatusCount, error) {
	rows, e := r.DB.SQL.QueryContext(ctx, `SELECT status,count(*) FROM missions WHERE tenant_id=? GROUP BY status ORDER BY status`, tenant)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []StatusCount{}
	for rows.Next() {
		var x StatusCount
		if e := rows.Scan(&x.Status, &x.Count); e != nil {
			return nil, e
		}
		out = append(out, x)
	}
	return out, rows.Err()
}
func (r *Repositories) RobotLastSeen(ctx context.Context, tenant string) (map[string]time.Time, error) {
	rows, e := r.DB.SQL.QueryContext(ctx, `SELECT robot_id,max(recorded_at) FROM telemetry_events WHERE tenant_id=? GROUP BY robot_id`, tenant)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := map[string]time.Time{}
	for rows.Next() {
		var id, ts string
		if e := rows.Scan(&id, &ts); e != nil {
			return nil, e
		}
		t, e := time.Parse(time.RFC3339Nano, ts)
		if e != nil {
			return nil, e
		}
		out[id] = t
	}
	return out, rows.Err()
}
func (r *Repositories) DeleteExpiredSessions(ctx context.Context, now time.Time) (int64, error) {
	res, e := r.DB.SQL.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at<?`, now.UTC().Format(time.RFC3339Nano))
	if e != nil {
		return 0, fmt.Errorf("delete sessions: %w", e)
	}
	return res.RowsAffected()
}
func (r *Repositories) Health(ctx context.Context) error {
	var one int
	return r.DB.SQL.QueryRowContext(ctx, `SELECT 1`).Scan(&one)
}

var _ sql.Result
