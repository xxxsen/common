package codec

import (
	"github.com/xxxsen/common/errs"

	"github.com/gin-gonic/gin"
)

var JsonCodec = &jsonCodec{}

type jsonCodec struct {
}

type jsonMessageFrame struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func (c *jsonCodec) Decode(gctx *gin.Context, request interface{}) error {
	if err := gctx.ShouldBindJSON(request); err != nil {
		return errs.Wrap(errs.ErrUnmarshal, "decode err", err)
	}
	return nil
}

func (c *jsonCodec) Encode(gctx *gin.Context, statuscode int, xerr error, response interface{}) error {
	err := errs.FromError(xerr)
	var (
		code int64 = 0
		msg        = ""
	)
	if err != nil {
		code = err.Code()
		msg = err.Message()
	}

	frame := &jsonMessageFrame{}
	frame.Code = code
	frame.Message = msg
	frame.Data = response
	gctx.JSON(statuscode, frame)
	return nil
}
