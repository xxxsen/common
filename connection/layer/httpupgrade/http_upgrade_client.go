package httpupgrade

import (
	"bufio"
	"context"
	"fmt"
	"net"
	"net/http"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/iotool"

	"github.com/xxxsen/common/utils"

	"github.com/xxxsen/common/errs"
)

func init() {
	layer.Register(httpUpgradeClientLayerName, createHttpUpgradeClientLayer)
}

func createHttpUpgradeClientLayer(params interface{}) (layer.ILayer, error) {
	c := &cliConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	if len(c.Path) == 0 {
		c.Path = "/"
	}
	if len(c.Host) == 0 {
		return nil, errs.New(errs.ErrParam, "invalid host name")
	}
	return &httpUpgradeClientLayer{
		c: c,
	}, nil
}

type httpUpgradeClientLayer struct {
	c *cliConfig
}

func (c *httpUpgradeClientLayer) Name() string {
	return httpUpgradeClientLayerName
}

func (d *httpUpgradeClientLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	req, err := http.NewRequest(http.MethodGet, fmt.Sprintf("http://%s%s", d.c.Host, d.c.Path), nil)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "create http upgrade request fail", err)
	}
	if d.c.PaddingMax > 0 {
		req.Header.Set(httpPaddingKey, utils.RandString(int(d.c.PaddingMin), int(d.c.PaddingMax)))
	}
	protocol := defaultHTTPUpgradeProtocol
	if len(d.c.UpgradeProtocol) != 0 {
		protocol = d.c.UpgradeProtocol
	}
	req.Header.Add("Upgrade", protocol)
	req.Header.Add("Connection", "upgrade")
	if err := req.Write(conn); err != nil {
		return nil, errs.Wrap(errs.ErrIO, "write http request fail", err)
	}
	bior := bufio.NewReader(conn)
	rsp, err := http.ReadResponse(bior, req)
	if err != nil {
		return nil, errs.Wrap(errs.ErrIO, "read body fail", err)
	}
	if rsp.StatusCode != http.StatusSwitchingProtocols {
		rsp.Body.Close()
		return nil, errs.Wrap(errs.ErrUnknown, "unswitchable response", err)
	}

	return iotool.WrapConn(conn, bior, nil, nil), nil
}
