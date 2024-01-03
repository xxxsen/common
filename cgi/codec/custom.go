package codec

import (
	"github.com/gin-gonic/gin"
)

type customCodec struct {
	enc ICodec
	dec ICodec
}

func CustomCodec(enc, dec ICodec) *customCodec {
	return &customCodec{
		enc: enc,
		dec: dec,
	}
}

func (c *customCodec) Decode(ctx *gin.Context, request interface{}) error {
	return c.dec.Decode(ctx, request)
}

func (c *customCodec) Encode(ctx *gin.Context, statuscode int, err error, response interface{}) error {
	return c.enc.Encode(ctx, statuscode, err, response)
}
