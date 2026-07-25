package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/webapi/auth"
)

type RegisterFunc func(c *gin.RouterGroup)
type NoRouteFunc func(c *gin.Context)

type config struct {
	mds            []gin.HandlerFunc
	userQuery      auth.UserQueryFunc
	regFn          RegisterFunc
	authenticators []auth.IAuth
	noRoute        NoRouteFunc
}

type Option func(c *config)

func WithExtraMiddlewares(mds ...gin.HandlerFunc) Option {
	return func(c *config) {
		c.mds = mds
	}
}

func WithAuth(uq auth.UserQueryFunc) Option {
	return func(c *config) {
		c.userQuery = uq
	}
}

// WithAuthenticators replaces the default Basic authenticator with an
// explicitly ordered, engine-local authenticator list.
func WithAuthenticators(authenticators ...auth.IAuth) Option {
	isolated := append([]auth.IAuth(nil), authenticators...)
	return func(c *config) {
		c.authenticators = isolated
	}
}

func WithRegister(fn RegisterFunc) Option {
	return func(c *config) {
		c.regFn = fn
	}
}

// WithNoRoute installs an engine-local handler for otherwise unmatched routes.
func WithNoRoute(fn NoRouteFunc) Option {
	return func(c *config) {
		c.noRoute = fn
	}
}
