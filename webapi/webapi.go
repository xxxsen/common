package webapi

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/webapi/middleware"
)

var (
	defaultMiddlewareList = []gin.HandlerFunc{
		middleware.PanicRecoverMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LogRequestMiddleware(),
		middleware.NonLengthIOLimitMiddleware(),
	}
)

func init() {
	gin.SetMode(gin.ReleaseMode)
}

func nonUserQuery(ctx context.Context, ak string) (string, bool, error) {
	return "", false, fmt.Errorf("no user query fn provide")
}

func nonRegister(c *gin.RouterGroup) {

}

func RunEngine(root string, addr string, opts ...Option) error {
	c := &config{
		userQuery: nonUserQuery,
		regFn:     nonRegister,
	}
	for _, opt := range opts {
		opt(c)
	}
	middlewares := make([]gin.HandlerFunc, 0, len(defaultMiddlewareList)+1+len(c.mds))
	middlewares = append(middlewares, defaultMiddlewareList...)
	middlewares = append(middlewares, middleware.TryAuthMiddleware(c.userQuery))
	middlewares = append(middlewares, c.mds...)

	engine := gin.New()
	engine.Use(middlewares...)
	gp := &engine.RouterGroup
	if len(root) > 0 {
		gp = gp.Group(root)
	}
	c.regFn(gp)
	return engine.Run(addr)
}
