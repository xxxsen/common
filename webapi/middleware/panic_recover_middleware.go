package middleware

import (
	"errors"
	"net/http"
	"runtime/debug"

	"github.com/xxxsen/common/webapi/proxyutil"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

var errServicePanic = errors.New("service panic")

func PanicRecoverMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if err := recover(); err != nil {
				logutil.GetLogger(ctx).Error(
					"service panic",
					zap.Any("panic", err),
					zap.String("path", ctx.Request.URL.Path),
					zap.String("stack", string(debug.Stack())),
				)
				proxyutil.FailJson(ctx, http.StatusInternalServerError, errServicePanic)
			}
		}()
		ctx.Next()
	}
}
