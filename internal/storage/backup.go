package storage

import (
	"context"
	"database/sql"
	"fmt"
	"strings"
)

type Backup struct{ DB *DB }

func (b Backup) Tables(ctx context.Context) ([]string, error) {
	rows, e := b.DB.SQL.QueryContext(ctx, `SELECT name FROM sqlite_master WHERE type='table' ORDER BY name`)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	out := []string{}
	for rows.Next() {
		var n string
		if e := rows.Scan(&n); e != nil {
			return nil, e
		}
		out = append(out, n)
	}
	return out, rows.Err()
}
func (b Backup) ValidateTables(ctx context.Context, required []string) error {
	got, e := b.Tables(ctx)
	if e != nil {
		return e
	}
	set := map[string]bool{}
	for _, n := range got {
		set[n] = true
	}
	for _, n := range required {
		if !set[n] {
			return fmt.Errorf("missing table %s", n)
		}
	}
	return nil
}
func (b Backup) ExportRows(ctx context.Context, table string) ([][]string, error) {
	if strings.TrimSpace(table) == "" || strings.ContainsAny(table, " ;\"") {
		return nil, fmt.Errorf("invalid table")
	}
	rows, e := b.DB.SQL.QueryContext(ctx, `SELECT * FROM `+table)
	if e != nil {
		return nil, e
	}
	defer rows.Close()
	cols, e := rows.Columns()
	if e != nil {
		return nil, e
	}
	out := [][]string{cols}
	for rows.Next() {
		vals := make([]any, len(cols))
		ptrs := make([]any, len(cols))
		for i := range vals {
			ptrs[i] = &vals[i]
		}
		if e := rows.Scan(ptrs...); e != nil {
			return nil, e
		}
		line := make([]string, len(cols))
		for i, v := range vals {
			line[i] = fmt.Sprint(v)
		}
		out = append(out, line)
	}
	return out, rows.Err()
}
func (b Backup) Close(ctx context.Context) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
		return b.DB.Close()
	}
}

var _ *sql.DB
