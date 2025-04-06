package webapi

import (
	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/webapi/auth"
)

type RegisterFunc func(c *gin.RouterGroup)

type config struct {
	mds       []gin.HandlerFunc
	userQuery auth.UserQueryFunc
	regFn     RegisterFunc
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

func WithRegister(fn RegisterFunc) Option {
	return func(c *config) {
		c.regFn = fn
	}
}
