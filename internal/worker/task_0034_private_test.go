package worker

import (
	"context"
	"errors"
	"testing"
)

func TestPipelineStopsWhenContextIsCancelledBetweenStages(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	p := NewPipeline([]PipelineStage{StageValidate, StageExecute})
	calls := 0
	err := p.Execute(ctx, func(context.Context, PipelineStage) error {
		calls++
		if calls == 1 {
			cancel()
		}
		return nil
	})
	if !errors.Is(err, context.Canceled) || calls != 1 {
		t.Fatalf("pipeline ignored cancellation: err=%v calls=%d", err, calls)
	}
}
