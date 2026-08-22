package service

import (
	"context"
	"errors"
	"testing"
)

func TestChecklistApprovalHonorsCancellation(t *testing.T) {
	c := Checklist{MissionID: "mission", Items: []ChecklistItem{{ID: "operator", Required: true, Passed: true}}}
	ctx, cancel := context.WithCancel(context.Background())
	cancel()
	_, err := (&Service{}).ApproveChecklist(ctx, c)
	if !errors.Is(err, context.Canceled) {
		t.Fatalf("cancelled approval returned %v", err)
	}
}
