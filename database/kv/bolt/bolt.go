package bolt

import (
	"context"
	"fmt"

	"github.com/xxxsen/common/database/kv"
	"go.etcd.io/bbolt"
)

var (
	errTableNotFound = fmt.Errorf("table not found")
)

type boltDB struct {
	db *bbolt.DB
}

func New(f string, tables ...string) (kv.IKvDataBase, error) {
	db, err := bbolt.Open(f, 0644, bbolt.DefaultOptions)
	if err != nil {
		return nil, err
	}

	if err := db.Batch(func(tx *bbolt.Tx) error {
		for _, tab := range tables {
			if _, err := tx.CreateBucketIfNotExists([]byte(tab)); err != nil {
				return err
			}
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return &boltDB{db: db}, nil
}

func (db *boltDB) copyBytes(dst *[]byte, src []byte) {
	*dst = make([]byte, len(src))
	copy(*dst, src)
}

func (db *boltDB) Close() error {
	return db.db.Close()
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

func (db *boltDB) Set(ctx context.Context, table string, key string, value []byte) error {
	return db.MultiSet(ctx, table, map[string][]byte{key: value})
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
			return err
		}
		if !next {
			break
		}
	}
	return nil
}

func (db *boltDB) Iter(ctx context.Context, table string, prefix string, cb kv.IterFunc) error {
	if err := db.db.View(func(tx *bbolt.Tx) error {
		return db.iter(ctx, tx, table, prefix, cb)
	}); err != nil {
		return err
	}
	return nil
}

func (b *boltDB) multiGet(ctx context.Context, tx *bbolt.Tx, table string, keys []string, m *map[string][]byte) error {
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
		b.copyBytes(&newv, val)
		(*m)[key] = newv
	}
	return nil
}

func (b *boltDB) MultiGet(ctx context.Context, table string, keys []string) (map[string][]byte, error) {
	m := make(map[string][]byte, len(keys))
	if err := b.db.View(func(tx *bbolt.Tx) error {
		return b.multiGet(ctx, tx, table, keys, &m)
	}); err != nil {
		return nil, err
	}
	return m, nil
}

func (b *boltDB) multiSet(ctx context.Context, tx *bbolt.Tx, table string, kvs map[string][]byte) error {
	bk := tx.Bucket([]byte(table))
	if bk == nil {
		return errTableNotFound
	}
	for k, v := range kvs {
		if err := bk.Put([]byte(k), v); err != nil {
			return err
		}
	}
	return nil
}

func (b *boltDB) MultiSet(ctx context.Context, table string, kvs map[string][]byte) error {
	if err := b.db.Update(func(tx *bbolt.Tx) error {
		return b.multiSet(ctx, tx, table, kvs)
	}); err != nil {
		return err
	}
	return nil
}

func (b *boltDB) Del(ctx context.Context, table string, key string) error {
	return b.MultiDel(ctx, table, []string{key})
}

func (b *boltDB) multiDel(ctx context.Context, tx *bbolt.Tx, table string, keys []string) error {
	bk := tx.Bucket([]byte(table))
	if bk == nil {
		return errTableNotFound
	}
	for _, key := range keys {
		if err := bk.Delete([]byte(key)); err != nil {
			return err
		}
	}
	return nil
}

func (b *boltDB) MultiDel(ctx context.Context, table string, keys []string) error {
	if err := b.db.Update(func(tx *bbolt.Tx) error {
		return b.multiDel(ctx, tx, table, keys)
	}); err != nil {
		return err
	}
	return nil
}

func (b *boltDB) OnTranscation(ctx context.Context, cb kv.KvTranscationFunc) error {
	tx, err := b.db.Begin(true)
	if err != nil {
		return err
	}
	defer tx.Rollback()
	if err := cb(ctx, &boltQueryExecutor{db: b, tx: tx}); err != nil {
		return err
	}
	return tx.Commit()
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

func (b *boltQueryExecutor) MultiGet(ctx context.Context, table string, keys []string) (map[string][]byte, error) {
	rs := make(map[string][]byte)
	if err := b.db.multiGet(ctx, b.tx, table, keys, &rs); err != nil {
		return nil, err
	}
	return rs, nil
}

func (b *boltQueryExecutor) Iter(ctx context.Context, table string, prefix string, cb kv.IterFunc) error {
	return b.db.iter(ctx, b.tx, table, prefix, cb)
}

func (b *boltQueryExecutor) Set(ctx context.Context, table string, key string, value []byte) error {
	return b.MultiSet(ctx, table, map[string][]byte{key: value})
}

func (b *boltQueryExecutor) MultiSet(ctx context.Context, table string, kvs map[string][]byte) error {
	return b.db.multiSet(ctx, b.tx, table, kvs)
}

func (b *boltQueryExecutor) Del(ctx context.Context, table string, key string) error {
	return b.MultiDel(ctx, table, []string{key})
}

func (b *boltQueryExecutor) MultiDel(ctx context.Context, table string, keys []string) error {
	return b.db.multiDel(ctx, b.tx, table, keys)
}
