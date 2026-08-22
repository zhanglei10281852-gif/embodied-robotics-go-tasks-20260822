package repository

import (
	"context"
	"database/sql"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"time"
)

func (r *Repositories) SavePolicy(ctx context.Context, p domain.Policy) error {
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO policies(id,tenant_id,name,revision,expression,state,updated_at) VALUES(?,?,?,?,?,?,?)`, p.ID, p.TenantID, p.Name, p.Revision, p.Expression, p.State, r.Clock.Now().UTC().Format(time.RFC3339Nano))
	return wrap("save policy", e)
}
func (r *Repositories) LoadPolicy(ctx context.Context, tenant, name string, revision int) (domain.Policy, error) {
	var p domain.Policy
	var ts string
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT id,tenant_id,name,revision,expression,state,updated_at FROM policies WHERE tenant_id=? AND name=? AND revision=?`, tenant, name, revision).Scan(&p.ID, &p.TenantID, &p.Name, &p.Revision, &p.Expression, &p.State, &ts)
	if e == sql.ErrNoRows {
		return p, domain.ErrNotFound
	}
	p.UpdatedAt, _ = time.Parse(time.RFC3339Nano, ts)
	return p, wrap("load policy", e)
}
func (r *Repositories) SaveApproval(ctx context.Context, a domain.Approval) error {
	var d any
	if a.DecidedAt != nil {
		d = a.DecidedAt.UTC().Format(time.RFC3339Nano)
	}
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO approvals(id,tenant_id,mission_id,policy_id,policy_revision,decision,decided_by,decided_at) VALUES(?,?,?,?,?,?,?,?)`, a.ID, a.TenantID, a.MissionID, a.PolicyID, a.PolicyRevision, a.Decision, a.DecidedBy, d)
	return wrap("save approval", e)
}
