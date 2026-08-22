package repository

import (
	"context"
	"database/sql"
	"time"
)

func (r *Repositories) FindIdempotent(ctx context.Context, tenant, key string) (string, error) {
	var id string
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT resource_id FROM idempotency_keys WHERE key=?`, key).Scan(&id)
	if e == sql.ErrNoRows {
		return "", nil
	}
	return id, e
}
func (r *Repositories) PutIdempotent(ctx context.Context, tenant, key, id string) error {
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO idempotency_keys(tenant_id,key,resource_id,created_at) VALUES(?,?,?,?)`, tenant, key, id, r.Clock.Now().UTC().Format(time.RFC3339Nano))
	return e
}
