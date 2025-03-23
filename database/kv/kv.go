package kv

import (
	"context"
	"io"
	"time"
)

type IKvQueryer interface {
	Get(ctx context.Context, table string, key string) ([]byte, bool, error)
	MultiGet(ctx context.Context, table string, keys []string) (map[string][]byte, error)
	Iter(ctx context.Context, table string, prefix string, cb IterFunc) error
}

type IKvExecutor interface {
	Set(ctx context.Context, table string, key string, value []byte, ttl time.Duration) error
	MultiSet(ctx context.Context, table string, kvs map[string][]byte, ttl time.Duration) error
	Del(ctx context.Context, table string, key string) error
	MultiDel(ctx context.Context, table string, keys []string) error
}

type IkvTx interface {
	OnTranscation(ctx context.Context, cb KvTranscationFunc) error
}

type IKvQueryExecutor interface {
	IKvQueryer
	IKvExecutor
}

type IKvDataBase interface {
	IKvQueryExecutor
	IkvTx
	io.Closer
}

type IterFunc func(ctx context.Context, key string, value []byte) (bool, error)

type KvTranscationFunc func(ctx context.Context, db IKvQueryExecutor) error
