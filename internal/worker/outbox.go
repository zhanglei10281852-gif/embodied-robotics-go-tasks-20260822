package worker

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/telemetry"
	"log/slog"
)

type Outbox struct {
	Repos  *repository.Repositories
	Client telemetry.Client
	Log    *slog.Logger
}

func (o *Outbox) Drain(ctx context.Context, limit int) error {
	jobs, e := o.Repos.ClaimOutbox(ctx, limit)
	if e != nil {
		return e
	}
	for _, j := range jobs {
		e = o.Client.Send(ctx, j.Topic, []byte(j.PayloadJSON))
		if fe := o.Repos.FinishOutbox(ctx, j.ID, e); fe != nil {
			return fe
		}
		if e != nil && !errors.Is(e, context.Canceled) {
			return e
		}
	}
	return nil
}
func (o *Outbox) Run(ctx context.Context, limit int) <-chan error {
	done := make(chan error, 1)
	go func() { defer close(done); done <- o.Drain(ctx, limit) }()
	return done
}
