package cgi

import (
	"net/http"
	"reflect"
	"time"

	"github.com/xxxsen/common/cgi/codec"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

func WrapHandler(ptr interface{}, cc codec.ICodec, pfunc ProcessFunc) gin.HandlerFunc {
	creator := func() interface{} {
		if ptr == nil {
			return nil
		}
		typ := reflect.TypeOf(ptr)
		val := reflect.New(typ.Elem())
		return val.Interface()
	}
	return wrapHandler(NewHandler(creator, cc, pfunc))
}

func wrapHandler(h IHandler) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		logger := logutil.GetLogger(ctx).With(
			zap.String("method", ctx.Request.Method),
			zap.String("path", ctx.Request.URL.Path),
			zap.String("ua", ctx.Request.UserAgent()),
		)
		start := time.Now()
		err := handleRequest(h, ctx)
		cost := time.Since(start)
		logger = logger.With(zap.Duration("cost", cost))
		if err != nil {
			logger.Error("process request failed", zap.Error(err))
			return
		}
		logger.Debug("process request succ")
	}
}

func handleRequest(h IHandler, ctx *gin.Context) error {
	req := h.Request()
	codec := h.Codec()
	var err error
	if req != nil {
		err = codec.Decode(ctx, req)
	}
	if !errs.IsErrOK(err) {
		writeJson(ctx, http.StatusBadRequest, err)
		return err
	}
	statusCode, rsp, err := h.Process(ctx, req)
	if !errs.IsErrOK(err) {
		writeJson(ctx, statusCode, err)
		return err
	}
	err = codec.Encode(ctx, statusCode, err, rsp)
	if !errs.IsErrOK(err) {
		return err
	}
	return nil
}

func writeJson(ctx *gin.Context, status int, err error) {
	m := make(map[string]interface{})
	code := 0
	msg := ""
	if err != nil {
		ierr := errs.FromError(err)
		code = int(ierr.Code())
		msg = ierr.Message()
	}
	m["code"] = code
	m["message"] = msg
	ctx.AbortWithStatusJSON(status, m)
}
