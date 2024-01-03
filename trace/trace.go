package trace

import (
	"context"
)

const (
	keyStrTraceId = "x-trace-id"
)

type traceIdType struct{}
type IKvGetter interface {
	Get(string) (interface{}, bool)
}

type IKvSetter interface {
	Set(k string, v interface{})
}

var (
	keyTraceId = traceIdType{}
)

func SetTraceId(ctx IKvSetter, traceid string) {
	ctx.Set(keyStrTraceId, traceid)
}

func WithTraceId(ctx context.Context, traceid string) context.Context {
	return context.WithValue(ctx, keyTraceId, traceid)
}

func GetTraceId(ctx context.Context) (string, bool) {
	if v := ctx.Value(keyTraceId); v != nil {
		return v.(string), true
	}
	getter, ok := ctx.(IKvGetter)
	if ok {
		if v, exist := getter.Get(keyStrTraceId); exist {
			return v.(string), true
		}
	}
	return "", false
}

func MustGetTraceId(ctx context.Context) string {
	if v, ok := GetTraceId(ctx); ok {
		return v
	}
	panic("not found trace id")
}
