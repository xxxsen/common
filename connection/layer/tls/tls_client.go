package tls

import (
	"context"
	ctls "crypto/tls"
	"net"

	"github.com/xxxsen/common/connection/layer"

	utls "github.com/refraction-networking/utls"

	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/utils"
)

func init() {
	layer.Register(tlsDialLayerName, createTLSDialLayer)
}

type itlsconn interface {
	net.Conn
	VerifyHostname(host string) error
	HandshakeContext(ctx context.Context) error
}

func createTLSDialLayer(params interface{}) (layer.ILayer, error) {
	c := &cliConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	if len(c.FingerPrint) > 0 {
		if _, ok := mappingHelloID[c.FingerPrint]; !ok {
			return nil, errs.New(errs.ErrParam, "invalid finger print str:%s", c.FingerPrint)
		}
	}
	return &tlsDialLayer{c: c}, nil
}

type tlsDialLayer struct {
	c *cliConfig
}

func (d *tlsDialLayer) Name() string {
	return tlsDialLayerName
}

func (d *tlsDialLayer) createTlsConn(conn net.Conn, fingerprint string) itlsconn {
	if len(d.c.FingerPrint) == 0 {
		tlsconn := ctls.Client(conn, &ctls.Config{
			ServerName:         d.c.SNI,
			InsecureSkipVerify: d.c.SkipInsecureVerify,
			MinVersion:         GetVersionIdByNameOrDefault(d.c.MinTLSVersion),
			MaxVersion:         GetVersionIdByNameOrDefault(d.c.MaxTLSVersion),
		})
		return tlsconn
	}
	return utls.UClient(conn, &utls.Config{
		ServerName:         d.c.SNI,
		InsecureSkipVerify: d.c.SkipInsecureVerify,
	}, mappingHelloID[fingerprint])
}

func (d *tlsDialLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	tlsconn := d.createTlsConn(conn, d.c.FingerPrint)
	if err := tlsconn.HandshakeContext(ctx); err != nil {
		return nil, errs.Wrap(errs.ErrIO, "handshake tls fail", err)
	}
	if !d.c.SkipInsecureVerify {
		if err := tlsconn.VerifyHostname(d.c.SNI); err != nil {
			return nil, errs.Wrap(errs.ErrUnknown, "verify tls host name fail", err)
		}
	}
	return tlsconn, nil
}
