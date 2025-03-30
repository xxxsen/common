package sqlite

import (
	"context"
	"database/sql"
	"time"

	_ "github.com/glebarez/go-sqlite"
	"github.com/xxxsen/common/database"
)

type sqliteDBWrap struct {
	db *sql.DB
}

type OnDBCreateSuccFunc func(db database.IDatabase) error

func (s *sqliteDBWrap) OnTransation(ctx context.Context, cb database.OnTxFunc) error {
	tx, err := s.db.Begin()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := cb(ctx, tx); err != nil {
		return err
	}
	if err := tx.Commit(); err != nil {
		return err
	}
	return nil
}

func (s *sqliteDBWrap) QueryContext(ctx context.Context, query string, args ...interface{}) (*sql.Rows, error) {
	return s.db.QueryContext(ctx, query, args...)
}

func (s *sqliteDBWrap) ExecContext(ctx context.Context, query string, args ...interface{}) (sql.Result, error) {
	return s.db.ExecContext(ctx, query, args...)
}

func (s *sqliteDBWrap) Close() error {
	return s.db.Close()
}

func New(f string, fns ...OnDBCreateSuccFunc) (database.IDatabase, error) {
	db, err := sql.Open("sqlite", f)
	if err != nil {
		return nil, err
	}
	db.SetMaxOpenConns(1)
	db.SetMaxIdleConns(1)
	db.SetConnMaxLifetime(time.Hour)
	w := &sqliteDBWrap{
		db: db,
	}
	for _, fn := range fns {
		if err := fn(w); err != nil {
			return nil, err
		}
	}
	return w, nil
}
