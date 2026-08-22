package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

type MaintenanceReport struct{ Sessions, Outbox, Alerts int64 }

func (r *Repositories) PurgeExpired(ctx context.Context, now time.Time) (MaintenanceReport, error) {
	var out MaintenanceReport
	err := r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		res, e := tx.ExecContext(ctx, `DELETE FROM sessions WHERE expires_at<?`, now.UTC().Format(time.RFC3339Nano))
		if e != nil {
			return e
		}
		out.Sessions, _ = res.RowsAffected()
		res, e = tx.ExecContext(ctx, `DELETE FROM outbox_jobs WHERE state=? AND available_at<?`, domain.OutboxDone, now.Add(-24*time.Hour).UTC().Format(time.RFC3339Nano))
		if e != nil {
			return e
		}
		out.Outbox, _ = res.RowsAffected()
		res, e = tx.ExecContext(ctx, `DELETE FROM alerts WHERE state=? AND created_at<?`, domain.AlertClosed, now.Add(-30*24*time.Hour).UTC().Format(time.RFC3339Nano))
		if e != nil {
			return e
		}
		out.Alerts, _ = res.RowsAffected()
		return nil
	})
	return out, err
}
func (r *Repositories) Vacuum(ctx context.Context) error {
	if _, e := r.DB.SQL.ExecContext(ctx, `PRAGMA optimize`); e != nil {
		return fmt.Errorf("optimize: %w", e)
	}
	return nil
}
func (r *Repositories) CheckForeignKeys(ctx context.Context) error {
	rows, e := r.DB.SQL.QueryContext(ctx, `PRAGMA foreign_key_check`)
	if e != nil {
		return e
	}
	defer rows.Close()
	if rows.Next() {
		return fmt.Errorf("foreign key violation")
	}
	return rows.Err()
}
