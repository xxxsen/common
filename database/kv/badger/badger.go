package badger

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/dgraph-io/badger/v4"
	"github.com/xxxsen/common/database/kv"
)

type badgerDB struct {
	db *badger.DB
}

func (b *badgerDB) OnTranscation(ctx context.Context, cb kv.KvTranscationFunc) error {
	tx := b.db.NewTransaction(true)
	defer tx.Discard()
	qe := &badgerQueryExecutor{
		db: b,
		tx: tx,
	}
	if err := cb(ctx, qe); err != nil {
		return fmt.Errorf("run transaction callback: %w", err)
	}
	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit badger transaction: %w", err)
	}
	return nil
}

func (b *badgerDB) Close() error {
	if err := b.db.Close(); err != nil {
		return fmt.Errorf("close badger database: %w", err)
	}
	return nil
}

func (b *badgerDB) Get(ctx context.Context, table string, key string) ([]byte, bool, error) {
	res, err := b.MultiGet(ctx, table, []string{key})
	if err != nil {
		return nil, false, err
	}
	if v, ok := res[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (b *badgerDB) Set(ctx context.Context, table string, key string, value []byte, ttl time.Duration) error {
	return b.MultiSet(ctx, table, map[string][]byte{key: value}, ttl)
}

func (b *badgerDB) Del(ctx context.Context, table string, key string) error {
	return b.MultiDel(ctx, table, []string{key})
}

func (b *badgerDB) copyBytes(dst *[]byte, src []byte) {
	*dst = make([]byte, len(src))
	copy(*dst, src)
}

func (b *badgerDB) multiGet(
	ctx context.Context,
	txn *badger.Txn,
	table string,
	keys []string,
	rs map[string][]byte,
) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("get context: %w", err)
	}
	for _, k := range keys {
		key := b.buildKey(table, k)
		item, err := txn.Get([]byte(key))
		if err != nil {
			if errors.Is(err, badger.ErrKeyNotFound) {
				continue
			}
			return fmt.Errorf("get badger item: %w", err)
		}
		err = item.Value(func(val []byte) error {
			var res []byte
			b.copyBytes(&res, val)
			rs[k] = res
			return nil
		})
		if err != nil {
			return fmt.Errorf("read badger item value: %w", err)
		}
	}
	return nil
}

func (b *badgerDB) MultiGet(ctx context.Context, table string, keys []string) (map[string][]byte, error) {
	rs := make(map[string][]byte)
	err := b.db.View(func(txn *badger.Txn) error {
		return b.multiGet(ctx, txn, table, keys, rs)
	})
	if err != nil {
		return nil, fmt.Errorf("view badger values: %w", err)
	}
	return rs, nil
}

func (b *badgerDB) buildKey(table string, k string) string {
	return fmt.Sprintf("bk:%s:%s", table, k)
}

func (b *badgerDB) multiSet(
	ctx context.Context,
	txn *badger.Txn,
	table string,
	kvs map[string][]byte,
	ttl time.Duration,
) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("set context: %w", err)
	}
	for k, v := range kvs {
		key := b.buildKey(table, k)
		entry := badger.NewEntry([]byte(key), v)
		if ttl > 0 {
			entry.WithTTL(ttl)
		}
		if err := txn.SetEntry(entry); err != nil {
			return fmt.Errorf("set badger item: %w", err)
		}
	}
	return nil
}

func (b *badgerDB) MultiSet(ctx context.Context, table string, kvs map[string][]byte, ttl time.Duration) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		return b.multiSet(ctx, txn, table, kvs, ttl)
	})
	if err != nil {
		return fmt.Errorf("update badger values: %w", err)
	}
	return nil
}

func (b *badgerDB) iter(
	ctx context.Context,
	txn *badger.Txn,
	table string,
	prefix string,
	cb kv.IterFunc,
) error {
	prefixKey := b.buildKey(table, prefix)
	it := txn.NewIterator(badger.DefaultIteratorOptions)
	defer it.Close()

	for it.Seek([]byte(prefixKey)); it.ValidForPrefix([]byte(prefixKey)); it.Next() {
		item := it.Item()
		key := string(item.Key())
		var res []byte
		if err := item.Value(func(val []byte) error {
			b.copyBytes(&res, val)
			return nil
		}); err != nil {
			return fmt.Errorf("read badger iterator value: %w", err)
		}
		next, err := cb(ctx, key, res)
		if err != nil {
			return fmt.Errorf("run badger iterator callback: %w", err)
		}
		if !next {
			break
		}
	}
	return nil
}

func (b *badgerDB) Iter(ctx context.Context, table string, prefix string, cb kv.IterFunc) error {
	err := b.db.View(func(txn *badger.Txn) error {
		return b.iter(ctx, txn, table, prefix, cb)
	})
	if err != nil {
		return fmt.Errorf("iterate badger values: %w", err)
	}
	return nil
}

func (b *badgerDB) multiDel(txn *badger.Txn, table string, keys []string) error {
	for _, k := range keys {
		key := b.buildKey(table, k)
		err := txn.Delete([]byte(key))
		if err != nil {
			return fmt.Errorf("delete badger item: %w", err)
		}
	}
	return nil
}

func (b *badgerDB) MultiDel(ctx context.Context, table string, keys []string) error {
	err := b.db.Update(func(txn *badger.Txn) error {
		if err := ctx.Err(); err != nil {
			return fmt.Errorf("delete context: %w", err)
		}
		return b.multiDel(txn, table, keys)
	})
	if err != nil {
		return fmt.Errorf("delete badger values: %w", err)
	}
	return nil
}

func New(f string) (kv.IKvDataBase, error) {
	db, err := badger.Open(badger.DefaultOptions(f).WithLoggingLevel(badger.ERROR))
	if err != nil {
		return nil, fmt.Errorf("open badger database: %w", err)
	}
	return &badgerDB{db: db}, nil
}

type badgerQueryExecutor struct {
	db *badgerDB
	tx *badger.Txn
}

func (b *badgerQueryExecutor) Get(ctx context.Context, table string, key string) ([]byte, bool, error) {
	res, err := b.MultiGet(ctx, table, []string{key})
	if err != nil {
		return nil, false, err
	}
	if v, ok := res[key]; ok {
		return v, true, nil
	}
	return nil, false, nil
}

func (b *badgerQueryExecutor) MultiGet(
	ctx context.Context,
	table string,
	keys []string,
) (map[string][]byte, error) {
	rs := make(map[string][]byte, len(keys))
	if err := b.db.multiGet(ctx, b.tx, table, keys, rs); err != nil {
		return nil, fmt.Errorf("get transaction values: %w", err)
	}
	return rs, nil
}

func (b *badgerQueryExecutor) Iter(ctx context.Context, table string, prefix string, cb kv.IterFunc) error {
	return b.db.iter(ctx, b.tx, table, prefix, cb)
}

func (b *badgerQueryExecutor) Set(
	ctx context.Context,
	table string,
	key string,
	value []byte,
	ttl time.Duration,
) error {
	return b.MultiSet(ctx, table, map[string][]byte{key: value}, ttl)
}

func (b *badgerQueryExecutor) MultiSet(
	ctx context.Context,
	table string,
	kvs map[string][]byte,
	ttl time.Duration,
) error {
	return b.db.multiSet(ctx, b.tx, table, kvs, ttl)
}

func (b *badgerQueryExecutor) Del(ctx context.Context, table string, key string) error {
	return b.MultiDel(ctx, table, []string{key})
}

func (b *badgerQueryExecutor) MultiDel(ctx context.Context, table string, keys []string) error {
	if err := ctx.Err(); err != nil {
		return fmt.Errorf("delete context: %w", err)
	}
	return b.db.multiDel(b.tx, table, keys)
}
