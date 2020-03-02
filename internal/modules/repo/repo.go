// Package repo is an implements app.Repo.
package repo

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"github.com/ZergsLaw/tg-translate/internal/core"

	_ "github.com/lib/pq"
)

// Repo is an implements app.Repo.
// Responsible for working with database.
type Repo struct {
	db *sql.DB
}

// Close closes database connections.
func (repo *Repo) Close() error {
	return repo.db.Close()
}

// New creates and returns new app.Repo.
func New(conn *sql.DB) core.Repo {
	return &Repo{db: conn}
}

const (
	dbMaxOpenConns  = 30 // about ⅓ of server's max_connections
	dbParallelConns = 5  // a bit more than average
)

// Connect to database by options.
func Connect(ctx context.Context, options ...Option) (*sql.DB, error) {
	opt := defaultConfig()

	for i := range options {
		options[i](opt)
	}

	dbConn, err := sql.Open("postgres", opt.FormatDSN())
	if err != nil {
		return nil, fmt.Errorf("sqlx open: %w", err)
	}
	dbConn.SetMaxOpenConns(dbMaxOpenConns)
	dbConn.SetMaxIdleConns(dbParallelConns)

	err = dbConn.PingContext(ctx)
	for err != nil {
		nextErr := dbConn.PingContext(ctx)
		if errors.Is(nextErr, context.DeadlineExceeded) || errors.Is(nextErr, context.Canceled) {
			if errClose := dbConn.Close(); errClose != nil {
				return nil, fmt.Errorf("db ping: %w, db close: %s", err, errClose)
			}
			return nil, fmt.Errorf("db ping: %w", err)
		}
		err = nextErr
	}

	return dbConn, nil
}
