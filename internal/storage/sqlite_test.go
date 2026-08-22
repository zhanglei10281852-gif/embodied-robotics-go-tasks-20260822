package storage

import (
	"context"
	"database/sql"
	"testing"
)

func TestOpenMigrateAndRollback(t *testing.T) {
	db, err := Open(context.Background(), "file:storage-test?mode=memory&cache=shared")
	if err != nil {
		t.Fatal(err)
	}
	defer db.Close()
	if err := db.Ping(context.Background()); err != nil {
		t.Fatal(err)
	}
	if err := db.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO tenants(id,name,created_at) VALUES('t','Fleet',datetime('now'))`)
		return err
	}); err != nil {
		t.Fatal(err)
	}
	if err := db.WithTx(context.Background(), func(tx *sql.Tx) error {
		_, err := tx.Exec(`INSERT INTO tenants(id,name,created_at) VALUES('t','Duplicate',datetime('now'))`)
		return err
	}); err == nil {
		t.Fatal("expected unique failure")
	}
}
