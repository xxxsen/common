package trace

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/naivesvr/constants"
)

type traceIdType struct{}

var (
	keyTraceId = traceIdType{}
)

func SetTraceId(ctx *gin.Context, traceid string) {
	ctx.Set(constants.KeyTraceID, traceid)
}

func WithTraceId(ctx context.Context, traceid string) context.Context {
	return context.WithValue(ctx, keyTraceId, traceid)
}

func GetTraceId(ctx context.Context) (string, bool) {
	v := ctx.Value(keyTraceId)
	if v == nil {
		return "", false
	}
	return v.(string), true
}

func MustGetTraceId(ctx context.Context) string {
	if v, ok := GetTraceId(ctx); ok {
		return v
	}
	panic("not found trace id")
}
