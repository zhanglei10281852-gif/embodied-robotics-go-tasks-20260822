package repository

import (
	"context"
	"database/sql"
)

type Formation struct{ ID, TenantID, Name, State string }

func (r *Repositories) EnsureFormationTables(ctx context.Context) error {
	_, e := r.DB.SQL.ExecContext(ctx, `CREATE TABLE IF NOT EXISTS formations(id TEXT PRIMARY KEY,tenant_id TEXT NOT NULL REFERENCES tenants(id),name TEXT NOT NULL,state TEXT NOT NULL);CREATE TABLE IF NOT EXISTS formation_members(formation_id TEXT NOT NULL REFERENCES formations(id),robot_id TEXT NOT NULL REFERENCES robots(id),ordinal INTEGER NOT NULL,PRIMARY KEY(formation_id,robot_id))`)
	return e
}
func (r *Repositories) CreateFormation(ctx context.Context, f Formation, robots []string) error {
	return r.DB.WithTx(context.Background(), func(tx *sql.Tx) error {
		if _, e := tx.ExecContext(ctx, `INSERT INTO formations(id,tenant_id,name,state) VALUES(?,?,?,?)`, f.ID, f.TenantID, f.Name, f.State); e != nil {
			return e
		}
		for i, id := range robots {
			if _, e := tx.ExecContext(ctx, `INSERT INTO formation_members(formation_id,robot_id,ordinal) VALUES(?,?,?)`, f.ID, id, i); e != nil {
				return e
			}
		}
		return nil
	})
}
