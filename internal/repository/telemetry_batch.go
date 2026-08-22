package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
)

func (r *Repositories) AppendTelemetryBatch(ctx context.Context, events []domain.TelemetryEvent) error {
	if len(events) == 0 {
		return nil
	}
	return r.DB.WithTx(ctx, func(tx *sql.Tx) error {
		for _, e := range events {
			if err := e.Validate(); err != nil {
				return err
			}
			if _, err := tx.ExecContext(ctx, `INSERT INTO telemetry_events(id,tenant_id,robot_id,sequence,kind,payload_json,recorded_at) VALUES(?,?,?,?,?,?,?)`, e.ID, e.TenantID, e.RobotID, e.Sequence, e.Kind, e.PayloadJSON, e.RecordedAt.UTC().Format("2006-01-02T15:04:05.999999999Z07:00")); err != nil {
				return fmt.Errorf("telemetry batch: %w", err)
			}
		}
		return nil
	})
}
func (r *Repositories) DeleteRobotTelemetry(ctx context.Context, tenant, robot string, before string) (int64, error) {
	res, e := r.DB.SQL.ExecContext(ctx, `DELETE FROM telemetry_events WHERE tenant_id=? AND robot_id=? AND recorded_at<?`, tenant, robot, before)
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
func (r *Repositories) SequenceGap(ctx context.Context, robot string) (bool, error) {
	var min, max, count int64
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT COALESCE(min(sequence),0),COALESCE(max(sequence),0),count(*) FROM telemetry_events WHERE robot_id=?`, robot).Scan(&min, &max, &count)
	if e != nil {
		return false, e
	}
	return max-min+1 != count, nil
}
