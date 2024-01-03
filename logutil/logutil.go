package logutil

import (
	"context"

	"github.com/xxxsen/common/logger"
	"github.com/xxxsen/common/trace"
	"go.uber.org/zap"
)

func GetLogger(ctx context.Context) *zap.Logger {
	l := logger.Logger()
	traceid, exist := trace.GetTraceId(ctx)
	if !exist || len(traceid) == 0 {
		return l
	}
	return l.With(zap.String("traceid", traceid))
}
