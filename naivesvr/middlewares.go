package naivesvr

import (
	"log"

	"github.com/xxxsen/common/naivesvr/constants"
	"github.com/xxxsen/common/trace"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func PanicRecoverMiddleware(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		if err := recover(); err != nil {
			log.Printf("svr panic, path:%s, err:%v", ctx.Request.URL.Path, err)
		}
		ctx.Next()
	}
}

func SupportServerGetterMiddleware(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constants.KeyServer, svr)
	}
}

func SupportAttachMiddleware(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Set(constants.KeyServerAttach, svr.c.attach)
	}
}

func EnableServerTraceMiddleware(svr *Server) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		requestid := ctx.GetHeader("x-request-id")
		if len(requestid) == 0 {
			requestid = uuid.NewString()
		}
		ctx.Writer.Header().Set("x-request-id", requestid)
		trace.SetTraceID(ctx, requestid)
	}
}
