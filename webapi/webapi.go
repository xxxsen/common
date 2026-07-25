package webapi

import (
	"context"
	"errors"
	"fmt"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/webapi/middleware"
)

var errUserQueryMissing = errors.New("no user query function provided")

var (
	setGinReleaseMode     sync.Once
	defaultMiddlewareList = []gin.HandlerFunc{
		middleware.PanicRecoverMiddleware(),
		middleware.TraceMiddleware(),
		middleware.LogRequestMiddleware(),
		middleware.NonLengthIOLimitMiddleware(),
	}
)

func nonUserQuery(context.Context, string) (string, bool, error) {
	return "", false, errUserQueryMissing
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
	if err := w.engine.Run(w.bind); err != nil {
		return fmt.Errorf("run web engine: %w", err)
	}
	return nil
}

func (w *webEngine) ServeHTTP(wr http.ResponseWriter, r *http.Request) {
	w.engine.ServeHTTP(wr, r)
}

func NewEngine(root string, addr string, opts ...Option) (IWebEngine, error) {
	setGinReleaseMode.Do(func() {
		gin.SetMode(gin.ReleaseMode)
	})
	c := &config{
		userQuery: nonUserQuery,
		regFn:     nonRegister,
	}
	for _, opt := range opts {
		opt(c)
	}
	middlewares := make([]gin.HandlerFunc, 0, len(defaultMiddlewareList)+1+len(c.mds))
	middlewares = append(middlewares, defaultMiddlewareList...)
	middlewares = append(
		middlewares,
		middleware.TryAuthMiddleware(c.userQuery, c.authenticators...),
	)
	middlewares = append(middlewares, c.mds...)

	engine := gin.New()
	engine.Use(middlewares...)
	gp := &engine.RouterGroup
	if len(root) > 0 && root != "/" {
		gp = gp.Group(root)
	}
	c.regFn(gp)
	if c.noRoute != nil {
		engine.NoRoute(gin.HandlerFunc(c.noRoute))
	}
	return &webEngine{
		bind:   addr,
		engine: engine,
	}, nil
}
