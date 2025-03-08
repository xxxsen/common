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
	var rs []byte
	var ok bool
	if err := db.db.View(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte(table))
		if bk == nil {
			return errTableNotFound
		}
		res := bk.Get([]byte(key))
		if len(res) == 0 {
			return nil
		}
		ok = true
		db.copyBytes(&rs, res)
		return nil
	}); err != nil {
		return nil, false, err
	}
	return rs, ok, nil

}

func (db *boltDB) Set(ctx context.Context, table string, key string, value []byte) error {
	if err := db.db.Update(func(tx *bbolt.Tx) error {
		bk := tx.Bucket([]byte(table))
		if bk == nil {
			return errTableNotFound
		}
		return bk.Put([]byte(key), value)
	}); err != nil {
		return err
	}
	return nil
}

func (db *boltDB) Iter(ctx context.Context, table string, prefix string, cb kv.IterFunc) error {
	if err := db.db.View(func(tx *bbolt.Tx) error {
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
	}); err != nil {
		return err
	}
	return nil
}

func (b *boltDB) MultiGet(ctx context.Context, table string, keys []string) (map[string][]byte, error) {
	m := make(map[string][]byte, len(keys))
	if err := b.db.View(func(tx *bbolt.Tx) error {
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
			m[key] = newv
		}
		return nil
	}); err != nil {
		return nil, err
	}
	return m, nil
}

func (b *boltDB) MultiSet(ctx context.Context, table string, kvs map[string][]byte) error {
	if err := b.db.Update(func(tx *bbolt.Tx) error {
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
	}); err != nil {
		return err
	}
	return nil
}
