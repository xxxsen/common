package database

import (
	"context"
	"database/sql"
	"io"
)

type OnTxFunc func(ctx context.Context, qe IQueryExecer) error

type IQueryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
}

type IExecer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type ITxer interface {
	OnTransation(ctx context.Context, cb OnTxFunc) error
}

type IQueryExecer interface {
	IQueryer
	IExecer
}

type IDatabase interface {
	IQueryExecer
	ITxer
	io.Closer
}
