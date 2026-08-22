package service

import (
	"context"
	"testing"
	"time"
)

func TestMaintenanceReport(t *testing.T) {
	s := newService(t)
	report, e := s.RunMaintenance(context.Background(), time.Now())
	if e != nil || !report.ForeignKeysOK {
		t.Fatalf("%v %+v", e, report)
	}
	if _, e = s.HealthReport(context.Background(), "tenant", time.Minute); e != nil {
		t.Fatal(e)
	}
}
