package cache

import (
	"context"
	"time"

	lru "github.com/hnlq715/golang-lru"
)

type LocalCache struct {
	c *lru.Cache
}

func New(size int) (*LocalCache, error) {
	if size <= 0 {
		size = 10000
	}
	c, err := lru.New(size)
	if err != nil {
		return nil, err
	}
	return &LocalCache{
		c: c,
	}, nil
}

func (c *LocalCache) Set(ctx context.Context, key interface{}, val interface{}, timeout time.Duration) error {
	_ = c.c.AddEx(key, val, timeout)
	return nil
}

func (c *LocalCache) Get(ctx context.Context, key interface{}) (interface{}, bool, error) {
	val, ok := c.c.Get(key)
	return val, ok, nil
}

func (c *LocalCache) Del(ctx context.Context, key interface{}) error {
	c.c.Remove(key)
	return nil
}
