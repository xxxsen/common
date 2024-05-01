package tls

import (
	"context"
	"crypto/tls"
	"net"

	"github.com/xxxsen/common/connection/layer"

	"github.com/xxxsen/common/errs"
	"github.com/xxxsen/common/utils"
)

func init() {
	layer.Regist(tlsServerLayerName, createTlsServerLayer)
}

func createTlsServerLayer(params interface{}) (layer.ILayer, error) {
	c := &svrConfig{}
	if err := utils.ConvStructJson(params, c); err != nil {
		return nil, err
	}
	cer, err := tls.LoadX509KeyPair(c.CertFile, c.KeyFile)
	if err != nil {
		return nil, errs.Wrap(errs.ErrParam, "invalid cert info", err)
	}
	return &tlsServerLayer{
		c: c,
		tlsc: &tls.Config{
			Certificates: []tls.Certificate{cer},
			MinVersion:   GetVersionIdByNameOrDefault(c.MinTLSVersion),
			MaxVersion:   GetVersionIdByNameOrDefault(c.MaxTLSVersion),
		},
	}, nil
}

type tlsServerLayer struct {
	c    *svrConfig
	tlsc *tls.Config
}

func (d *tlsServerLayer) Name() string {
	return tlsServerLayerName
}

func (d *tlsServerLayer) MakeLayerContext(ctx context.Context, conn net.Conn) (net.Conn, error) {
	tlsconn := tls.Server(conn, d.tlsc)
	if err := tlsconn.HandshakeContext(ctx); err != nil {
		return nil, errs.Wrap(errs.ErrIO, "handshake tls fail", err)
	}
	return tlsconn, nil
}
