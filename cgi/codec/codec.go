package codec

import (
	"github.com/gin-gonic/gin"
)

type ICodec interface {
	Decode(c *gin.Context, request interface{}) error
	Encode(c *gin.Context, statuscode int, err error, response interface{}) error
}
