package storage

import (
	"context"
	"database/sql"
	"embed"
	"fmt"
	"strings"

	_ "modernc.org/sqlite"
)

//go:embed migration.sql
var migrationFS embed.FS

type DB struct{ SQL *sql.DB }

func Open(ctx context.Context, dsn string) (*DB, error) {
	db, err := sql.Open("sqlite", dsn)
	if err != nil {
		return nil, fmt.Errorf("open sqlite: %w", err)
	}
	db.SetMaxOpenConns(8)
	db.SetMaxIdleConns(8)
	store := &DB{SQL: db}
	if err := store.Migrate(ctx); err != nil {
		_ = db.Close()
		return nil, err
	}
	return store, nil
}

func (d *DB) Migrate(ctx context.Context) error {
	b, err := migrationFS.ReadFile("migration.sql")
	if err != nil {
		return fmt.Errorf("read migration: %w", err)
	}
	for _, stmt := range strings.Split(string(b), ";") {
		stmt = strings.TrimSpace(stmt)
		if stmt == "" {
			continue
		}
		if _, err := d.SQL.ExecContext(ctx, stmt); err != nil {
			return fmt.Errorf("migration: %w", err)
		}
	}
	return nil
}

func (d *DB) Close() error                   { return d.SQL.Close() }
func (d *DB) Ping(ctx context.Context) error { return d.SQL.PingContext(ctx) }

func (d *DB) WithTx(ctx context.Context, fn func(*sql.Tx) error) error {
	tx, err := d.SQL.BeginTx(ctx, &sql.TxOptions{Isolation: sql.LevelSerializable})
	if err != nil {
		return fmt.Errorf("begin tx: %w", err)
	}
	if err := fn(tx); err != nil {
		_ = tx.Rollback()
		return err
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit tx: %w", err)
	}
	return nil
}
