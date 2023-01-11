package iotool

import (
	"context"
	"io"

	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func makeBuf(sz int64) []byte {
	if sz == 0 {
		return nil
	}
	return make([]byte, sz)
}

func ProxyStreamWithBuffer(ctx context.Context, left, right io.ReadWriteCloser, bufsz int64) error {
	defer func() {
		el := left.Close()
		er := right.Close()
		if el != nil || er != nil {
			logutil.GetLogger(ctx).Error("proxy stream may fail", zap.Any("left_err", el), zap.Any("right_err", er))
		}
	}()
	ch := make(chan error, 2)
	go func() {
		_, err := io.CopyBuffer(left, right, makeBuf(bufsz))
		if err != nil {
			ch <- errs.Wrap(errs.ErrIO, "copy right to left fail", err)
			return
		}
		ch <- nil
	}()
	go func() {
		_, err := io.CopyBuffer(right, left, makeBuf(bufsz))
		if err != nil {
			ch <- errs.Wrap(errs.ErrIO, "copy left to right fail", err)
			return
		}
		ch <- nil
	}()
	err := <-ch
	logutil.GetLogger(ctx).With(zap.Error(err)).Debug("proxy thread exit")
	return err
}

func ProxyStream(ctx context.Context, left, right io.ReadWriteCloser) error {
	return ProxyStreamWithBuffer(ctx, left, right, 0)
}
