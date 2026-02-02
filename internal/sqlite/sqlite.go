package sqlite

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/garrettladley/peli/internal/db"
	migrations "github.com/garrettladley/peli/internal/migrations/sqlite"
	_ "modernc.org/sqlite"
)

const driverName = "sqlite"

type DB interface {
	db.Querier
	Close() error
}

type database struct {
	conn *sql.DB
	*db.Queries
}

var _ DB = (*database)(nil)

func New(ctx context.Context, path string) (DB, error) {
	conn, err := sql.Open(driverName, path)
	if err != nil {
		return nil, fmt.Errorf("open: %w", err)
	}

	if err := migrations.Apply(ctx, conn); err != nil {
		_ = conn.Close()
		return nil, fmt.Errorf("migrations: %w", err)
	}

	return &database{
		conn:    conn,
		Queries: db.New(conn),
	}, nil
}

func (d *database) Close() error {
	if err := d.conn.Close(); err != nil {
		return fmt.Errorf("close: %w", err)
	}
	return nil
}
