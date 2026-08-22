package repository

import (
	"context"
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"time"

	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/storage"
)

type Repositories struct {
	DB    *storage.DB
	Clock storage.Clock
}

func New(db *storage.DB) *Repositories { return &Repositories{DB: db, Clock: storage.RealClock{}} }

func scanRobot(row interface{ Scan(...any) error }) (domain.Robot, error) {
	var r domain.Robot
	var created string
	var lease sql.NullString
	err := row.Scan(&r.ID, &r.TenantID, &r.Serial, &r.Name, &r.Status, &r.Version, &lease, &created)
	if err != nil {
		return r, err
	}
	r.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	if lease.Valid {
		t, _ := time.Parse(time.RFC3339Nano, lease.String)
		r.LeaseUntil = &t
	}
	return r, nil
}

func scanMission(row interface{ Scan(...any) error }) (domain.Mission, error) {
	var m domain.Mission
	var created, updated string
	err := row.Scan(&m.ID, &m.TenantID, &m.RobotID, &m.Status, &m.Priority, &m.PolicyVersion, &m.Owner, &m.Attempt, &created, &updated, &m.Version)
	if err != nil {
		return m, err
	}
	m.CreatedAt, _ = time.Parse(time.RFC3339Nano, created)
	m.UpdatedAt, _ = time.Parse(time.RFC3339Nano, updated)
	return m, nil
}

func marshal(v any) (string, error) { b, err := json.Marshal(v); return string(b), err }
func isNotFound(err error) bool     { return errors.Is(err, sql.ErrNoRows) }
func wrap(op string, err error) error {
	if err == nil {
		return nil
	}
	return fmt.Errorf("%s: %w", op, err)
}
func now(c storage.Clock) string { return c.Now().UTC().Format(time.RFC3339Nano) }

type RobotRepository struct{ Parent *Repositories }

func (r *Repositories) Robots() RobotRepository { return RobotRepository{Parent: r} }

func (r RobotRepository) Create(ctx context.Context, robot domain.Robot) error {
	if err := robot.Validate(); err != nil {
		return err
	}
	_, err := r.Parent.DB.SQL.ExecContext(ctx, `INSERT INTO robots(id,tenant_id,serial,name,status,version,created_at) VALUES(?,?,?,?,?,?,?)`, robot.ID, robot.TenantID, robot.Serial, robot.Name, robot.Status, robot.Version, now(r.Parent.Clock))
	return wrap("create robot", err)
}
func (r RobotRepository) Get(ctx context.Context, tenant, id string) (domain.Robot, error) {
	row := r.Parent.DB.SQL.QueryRowContext(ctx, `SELECT id,tenant_id,serial,name,status,version,lease_until,created_at FROM robots WHERE tenant_id=? AND id=?`, tenant, id)
	v, err := scanRobot(row)
	if isNotFound(err) {
		return v, domain.ErrNotFound
	}
	return v, wrap("get robot", err)
}
func (r RobotRepository) BySerial(ctx context.Context, tenant, serial string) (domain.Robot, error) {
	row := r.Parent.DB.SQL.QueryRowContext(ctx, `SELECT id,tenant_id,serial,name,status,version,lease_until,created_at FROM robots WHERE tenant_id=? AND serial=?`, tenant, serial)
	v, err := scanRobot(row)
	if isNotFound(err) {
		return v, domain.ErrNotFound
	}
	return v, wrap("get robot by serial", err)
}
func (r RobotRepository) UpdateStatus(ctx context.Context, tenant, id, status string, version int64) error {
	res, err := r.Parent.DB.SQL.ExecContext(ctx, `UPDATE robots SET status=?,version=version+1 WHERE tenant_id=? AND id=? AND version=?`, status, tenant, id, version)
	if err != nil {
		return wrap("update status", err)
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r RobotRepository) SaveLease(ctx context.Context, tenant, id string, until *time.Time, version int64) error {
	var v any
	if until != nil {
		v = until.UTC().Format(time.RFC3339Nano)
	}
	res, err := r.Parent.DB.SQL.ExecContext(ctx, `UPDATE robots SET lease_until=?,version=version+1 WHERE tenant_id=? AND id=? AND version=?`, v, tenant, id, version)
	if err != nil {
		return wrap("save lease", err)
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r RobotRepository) List(ctx context.Context, tenant string, limit int) ([]domain.Robot, error) {
	rows, err := r.Parent.DB.SQL.QueryContext(ctx, `SELECT id,tenant_id,serial,name,status,version,lease_until,created_at FROM robots WHERE tenant_id=? ORDER BY created_at,id LIMIT ?`, tenant, limit)
	if err != nil {
		return nil, wrap("list robots", err)
	}
	defer rows.Close()
	out := make([]domain.Robot, 0, limit)
	for rows.Next() {
		v, e := scanRobot(rows)
		if e != nil {
			return nil, e
		}
		out = append(out, v)
	}
	return out, rows.Err()
}

type MissionRepository struct{ Parent *Repositories }

func (r *Repositories) Missions() MissionRepository { return MissionRepository{Parent: r} }
func (r MissionRepository) Create(ctx context.Context, m domain.Mission, steps []domain.MissionStep) error {
	return r.Parent.DB.WithTx(ctx, func(tx *sql.Tx) error { return r.createTx(ctx, tx, m, steps) })
}
func (r MissionRepository) createTx(ctx context.Context, tx *sql.Tx, m domain.Mission, steps []domain.MissionStep) error {
	n := now(r.Parent.Clock)
	if _, err := tx.ExecContext(ctx, `INSERT INTO missions(id,tenant_id,robot_id,status,priority,policy_version,owner,attempt,created_at,updated_at,version) VALUES(?,?,?,?,?,?,?,?,?,?,?)`, m.ID, m.TenantID, m.RobotID, m.Status, m.Priority, m.PolicyVersion, m.Owner, m.Attempt, n, n, m.Version); err != nil {
		return wrap("insert mission", err)
	}
	for _, s := range steps {
		if _, err := tx.ExecContext(ctx, `INSERT INTO mission_steps(id,mission_id,ordinal,action,payload_json,status) VALUES(?,?,?,?,?,?)`, s.ID, m.ID, s.Ordinal, s.Action, s.PayloadJSON, s.Status); err != nil {
			return wrap("insert mission step", err)
		}
	}
	return nil
}
func (r MissionRepository) Get(ctx context.Context, tenant, id string) (domain.Mission, error) {
	row := r.Parent.DB.SQL.QueryRowContext(ctx, `SELECT id,tenant_id,robot_id,status,priority,policy_version,owner,attempt,created_at,updated_at,version FROM missions WHERE tenant_id=? AND id=?`, tenant, id)
	m, e := scanMission(row)
	if isNotFound(e) {
		return m, domain.ErrNotFound
	}
	return m, wrap("get mission", e)
}
func (r MissionRepository) Transition(ctx context.Context, tenant, id, from, to string, version int64) error {
	res, e := r.Parent.DB.SQL.ExecContext(ctx, `UPDATE missions SET status=?,version=version+1,updated_at=? WHERE tenant_id=? AND id=? AND status=? AND version=?`, to, now(r.Parent.Clock), tenant, id, from, version)
	if e != nil {
		return wrap("transition mission", e)
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}
func (r MissionRepository) ClaimNext(ctx context.Context, tenant, owner string) (domain.Mission, error) {
	var out domain.Mission
	err := r.Parent.DB.WithTx(ctx, func(tx *sql.Tx) error {
		row := tx.QueryRowContext(ctx, `SELECT id,tenant_id,robot_id,status,priority,policy_version,owner,attempt,created_at,updated_at,version FROM missions WHERE tenant_id=? AND status=? ORDER BY priority DESC,created_at,id LIMIT 1`, tenant, domain.MissionQueued)
		m, e := scanMission(row)
		if e != nil {
			return e
		}
		res, e := tx.ExecContext(ctx, `UPDATE missions SET status=?,owner=?,version=version+1,updated_at=? WHERE id=? AND status=? AND version=?`, domain.MissionRunning, owner, now(r.Parent.Clock), m.ID, domain.MissionQueued, m.Version)
		if e != nil {
			return e
		}
		n, _ := res.RowsAffected()
		if n != 1 {
			return domain.ErrConflict
		}
		m.Status = domain.MissionRunning
		m.Owner = owner
		m.Version++
		out = m
		return nil
	})
	if isNotFound(err) {
		return out, domain.ErrNotFound
	}
	return out, err
}
func (r MissionRepository) IncrementAttempt(ctx context.Context, tenant, id string, version int64) error {
	res, e := r.Parent.DB.SQL.ExecContext(ctx, `UPDATE missions SET attempt=attempt+1,version=version+1,updated_at=? WHERE tenant_id=? AND id=? AND version=?`, now(r.Parent.Clock), tenant, id, version)
	if e != nil {
		return e
	}
	n, _ := res.RowsAffected()
	if n != 1 {
		return domain.ErrConflict
	}
	return nil
}

type CapabilityRepository struct{ Parent *Repositories }

func (r *Repositories) Capabilities() CapabilityRepository { return CapabilityRepository{Parent: r} }
func (r CapabilityRepository) Upsert(ctx context.Context, c domain.Capability) error {
	_, e := r.Parent.DB.SQL.ExecContext(ctx, `INSERT INTO robot_capabilities(robot_id,name,revision,schema_json,updated_at) VALUES(?,?,?,?,?) ON CONFLICT(robot_id,name) DO UPDATE SET revision=excluded.revision,schema_json=excluded.schema_json,updated_at=excluded.updated_at`, c.RobotID, c.Name, c.Revision, c.SchemaJSON, now(r.Parent.Clock))
	return wrap("upsert capability", e)
}
func (r CapabilityRepository) Get(ctx context.Context, robot, name string) (domain.Capability, error) {
	var c domain.Capability
	var ts string
	e := r.Parent.DB.SQL.QueryRowContext(ctx, `SELECT robot_id,name,revision,schema_json,updated_at FROM robot_capabilities WHERE robot_id=? AND name=?`, robot, name).Scan(&c.RobotID, &c.Name, &c.Revision, &c.SchemaJSON, &ts)
	if isNotFound(e) {
		return c, domain.ErrNotFound
	}
	c.UpdatedAt, _ = time.Parse(time.RFC3339Nano, ts)
	return c, wrap("get capability", e)
}

func (r *Repositories) CreateTenant(ctx context.Context, t domain.Tenant) error {
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO tenants(id,name,quota,created_at) VALUES(?,?,?,?)`, t.ID, t.Name, t.Quota, now(r.Clock))
	return wrap("create tenant", e)
}
func (r *Repositories) CreateUser(ctx context.Context, u domain.User) error {
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO users(id,tenant_id,email,role) VALUES(?,?,?,?)`, u.ID, u.TenantID, u.Email, u.Role)
	return wrap("create user", e)
}
func (r *Repositories) CreateSession(ctx context.Context, s domain.Session) error {
	_, e := r.DB.SQL.ExecContext(ctx, `INSERT INTO sessions(token,user_id,tenant_id,expires_at) VALUES(?,?,?,?)`, s.Token, s.UserID, s.TenantID, s.ExpiresAt.UTC().Format(time.RFC3339Nano))
	return wrap("create session", e)
}
