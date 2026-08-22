package audit

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"time"
)

type Retention struct {
	Repos *repository.Repositories
	Keep  time.Duration
}

func (r Retention) Purge(ctx context.Context, tenant string, now time.Time) (int64, error) {
	if r.Keep <= 0 {
		return 0, fmt.Errorf("retention must be positive")
	}
	cut := now.Add(-r.Keep).UTC().Format(time.RFC3339Nano)
	res, e := r.Repos.DB.SQL.ExecContext(ctx, `DELETE FROM audit_events WHERE tenant_id=? AND occurred_at<?`, tenant, cut)
	if e != nil {
		return 0, e
	}
	return res.RowsAffected()
}
func (r Retention) Window(now time.Time) time.Time { return now.Add(-r.Keep) }
