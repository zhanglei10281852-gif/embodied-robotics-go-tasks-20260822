package worker

import (
	"context"
	"errors"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sync"
	"time"
)

type DispatchState string

const (
	DispatchIdle    DispatchState = "idle"
	DispatchRunning DispatchState = "running"
	DispatchStopped DispatchState = "stopped"
)

type DispatchJob struct {
	ID, TenantID      string
	Mission           domain.Mission
	State             DispatchState
	Started, Finished time.Time
	Error             error
}
type Dispatcher struct {
	mu   sync.Mutex
	jobs map[string]DispatchJob
	sem  chan struct{}
}

func NewDispatcher(limit int) *Dispatcher {
	if limit < 1 {
		limit = 1
	}
	return &Dispatcher{jobs: map[string]DispatchJob{}, sem: make(chan struct{}, limit)}
}
func (d *Dispatcher) Submit(ctx context.Context, job DispatchJob, fn func(context.Context, domain.Mission) error) <-chan error {
	done := make(chan error, 1)
	d.mu.Lock()
	job.State = DispatchRunning
	job.Started = time.Now().UTC()
	d.jobs[job.ID] = job
	d.mu.Unlock()
	go func() {
		select {
		case d.sem <- struct{}{}:
		case <-ctx.Done():
			done <- ctx.Err()
			return
		}
		e := fn(ctx, job.Mission)
		d.mu.Lock()
		job.Finished = time.Now().UTC()
		job.State = DispatchStopped
		job.Error = e
		d.jobs[job.ID] = job
		d.mu.Unlock()
		done <- e
	}()
	return done
}
func (d *Dispatcher) Get(id string) (DispatchJob, bool) {
	d.mu.Lock()
	defer d.mu.Unlock()
	j, ok := d.jobs[id]
	return j, ok
}
func (d *Dispatcher) Cancel(id string) error {
	d.mu.Lock()
	defer d.mu.Unlock()
	j, ok := d.jobs[id]
	if !ok {
		return errors.New("dispatch job missing")
	}
	if j.State != DispatchRunning {
		return errors.New("dispatch job not running")
	}
	j.State = DispatchStopped
	j.Error = context.Canceled
	d.jobs[id] = j
	return nil
}
func (d *Dispatcher) Snapshot() []DispatchJob {
	d.mu.Lock()
	defer d.mu.Unlock()
	out := make([]DispatchJob, 0, len(d.jobs))
	for _, j := range d.jobs {
		out = append(out, j)
	}
	return out
}
func ValidateDispatch(job DispatchJob) error {
	if job.ID == "" || job.TenantID == "" || job.Mission.ID == "" {
		return fmt.Errorf("dispatch identity missing")
	}
	if job.Mission.TenantID != job.TenantID {
		return domain.ErrForbidden
	}
	return nil
}
