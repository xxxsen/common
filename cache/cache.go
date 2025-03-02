package cache

import (
	"context"
	"time"
)

// Deprecated: should not use this
type ICache interface {
	Set(ctx context.Context, key interface{}, val interface{}, timeout time.Duration) error
	Get(ctx context.Context, key interface{}) (interface{}, bool, error)
	Del(ctx context.Context, key interface{}) error
}
