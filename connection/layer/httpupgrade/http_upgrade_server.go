package httpupgrade

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"io"
	"net"
	"net/http"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/iotool"

	"github.com/xxxsen/common/utils"

	"github.com/xxxsen/common/errs"
)

func init() {
	layer.Register(httpUpgradeServerLayerName, createHttpUpgradeServerLayer)
}

func createHttpUpgradeServerLayer(params interface{}) (layer.ILayer, error) {
	c := &svrConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	if c.FailCode == 0 {
		c.FailCode = 404
		c.FailReason = "Not Found"
	}
	if len(c.Path) == 0 {
		c.Path = "/"
	}
	return &httpUpgradeServerLayer{c: c}, nil
}

type httpUpgradeServerLayer struct {
	c *svrConfig
}

func (d *httpUpgradeServerLayer) Name() string {
	return httpUpgradeServerLayerName
}

func (d *httpUpgradeServerLayer) writeFailResponse(conn net.Conn) error {
	_, err := io.WriteString(conn, fmt.Sprintf("HTTP/1.1 %d %s\r\n\r\n", d.c.FailCode, d.c.FailReason))
	if err != nil {
		return err
	}
	return nil
}

func (d *httpUpgradeServerLayer) writeUpgradeResponse(conn net.Conn, protocol string) error {
	buf := bytes.NewBuffer(nil)
	_, _ = io.WriteString(buf, fmt.Sprintf("HTTP/1.1 %d %s\r\n", http.StatusSwitchingProtocols, "Switching Protocols"))
	h := http.Header{}
	h.Add("Upgrade", protocol)
	h.Add("Connection", "upgrade")
	if d.c.PaddingMax > 0 {
		h.Set(httpPaddingKey, utils.RandString(int(d.c.PaddingMin), int(d.c.PaddingMax)))
		_ = h.Write(buf)
	}
	_, _ = io.WriteString(buf, "\r\n")
	if _, err := conn.Write(buf.Bytes()); err != nil {
		return err
	}
	return nil
}

func (d *httpUpgradeServerLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	bio := bufio.NewReader(conn)
	req, err := http.ReadRequest(bio)
	if err != nil {
		_ = d.writeFailResponse(conn)
		return nil, errs.Wrap(errs.ErrIO, "invalid http request", err)
	}
	headerProtocol := req.Header.Get("Upgrade")
	//如果本地没有配置upgrade协议, 那么不进行协议校验
	if len(d.c.UpgradeProtocol) > 0 && headerProtocol != d.c.UpgradeProtocol {
		_ = d.writeFailResponse(conn)
		return nil, errs.New(errs.ErrParam, "invalid upgrade protocol:%s", headerProtocol)
	}
	if req.URL.Path != d.c.Path {
		_ = d.writeFailResponse(conn)
		return nil, errs.New(errs.ErrParam, "invalid path:%s", req.URL.Path)
	}
	if err := d.writeUpgradeResponse(conn, headerProtocol); err != nil {
		return nil, errs.Wrap(errs.ErrIO, "write upgrade response fail", err)
	}
	return iotool.WrapConn(conn, bio, nil, nil), nil
}
