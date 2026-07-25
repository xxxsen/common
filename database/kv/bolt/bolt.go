package bolt

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/xxxsen/common/database/kv"
	"go.etcd.io/bbolt"
)

var (
	errTableNotFound = errors.New("table not found")
)

type boltDB struct {
	db *bbolt.DB
}

func New(f string, tables ...string) (kv.IKvDataBase, error) {
	db, err := bbolt.Open(f, 0644, bbolt.DefaultOptions)
	if err != nil {
		return nil, fmt.Errorf("open bolt database: %w", err)
	}

	if err := db.Batch(func(tx *bbolt.Tx) error {
		for _, tab := range tables {
			if _, err := tx.CreateBucketIfNotExists([]byte(tab)); err != nil {
				return fmt.Errorf("create bolt bucket: %w", err)
			}
		}
		return nil
	}); err != nil {
		_ = db.Close()
		return nil, fmt.Errorf("initialize bolt database: %w", err)
	}
	return &boltDB{db: db}, nil
}

func (db *boltDB) copyBytes(dst *[]byte, src []byte) {
	*dst = make([]byte, len(src))
	copy(*dst, src)
}

func (db *boltDB) Close() error {
	if err := db.db.Close(); err != nil {
		return fmt.Errorf("close bolt database: %w", err)
	}
	return nil
}

func (db *boltDB) Get(ctx context.Context, table string, key string) ([]byte, bool, error) {
	m, err := db.MultiGet(ctx, table, []string{key})
	if err != nil {
		return nil, false, err
	}
	if v, ok := m[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (db *boltDB) Set(ctx context.Context, table string, key string, value []byte, ttl time.Duration) error {
	return db.MultiSet(ctx, table, map[string][]byte{key: value}, ttl)
}

func (db *boltDB) iter(ctx context.Context, tx *bbolt.Tx, table string, prefix string, cb kv.IterFunc) error {
	bk := tx.Bucket([]byte(table))
	if bk == nil {
		return errTableNotFound
	}
	cursor := bk.Cursor()
	for k, v := cursor.Seek([]byte(prefix)); k != nil; k, v = cursor.Next() {
		if v == nil {
			continue
		}
		next, err := cb(ctx, string(k), v)
		if err != nil {
			return fmt.Errorf("run bolt iterator callback: %w", err)
		}
		if !next {
			break
		}
	}
	return nil
}

func (db *boltDB) Iter(ctx context.Context, table string, prefix string, cb kv.IterFunc) error {
	err := db.db.View(func(tx *bbolt.Tx) error {
		return db.iter(ctx, tx, table, prefix, cb)
	})
	if err != nil {
		return fmt.Errorf("iterate bolt values: %w", err)
	}
	return nil
}

func (db *boltDB) multiGet(
	ctx context.Context,
	tx *bbolt.Tx,
	table string,
	keys []string,
	values *map[string][]byte,
) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("get context: %w", err)
	}
	bk := tx.Bucket([]byte(table))
	if bk == nil {
		return errTableNotFound
	}
	for _, key := range keys {
		val := bk.Get([]byte(key))
		if len(val) == 0 {
			continue
		}
		var newv []byte
		db.copyBytes(&newv, val)
		(*values)[key] = newv
	}
	return nil
}

func (db *boltDB) MultiGet(
	ctx context.Context,
	table string,
	keys []string,
) (map[string][]byte, error) {
	m := make(map[string][]byte, len(keys))
	if err := db.db.View(func(tx *bbolt.Tx) error {
		return db.multiGet(ctx, tx, table, keys, &m)
	}); err != nil {
		return nil, fmt.Errorf("get bolt values: %w", err)
	}
	return m, nil
}

func (db *boltDB) multiSet(
	ctx context.Context,
	tx *bbolt.Tx,
	table string,
	kvs map[string][]byte,
	ttl time.Duration,
) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("set context: %w", err)
	}
	bk := tx.Bucket([]byte(table))
	if bk == nil {
		return errTableNotFound
	}
	// Bolt does not support TTL; retain the parameter for interface parity.
	_ = ttl
	for k, v := range kvs {
		if err := bk.Put([]byte(k), v); err != nil {
			return fmt.Errorf("set bolt value: %w", err)
		}
	}
	return nil
}

func (db *boltDB) MultiSet(
	ctx context.Context,
	table string,
	kvs map[string][]byte,
	ttl time.Duration,
) error {
	err := db.db.Update(func(tx *bbolt.Tx) error {
		return db.multiSet(ctx, tx, table, kvs, ttl)
	})
	if err != nil {
		return fmt.Errorf("set bolt values: %w", err)
	}
	return nil
}

func (db *boltDB) Del(ctx context.Context, table string, key string) error {
	return db.MultiDel(ctx, table, []string{key})
}

func (db *boltDB) multiDel(ctx context.Context, tx *bbolt.Tx, table string, keys []string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("delete context: %w", err)
	}
	bk := tx.Bucket([]byte(table))
	if bk == nil {
		return errTableNotFound
	}
	for _, key := range keys {
		if err := bk.Delete([]byte(key)); err != nil {
			return fmt.Errorf("delete bolt value: %w", err)
		}
	}
	return nil
}

func (db *boltDB) MultiDel(ctx context.Context, table string, keys []string) error {
	err := db.db.Update(func(tx *bbolt.Tx) error {
		return db.multiDel(ctx, tx, table, keys)
	})
	if err != nil {
		return fmt.Errorf("delete bolt values: %w", err)
	}
	return nil
}

func (db *boltDB) OnTranscation(ctx context.Context, cb kv.KvTranscationFunc) error {
	tx, err := db.db.Begin(true)
	if err != nil {
		return fmt.Errorf("begin bolt transaction: %w", err)
	}
	defer func() {
		_ = tx.Rollback()
	}()
	if err := cb(ctx, &boltQueryExecutor{db: db, tx: tx}); err != nil {
		return fmt.Errorf("run bolt transaction callback: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit bolt transaction: %w", err)
	}
	return nil
}

type boltQueryExecutor struct {
	db *boltDB
	tx *bbolt.Tx
}

func (b *boltQueryExecutor) Get(ctx context.Context, table string, key string) ([]byte, bool, error) {
	ms, err := b.MultiGet(ctx, table, []string{key})
	if err != nil {
		return nil, false, err
	}
	if v, ok := ms[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (b *boltQueryExecutor) MultiGet(
	ctx context.Context,
	table string,
	keys []string,
) (map[string][]byte, error) {
	rs := make(map[string][]byte)
	if err := b.db.multiGet(ctx, b.tx, table, keys, &rs); err != nil {
		return nil, fmt.Errorf("get transaction values: %w", err)
	}
	return rs, nil
}

func (b *boltQueryExecutor) Iter(ctx context.Context, table string, prefix string, cb kv.IterFunc) error {
	return b.db.iter(ctx, b.tx, table, prefix, cb)
}

func (b *boltQueryExecutor) Set(
	ctx context.Context,
	table string,
	key string,
	value []byte,
	ttl time.Duration,
) error {
	return b.MultiSet(ctx, table, map[string][]byte{key: value}, ttl)
}

func (b *boltQueryExecutor) MultiSet(
	ctx context.Context,
	table string,
	kvs map[string][]byte,
	ttl time.Duration,
) error {
	return b.db.multiSet(ctx, b.tx, table, kvs, ttl)
}

func (b *boltQueryExecutor) Del(ctx context.Context, table string, key string) error {
	return b.MultiDel(ctx, table, []string{key})
}

func (b *boltQueryExecutor) MultiDel(ctx context.Context, table string, keys []string) error {
	return b.db.multiDel(ctx, b.tx, table, keys)
}
