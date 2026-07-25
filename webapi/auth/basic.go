package auth

import (
	"crypto/hmac"
	"errors"

	"github.com/gin-gonic/gin"
)

const BasicAuthName = "basic"

var (
	errBasicCredentialsMissing = errors.New("basic credentials are missing")
	errBasicUserUnknown        = errors.New("basic user is unknown")
	errBasicSecretMismatch     = errors.New("basic secret does not match")
)

type basicAuth struct {
}

// NewBasic returns an isolated Basic authenticator.
func NewBasic() IAuth {
	return &basicAuth{}
}

func (b *basicAuth) Name() string {
	return BasicAuthName
}

func (b *basicAuth) Auth(ctx *gin.Context, fn UserQueryFunc) (string, error) {
	uak, usk, ok := ctx.Request.BasicAuth()
	if !ok {
		return "", errBasicCredentialsMissing
	}

	sk, ok, err := fn(ctx, uak)
	if err != nil {
		return "", err
	}
	if !ok {
		return "", errBasicUserUnknown
	}
	if !hmac.Equal([]byte(sk), []byte(usk)) {
		return "", errBasicSecretMismatch
	}
	return uak, nil
}
