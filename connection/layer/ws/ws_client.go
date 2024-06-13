package ws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/iotool"

	"strings"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"github.com/xxxsen/common/utils"
	"go.uber.org/zap"
)

func init() {
	layer.Register(wsDialLayerName, createWsDialLayer)
}

type wsDialLayer struct {
	c   *cliConfig
	uri *url.URL
}

func createWsDialLayer(params interface{}) (layer.ILayer, error) {
	c := &cliConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	schema := c.Schema
	if len(schema) == 0 {
		schema = "http"
	}
	if !strings.EqualFold(schema, "http") && !strings.EqualFold(schema, "https") {
		return nil, errs.New(errs.ErrParam, "invalid schema:%s", schema)
	}
	uri, err := url.Parse(fmt.Sprintf("%s://%s%s", schema, c.Host, c.Path))
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "invalid uri", err)
	}
	return &wsDialLayer{
		c:   c,
		uri: uri,
	}, nil
}

func (d *wsDialLayer) Name() string {
	return wsDialLayerName
}

func (d *wsDialLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	dialer := ws.Dialer{}
	if d.c.HandshakePaddingMax > 0 {
		header := http.Header{}
		header.Add(randomHeaderKey, utils.RandString(d.c.HandshakePaddingMin, d.c.HandshakePaddingMax))
		dialer.Header = ws.HandshakeHeaderHTTP(header)
	}
	bioReader, _, err := dialer.Upgrade(conn, d.uri)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "upgrade ws fail", err)
	}
	pairIO := conn
	if bioReader != nil {
		pairIO = iotool.WrapConn(conn, bioReader, nil, nil)
	}
	cliio := newWsClientIO(pairIO)
	return iotool.WrapConn(conn, cliio, cliio, nil), nil
}

type wsClientIO struct {
	buf bytes.Buffer
	rw  io.ReadWriter
}

func newWsClientIO(rw io.ReadWriter) *wsClientIO {
	return &wsClientIO{
		rw: rw,
	}
}

func (w *wsClientIO) Read(b []byte) (int, error) {
	if w.buf.Len() == 0 {
		data, err := wsutil.ReadServerBinary(w.rw)
		if err != nil {
			//打个日志看看是否真会出现这种场景
			if len(data) > 0 {
				logutil.GetLogger(context.Background()).Error("client io spare data but get err", zap.Int("data_len", len(data)), zap.Error(err))
			}
			return 0, err
		}
		if len(data) == 0 {
			return 0, errs.New(errs.ErrParam, "read data count == 0")
		}
		w.buf.Write(data)
	}
	return w.buf.Read(b)
}

func (w *wsClientIO) Write(b []byte) (int, error) {
	return len(b), wsutil.WriteClientBinary(w.rw, b)
}
