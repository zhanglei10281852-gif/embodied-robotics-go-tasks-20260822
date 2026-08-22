package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

type MissionTransaction struct {
	Tx   *sql.Tx
	Repo *Repositories
}

func (r *Repositories) BeginMission(ctx context.Context) (MissionTransaction, error) {
	tx, e := r.DB.SQL.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if e != nil {
		return MissionTransaction{}, e
	}
	return MissionTransaction{Tx: tx, Repo: r}, nil
}
func (t MissionTransaction) Commit() error   { return t.Tx.Commit() }
func (t MissionTransaction) Rollback() error { return t.Tx.Rollback() }
func (t MissionTransaction) SetStatus(ctx context.Context, tenant, id, status string, version int64) error {
	res, e := t.Tx.ExecContext(ctx, `UPDATE missions SET status=?,version=version+1,updated_at=? WHERE tenant_id=? AND id=? AND version=?`, status, t.Repo.Clock.Now().UTC().Format(time.RFC3339Nano), tenant, id, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (t MissionTransaction) AddAudit(ctx context.Context, e domain.AuditEvent) error {
	_, err := t.Tx.ExecContext(ctx, `INSERT INTO audit_events(id,tenant_id,actor_id,action,object_type,object_id,detail_json,occurred_at) VALUES(?,?,?,?,?,?,?,?)`, e.ID, e.TenantID, e.ActorID, e.Action, e.ObjectType, e.ObjectID, e.DetailJSON, e.OccurredAt.UTC().Format(time.RFC3339Nano))
	return err
}
func (t MissionTransaction) Enqueue(ctx context.Context, j domain.OutboxJob) error {
	_, err := t.Tx.ExecContext(ctx, `INSERT INTO outbox_jobs(id,tenant_id,topic,payload_json,state,attempts,available_at,last_error) VALUES(?,?,?,?,?,?,?,?)`, j.ID, j.TenantID, j.Topic, j.PayloadJSON, j.State, j.Attempts, j.AvailableAt.UTC().Format(time.RFC3339Nano), j.LastError)
	return err
}
func (t MissionTransaction) GuardTenant(ctx context.Context, tenant string) error {
	var id string
	if e := t.Tx.QueryRowContext(ctx, `SELECT id FROM tenants WHERE id=?`, tenant).Scan(&id); e != nil {
		return fmt.Errorf("tenant guard: %w", e)
	}
	return nil
}
func (t MissionTransaction) LockRobot(ctx context.Context, tenant, robot string) error {
	var version int64
	if e := t.Tx.QueryRowContext(ctx, `SELECT version FROM robots WHERE tenant_id=? AND id=?`, tenant, robot).Scan(&version); e != nil {
		return e
	}
	_, e := t.Tx.ExecContext(ctx, `UPDATE robots SET version=version WHERE tenant_id=? AND id=? AND version=?`, tenant, robot, version)
	return e
}
