package service

import (
	"context"
	"testing"
)

func TestRegisterFleet(t *testing.T) {
	s := newService(t)
	ctx := context.Background()
	fleet, e := s.RegisterFleet(ctx, "tenant", []BulkRegistration{{Serial: "F1", Name: "one"}, {Serial: "F2", Name: "two"}})
	if e != nil || len(fleet) != 2 {
		t.Fatalf("%v %d", e, len(fleet))
	}
	if e = s.ValidateFleet(ctx, "tenant", []BulkRegistration{{Serial: "A"}, {Serial: "A"}}); e == nil {
		t.Fatal("duplicates")
	}
}
