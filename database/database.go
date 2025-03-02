package database

import (
	"context"
	"database/sql"
	"fmt"
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

func buildSqlDataSource(c *DBConfig) string {
	return fmt.Sprintf("%s:%s@%s(%s:%d)/%s?charset=utf8mb4", c.User, c.Pwd, "tcp", c.Host, c.Port, c.DB)
}

func InitDatabase(c *DBConfig) (*sql.DB, error) {
	client, err := sql.Open("mysql", buildSqlDataSource(c))
	if err != nil {
		return nil, fmt.Errorf("open db failed, err:%w", err)
	}
	if err := client.Ping(); err != nil {
		return nil, fmt.Errorf("ping failed, err:%w", err)
	}
	return client, nil
}
