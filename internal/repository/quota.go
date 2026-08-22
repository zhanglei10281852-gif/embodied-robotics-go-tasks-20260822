package repository

import (
	"context"
	"database/sql"
)

func (r *Repositories) ReserveQuota(ctx context.Context, tenant string) error {
	res, e := r.DB.SQL.ExecContext(ctx, `UPDATE tenants SET quota=quota-1 WHERE id=?`, tenant)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return sql.ErrNoRows
	}
	return nil
}
func (r *Repositories) ReleaseQuota(ctx context.Context, tenant string) error {
	_, e := r.DB.SQL.ExecContext(ctx, `UPDATE tenants SET quota=quota+1 WHERE id=?`, tenant)
	return e
}
func (r *Repositories) Quota(ctx context.Context, tenant string) (int, error) {
	var n int
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT quota FROM tenants WHERE id=?`, tenant).Scan(&n)
	return n, e
}
