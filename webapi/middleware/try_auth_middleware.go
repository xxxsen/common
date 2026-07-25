package middleware

import (
	"github.com/xxxsen/common/webapi/auth"
	"github.com/xxxsen/common/webapi/proxyutil"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func TryAuthMiddleware(fn auth.UserQueryFunc, authenticators ...auth.IAuth) gin.HandlerFunc {
	if len(authenticators) == 0 {
		authenticators = []auth.IAuth{auth.NewBasic()}
	}
	isolated := append([]auth.IAuth(nil), authenticators...)
	return tryAuthMiddleware(fn, isolated...)
}

func tryAuthMiddleware(matchfn auth.UserQueryFunc, ats ...auth.IAuth) gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx := c.Request.Context()
		logger := logutil.GetLogger(ctx).With(zap.String("method", c.Request.Method),
			zap.String("path", c.Request.URL.Path), zap.String("ip", c.ClientIP()))

		for _, authenticator := range ats {
			if authenticator == nil {
				continue
			}
			ak, err := authenticator.Auth(c, matchfn)
			if err != nil {
				continue
			}
			logger.Debug(
				"user auth succ",
				zap.String("auth", authenticator.Name()),
				zap.String("ak", ak),
			)
			authContext := proxyutil.SetUserInfo(c.Request.Context(), &proxyutil.UserInfo{
				AuthType: authenticator.Name(),
				Username: ak,
			})
			c.Request = c.Request.WithContext(authContext)
			return
		}
	}
}
