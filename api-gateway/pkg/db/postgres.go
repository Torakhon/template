// Package postgres implements postgres connection.
package db

import (
	"time"

	"github.com/Masterminds/squirrel"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

// Postgres -.
type Postgres struct {
	connAttempts int
	connTimeout  time.Duration

	Builder squirrel.StatementBuilderType
	Db      *sqlx.DB
}
