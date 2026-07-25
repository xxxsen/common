package kv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"time"
)

type IterObjectFunc[T interface{}] func(ctx context.Context, key string, val *T) (bool, error)
type OnGetKeyForUpdateFunc[T interface{}] func(ctx context.Context, key string, val *T) (*T, bool, error)

var errKeyNotFound = errors.New("key not found")

func GetJsonObject[T interface{}](
	ctx context.Context,
	db IKvQueryExecutor,
	table string,
	key string,
) (*T, bool, error) {
	raw, ok, err := db.Get(ctx, table, key)
	if err != nil {
		return nil, false, fmt.Errorf("get JSON object: %w", err)
	}
	if !ok {
		return nil, false, nil
	}
	inst := new(T)
	if err := json.Unmarshal(raw, inst); err != nil {
		return nil, false, err
	}
	return inst, true, nil
}

func MultiGetJsonObject[T interface{}](
	ctx context.Context,
	db IKvQueryExecutor,
	table string,
	keys []string,
) (map[string]*T, error) {
	ms, err := db.MultiGet(ctx, table, keys)
	if err != nil {
		return nil, fmt.Errorf("get JSON objects: %w", err)
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

func SetJsonObjectWithTTL(
	ctx context.Context,
	db IKvQueryExecutor,
	table string,
	key string,
	val interface{},
	ttl time.Duration,
) error {
	raw, err := json.Marshal(val)
	if err != nil {
		return err
	}
	if err := db.Set(ctx, table, key, raw, ttl); err != nil {
		return fmt.Errorf("set JSON object: %w", err)
	}
	return nil
}

func MultiSetJsonObject[T interface{}](
	ctx context.Context,
	db IKvQueryExecutor,
	table string,
	objects map[string]*T,
) error {
	return MultiSetJsonObjectWithTTL[T](ctx, db, table, objects, 0)
}

func MultiSetJsonObjectWithTTL[T interface{}](
	ctx context.Context,
	db IKvQueryExecutor,
	table string,
	objects map[string]*T,
	ttl time.Duration,
) error {
	ms := make(map[string][]byte, len(objects))
	for k, v := range objects {
		raw, err := json.Marshal(v)
		if err != nil {
			return err
		}
		ms[k] = raw
	}
	if err := db.MultiSet(ctx, table, ms, ttl); err != nil {
		return fmt.Errorf("set JSON objects: %w", err)
	}
	return nil
}

func IterJsonObject[T interface{}](
	ctx context.Context,
	db IKvQueryExecutor,
	table string,
	prefix string,
	cb IterObjectFunc[T],
) error {
	err := db.Iter(ctx, table, prefix, func(ctx context.Context, key string, value []byte) (bool, error) {
		v := new(T)
		if err := json.Unmarshal(value, v); err != nil {
			return false, err
		}
		return cb(ctx, key, v)
	})
	if err != nil {
		return fmt.Errorf("iterate JSON objects: %w", err)
	}
	return nil
}

func OnGetJsonKeyForUpdate[T interface{}](
	ctx context.Context,
	db IKvDataBase,
	table string,
	key string,
	ttl time.Duration,
	cb OnGetKeyForUpdateFunc[T],
) error {
	err := db.OnTranscation(ctx, func(ctx context.Context, executor IKvQueryExecutor) error {
		res, ok, err := executor.Get(ctx, table, key)
		if err != nil {
			return fmt.Errorf("get update object: %w", err)
		}
		if !ok {
			return errKeyNotFound
		}
		v := new(T)
		if err := json.Unmarshal(res, v); err != nil {
			return err
		}
		newV, ok, err := cb(ctx, key, v)
		if err != nil {
			return fmt.Errorf("run update callback: %w", err)
		}
		if !ok {
			return nil
		}
		res, err = json.Marshal(newV)
		if err != nil {
			return err
		}
		if err := executor.Set(ctx, table, key, res, ttl); err != nil {
			return fmt.Errorf("set update object: %w", err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("update JSON object transaction: %w", err)
	}
	return nil
}
