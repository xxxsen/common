package middleware

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"net/http"

	"github.com/xxxsen/common/webapi/proxyutil"

	"github.com/gin-gonic/gin"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

const (
	defaultMaxAllowChunkStreamLength = 5 * 1024 * 1024 // 5 MB
)

var (
	errChunkedEncodingRequired = errors.New("unknown-length body must use chunked encoding")
	errChunkStreamTooLarge     = errors.New("chunk stream exceeds the length limit")
)

func NonLengthIOLimitMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		if c.Request.ContentLength >= 0 {
			c.Next()
			return
		}
		ctx := c.Request.Context()
		logutil.GetLogger(ctx).Debug("recv non-content-length io request")
		if len(c.Request.TransferEncoding) == 0 || c.Request.TransferEncoding[0] != "chunked" {
			proxyutil.FailStatus(c, http.StatusBadRequest, errChunkedEncodingRequired)
			return
		}
		data, err := io.ReadAll(io.LimitReader(c.Request.Body, defaultMaxAllowChunkStreamLength+1))
		if err != nil {
			proxyutil.FailStatus(c, http.StatusBadRequest, fmt.Errorf("read client data failed, err:%w", err))
			return
		}
		logutil.GetLogger(ctx).Debug("read chunk stream from client", zap.Int("length", len(data)))
		if len(data) > defaultMaxAllowChunkStreamLength {
			proxyutil.FailStatus(c, http.StatusBadRequest, errChunkStreamTooLarge)
			return
		}
		r := bytes.NewReader(data)
		clr := c.Request.Body
		rc := &readCloserWrap{
			r: r,
			c: clr,
		}
		c.Request.Body = rc
		c.Request.ContentLength = int64(len(data))
	}
}

type readCloserWrap struct {
	r io.Reader
	c io.Closer
}

func (c *readCloserWrap) Read(p []byte) (int, error) {
	count, err := c.r.Read(p)
	if err == nil {
		return count, nil
	}
	if errors.Is(err, io.EOF) {
		return count, io.EOF
	}
	return count, fmt.Errorf("read buffered request: %w", err)
}

func (c *readCloserWrap) Close() error {
	if err := c.c.Close(); err != nil {
		return fmt.Errorf("close request body: %w", err)
	}
	return nil
}
