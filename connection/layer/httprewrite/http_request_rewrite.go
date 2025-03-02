package httprewrite

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/iotool"
	"github.com/xxxsen/common/utils"
)

const (
	httpRequestRewriteDialerName = "http_request_rewrite"
)

func init() {
	layer.Register(httpRequestRewriteDialerName, createHTTPRequestRewriteLayer)
}

func createHTTPRequestRewriteLayer(params interface{}) (layer.ILayer, error) {
	c := &config{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	return &httpRequestRewriteLayer{c: c}, nil
}

type httpRequestRewriteLayer struct {
	c *config
}

func (d *httpRequestRewriteLayer) Name() string {
	return httpRequestRewriteDialerName
}

func (d *httpRequestRewriteLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	bio := bufio.NewReader(conn)
	httpctx, err := ParseBasicHTTPRequestContext(bio)
	if err != nil {
		return nil, fmt.Errorf("read basic http context failed, err:%w", err)
	}
	if len(d.c.RewritePath) != 0 {
		httpctx.URL.Path = d.c.RewritePath
	}
	if len(d.c.RewriteQuery) > 0 {
		q := httpctx.URL.Query()
		for k, v := range d.c.RewriteQuery {
			q.Set(k, v)
		}
		httpctx.URL.RawQuery = q.Encode()
	}
	for k, v := range d.c.RewriteHeader {
		httpctx.Header.Set(k, v)
	}
	if v, ok := d.c.RewriteHeader["host"]; ok {
		httpctx.URL.Host = v
	}
	reader, err := httpctx.ToReader(d.c.ForceUseProxy)
	if err != nil {
		return nil, fmt.Errorf("http to reader failed, err:%w", err)
	}
	reader = io.MultiReader(reader, bio)
	return iotool.WrapConn(conn, reader, nil, nil), nil
}
