package cgi

import (
	"context"
	"log"
	"strings"

	"github.com/xxxsen/common/trace"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

const (
	keyTracePrefix  = "TID-"
	keyRequestId    = "x-request-id"
	keyServerAttach = "server_attach"
)

func PanicRecoverMiddleware(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := recover(); err != nil {
			log.Printf("svr panic, path:%s, err:%v", ctx.Request.URL.Path, err)
		}
		ctx.Next()
	}
}

func EnableServerTraceMiddleware(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestid := ctx.GetHeader(keyRequestId)
		if !strings.HasPrefix(requestid, keyTracePrefix) {
			requestid = keyTracePrefix + uuid.NewString()
		}
		ctx.Writer.Header().Set(keyRequestId, requestid)
		ctx.Request = ctx.Request.WithContext(trace.WithTraceId(ctx.Request.Context(), requestid))
		trace.SetTraceId(ctx, requestid)
	}
}

func EnableAttachMiddleware(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(keyServerAttach, svr.c.attach)
	}
}

func GetAttach(ctx context.Context) map[string]interface{} {
	iVal := ctx.Value(keyServerAttach)
	if iVal == nil {
		return nil
	}
	return iVal.(map[string]interface{})
}

func GetAttachKey(ctx context.Context, key string) (interface{}, bool) {
	m := GetAttach(ctx)
	if v, ok := m[key]; ok {
		return v, true
	}
	return nil, false
}
