package auth

import (
	"context"

	"github.com/gin-gonic/gin"
)

// UserQueryFunc resolves the stored secret for an access key.
type UserQueryFunc func(ctx context.Context, ak string) (string, bool, error)

func MapUserMatch(ud map[string]string) UserQueryFunc {
	return func(ctx context.Context, ak string) (string, bool, error) {
		usk, ok := ud[ak]
		if !ok {
			return "", false, nil
		}
		return usk, true, nil
	}
}

type IAuth interface {
	Name() string
	Auth(ctx *gin.Context, userdata UserQueryFunc) (string, error)
}
