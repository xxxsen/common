package naivesvr

import (
	"context"
	"fmt"

	"github.com/xxxsen/common/naivesvr/constants"
)

var emptyAttach = map[string]interface{}{}

func GetAttach(ctx context.Context) map[string]interface{} {
	iVal := ctx.Value(constants.KeyServerAttach)
	if iVal == nil {
		return emptyAttach
	}
	return iVal.(map[string]interface{})
}

func GetServer(ctx context.Context) (*Server, bool) {
	v := ctx.Value(constants.KeyServer)
	if v == nil {
		return nil, false
	}
	return v.(*Server), true
}

func MustGetServer(ctx context.Context) *Server {
	if svr, ok := GetServer(ctx); ok {
		return svr
	}
	panic(fmt.Errorf("svr not found in ctx"))
}
