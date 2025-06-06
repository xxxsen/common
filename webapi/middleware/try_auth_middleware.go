package middleware

import (
	"github.com/xxxsen/common/webapi/auth"
	"github.com/xxxsen/common/webapi/proxyutil"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func TryAuthMiddleware(fn auth.UserQueryFunc) gin.HandlerFunc {
	return tryAuthMiddleware(fn, auth.AuthList()...)
}

func tryAuthMiddleware(matchfn auth.UserQueryFunc, ats ...auth.IAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logutil.GetLogger(ctx).With(zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path), zap.String("ip", c.ClientIP()))

		for _, fn := range ats {
			ak, err := fn.Auth(c, matchfn)
			if err != nil {
				continue
			}
			logger.Debug("user auth succ", zap.String("auth", fn.Name()), zap.String("ak", ak))
			ctx := c.Request.Context()
			ctx = proxyutil.SetUserInfo(ctx, &proxyutil.UserInfo{
				AuthType: fn.Name(),
				Username: ak,
			})
			c.Request = c.Request.WithContext(ctx)
			return
		}
	}
}
