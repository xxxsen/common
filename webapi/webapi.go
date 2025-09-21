package webapi

import (
	"context"
	"fmt"
	"net/http"

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

type IWebEngine interface {
	http.Handler
	Run() error
}

type webEngine struct {
	bind   string
	engine *gin.Engine
}

func (w *webEngine) Run() error {
	return w.engine.Run(w.bind)
}

func (w *webEngine) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	w.engine.ServeHTTP(wr, r)
}

func NewEngine(root string, addr string, opts ...Option) (IWebEngine, error) {
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
	if len(root) > 0 && root != "/" {
		gp = gp.Group(root)
	}
	c.regFn(gp)
	return &webEngine{
		bind:   addr,
		engine: engine,
	}, nil
}
