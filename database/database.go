package database

import (
	"context"
	"database/sql"
)

type IQueryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type IExecer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type IQueryExecer interface {
	IQueryer
	IExecer
}
