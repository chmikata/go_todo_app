package store

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	"github.com/chmikata/go_todo_app/config"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

func New(ctx context.Context, cfg *config.Config) (*sqlx.DB, func(), error) {
	db, err := sql.Open("postgres",
		fmt.Sprintf(
			"host=%s port=%s dbname=%s user=%s password=%s",
			cfg.DBHost, cfg.DBPort,
			cfg.DBName, cfg.DBUser,
			cfg.DBPassword,
		),
	)
	if err != nil {
		return nil, nil, err
	}

	ctx, cancel := context.WithTimeout(ctx, 2*time.Second)
	defer cancel()
	if err := db.PingContext(ctx); err != nil {
		return nil, func() { _ = db.Close() }, err
	}

	xdb := sqlx.NewDb(db, "postgres")
	return xdb, func() { _ = db.Close() }, nil
}

type Beginner interface {
	BeginTx(ctx context.Context, opts *sql.TxOptions) (*sql.Tx, error)
}

type Preparer interface {
	PreparexContext(ctx context.Context, query string) (*sqlx.Stmt, error)
}

type Execer interface {
	ExecContext(ctx context.Context, query string, args ...any) (sql.Result, error)
	NamedExecContext(ctx context.Context, query string, arg interface{})
}

type Queryer interface {
	QueryContext(ctx context.Context, query string, args ...any) (*sqlx.Rows, error)
}
