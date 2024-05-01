package ws

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"
	"net/url"
	"strings"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/iotool"

	"github.com/xxxsen/common/utils"

	"github.com/gobwas/ws"
	"github.com/gobwas/ws/wsutil"
	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/logutil"
	"go.uber.org/zap"
)

func init() {
	layer.Regist(wsServerLayerName, createWsServerLayer)
}

func createWsServerLayer(params interface{}) (layer.ILayer, error) {
	c := &svrConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	return &wsServerLayer{c: c}, nil
}

type wsServerLayer struct {
	c *svrConfig
}

func (d *wsServerLayer) Name() string {
	return wsServerLayerName
}

func (d *wsServerLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	var errDetail error
	ug := ws.Upgrader{
		OnRequest: func(uri []byte) error {
			if len(d.c.Path) == 0 {
				return nil
			}
			parsedUri, err := url.ParseRequestURI(string(uri))
			if err != nil {
				errDetail = fmt.Errorf("parse uri fail, uri:%s", string(uri))
				return ws.RejectConnectionError(ws.RejectionStatus(http.StatusNotFound))
			}
			if !strings.HasPrefix(parsedUri.Path, d.c.Path) {
				errDetail = fmt.Errorf("path not math, path:%s", parsedUri.Path)
				return ws.RejectConnectionError(ws.RejectionStatus(http.StatusNotFound))
			}
			return nil
		},
		OnHost: func(host []byte) error {
			if len(d.c.Host) == 0 {
				return nil
			}
			if !strings.EqualFold(d.c.Host, string(host)) {
				errDetail = fmt.Errorf("host not match, host:%s", host)
				return ws.RejectConnectionError(ws.RejectionStatus(http.StatusNotFound))
			}
			return nil
		},
	}
	if d.c.HandshakePaddingMax > 0 {
		header := http.Header{}
		header.Add(randomHeaderKey, utils.RandString(d.c.HandshakePaddingMin, d.c.HandshakePaddingMax))
		ug.Header = ws.HandshakeHeaderHTTP(header)
	}
	_, err := ug.Upgrade(conn)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, fmt.Sprintf("upgrade conn fail, err detail:[%+v]", errDetail), err)
	}

	svrio := newWsServerIO(conn)
	return iotool.WrapConn(conn, svrio, svrio, nil), nil
}

type wsServerIO struct {
	buf bytes.Buffer
	rw  io.ReadWriter
}

func newWsServerIO(rw io.ReadWriter) *wsServerIO {
	return &wsServerIO{
		rw: rw,
	}
}

func (w *wsServerIO) Read(b []byte) (int, error) {
	if w.buf.Len() == 0 {
		data, err := wsutil.ReadClientBinary(w.rw)
		if err != nil {
			//打个日志看看是否真会出现这种场景
			if len(data) > 0 {
				logutil.GetLogger(context.Background()).Error("server io spare data but get err", zap.Int("data_len", len(data)), zap.Error(err))
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

func (w *wsServerIO) Write(b []byte) (int, error) {
	return len(b), wsutil.WriteServerBinary(w.rw, b)
}
