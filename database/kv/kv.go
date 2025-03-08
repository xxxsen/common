package kv

import (
	"context"
)

type IKvDataBase interface {
	Get(ctx context.Context, table string, key string) ([]byte, bool, error)
	MultiGet(ctx context.Context, table string, keys []string) (map[string][]byte, error)
	Set(ctx context.Context, table string, key string, value []byte) error
	MultiSet(ctx context.Context, table string, kvs map[string][]byte) error
	Close() error
	Iter(ctx context.Context, table string, prefix string, cb IterFunc) error
}

type IterFunc func(ctx context.Context, key string, value []byte) (bool, error)
