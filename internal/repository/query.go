package repository

import (
	"context"
	"database/sql"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"strings"
	"time"
)

type MissionFilter struct {
	Status      string
	RobotID     string
	MinPriority *int
	Before      time.Time
	Limit       int
}

func (r *Repositories) SearchMissions(ctx context.Context, tenant string, f MissionFilter) (domain.Page[domain.Mission], error) {
	limit := f.Limit
	if limit < 1 || limit > 500 {
		limit = 100
	}
	args := []any{tenant}
	where := []string{"tenant_id=?"}
	if f.Status != "" {
		where = append(where, "status=?")
		args = append(args, f.Status)
	}
	if f.RobotID != "" {
		where = append(where, "robot_id=?")
		args = append(args, f.RobotID)
	}
	if f.MinPriority != nil {
		where = append(where, "priority>=?")
		args = append(args, *f.MinPriority)
	}
	if !f.Before.IsZero() {
		where = append(where, "updated_at<?")
		args = append(args, f.Before.UTC().Format(time.RFC3339Nano))
	}
	args = append(args, limit)
	q := fmt.Sprintf(`SELECT id,tenant_id,robot_id,status,priority,policy_version,owner,attempt,created_at,updated_at,version FROM missions WHERE %s ORDER BY updated_at DESC,id DESC LIMIT ?`, strings.Join(where, " AND "))
	rows, e := r.DB.SQL.QueryContext(ctx, q, args...)
	if e != nil {
		return domain.Page[domain.Mission]{}, wrap("search missions", e)
	}
	defer rows.Close()
	out := domain.Page[domain.Mission]{Items: []domain.Mission{}}
	for rows.Next() {
		m, e := scanMission(rows)
		if e != nil {
			return out, e
		}
		out.Items = append(out.Items, m)
	}
	return out, rows.Err()
}
func (r *Repositories) ListSteps(ctx context.Context, tenant, mission string) ([]domain.MissionStep, error) {
	rows, e := r.DB.SQL.QueryContext(ctx, `SELECT s.id,s.mission_id,s.ordinal,s.action,s.payload_json,s.status FROM mission_steps s JOIN missions m ON m.id=s.mission_id WHERE m.tenant_id=? AND s.mission_id=? ORDER BY s.ordinal`, tenant, mission)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []domain.MissionStep{}
	for rows.Next() {
		var s domain.MissionStep
		if e := rows.Scan(&s.ID, &s.MissionID, &s.Ordinal, &s.Action, &s.PayloadJSON, &s.Status); e != nil {
			return nil, e
		}
		out = append(out, s)
	}
	return out, rows.Err()
}
func (r *Repositories) UpdateStep(ctx context.Context, tenant, step, status string) error {
	res, e := r.DB.SQL.ExecContext(ctx, `UPDATE mission_steps SET status=? WHERE id=? AND mission_id IN(SELECT id FROM missions WHERE tenant_id=?)`, status, step, tenant)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrNotFound
	}
	return nil
}
func (r *Repositories) FindSession(ctx context.Context, token string) (domain.Session, error) {
	var s domain.Session
	var expires, rev sql.NullString
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT token,user_id,tenant_id,expires_at,revoked_at FROM sessions WHERE token=?`, token).Scan(&s.Token, &s.UserID, &s.TenantID, &expires, &rev)
	if e == sql.ErrNoRows {
		return s, domain.ErrNotFound
	}
	s.ExpiresAt, _ = time.Parse(time.RFC3339Nano, expires.String)
	if rev.Valid {
		t, _ := time.Parse(time.RFC3339Nano, rev.String)
		s.RevokedAt = &t
	}
	return s, e
}
func (r *Repositories) RevokeSession(ctx context.Context, token string) error {
	res, e := r.DB.SQL.ExecContext(ctx, `UPDATE sessions SET revoked_at=? WHERE token=?`, r.Clock.Now().UTC().Format(time.RFC3339Nano), token)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrNotFound
	}
	return nil
}
func (r *Repositories) CountOpenAlerts(ctx context.Context, tenant string) (int, error) {
	var n int
	e := r.DB.SQL.QueryRowContext(ctx, `SELECT count(*) FROM alerts WHERE tenant_id=? AND state=?`, tenant, domain.AlertOpen).Scan(&n)
	return n, e
}
