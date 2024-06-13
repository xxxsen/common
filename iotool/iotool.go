package iotool

import (
	"context"
	"io"
	"sync"

	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

const (
	defaultBufSize = 32 * 1024
)

var bufPool = sync.Pool{
	New: func() interface{} {
		return make([]byte, defaultBufSize)
	},
}

func transfer(w io.Writer, r io.Reader, errs chan<- error) {
	buf := bufPool.Get().([]byte)
	defer bufPool.Put(buf)
	_, err := io.CopyBuffer(w, r, buf)
	if err == io.EOF {
		err = nil
	}
	errs <- err
}

func ProxyStream(ctx context.Context, left, right io.ReadWriter) error {
	ch := make(chan error, 2)
	go transfer(left, right, ch)
	go transfer(right, left, ch)
	err := <-ch
	if err != nil {
		logutil.GetLogger(ctx).With(zap.Error(err)).Error("proxy stream err")
	}
	return err
}
