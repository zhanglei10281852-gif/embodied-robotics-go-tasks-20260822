package worker

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
)

func TestPipelineStages(t *testing.T) {
	p := NewPipeline([]PipelineStage{StageValidate, StageReserve, StageExecute, StagePersist, StageNotify})
	if e := ValidatePipeline(p); e != nil {
		t.Fatal(e)
	}
	seen := 0
	e := p.Execute(context.Background(), func(context.Context, PipelineStage) error { seen++; return nil })
	if e != nil || seen != 5 {
		t.Fatalf("%v %d", e, seen)
	}
	if len(p.Events()) != 5 {
		t.Fatal("events")
	}
}
func TestPipelineCancelAndBatch(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	p := NewPipeline([]PipelineStage{StageValidate})
	if e := p.Execute(ctx, func(context.Context, PipelineStage) error { return nil }); e == nil {
		t.Fatal("cancel")
	}
	r := ProcessBatch(context.Background(), []domain.TelemetryEvent{{ID: "1"}, {ID: "2"}}, func(context.Context, domain.TelemetryEvent) error { return errors.New("bad") })
	if r.Rejected != 2 || r.Accepted != 0 {
		t.Fatal(r)
	}
}
