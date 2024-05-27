package models

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/blockloop/scan"
)

// SingleRowQuery executes a query that returns a single row.
//
// Call it passing ctx, db and columns of the structure.
func SingleRowQuery(ctx context.Context, db *sql.DB, query string, rows ...interface{}) {
	db.QueryRowContext(ctx, query).Scan(rows...)
}

// Query executes a query that returns one or more rows.
//
// It can handle passing ctx, db and output []structure, to be parset at the caller.
func Query(ctx context.Context, db *sql.DB, query string, rows interface{}) error {
	r, err := db.QueryContext(ctx, query)
	if err != nil {
		return err
	}
	defer r.Close()
	err = scan.Rows(rows, r)
	if err != nil {
		return fmt.Errorf("failed to scan rows: %w", err)
	}

	return nil
}

// func (d *Devices) Close() error {
// 	return d.Conn.Close()
// }
