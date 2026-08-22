package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (r *Repositories) Enqueue(ctx context.Context, j domain.OutboxJob) error {
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO outbox_jobs(id,tenant_id,topic,payload_json,state,attempts,available_at,last_error) VALUES(?,?,?,?,?,?,?,?)`, j.ID, j.TenantID, j.Topic, j.PayloadJSON, j.State, j.Attempts, j.AvailableAt.UTC().Format(time.RFC3339Nano), j.LastError)
	return wrap("enqueue outbox", e)
}
func (r *Repositories) ClaimOutbox(ctx context.Context, limit int) ([]domain.OutboxJob, error) {
	var out []domain.OutboxJob
	err := r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		rows, e := tx.QueryContext(ctx, `SELECT id,tenant_id,topic,payload_json,state,attempts,available_at,last_error FROM outbox_jobs WHERE state=? AND available_at<=? ORDER BY available_at,id LIMIT ?`, domain.OutboxPending, r.Clock.Now().UTC().Format(time.RFC3339Nano), limit)
		if e != nil {
			return e
		}
		defer rows.Close()
		for rows.Next() {
			var j domain.OutboxJob
			var ts, last sql.NullString
			if e := rows.Scan(&j.ID, &j.TenantID, &j.Topic, &j.PayloadJSON, &j.State, &j.Attempts, &ts, &last); e != nil {
				return e
			}
			j.AvailableAt, _ = time.Parse(time.RFC3339Nano, ts.String)
			j.LastError = last.String
			out = append(out, j)
			if _, e = tx.ExecContext(ctx, `UPDATE outbox_jobs SET state=? WHERE id=? AND state=?`, domain.OutboxSending, j.ID, domain.OutboxPending); e != nil {
				return e
			}
		}
		return rows.Err()
	})
	return out, err
}
func (r *Repositories) FinishOutbox(ctx context.Context, id string, errValue error) error {
	if errValue == nil {
		_, e := r.DB.SQL.ExecContext(ctx, `UPDATE outbox_jobs SET state=? WHERE id=?`, domain.OutboxDone, id)
		return e
	}
	_, e := r.DB.SQL.ExecContext(ctx, `UPDATE outbox_jobs SET state=?,attempts=attempts+1,last_error=?,available_at=? WHERE id=?`, domain.OutboxDone, errValue.Error(), r.Clock.Now().Add(time.Second).UTC().Format(time.RFC3339Nano), id)
	return e
}
