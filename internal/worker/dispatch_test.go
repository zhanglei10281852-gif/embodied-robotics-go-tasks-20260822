package worker

import (
	"context"
	"github.com/zhanglei10281852-gif/embodied-robotics-go-tasks-20260822/internal/domain"
	"testing"
	"time"
)

func TestDispatcherSubmit(t *testing.T) {
	d := NewDispatcher(1)
	job := DispatchJob{ID: "j", TenantID: "t", Mission: domain.Mission{ID: "m", TenantID: "t"}}
	if e := ValidateDispatch(job); e != nil {
		t.Fatal(e)
	}
	done := d.Submit(context.Background(), job, func(context.Context, domain.Mission) error { return nil })
	select {
	case e := <-done:
		if e != nil {
			t.Fatal(e)
		}
	case <-time.After(time.Second):
		t.Fatal("dispatch timeout")
	}
	got, ok := d.Get("j")
	if !ok || got.State != DispatchStopped {
		t.Fatal(got, ok)
	}
}
