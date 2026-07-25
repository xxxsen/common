package sqlite

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	// Register the pure-Go SQLite driver with database/sql.
	_ "github.com/glebarez/go-sqlite"
	"github.com/xxxsen/common/database"
)

type sqliteDBWrap struct {
	db *sql.DB
}

type OnDBCreateSuccFunc func(db database.IDatabase) error

func (s *sqliteDBWrap) OnTransation(ctx context.Context, cb database.OnTxFunc) error {
	tx, err := s.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if err := cb(ctx, tx); err != nil {
		return fmt.Errorf("run transaction callback: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}
	return nil
}

func (s *sqliteDBWrap) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	rows, err := s.db.QueryContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("query sqlite database: %w", err)
	}
	return rows, nil
}

func (s *sqliteDBWrap) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	result, err := s.db.ExecContext(ctx, query, args...)
	if err != nil {
		return nil, fmt.Errorf("execute sqlite statement: %w", err)
	}
	return result, nil
}

func (s *sqliteDBWrap) Close() error {
	if err := s.db.Close(); err != nil {
		return fmt.Errorf("close sqlite database: %w", err)
	}
	return nil
}

func New(f string, fns ...OnDBCreateSuccFunc) (database.IDatabase, error) {
	db, err := sql.Open("sqlite", f)
	if err != nil {
		return nil, fmt.Errorf("open sqlite database: %w", err)
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)
	w := &sqliteDBWrap{
		db: db,
	}
	for _, fn := range fns {
		if err := fn(w); err != nil {
			_ = db.Close()
			return nil, fmt.Errorf("initialize sqlite database: %w", err)
		}
	}
	return w, nil
}
