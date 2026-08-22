package worker

import (
	"context"
	"fmt"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/repository"
)

type Ingestor struct{ Repos *repository.Repositories }

func (i *Ingestor) Batch(ctx context.Context, events []domain.TelemetryEvent) error {
	if len(events) == 0 {
		return nil
	}
	for _, e := range events {
		if err := e.Validate(); err != nil {
			return fmt.Errorf("telemetry batch: %w", err)
		}
		if err := i.Repos.AppendTelemetry(ctx, e); err != nil {
			return err
		}
	}
	return nil
}
