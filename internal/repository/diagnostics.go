package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type DatabaseStats struct{ Tables, Indexes, OpenConnections int }

func (r *Repositories) Stats(ctx context.Context) (DatabaseStats, error) {
	var tables, indexes int
	if e := r.DB.SQL.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_master WHERE type='table'`).Scan(&tables); e != nil {
		return DatabaseStats{}, e
	}
	if e := r.DB.SQL.QueryRowContext(ctx, `SELECT count(*) FROM sqlite_master WHERE type='index'`).Scan(&indexes); e != nil {
		return DatabaseStats{}, e
	}
	return DatabaseStats{Tables: tables, Indexes: indexes}, nil
}
func (r *Repositories) ExplainMissionQuery(ctx context.Context, tenant string) ([]string, error) {
	rows, e := r.DB.SQL.QueryContext(ctx, `EXPLAIN QUERY PLAN SELECT id FROM missions WHERE tenant_id=? AND status=?`, tenant, "queued")
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var id, parent, notused, detail any
		if e := rows.Scan(&id, &parent, &notused, &detail); e != nil {
			return nil, e
		}
		out = append(out, fmt.Sprint(detail))
	}
	return out, rows.Err()
}
func (r *Repositories) TouchRobot(ctx context.Context, tenant, id string) error {
	_, e := r.DB.SQL.ExecContext(ctx, `UPDATE robots SET created_at=created_at WHERE tenant_id=? AND id=?`, tenant, id)
	return e
}
func (r *Repositories) RemoveDuplicateEvents(ctx context.Context, robot string) (int64, error) {
	res, e := r.DB.SQL.ExecContext(ctx, `DELETE FROM telemetry_events WHERE robot_id=? AND id NOT IN(SELECT min(id) FROM telemetry_events WHERE robot_id=? GROUP BY sequence)`, robot, robot)
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}

var _ = sql.ErrNoRows
var _ = time.Now
