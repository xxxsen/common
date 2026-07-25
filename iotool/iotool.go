package iotool

import (
	"context"
	"errors"
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
		buffer := make([]byte, defaultBufSize)
		return &buffer
	},
}

func transfer(w io.Writer, r io.Reader, errs chan<- error) {
	buffer, ok := bufPool.Get().(*[]byte)
	if !ok {
		panic("iotool: invalid buffer pool value")
	}
	defer bufPool.Put(buffer)
	_, err := io.CopyBuffer(w, r, *buffer)
	if errors.Is(err, io.EOF) {
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
