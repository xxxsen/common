package kv

import (
	"context"
	"encoding/json"
)

type IterObjectFunc[T any] func(ctx context.Context, key string, val *T) (bool, error)

func GetJsonObject[T interface{}](ctx context.Context, db IKvQueryExecutor, table string, key string) (*T, bool, error) {
	raw, ok, err := db.Get(ctx, table, key)
	if err != nil || !ok {
		return nil, ok, err
	}
	inst := new(T)
	if err := json.Unmarshal(raw, inst); err != nil {
		return nil, false, err
	}
	return inst, true, nil
}

func MultiGetJsonObject[T interface{}](ctx context.Context, db IKvQueryExecutor, table string, ks []string) (map[string]*T, error) {
	ms, err := db.MultiGet(ctx, table, ks)
	if err != nil {
		return nil, err
	}
	rs := make(map[string]*T, len(ms))
	for k, raw := range ms {
		v := new(T)
		if err := json.Unmarshal(raw, v); err != nil {
			return nil, err
		}
		rs[k] = v
	}
	return rs, nil
}

func SetJsonObject(ctx context.Context, db IKvQueryExecutor, table string, key string, val interface{}) error {
	raw, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err := db.Set(ctx, table, key, raw); err != nil {
		return err
	}
	return nil
}

func MultiSetJsonObject[T interface{}](ctx context.Context, db IKvQueryExecutor, table string, m map[string]*T) error {
	ms := make(map[string][]byte, len(m))
	for k, v := range m {
		raw, err := json.Marshal(v)
		if err != nil {
			return err
		}
		ms[k] = raw
	}
	return db.MultiSet(ctx, table, ms)
}

func IterJsonObject[T any](ctx context.Context, db IKvQueryExecutor, table string, prefix string, cb IterObjectFunc[T]) {
	db.Iter(ctx, table, prefix, func(ctx context.Context, key string, value []byte) (bool, error) {
		v := new(T)
		if err := json.Unmarshal(value, v); err != nil {
			return false, err
		}
		return cb(ctx, key, v)
	})
}
