package codec

import (
	"github.com/xxxsen/common/errs"

	"github.com/gin-gonic/gin"
)

var MultipartCodec = &multipartCodec{}

type multipartCodec struct {
	nopCodec
}

func (c *multipartCodec) Decode(ctx *gin.Context, request interface{}) error {
	if err := ctx.ShouldBind(request); err != nil {
		return errs.Wrap(errs.ErrParam, "decode multipart fail", err)
	}
	return nil
}
