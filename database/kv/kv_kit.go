package kv

import (
	"context"
	"encoding/json"
	"fmt"
	"time"
)

type IterObjectFunc[T interface{}] func(ctx context.Context, key string, val *T) (bool, error)
type OnGetKeyForUpdateFunc[T interface{}] func(ctx context.Context, key string, val *T) (*T, bool, error)

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
	return SetJsonObjectWithTTL(ctx, db, table, key, val, 0)
}

func SetJsonObjectWithTTL(ctx context.Context, db IKvQueryExecutor, table string, key string, val interface{}, ttl time.Duration) error {
	raw, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err := db.Set(ctx, table, key, raw, ttl); err != nil {
		return err
	}
	return nil
}

func MultiSetJsonObject[T interface{}](ctx context.Context, db IKvQueryExecutor, table string, m map[string]*T) error {
	return MultiSetJsonObjectWithTTL[T](ctx, db, table, m, 0)
}

func MultiSetJsonObjectWithTTL[T interface{}](ctx context.Context, db IKvQueryExecutor, table string, m map[string]*T, ttl time.Duration) error {
	ms := make(map[string][]byte, len(m))
	for k, v := range m {
		raw, err := json.Marshal(v)
		if err != nil {
			return err
		}
		ms[k] = raw
	}
	return db.MultiSet(ctx, table, ms, ttl)
}

func IterJsonObject[T interface{}](ctx context.Context, db IKvQueryExecutor, table string, prefix string, cb IterObjectFunc[T]) error {
	return db.Iter(ctx, table, prefix, func(ctx context.Context, key string, value []byte) (bool, error) {
		v := new(T)
		if err := json.Unmarshal(value, v); err != nil {
			return false, err
		}
		return cb(ctx, key, v)
	})
}

func OnGetJsonKeyForUpdate[T interface{}](ctx context.Context, db IKvDataBase, table string, key string, ttl time.Duration, cb OnGetKeyForUpdateFunc[T]) error {
	return db.OnTranscation(ctx, func(ctx context.Context, db IKvQueryExecutor) error {
		res, ok, err := db.Get(ctx, table, key)
		if err != nil {
			return err
		}
		if !ok {
			return fmt.Errorf("not found")
		}
		v := new(T)
		if err := json.Unmarshal(res, v); err != nil {
			return err
		}
		newV, ok, err := cb(ctx, key, v)
		if err != nil {
			return err
		}
		if !ok {
			return nil
		}
		res, err = json.Marshal(newV)
		if err != nil {
			return err
		}
		return db.Set(ctx, table, key, res, ttl)
	})
}
