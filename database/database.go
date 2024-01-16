package database

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/xxxsen/common/errs"
)

type IQueryer interface {
	QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error)
	QueryRowContext(ctx context.Context, query string, args ...interface{}) *sql.Row
}

type IExecer interface {
	ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error)
}

type IQueryExecer interface {
	IQueryer
	IExecer
}

func buildSqlDataSource(c *DBConfig) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4", c.User, c.Pwd, "tcp", c.Host, c.Port, c.DB)
}

func InitDatabase(c *DBConfig) (*sql.DB, error) {
	client, err := sql.Open("mysql", buildSqlDataSource(c))
	if err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "open db fail", err)
	}
	if err := client.Ping(); err != nil {
		return nil, errs.Wrap(errs.ErrDatabase, "ping fail", err)
	}
	return client, nil
}
