package storage

import (
	"context"
	"testing"
)

func TestBackupTables(t *testing.T) {
	d, e := Open(context.Background(), "file:backup-test?mode=memory&cache=shared")
	if e != nil {
		t.Fatal(e)
	}
	defer d.Close()
	b := Backup{DB: d}
	if e = b.ValidateTables(context.Background(), []string{"robots", "missions", "audit_events"}); e != nil {
		t.Fatal(e)
	}
	rows, e := b.ExportRows(context.Background(), "tenants")
	if e != nil || len(rows) == 0 {
		t.Fatalf("rows %v", e)
	}
}
