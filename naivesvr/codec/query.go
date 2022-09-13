package codec

import (
	"github.com/xxxsen/common/errs"

	"github.com/gin-gonic/gin"
)

var QueryCodec = &queryCodec{}

type queryCodec struct {
	nopCodec
}

func (c *queryCodec) Decode(ctx *gin.Context, request interface{}) error {
	if err := ctx.ShouldBind(request); err != nil {
		return errs.Wrap(errs.ErrParam, "bind query fail", err)
	}
	return nil
}
