package worker

import (
	"context"
	"errors"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"sync"
	"time"
)

type PipelineStage string

const (
	StageValidate PipelineStage = "validate"
	StageReserve  PipelineStage = "reserve"
	StageExecute  PipelineStage = "execute"
	StagePersist  PipelineStage = "persist"
	StageNotify   PipelineStage = "notify"
)

type PipelineEvent struct {
	Stage PipelineStage
	At    time.Time
	Err   error
}
type Pipeline struct {
	mu     sync.Mutex
	stages []PipelineStage
	events []PipelineEvent
}

func NewPipeline(stages []PipelineStage) *Pipeline {
	out := append([]PipelineStage(nil), stages...)
	return &Pipeline{stages: out}
}
func (p *Pipeline) Stages() []PipelineStage {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]PipelineStage(nil), p.stages...)
}
func (p *Pipeline) Execute(ctx context.Context, fn func(context.Context, PipelineStage) error) error {
	for _, stage := range p.Stages() {
		select {
		case <-ctx.Done():
			return ctx.Err()
		default:
		}
		e := fn(ctx, stage)
		p.mu.Lock()
		p.events = append(p.events, PipelineEvent{Stage: stage, At: time.Now().UTC(), Err: e})
		p.mu.Unlock()
		if e != nil {
			return e
		}
	}
	return nil
}
func (p *Pipeline) Events() []PipelineEvent {
	p.mu.Lock()
	defer p.mu.Unlock()
	return append([]PipelineEvent(nil), p.events...)
}
func ValidatePipeline(p *Pipeline) error {
	if p == nil || len(p.stages) == 0 {
		return errors.New("pipeline empty")
	}
	seen := map[PipelineStage]bool{}
	for _, s := range p.stages {
		if seen[s] {
			return errors.New("duplicate pipeline stage")
		}
		seen[s] = true
	}
	return nil
}
func ProcessBatch(ctx context.Context, items []domain.TelemetryEvent, fn func(context.Context, domain.TelemetryEvent) error) domain.BatchResult {
	out := domain.BatchResult{}
	for _, item := range items {
		if e := fn(ctx, item); e != nil {
			out.Add(e)
		} else {
			out.Add(nil)
		}
	}
	return out
}
