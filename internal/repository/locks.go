package repository

import (
	"context"
	"database/sql"
	"fmt"
	"time"
)

type LockLease struct {
	Key, Owner string
	Until      time.Time
}

func (r *Repositories) AcquireLock(ctx context.Context, tenant, key, owner string, ttl time.Duration) error {
	if ttl <= 0 {
		return fmt.Errorf("lock ttl positive")
	}
	now := r.Clock.Now().UTC()
	_, e := r.DB.SQL.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS resource_locks(tenant_id TEXT NOT NULL,key TEXT NOT NULL,owner TEXT NOT NULL,until_at TEXT NOT NULL,PRIMARY KEY(tenant_id,key))`)
	if e != nil {
		return e
	}
	tx, e := r.DB.SQL.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if e != nil {
		return e
	}
	defer tx.Rollback()
	var until string
	scan := tx.QueryRowContext(ctx, `SELECT until_at FROM resource_locks WHERE tenant_id=? AND key=?`, tenant, key).Scan(&until)
	if scan == nil {
		when, pe := time.Parse(time.RFC3339Nano, until)
		if pe == nil && when.After(now) {
			return fmt.Errorf("lock held")
		}
	}
	_, e = tx.ExecContext(ctx, `INSERT INTO resource_locks(tenant_id,key,owner,until_at) VALUES(?,?,?,?) ON CONFLICT(tenant_id,key) DO UPDATE SET owner=excluded.owner,until_at=excluded.until_at`, tenant, key, owner, now.Add(ttl).Format(time.RFC3339Nano))
	if e != nil {
		return e
	}
	return tx.Commit()
}
func (r *Repositories) ReleaseLock(ctx context.Context, tenant, key, owner string) error {
	_, e := r.DB.SQL.ExecContext(ctx, `DELETE FROM resource_locks WHERE tenant_id=? AND key=? AND owner=?`, tenant, key, owner)
	return e
}
func (r *Repositories) ListLocks(ctx context.Context, tenant string) ([]LockLease, error) {
	rows, e := r.DB.SQL.QueryContext(ctx, `SELECT key,owner,until_at FROM resource_locks WHERE tenant_id=? ORDER BY key`, tenant)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []LockLease{}
	for rows.Next() {
		var l LockLease
		var ts string
		if e := rows.Scan(&l.Key, &l.Owner, &ts); e != nil {
			return nil, e
		}
		l.Until, _ = time.Parse(time.RFC3339Nano, ts)
		out = append(out, l)
	}
	return out, rows.Err()
}
