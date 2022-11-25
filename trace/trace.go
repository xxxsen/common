package trace

import (
	"context"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/naivesvr/constants"
)

func SetTraceId(ctx *gin.Context, traceid string) {
	ctx.Set(constants.KeyTraceID, traceid)
}

func WithTraceID(ctx context.Context, traceid string) context.Context {
	//lint:ignore SA1029 ignore it
	return context.WithValue(ctx, constants.KeyTraceID, traceid)
}

func GetTraceId(ctx context.Context) (string, bool) {
	v := ctx.Value(constants.KeyTraceID)
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
