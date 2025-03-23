package badger

import (
	"context"
	"fmt"

	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

type loggerWrap struct {
}

func (l *loggerWrap) Errorf(p0 string, p1 ...interface{}) {
	msg := fmt.Sprintf(p0, p1...)
	logutil.GetLogger(context.Background()).Error("log from badger db", zap.String("msg", msg))
}

func (l *loggerWrap) Warningf(p0 string, p1 ...interface{}) {
	msg := fmt.Sprintf(p0, p1...)
	logutil.GetLogger(context.Background()).Warn("log from badger db", zap.String("msg", msg))
}

func (l *loggerWrap) Infof(p0 string, p1 ...interface{}) {
	msg := fmt.Sprintf(p0, p1...)
	logutil.GetLogger(context.Background()).Info("log from badger db", zap.String("msg", msg))
}

func (l *loggerWrap) Debugf(p0 string, p1 ...interface{}) {
	msg := fmt.Sprintf(p0, p1...)
	logutil.GetLogger(context.Background()).Debug("log from badger db", zap.String("msg", msg))
}
