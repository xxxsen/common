package codec

import (
	"github.com/gin-gonic/gin"
)

var NopCodec = &nopCodec{}

type nopCodec struct {
}

func (c *nopCodec) Decode(ctx *gin.Context, request interface{}) error {
	return nil
}

func (c *nopCodec) Encode(ctx *gin.Context, statuscode int, err error, response interface{}) error {
	return nil
}
